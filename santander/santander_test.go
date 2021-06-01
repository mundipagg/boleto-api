package santander

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
	"BankNumber": 33,
	"Agreement": {
		"AgreementNumber": 11111111,		
		"Agency":"5555",
		"Account":"55555"
	},
	"Title": {
		"ExpireDate": "2035-08-01",
		"AmountInCents": 200,
		"OurNumber":10000000004		
	},
	"Buyer": {
		"Name": "TESTE",
		"Document": {
			"Type": "CPF",
			"Number": "12345678903"
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
	{Input: models.Title{BoletoType: ""}, Expected: "02"},
	{Input: models.Title{BoletoType: "NSA"}, Expected: "02"},
	{Input: models.Title{BoletoType: "BDP"}, Expected: "32"},
	{Input: models.Title{BoletoType: "DM"}, Expected: "02"},
	{Input: models.Title{BoletoType: "DS"}, Expected: "04"},
	{Input: models.Title{BoletoType: "NP"}, Expected: "12"},
	{Input: models.Title{BoletoType: "RC"}, Expected: "17"},
	{Input: models.Title{BoletoType: "OUT"}, Expected: "99"},
}

func TestProcessBoleto_WhenServiceRespondsSuccessfully_ShouldHasSuccessfulBoletoResponse(t *testing.T) {
	mock.StartMockService("9098")
	input := new(models.BoletoRequest)
	util.FromJSON(baseMockJSON, input)
	bank, _ := New()

	output, _ := bank.ProcessBoleto(input)

	test.AssertProcessBoletoWithSuccess(t, output)
}

func TestGetBoletoType_WhenCalled_ShouldBeMapTypeSuccessful(t *testing.T) {
	request := new(models.BoletoRequest)
	for _, fact := range boletoTypeParameters {
		request.Title = fact.Input.(models.Title)
		_, result := getBoletoType(request)
		assert.Equal(t, fact.Expected, result, "Deve mapear o boleto type corretamente")
	}
}
