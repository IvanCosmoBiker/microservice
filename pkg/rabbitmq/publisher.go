package rabbitmq

import (
	"amqp"
	"context"
	"encoding/json"
	"log"
	"sync"

	"github.com/pkg/errors"
)

//go:generate mockgen --source=./publisher.go -destination=./publisher_mocks_test.go -package=rabbitmqpub_test

type ConnectorPublisher interface {
	ExchangeDeclare(name, kind string, durable, autoDelete, internal, noWait bool, args amqp.Table) error
	Publish(exchange, key string, mandatory, immediate bool, msg amqp.Publishing) error
}

type ConfigPublisher struct {
	ExchangeName string
}

// Publisher is a RabbitPublisher
type Publisher struct {
	config      ConfigPublisher
	conn        ConnectorPublisher
	isConnected bool
	name        string
	muConn      sync.Mutex

	okMessages  func(ctx context.Context) error
	badMessages func(ctx context.Context) error
}

func NewPublisher(ctx context.Context, name, exchangeName string, ch ConnectorPublisher) (*Publisher, error) {
	if len(exchangeName) <= 0 {
		exchangeName = "amq.topic"
	}
	enity := &Publisher{
		config: ConfigPublisher{
			ExchangeName: exchangeName,
		},
		conn: ch,
		name: name,
	}
	return enity, nil
}

func (p *Publisher) connect(_ context.Context) error {
	p.muConn.Lock()
	defer p.muConn.Unlock()
	if p.isConnected {
		return nil
	}
	if len(p.config.ExchangeName) > 0 {
		if err := p.conn.ExchangeDeclare(p.config.ExchangeName, "direct", true,
			false, false,
			false, nil); err != nil {
			return errors.Wrap(err, "declare a exchange")
		}
	}

	p.isConnected = true
	return nil
}

// SendMessage publish message to exchange
func (p *Publisher) SendMessage(ctx context.Context, RoutingKey string, message interface{}) error {
	body, err := json.Marshal(message)
	if err != nil {
		return errors.Wrap(err, "marshal message")
	}
	ampqMsg := buildMessage(body)
	log.Printf("send message: %s", string(body))
	log.Printf("%v", p.isConnected)
	if !p.isConnected {
		if err := p.connect(ctx); err != nil {
			log.Printf("%v", err)
			log.Println("connect publisher to rabbitMQ")
		}
	}
	// We try to send message twice. Between attempts we try to reconnect.
	if err := p.sendMessage(ctx, RoutingKey, ampqMsg); err != nil {
		if errRetryPub := p.sendMessage(ctx, RoutingKey, ampqMsg); err != nil {
			return errors.Wrap(errRetryPub, "retry publish a message")
		}
	}
	log.Printf("%s", RoutingKey)
	return nil
}

func (p *Publisher) sendMessage(ctx context.Context, RoutingKey string, ampqMsg *amqp.Publishing) error {
	if !p.isConnected {
		if err := p.connect(ctx); err != nil {
			log.Println("connect publisher to rabbitMQ")
		}
	}
	log.Printf("%s", RoutingKey)
	if err := p.conn.Publish(
		p.config.ExchangeName,
		RoutingKey,
		false,
		false,
		*ampqMsg,
	); err != nil {
		log.Printf("%v", err)
		p.muConn.Lock()
		p.isConnected = false
		p.muConn.Unlock()
		return errors.Wrap(err, "publish a message")
	}
	return nil
}

func buildMessage(body []byte) *amqp.Publishing {
	return &amqp.Publishing{
		ContentType: "application/json",
		Body:        body,
	}
}
