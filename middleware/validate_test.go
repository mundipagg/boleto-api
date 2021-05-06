package middleware

import (
	"bytes"
	"net/http"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/mundipagg/boleto-api/models"
	"github.com/mundipagg/boleto-api/test"
	"github.com/mundipagg/boleto-api/util"
	"github.com/stretchr/testify/assert"
)

func Test_ValidateRegisterV1_WhenWithoutRules_PassSuccessful(t *testing.T) {
	router, w := arrangeMiddlewareRoute("/validateV1", ParseBoleto, ValidateRegisterV1)
	body := test.NewStubBoletoRequest(models.BancoDoBrasil).WithExpirationDate(time.Now()).Build()
	req, _ := http.NewRequest("POST", "/validateV1", bytes.NewBuffer([]byte(util.ToJSON(body))))

	router.ServeHTTP(w, req)

	assert.Equal(t, 200, w.Code)
}

func Test_ValidateRegisterV1_WhenHasRules_ReturnBadRequest(t *testing.T) {
	router, w := arrangeMiddlewareRoute("/validateV1", ParseBoleto, ValidateRegisterV1)
	body := test.NewStubBoletoRequest(models.Caixa).WithExpirationDate(time.Now()).WithAcceptDivergentAmount(true).Build()
	req, _ := http.NewRequest("POST", "/validateV1", bytes.NewBuffer([]byte(util.ToJSON(body))))

	router.ServeHTTP(w, req)

	assert.Equal(t, 400, w.Code)
	assert.Equal(t, `[{"code":"MP400","message":"title.rules not available in this version"}]`, w.Body.String())
}

func Test_ValidateRegisterV2_WhenWithoutRules_PassSuccessful(t *testing.T) {
	router, w := arrangeMiddlewareRoute("/validateV2", ParseBoleto, ValidateRegisterV2)
	body := test.NewStubBoletoRequest(models.BancoDoBrasil).WithExpirationDate(time.Now()).Build()
	req, _ := http.NewRequest("POST", "/validateV2", bytes.NewBuffer([]byte(util.ToJSON(body))))

	router.ServeHTTP(w, req)

	assert.Equal(t, 200, w.Code)
}

func Test_ValidateRegisterV2_WhenHasRulesAndNotCaixaBank_ReturnBadRequest(t *testing.T) {
	router, w := arrangeMiddlewareRoute("/validateV2", ParseBoleto, ValidateRegisterV2)
	body := test.NewStubBoletoRequest(models.BancoDoBrasil).WithExpirationDate(time.Now()).WithAcceptDivergentAmount(true).Build()
	req, _ := http.NewRequest("POST", "/validateV2", bytes.NewBuffer([]byte(util.ToJSON(body))))

	router.ServeHTTP(w, req)

	assert.Equal(t, 400, w.Code)
	assert.Equal(t, `[{"code":"MP400","message":"title.rules not available for this bank"}]`, w.Body.String())
}

func Test_ValidateRegisterV2_WhenHasRulesAndCaixaBank_PassSuccessful(t *testing.T) {
	router, w := arrangeMiddlewareRoute("/validateV2", ParseBoleto, ValidateRegisterV2)
	body := test.NewStubBoletoRequest(models.Caixa).WithExpirationDate(time.Now()).WithAcceptDivergentAmount(true).Build()
	req, _ := http.NewRequest("POST", "/validateV2", bytes.NewBuffer([]byte(util.ToJSON(body))))

	router.ServeHTTP(w, req)

	assert.Equal(t, 200, w.Code)
}

func Test_GetBoletoFromContext_WhenNotSetBoletoRequestInContext_ReturnEmptyRequest(t *testing.T) {
	c := new(gin.Context)

	result := getBoletoFromContext(c)

	assert.Equal(t, models.BoletoRequest{}, result)
}

func Test_GetBankFromContext_WhenNotSetBankInContext_ReturnNil(t *testing.T) {
	c := new(gin.Context)

	result := getBankFromContext(c)

	assert.Equal(t, nil, result)
}
