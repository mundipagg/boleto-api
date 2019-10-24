package bank

import (
	"github.com/mundipagg/boleto-api/models"
	"github.com/mundipagg/boleto-api/santander"
)

func getIntegrationSantander(boleto models.BoletoRequest) (Bank, error) {
	return santander.New()
}
