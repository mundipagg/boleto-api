package middleware

import (
	"bytes"
	"net/http"
	"testing"
	"time"

	"github.com/mundipagg/boleto-api/models"
	"github.com/mundipagg/boleto-api/test"
	"github.com/mundipagg/boleto-api/util"
	"github.com/stretchr/testify/assert"
)

func Test_ParseBoleto_WhenInvalidBody_ReturnBadRequest(t *testing.T) {
	router, w := arrangeMiddlewareRoute("/parseboleto", ParseBoleto)
	req, _ := http.NewRequest("POST", "/parseboleto", bytes.NewBuffer([]byte(``)))

	router.ServeHTTP(w, req)

	assert.Equal(t, 400, w.Code)
	assert.Equal(t, `{"errors":[{"code":"MP400","message":"EOF"}]}`, w.Body.String())
}

func Test_ParseBoleto_WhenInvalidBank_ReturnBadRequest(t *testing.T) {
	router, w := arrangeMiddlewareRoute("/parseboleto", ParseBoleto)
	req, _ := http.NewRequest("POST", "/parseboleto", bytes.NewBuffer([]byte(`{"bankNumber": 999}`)))

	router.ServeHTTP(w, req)

	assert.Equal(t, 400, w.Code)
	assert.Equal(t, `{"errors":[{"code":"MPBankNumber","message":"Banco 999 não existe"}]}`, w.Body.String())
}

func Test_ParseBoleto_WhenInvalidExpirationDate_ReturnBadRequest(t *testing.T) {
	router, w := arrangeMiddlewareRoute("/parseboleto", ParseBoleto)
	body := test.NewStubBoletoRequest(models.BancoDoBrasil).Build()
	req, _ := http.NewRequest("POST", "/parseboleto", bytes.NewBuffer([]byte(util.ToJSON(body))))

	router.ServeHTTP(w, req)

	assert.Equal(t, 400, w.Code)
	assert.Equal(t, `{"errors":[{"code":"MP400","message":"parsing time \"\" as \"2006-01-02\": cannot parse \"\" as \"2006\""}]}`, w.Body.String())
}

func Test_ParseBoleto_WhenValidRequest_PassSuccessful(t *testing.T) {
	router, w := arrangeMiddlewareRoute("/parseboleto", ParseBoleto)
	body := test.NewStubBoletoRequest(models.BancoDoBrasil).WithExpirationDate(time.Now()).Build()
	req, _ := http.NewRequest("POST", "/parseboleto", bytes.NewBuffer([]byte(util.ToJSON(body))))

	router.ServeHTTP(w, req)

	assert.Equal(t, 200, w.Code)
}