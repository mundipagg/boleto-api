package bank

import (
	"github.com/mundipagg/boleto-api/models"
	"github.com/mundipagg/boleto-api/stonebank"
)

func getIntegrationStoneBank(boleto models.BoletoRequest) (Bank, error) {
	return stonebank.New()
}
