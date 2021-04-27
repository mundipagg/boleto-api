package queue

import (
	"crypto/tls"
	"errors"
	"fmt"
	"strconv"
	"time"

	"github.com/mundipagg/boleto-api/config"
	"github.com/mundipagg/boleto-api/util"

	"github.com/streadway/amqp"
)

const HEARTBEAT_DEFAULT = 10

func openChannel(conn *amqp.Connection, op string) (*amqp.Channel, error) {
	var channel *amqp.Channel
	var err error

	if conn == nil {
		err = errors.New("failed to open a channel. The connection is closed")
		failOnError(err, "Failed to open a channel. The connection is closed", op)
		return nil, err
	}

	if channel, err = conn.Channel(); err != nil {
		failOnError(err, "Failed to open a channel", op)
		return nil, err
	}

	if err = channel.Confirm(false); err != nil {
		failOnError(err, "Failed to set channel into Confirm Mode", op)
		return nil, err
	}

	if err = channel.Qos(1, 0, false); err != nil {
		failOnError(err, "Error setting qos", op)
		return nil, err
	}

	return channel, err
}

func closeChannel(channel *amqp.Channel, op string) {
	if channel != nil {
		err := channel.Close()
		failOnError(err, "Failed to close a channel", op)
	}
}

func exchangeDeclare(channel *amqp.Channel, exchange, kind string) bool {
	err := channel.ExchangeDeclare(exchange, kind, true, false, false, false, nil)
	hdr := fmt.Sprintf("[{Application}: {Operation}] - Error Declaring RabbitMQ Exchange %s", exchange)
	failOnError(err, hdr, "ExchangeDeclare")
	return err == nil
}

func queueDeclare(channel *amqp.Channel, q amqp.Queue, queueName string) bool {

	q, err := channel.QueueDeclare(queueName, true, false, false, false, nil)
	hdr := fmt.Sprintf("[{Application}: {Operation}] - Error Declaring RabbitMQ Queue %s", queueName)
	failOnError(err, hdr, "QueueDeclare")
	return err == nil
}

func queueBinding(channel *amqp.Channel, queue, exchange, key string) bool {
	err := channel.QueueBind(queue, key, exchange, false, nil)
	hdr := fmt.Sprintf("[{Application}: {Operation}] - Error Binding RabbitMQ Queue %s into Exchange %s", queue, exchange)
	failOnError(err, hdr, "QueueBinding")
	return err == nil
}

func writeMessage(channel *amqp.Channel, p PublisherInterface) error {
	notifyConfirm := make(chan amqp.Confirmation)
	channel.NotifyPublish(notifyConfirm)

	err := channel.Publish(
		p.GetExchangeName(),
		p.GetRoutingKey(), // queue
		false,             // mandatory
		false,             // immediate
		amqp.Publishing{
			DeliveryMode: amqp.Persistent,
			ContentType:  "text/plain, charset=UTF-8",
			Body:         p.GetMessageToPublish(),
		})

	if err == nil {
		if confirm := <-notifyConfirm; !confirm.Ack {
			err = errors.New("nack received from the server during message posting")
		}
	}

	failOnError(err, "Failed to publish a message", "WriteMessage")
	return err
}

func openConnection() (*amqp.Connection, error) {

	var hb int
	var err error

	if hb, err = strconv.Atoi(config.Get().Heartbeat); err != nil {
		hb = HEARTBEAT_DEFAULT
	}

	conn, err := amqp.DialConfig(config.Get().ConnQueue, amqp.Config{
		Heartbeat: time.Duration(hb) * time.Second,
		TLSClientConfig: &tls.Config{
			MinVersion:         util.GetTLSVersion(config.Get().QueueMinTLS),
			MaxVersion:         util.GetTLSVersion(config.Get().QueueMaxTLS),
			InsecureSkipVerify: config.Get().QueueByPassCertificate,
		},
	})

	return conn, err
}

func closeConnection(conn *amqp.Connection, op string) {
	if conn != nil {
		err := conn.Close()
		failOnError(err, "Failed to close connection to RabbitMQ", op)
	}
}
