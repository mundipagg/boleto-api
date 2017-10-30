package tmpl

import (
	"testing"

	"github.com/mundipagg/boleto-api/models"

	. "github.com/smartystreets/goconvey/convey"
)

func TestShouldPadLeft(t *testing.T) {
	Convey("O texto deve ter zeros a esqueda e até 5 caracteres", t, func() {
		s := padLeft("5", "0", 5)
		So(len(s), ShouldEqual, 5)
		So(s, ShouldEqual, "00005")
	})
}

func TestShouldReturnString(t *testing.T) {
	Convey("O número deve ser uma string", t, func() {
		So(toString(5), ShouldEqual, "5")
	})
}
func TestFormatDigitableLine(t *testing.T) {
	Convey("A linha digitável deve ser formatada corretamente", t, func() {
		s := "34191123456789010111213141516171812345678901112"
		So(fmtDigitableLine(s), ShouldEqual, "34191.12345 67890.101112 13141.516171 8 12345678901112")
	})
}

func TestTruncate(t *testing.T) {
	Convey("Deve-se truncar uma string", t, func() {
		s := "00000000000000000000"
		So(truncateString(s, 5), ShouldEqual, "00000")
		So(truncateString(s, 50), ShouldEqual, "00000000000000000000")
		So(truncateString("", 50), ShouldEqual, "")
	})
}

func TestClearString(t *testing.T) {
	Convey("Deve-se limpar uma string", t, func() {
		So(clearString("óláçñê"), ShouldEqual, "olacne")
		So(clearString("ola"), ShouldEqual, "ola")
		So(clearString(""), ShouldEqual, "")
		So(clearString("Jardim Novo Cambuí "), ShouldEqual, "Jardim Novo Cambui")
		So(clearString("Jardim Novo Cambuí�"), ShouldEqual, "Jardim Novo Cambui")
	})
}

func TestJoinStringSpace(t *testing.T) {
	Convey("Deve-se fazer um join em uma string com espaços", t, func() {
		So(joinSpace("a", "b", "c"), ShouldEqual, "a b c")
	})
}

func TestFormatCNPJ(t *testing.T) {
	Convey("O CNPJ deve ser formatado corretamente", t, func() {
		s := "01000000000100"
		So(fmtCNPJ(s), ShouldEqual, "01.000.000/0001-00")
	})
}

func TestFormatCPF(t *testing.T) {
	Convey("O CPF deve ser formatado corretamente", t, func() {
		s := "12312100100"
		So(fmtCPF(s), ShouldEqual, "123.121.001-00")
	})
}

func TestFormatNumber(t *testing.T) {
	Convey("O valor em inteiro deve ser convertido para uma string com duas casas decimais separado por vírgula (0,00)", t, func() {
		So(fmtNumber(50332), ShouldEqual, "503,32")
		So(fmtNumber(55), ShouldEqual, "0,55")
		So(fmtNumber(0), ShouldEqual, "0,00")
	})
}

func TestMod11OurNumber(t *testing.T) {
	Convey("Deve-se calcular o mod11 do nosso número e retornar o digito à esquerda", t, func() {
		So(calculateOurNumberMod11(12000000114), ShouldEqual, 120000001148)
		So(calculateOurNumberMod11(8423657), ShouldEqual, 84236574)
	})
}

func TestToFloatStr(t *testing.T) {
	Convey("O valor em inteiro deve ser convertido para uma string com duas casas decimais separado por ponto (0.00)", t, func() {
		So(toFloatStr(50332), ShouldEqual, "503.32")
		So(toFloatStr(55), ShouldEqual, "0.55")
		So(toFloatStr(0), ShouldEqual, "0.00")
	})
}

func TestFormatDoc(t *testing.T) {
	Convey("O CPF deve ser formatado corretamente", t, func() {
		d := models.Document{
			Type:   "CPF",
			Number: "12312100100",
		}
		So(fmtDoc(d), ShouldEqual, "123.121.001-00")
		Convey("O CNPJ deve ser formatado corretamente", func() {
			d.Type = "CNPJ"
			d.Number = "01000000000100"
			So(fmtDoc(d), ShouldEqual, "01.000.000/0001-00")
		})
	})
}

func TestDocType(t *testing.T) {
	Convey("O tipo retornardo deve ser CPF", t, func() {
		d := models.Document{
			Type:   "CPF",
			Number: "12312100100",
		}
		So(docType(d), ShouldEqual, 1)
		Convey("O tipo retornardo deve ser CNPJ", func() {
			d.Type = "CNPJ"
			d.Number = "01000000000100"
			So(docType(d), ShouldEqual, 2)
		})
	})
}

func TestTrim(t *testing.T) {
	Convey("O texto não deve ter espaços no início e no final", t, func() {
		d := " hue br festa "
		So(trim(d), ShouldEqual, "hue br festa")
	})
}
