package strings

import (
	"github.com/iancoleman/strcase"
	"regexp"
	"strings"
)

var emptyOrWhitespacePattern = regexp.MustCompile(`^\s*$`)

func IsBlank(str string) bool {
	return emptyOrWhitespacePattern.MatchString(str)
}

func ToPascalCase(str string) string {
	return Capitalize(strcase.ToCamel(str))
}

func Capitalize(str string) string {
	if len(str) <= 1 {
		return strings.ToUpper(str)
	}
	return strings.ToUpper(string(str[0])) + str[1:]
}
