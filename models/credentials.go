package models

import "gopkg.in/mgo.v2/bson"

//Credentials Credenciais para requisição de Registro de Boleto
type Credentials struct {
	ID       bson.ObjectId `bson:"_id,omitempty"`
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
