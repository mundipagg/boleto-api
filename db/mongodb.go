package db

import (
	"errors"
	"fmt"
	"net"
	"strings"
	"sync"
	"time"

	"github.com/mundipagg/boleto-api/config"
	"github.com/mundipagg/boleto-api/log"
	"github.com/mundipagg/boleto-api/models"
	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

//MongoDb Struct
type MongoDb struct {
	m sync.RWMutex
}

var (
	dbSession *mgo.Session
	err       error
)

const (
	NotFoundDoc = "not found"
	InvalidPK   = "invalid pk"
)

//CreateMongo cria uma nova intancia de conexão com o mongodb
func CreateMongo(l *log.Log) (*MongoDb, error) {

	if dbSession == nil {
		dbSession, err = mgo.Dial(config.Get().MongoURL)

		if err != nil {
			l.Warn(err.Error(), fmt.Sprintf("Error create connection with mongo %s", err.Error()))
			return nil, err
		}
	}

	db := new(MongoDb)

	return db, nil
}

func getInfo() *mgo.DialInfo {
	connMgo := strings.Split(config.Get().MongoURL, ",")
	return &mgo.DialInfo{
		Addrs:     connMgo,
		Timeout:   5 * time.Second,
		Database:  config.Get().MongoDatabase,
		PoolLimit: 512,
		Username:  config.Get().MongoUser,
		Password:  config.Get().MongoPassword,
	}
}

//SaveBoleto salva um boleto no mongoDB
func (e *MongoDb) SaveBoleto(boleto models.BoletoView) error {

	e.m.Lock()
	defer e.m.Unlock()

	session := dbSession.Copy()

	defer session.Close()

	c := session.DB(config.Get().MongoDatabase).C(config.Get().MongoBoletoCollection)
	err = c.Insert(boleto)

	return err
}

//GetBoletoByID busca um boleto pelo ID que vem na URL
//O retorno será um objeto BoletoView, o tempo decorrido da operação (em milisegundos) e algum erro ocorrido durante a operação
func (e *MongoDb) GetBoletoByID(id, pk string) (models.BoletoView, int64, error) {

	start := time.Now()

	e.m.Lock()
	defer e.m.Unlock()
	result := models.BoletoView{}

	session := dbSession.Copy()

	defer session.Close()

	c := session.DB(config.Get().MongoDatabase).C(config.Get().MongoBoletoCollection)

	for i := 0; i <= config.Get().RetryNumberGetBoleto; i++ {

		if len(id) == 24 {
			d := bson.ObjectIdHex(id)
			err = c.Find(bson.M{"_id": d}).One(&result)
		} else {
			err = c.Find(bson.M{"id": id}).One(&result)
		}

		if opErr, ok := err.(*net.OpError); ok && opErr.Timeout() {
			continue
		} else {
			break
		}
	}

	if err != nil {
		return models.BoletoView{}, time.Since(start).Milliseconds(), err
	} else if !hasValidKey(result, pk) {
		return models.BoletoView{}, time.Since(start).Milliseconds(), errors.New(InvalidPK)
	}

	return result, time.Since(start).Milliseconds(), nil
}

//GetUserCredentials Busca as Credenciais dos Usuários
func (e *MongoDb) GetUserCredentials() ([]models.Credentials, error) {
	e.m.Lock()
	defer e.m.Unlock()
	result := []models.Credentials{}
	session := dbSession.Copy()
	defer session.Close()

	c := session.DB(config.Get().MongoDatabase).C(config.Get().MongoCredentialsCollection)
	err = c.Find(nil).All(&result)

	if err != nil {
		return nil, err
	}
	return result, nil
}

//GetBoletoByDocumentNumber busca um boleto pelo nosso numero
//O retorno será um objeto BoletoView, o tempo decorrido da operação (em milisegundos) e algum erro ocorrido durante a operação
//Este método foi criado para fins de teste
func (e *MongoDb) GetBoletoByDocumentNumber(docNum string) (models.BoletoView, error) {
	e.m.Lock()
	defer e.m.Unlock()
	result := models.BoletoView{}

	session := dbSession.Copy()

	defer session.Close()

	c := session.DB(config.Get().MongoDatabase).C(config.Get().MongoBoletoCollection)

	for i := 0; i <= config.Get().RetryNumberGetBoleto; i++ {
		c.Find(bson.M{"boleto.title.documentnumber": docNum}).One(&result)

		if opErr, ok := err.(*net.OpError); ok && opErr.Timeout() {
			continue
		} else {
			break
		}
	}

	if err != nil {
		return models.BoletoView{}, err
	}

	return result, nil
}

//Close Fecha a conexão
func (e *MongoDb) Close() {
	fmt.Println("Close Database Connection")
}

func hasValidKey(r models.BoletoView, pk string) bool {
	return r.SecretKey == "" || r.PublicKey == pk
}
