package bank

import (
	"github.com/mundipagg/boleto-api/citibank"
	"github.com/mundipagg/boleto-api/models"
)

func getIntegrationCitibank(boleto models.BoletoRequest) (Bank, error) {
	return citibank.New()
}
