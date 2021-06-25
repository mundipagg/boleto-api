package api

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/mundipagg/boleto-api/config"
	"github.com/mundipagg/boleto-api/models"
	"github.com/stretchr/testify/assert"
)

func TestGetBoleto_WhenInvalidKeys_ShouldReturnNotFound(t *testing.T) {
	expected := models.ErrorResponse{Code: "MP404", Message: "Not Found"}

	c, r, w := arrangeGetBoleto()
	url := "http://localhost:3000/boleto?fmt=html&pk=1234567890"
	c.Request, _ = http.NewRequest(http.MethodGet, url, nil)

	r.ServeHTTP(w, c.Request)

	var response models.BoletoResponse
	json.Unmarshal(w.Body.Bytes(), &response)

	assert.Equal(t, 404, w.Code)
	assert.Equal(t, 1, len(response.Errors))
	assert.Equal(t, expected.Code, response.Errors[0].Code, "O erro code deverá ser mapeado corretamente")
	assert.Equal(t, expected.Message, response.Errors[0].Message, "O erro message deverá ser mapeado corretamente")

}

func TestGetBoleto_WhenFail_ShouldReturnInternalError(t *testing.T) {
	expected := models.ErrorResponse{Code: "MP500", Message: "Internal Error"}

	c, r, w := arrangeGetBoleto()
	url := "http://localhost:3000/boleto?fmt=html&id=1234567890&pk=1234567890"
	c.Request, _ = http.NewRequest(http.MethodGet, url, nil)

	r.ServeHTTP(w, c.Request)

	var response models.BoletoResponse
	json.Unmarshal(w.Body.Bytes(), &response)

	assert.Equal(t, 500, w.Code)
	assert.Equal(t, 1, len(response.Errors))
	assert.Equal(t, expected.Code, response.Errors[0].Code, "O erro code deverá ser mapeado corretamente")
	assert.Equal(t, expected.Message, response.Errors[0].Message, "O erro message deverá ser mapeado corretamente")
}

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

func arrangeGetBoleto() (*gin.Context, *gin.Engine, *httptest.ResponseRecorder) {
	gin.SetMode(gin.TestMode)
	config.Install(true, false, true)
	w := httptest.NewRecorder()
	c, r := gin.CreateTestContext(w)
	r.GET("/boleto", getBoleto)
	return c, r, w
}
