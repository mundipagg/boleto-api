package models

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

type testInterestParameter struct {
	Input         interface{}
	Expected      interface{}
	AmountInCents interface{}
}

var interestParameters = []testInterestParameter{
	{Input: Interest{Type: "Percentual", Days: 0, Rate: 0.020}, AmountInCents: 0, Expected: true},
	{Input: Interest{Type: "Percentual ", Days: 0, Rate: 0.020}, AmountInCents: 0, Expected: false},
	{Input: Interest{Type: "PERCENTUAL", Days: 0, Rate: 0.020}, AmountInCents: 0, Expected: true},
	{Input: Interest{Type: "pErCeNtUaL", Days: 0, Rate: 0.020}, AmountInCents: 0, Expected: true},
	{Input: Interest{Type: "Porcentagem", Days: 0, Rate: 0.020}, AmountInCents: 0, Expected: false},
	{Input: Interest{Type: "Percentual", Days: -1, Rate: 0.020}, AmountInCents: 0, Expected: false},
	{Input: Interest{Type: "Percentual", Days: 0, Rate: 0.021}, AmountInCents: 0, Expected: false},
	{Input: Interest{Type: "Percentual", Days: 0, Rate: 0.0201}, AmountInCents: 0, Expected: false},
	{Input: Interest{Type: "Percentual", Days: 0, Rate: 0.02001}, AmountInCents: 0, Expected: false},
	{Input: Interest{Type: "Percentual", Days: 0, Rate: 0.020001}, AmountInCents: 0, Expected: false},
	{Input: Interest{Type: "Nominal", Days: 0, Rate: 2.00}, AmountInCents: 100, Expected: true},
	{Input: Interest{Type: "NOMINAL", Days: 0, Rate: 2.00}, AmountInCents: 100, Expected: true},
	{Input: Interest{Type: "nOmInAl", Days: 0, Rate: 2.00}, AmountInCents: 100, Expected: true},
	{Input: Interest{Type: "Absoluto", Days: 0, Rate: 2.00}, AmountInCents: 100, Expected: false},
	{Input: Interest{Type: "Nominal", Days: 0, Rate: 3.00}, AmountInCents: 100, Expected: false},
	{Input: Interest{Type: "Nominal", Days: 0, Rate: 3.00}, AmountInCents: 0, Expected: false},
}

func TestInterest_IsValid(t *testing.T) {
	i := 1
	for _, fact := range interestParameters {
		input := fact.Input.(Interest)

		result := input.IsValid(fact.AmountInCents.(int))

		assert.Equal(t, fact.Expected, result, fmt.Sprintf("Teste %d: A multa não foi validada corretamente", i))
		i++
	}
}
