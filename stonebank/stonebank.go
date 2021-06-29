package stonebank

import (
	"fmt"

	"github.com/mundipagg/boleto-api/log"
	"github.com/mundipagg/boleto-api/models"
	"github.com/mundipagg/boleto-api/validations"
)

type bankStoneBank struct {
	validate *models.Validator
	log      *log.Log
}

//New Create a new Santander Integration Instance
func New() (bankStoneBank, error) {
	var err error
	b := bankStoneBank{
		validate: models.NewValidator(),
		log:      log.CreateLog(),
	}

	if err != nil {
		return bankStoneBank{}, err
	}

	b.validate.Push(validations.ValidateAmount)
	b.validate.Push(validations.ValidateExpireDate)
	b.validate.Push(validations.ValidateBuyerDocumentNumber)
	b.validate.Push(validations.ValidateRecipientDocumentNumber)

	b.validate.Push(stoneBankValidateAccessKeyNotEmpty)

	return b, nil
}

func (b bankStoneBank) ProcessBoleto(boleto *models.BoletoRequest) (models.BoletoResponse, error) {
	errs := b.ValidateBoleto(boleto)
	if len(errs) > 0 {
		return models.BoletoResponse{Errors: errs}, nil
	}

	if accToken, err := authenticate(boleto.Authentication.AccessKey); err != nil {
		return models.BoletoResponse{Errors: errs}, err
	} else {
		boleto.Authentication.AuthorizationToken = accToken
	}

	return b.RegisterBoleto(boleto)
}

func (b bankStoneBank) RegisterBoleto(request *models.BoletoRequest) (models.BoletoResponse, error) {
	b.log.Info(fmt.Sprintf("StoneBank Register Boleto %v", request))
	return models.BoletoResponse{}, nil
}

func (b bankStoneBank) ValidateBoleto(request *models.BoletoRequest) models.Errors {
	return models.Errors(b.validate.Assert(request))
}

func (b bankStoneBank) GetBankNumber() models.BankNumber {
	return models.StoneBank
}

func (b bankStoneBank) GetBankNameIntegration() string {
	return "stonebank"
}

func (b bankStoneBank) Log() *log.Log {
	return b.log
}
