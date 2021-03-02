package models

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

type ModelTestParameter struct {
	Input    interface{}
	Expected interface{}
	Length   interface{}
}

var isCPFParameters = []ModelTestParameter{
	{Input: Document{Type: "CPF", Number: "13245678901"}, Expected: true},
	{Input: Document{Type: "AAA", Number: "13245678901"}, Expected: false},
	{Input: Document{Type: "cPf", Number: "13245678901ssa"}, Expected: true},
	{Input: Document{Type: "CPF", Number: "lasjdlf019239098adjal9390jflsadjf9309jfsl"}, Expected: true},
}

var validateCPFParameters = []ModelTestParameter{
	{Input: Document{Type: "CPF", Number: "13245678901"}, Expected: true},
	{Input: Document{Type: "CPF", Number: "lasjdlf019239098adjal9390jflsadjf9309jfsl"}, Expected: false},
}

var isCNPJParameters = []ModelTestParameter{
	{Input: Document{Type: "CNPJ", Number: "12123123000112"}, Expected: true},
	{Input: Document{Type: "AAA", Number: "12123123000112"}, Expected: false},
	{Input: Document{Type: "cNpJ", Number: "13245678901ssa"}, Expected: true},
	{Input: Document{Type: "CNPJ", Number: "lasjdlf019239098adjal9390jflsadjf9309jfsl"}, Expected: true},
}

var validateCNPJParameters = []ModelTestParameter{
	{Input: Document{Type: "CNPJ", Number: "12123123000112"}, Expected: true},
	{Input: Document{Type: "CNPJ", Number: "lasjdlf019239098adjal9390jflsadjf9309jfsl"}, Expected: false},
}

var titleDocumentNumberValidParameters = []ModelTestParameter{
	{Input: Title{DocumentNumber: "1234567891011"}, Expected: true},
	{Input: Title{DocumentNumber: "123x"}, Expected: true},
}

var titleDocumentNumberInvalidParameters = []ModelTestParameter{
	{Input: Title{DocumentNumber: "xx"}, Expected: false},
	{Input: Title{DocumentNumber: ""}, Expected: false},
}

var titleInstructionParameters = []ModelTestParameter{
	{Input: Title{Instructions: "Some instructions"}, Length: 100, Expected: true},
	{Input: Title{Instructions: "Some instructions"}, Length: 1, Expected: false},
}

var titleAmountInCentsParameters = []ModelTestParameter{
	{Input: Title{AmountInCents: 100}, Expected: true},
	{Input: Title{AmountInCents: 0}, Expected: false},
}

var titleParseDateParameters = []ModelTestParameter{
	{Input: "2021-03-05", Expected: true},
	{Input: "05/03/2021", Expected: false},
}

var titleExpirationDateParameters = []ModelTestParameter{
	{Input: Title{ExpireDate: time.Now().AddDate(0, 0, 5).Format("2006-01-02")}, Expected: true},
	{Input: Title{ExpireDate: time.Now().AddDate(0, 0, -5).Format("2006-01-02")}, Expected: false},
	{Input: Title{ExpireDate: "05/03/2021"}, Expected: false},
}

var agreementAgencyParameters = []ModelTestParameter{
	{Input: Agreement{Agency: "321"}, Expected: true},
	{Input: Agreement{Agency: "234-2222a"}, Expected: false},
}

var agreementAgencyDigitParameters = []ModelTestParameter{
	{Input: Agreement{AgencyDigit: "2sssss"}, Expected: "2"},
	{Input: Agreement{AgencyDigit: "332sssss"}, Expected: "1"},
}

var agreementAccountDigitParameters = []ModelTestParameter{
	{Input: Agreement{AccountDigit: "2sssss"}, Expected: "2"},
	{Input: Agreement{AccountDigit: "332sssss"}, Expected: "1"},
}

var agreementAccountParameters = []ModelTestParameter{
	{Input: Agreement{Account: "1234fff"}, Length: 8, Expected: true},
	{Input: Agreement{Account: "654654654654654654654654654564"}, Length: 8, Expected: false},
}

func TestIsCpf(t *testing.T) {
	for _, fact := range isCPFParameters {
		input := fact.Input.(Document)
		result := input.IsCPF()
		assert.Equal(t, fact.Expected, result, "Espera que o tipo de documento passado seja um CPF")
	}
}

func TestValidateCPF(t *testing.T) {
	for _, fact := range validateCPFParameters {
		input := fact.Input.(Document)
		result := input.ValidateCPF() == nil
		assert.Equal(t, fact.Expected, result, "Espera que CPF seja válido")
	}
}

func TestIsCnpj(t *testing.T) {
	for _, fact := range isCNPJParameters {
		input := fact.Input.(Document)
		result := input.IsCNPJ()
		assert.Equal(t, fact.Expected, result, "Espera que o tipo de documento passado seja um CNPJ")
	}
}

