package queue

import (
	"fmt"
	"github.com/mundipagg/boleto-api/config"
	"github.com/mundipagg/boleto-api/log"

	"strconv"
	"time"

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

	if err == nil {
		watchConnection(conn)
	}
	return err
}

func CloseConnection() {
	closeConnection(conn, "CloseConnection")
}

// GetConnectionReadOnly Obtém conexão de leitura com RabbitMQ
func GetConnection() *amqp.Connection {
	return conn
}

func watchConnection(c *amqp.Connection) {
	errs := make(chan *amqp.Error, 1)
	c.NotifyClose(errs)

	go func() {
		for err := range errs {
			if err != nil {

				l.Warn(err, "WatchConnection - AMQP connection is closing")

				rErr := OpenConnection()
				for rErr != nil {
					t, _ := strconv.Atoi(config.Get().TimeToRecoveryWithQueueInSeconds)
					l.Warn(rErr, fmt.Sprintf("WatchConnection - Retry write connection into %d seconds", t))
					time.Sleep(time.Second * time.Duration(t))
				}
			}
		}
	}()
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
