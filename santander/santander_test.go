package santander

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
	"BankNumber": 33,
	"Agreement": {
		"AgreementNumber": 11111111,		
		"Agency":"5555",
		"Account":"55555"
	},
	"Title": {
		"ExpireDate": "2035-08-01",
		"AmountInCents": 200,
		"OurNumber":10000000004		
	},
	"Buyer": {
		"Name": "TESTE",
		"Document": {
			"Type": "CPF",
			"Number": "12345678903"
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

func TestShouldProcessBoletoSantander(t *testing.T) {
	env.Config(true, true, true)
	input := new(models.BoletoRequest)
	if err := util.FromJSON(baseMockJSON, input); err != nil {
		t.Fail()
	}
	bank := New()
	go mock.Run("9097")
	time.Sleep(2 * time.Second)
	Convey("deve-se processar um boleto santander com sucesso", t, func() {
		output, err := bank.ProcessBoleto(input)
		So(err, ShouldBeNil)
		So(output.BarCodeNumber, ShouldNotBeEmpty)
		So(output.DigitableLine, ShouldNotBeEmpty)
		So(output.Errors, ShouldBeEmpty)
	})

	input.Title.BoletoType = "BP"
	Convey("deve-se mapear corretamente o BoletoType de boleto de proposta", t, func() {
		output := bank.GetBoletoType(input)
		So(output, ShouldEqual, "32")
	})

	input.Title.BoletoType = "Santander"
	Convey("deve-se mapear corretamente o BoletoType quando valor enviado não existir", t, func() {
		output := bank.GetBoletoType(input)
		So(output, ShouldEqual, "02")
	})
}
