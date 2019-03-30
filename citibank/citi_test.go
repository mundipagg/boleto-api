package citibank

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
    "BankNumber": 745,
    "Authentication": {
        "Username": "55555555555555555555"
    },
    "Agreement": {
        "AgreementNumber": 55555555,
        "Wallet" : 100,
        "Agency":"0011",
        "Account":"0088881323",
        "AccountDigit" : "2"        
    },
    "Title": {
        "ExpireDate": "2029-09-20",
        "AmountInCents": 200,
        "OurNumber": 10000000001,
        "DocumentNumber":"5555555555"
    },
    "Buyer": {
        "Name": "Fulano de Tal",
        "Document": {
            "Type": "CNPJ",
            "Number": "55555555555555"
        }
    },
    "Recipient": {
      "Name": "TESTE",
        "Document": {
            "Type": "CNPJ",
            "Number": "55555555555555"
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
	go mock.Run("9095")
	time.Sleep(2 * time.Second)
	output, err := bank.ProcessBoleto(input)
	Convey("deve-se processar um boleto citibank com sucesso", t, func() {
		So(err, ShouldBeNil)
		So(output.BarCodeNumber, ShouldNotBeEmpty)
		So(output.DigitableLine, ShouldNotBeEmpty)
		So(output.Errors, ShouldBeEmpty)
	})
}
func TestOurNumber(t *testing.T) {
	boleto := models.BoletoRequest{
		Title: models.Title{
			OurNumber: 8605970,
		},
	}
	Convey("deve-se calcular corretamente o nosso numero para o Citi", t, func() {
		So(calculateOurNumber(&boleto), ShouldEqual, 86059700)
	})
}

func TestGetBoletoType(t *testing.T) {

	input := new(models.BoletoRequest)
	if err := util.FromJSON(baseMockJSON, input); err != nil {
		t.Fail()
	}

	input.Title.BoletoType = ""
	expectBoletoTypeCode := "03"

	Convey("Quando não informado o BoletoType o retorno deve ser 03 - Duplicata Mercantil", t, func() {
		_, output := getBoletoType(input)
		So(output, ShouldEqual, expectBoletoTypeCode)
	})

	input.Title.BoletoType = "NSA"
	expectBoletoTypeCode = "03"

	Convey("Quando informado o BoletoType Inválido o retorno deve ser 03 - Duplicata Mercantil", t, func() {
		_, output := getBoletoType(input)
		So(output, ShouldEqual, expectBoletoTypeCode)
	})

	input.Title.BoletoType = "BDP"
	expectBoletoTypeCode = "03"

	Convey("Quando informado o BoletoType BDP o retorno deve ser 03 - Duplicata Mercantil", t, func() {
		_, output := getBoletoType(input)
		So(output, ShouldEqual, expectBoletoTypeCode)
	})
}
