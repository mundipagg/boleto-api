package citibank

import (
	"testing"

	"github.com/mundipagg/boleto-api/models"
	. "github.com/smartystreets/goconvey/convey"
)

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
