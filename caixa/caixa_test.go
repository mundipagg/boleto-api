package caixa

import (
	"testing"
	"time"

	"github.com/mundipagg/boleto-api/env"
	"github.com/mundipagg/boleto-api/mock"
	"github.com/mundipagg/boleto-api/models"
	"github.com/mundipagg/boleto-api/util"
	. "github.com/smartystreets/goconvey/convey"
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

func TestShouldProcessBoleto(t *testing.T) {
	env.Config(true, true, true)
	input := new(models.BoletoRequest)
	if err := util.FromJSON(baseMockJSON, input); err != nil {
		t.Fail()
	}
	bank := New()
	go mock.Run("9094")
	time.Sleep(2 * time.Second)
	Convey("deve-se processar um boleto Caixa com sucesso", t, func() {
		output, err := bank.ProcessBoleto(input)
		So(err, ShouldBeNil)
		So(output.BarCodeNumber, ShouldNotBeEmpty)
		So(output.DigitableLine, ShouldNotBeEmpty)
		So(output.Errors, ShouldBeEmpty)
	})
	input.Title.AmountInCents = 400
	Convey("deve-se tratar erro no boleto Caixa", t, func() {
		output, err := bank.ProcessBoleto(input)
		So(err, ShouldBeNil)
		So(output.BarCodeNumber, ShouldBeEmpty)
		So(output.DigitableLine, ShouldBeEmpty)
		So(output.Errors, ShouldNotBeEmpty)
	})

	input.Title.AmountInCents = 200
	input.Title.Instructions = "aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa"
	input.Title.OurNumber = 9999999999999999
	Convey("deve-se tratar erro no boleto Caixa", t, func() {
		output, err := bank.ProcessBoleto(input)
		So(err, ShouldBeNil)
		So(output.BarCodeNumber, ShouldBeEmpty)
		So(output.DigitableLine, ShouldBeEmpty)
		So(output.Errors, ShouldNotBeEmpty)
	})
}

func TestGetCaixaCheckSumInfo(t *testing.T) {
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
	Convey("Geração do token de autorização da Caixa", t, func() {
		Convey("Deve-se formar uma string seguindo o padrão da documentação", func() {
			So(caixa.getCheckSumCode(boleto), ShouldEqual, "0200656000000000000000003008201700000000000100000732159000109")
		})
		Convey("Deve-se fazer um hash sha256 e encodar com base64", func() {
			So(caixa.getAuthToken(caixa.getCheckSumCode(boleto)), ShouldEqual, "LvWr1op5Ayibn6jsCQ3/2bW4KwThVAlLK5ftxABlq20=")
		})
	})

}

func TestShouldCalculateAccountDigitCaixa(t *testing.T) {
	Convey("Deve-se calcular  e validar Agencia e Conta da Caixa", t, func() {
		boleto := models.BoletoRequest{
			Agreement: models.Agreement{
				Account: "100000448",
				Agency:  "2004",
			},
		}
		err := caixaValidateAccountAndDigit(&boleto)
		errAg := caixaValidateAgency(&boleto)
		So(err, ShouldBeNil)
		So(errAg, ShouldBeNil)
	})
}

func TestGetBoletoType(t *testing.T) {

	input := new(models.BoletoRequest)
	if err := util.FromJSON(baseMockJSON, input); err != nil {
		t.Fail()
	}

	input.Title.BoletoType = ""
	expectBoletoTypeCode := "99"

	Convey("Quando não informado o BoletoType o retorno deve ser 99 - Outros", t, func() {
		_, output := getBoletoType(input)
		So(output, ShouldEqual, expectBoletoTypeCode)
	})

	input.Title.BoletoType = "NSA"
	expectBoletoTypeCode = "99"

	Convey("Quando informado o BoletoType Inválido o retorno deve ser 99 - Outros", t, func() {
		_, output := getBoletoType(input)
		So(output, ShouldEqual, expectBoletoTypeCode)
	})

	input.Title.BoletoType = "BDP"
	expectBoletoTypeCode = "99"

	Convey("Quando informado o BoletoType BDP o retorno deve ser 99 - Outros", t, func() {
		_, output := getBoletoType(input)
		So(output, ShouldEqual, expectBoletoTypeCode)
	})
}
