package bb

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
	
	"BankNumber": 1,
	"Authentication": {
		"Username": "xxx",
		"Password": "xxx"
	},
	"Agreement": {
		"AgreementNumber": 5555555,
		"WalletVariation": 19,
		"Wallet":17,
		"Agency":"5555",
		"Account":"55555"
	},
	"Title": {
		"ExpireDate": "2029-10-01",
		"AmountInCents": 200,
		"OurNumber": 1,
		"Instructions": "Senhor caixa, não receber após o vencimento",
		"DocumentNumber": "123456"
	},
	"Buyer": {
		"Name": "TESTE",
		"Document": {
			"Type": "CNPJ",
			"Number": "55555555550140"
		},
		"Address": {
			"Street": "Teste",
			"Number": "123",
			"Complement": "Apto",
			"ZipCode": "55555555",
			"City": "Rio de Janeiro",
			"District": "Teste",
			"StateCode": "RJ"
		}
	},
	"Recipient": {
		"Name": "TESTE",
		"Document": {
			"Type": "CNPJ",
			"Number": "55555555555555"
		},
		"Address": {
			"Street": "TESTE",
			"Number": "555",
			"Complement": "Teste",
			"ZipCode": "0455555",
			"City": "São Paulo",
			"District": "",
			"StateCode": "SP"
		}
	}
}
`

var boletoTypeParameters = []test.Parameter{
	{Input: models.Title{BoletoType: ""}, Expected: "19"},
	{Input: models.Title{BoletoType: "NSA"}, Expected: "19"},
	{Input: models.Title{BoletoType: "BDP"}, Expected: "19"},
	{Input: models.Title{BoletoType: "CH"}, Expected: "01"},
	{Input: models.Title{BoletoType: "DM"}, Expected: "02"},
	{Input: models.Title{BoletoType: "DS"}, Expected: "04"},
	{Input: models.Title{BoletoType: "NP"}, Expected: "12"},
	{Input: models.Title{BoletoType: "RC"}, Expected: "17"},
	{Input: models.Title{BoletoType: "ND"}, Expected: "19"},
}

func TestProcessBoleto_WhenServiceRespondsSuccessfully_ShouldHasSuccessfulBoletoResponse(t *testing.T) {
	mock.StartMockService("9091")
	input := new(models.BoletoRequest)
	util.FromJSON(baseMockJSON, input)
	bank := New()

	output, _ := bank.ProcessBoleto(input)

	test.AssertProcessBoletoWithSuccess(t, output)
}

func TestProcessBoleto_WhenServiceRespondsFailed_ShouldHasFailedBoletoResponse(t *testing.T) {
	mock.StartMockService("9091")
	input := new(models.BoletoRequest)
	util.FromJSON(baseMockJSON, input)
	input.Title.AmountInCents = 400
	bank := New()

	output, _ := bank.ProcessBoleto(input)

	test.AssertProcessBoletoFailed(t, output)
}

func TestProcessBoleto_WhenAccountInvalid_ShouldHasFailedBoletoResponse(t *testing.T) {
	mock.StartMockService("9091")
	input := new(models.BoletoRequest)
	util.FromJSON(baseMockJSON, input)
	input.Agreement.Account = ""
	bank := New()

	output, _ := bank.ProcessBoleto(input)

	test.AssertProcessBoletoFailed(t, output)
}

func TestShouldCalculateAgencyDigitFromBb(t *testing.T) {
	test.ExpectTrue(bbAgencyDigitCalculator("0137") == "6", t)

	test.ExpectTrue(bbAgencyDigitCalculator("3418") == "5", t)

	test.ExpectTrue(bbAgencyDigitCalculator("3324") == "3", t)

	test.ExpectTrue(bbAgencyDigitCalculator("5797") == "5", t)
}

func TestShouldCalculateAccountDigitFromBb(t *testing.T) {
	test.ExpectTrue(bbAccountDigitCalculator("", "00006685") == "0", t)

	test.ExpectTrue(bbAccountDigitCalculator("", "00025619") == "6", t)

	test.ExpectTrue(bbAccountDigitCalculator("", "00006842") == "X", t)

	test.ExpectTrue(bbAccountDigitCalculator("", "00000787") == "0", t)
}

func TestGetBoletoType_WhenCalled_ShouldBeMapTypeSuccessful(t *testing.T) {
	request := new(models.BoletoRequest)
	for _, fact := range boletoTypeParameters {
		request.Title = fact.Input.(models.Title)
		_, result := getBoletoType(request)
		assert.Equal(t, fact.Expected, result, "Deve mapear o boleto type corretamente")
	}
}
