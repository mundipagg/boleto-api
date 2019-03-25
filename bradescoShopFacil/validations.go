package bradescoShopFacil

import (
	"github.com/mundipagg/boleto-api/models"
	"github.com/mundipagg/boleto-api/validations"
)

func bradescoShopFacilValidateAgency(b interface{}) error {
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

func bradescoShopFacilValidateAccount(b interface{}) error {
	switch t := b.(type) {
	case *models.BoletoRequest:
		if t.Agreement.Account == "" {
			return models.NewErrorResponse("MP400", "a conta deve ser preenchida")
		}
		return nil
	default:
		return validations.InvalidType(t)
	}
}

func bradescoShopFacilValidateWallet(b interface{}) error {
	switch t := b.(type) {
	case *models.BoletoRequest:
		if t.Agreement.Wallet != 25 && t.Agreement.Wallet != 26 {
			return models.NewErrorResponse("MP400", "a carteira deve ser 25 ou 26 para o BradescoShopFacil")
		}
		return nil
	default:
		return validations.InvalidType(t)
	}
}

func bradescoShopFacilValidateAuth(b interface{}) error {
	switch t := b.(type) {
	case *models.BoletoRequest:
		if t.Authentication.Username == "" || t.Authentication.Password == "" {
			return models.NewErrorResponse("MP400", "o nome de usuário e senha devem ser preenchidos")
		}
		return nil
	default:
		return validations.InvalidType(t)
	}
}

func bradescoShopFacilValidateAgreement(b interface{}) error {
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

func bradescoShopFacilBoletoTypeValidate(b interface{}) error {
	bt := bradescoShopFacilBoletoTypes()

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
