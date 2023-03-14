package publisher

import (
	"context"
	"encoding/json"
	"log"
	"sync"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/pkg/errors"
)

//go:generate mockgen --source=./publisher.go -destination=./publisher_mocks_test.go -package=rabbitmqpub_test

type ConnectorPublisher interface {
	Publish(topic string, qos byte, retained bool, payload interface{}) mqtt.Token
}

type ConfigPublisher struct {
	Name string
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

func NewPublisher(ctx context.Context, name string, ch ConnectorPublisher) (*Publisher, error) {
	enity := &Publisher{
		config: ConfigPublisher{
			Name: name,
		},
		conn: ch,
		name: name,
	}
	return enity, nil
}

// SendMessage publish message to exchange
func (p *Publisher) SendMessage(ctx context.Context, topic string, message interface{}) error {
	defer p.checkError()
	log.Printf("%s", "Get to Send Message")
	body, err := json.Marshal(message)
	if err != nil {
		return errors.Wrap(err, "marshal message")
	}
	// We try to send message twice. Between attempts we try to reconnect.
	if err := p.sendMessage(ctx, topic, body); err != nil {
		if errRetryPub := p.sendMessage(ctx, topic, body); err != nil {
			return errors.Wrap(errRetryPub, "retry publish a message")
		}
	}
	return nil
}

func (p *Publisher) sendMessage(ctx context.Context, topic string, payload interface{}) (err error) {
	var token mqtt.Token
	if token = p.conn.Publish(
		topic,
		1,
		false,
		payload,
	); token.Error() != nil {
		log.Printf("%v", token.Error())
		p.muConn.Lock()
		p.isConnected = false
		p.muConn.Unlock()
		return token.Error()
	}
	log.Printf("%v", payload)
	log.Printf("%+v", token)
	return err
}

func (p *Publisher) checkError() (err error) {
	if r := recover(); r != nil {
		err = r.(error)
		return err
	}
	return err
}
