package usermanagement

import (
	"fmt"
	"sync"

	"github.com/mundipagg/boleto-api/db"
	"github.com/mundipagg/boleto-api/log"
	"github.com/mundipagg/boleto-api/models"
	"gopkg.in/mgo.v2/bson"
)

var userCredentialStorage = sync.Map{}

func addUser(key string, value interface{}) {
	userCredentialStorage.Store(key, value)
}

//GetUser Busca credenciais de um usuário
func GetUser(key string) (interface{}, bool) {
	if value, ok := userCredentialStorage.Load(key); ok {
		return value, true
	}
	return nil, false
}

//LoadUserCredentials Carrega credenciais salvas no banco de dados
func LoadUserCredentials() {
	log := log.CreateLog()

	mongo, errMongo := db.CreateMongo(log)
	if errMongo != nil {
		log.Error(errMongo.Error(), "Error in connection to MongoDB")
		return
	}

	c, err := mongo.GetUserCredentials()
	if err != nil {
		log.Error(err.Error(), fmt.Sprintf("Error in get user crendentials - %s", err.Error()))
		return
	}

	for _, u := range c {
		uk, _ := u.ID.MarshalText()
		u.UserKey = string(uk)
		addUser(u.UserKey, u)
	}
}

//LoadMockUserCredentials Cria credenciais do mock
func LoadMockUserCredentials() (string, string) {
	c := models.Credentials{
		ID:       bson.NewObjectId(),
		Username: "user",
		Password: "pass",
	}
	uk, _ := c.ID.MarshalText()
	c.UserKey = string(uk)
	addUser(c.UserKey, c)
	return c.UserKey, c.Password
}
