package bank

import (
	"fmt"

	"github.com/mundipagg/boleto-api/bb"
	"github.com/mundipagg/boleto-api/log"
	"github.com/mundipagg/boleto-api/models"
)

//Bank é a interface que vai oferecer os serviços em comum entre os bancos
type Bank interface {
	ProcessBoleto(*models.BoletoRequest) (models.BoletoResponse, error)
	RegisterBoleto(*models.BoletoRequest) (models.BoletoResponse, error)
	ValidateBoleto(*models.BoletoRequest) models.Errors
	GetBankNumber() models.BankNumber
	Log() *log.Log
}

//Get retorna estrategia de acordo com o banco ou erro caso o banco não exista
func Get(number models.BankNumber) (Bank, error) {
	switch number {
	case models.BancoDoBrasil:
		return bb.New(), nil
	case models.Caixa:
		return caixa.New(), nil
	case models.Citibank:
		return citibank.New(), nil
	default:
		return nil, models.NewErrorResponse("MPBankNumber", fmt.Sprintf("Banco %d não existe", number))
	}
}
