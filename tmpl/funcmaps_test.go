package tmpl

import (
	"html/template"
	"testing"

	"github.com/mundipagg/boleto-api/models"
	"github.com/mundipagg/boleto-api/test"
	"github.com/stretchr/testify/assert"
)

var formatDigitableLineParameters = []test.Parameter{
	{Input: "34191123456789010111213141516171812345678901112", Expected: "34191.12345 67890.101112 13141.516171 8 12345678901112"},
}

var truncateParameters = []test.Parameter{
	{Input: "00000000000000000000", Length: 5, Expected: "00000"},
	{Input: "00000000000000000000", Length: 50, Expected: "00000000000000000000"},
	{Input: "Rua de teste para o truncate", Length: 20, Expected: "Rua de teste para o "},
	{Input: "", Length: 50, Expected: ""},
}

var clearStringParameters = []test.Parameter{
	{Input: "óláçñê", Expected: "olacne"},
	{Input: "ola", Expected: "ola"},
	{Input: "", Expected: ""},
	{Input: "Jardim Novo Cambuí ", Expected: "Jardim Novo Cambui"},
	{Input: "Jardim Novo Cambuí�", Expected: "Jardim Novo Cambui"},
	{Input: "CaixaAccepted:ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789,/()*&=-+!:?<>.;_\"", Expected: "CaixaAccepted:ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789,/()*&=-+!:?<>.;_\""},
	{Input: "CaixaNotAccepted:ÁÉÍÓÚÀÈÌÒÙÂÊÎÔÛÃÕáéíóúàèìòùâêîôûãõç", Expected: "CaixaNotAccepted:AEIOUAEIOUAEIOUAOaeiouaeiouaeiouaoc"},
	{Input: "{|}", Expected: ""},
}

var clearStringCaixaParameters = []test.Parameter{
	{Input: "CaixaAccepted:ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789,/()*&=-+!:?<>.;_\"", Expected: "CaixaAccepted:ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789,/()*&=-+!:?<>.;_\""},
	{Input: "CaixaClearCharacter:ÁÉÍÓÚÀÈÌÒÙÂÊÎÔÛÃÕáéíóúàèìòùâêîôûãõç", Expected: "CaixaClearCharacter:AEIOUAEIOUAEIOUAOaeiouaeiouaeiouaoc"},
	{Input: "@#$%¨{}[]^~|ºª§°¹²³£¢¬\\�", Expected: "                        "},
}

var formatNumberParameters = []test.UInt64TestParameter{
	{Input: 50332, Expected: "503,32"},
	{Input: 55, Expected: "0,55"},
	{Input: 0, Expected: "0,00"},
}

var toFloatStrParameters = []test.UInt64TestParameter{
	{Input: 50332, Expected: "503.32"},
	{Input: 55, Expected: "0.55"},
	{Input: 0, Expected: "0.00"},
}

var formatDocParameters = []test.Parameter{
	{Input: models.Document{Type: "CPF", Number: "12312100100"}, Expected: "123.121.001-00"},
	{Input: models.Document{Type: "CNPJ", Number: "12123123000112"}, Expected: "12.123.123/0001-12"},
}

var docTypeParameters = []test.Parameter{
	{Input: models.Document{Type: "CPF", Number: "12312100100"}, Expected: 1},
	{Input: models.Document{Type: "CNPJ", Number: "12123123000112"}, Expected: 2},
}

var sanitizeCepParameters = []test.Parameter{
	{Input: "25368-100", Expected: "25368100"},
	{Input: "25368100", Expected: "25368100"},
}

var mod11BradescoShopFacilDvParameters = []test.Parameter{
	{Input: "00000000006", Expected: "0"},
	{Input: "00000000001", Expected: "P"},
	{Input: "00000000002", Expected: "8"},
}

var sanitizeCitibankSpecialCharacteresParameters = []test.Parameter{
	{Input: "", Length: 0, Expected: ""},       //Default string value
	{Input: "   ", Length: 3, Expected: "   "}, //Whitespaces
	{Input: "a b", Length: 3, Expected: "a b"},
	{Input: "/-;@", Length: 4, Expected: "/-;@"}, //Caracteres especiais aceitos pelo Citibank
	{Input: "???????????????????????????a-zA-Z0-9ÁÉÍÓÚÀÈÌÒÙÂÊÎÔÛÃÕáéíóúàèìòùâêîôûãõç.", Length: 45, Expected: "a-zA-Z0-9AEIOUAEIOUAEIOUAOaeiouaeiouaeiouaoc."},
	{Input: "Ol@ Mundo. você pode ver uma barra /, mas não uma exclamação!?; Nem Isso", Length: 60, Expected: "Ol@ Mundo. voce pode ver uma barra / mas nao uma exclamacao;"},
	{Input: "Avenida Andr? Rodrigues de Freitas", Length: 33, Expected: "Avenida Andr Rodrigues de Freitas"},
}

func TestShouldPadLeft(t *testing.T) {
	expected := "00005"

	result := padLeft("5", "0", 5)

	assert.Equal(t, expected, result, "O texto deve ter zeros a esqueda e até 5 caracteres")
}

func TestShouldReturnString(t *testing.T) {
	expected := "5"

	result := toString(5)

	assert.Equal(t, expected, result, "O número deve ser uma string")
}

func TestFormatDigitableLine(t *testing.T) {
	for _, fact := range formatDigitableLineParameters {
		result := fmtDigitableLine(fact.Input.(string))
		assert.Equal(t, fact.Expected, result, "A linha digitável deve ser formatada corretamente")
	}
}

