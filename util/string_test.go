package util

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

type UtilTestParameter struct {
	Input    interface{}
	Expected interface{}
}

var removeDiacriticsParameters = []UtilTestParameter{
	{Input: "maçã", Expected: "maca"},
	{Input: "áÉçãẽś", Expected: "aEcaes"},
	{Input: "Týr", Expected: "Tyr"},
	{Input: "párãlèlëpípêdö", Expected: "paralelepipedo"},
}

var padLeftParameters = []UtilTestParameter{
	{Input: "123", Expected: "0000000123"},
	{Input: "1234567890", Expected: "1234567890"},
}

func TestRemoveDiacritics(t *testing.T) {
	for _, fact := range removeDiacriticsParameters {
		result := RemoveDiacritics(fact.Input.(string))
		assert.Equal(t, fact.Expected, result, "Deve receber um texto com acentos e retornar o texto sem acentos")
	}
}

func TestPadLeft(t *testing.T) {
	length := 10
	paddingCaracter := "0"

	for _, fact := range padLeftParameters {
		result := PadLeft(fact.Input.(string), paddingCaracter, uint(length))
		assert.Equal(t, fact.Expected, result, "O numero deve ser ajustado corretamente")
	}
}
