package usermanagement

import (
	"github.com/mundipagg/boleto-api/models"
	"gopkg.in/mgo.v2/bson"
)

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
