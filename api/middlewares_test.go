package api

import (
	"bytes"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/mundipagg/boleto-api/bank"
	"github.com/mundipagg/boleto-api/log"
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

func Test_ParseExpirationDate(t *testing.T) {
	expectedExpireDate := time.Now().Format("2006-01-02")
	expectedExpireDateTime, _ := time.Parse("2006-01-02", expectedExpireDate)

	boleto := models.BoletoRequest{BankNumber: models.BancoDoBrasil, Title: models.Title{ExpireDate: expectedExpireDate}}
	bank, _ := bank.Get(boleto)
	parseExpirationDate(nil, &boleto, bank)

	assert.Equal(t, expectedExpireDateTime, boleto.Title.ExpireDateTime)
}

func Test_LoadLog(t *testing.T) {
	expectedIP := "127.0.0.1"
	expectedUser := "user"
	requestOurNumber := uint(1234567890)
	expectedOurNumber := strconv.FormatUint(uint64(requestOurNumber), 10)

	boleto := test.NewStubBoletoRequest(models.BancoDoBrasil).WithOurNumber(requestOurNumber).Build()
	bank, _ := bank.Get(*boleto)

	ginCtx, _ := gin.CreateTestContext(httptest.NewRecorder())
	ginCtx.Set(serviceUserKey, expectedUser)
	ginCtx.Set(boletoKey, *boleto)
	ginCtx.Set(bankKey, bank)
	ginCtx.Request, _ = http.NewRequest("POST", "/", nil)
	ginCtx.Request.Header.Set("X-Forwarded-For", expectedIP)

	l := loadBankLog(ginCtx)

	assert.NotNil(t, l)
	assert.Equal(t, expectedIP, l.IPAddress)
	assert.Equal(t, expectedUser, l.ServiceUser)
	assert.Equal(t, expectedOurNumber, l.NossoNumero)
	assert.Equal(t, bank.GetBankNameIntegration(), l.BankName)
}

func Test_CheckError_WhenNotFoundError(t *testing.T) {
	_, w := arrangeMiddlewareRoute("/err", gin.Default().HandleContext)
	ginCtx, _ := gin.CreateTestContext(w)
	err := models.NewHTTPNotFound("404", "objeto não encontrado")
	l := log.CreateLog()

	checkError(ginCtx, err, l)

	assert.Equal(t, http.StatusNotFound, w.Result().StatusCode)
	assert.Equal(t, `{"errors":[{"code":"MP404","message":"objeto não encontrado"}]}`, w.Body.String())
}

func Test_CheckError_WhenInternalServerError(t *testing.T) {
	_, w := arrangeMiddlewareRoute("/err", gin.Default().HandleContext)
	ginCtx, _ := gin.CreateTestContext(w)
	err := models.NewInternalServerError("500", "erro interno")
	l := log.CreateLog()

	checkError(ginCtx, err, l)

	assert.Equal(t, http.StatusInternalServerError, w.Result().StatusCode)
	assert.Equal(t, `{"errors":[{"code":"MP500","message":"erro interno"}]}`, w.Body.String())
}

func Test_CheckError_WhenBadGatewayError(t *testing.T) {
	_, w := arrangeMiddlewareRoute("/err", gin.Default().HandleContext)
	ginCtx, _ := gin.CreateTestContext(w)
	err := models.NewBadGatewayError("erro externo")
	l := log.CreateLog()

	checkError(ginCtx, err, l)

	assert.Equal(t, http.StatusBadGateway, w.Result().StatusCode)
	assert.Equal(t, `{"errors":[{"code":"MP502","message":"erro externo"}]}`, w.Body.String())
}

func Test_CheckError_WhenGenericError(t *testing.T) {
	_, w := arrangeMiddlewareRoute("/err", gin.Default().HandleContext)
	ginCtx, _ := gin.CreateTestContext(w)
	err := fmt.Errorf("erro generico")
	l := log.CreateLog()

	checkError(ginCtx, err, l)

	assert.Equal(t, http.StatusInternalServerError, w.Result().StatusCode)
	assert.Equal(t, `{"errors":[{"code":"MP500","message":"Internal Error"}]}`, w.Body.String())
}

func Test_GetOurNumberFromContext_WhenOurNumberInRequest(t *testing.T) {
	requestOurNumber := uint(12345678901234567890)
	expectedOurNumber := strconv.FormatUint(uint64(requestOurNumber), 10)

	c, _ := gin.CreateTestContext(httptest.NewRecorder())
	request := test.NewStubBoletoRequest(models.BancoDoBrasil).WithOurNumber(requestOurNumber).Build()
	c.Set(boletoKey, *request)

	result := getNossoNumeroFromContext(c)

	assert.Equal(t, expectedOurNumber, result)
}

func Test_GetOurNumberFromContext_WhenOurNumberInResponse(t *testing.T) {
	expectedOurNumber := "123456789012345678901234567890"
	c, _ := gin.CreateTestContext(httptest.NewRecorder())
	response := models.BoletoResponse{OurNumber: expectedOurNumber}
	c.Set(responseKey, response)

	result := getNossoNumeroFromContext(c)

	assert.Equal(t, expectedOurNumber, result)
}

func arrangeMiddlewareRoute(route string, handlers ...gin.HandlerFunc) (*gin.Engine, *httptest.ResponseRecorder) {
	router := mockInstallApi()
	router.POST(route, handlers...)
	w := httptest.NewRecorder()
	return router, w
}
