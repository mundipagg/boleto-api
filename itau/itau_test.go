package itau

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
	"BankNumber": 341,
	"Authentication": {
		"Username": "a",
		"Password": "b",
		"AccessKey":"c"
	},
	"Agreement": {
		"Wallet":109,
		"Agency":"0407",
		"Account":"55292",
		"AccountDigit":"6"
	},
	"Title": {
		"ExpireDate": "2999-12-31",
		"AmountInCents": 200			
	},
	"Buyer": {
		"Name": "TESTE",
		"Document": {
			"Type": "CNPJ",
			"Number": "00001234567890"
		}
	},
	"Recipient": {
		"Name": "TESTE",
		"Document": {
			"Type": "CNPJ",
			"Number": "00123456789067"
		}
	}
}
`

var boletoTypeParameters = []test.Parameter{
	{Input: models.Title{BoletoType: ""}, Expected: "01"},
	{Input: models.Title{BoletoType: "NSA"}, Expected: "01"},
	{Input: models.Title{BoletoType: "BDP"}, Expected: "18"},
	{Input: models.Title{BoletoType: "DM"}, Expected: "01"},
	{Input: models.Title{BoletoType: "DS"}, Expected: "08"},
	{Input: models.Title{BoletoType: "NP"}, Expected: "02"},
	{Input: models.Title{BoletoType: "RC"}, Expected: "05"},
	{Input: models.Title{BoletoType: "OUT"}, Expected: "99"},
}

func TestProcessBoleto_WhenServiceRespondsSuccessfully_ShouldHasSuccessfulBoletoResponse(t *testing.T) {
	mock.StartMockService("9096")
	input := new(models.BoletoRequest)
	util.FromJSON(baseMockJSON, input)
	bank := New()

	output, _ := bank.ProcessBoleto(input)

	test.AssertProcessBoletoWithSuccess(t, output)
}

func TestProcessBoleto_WhenServiceRespondsFailed_ShouldHasFailedBoletoResponse(t *testing.T) {
	mock.StartMockService("9096")
	input := new(models.BoletoRequest)
	util.FromJSON(baseMockJSON, input)
	input.Title.AmountInCents = 400
	bank := New()

	output, _ := bank.ProcessBoleto(input)

	test.AssertProcessBoletoFailed(t, output)
}

func TestProcessBoleto_WhenRequestHasInvalidAccountParameters_ShouldHasFailedBoletoResponse(t *testing.T) {
	mock.StartMockService("9096")
	input := new(models.BoletoRequest)
	util.FromJSON(baseMockJSON, input)
	input.Title.AmountInCents = 200
	input.Agreement.Account = ""
	bank := New()

	output, _ := bank.ProcessBoleto(input)

	test.AssertProcessBoletoFailed(t, output)
}

func TestProcessBoleto_WhenRequestHasInvalidUserNameParameter_ShouldHasFailedBoletoResponse(t *testing.T) {
	mock.StartMockService("9096")
	input := new(models.BoletoRequest)
	util.FromJSON(baseMockJSON, input)
	input.Title.AmountInCents = 200
	input.Authentication.Username = ""
	bank := New()

	_, err := bank.ProcessBoleto(input)

	assert.NotNil(t, err, "Deve ocorrer um erro")
}

func TestGetBoletoType_WhenCalled_ShouldBeMapTypeSuccessful(t *testing.T) {
	request := new(models.BoletoRequest)
	for _, fact := range boletoTypeParameters {
		request.Title = fact.Input.(models.Title)
		_, result := getBoletoType(request)
		assert.Equal(t, fact.Expected, result, "Deve mapear o boleto type corretamente")
	}
}

