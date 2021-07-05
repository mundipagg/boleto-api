package stone

import (
	"github.com/mundipagg/boleto-api/models"
	"github.com/mundipagg/boleto-api/validations"
)

func stoneValidateAccessKey(b interface{}) error {
	switch t := b.(type) {
	case *models.BoletoRequest:
		if t.Authentication.AccessKey == "" {
			return models.NewErrorResponse("MP400", "accessKey is required")
		}
		return nil
	default:
		return validations.InvalidType(t)
	}
}
