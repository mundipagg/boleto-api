package stone

import (
	"github.com/mundipagg/boleto-api/log"
	"github.com/mundipagg/boleto-api/models"
	"github.com/mundipagg/boleto-api/validations"
)

type stone struct {
	validate *models.Validator
	log      *log.Log
}

//New Create a new Stone Integration Instance
func New() (stone, error) {
	var err error
	b := stone{
		validate: models.NewValidator(),
		log:      log.CreateLog(),
	}

	if err != nil {
		return stone{}, err
	}

	b.validate.Push(validations.ValidateAmount)
	b.validate.Push(validations.ValidateExpireDate)
	b.validate.Push(validations.ValidateBuyerDocumentNumber)
	b.validate.Push(validations.ValidateRecipientDocumentNumber)

	b.validate.Push(stoneValidateAccessKeyNotEmpty)

	return b, nil
}

func (b stone) ProcessBoleto(boleto *models.BoletoRequest) (models.BoletoResponse, error) {
	errs := b.ValidateBoleto(boleto)
	if len(errs) > 0 {
		return models.BoletoResponse{Errors: errs}, nil
	}

	if accToken, err := authenticate(boleto.Authentication.AccessKey, b.log); err != nil {
		return models.BoletoResponse{Errors: errs}, err
	} else {
		boleto.Authentication.AuthorizationToken = accToken
	}

	return b.RegisterBoleto(boleto)
}

func (b stone) RegisterBoleto(request *models.BoletoRequest) (models.BoletoResponse, error) {
	return models.BoletoResponse{}, nil
}

func (b stone) ValidateBoleto(request *models.BoletoRequest) models.Errors {
	return models.Errors(b.validate.Assert(request))
}

func (b stone) GetBankNumber() models.BankNumber {
	return models.Stone
}

func (b stone) GetBankNameIntegration() string {
	return "stone"
}

func (b stone) Log() *log.Log {
	return b.log
}
