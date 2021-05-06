package middleware

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/mundipagg/boleto-api/bank"
	"github.com/mundipagg/boleto-api/log"
	"github.com/mundipagg/boleto-api/models"
)

const (
	boletoKey = "boleto"
	bankKey   = "bank"
)

//ParseBoleto Middleware de tratamento do request de registro de boleto
func ParseBoleto(c *gin.Context) {
	var boleto models.BoletoRequest
	var bank bank.Bank
	var err error

	if boleto, err = getBoletoRequest(c); err != nil {
		response := errorResponse("MP400", err)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	if bank, err = getBank(c, boleto); err != nil {
		response := errorResponse("MPBankNumber", err)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	if err = parseExpirationDate(c, boleto, bank); err != nil {
		response := errorResponse("MP400", err)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	c.Set(boletoKey, boleto)
	c.Set(bankKey, bank)
}

func getBoletoRequest(c *gin.Context) (models.BoletoRequest, error) {
	boleto := models.BoletoRequest{}
	errBind := c.BindJSON(&boleto)
	if errBind != nil {
		return models.BoletoRequest{}, errBind
	}
	return boleto, nil
}

func getBank(c *gin.Context, boleto models.BoletoRequest) (bank.Bank, error) {
	bank, err := bank.Get(boleto)
	if err != nil {
		c.Set("error", err)
		return nil, err
	}
	return bank, nil
}

func parseExpirationDate(c *gin.Context, boleto models.BoletoRequest, bank bank.Bank) error {
	d, errFmt := time.Parse("2006-01-02", boleto.Title.ExpireDate)
	boleto.Title.ExpireDateTime = d
	return errFmt
}

func errorResponse(code string, err error) models.BoletoResponse {
	response := models.BoletoResponse{
		Errors: models.NewErrors(),
	}
	l := log.CreateLog()
	l.Warn(err, err.Error())
	response.Errors.Append(code, err.Error())
	return response
}
