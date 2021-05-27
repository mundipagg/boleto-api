package models

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

var keysParameters = []ModelTestParameter{
	{Input: GetBoletoResult{Id: "", PrivateKey: ""}, Expected: false},
	{Input: GetBoletoResult{Id: "123", PrivateKey: ""}, Expected: false},
	{Input: GetBoletoResult{Id: "", PrivateKey: "123"}, Expected: false},
	{Input: GetBoletoResult{Id: "111", PrivateKey: "111"}, Expected: true},
}

func TestNewGetBoletoResult_WhenCalled_CreateNewGetBoletoResultSuccessful(t *testing.T) {
	expectedId := "60a56171ea46f83aecb67af9"
	expectedPk := "855a960c435a01863ad13c71d988e56e1f3931980a749b2a69ac982b7e03d16c"
	expectedFormat := "html"
	expectedSource := "none"
	expectedUri := fmt.Sprintf("/boleto?fmt=%s&id=%s&pk=%s", expectedFormat, expectedId, expectedPk)
	url := fmt.Sprintf("http://localhost:3000/boleto?fmt=%s&id=%s&pk=%s", expectedFormat, expectedId, expectedPk)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request, _ = http.NewRequest(http.MethodGet, url, nil)
	c.Request.RequestURI = expectedUri

	r := NewGetBoletoResult(c)

	assert.Equal(t, expectedId, r.Id, "O id deve ser igual ao do contexto")
	assert.Equal(t, expectedPk, r.PrivateKey, "A pk deve ser igual ao do contexto")
	assert.Equal(t, expectedFormat, r.Format, "O formato deve ser igual ao do contexto")
	assert.Equal(t, expectedUri, r.URI, "A uri deve ser igual ao do contexto")
	assert.Equal(t, expectedSource, r.BoletoSource, "O source deve ser none")
}

func TestHasValidKey_WhenCalled_CheckQueryParametersSuccessful(t *testing.T) {
	for _, fact := range keysParameters {
		input := fact.Input.(GetBoletoResult)
		result := input.HasValidKeys()
		assert.Equal(t, fact.Expected.(bool), result, fmt.Sprintf("Os parametros não foram validados corretamente: %v", result))
	}
}

func TestSetErrorResponse_WhenNotFoundOccurs_ShouldBeGeneratedWarningError(t *testing.T) {
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	var r GetBoletoResult

	expectedError := NewErrorResponse("MP404", "Not Found")
	expectedStatusCode := 404

	r.SetErrorResponse(c, expectedError, expectedStatusCode)

	assert.Equal(t, expectedStatusCode, c.Writer.Status(), "O StausCode deve ser mapeado para o contexto corretamente")
	assert.NotEmpty(t, r.ErrorResponse.Errors, "O result deve conter um objeto de erro")
	assert.Equal(t, expectedError.Code, r.ErrorResponse.Errors[0].Code, "O erro code deverá ser mapeado corretamente")
	assert.Equal(t, expectedError.Message, r.ErrorResponse.Errors[0].Message, "O erro message deverá ser mapeado corretamente")
}

func TestSetErrorResponse_WhenInternalErrosOccurs_ShouldBeGeneratedError(t *testing.T) {
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	var r GetBoletoResult

	expectedError := NewErrorResponse("MP500", "Internal Error")
	expectedStatusCode := 500

	r.SetErrorResponse(c, expectedError, expectedStatusCode)

	assert.Equal(t, expectedStatusCode, c.Writer.Status(), "O StausCode deve ser mapeado para o contexto corretamente")
	assert.NotEmpty(t, r.ErrorResponse.Errors, "O result deve conter um objeto de erro")
	assert.Equal(t, expectedError.Code, r.ErrorResponse.Errors[0].Code, "O erro code deverá ser mapeado corretamente")
	assert.Equal(t, expectedError.Message, r.ErrorResponse.Errors[0].Message, "O erro message deverá ser mapeado corretamente")
}
