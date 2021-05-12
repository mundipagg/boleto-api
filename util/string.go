package util

import (
	"fmt"
	"strconv"
	"unicode"

	"encoding/json"

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
//	esp	32	&	38	 :	58	 _  95
//	!	33	'	39   ;	59
//			(	40   <	60
//			)	41   =	61
//			*	42   >	62
//			+	43   ?	63
//			,	44
//			-	45
//			.	46
//			/	47
func IsCaixaSpecialCharacter(r rune) bool {
	return (r >= 32 && r <= 33) || (r >= 38 && r <= 47) || (r >= 58 && r <= 63) || r == 95
}
