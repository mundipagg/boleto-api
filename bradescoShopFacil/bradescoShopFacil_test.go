package bradescoShopFacil

import (
	"testing"
	"time"

	"github.com/mundipagg/boleto-api/mock"
	"github.com/mundipagg/boleto-api/models"
	"github.com/mundipagg/boleto-api/test"
	"github.com/mundipagg/boleto-api/util"
	"github.com/stretchr/testify/assert"
)

const baseMockJSON = `
{
    "BankNumber": 237,
     "Authentication": {
        "Username": "55555555555",
        "Password": "55555555555555555"
    },
    "Agreement": {
        "AgreementNumber": 55555555,
        "Wallet": 25,
        "Agency":"5555",
        "Account":"55555"
    },
    "Title": {
        "ExpireDate": "2029-08-01",
        "AmountInCents": 200,
        "OurNumber": 12446688,
        "Instructions": "Senhor caixa, não receber após o vencimento",
        "DocumentNumber": "1234566"
    },
    "Buyer": {
        "Name": "Luke Skywalker",
        "Document": {
            "Type": "CPF",
            "Number": "01234567890"
        },
        "Address": {
            "Street": "Mos Eisley Cantina",
            "Number": "123",
            "Complement": "Apto",
            "ZipCode": "20001-000",
            "City": "Tatooine",
            "District": "Tijuca",
            "StateCode": "RJ"
        }
    },
    "Recipient": {
      "Name": "TESTE",
        "Document": {
            "Type": "CNPJ",
            "Number": "00555555000109"
        },

        "Address": {
            "Street": "TESTE",
            "Number": "111",
            "Complement": "TESTE",
            "ZipCode": "11111111",
            "City": "Teste",
            "District": "",
            "StateCode": "SP"
        }

    }
}
`

var boletoTypeParameters = []test.Parameter{
	{Input: models.Title{BoletoType: ""}, Expected: "01"},
	{Input: models.Title{BoletoType: "NSA"}, Expected: "01"},
	{Input: models.Title{BoletoType: "BDP"}, Expected: "01"},
	{Input: models.Title{BoletoType: "DM"}, Expected: "01"},
	{Input: models.Title{BoletoType: "DS"}, Expected: "12"},
	{Input: models.Title{BoletoType: "NP"}, Expected: "02"},
	{Input: models.Title{BoletoType: "RC"}, Expected: "05"},
	{Input: models.Title{BoletoType: "OUT"}, Expected: "99"},
}

func TestProcessBoleto_WhenServiceRespondsSuccessfully_ShouldHasSuccessfulBoletoResponse(t *testing.T) {
	mock.StartMockService("9093")
	input := new(models.BoletoRequest)
	util.FromJSON(baseMockJSON, input)
	bank := New()

	output, _ := bank.ProcessBoleto(input)

	test.AssertProcessBoletoWithSuccess(t, output)
}

func TestProcessBoleto_WhenServiceRespondsFailed_ShouldHasFailedBoletoResponse(t *testing.T) {
	mock.StartMockService("9093")
	input := new(models.BoletoRequest)
	util.FromJSON(baseMockJSON, input)
	input.Title.AmountInCents = 400
	bank := New()

	output, _ := bank.ProcessBoleto(input)

	test.AssertProcessBoletoFailed(t, output)
}

func TestBarcodeGenerationBradescoShopFacil(t *testing.T) {
	const expected = "23795796800000001990001250012446693212345670"

	boleto := models.BoletoRequest{}
	boleto.BankNumber = models.Bradesco
	boleto.Agreement = models.Agreement{
		Account: "1234567",
		Agency:  "1",
		Wallet:  25,
	}
	expireDate, _ := time.Parse("02-01-2006", "01-08-2019")
	boleto.Title = models.Title{
		AmountInCents:  199,
		OurNumber:      124466932,
		ExpireDateTime: expireDate,
	}
	bc := getBarcode(boleto)

	assert.Equal(t, expected, bc.toString(), "Deve-se montar o código de barras do BradescoShopFacil")
}

func TestRemoveDigitFromAccount(t *testing.T) {
	const expected = "23791796800000001992372250012446693300056000"

	bc := barcode{
		account:       "0005600",
		bankCode:      "237",
		currencyCode:  "9",
		agency:        "2372",
		dateDueFactor: "7968",
		ourNumber:     "00124466933",
		zero:          "0",
		wallet:        "25",
		value:         "0000000199",
	}

	assert.Equal(t, expected, bc.toString(), "Deve-se montar identificar e remover o digito da conta")
}

func TestGetBoletoType_WhenCalled_ShouldBeMapTypeSuccessful(t *testing.T) {
	request := new(models.BoletoRequest)
	for _, fact := range boletoTypeParameters {
		request.Title = fact.Input.(models.Title)
		_, result := getBoletoType(request)
		assert.Equal(t, fact.Expected, result, "Deve mapear o boleto type corretamente")
	}
}
