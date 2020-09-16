package queue

import (
	"crypto/tls"
	"errors"
	"fmt"
	"github.com/mundipagg/boleto-api/config"
	"strconv"
	"time"

	"github.com/streadway/amqp"
)

func openChannel(conn *amqp.Connection, op string) (*amqp.Channel, error) {
	if conn != nil {
		channel, err := conn.Channel()
		failOnError(err, "Failed to open a channel", op)

		if err == nil {
			err = channel.Qos(1, 0, false)
			failOnError(err, "Error setting qos", op)
		}
		return channel, err
	}
	err := errors.New("Failed to open a channel. The connection is closed")
	failOnError(err, "Failed to open a channel. The connection is closed", op)
	return nil, err
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
	failOnError(err, "Failed to publish a message", "WriteMessage")
	return err
}

func openConnection() (*amqp.Connection, error) {

	hb, err := strconv.Atoi(config.Get().Heartbeat)

	if err != nil {
		hb = 60
	}

	conn, err := amqp.DialConfig(config.Get().ConnQueue, amqp.Config{
		Heartbeat:       time.Duration(hb) * time.Second,
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	})

	return conn, err
}

func closeConnection(conn *amqp.Connection, op string) {
	if conn != nil {
		err := conn.Close()
		failOnError(err, "Failed to close connection to RabbitMQ", op)
	}
}
