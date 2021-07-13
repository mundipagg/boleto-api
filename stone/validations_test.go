package stone

import (
	"testing"

	"github.com/mundipagg/boleto-api/models"
	"github.com/stretchr/testify/assert"
)

func Test_ValidationAccessKey_WithSucessful(t *testing.T) {
	input := newStubBoletoRequestStone().Build()

	result := stoneValidateAccessKeyNotEmpty(input)

	assert.Nil(t, result)
}

func Test_ValidationAccessKey_WhenNotBoletoRequest_ReturnInvalidType(t *testing.T) {

	result := stoneValidateAccessKeyNotEmpty("input")

	assert.NotNil(t, result)
	assert.IsType(t, models.ErrorResponse{}, result)
	assert.Equal(t, "MP500", result.(models.ErrorResponse).Code)
}

func Test_ValidationAccessKey_WhenNotFill_ReturnBadRequestError(t *testing.T) {
	input := newStubBoletoRequestStone().WithAccessKey("").Build()

	result := stoneValidateAccessKeyNotEmpty(input)

	assert.NotNil(t, result)
	assert.IsType(t, models.ErrorResponse{}, result)
	assert.Equal(t, "MP400", result.(models.ErrorResponse).Code)
}
