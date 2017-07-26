package validations

import (
	"strconv"

	"github.com/mundipagg/boleto-api/models"
)

func BaseValidateRecipientDocumentNumber(b interface{}) error {
	switch t := b.(type) {
	case *models.BoletoRequest:
		if t.Recipient.Document.IsCPF() {
			return t.Recipient.Document.ValidateCPF()
		}
		if t.Recipient.Document.IsCNPJ() {
			return t.Recipient.Document.ValidateCNPJ()
		}
		return models.NewErrorResponse("MPRecipientDocumentType", "Tipo de Documento inválido")
	default:
		return invalidType(t)
	}
}

func SumAccountDigits(a string, m []int) int {
	sum := 0
	for idx, c := range a {
		i, _ := strconv.Atoi(string(c))

		sum += i * m[idx]
	}
	return sum
}

func BaseValidateBuyerDocumentNumber(b interface{}) error {
	switch t := b.(type) {
	case *models.BoletoRequest:
		if t.Buyer.Document.IsCPF() {
			return t.Buyer.Document.ValidateCPF()
		}
		if t.Buyer.Document.IsCNPJ() {
			return t.Buyer.Document.ValidateCNPJ()
		}
		return models.NewErrorResponse("MPBuyerDocumentType", "Tipo de Documento inválido")
	default:
		return invalidType(t)
	}
}

func BaseValidateExpireDate(b interface{}) error {
	switch t := b.(type) {
	case *models.BoletoRequest:
		return t.Title.IsExpireDateValid()
	default:
		return invalidType(t)
	}
}

func BaseValidateAmountInCents(b interface{}) error {
	switch t := b.(type) {
	case *models.BoletoRequest:
		return t.Title.IsAmountInCentsValid()
	default:
		return invalidType(t)
	}
}
