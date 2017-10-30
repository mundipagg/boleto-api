package integrationTests

import (
	"encoding/json"
	"strings"
	"testing"

	"github.com/mundipagg/boleto-api/app"
	"github.com/mundipagg/boleto-api/models"
	"github.com/mundipagg/boleto-api/util"
	. "github.com/smartystreets/goconvey/convey"
)

func TestShouldRegisterBoletoSantander(t *testing.T) {
	param := app.NewParams()
	param.DevMode = true
	param.DisableLog = true
	param.MockMode = true
	go app.Run(param)
	Convey("Deve-se registrar um boleto no Santander", t, func() {
		response, st, err := util.Post("http://localhost:3000/v1/boleto/register", getBody(models.Santander, 200), nil)
		So(err, ShouldEqual, nil)
		So(st, ShouldEqual, 200)
		boleto := models.BoletoResponse{}
		errJSON := json.Unmarshal([]byte(response), &boleto)
		So(errJSON, ShouldEqual, nil)
		Convey("Se o boleto foi registrado então ele tem que está disponível no formato HTML", func() {
			So(len(boleto.Links), ShouldBeGreaterThan, 0)
			html, st, err := util.Get(boleto.Links[0].Href, "", nil)
			So(err, ShouldEqual, nil)
			So(st, ShouldEqual, 200)
			htmlFromBoleto := strings.Contains(html, boleto.DigitableLine)
			So(htmlFromBoleto, ShouldBeTrue)
		})
	})

	Convey("Deve-se retornar uma falha ao registrar boleto no Santander", t, func() {
		response, st, err := util.Post("http://localhost:3000/v1/boleto/register", getBody(models.Santander, 3), nil)
		So(err, ShouldEqual, nil)
		So(st, ShouldEqual, 400)
		boleto := models.BoletoResponse{}
		errJSON := json.Unmarshal([]byte(response), &boleto)
		So(errJSON, ShouldEqual, nil)
		So(len(boleto.Errors), ShouldBeGreaterThan, 0)
	})

	Convey("Deve-se retornar uma falha ao registrar boleto devido ao código do convenio invalido no Santander", t, func() {
		req := getModelBody(models.Santander, 200)
		req.Agreement.AgreementNumber = 0
		response, st, err := util.Post("http://localhost:3000/v1/boleto/register", stringify(req), nil)
		So(err, ShouldEqual, nil)
		So(st, ShouldEqual, 400)
		boleto := models.BoletoResponse{}
		errJSON := json.Unmarshal([]byte(response), &boleto)
		So(errJSON, ShouldEqual, nil)
		So(len(boleto.Errors), ShouldBeGreaterThan, 0)
		So(boleto.Errors[0].Message, ShouldEqual, "O código do convênio deve ser preenchido")
	})
}
