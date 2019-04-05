package itau

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
	"BankNumber": 341,
	"Authentication": {
		"Username": "a",
		"Password": "b",
		"AccessKey":"c"
	},
	"Agreement": {
		"Wallet":109,
		"Agency":"0407",
		"Account":"55292",
		"AccountDigit":"6"
	},
	"Title": {
		"ExpireDate": "2999-12-31",
		"AmountInCents": 200			
	},
	"Buyer": {
		"Name": "TESTE",
		"Document": {
			"Type": "CNPJ",
			"Number": "00001234567890"
		}
	},
	"Recipient": {
		"Name": "TESTE",
		"Document": {
			"Type": "CNPJ",
			"Number": "00123456789067"
		}
	}
}
`

func TestRegiterBoletoItau(t *testing.T) {
	env.Config(true, true, true)
	input := new(models.BoletoRequest)
	if err := util.FromJSON(baseMockJSON, input); err != nil {
		t.Fail()
	}
	bank := New()
	go mock.Run("9096")
	time.Sleep(2 * time.Second)
	Convey("deve-se processar um boleto itau com sucesso", t, func() {
		output, err := bank.ProcessBoleto(input)
		So(err, ShouldBeNil)
		So(output.BarCodeNumber, ShouldNotBeEmpty)
		So(output.DigitableLine, ShouldNotBeEmpty)
		So(output.Errors, ShouldBeEmpty)
	})
	input.Title.AmountInCents = 400
	Convey("deve-se processar uma falha no registro de boleto no itau", t, func() {
		output, err := bank.ProcessBoleto(input)
		So(err, ShouldBeNil)
		So(output.BarCodeNumber, ShouldBeEmpty)
		So(output.DigitableLine, ShouldBeEmpty)
		So(output.Errors, ShouldNotBeEmpty)
	})
	input.Title.AmountInCents = 200
	ac := input.Agreement.Account
	input.Agreement.Account = ""
	Convey("deve-se tratar uma validacao de conta no itau", t, func() {
		output, err := bank.ProcessBoleto(input)
		So(err, ShouldBeNil)
		So(output.BarCodeNumber, ShouldBeEmpty)
		So(output.DigitableLine, ShouldBeEmpty)
		So(output.Errors, ShouldNotBeEmpty)
	})
	input.Agreement.Account = ac
	input.Authentication.Username = ""
	Convey("deve-se tratar uma falha de login no itau", t, func() {
		output, err := bank.ProcessBoleto(input)
		So(err, ShouldNotBeNil)
		So(output.BarCodeNumber, ShouldBeEmpty)
		So(output.DigitableLine, ShouldBeEmpty)
	})

	input.Title.BoletoType = "BDP"
	Convey("deve-se mapear corretamente o BoletoType de boleto de proposta", t, func() {
		_, output := getBoletoType(input)
		So(output, ShouldEqual, "18")
	})

	input.Title.BoletoType = "ITAU"
	Convey("deve-se mapear corretamente o BoletoType quando valor enviado não existir", t, func() {
		_, output := getBoletoType(input)
		So(output, ShouldEqual, "01")
	})
}

func TestGetBoletoType(t *testing.T) {

	input := new(models.BoletoRequest)
	if err := util.FromJSON(baseMockJSON, input); err != nil {
		t.Fail()
	}

	input.Title.BoletoType = ""
	expectBoletoTypeCode := "01"

	Convey("Quando não informado o BoletoType o retorno deve ser 01 - Duplicata Mercantil", t, func() {
		_, output := getBoletoType(input)
		So(output, ShouldEqual, expectBoletoTypeCode)
	})

	input.Title.BoletoType = "NSA"
	expectBoletoTypeCode = "01"

	Convey("Quando informado o BoletoType Inválido o retorno deve ser 01 - Duplicata Mercantil", t, func() {
		_, output := getBoletoType(input)
		So(output, ShouldEqual, expectBoletoTypeCode)
	})

	input.Title.BoletoType = "BDP"
	expectBoletoTypeCode = "18"

	Convey("Quando informado o BoletoType BDP o retorno deve ser 18 - Boleto de Proposta", t, func() {
		_, output := getBoletoType(input)
		So(output, ShouldEqual, expectBoletoTypeCode)
	})

}

func TestShouldMapItauBoletoType(t *testing.T) {
	env.Config(true, true, true)
	input := new(models.BoletoRequest)
	if err := util.FromJSON(baseMockJSON, input); err != nil {
		t.Fail()
	}

	go mock.Run("9097")
	time.Sleep(2 * time.Second)

	Convey("deve-se mapear corretamente o BoletoType quando informação for vazia", t, func() {
		_, output := getBoletoType(input)
		So(input.Title.BoletoType, ShouldEqual, "")
		So(output, ShouldEqual, "01")
	})

	input.Title.BoletoType = "BDP"
	Convey("deve-se mapear corretamente o BoletoType de boleto de proposta", t, func() {
		_, output := getBoletoType(input)
		So(output, ShouldEqual, "18")
	})

	input.Title.BoletoType = "ITAU"
	Convey("deve-se mapear corretamente o BoletoType quando valor enviado não existir", t, func() {
		_, output := getBoletoType(input)
		So(output, ShouldEqual, "01")
	})
}
