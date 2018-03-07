package db

import (
	"sync"

	"github.com/mundipagg/boleto-api/config"
	"github.com/mundipagg/boleto-api/models"
)

//DB é a interface basica para os métodos de persistência
type DB interface {
	SaveBoleto(models.BoletoView) error
	GetBoletoByID(string) (models.BoletoView, error)
	Close()
}

var db DB = nil
var _create sync.Mutex
var err error

//GetDB retorna o objeto concreto que implementa as funções de persistência
func GetDB() (DB, error) {
	if config.Get().MockMode || config.Get().DevMode {
		return new(mock), nil
	}
	if config.Get().EnablePostgres {
		return CreatePostgreSQL()
	}
	return CreateMongo()
}
