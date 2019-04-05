package bradescoNetEmpresa

import (
	"github.com/mundipagg/boleto-api/models"
	"github.com/mundipagg/boleto-api/validations"
)

func bradescoNetEmpresaValidateAgency(b interface{}) error {
	switch t := b.(type) {
	case *models.BoletoRequest:
		err := t.Agreement.IsAgencyValid()
		if err != nil {
			return models.NewErrorResponse("MP400", err.Error())
		}
		return nil
	default:
		return validations.InvalidType(t)
	}
}

func bradescoNetEmpresaValidateAccount(b interface{}) error {
	switch t := b.(type) {
	case *models.BoletoRequest:
		err := t.Agreement.IsAccountValid(7)
		if err != nil {
			return models.NewErrorResponse("MP400", err.Error())
		}
		return nil
	default:
		return validations.InvalidType(t)
	}
}

func bradescoNetEmpresaValidateWallet(b interface{}) error {
	switch t := b.(type) {
	case *models.BoletoRequest:
		if t.Agreement.Wallet != 4 && t.Agreement.Wallet != 9 && t.Agreement.Wallet != 19 {
			return models.NewErrorResponse("MP400", "a carteira deve ser 4, 9 ou 19 para o bradescoNetEmpresa")
		}
		return nil
	default:
		return validations.InvalidType(t)
	}
}

func bradescoNetEmpresaValidateAgreement(b interface{}) error {
	switch t := b.(type) {
	case *models.BoletoRequest:
		if t.Agreement.AgreementNumber == 0 {
			return models.NewErrorResponse("MP400", "o código do contrato deve ser preenchido")
		}
		return nil
	default:
		return validations.InvalidType(t)
	}
}

func bradescoNetEmpresaBoletoTypeValidate(b interface{}) error {
	bt := bradescoNetEmpresaBoletoTypes()

	switch t := b.(type) {

	case *models.BoletoRequest:
		if len(t.Title.BoletoType) > 0 && bt[t.Title.BoletoType] == "" {
			return models.NewErrorResponse("MP400", "espécie de boleto informada não existente")
		}
		return nil
	default:
		return validations.InvalidType(t)
	}
}
