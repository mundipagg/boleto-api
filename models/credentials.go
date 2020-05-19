package models

//Credentials Credenciais para requisição de Registro de Boleto
type Credentials struct {
	Username string `bson:"username"`
	Password string `bson:"password"`
}

//NewCredentials Cria uma instância de Credential
func NewCredentials(u, p string) *Credentials {
	return &Credentials{
		Username: u,
		Password: p,
	}
}
