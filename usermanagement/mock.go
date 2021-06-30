package usermanagement

import (
	"github.com/mundipagg/boleto-api/models"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

//LoadMockUserCredentials Cria credenciais do mock
func LoadMockUserCredentials() (string, string) {
	c := models.Credentials{
		ID:       primitive.NewObjectID(),
		Username: "user",
		Password: "pass",
	}
	uk := c.ID.Hex()
	c.UserKey = string(uk)
	addUser(c.UserKey, c)
	return c.UserKey, c.Password
}
