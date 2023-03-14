package rabbitmq

import (
	"amqp"
	"context"
	"log"
	"time"

	"github.com/pkg/errors"
	"golang.org/x/sync/errgroup"
)

//go:generate mockgen --source=./consumer.go -destination=./consumer_mocks_test.go -package=rabbitmqconsum_test
type ConnectorCunsumer interface {
	ExchangeDeclare(name, kind string, durable, autoDelete, internal, noWait bool, args amqp.Table) error
	QueueDeclare(name string, durable, autoDelete, exclusive, noWait bool, args amqp.Table) (amqp.Queue, error)
	QueueBind(name, key, exchange string, noWait bool, args amqp.Table) error
	Consume(queue, consumer string, autoAck, exclusive, noLocal, noWait bool, args amqp.Table) (<-chan amqp.Delivery, error)
	OrinChannelReconnect() chan int
	IsClosed() bool
}

type ConfigConsumer struct {
	Name, ExchangeName, RoutingKey, ConsumerName string
}

// Consumer is a RabbitConsumer
type Consumer struct {
	config    ConfigConsumer
	conn      ConnectorCunsumer
	Reconnect chan int
}

func NewConsumer(ctx context.Context, name, consumerName, exchangeName, routingKey string, ch ConnectorCunsumer) (*Consumer, error) {
	c := &Consumer{
		config: ConfigConsumer{Name: name,
			ConsumerName: consumerName,
			ExchangeName: exchangeName,
			RoutingKey:   routingKey,
		},
		conn:      ch,
		Reconnect: make(chan int),
	}
	log.Println("New Consumer")
	return c, nil
}

func (c *Consumer) connect(_ context.Context) (<-chan amqp.Delivery, error) {
	var err error
	var msg <-chan amqp.Delivery
	exhangeName := c.config.ExchangeName
	if len(exhangeName) > 0 {
		if err = c.conn.ExchangeDeclare(c.config.ExchangeName, "direct", true,
			false, false,
			false, nil); err != nil {
			log.Printf("%v", err)
			return nil, errors.Wrap(err, "declare a exchange")
		}
	} else {
		exhangeName = "amq.topic"
	}
	_, err = c.conn.QueueDeclare(
		c.config.Name, // name
		true,          // durable
		false,         // delete when unused
		false,         // exclusive
		false,         // no-wait
		nil,           // arguments
	)
	if err != nil {
		log.Printf("%v", err)
		return nil, errors.Wrap(err, "QueueDeclare")
	}
	if err = c.conn.QueueBind(
		c.config.Name,       // queue name
		c.config.RoutingKey, // routing key
		exhangeName,         // exchange
		false,
		nil,
	); err != nil {
		log.Printf("%v", err)
		return nil, errors.Wrap(err, "bind to queue")
	}
	msg, err = c.conn.Consume(
		c.config.Name,         // queue
		c.config.ConsumerName, // consume
		false,                 // auto-ack
		false,                 // exclusive
		false,                 // no-local
		false,                 // no-wait
		nil,                   // args
	)
	if err != nil {
		log.Printf("%v", err)
		return nil, errors.Wrap(err, "consume message")
	}
	return msg, nil
}

// Subscriber describes interface with methods for subscriber
type Subscriber interface {
	Consume(ctx context.Context, data []byte) error
	Shutdown(ctx context.Context) error
}

func (c *Consumer) subscribe(ctx context.Context, errorGroup *errgroup.Group, subscriber Subscriber) error {
	var msg <-chan amqp.Delivery
	var err error
	for {
		if msg, err = c.connect(ctx); err != nil {
			log.Println("connect consumer to rabbitMQ")
			time.Sleep(10 * time.Second)
			continue
		}
		break
	}
	log.Println("consumer connected")
	for {
		select {
		case <-ctx.Done():
			log.Println("connection watcher stopped")
			if err := subscriber.Shutdown(ctx); err != nil {
				log.Println("shutdown handler")
				return err
			}
			return ctx.Err()
		case d, ok := <-msg:
			if ok {
				log.Printf("got new event %+v", string(d.Body))
				if errConsume := subscriber.Consume(ctx, d.Body); errConsume != nil {
					log.Println(errConsume)
				}
				if err := d.Ack(true); err != nil {
					log.Println("ack")
				}
			} else {
				log.Printf("%#v", c.conn)
				for {
					select {
					case <-ctx.Done():
						{
							log.Println("connection watcher stopped")
							if err := subscriber.Shutdown(ctx); err != nil {
								log.Println("shutdown handler")
								return err
							}
							return ctx.Err()
						}
					case _, ok := <-c.Reconnect:
						{
							if !ok {
								return nil
							}
							return c.subscribe(ctx, errorGroup, subscriber)
						}
					}
				}
				log.Println("RESTART CONSUMER")
			}
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
