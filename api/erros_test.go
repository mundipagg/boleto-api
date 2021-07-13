package api

import (
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/mundipagg/boleto-api/bank"
	"github.com/mundipagg/boleto-api/models"
	"github.com/stretchr/testify/assert"
)

func Test_QualifiedForNewErrorHandling_WhenBankStoneWithError_ReturnTrue(t *testing.T) {
	request := models.BoletoRequest{BankNumber: models.Stone}
	response := models.GetBoletoResponseError("MP000", "error")
	bank, _ := bank.Get(request)

	c, _ := gin.CreateTestContext(httptest.NewRecorder())
	c.Set(boletoKey, request)
	c.Set(bankKey, bank)

	result := qualifiedForNewErrorHandling(c, response)

	assert.True(t, result)
}

func Test_QualifiedForNewErrorHandling_WhenBankStoneWithoutError_ReturnFalse(t *testing.T) {
	request := models.BoletoRequest{BankNumber: models.Stone}
	response := models.BoletoResponse{}
	bank, _ := bank.Get(request)

	c, _ := gin.CreateTestContext(httptest.NewRecorder())
	c.Set(boletoKey, request)
	c.Set(bankKey, bank)

	result := qualifiedForNewErrorHandling(c, response)

	assert.False(t, result)
}

func Test_QualifiedForNewErrorHandling_WhenAnotherBankWithError_ReturnFalse(t *testing.T) {
	request := models.BoletoRequest{BankNumber: models.Caixa}
	response := models.GetBoletoResponseError("MP000", "error")
	bank, _ := bank.Get(request)

	c, _ := gin.CreateTestContext(httptest.NewRecorder())
	c.Set(boletoKey, request)
	c.Set(bankKey, bank)

	result := qualifiedForNewErrorHandling(c, response)

	assert.False(t, result)
}

func Test_QualifiedForNewErrorHandling_WhenAnotherBankWithoutError_ReturnFalse(t *testing.T) {
	request := models.BoletoRequest{BankNumber: models.Caixa}
	response := models.BoletoResponse{}
	bank, _ := bank.Get(request)

	c, _ := gin.CreateTestContext(httptest.NewRecorder())
	c.Set(boletoKey, request)
	c.Set(bankKey, bank)

	result := qualifiedForNewErrorHandling(c, response)

	assert.False(t, result)
}
