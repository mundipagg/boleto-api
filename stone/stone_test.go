package stone

import (
	"github.com/mundipagg/boleto-api/env"
	"github.com/mundipagg/boleto-api/mock"
	"github.com/mundipagg/boleto-api/models"
	"github.com/mundipagg/boleto-api/util"
	. "github.com/smartystreets/goconvey/convey"
)

import (
	"testing"
	"time"
)

const baseMockJSON = `
{
    "Agreement":{
        "AccountID":"2dd5ea76-80d2-41be-a991-fe30713ab7ed"
    },
    "Title":{
      "ExpireDate": "2020-06-30",
      "AmountInCents":5000,
      "BoletoType":"proposal"
    },
    "Buyer":{
        "Name":"Cliente de teste Stone",
				"TradeName": "Trade Name teste",
        "Document": {
            "Type":"CPF",
            "Number":"96558156415"
        }
    },
    "BankNumber":197
}
`

func TestRegiterBoletoStone(t *testing.T) {
	env.Config(true, true, true)
	input := new(models.BoletoRequest)
	if err := util.FromJSON(baseMockJSON, input); err != nil {
		t.Fail()
	}

	bank := New()

	go mock.Run("9096")
	time.Sleep(2 * time.Second)

	Convey("deve-se processar um boleto stone com sucesso", t, func() {
		output, err := bank.ProcessBoleto(input)
		So(err, ShouldBeNil)
		So(output.BarCodeNumber, ShouldNotBeEmpty)
		So(output.DigitableLine, ShouldNotBeEmpty)
		So(output.Errors, ShouldBeEmpty)
	})

	Convey("Deve-se exibir uma mensagem de erro, caso o registro não aconteça com sucesso", t, func() {
		input.Title.AmountInCents = 10
		output, err := bank.ProcessBoleto(input)
		So(err, ShouldBeNil)
		So(output.BarCodeNumber, ShouldBeEmpty)
		So(output.DigitableLine, ShouldBeEmpty)
		So(output.Errors, ShouldNotBeEmpty)
	})
}



