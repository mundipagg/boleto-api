package bradescoNetEmpresa

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
    "bankNumber": 237,
   "authentication": {
            "Username": "",
            "Password": ""
        },
        "agreement": {
            "agreementNumber": 5822351,
            "wallet": 9,
            "agency": "1111",
            "account": "0062145"
        },
        "title": {
           
            "expireDate": "2050-12-30",
            "amountInCents": 200,
            "ourNumber": 12345678901,
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
	{Input: models.Title{BoletoType: ""}, Expected: "02"},
	{Input: models.Title{BoletoType: "NSA"}, Expected: "02"},
	{Input: models.Title{BoletoType: "BDP"}, Expected: "02"},
	{Input: models.Title{BoletoType: "CH"}, Expected: "01"},
	{Input: models.Title{BoletoType: "DM"}, Expected: "02"},
	{Input: models.Title{BoletoType: "DS"}, Expected: "04"},
	{Input: models.Title{BoletoType: "NP"}, Expected: "12"},
	{Input: models.Title{BoletoType: "RC"}, Expected: "17"},
	{Input: models.Title{BoletoType: "OUT"}, Expected: "99"},
}

func TestProcessBoleto_WhenServiceRespondsSuccessfully_ShouldHasSuccessfulBoletoResponse(t *testing.T) {
	mock.StartMockService("9092")
	input := new(models.BoletoRequest)
	util.FromJSON(baseMockJSON, input)
	bank := New()

	output, _ := bank.ProcessBoleto(input)

	test.AssertProcessBoletoWithSuccess(t, output)
}

func TestProcessBoleto_WhenServiceRespondsFailed_ShouldHasFailedBoletoResponse(t *testing.T) {
	mock.StartMockService("9092")
	input := new(models.BoletoRequest)
	util.FromJSON(baseMockJSON, input)
	input.Title.AmountInCents = 201
	bank := New()

	output, _ := bank.ProcessBoleto(input)

	test.AssertProcessBoletoFailed(t, output)
}

func TestProcessBoleto_WhenServiceRespondsCertificateFailed_ShouldHasFailedBoletoResponse(t *testing.T) {
	mock.StartMockService("9092")
	input := new(models.BoletoRequest)
	util.FromJSON(baseMockJSON, input)
	input.Title.AmountInCents = 202
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
