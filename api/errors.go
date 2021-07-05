package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/mundipagg/boleto-api/models"
)

var validate = map[string]int{
	"MP400":                   http.StatusBadRequest,
	"MPAmountInCents":         http.StatusBadRequest,
	"MPExpireDate":            http.StatusBadRequest,
	"MPBuyerDocumentType":     http.StatusBadRequest,
	"MPDocumentNumber":        http.StatusBadRequest,
	"MPRecipientDocumentType": http.StatusBadRequest,
	"MPTimeout":               http.StatusGatewayTimeout,
	"MPOurNumberFail":         http.StatusBadGateway,
}

//Esse mapper poderá ser movido para interface IBank posteriomente
var stone = map[string]int{
	"srn:error:validation":          http.StatusBadRequest,
	"srn:error:unauthenticated":     http.StatusInternalServerError,
	"srn:error:unauthorized":        http.StatusBadGateway,
	"srn:error:not_found":           http.StatusBadGateway,
	"srn:error:conflict":            http.StatusBadGateway,
	"srn:error:product_not_enabled": http.StatusBadRequest,
}

func handleErrors(c *gin.Context) {
	c.Next()

	var status int
	var exist bool

	response := getResponseFromContext(c)

	if !qualifiedForNewErrorHandling(c, response) {
		return
	}

	bank := getBankFromContext(c).GetBankNumber()
	bankcode := response.Errors[0].Code

	if status, exist = validate[bankcode]; !exist {
		status, exist = getMapper(bank)[bankcode]
	}

	switch status {
	case http.StatusBadRequest:
		response.Errors[0].Code = "MP400"
		c.JSON(http.StatusBadRequest, response)
	case http.StatusBadGateway:
		response.Errors[0].Code = "MP502"
		c.JSON(http.StatusBadGateway, response)
	case http.StatusGatewayTimeout:
		response.Errors[0].Code = "MP504"
		c.JSON(http.StatusGatewayTimeout, response)
	default:
		response.Errors[0].Code = "MP500"
		c.JSON(http.StatusInternalServerError, response)
	}

	c.Set(responseKey, response)
}

func getMapper(bank models.BankNumber) map[string]int {
	switch bank {
	default:
		return stone
	}
}

func qualifiedForNewErrorHandling(c *gin.Context, response models.BoletoResponse) bool {
	bankNumber := getBankFromContext(c).GetBankNumber()
	if bankNumber == models.Stone && response.HasErrors() {
		return true
	}
	return false
}

func getErrorCodeToLog(c *gin.Context) string {
	response := getResponseFromContext(c)
	if response.HasErrors() {
		return response.Errors[0].ErrorCode()
	}
	return ""
}
