package db

import (
	"errors"

	"github.com/mundipagg/boleto-api/cache"
	"github.com/mundipagg/boleto-api/models"
)


//SaveBoleto salva o boleto num cache local em memoria
func SaveBoletoMock(boleto models.BoletoView) error {
	idBson, _ := boleto.ID.MarshalText()
	cache.Set(string(idBson), boleto)
	return nil
}

//GetBoletoById retorna o boleto por id do cache em memoria
func GetBoletoByIDMock(id string) (models.BoletoView, error) {
	c, ok := cache.Get(string(id))
	if !ok {
		return models.BoletoView{}, errors.New("Boleto não encontrado")
	}
	return c.(models.BoletoView), nil
}

func Close() {}
