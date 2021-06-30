package api

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/mundipagg/boleto-api/bank"
	"github.com/mundipagg/boleto-api/log"
	"github.com/mundipagg/boleto-api/metrics"
	"github.com/mundipagg/boleto-api/models"
	"github.com/mundipagg/boleto-api/usermanagement"
	"github.com/mundipagg/boleto-api/util"
)

// ReturnHeaders 'seta' os headers padrões de resposta
func ReturnHeaders() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Header("Content-Type", "application/json")
		c.Next()
	}
}

func timingMetrics() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		c.Next()
		end := time.Now()
		total := end.Sub(start)
		s := float64(total.Seconds())
		metrics.PushTimingMetric("request-time", s)
	}
}

//ParseBoleto trata a entrada de boleto em todos os requests
func ParseBoleto() gin.HandlerFunc {
	return func(c *gin.Context) {

		boleto := models.BoletoRequest{}
		errBind := c.BindJSON(&boleto)
		if errBind != nil {
			e := models.NewFormatError(errBind.Error())
			checkError(c, e, log.CreateLog())
			metrics.PushBusinessMetric("json_error", 1)
			return
		}
		bank, err := bank.Get(boleto)
		if checkError(c, err, log.CreateLog()) {
			c.Set("error", err)
			return
		}
		c.Set("bank", bank)
		d, errFmt := time.Parse("2006-01-02", boleto.Title.ExpireDate)
		boleto.Title.ExpireDateTime = d
		if errFmt != nil {
			e := models.NewFormatError(errFmt.Error())
			checkError(c, e, log.CreateLog())
			metrics.PushBusinessMetric(bank.GetBankNameIntegration()+"-bad-request", 1)
			return
		}

		l := log.CreateLog()
		user, _ := c.Get("serviceuser")
		l.NossoNumero = boleto.Title.OurNumber
		l.Operation = "RegisterBoleto"
		l.Recipient = boleto.Recipient.Name
		l.RequestKey = boleto.RequestKey
		l.BankName = bank.GetBankNameIntegration()
		l.IPAddress = c.ClientIP()
		l.ServiceUser = user.(string)
		l.RequestApplication(boleto, c.Request.URL.RequestURI(), util.HeaderToMap(c.Request.Header))
		c.Set("boleto", boleto)
		c.Next()
		resp, _ := c.Get("boletoResponse")
		l.ResponseApplication(resp, c.Request.URL.RequestURI())
		tag := bank.GetBankNameIntegration() + "-status"
		metrics.PushBusinessMetric(tag, c.Writer.Status())
	}
}

//Authentication Trata a autenticação para registro de boleto
func Authentication() gin.HandlerFunc {
	return func(c *gin.Context) {

		cred := getHeaderCredentials(c)

		if cred == nil {
			c.AbortWithStatusJSON(401, models.GetBoletoResponseError("MP401", "Unauthorized"))
			return
		}

		if !hasValidCredentials(cred) {
			c.AbortWithStatusJSON(401, models.GetBoletoResponseError("MP401", "Unauthorized"))
			return
		}

		c.Set("serviceuser", cred.Username)

		c.Next()
	}
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
