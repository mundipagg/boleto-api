package validations

import "github.com/mundipagg/boleto-api/models"

//ValidateBuyerDocumentNumber verifica se o número do documento do pagador é válido
func ValidateBuyerDocumentNumber(b interface{}) error {
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
		return InvalidType(t)
	}
}
