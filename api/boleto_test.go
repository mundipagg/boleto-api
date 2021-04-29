package api

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_GetBoletoConfirmation_ReturnOkSuccessful(t *testing.T) {
	router := mockInstallApi()

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/boleto/confirmation", nil)

	router.ServeHTTP(w, req)

	assert.Equal(t, 200, w.Code)
	assert.Equal(t, "OK", w.Body.String())
}

func Test_PostBoletoConfirmation_ReturnOkSuccessful(t *testing.T) {
	router := mockInstallApi()

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/boleto/confirmation", nil)

	router.ServeHTTP(w, req)

	assert.Equal(t, 200, w.Code)
	assert.Equal(t, "OK", w.Body.String())
}
