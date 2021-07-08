package stone

import (
	"github.com/mundipagg/boleto-api/models"
	"github.com/mundipagg/boleto-api/validations"
)

func stoneValidateAccessKeyNotEmpty(b interface{}) error {
	switch t := b.(type) {
	case *models.BoletoRequest:
		if t.Authentication.AccessKey == "" {
			return models.NewErrorResponse("MP400", "o campo AccessKey não pode ser vazio")
		}
		return nil
	default:
		return validations.InvalidType(t)
	}
}
