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

var digitParameters = []UtilTestParameter{
	{Input: "0123456789", Expected: true},
	{Input: " ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz,/()*&=-+!:?<>.;_\"", Expected: false},
}

var basicCharacter = []UtilTestParameter{
	{Input: " 0123456789,/()*&=-+!:?<>.;_\"", Expected: false},
	{Input: "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz", Expected: true},
}

var caixaSpecialCharacter = []UtilTestParameter{
	{Input: " ,/()*=-+!:?.;_'", Expected: true},
	{Input: "01223456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz@#$%¨{}[]^~Çç\"&<>\\", Expected: false},
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

func TestIsDigit(t *testing.T) {
	for _, fact := range digitParameters {
		s := fact.Input.(string)
		for _, c := range s {
			result := IsDigit(c)
			assert.Equal(t, fact.Expected, result, "A verificação de dígito deve ocorrer corretamente")
		}
	}
}

func TestIsBasicCharacter(t *testing.T) {
	for _, fact := range basicCharacter {
		s := fact.Input.(string)
		for _, c := range s {
			result := IsBasicCharacter(c)
			assert.Equal(t, fact.Expected, result, "A verificação de caracter deve ocorrer corretamente")
		}
	}
}

func TestIsSpecialCharacterCaixa(t *testing.T) {
	for _, fact := range caixaSpecialCharacter {
		s := fact.Input.(string)
		for _, c := range s {
			result := IsCaixaSpecialCharacter(c)
			assert.Equal(t, fact.Expected, result, "A verificação de caracter deve ocorrer corretamente")
		}
	}
}
