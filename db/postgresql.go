package db

import (
	"database/sql"
	"encoding/json"
	"fmt"

	"github.com/mundipagg/boleto-api/config"
	"github.com/mundipagg/boleto-api/models"
)

type postgresql struct {
}

//CreatePostgreSQL cria uma nova intancia de conexão com o PostgreSQL
func CreatePostgreSQL() (DB, error) {
	db := new(postgresql)
	return db, nil
}

func createConnection() (*sql.DB, error) {
	cfg := config.Get()
	conString := fmt.Sprintf(
		"user=%s password=%s dbname=%s host=%s port=%d sslmode=%s",
		cfg.PostgresUser, cfg.PostgresPassword, cfg.PostgresDBName, cfg.PostgresHost, cfg.PostgresPort, cfg.PostgresSSLMode)
	db, err := sql.Open("postgres", conString)
	if err != nil {
		return nil, err
	}
	return db, db.Ping()
}

//SaveBoleto salva um boleto no PostgreSQL
func (e *postgresql) SaveBoleto(boleto models.BoletoView) error {
	con, err := createConnection()
	if err != nil {
		return models.NewInternalServerError(err.Error(), "Falha ao conectar com o banco de dados")
	}
	defer con.Close()
	_, err = con.Exec("INSERT INTO boletos(id, boleto) VALUES ($1, $2)", boleto.ID, boleto.ToJSON())
	return err
}

//GetBoletoById busca um boleto pelo ID que vem na URL
func (e *postgresql) GetBoletoByID(id string) (models.BoletoView, error) {
	var rid, rjson string
	var result models.BoletoView
	con, err := createConnection()
	if err != nil {
		return result, models.NewInternalServerError(err.Error(), "Falha ao conectar com o banco de dados")
	}
	defer con.Close()
	row := con.QueryRow("SELECT id, boleto FROM boletos WHERE id = $1 ", id)
	err = row.Scan(&rid, &rjson)
	if err != nil {
		return result, err
	}
	err = json.Unmarshal([]byte(rjson), &result)
	return result, err
}

//Close fecha a conexão
func (e *postgresql) Close() {
	fmt.Println("Close Database Connection")
}
