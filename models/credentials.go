package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

//Credentials Credenciais para requisição de Registro de Boleto
type Credentials struct {
	ID       primitive.ObjectID `bson:"_id,omitempty"`
	UserKey  string
	Username string `bson:"username,omitempty"`
	Password string `bson:"password,omitempty"`
}

//NewCredentials Cria uma instância de Credential
func NewCredentials(k, p string) *Credentials {
	return &Credentials{
		UserKey:  k,
		Password: p,
	}
}
