package bank

import (
	"github.com/mundipagg/boleto-api/models"
	"github.com/mundipagg/boleto-api/pefisa"
)

func getIntegrationPefisa(boleto models.BoletoRequest) (Bank, error) {
	return pefisa.New(), nil
}
