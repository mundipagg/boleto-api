package api

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/mundipagg/boleto-api/bank"
	"github.com/mundipagg/boleto-api/models"
	"github.com/mundipagg/boleto-api/test"
	"github.com/mundipagg/boleto-api/usermanagement"
	"github.com/mundipagg/boleto-api/util"
	"github.com/stretchr/testify/assert"
)

func Test_Authentication_WhenEmptyCredentials_ReturnUnauthorized(t *testing.T) {
	router, w := arrangeMiddlewareRoute("/authentication", authentication)
	req, _ := http.NewRequest("POST", "/authentication", bytes.NewBuffer([]byte(`{"type":"without_credentials"}`)))

	router.ServeHTTP(w, req)

	assert.Equal(t, 401, w.Code)
	assert.Equal(t, `{"errors":[{"code":"MP401","message":"Unauthorized"}]}`, w.Body.String())
}

func Test_Authentication_WhenInvalidCredentials_ReturnUnauthorized(t *testing.T) {
	router, w := arrangeMiddlewareRoute("/authentication", authentication)
	req, _ := http.NewRequest("POST", "/authentication", bytes.NewBuffer([]byte(`{"type":"invalid_credentials"}`)))
	req.SetBasicAuth("invalid_user", "invalid_pass")

	router.ServeHTTP(w, req)

	assert.Equal(t, 401, w.Code)
	assert.Equal(t, `{"errors":[{"code":"MP401","message":"Unauthorized"}]}`, w.Body.String())
}

func Test_Authentication_WhenInvalidPassword_ReturnUnauthorized(t *testing.T) {
	router, w := arrangeMiddlewareRoute("/authentication", authentication)
	user, _ := usermanagement.LoadMockUserCredentials()
	req, _ := http.NewRequest("POST", "/authentication", bytes.NewBuffer([]byte(`{"type":"valid_credentials"}`)))
	req.SetBasicAuth(user, "invalid pass")

	router.ServeHTTP(w, req)

	assert.Equal(t, 401, w.Code)
	assert.Equal(t, `{"errors":[{"code":"MP401","message":"Unauthorized"}]}`, w.Body.String())
}

func Test_Authentication_WhenValidCredentials_AuthorizedRequestSuccessful(t *testing.T) {
	router, w := arrangeMiddlewareRoute("/authentication", authentication)
	user, pass := usermanagement.LoadMockUserCredentials()
	req, _ := http.NewRequest("POST", "/authentication", bytes.NewBuffer([]byte(`{"type":"valid_credentials"}`)))
	req.SetBasicAuth(user, pass)

	router.ServeHTTP(w, req)

	assert.Equal(t, 200, w.Code)
}

func Test_ParseBoleto_WhenInvalidBody_ReturnBadRequest(t *testing.T) {
	router, w := arrangeMiddlewareRoute("/parseboleto", parseBoleto)
	req, _ := http.NewRequest("POST", "/parseboleto", bytes.NewBuffer([]byte(``)))

	router.ServeHTTP(w, req)

	assert.Equal(t, 400, w.Code)
	assert.Equal(t, `{"errors":[{"code":"MP400","message":"EOF"}]}`, w.Body.String())
}

func Test_ParseBoleto_WhenInvalidBank_ReturnBadRequest(t *testing.T) {
	router, w := arrangeMiddlewareRoute("/parseboleto", parseBoleto)
	req, _ := http.NewRequest("POST", "/parseboleto", bytes.NewBuffer([]byte(`{"bankNumber": 999}`)))

	router.ServeHTTP(w, req)

	assert.Equal(t, 400, w.Code)
	assert.Equal(t, `{"errors":[{"code":"MPBankNumber","message":"Banco 999 não existe"}]}`, w.Body.String())
}

func Test_ParseBoleto_WhenInvalidExpirationDate_ReturnBadRequest(t *testing.T) {
	router, w := arrangeMiddlewareRoute("/parseboleto", parseBoleto)
	body := test.NewStubBoletoRequest(models.BancoDoBrasil).Build()
	req, _ := http.NewRequest("POST", "/parseboleto", bytes.NewBuffer([]byte(util.ToJSON(body))))

	router.ServeHTTP(w, req)

	assert.Equal(t, 400, w.Code)
	assert.Equal(t, `{"errors":[{"code":"MP400","message":"parsing time \"\" as \"2006-01-02\": cannot parse \"\" as \"2006\""}]}`, w.Body.String())
}

