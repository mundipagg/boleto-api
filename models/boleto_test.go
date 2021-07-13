package models

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

var isBankNumberValidParameters = []ModelTestParameter{
	{Input: BoletoRequest{BankNumber: BancoDoBrasil}, Expected: true},
	{Input: BoletoRequest{BankNumber: Bradesco}, Expected: true},
	{Input: BoletoRequest{BankNumber: Caixa}, Expected: true},
	{Input: BoletoRequest{BankNumber: Citibank}, Expected: true},
	{Input: BoletoRequest{BankNumber: Itau}, Expected: true},
	{Input: BoletoRequest{BankNumber: Pefisa}, Expected: true},
	{Input: BoletoRequest{BankNumber: Santander}, Expected: true},
	{Input: BoletoRequest{BankNumber: Stone}, Expected: true},
	{Input: BoletoRequest{BankNumber: 0}, Expected: false},
}

var bankNumberAndDigitParameters = []ModelTestParameter{
	{Input: BoletoRequest{BankNumber: BancoDoBrasil}, Expected: "001-9"},
	{Input: BoletoRequest{BankNumber: Bradesco}, Expected: "237-2"},
	{Input: BoletoRequest{BankNumber: Caixa}, Expected: "104-0"},
	{Input: BoletoRequest{BankNumber: Citibank}, Expected: "745-5"},
	{Input: BoletoRequest{BankNumber: Itau}, Expected: "341-7"},
	{Input: BoletoRequest{BankNumber: Pefisa}, Expected: "174"},
	{Input: BoletoRequest{BankNumber: Santander}, Expected: "033-7"},
	{Input: BoletoRequest{BankNumber: Stone}, Expected: "197-1"},
	{Input: BoletoRequest{BankNumber: 0}, Expected: ""},
}

func Test_IsBankNumberValid(t *testing.T) {
	for _, fact := range isBankNumberValidParameters {
		input := fact.Input.(BoletoRequest)
		result := input.BankNumber.IsBankNumberValid()
		assert.Equal(t, fact.Expected.(bool), result)
	}
}

func Test_GetBoletoBankNumberAndDigit(t *testing.T) {
	for _, fact := range bankNumberAndDigitParameters {
		input := fact.Input.(BoletoRequest)
		result := input.BankNumber.GetBoletoBankNumberAndDigit()
		assert.Equal(t, fact.Expected.(string), result)
	}
}

func Test_HasErrorResponse(t *testing.T) {
	withError := GetBoletoResponseError("CODE", "message")
	resultWithError := withError.HasErrors()
	assert.True(t, resultWithError)

	withoutError := BoletoResponse{}
	resultWithputError := withoutError.HasErrors()
	assert.False(t, resultWithputError)
}

func Test_ViewToJson(t *testing.T) {
	view := arrangeBoletoView()

	result := view.ToJSON()

	assert.Contains(t, result, `"bankNumber":"001-9"`)
	assert.Contains(t, result, `"barcode":"123456789012345678901234567890"`)
	assert.Contains(t, result, `"digitableLine":"123467890123456790134567890"`)
	assert.Contains(t, result, `"ourNumber":"1234567890"`)
}

func Test_ViewToMinifyJSON(t *testing.T) {
	view := arrangeBoletoView()

	result := view.ToMinifyJSON()

	assert.Contains(t, result, `"bankNumber":"001-9"`)
	assert.Contains(t, result, `"barcode":"123456789012345678901234567890"`)
	assert.Contains(t, result, `"digitableLine":"123467890123456790134567890"`)
	assert.Contains(t, result, `"ourNumber":"1234567890"`)
}

func arrangeBoletoView() BoletoView {
	request := BoletoRequest{BankNumber: BancoDoBrasil}
	response := BoletoResponse{BarCodeNumber: "123456789012345678901234567890", DigitableLine: "123467890123456790134567890", OurNumber: "1234567890"}

	return NewBoletoView(request, response, "BancoDoBrasil")
}
