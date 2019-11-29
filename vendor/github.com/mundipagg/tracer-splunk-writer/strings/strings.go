package strings

import (
	"bytes"
	"fmt"
	"html/template"
	"regexp"
	"strings"

	"github.com/iancoleman/strcase"
)

var emptyOrWhitespacePattern = regexp.MustCompile(`^\s*$`)

func IsBlank(str string) bool {
	return emptyOrWhitespacePattern.MatchString(str)
}

func ToPascalCase(str string) string {
	return Capitalize(strcase.ToCamel(str))
}

func ToLowerCamelCase(str string) string {
	return strcase.ToLowerCamel(str)
}

func UseAnnotation(str string) string {
	return str
}

func Capitalize(str string) string {
	if len(str) <= 1 {
		return strings.ToUpper(str)
	}
	return strings.ToUpper(string(str[0])) + str[1:]
}

func ProcessString(str string, vars interface{}) string {
	tmpl, err := template.New("tmpl").Parse(str)

	if err != nil {
		fmt.Println(err.Error())
	}
	return process(tmpl, vars)
}

func process(t *template.Template, vars interface{}) string {
	var tmplBytes bytes.Buffer

	err := t.Execute(&tmplBytes, vars)
	if err != nil {
		fmt.Println(err.Error())
	}
	return tmplBytes.String()
}
