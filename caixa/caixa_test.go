package caixa

import (
	"fmt"
	"testing"
	"time"

	"github.com/PMoneda/flow"
	"github.com/mundipagg/boleto-api/mock"
	"github.com/mundipagg/boleto-api/models"
	"github.com/mundipagg/boleto-api/test"
	"github.com/mundipagg/boleto-api/tmpl"
	"github.com/mundipagg/boleto-api/util"
	"github.com/stretchr/testify/assert"
)

const baseMockJSON = `
{
	"BankNumber": 104,
	"Agreement": {
		"AgreementNumber": 555555,
		"Agency":"5555"
	},
	"Title": {
		"ExpireDate": "2029-08-30",
		"AmountInCents": 200,
		"OurNumber": 0,
		"Instructions": "Mensagem",
		"DocumentNumber": "NPC160517"
	},
	"Buyer": {
		"Name": "TESTE PAGADOR 001",
		"Document": {
			"Type": "CPF",
			"Number": "57962014849"
		},
		"Address": {
			"Street": "SAUS QUADRA 03",
			"Number": "",
			"Complement": "",
			"ZipCode": "20520051",
			"City": "Rio de Janeiro",
			"District": "Tijuca",
			"StateCode": "RJ"
		}
	},
	"Recipient": {
		"Document": {
			"Type": "CNPJ",
			"Number": "00555555000109"
		}
	}
}
`

var boletoTypeParameters = []test.Parameter{
	{Input: models.Title{BoletoType: ""}, Expected: "99"},
	{Input: models.Title{BoletoType: "NSA"}, Expected: "99"},
	{Input: models.Title{BoletoType: "BDP"}, Expected: "99"},
}

var boletoBuyerNameParameters = []test.Parameter{
	{Input: "Leonardo Jasmim", Expected: "<NOME>Leonardo Jasmim</NOME>"},
	{Input: "Ântôníõ Tùpìnâmbáú", Expected: "<NOME>Antonio Tupinambau</NOME>"},
	{Input: "Accepted , / ( ) * = - + ! : ? . ; _ ' ", Expected: "<NOME>Accepted , / ( ) * = - &#43; ! : ? . ; _ &#39; </NOME>"},
	{Input: "NotAccepted @#$%¨{}[]^~\"&<>\\", Expected: "<NOME>NotAccepted                 </NOME>"},
}

var boletoInstructionsParameters = []test.Parameter{
	{Input: ", / ( ) * = - + ! : ? . ; _ ' ", Expected: "<MENSAGEM>, / ( ) * = - &#43; ! : ? . ; _ &#39; </MENSAGEM>"},
	{Input: "@ # $ % ¨ { } [ ] ^ ~ \" & < > \\", Expected: "                              "},
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
	mock.StartMockService("9092")
	input := new(models.BoletoRequest)
	util.FromJSON(baseMockJSON, input)
	input.Title.AmountInCents = 400
	bank := New()

	output, _ := bank.ProcessBoleto(input)

	test.AssertProcessBoletoFailed(t, output)
}

func TestProcessBoleto_WhenRequestContainsInvalidParameters_ShouldHasFailedBoletoResponse(t *testing.T) {
	mock.StartMockService("9092")
	input := new(models.BoletoRequest)
	util.FromJSON(baseMockJSON, input)
	input.Title.Instructions = "aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa"
	input.Title.OurNumber = 9999999999999999
	bank := New()

	output, _ := bank.ProcessBoleto(input)

	test.AssertProcessBoletoFailed(t, output)
}

func TestGetCaixaCheckSumInfo(t *testing.T) {
	const expectedSumCode = "0200656000000000000000003008201700000000000100000732159000109"
	const expectedToken = "LvWr1op5Ayibn6jsCQ3/2bW4KwThVAlLK5ftxABlq20="

	boleto := models.BoletoRequest{
		Agreement: models.Agreement{
			AgreementNumber: 200656,
		},
		Title: models.Title{
			OurNumber:      0,
			ExpireDateTime: time.Date(2017, 8, 30, 12, 12, 12, 12, time.Local),
			AmountInCents:  1000,
		},
		Recipient: models.Recipient{
			Document: models.Document{
				Number: "00732159000109",
			},
		},
	}
	caixa := New()

	assert.Equal(t, expectedSumCode, caixa.getCheckSumCode(boleto), "Deve-se formar uma string seguindo o padrão da documentação")
	assert.Equal(t, expectedToken, caixa.getAuthToken(caixa.getCheckSumCode(boleto)), "Deve-se fazer um hash sha256 e encodar com base64")
}

func TestShouldCalculateAccountDigitCaixa(t *testing.T) {
	boleto := models.BoletoRequest{
		Agreement: models.Agreement{
			Account: "100000448",
			Agency:  "2004",
		},
	}

	assert.Nil(t, caixaValidateAccountAndDigit(&boleto))
	assert.Nil(t, caixaValidateAgency(&boleto))
}

func TestGetBoletoType_WhenCalled_ShouldBeMapTypeSuccessful(t *testing.T) {
	request := new(models.BoletoRequest)
	for _, fact := range boletoTypeParameters {
		request.Title = fact.Input.(models.Title)
		_, result := getBoletoType(request)
		assert.Equal(t, fact.Expected, result, "Deve mapear o boleto type corretamente")
	}
}

func TestTempletaRequestCaixa_WhenParseBuyerName_ShouldParseSucessfully(t *testing.T) {
	f := flow.NewFlow()

	request := new(models.BoletoRequest)
	util.FromJSON(baseMockJSON, request)

	for _, fact := range boletoBuyerNameParameters {
		request.Buyer.Name = fact.Input.(string)
		result := fmt.Sprintf("%v", f.From("message://?source=inline", request, getRequestCaixa(), tmpl.GetFuncMaps()).GetBody())
		assert.Contains(t, result, fact.Expected, "Conversão não realizada como esperado")
	}
}

func TestTempletaRequestCaixa_WhenParseInstruction_ShouldParseSucessfully(t *testing.T) {
	f := flow.NewFlow()

	request := new(models.BoletoRequest)
	util.FromJSON(baseMockJSON, request)

	for _, fact := range boletoInstructionsParameters {
		request.Title.Instructions = fact.Input.(string)
		result := fmt.Sprintf("%v", f.From("message://?source=inline", request, getRequestCaixa(), tmpl.GetFuncMaps()).GetBody())
		assert.Contains(t, result, fact.Expected, "Conversão não realizada como esperado")
	}
}
