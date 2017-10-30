package santander

import (
	"fmt"

	"github.com/mundipagg/boleto-api/models"
	"github.com/mundipagg/boleto-api/validations"
)

func santanderValidateAgreementNumber(b interface{}) error {
	switch t := b.(type) {
	case *models.BoletoRequest:
		if t.Agreement.AgreementNumber == 0 {
			return models.NewErrorResponse("MP400", fmt.Sprintf("O código do convênio deve ser preenchido"))
		}
		return nil
	default:
		return validations.InvalidType(t)
	}
}
