package citibank

import (
	"testing"

	"github.com/mundipagg/boleto-api/mock"
	"github.com/mundipagg/boleto-api/models"
	"github.com/mundipagg/boleto-api/test"
	"github.com/mundipagg/boleto-api/util"
	"github.com/stretchr/testify/assert"
)

const baseMockJSON = `
{
    "BankNumber": 745,
    "Authentication": {
        "Username": "55555555555555555555"
    },
    "Agreement": {
        "AgreementNumber": 55555555,
        "Wallet" : 100,
        "Agency":"0011",
        "Account":"0088881323",
        "AccountDigit" : "2"        
    },
    "Title": {
        "ExpireDate": "2029-09-20",
        "AmountInCents": 200,
        "OurNumber": 10000000001,
        "DocumentNumber":"5555555555"
    },
    "Buyer": {
        "Name": "Fulano de Tal",
        "Document": {
            "Type": "CNPJ",
            "Number": "55555555555555"
        }
    },
    "Recipient": {
      "Name": "TESTE",
        "Document": {
            "Type": "CNPJ",
            "Number": "55555555555555"
        }
    }
}
`

var boletoTypeParameters = []test.Parameter{
	{Input: models.Title{BoletoType: ""}, Expected: "03"},
	{Input: models.Title{BoletoType: "NSA"}, Expected: "03"},
	{Input: models.Title{BoletoType: "BDP"}, Expected: "03"},
	{Input: models.Title{BoletoType: "DMI"}, Expected: "03"},
}

func TestProcessBoleto_WhenServiceRespondsSuccessfully_ShouldHasSuccessfulBoletoResponse(t *testing.T) {
	mock.StartMockService("9095")
	input := new(models.BoletoRequest)
	util.FromJSON(baseMockJSON, input)
	bank, _ := New()

	output, _ := bank.ProcessBoleto(input)

	test.AssertProcessBoletoWithSuccess(t, output)
}

func TestCalculateOurNumber_WhenCalled_ShouldBeCalcutateOurNumberWithSuccess(t *testing.T) {
	boleto := models.BoletoRequest{
		Title: models.Title{
			OurNumber: 8605970,
		},
	}
	var expected uint = 86059700

	result := calculateOurNumber(&boleto)

	assert.Equal(t, expected, result, "Deve-se calcular corretamente o nosso numero para o Citi")
}

func TestGetBoletoType_WhenCalled_ShouldBeMapTypeSuccessful(t *testing.T) {
	for _, fact := range boletoTypeParameters {
		_, result := getBoletoType()
		assert.Equal(t, fact.Expected, result, "Deve mapear o boleto type corretamente")
	}
}
