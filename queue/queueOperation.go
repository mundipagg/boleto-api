package queue

import (
	"fmt"
	"github.com/mundipagg/boleto-api/log"

	"github.com/streadway/amqp"
)

var (
	conn *amqp.Connection
	l    *log.Log
)

// OpenConnectionReadOnly Abre conexão de leitura com RabbitMQ
func OpenConnection() error {
	l = log.CreateLog()
	var err error
	conn, err = openConnection()

	return err
}

// GetConnectionReadOnly Obtém conexão de leitura com RabbitMQ
func GetConnection() *amqp.Connection {
	return conn
}

func failOnError(err error, title string, op string, msg ...string) {
	if err != nil {
		if msg != nil && msg[0] != "" {
			l.Warn(msg[0], fmt.Sprintf("%s - %s - %s", op, title, err))
		} else {
			l.Error(err, fmt.Sprintf("%s - %s", op, title))
		}
	}
}