func TestValidateCNPJ(t *testing.T) {
	for _, fact := range validateCNPJParameters {
		input := fact.Input.(Document)
		result := input.ValidateCNPJ() == nil
		assert.Equal(t, fact.Expected, result, "Espera que CNPJ seja válido")
	}
}

func TestTitleValidateDocumentNumberSuccess(t *testing.T) {
	for _, fact := range titleDocumentNumberValidParameters {
		input := fact.Input.(Title)
		result := input.ValidateDocumentNumber() == nil
		assert.Equal(t, fact.Expected, result, "Espera que DocumentNumber seja válido")
		assert.Equal(t, len(input.DocumentNumber), 10)
	}
}

func TestTitleValidateDocumentNumberFailed(t *testing.T) {
	for _, fact := range titleDocumentNumberInvalidParameters {
		input := fact.Input.(Title)
		result := input.ValidateDocumentNumber() != nil
		assert.Equal(t, fact.Expected, result, "Espera que DocumentNumber seja inválido")
		assert.Equal(t, len(input.DocumentNumber), 0)
	}
}

func TestTitleValidateInstructions(t *testing.T) {
	for _, fact := range titleInstructionParameters {
		input := fact.Input.(Title)
		result := input.ValidateInstructionsLength(fact.Length.(int)) == nil
		assert.Equal(t, fact.Expected, result, "Deve validar corretamente as instruções")
	}
}

func TestTitleValidateAmountInCents(t *testing.T) {
	for _, fact := range titleAmountInCentsParameters {
		input := fact.Input.(Title)
		result := input.IsAmountInCentsValid() == nil
		assert.Equal(t, fact.Expected, result, "Deve validar corretamente o valor do boleto")
	}
}

func TestTitleParseDate(t *testing.T) {
	for _, fact := range titleParseDateParameters {
		_, err := parseDate(fact.Input.(string))
		result := err == nil
		assert.Equal(t, fact.Expected, result, "Deve transformar uma string no padrão para um tipo time.Time")
	}
}

func TestTitleExpireDate(t *testing.T) {
	for _, fact := range titleExpirationDateParameters {
		input := fact.Input.(Title)
		result := input.IsExpireDateValid() == nil
		assert.Equal(t, fact.Expected, result, "O ExpireDate deve ser válido")
	}
}

func TestNewBoletoView(t *testing.T) {
	expectedBarCode := "1234"
	expectedDigitableLine := "12345"

	response := BoletoResponse{
		BarCodeNumber: expectedBarCode,
		DigitableLine: expectedDigitableLine,
	}

	result := NewBoletoView(BoletoRequest{}, response, "BradescoShopFacil")

	assert.NotEmpty(t, result.UID)
	assert.Equal(t, expectedBarCode, result.Barcode)
	assert.Equal(t, expectedDigitableLine, result.DigitableLine)
}

func TestNewErrorResponse(t *testing.T) {
	result := NewErrors()
	assert.Empty(t, result, "Deve criar uma coleção de Errors vazia")
}

func TestAppendErrorResponse(t *testing.T) {
	result := NewErrors()

	result.Append("100", "FirstErrorMessage")

	assert.NotEmpty(t, result, "Deve incrementar a coleção com um item")
	assert.Equal(t, len(result), 1, "Deve conter um erro")

	result.Append("200", "SecondErrorMessage")

	assert.NotEmpty(t, result, "Deve incrementar a coleção com um item")
	assert.Equal(t, len(result), 2, "Deve conter dois erros")
}

func TestAgreementIsValidAgency(t *testing.T) {
	for _, fact := range agreementAgencyParameters {
		input := fact.Input.(Agreement)
		result := input.IsAgencyValid() == nil
		assert.Equal(t, fact.Expected, result, "Deve validar corretamente a agência")
	}
}

func TestAgreementCalculateAgencyDigit(t *testing.T) {
	c := func(s string) string {
		return "1"
	}
	for _, fact := range agreementAgencyDigitParameters {
		input := fact.Input.(Agreement)
		input.CalculateAgencyDigit(c)
		assert.Equal(t, fact.Expected, input.AgencyDigit, "Deve calcular corretamente o dígito da agência")
	}
}

func TestAgreementCalculateAccountDigit(t *testing.T) {
	c := func(s, y string) string {
		return "1"
	}
	for _, fact := range agreementAccountDigitParameters {
		input := fact.Input.(Agreement)
		input.CalculateAccountDigit(c)
		assert.Equal(t, fact.Expected, input.AccountDigit, "Deve calcular corretamente o dígito da conta")
	}
}

func TestAgreementIsAccountValid(t *testing.T) {
	for _, fact := range agreementAccountParameters {
		input := fact.Input.(Agreement)
		result := input.IsAccountValid(fact.Length.(int)) == nil
		assert.Equal(t, fact.Expected, result, "Verifica se a conta é valida")
	}
}
