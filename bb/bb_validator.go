package bb

import (
	"strconv"

	"github.com/mundipagg/boleto-api/models"
)

var bb bankBB

func modElevenCalculator(a string, m []int) string {
	sum := sumAccountDigits(a, m)

	digit := 11 - sum%11

	if digit == 10 {
		return "X"
	}

	if digit == 11 {
		return "0"
	}

	return strconv.Itoa(digit)
}

func bbAgencyDigitCalculator(agency string) string {
	multiplier := []int{5, 4, 3, 2}
	return modElevenCalculator(agency, multiplier)
}

func bbAccountDigitCalculator(agency, account string) string {
	multiplier := []int{9, 8, 7, 6, 5, 4, 3, 2}
	return modElevenCalculator(account, multiplier)
}

func bbValidateAgencyAndDigit(b interface{}) error {
	switch t := b.(type) {
	case *models.BoletoRequest:
		err := t.Agreement.IsAgencyValid()
		if err != nil {
			return err
		}
		t.Agreement.CalculateAgencyDigit(bbAgencyDigitCalculator)
		return nil
	default:
		return invalidType(t)
	}
}

func invalidType(t interface{}) error {
	return models.NewErrorResponse("MP500", "Tipo inválido")
}

func bbValidateAccountAndDigit(b interface{}) error {
	switch t := b.(type) {
	case *models.BoletoRequest:
		err := t.Agreement.IsAccountValid(8)
		if err != nil {
			return err
		}
		t.Agreement.CalculateAccountDigit(bbAccountDigitCalculator)
		return nil
	default:
		return invalidType(t)
	}
}

func bbValidateOurNumber(b interface{}) error {
	switch t := b.(type) {
	case *models.BoletoRequest:
		if t.Title.OurNumber > 9999999999 {
			return models.NewErrorResponse("MPOurNumber", "Nosso número inválido")
		}
		return nil
	default:
		return invalidType(t)
	}
}

func bbValidateWalletVariation(b interface{}) error {
	switch t := b.(type) {
	case *models.BoletoRequest:
		if t.Agreement.WalletVariation < 1 {
			return models.NewErrorResponse("MPWalletVariation", "Variação da carteira inválida")
		}
		return nil
	default:
		return invalidType(t)
	}
}

func bbValidateTitleInstructions(b interface{}) error {
	switch t := b.(type) {
	case *models.BoletoRequest:
		return t.Title.ValidateInstructionsLength(220)
	default:
		return invalidType(t)
	}
}

func bbValidateTitleDocumentNumber(b interface{}) error {
	switch t := b.(type) {
	case *models.BoletoRequest:
		return t.Title.ValidateDocumentNumber()
	default:
		return invalidType(t)
	}
}
