package api

import (
	"errors"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/mundipagg/boleto-api/bank"
	"github.com/mundipagg/boleto-api/config"
	"github.com/mundipagg/boleto-api/db"
	"github.com/mundipagg/boleto-api/log"
	"github.com/mundipagg/boleto-api/metrics"
	"github.com/mundipagg/boleto-api/models"
	"github.com/mundipagg/boleto-api/util"
)

// ReturnHeaders 'seta' os headers padrões de resposta
func ReturnHeaders() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Header("Content-Type", "application/json")
		c.Next()
	}
}

func executionController() gin.HandlerFunc {
	return func(c *gin.Context) {
		if config.IsRunning() {
			c.AbortWithError(500, errors.New("A aplicação está sendo finalizada"))
			return
		}
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

		user, _, _ := c.Request.BasicAuth()

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
		l.NossoNumero = boleto.Title.OurNumber
		l.Operation = "RegisterBoleto"
		l.Recipient = boleto.Recipient.Name
		l.RequestKey = boleto.RequestKey
		l.BankName = bank.GetBankNameIntegration()
		l.IPAddress = c.ClientIP()
		l.ServiceRefererName = user
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
		log := log.CreateLog()

		cred := getHeaderCredentials(c)

		if cred == nil {
			c.AbortWithStatusJSON(401, models.GetBoletoResponseError("MP401", "Unauthorized"))
			return
		}

		mongo, errMongo := db.CreateMongo(log)
		if checkError(c, errMongo, log) {
			c.AbortWithStatusJSON(500, models.GetBoletoResponseError("MP500", "InternalError"))
			return
		}

		if !mongo.HasValidCredentials(cred) {
			c.AbortWithStatusJSON(401, models.GetBoletoResponseError("MP401", "Unauthorized"))
			return
		}
		c.Next()
	}
}

func getHeaderCredentials(c *gin.Context) *models.Credentials {
	user, pass, hasAuth := c.Request.BasicAuth()
	if user == "" || pass == "" || !hasAuth {
		return nil
	}
	return models.NewCredentials(user, pass)
}