func Test_ParseBoleto_WhenValidRequest_PassSuccessful(t *testing.T) {
	router, w := arrangeMiddlewareRoute("/parseboleto", parseBoleto)
	body := test.NewStubBoletoRequest(models.BancoDoBrasil).WithExpirationDate(time.Now()).Build()
	req, _ := http.NewRequest("POST", "/parseboleto", bytes.NewBuffer([]byte(util.ToJSON(body))))

	router.ServeHTTP(w, req)

	assert.Equal(t, 200, w.Code)
}

func Test_ValidateRegisterV1_WhenWithoutRules_PassSuccessful(t *testing.T) {
	router, w := arrangeMiddlewareRoute("/validateV1", parseBoleto, validateRegisterV1)
	body := test.NewStubBoletoRequest(models.BancoDoBrasil).WithExpirationDate(time.Now()).Build()
	req, _ := http.NewRequest("POST", "/validateV1", bytes.NewBuffer([]byte(util.ToJSON(body))))

	router.ServeHTTP(w, req)

	assert.Equal(t, 200, w.Code)
}

func Test_ValidateRegisterV1_WhenHasRules_ReturnBadRequest(t *testing.T) {
	router, w := arrangeMiddlewareRoute("/validateV1", parseBoleto, validateRegisterV1)
	body := test.NewStubBoletoRequest(models.Caixa).WithExpirationDate(time.Now()).WithAcceptDivergentAmount(true).Build()
	req, _ := http.NewRequest("POST", "/validateV1", bytes.NewBuffer([]byte(util.ToJSON(body))))

	router.ServeHTTP(w, req)

	assert.Equal(t, 400, w.Code)
	assert.Equal(t, `[{"code":"MP400","message":"title.rules not available in this version"}]`, w.Body.String())
}

func Test_ValidateRegisterV2_WhenWithoutRules_PassSuccessful(t *testing.T) {
	router, w := arrangeMiddlewareRoute("/validateV2", parseBoleto, validateRegisterV2)
	body := test.NewStubBoletoRequest(models.BancoDoBrasil).WithExpirationDate(time.Now()).Build()
	req, _ := http.NewRequest("POST", "/validateV2", bytes.NewBuffer([]byte(util.ToJSON(body))))

	router.ServeHTTP(w, req)

	assert.Equal(t, 200, w.Code)
}

func Test_ValidateRegisterV2_WhenHasRulesAndNotCaixaBank_ReturnBadRequest(t *testing.T) {
	router, w := arrangeMiddlewareRoute("/validateV2", parseBoleto, validateRegisterV2)
	body := test.NewStubBoletoRequest(models.BancoDoBrasil).WithExpirationDate(time.Now()).WithAcceptDivergentAmount(true).Build()
	req, _ := http.NewRequest("POST", "/validateV2", bytes.NewBuffer([]byte(util.ToJSON(body))))

	router.ServeHTTP(w, req)

	assert.Equal(t, 400, w.Code)
	assert.Equal(t, `[{"code":"MP400","message":"title.rules not available for this bank"}]`, w.Body.String())
}

func Test_ValidateRegisterV2_WhenHasRulesAndCaixaBank_PassSuccessful(t *testing.T) {
	router, w := arrangeMiddlewareRoute("/validateV2", parseBoleto, validateRegisterV2)
	body := test.NewStubBoletoRequest(models.Caixa).WithExpirationDate(time.Now()).WithAcceptDivergentAmount(true).Build()
	req, _ := http.NewRequest("POST", "/validateV2", bytes.NewBuffer([]byte(util.ToJSON(body))))

	router.ServeHTTP(w, req)

	assert.Equal(t, 200, w.Code)
}

func Test_ParseExpirationDate(t *testing.T) {
	expectedExpireDate := time.Now().Format("2006-01-02")
	expectedExpireDateTime, _ := time.Parse("2006-01-02", expectedExpireDate)

	boleto := models.BoletoRequest{BankNumber: models.BancoDoBrasil, Title: models.Title{ExpireDate: expectedExpireDate}}
	bank, _ := bank.Get(boleto)
	parseExpirationDate(nil, &boleto, bank)

	assert.Equal(t, expectedExpireDateTime, boleto.Title.ExpireDateTime)
}

func arrangeMiddlewareRoute(route string, handlers ...gin.HandlerFunc) (*gin.Engine, *httptest.ResponseRecorder) {
	router := mockInstallApi()
	router.POST(route, handlers...)
	w := httptest.NewRecorder()
	return router, w
}
