package bank

import (
	"testing"

	"github.com/mundipagg/boleto-api/bb"
	"github.com/mundipagg/boleto-api/bradescoNetEmpresa"
	"github.com/mundipagg/boleto-api/bradescoShopFacil"
	"github.com/mundipagg/boleto-api/caixa"
	"github.com/mundipagg/boleto-api/citibank"
	"github.com/mundipagg/boleto-api/itau"
	"github.com/mundipagg/boleto-api/models"
	"github.com/mundipagg/boleto-api/pefisa"
	"github.com/mundipagg/boleto-api/santander"
	"github.com/mundipagg/boleto-api/stone"

	"github.com/stretchr/testify/assert"
)

type dataTest struct {
	request    models.BoletoRequest
	bankNumber models.BankNumber
	bank       Bank
}

var bradescoNetEmpresaInstance = bradescoNetEmpresa.New()
var bradescoShopFacilInstance = bradescoShopFacil.New()
var bancoDoBrasilInstance = bb.New()
var citibankInstance, _ = citibank.New()
var santanderInstance, _ = santander.New()
var itauInstance = itau.New()
var caixaInstance = caixa.New()
var pefisaInstance = pefisa.New()
var stoneInstance = stone.New()

var getBankTestData = []dataTest{
	{models.BoletoRequest{BankNumber: models.Bradesco, Agreement: models.Agreement{Wallet: 9}}, models.Bradesco, bradescoNetEmpresaInstance},
	{models.BoletoRequest{BankNumber: models.Bradesco, Agreement: models.Agreement{Wallet: 25}}, models.Bradesco, bradescoShopFacilInstance},
	{models.BoletoRequest{BankNumber: models.BancoDoBrasil}, models.BancoDoBrasil, bancoDoBrasilInstance},
	{models.BoletoRequest{BankNumber: models.Citibank}, models.Citibank, citibankInstance},
	{models.BoletoRequest{BankNumber: models.Santander}, models.Santander, santanderInstance},
	{models.BoletoRequest{BankNumber: models.Itau}, models.Itau, itauInstance},
	{models.BoletoRequest{BankNumber: models.Caixa}, models.Caixa, caixaInstance},
	{models.BoletoRequest{BankNumber: models.Pefisa}, models.Pefisa, pefisaInstance},
	{models.BoletoRequest{BankNumber: models.Stone}, models.Stone, stoneInstance},
}

func TestShouldExecuteBankStrategy(t *testing.T) {
	for _, fact := range getBankTestData {
		bank, err := Get(fact.request)
		number := bank.GetBankNumber()

		assert.Nil(t, err, "Não deve haver erro")
		assert.True(t, number.IsBankNumberValid(), "Deve ser um banco válido")
		assert.Equal(t, fact.bankNumber, number, "Deve conter o bankNumber correto")
		assert.IsType(t, fact.bank, bank, "Deve ter instaciado o banco correto")
	}
}

func TestGetBank_WhenInvalidBank_ReturnError(t *testing.T) {
	request := models.BoletoRequest{BankNumber: 0}

	result, err := Get(request)

	assert.Nil(t, result, "Não deve haver banco")
	assert.NotNil(t, err, "Deve haver erro")
	assert.Equal(t, err.(models.ErrorResponse).Code, "MPBankNumber")
	assert.Equal(t, err.(models.ErrorResponse).Message, "Banco 0 não existe")
}
