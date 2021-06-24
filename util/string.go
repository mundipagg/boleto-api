package util

import (
	"fmt"
	"strconv"
	"unicode"

	"encoding/json"

	"github.com/tdewolff/minify"
	"github.com/tdewolff/minify/css"
	"github.com/tdewolff/minify/html"
	"github.com/tdewolff/minify/js"
	jm "github.com/tdewolff/minify/json"

	"golang.org/x/text/transform"
	"golang.org/x/text/unicode/norm"
)

func isMn(r rune) bool {
	return unicode.Is(unicode.Mn, r) // Mn: nonspacing marks
}

//RemoveDiacritics remove caracteres especiais de uma string
func RemoveDiacritics(s string) string {
	t := transform.Chain(norm.NFD, transform.RemoveFunc(isMn), norm.NFC)
	result, _, _ := transform.String(t, s)
	return result
}

// PadLeft insere um caractere a esquerda de um texto
func PadLeft(value, char string, total uint) string {
	s := "%" + char + strconv.Itoa(int(total)) + "s"
	return fmt.Sprintf(s, value)
}

//Stringify convete objeto para JSON
func Stringify(o interface{}) string {
	b, _ := json.Marshal(o)
	return string(b)
}

//ParseJSON converte string para um objeto GO
func ParseJSON(s string, o interface{}) interface{} {
	json.Unmarshal([]byte(s), o)
	return o
}

//IsDigit Verifica se um caracter é um dígito numérico de acordo com o código decimal da Tabela ASCII,
// onde o '0' representa o valor 48 e o '9' o valor 57
func IsDigit(r rune) bool {
	return (r >= 48 && r <= 57)
}

//IsBasicCharacter Verifica se um caracter é uma letra sem acento, maiúscula ou minúscula, de acordo com o código decimal da Tabela ASCII
// onde o 'A' representa o valor 65 e o 'Z' o valor 90 e o 'a' representa o valor 97 e 'z' o valor 122
//true para 0123456789
func IsBasicCharacter(r rune) bool {
	return (r >= 65 && r <= 90) || (r >= 97 && r <= 122)
}

//IsCaixaSpecialCharacter Verifica se um caracter especial é aceito Caixa Econômica, de acordo com o  código decimal da Tabela ASCII
// sendo aceito os seguinte caracteres:
//	esp	32	'	39   :	58
//	!	33	(	40   ;	59
//			)	41   =	61
//			*	42
//			+	43   ?	63
//			,	44
//			-	45   _  95
//			.	46
//			/	47
//
// OBS: Apesar de descritos como aceitos, os caracteres & (38)  < (60) e > (62) foram removidos pois não
// estão disponíveis para XML. Testamos seus respectivos encodes: &amp; &lt; &gt; entretanto recebemos a
// resposta (66) CARACTER INVALIDO.
func IsCaixaSpecialCharacter(r rune) bool {
	caixaSpecialCharacters := []rune{32, 33, 39, 40, 41, 42, 43, 44, 45, 46, 47, 58, 59, 61, 63, 95}
	for _, c := range caixaSpecialCharacters {
		if r == c {
			return true
		}
	}
	return false
}

//MinifyString Minifica uma string de acordo com um determinado formato
func MinifyString(mString, tp string) string {
	m := minify.New()
	m.Add("text/html", &html.Minifier{
		KeepDocumentTags:        true,
		KeepEndTags:             true,
		KeepWhitespace:          false,
		KeepConditionalComments: true,
	})
	m.AddFunc("text/css", css.Minify)
	m.AddFunc("text/javascript", js.Minify)
	m.AddFunc("application/json", jm.Minify)

	s, err := m.String(tp, mString)

	if err != nil {
		return mString
	}

	return s
}