func TestTruncate(t *testing.T) {
	for _, fact := range truncateParameters {
		result := truncateString(fact.Input.(string), fact.Length)
		assert.Equal(t, fact.Expected, result, "Deve-se truncar uma string corretamente")
	}
}

func TestClearString(t *testing.T) {
	for _, fact := range clearStringParameters {
		result := clearString(fact.Input.(string))
		assert.Equal(t, fact.Expected, result, "Deve-se limpar uma string corretamente")
	}
}

func TestClearStringCaixa(t *testing.T) {
	for _, fact := range clearStringCaixaParameters {
		result := clearStringCaixa(fact.Input.(string))
		assert.Equal(t, fact.Expected, result, "Deve-se limpar uma string corretamente")
	}
}

func TestJoinStringSpace(t *testing.T) {
	expected := "a b c"

	result := joinSpace("a", "b", "c")

	assert.Equal(t, expected, result, "Deve-se fazer um join em uma string com espaços")
}

func TestFormatCNPJ(t *testing.T) {
	expected := "01.000.000/0001-00"

	result := fmtCNPJ("01000000000100")

	assert.Equal(t, expected, result, "O CNPJ deve ser formatado corretamente")
}

func TestFormatCPF(t *testing.T) {
	expected := "123.121.001-00"

	result := fmtCPF("12312100100")

	assert.Equal(t, expected, result, "O CPF deve ser formatado corretamente")
}

func TestFormatNumber(t *testing.T) {
	for _, fact := range formatNumberParameters {
		result := fmtNumber(fact.Input)
		assert.Equal(t, fact.Expected, result, "O valor em inteiro deve ser convertido para uma string com duas casas decimais separado por vírgula (0,00)")
	}
}

func TestMod11OurNumber(t *testing.T) {
	var expected, onlyDigitExpected uint
	expected = 120000001148
	onlyDigitExpected = 8

	result := calculateOurNumberMod11(12000000114, false)
	onlyDigitResult := calculateOurNumberMod11(12000000114, true)

	assert.Equal(t, expected, result, "Deve-se calcular o mod11 do nosso número e retornar o digito à esquerda")
	assert.Equal(t, onlyDigitExpected, onlyDigitResult, "Deve-se calcular o mod11 do nosso número e retornar o digito à esquerda")
}

func TestToFloatStr(t *testing.T) {
	for _, fact := range toFloatStrParameters {
		result := toFloatStr(fact.Input)
		assert.Equal(t, fact.Expected, result, "O valor em inteiro deve ser convertido para uma string com duas casas decimais separado por ponto (0.00)")
	}
}

func TestFormatDoc(t *testing.T) {
	for _, fact := range formatDocParameters {
		result := fmtDoc(fact.Input.(models.Document))
		assert.Equal(t, fact.Expected, result, "O documento deve ser formatado corretamente")
	}
}

func TestDocType(t *testing.T) {
	for _, fact := range docTypeParameters {
		result := docType(fact.Input.(models.Document))
		assert.Equal(t, fact.Expected, result, "O documento deve ser do tipo correto")
	}
}

func TestTrim(t *testing.T) {
	expected := "hue br festa"

	result := trim(" hue br festa ")

	assert.Equal(t, expected, result, "O texto não deve ter espaços no início e no final")
}

func TestSanitizeHtml(t *testing.T) {
	expected := "hu3 br festa"

	result := sanitizeHtmlString("<b>hu3 br festa</b>")

	assert.Equal(t, expected, result, "O texto não deve conter HTML tags")
}

func TestUnscapeHtml(t *testing.T) {
	var expected template.HTML
	expected = "ó"

	result := unescapeHtmlString("&#243;")

	assert.Equal(t, expected, result, "A string não deve ter caracteres Unicode")
}

func TestSanitizeCep(t *testing.T) {
	for _, fact := range sanitizeCepParameters {
		result := extractNumbers(fact.Input.(string))
		assert.Equal(t, fact.Expected, result, "o zipcode deve conter apenas números")
	}
}

func TestDVOurNumberMod11BradescoShopFacil(t *testing.T) {
	wallet := "19"
	for _, fact := range mod11BradescoShopFacilDvParameters {
		result := mod11BradescoShopFacilDv(fact.Input.(string), wallet)
		assert.Equal(t, fact.Expected, result, "o dígito verificador deve ser equivalente ao OurNumber")
	}
}

func TestEscape(t *testing.T) {
	expected := "KM 5,00    "

	result := escapeStringOnJson("KM 5,00 \t \f \r \b")

	assert.Equal(t, expected, result, "O texto deve ser escapado")
}

func TestRemoveCharacterSpecial(t *testing.T) {
	expected := "Texto com carácter especial   -"

	result := removeSpecialCharacter("Texto? com \"carácter\" especial * ' -")

	assert.Equal(t, expected, result, "Os caracteres especiais devem ser removidos")
}

func TestCitiBankSanitizeString(t *testing.T) {
	for _, fact := range sanitizeCitibankSpecialCharacteresParameters {
		input := fact.Input.(string)
		result := sanitizeCitibankSpecialCharacteres(input, fact.Length)
		assert.Equal(t, fact.Expected, result, "Caracteres especiais e acentos devem ser removidos")
		assert.Equal(t, fact.Length, len(result), "O texto deve ser devidamente truncado")
	}
}
