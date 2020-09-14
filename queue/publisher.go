package queue

import (
	"github.com/streadway/amqp"
)

//PublisherInterface Interface do Publicador
type PublisherInterface interface {
	GetExchangeName() string
	GetQueueName() string
	GetRoutingKey() string
	GetMessageToPublish() []byte
}

//WriteMessage Publica um messagem na fila
func WriteMessage(queuePublisher PublisherInterface) bool {

	var channel *amqp.Channel
	var queue amqp.Queue

	err := OpenConnection()
	if err != nil {
		return false
	}

	channel, err = openChannel(GetConnection(), "WriteMessage")
	if err != nil {
		return false
	}

	defer closeConnection(GetConnection(), "WriteMessage")
	defer closeChannel(channel, "WriteMessage")

	if exchangeDeclare(channel, queuePublisher.GetExchangeName(), "topic") &&
		queueDeclare(channel, queue, queuePublisher.GetQueueName()) &&
		queueBinding(channel, queuePublisher.GetQueueName(), queuePublisher.GetExchangeName(), queuePublisher.GetRoutingKey()) {
		return writeMessage(channel, queuePublisher) == nil
	}
	return false
}
