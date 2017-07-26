package validations

import "github.com/mundipagg/boleto-api/models"

//ValidateRecipientDocumentNumber Verifica se o número do documento do recebedor é válido
func ValidateRecipientDocumentNumber(b interface{}) error {
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
		return InvalidType(t)
	}
}
