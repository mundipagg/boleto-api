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
	param.MockMode = true
	go app.Run(param)
	Convey("Deve-se registrar um boleto no Citi", t, func() {
		boletoReq := getModelBody(models.Citibank, 200)
		boletoReq.Agreement.Account = "0887214891"
		boletoReq.Agreement.Wallet = 100
		req := util.Stringify(boletoReq)
		resp, st, err := util.Post("http://localhost:3000/v1/boleto/register", req, nil)
		boletoResp := util.ParseJSON(resp, new(models.BoletoResponse)).(*models.BoletoResponse)
		So(err, ShouldBeNil)
		So(st, ShouldEqual, 200)
		savedBoleto := util.ParseJSON(resp, new(models.BoletoView)).(*models.BoletoView)
		So(strings.Contains(boletoResp.Links[0].Href, fmt.Sprintf("%d", savedBoleto.Boleto.Title.OurNumber)), ShouldBeTrue)
		So(strings.Contains(boletoResp.Links[0].Href, savedBoleto.Boleto.Recipient.Document.Number), ShouldBeTrue)
		So(strings.Contains(boletoResp.Links[0].Href, savedBoleto.Boleto.Buyer.Document.Number), ShouldBeTrue)
	})
	Convey("Deve-se validar os campos do agreement", t, func() {
		boletoReq := getModelBody(models.Citibank, 200)
		boletoReq.Agreement.Account = "0887214811111"
		boletoReq.Agreement.Wallet = 10
		req := util.Stringify(boletoReq)
		resp, st, err := util.Post("http://localhost:3000/v1/boleto/register", req, nil)
		boletoResp := util.ParseJSON(resp, new(models.BoletoResponse)).(*models.BoletoResponse)
		So(boletoResp, ShouldNotBeNil)
		So(err, ShouldBeNil)
		So(st, ShouldEqual, 400)
	})
}
