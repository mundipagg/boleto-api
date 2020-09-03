package queue

import (
	"github.com/mundipagg/boleto-api/config"
)

//Publisher Implementação da interface publisher
type Publisher struct {
	ExchangeName string
	QueueName    string
	RoutingKey   string
	Message      string
}

func NewPublisher(message string) *Publisher {
	p := new(Publisher)
	p.ExchangeName = config.Get().OriginExchange
	p.QueueName = config.Get().OriginQueue
	p.RoutingKey = config.Get().OriginRoutingKey
	p.Message = message

	return p
}

//GetExchangeName Retorna o nome da fila
func (p *Publisher) GetExchangeName() string {
	return p.ExchangeName
}

//GetQueueName Retorna o nome da fila com sufixo identificador
func (p *Publisher) GetQueueName() string {
	return p.QueueName
}

//GetRoutingKey Retorna a RoutingKey para direcionamento da mensagem
func (p *Publisher) GetRoutingKey() string {
	return p.RoutingKey
}

//GetMessageToPublish Retorna a mensagem convertida para publicação na fila
func (p *Publisher) GetMessageToPublish() []byte {
	return []byte(p.Message)
}
