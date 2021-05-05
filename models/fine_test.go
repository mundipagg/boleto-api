package models

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

type testFineParameter struct {
	Input         interface{}
	Expected      interface{}
	AmountInCents interface{}
}

var fineParameters = []testFineParameter{
	{Input: Fine{Type: "Percentual", Days: 0, Rate: 0.010}, AmountInCents: 0, Expected: true},
	{Input: Fine{Type: "Percentual ", Days: 0, Rate: 0.010}, AmountInCents: 0, Expected: false},
	{Input: Fine{Type: "PERCENTUAL", Days: 0, Rate: 0.010}, AmountInCents: 0, Expected: true},
	{Input: Fine{Type: "pErCeNtUaL", Days: 0, Rate: 0.010}, AmountInCents: 0, Expected: true},
	{Input: Fine{Type: "Porcentagem", Days: 0, Rate: 0.010}, AmountInCents: 0, Expected: false},
	{Input: Fine{Type: "Percentual", Days: -1, Rate: 0.010}, AmountInCents: 0, Expected: false},
	{Input: Fine{Type: "Percentual", Days: 0, Rate: 0.011}, AmountInCents: 0, Expected: false},
	{Input: Fine{Type: "Percentual", Days: 0, Rate: 0.0101}, AmountInCents: 0, Expected: false},
	{Input: Fine{Type: "Percentual", Days: 0, Rate: 0.01001}, AmountInCents: 0, Expected: false},
	{Input: Fine{Type: "Percentual", Days: 0, Rate: 0.010001}, AmountInCents: 0, Expected: false},
	{Input: Fine{Type: "Nominal", Days: 0, Rate: 1.00}, AmountInCents: 100, Expected: true},
	{Input: Fine{Type: "NOMINAL", Days: 0, Rate: 1.00}, AmountInCents: 100, Expected: true},
	{Input: Fine{Type: "nOmInAl", Days: 0, Rate: 1.00}, AmountInCents: 100, Expected: true},
	{Input: Fine{Type: "Absoluto", Days: 0, Rate: 1.00}, AmountInCents: 100, Expected: false},
	{Input: Fine{Type: "Nominal", Days: 0, Rate: 2.00}, AmountInCents: 100, Expected: false},
	{Input: Fine{Type: "Nominal", Days: 0, Rate: 2.00}, AmountInCents: 0, Expected: false},
}

func TestFine_IsValid(t *testing.T) {
	i := 1
	for _, fact := range fineParameters {
		input := fact.Input.(Fine)

		result := input.IsValid(fact.AmountInCents.(int))

		assert.Equal(t, fact.Expected, result, fmt.Sprintf("Teste %d: Os juros não foram validados corretamente", i))
		i++
	}
}
