package bank

import (
	"github.com/mundipagg/boleto-api/models"
	"github.com/mundipagg/boleto-api/stone"
)

func getIntegrationStone(boleto models.BoletoRequest) (Bank, error) {
	return stone.New(), nil
}
