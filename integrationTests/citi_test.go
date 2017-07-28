package integrationTests

import (
	"fmt"
	"strings"
	"testing"

	"github.com/mundipagg/boleto-api/app"
	"github.com/mundipagg/boleto-api/models"
	"github.com/mundipagg/boleto-api/util"
	. "github.com/smartystreets/goconvey/convey"
)

func TestRegisterBoletoCiti(t *testing.T) {
	param := app.NewParams()
	param.DevMode = true
	param.DisableLog = true
	param.HTTPOnly = true
	param.MockMode = true
	go app.Run(param)
	Convey("Deve-se registrar um boleto no Citi", t, func() {
		resp, st, err := util.Post("http://localhost:3000/v1/boleto/register", getBody(models.Citibank, 200), nil)
		boleto := util.ParseJSON(resp, new(models.BoletoResponse)).(*models.BoletoResponse)
		So(err, ShouldBeNil)
		So(st, ShouldEqual, 200)
		Convey("Deve-se gerar uma url de acesso ao boleto seguindo o padrão do Citibank", func() {
			So(len(boleto.Links), ShouldBeGreaterThan, 0)
			resp, st, err := util.Post("http://localhost:3000/v1/boleto/register", getBody(models.Citibank, 200), nil)
			So(err, ShouldBeNil)
			So(st, ShouldEqual, 200)
			savedBoleto := util.ParseJSON(resp, new(models.BoletoView)).(*models.BoletoView)
			So(strings.Contains(boleto.Links[0].Href, fmt.Sprintf("%d", savedBoleto.Boleto.Title.OurNumber)), ShouldBeTrue)
			So(strings.Contains(boleto.Links[0].Href, savedBoleto.Boleto.Recipient.Document.Number), ShouldBeTrue)
			So(strings.Contains(boleto.Links[0].Href, savedBoleto.Boleto.Buyer.Document.Number), ShouldBeTrue)
		})
	})
}
