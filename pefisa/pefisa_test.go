package pefisa

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
    "bankNumber": 174,
    "authentication": {
            "Username": "altsa",
            "Password": "altsa"
	},
	"agreement": {
		"agreementNumber": 267,
		"wallet": 36,
		"agency": "00000"
	},
	"title": {           
		"expireDate": "2050-12-30",
		"amountInCents": 200,
		"ourNumber": 1,
		"instructions": "Não receber após a data de vencimento.",
		"documentNumber": "1234567890"
	},
	"recipient": {
		"name": "Empresa - Boletos",
		"document": {
			"type": "CNPJ",
			"number": "29799428000128"
		},
		"address": {
			"street": "Avenida Miguel Estefno, 2394",
			"complement": "Água Funda",
			"zipCode": "04301-002",
			"city": "São Paulo",
			"stateCode": "SP"
		}
	},
	"buyer": {
		"name": "Usuario Teste",
		"email": "p@p.com",
		"document": {
			"type": "CNPJ",
			"number": "29.799.428/0001-28"
		},
		"address": {
			"street": "Rua Teste",
			"number": "2",
			"complement": "SALA 1",
			"zipCode": "20931-001",
			"district": "Centro",
			"city": "Rio de Janeiro",
			"stateCode": "RJ"
		}
	}
}
`

var boletoTypeParameters = []test.Parameter{
	{Input: models.Title{BoletoType: ""}, Expected: "1"},
	{Input: models.Title{BoletoType: "NSA"}, Expected: "1"},
	{Input: models.Title{BoletoType: "DM"}, Expected: "1"},
	{Input: models.Title{BoletoType: "DS"}, Expected: "2"},
	{Input: models.Title{BoletoType: "NP"}, Expected: "3"},
	{Input: models.Title{BoletoType: "SE"}, Expected: "4"},
	{Input: models.Title{BoletoType: "CH"}, Expected: "10"},
	{Input: models.Title{BoletoType: "OUT"}, Expected: "99"},
}

func TestProcessBoleto_WhenServiceRespondsSuccessfully_ShouldHasSuccessfulBoletoResponse(t *testing.T) {
	mock.StartMockService("9097")
	input := new(models.BoletoRequest)
	util.FromJSON(baseMockJSON, input)
	bank := New()

	output, _ := bank.ProcessBoleto(input)

	test.AssertProcessBoletoWithSuccess(t, output)
}

func TestProcessBoleto_WhenServiceRespondsFailed_ShouldHasFailedBoletoResponse(t *testing.T) {
	mock.StartMockService("9097")
	input := new(models.BoletoRequest)
	util.FromJSON(baseMockJSON, input)
	input.Title.AmountInCents = 201
	bank := New()

	output, _ := bank.ProcessBoleto(input)

	test.AssertProcessBoletoFailed(t, output)
}

func TestGetBoletoType_WhenCalled_ShouldBeMapTypeSuccessful(t *testing.T) {
	request := new(models.BoletoRequest)
	for _, fact := range boletoTypeParameters {
		request.Title = fact.Input.(models.Title)
		_, result := getBoletoType(request)
		assert.Equal(t, fact.Expected, result, "Deve mapear o boleto type corretamente")
	}
}
