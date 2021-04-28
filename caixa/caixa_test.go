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
	"github.com/stretchr/testify/assert"
)

var boletoTypeParameters = []test.Parameter{
	{Input: models.Title{BoletoType: ""}, Expected: "99"},
	{Input: models.Title{BoletoType: "NSA"}, Expected: "99"},
	{Input: models.Title{BoletoType: "BDP"}, Expected: "99"},
}

func TestProcessBoleto_WhenServiceRespondsSuccessfully_ShouldHasSuccessfulBoletoResponse(t *testing.T) {
	mock.StartMockService("9093")

	input := newStubBoletoRequestCaixa().Build()
	bank := New()

	output, _ := bank.ProcessBoleto(input)

	test.AssertProcessBoletoWithSuccess(t, output)
}

func TestProcessBoleto_WhenServiceRespondsFailed_ShouldHasFailedBoletoResponse(t *testing.T) {
	mock.StartMockService("9092")

	input := newStubBoletoRequestCaixa().withAmountIsCents(400).Build()
	bank := New()

	output, _ := bank.ProcessBoleto(input)

	test.AssertProcessBoletoFailed(t, output)
}

func TestProcessBoleto_WhenRequestContainsInvalidOurNumberParameter_ShouldHasFailedBoletoResponse(t *testing.T) {
	largeOurNumber := uint(9999999999999999)
	mock.StartMockService("9092")
	input := newStubBoletoRequestCaixa().withOurNumber(largeOurNumber).Build()

	bank := New()

	output, _ := bank.ProcessBoleto(input)

	test.AssertProcessBoletoFailed(t, output)
}

func TestGetCaixaCheckSumInfo(t *testing.T) {
	const expectedSumCode = "0200656000000000000000003008201700000000000100000732159000109"
	const expectedToken = "LvWr1op5Ayibn6jsCQ3/2bW4KwThVAlLK5ftxABlq20="

	bank := New()

	agreement := uint(200656)
	expiredAt := time.Date(2017, 8, 30, 12, 12, 12, 12, time.Local)
	doc := "00732159000109"

	input := newStubBoletoRequestCaixa().withAgreementNumber(agreement).withOurNumber(0).withAmountIsCents(1000).withExpirationDate(expiredAt).withRecipientDocumentNumber(doc).Build()

	assert.Equal(t, expectedSumCode, bank.getCheckSumCode(*input), "Deve-se formar uma string seguindo o padrão da documentação")
	assert.Equal(t, expectedToken, bank.getAuthToken(bank.getCheckSumCode(*input)), "Deve-se fazer um hash sha256 e encodar com base64")
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

func TestTemplateRequestCaixa_WhenRequestV1_ParseSuccessful(t *testing.T) {
	f := flow.NewFlow()
	input := newStubBoletoRequestCaixa().Build()

	b := fmt.Sprintf("%v", f.From("message://?source=inline", input, getRequestCaixa(), tmpl.GetFuncMaps()).GetBody())

	for _, expected := range expectedBasicTitleRequestFields {
		assert.Contains(t, b, expected, "Erro no mapeamento dos campos básicos do Título")
	}

	for _, expected := range expectedBuyerRequestFields {
		assert.Contains(t, b, expected, "Erro no mapeamento dos campos básicos do Comprador")
	}

	for _, notExpected := range expectedStrictRulesFieldsV2 {
		assert.NotContains(t, b, notExpected, "Não devem haver campos de regras de pagamento na V1")
	}

	for _, notExpected := range expectedFlexRulesFieldsV2 {
		assert.NotContains(t, b, notExpected, "Não devem haver campos de regras de pagamento na V1")
	}
}

func TestTemplateRequestCaixa_WhenRequestWithStrictRulesV2_ParseSuccessful(t *testing.T) {
	f := flow.NewFlow()
	input := newStubBoletoRequestCaixa().withStrictRules().Build()

	b := fmt.Sprintf("%v", f.From("message://?source=inline", input, getRequestCaixa(), tmpl.GetFuncMaps()).GetBody())

	for _, expected := range expectedBasicTitleRequestFields {
		assert.Contains(t, b, expected, "Erro no mapeamento dos campos básicos do Título")
	}

	for _, expected := range expectedBuyerRequestFields {
		assert.Contains(t, b, expected, "Erro no mapeamento dos campos básicos do Comprador")
	}

	for _, expected := range expectedStrictRulesFieldsV2 {
		assert.Contains(t, b, expected, "Erro no mapeamento das regras de pagamento")
	}
}

func TestTemplateRequestCaixa_WhenRequestWithFlexRulesV2_ParseSuccessful(t *testing.T) {
	f := flow.NewFlow()
	input := newStubBoletoRequestCaixa().withFlexRules().Build()

	b := fmt.Sprintf("%v", f.From("message://?source=inline", input, getRequestCaixa(), tmpl.GetFuncMaps()).GetBody())

	for _, expected := range expectedBasicTitleRequestFields {
		assert.Contains(t, b, expected, "Erro no mapeamento dos campos básicos do Título")
	}

	for _, expected := range expectedBuyerRequestFields {
		assert.Contains(t, b, expected, "Erro no mapeamento dos campos básicos do Comprador")
	}

	for _, expected := range expectedFlexRulesFieldsV2 {
		assert.Contains(t, b, expected, "Erro no mapeamento das regras de pagamento")
	}
}
