package caixa

import (
	"fmt"
	"strconv"

	"github.com/mundipagg/boleto-api/models"
)

func caixaAccountDigitCalculator(agency, account string) string {
	multiplier := []int{8, 7, 6, 5, 4, 3, 2, 9, 8, 7, 6, 5, 4, 3, 2}
	toCheck := fmt.Sprintf("%04s%011s", agency, account)
	return caixaModElevenCalculator(toCheck, multiplier)
}

func caixaModElevenCalculator(a string, m []int) string {
	sum := sumAccountDigits(a, m)
	digit := (sum * 10) % 11
	if digit == 10 {
		return "0"
	}
	return strconv.Itoa(digit)
}

func validateInstructions(b interface{}) error {
	switch t := b.(type) {
	case *models.BoletoRequest:
		if len(t.Title.Instructions) > 40 {
			return models.NewErrorResponse("MP400", "O número máximo permitido para instruções é de 40 caracteres")
		}
		return nil
	default:
		return invalidType(t)
	}
}

func caixaValidateAccountAndDigit(b interface{}) error {
	switch t := b.(type) {
	case *models.BoletoRequest:
		err := t.Agreement.IsAccountValid(11)
		if err != nil {
			return err
		}
		errAg := t.Agreement.IsAgencyValid()
		if errAg != nil {
			return errAg
		}
		t.Agreement.CalculateAccountDigit(caixaAccountDigitCalculator)
		return nil
	default:
		return invalidType(t)
	}
}

func caixaAgencyDigitCalculator(agency string) string {
	multiplier := []int{5, 4, 3, 2}
	return modElevenCalculator(agency, multiplier)
}

func caixaValidateAgency(b interface{}) error {
	switch t := b.(type) {
	case *models.BoletoRequest:
		err := t.Agreement.IsAgencyValid()
		if err != nil {
			return err
		}
		return nil
	default:
		return invalidType(t)
	}
}
