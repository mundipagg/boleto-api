package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/mundipagg/boleto-api/bank"
	"github.com/mundipagg/boleto-api/log"
	"github.com/mundipagg/boleto-api/metrics"
	"github.com/mundipagg/boleto-api/models"
	"github.com/mundipagg/boleto-api/util"
)

//Logger Middleware de log do request e response da BoletoAPI
func Logger(c *gin.Context) {
	boleto := getBoletoFromContext(c)
	bank := getBankFromContext(c)

	l := loadLog(c, boleto, bank)

	l.RequestApplication(boleto, c.Request.URL.RequestURI(), util.HeaderToMap(c.Request.Header))

	c.Next()

	resp, _ := c.Get("boletoResponse")
	l.ResponseApplication(resp, c.Request.URL.RequestURI())

	tag := bank.GetBankNameIntegration() + "-status"
	metrics.PushBusinessMetric(tag, c.Writer.Status())
}

func loadLog(c *gin.Context, boleto models.BoletoRequest, bank bank.Bank) *log.Log {
	l := log.CreateLog()
	user, _ := c.Get(serviceUserKey)
	l.NossoNumero = boleto.Title.OurNumber
	l.Operation = "RegisterBoleto"
	l.Recipient = boleto.Recipient.Name
	l.RequestKey = boleto.RequestKey
	l.BankName = bank.GetBankNameIntegration()
	l.IPAddress = c.ClientIP()
	l.ServiceUser = user.(string)
	return l
}
