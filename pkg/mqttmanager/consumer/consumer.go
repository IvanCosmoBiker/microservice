package consumer

import (
	"context"
	"fmt"
	"log"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	"golang.org/x/sync/errgroup"
)

//go:generate mockgen --source=./consumer.go -destination=./consumer_mocks_test.go -package=rabbitmqconsum_test
type ConnectorConsumer interface {
	Subscribe(topic string, qos byte, callback mqtt.MessageHandler) mqtt.Token
}

type ConfigConsumer struct {
	Name, Topic string
}

// Consumer is a RabbitConsumer
type Consumer struct {
	config            ConfigConsumer
	conn              ConnectorConsumer
	Reconnect         chan int
	SubscriderMessage Subscriber
	ctx               context.Context
}

func NewConsumer(ctx context.Context, name, topic string, ch ConnectorConsumer) (*Consumer, error) {
	c := &Consumer{
		config: ConfigConsumer{
			Name:  name,
			Topic: topic,
		},
		conn:      ch,
		Reconnect: make(chan int),
		ctx:       ctx,
	}
	log.Println("New Consumer")
	return c, nil
}

func (c *Consumer) GetMassage(client mqtt.Client, message mqtt.Message) {
	fmt.Printf("%s", "Start Message")
	fmt.Printf("%s", message.Payload())
	err := c.SubscriderMessage.Consume(c.ctx, message.Payload())
	if err != nil {
		log.Printf("%v", err)
	}
}

// Subscriber describes interface with methods for subscriber
type Subscriber interface {
	Consume(ctx context.Context, data []byte) error
	Shutdown(ctx context.Context) error
}

func (c *Consumer) subscribe(ctx context.Context, errorGroup *errgroup.Group, subscriber Subscriber) (err error) {
	defer c.checkError()
	c.SubscriderMessage = subscriber
	token := c.conn.Subscribe(c.config.Topic, 1, c.GetMassage)
	for {
		if token.Wait() && token.Error() != nil {
			fmt.Println(token.Error())
			return token.Error()
		}
		select {
		case <-ctx.Done():
			log.Println("connection watcher stopped")
			if err := subscriber.Shutdown(ctx); err != nil {
				log.Println("shutdown handler")
				return err
			}
			return ctx.Err()
		}
	}
}

// Subscribe to channel for receiving message
func (c *Consumer) Subscribe(ctx context.Context, errorGroup *errgroup.Group, subscriber Subscriber) error {
	log.Println("New Subscribe")
	errorGroup.Go(func() error {
		return c.subscribe(ctx, errorGroup, subscriber)
	})
	return nil
}

func (c *Consumer) checkError() (err error) {
	if r := recover(); r != nil {
		err = r.(error)
		return err
	}
	return err
}
