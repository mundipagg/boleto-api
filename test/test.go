package test

import (
	"testing"

	"github.com/mundipagg/boleto-api/models"
	"github.com/stretchr/testify/assert"
)

//Parameter Parâmetro de teste com input generico
type Parameter struct {
	Input    interface{}
	Expected interface{}
	Length   int
}

//UInt64TestParameter Parâmetro de teste com input do tipo uint64
type UInt64TestParameter struct {
	Input    uint64
	Expected string
}

// ExpectNoError falha o teste se e != nil
func ExpectNoError(e error, t *testing.T) {
	if e != nil {
		t.Fail()
	}
}

// ExpectError falha o teste se e == nil
func ExpectError(e error, t *testing.T) {
	if e == nil {
		t.Fail()
	}
}

// ExpectTrue falha o teste caso a condição não seja verdadeira
func ExpectTrue(condition bool, t *testing.T) {
	if !condition {
		t.Fail()
	}
}

// ExpectFalse falha o teste caso a condição não seja falsa
func ExpectFalse(condition bool, t *testing.T) {
	if condition {
		t.Fail()
	}
}

// ExpectNil falha o teste caso obj seja diferente de nil
func ExpectNil(obj interface{}, t *testing.T) {
	if obj != nil {
		t.Fail()
	}
}

//AssertProcessBoletoWithSuccess Valida se o boleto foi gerado com sucesso
func AssertProcessBoletoWithSuccess(t *testing.T, response models.BoletoResponse) {
	assert.Empty(t, response.Errors, "Não deve ocorrer erros")
	assert.NotNil(t, response.BarCodeNumber, "Deve haver um Barcode")
	assert.NotNil(t, response.DigitableLine, "Deve haver uma linha digitável")
}

//AssertProcessBoletoFailed Valida se o houve um erro no processamento do boleto
func AssertProcessBoletoFailed(t *testing.T, response models.BoletoResponse) {
	assert.NotEmpty(t, response.Errors, "Devem ocorrer erros ")
	assert.Empty(t, response.BarCodeNumber, "Não deve haver um Barcode")
	assert.Empty(t, response.DigitableLine, "Não deve haver uma linha digitável")
}
