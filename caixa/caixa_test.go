package caixa

import (
	"testing"

	"github.com/mundipagg/boleto-api/models"

	"time"

	. "github.com/smartystreets/goconvey/convey"
)

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
