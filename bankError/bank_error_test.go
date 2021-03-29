package bankError

import (
	"testing"

	"github.com/mundipagg/boleto-api/models"
	"github.com/stretchr/testify/assert"
)

func TestParseError_WhenErrorCodeIsAKnownErrorAndKnowBank_ThenReturnShouldNotBeNull(t *testing.T) {

	err := models.NewErrorResponse("810", "Erro de certificado ou formatação de campos")
	errorResponse := ParseError(err, "BradescoNetEmpresa")

	assert.NotNil(t, errorResponse)
}

func TestParseError_WhenErrorCodeIsAnUnknownError_ThenReturnShouldBeNull(t *testing.T) {

	err := models.NewErrorResponse("UnknowError", "This is an unknown error")
	errorResponse := ParseError(err, "BradescoNetEmpresa")

	assert.Nil(t, errorResponse)
}

func TestParseError_WhenBankIsNotMapped_ThenReturnShouldBeNull(t *testing.T) {

	err := models.NewErrorResponse("0", "Fim anormal do programa")
	errorResponse := ParseError(err, "UnknowBank")

	assert.Nil(t, errorResponse)
}
