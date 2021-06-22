package stonebank

import (
	"fmt"
	"net/http"

	"github.com/mundipagg/boleto-api/log"
	"github.com/mundipagg/boleto-api/models"
	"github.com/mundipagg/boleto-api/util"
	"github.com/mundipagg/boleto-api/validations"
)

type bankStoneBank struct {
	validate  *models.Validator
	log       *log.Log
	transport *http.Transport
}

//New Create a new Santander Integration Instance
func New() (bankStoneBank, error) {
	var err error
	b := bankStoneBank{
		validate: models.NewValidator(),
		log:      log.CreateLog(),
	}

	b.transport, err = util.BuildTLSTransport()
	if err != nil {
		return bankStoneBank{}, err
	}

	b.validate.Push(validations.ValidateAmount)
	b.validate.Push(validations.ValidateExpireDate)
	b.validate.Push(validations.ValidateBuyerDocumentNumber)
	b.validate.Push(validations.ValidateRecipientDocumentNumber)

	return b, nil
}

func (b bankStoneBank) ProcessBoleto(request *models.BoletoRequest) (models.BoletoResponse, error) {
	b.log.Info(fmt.Sprintf("StoneBank ProcessBoleto %v", request))
	return models.BoletoResponse{}, nil
}

func (b bankStoneBank) RegisterBoleto(request *models.BoletoRequest) (models.BoletoResponse, error) {
	b.log.Info(fmt.Sprintf("StoneBank Register Boleto %v", request))
	return models.BoletoResponse{}, nil
}

func (b bankStoneBank) ValidateBoleto(request *models.BoletoRequest) models.Errors {
	b.log.Info(fmt.Sprintf("StoneBank ValidateBoleto %v", request))
	return nil
}

func (b bankStoneBank) GetBankNumber() models.BankNumber {
	b.log.Info("StoneBank GetBankNumber")
	return 197
}

func (b bankStoneBank) GetBankNameIntegration() string {
	b.log.Info("StoneBank GetBankNameIntegration")
	return "stonebank"
}

func (b bankStoneBank) Log() *log.Log {
	b.log.Info("StoneBank Log")
	return b.log
}
