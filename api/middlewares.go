package api

import (
	"errors"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/mundipagg/boleto-api/bank"
	"github.com/mundipagg/boleto-api/config"
	"github.com/mundipagg/boleto-api/log"
	"github.com/mundipagg/boleto-api/metrics"
	"github.com/mundipagg/boleto-api/models"
	"github.com/mundipagg/boleto-api/usermanagement"
	"github.com/mundipagg/boleto-api/util"
)

const (
	boletoKey      = "boleto"
	bankKey        = "bank"
	serviceUserKey = "serviceuser"
	responseKey    = "boletoResponse"
)

func returnHeaders() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Header("Content-Type", "application/json")
		c.Next()
	}
}

func executionController() gin.HandlerFunc {
	return func(c *gin.Context) {
		if config.IsRunning() {
			c.AbortWithError(500, errors.New("a aplicação está sendo finalizada"))
			return
		}
	}
}

//parseBoleto Middleware de tratamento do request de registro de boleto
func parseBoleto(c *gin.Context) {
	var ok bool
	var boleto models.BoletoRequest
	var bank bank.Bank

	if boleto, ok = getBoletoRequest(c); !ok {
		return
	}

	if bank, ok = getBank(c, boleto); !ok {
		return
	}

	if !parseExpirationDate(c, &boleto, bank) {
		return
	}

	c.Set(boletoKey, boleto)
	c.Set(bankKey, bank)
}

//authentication Middleware de autenticação para registro de boleto
func authentication(c *gin.Context) {

	cred := getHeaderCredentials(c)

	if cred == nil {
		c.AbortWithStatusJSON(401, models.GetBoletoResponseError("MP401", "Unauthorized"))
		return
	}

	if !hasValidCredentials(cred) {
		c.AbortWithStatusJSON(401, models.GetBoletoResponseError("MP401", "Unauthorized"))
		return
	}

	c.Set(serviceUserKey, cred.Username)
}

//logger Middleware de log do request e response da BoletoAPI
func logger(c *gin.Context) {
	boleto := getBoletoFromContext(c)
	bank := getBankFromContext(c)

	l := loadBankLog(c)

	l.RequestApplication(boleto, c.Request.URL.RequestURI(), util.HeaderToMap(c.Request.Header))

	c.Next()

	l = loadBankLog(c)

	resp, _ := c.Get(responseKey)

	l.ResponseApplication(resp, c.Request.URL.RequestURI(), getErrorCodeToLog(c))

	tag := bank.GetBankNameIntegration() + "-status"
	metrics.PushBusinessMetric(tag, c.Writer.Status())
}

//checkError Middleware de verificação de erros
func checkError(c *gin.Context, err error, l *log.Log) bool {

	if err != nil {
		errResp := models.BoletoResponse{
			Errors: models.NewErrors(),
		}

		switch v := err.(type) {

		case models.ErrorResponse:
			errResp.Errors.Append(v.ErrorCode(), v.Error())
			c.JSON(http.StatusBadRequest, errResp)

		case models.HttpNotFound:
			errResp.Errors.Append("MP404", v.Error())
			l.Warn(errResp, v.Error())
			c.JSON(http.StatusNotFound, errResp)

		case models.InternalServerError:
			errResp.Errors.Append("MP500", v.Error())
			l.Warn(errResp, v.Error())
			c.JSON(http.StatusInternalServerError, errResp)

		case models.BadGatewayError:
			errResp.Errors.Append("MP502", v.Error())
			l.Warn(errResp, v.Error())
			c.JSON(http.StatusBadGateway, errResp)

		case models.FormatError:
			errResp.Errors.Append("MP400", v.Error())
			l.Warn(errResp, v.Error())
			c.JSON(http.StatusBadRequest, errResp)

		default:
			errResp.Errors.Append("MP500", "Internal Error")
			l.Fatal(errResp, v.Error())
			c.JSON(http.StatusInternalServerError, errResp)
		}

		c.Set("boletoResponse", errResp)
		return true
	}
	return false
}

func hasValidCredentials(c *models.Credentials) bool {
	u, hasUser := usermanagement.GetUser(c.UserKey)

	if !hasUser {
		return false
	}

	user := u.(models.Credentials)

	if user.UserKey == c.UserKey && user.Password == c.Password {
		c.Username = user.Username
		return true
	}

	return false
}

func getHeaderCredentials(c *gin.Context) *models.Credentials {
	userkey, pass, hasAuth := c.Request.BasicAuth()
	if userkey == "" || pass == "" || !hasAuth {
		return nil
	}
	return models.NewCredentials(userkey, pass)
}

func getBoletoRequest(c *gin.Context) (models.BoletoRequest, bool) {
	boleto := models.BoletoRequest{}
	errBind := c.BindJSON(&boleto)
	if errBind != nil {
		e := models.NewFormatError(errBind.Error())
		checkError(c, e, log.CreateLog())
		return boleto, false
	}
	return boleto, true
}

func getBank(c *gin.Context, boleto models.BoletoRequest) (bank.Bank, bool) {
	bank, err := bank.Get(boleto)
	if checkError(c, err, log.CreateLog()) {
		c.Set("error", err)
		return bank, false
	}
	return bank, true
}

func parseExpirationDate(c *gin.Context, boleto *models.BoletoRequest, bank bank.Bank) bool {
	d, errFmt := time.Parse("2006-01-02", boleto.Title.ExpireDate)
	boleto.Title.ExpireDateTime = d
	if errFmt != nil {
		e := models.NewFormatError(errFmt.Error())
		checkError(c, e, log.CreateLog())
		metrics.PushBusinessMetric(bank.GetBankNameIntegration()+"-bad-request", 1)
		return false
	}
	return true
}

func loadBankLog(c *gin.Context) *log.Log {
	boleto := getBoletoFromContext(c)
	bank := getBankFromContext(c)
	l := bank.Log()
	l.Operation = "RegisterBoleto"
	l.NossoNumero = getNossoNumeroFromContext(c)
	l.Recipient = boleto.Recipient.Name
	l.RequestKey = boleto.RequestKey
	l.BankName = bank.GetBankNameIntegration()
	l.IPAddress = c.ClientIP()
	l.ServiceUser = getUserFromContext(c)
	return l
}

func getBoletoFromContext(c *gin.Context) models.BoletoRequest {
	if boleto, exists := c.Get(boletoKey); exists {
		return boleto.(models.BoletoRequest)
	}
	return models.BoletoRequest{}
}

func getBankFromContext(c *gin.Context) bank.Bank {
	if banking, exists := c.Get(bankKey); exists {
		return banking.(bank.Bank)
	}
	return nil
}

func getUserFromContext(c *gin.Context) string {
	if user, exists := c.Get(serviceUserKey); exists {
		return user.(string)
	}
	return ""
}

func getResponseFromContext(c *gin.Context) models.BoletoResponse {
	if response, exists := c.Get(responseKey); exists {
		return response.(models.BoletoResponse)
	}
	return models.BoletoResponse{}
}

func getNossoNumeroFromContext(c *gin.Context) string {
	if respOurNumber := getResponseFromContext(c).OurNumber; respOurNumber != "" {
		return respOurNumber
	}

	return strconv.FormatUint(uint64(getBoletoFromContext(c).Title.OurNumber), 10)
}
