package db

import (
	"fmt"
	"sync"
	"time"

	"github.com/mundipagg/boleto-api/config"
	"github.com/mundipagg/boleto-api/models"
	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

type mongoDb struct {
	m sync.RWMutex
}

var dbName = "Boleto"

//CreateMongo cria uma nova intancia de conexão com o mongodb
func CreateMongo() (DB, error) {
	db := new(mongoDb)
	if config.Get().MockMode {
		dbName = "boletoapi_mock"
	}
	return db, nil
}

func getInfo() *mgo.DialInfo {
	return &mgo.DialInfo{
		Addrs:    []string{config.Get().MongoURL},
		Timeout:  10 * time.Second,
		Database: "Boleto",
		Username: config.Get().MongoUser,
		Password: config.Get().MongoPassword,
	}
}

//SaveBoleto salva um boleto no mongoDB
func (e *mongoDb) SaveBoleto(boleto models.BoletoView) error {
	var err error
	e.m.Lock()
	defer e.m.Unlock()

	session, err := mgo.DialWithInfo(getInfo())
	if err != nil {
		return models.NewInternalServerError(err.Error(), "Falha ao conectar com o banco de dados")
	}
	defer session.Close()
	c := session.DB(dbName).C("boletos")
	err = c.Insert(boleto)
	return err
}

//GetBoletoById busca um boleto pelo ID que vem na URL
func (e *mongoDb) GetBoletoByID(id string) (models.BoletoView, error) {
	e.m.Lock()
	defer e.m.Unlock()
	result := models.BoletoView{}
	session, err := mgo.DialWithInfo(getInfo())
	if err != nil {
		return result, models.NewInternalServerError(err.Error(), "Falha ao conectar com o banco de dados")
	}
	defer session.Close()
	c := session.DB(dbName).C("boletos")
	errF := c.Find(bson.M{"id": id}).One(&result)
	if errF != nil {
		return models.BoletoView{}, err
	}
	return result, nil
}

func (e *mongoDb) Close() {
	fmt.Println("Close Database Connection")
}
