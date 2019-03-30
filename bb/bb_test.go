package bb

import (
	"testing"
	"time"

	"github.com/mundipagg/boleto-api/env"
	"github.com/mundipagg/boleto-api/mock"
	"github.com/mundipagg/boleto-api/models"
	"github.com/mundipagg/boleto-api/test"
	"github.com/mundipagg/boleto-api/util"
	. "github.com/smartystreets/goconvey/convey"
)

const baseMockJSON = `
{
	
	"BankNumber": 1,
	"Authentication": {
		"Username": "xxx",
		"Password": "xxx"
	},
	"Agreement": {
		"AgreementNumber": 5555555,
		"WalletVariation": 19,
		"Wallet":17,
		"Agency":"5555",
		"Account":"55555"
	},
	"Title": {
		"ExpireDate": "2029-10-01",
		"AmountInCents": 200,
		"OurNumber": 1,
		"Instructions": "Senhor caixa, não receber após o vencimento",
		"DocumentNumber": "123456"
	},
	"Buyer": {
		"Name": "TESTE",
		"Document": {
			"Type": "CNPJ",
			"Number": "55555555550140"
		},
		"Address": {
			"Street": "Teste",
			"Number": "123",
			"Complement": "Apto",
			"ZipCode": "55555555",
			"City": "Rio de Janeiro",
			"District": "Teste",
			"StateCode": "RJ"
		}
	},
	"Recipient": {
		"Name": "TESTE",
		"Document": {
			"Type": "CNPJ",
			"Number": "55555555555555"
		},
		"Address": {
			"Street": "TESTE",
			"Number": "555",
			"Complement": "Teste",
			"ZipCode": "0455555",
			"City": "São Paulo",
			"District": "",
			"StateCode": "SP"
		}
	}
}
`

func TestRegiterBoleto(t *testing.T) {
	env.Config(true, true, true)
	input := new(models.BoletoRequest)
	if err := util.FromJSON(baseMockJSON, input); err != nil {
		t.Fail()
	}
	bank := New()
	go mock.Run("9092")
	time.Sleep(2 * time.Second)
	Convey("deve-se processar um boleto BB com sucesso", t, func() {
		output, err := bank.ProcessBoleto(input)
		So(err, ShouldBeNil)
		So(output.BarCodeNumber, ShouldNotBeEmpty)
		So(output.DigitableLine, ShouldNotBeEmpty)
		So(output.Errors, ShouldBeEmpty)
	})
	input.Title.AmountInCents = 400
	Convey("deve-se tratar um boleto BB com erro", t, func() {
		output, err := bank.ProcessBoleto(input)
		So(err, ShouldBeNil)
		So(output.BarCodeNumber, ShouldBeEmpty)
		So(output.DigitableLine, ShouldBeEmpty)
		So(output.Errors, ShouldNotBeEmpty)
	})
	input.Title.AmountInCents = 200
	input.Agreement.Account = ""
	Convey("deve-se tratar um boleto BB com erro na conta", t, func() {
		output, err := bank.ProcessBoleto(input)
		So(err, ShouldBeNil)
		So(output.BarCodeNumber, ShouldBeEmpty)
		So(output.DigitableLine, ShouldBeEmpty)
		So(output.Errors, ShouldNotBeEmpty)
	})
}

func TestShouldCalculateAgencyDigitFromBb(t *testing.T) {
	test.ExpectTrue(bbAgencyDigitCalculator("0137") == "6", t)

	test.ExpectTrue(bbAgencyDigitCalculator("3418") == "5", t)

	test.ExpectTrue(bbAgencyDigitCalculator("3324") == "3", t)

	test.ExpectTrue(bbAgencyDigitCalculator("5797") == "5", t)
}

func TestShouldCalculateAccountDigitFromBb(t *testing.T) {
	test.ExpectTrue(bbAccountDigitCalculator("", "00006685") == "0", t)

	test.ExpectTrue(bbAccountDigitCalculator("", "00025619") == "6", t)

	test.ExpectTrue(bbAccountDigitCalculator("", "00006842") == "X", t)

	test.ExpectTrue(bbAccountDigitCalculator("", "00000787") == "0", t)
}

func TestGetBoletoType(t *testing.T) {

	input := new(models.BoletoRequest)
	if err := util.FromJSON(baseMockJSON, input); err != nil {
		t.Fail()
	}

	input.Title.BoletoType = ""
	expectBoletoTypeCode := "19"

	Convey("Quando não informado o BoletoType o padrão deve ser 19 - Nota de Débito", t, func() {
		_, output := getBoletoType(input)
		So(output, ShouldEqual, expectBoletoTypeCode)
	})

	input.Title.BoletoType = "NSA"
	expectBoletoTypeCode = "19"

	Convey("Quando informado o BoletoType Invalido o padrão deve ser 19 - Nota de Débito", t, func() {
		_, output := getBoletoType(input)
		So(output, ShouldEqual, expectBoletoTypeCode)
	})

}
