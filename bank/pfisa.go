package bank

import (
	"github.com/mundipagg/boleto-api/models"
	"github.com/mundipagg/boleto-api/pfisa"
)

func getIntegrationPfisa(boleto models.BoletoRequest) (Bank, error) {
	return pfisa.New(), nil
}
