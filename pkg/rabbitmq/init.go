package rabbitmq

import (
	"amqp"
	"context"
	config "ephorservices/config"
	loggerEvent "ephorservices/pkg/logger"
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/pkg/errors"
	"golang.org/x/sync/errgroup"
)

type ChannelPoolItemKey struct {
	Queue    string
	Consumer string
	Exchange string
	Key      string
}

type ItemChannelPool struct {
	isUsed  bool
	Channel *amqp.Channel
}

type Connection struct {
	EventLogger          *loggerEvent.LoggerEvent
	Dsn                  string
	BackoffPolicy        map[string]int
	Conn                 *amqp.Connection
	ServiceChannel       *amqp.Channel
	Mu                   sync.RWMutex
	ChannelPoolPublisher map[int8]*ItemChannelPool
	ChannelPoolConsumer  map[ChannelPoolItemKey]*amqp.Channel
	ChannelPoolMu        sync.RWMutex
	isClosed             bool
	ErrorClose           bool
	ErrorGroup           *errgroup.Group
	ChanCtx              context.Context
	ReconnectRabbit      chan int
	cfg                  *config.Config
}

func (c *Connection) defaultBackoffPolicy() map[string]int {
	backoffPolicyDefault := make(map[string]int, 2)
	backoffPolicyDefault["time"] = 10
	backoffPolicyDefault["count"] = 400
	return backoffPolicyDefault
}

func (c *Connection) setStringDsn(login string, password string, address string, port string, maxAttempts int) string {
	return fmt.Sprintf("amqp://%s:%s@%s:%s", login, password, address, port)
}

func (c *Connection) SetConfig() {
	c.Dsn = c.setStringDsn(c.cfg.RabbitMq.Login, c.cfg.RabbitMq.Password, c.cfg.RabbitMq.Address, c.cfg.RabbitMq.Port, c.cfg.RabbitMq.MaxAttempts)
	c.BackoffPolicy = c.defaultBackoffPolicy()
}

func NewConnection(cfg *config.Config) (*Connection, error) {
	Conn := &Connection{
		isClosed:   true,
		ErrorClose: false,
		cfg:        cfg,
	}
	return Conn, nil
}

// OriConn returns original connection to rabbitmq
func (c *Connection) OriConn() *amqp.Connection {
	c.Mu.RLock()
	defer c.Mu.RUnlock()
	return c.Conn
}

// OriServiceChannel return original service channel
func (c *Connection) OrinServiceChannel() *amqp.Channel {
	c.Mu.RLock()
	defer c.Mu.RUnlock()
	return c.ServiceChannel
}

func (c *Connection) OrinChannelReconnect() chan int {
	c.Mu.RLock()
	defer c.Mu.RUnlock()
	return c.ReconnectRabbit
}

// Channel returns original channel to rabbitmq
func (c *Connection) Channel() (*amqp.Channel, error) {
	c.Mu.RLock()
	defer c.Mu.RUnlock()
	channel, err := c.Conn.Channel()
	if err != nil {
		return nil, errors.Wrap(err, "open a channel")
	}
	return channel, nil
}

func (c *Connection) MakePoolChannel() (err error) {
	if c.cfg.RabbitMq.PoolChannel < int8(0) {
		return errors.New("not set PoolChannel")
	}
	c.ChannelPoolPublisher = make(map[int8]*ItemChannelPool)
	for i := int8(0); i < c.cfg.RabbitMq.PoolChannel; i++ {
		channel, errChannel := c.Conn.Channel()
		log.Printf("%v", i)
		if errChannel != nil {
			return errChannel
		}
		itemChannel := &ItemChannelPool{
			Channel: channel,
			isUsed:  false,
		}
		c.ChannelPoolPublisher[i] = itemChannel
		go c.NotifyPublish(itemChannel.Channel)
	}
	return err
}

func (c *Connection) CloseChannel(ch *ItemChannelPool, key int8) (err error) {
	log.Printf("%s", "Channel Close")
	err = ch.Channel.Close()
	delete(c.ChannelPoolPublisher, key)
	return err
}

func (c *Connection) OpenChannel() (err error) {
	if int8(len(c.ChannelPoolPublisher)) == c.cfg.RabbitMq.PoolChannel {
		return nil
	}
	var ch *amqp.Channel
	ch, err = c.Channel()
	if err != nil {
		return err
	}
	c.ChannelPoolPublisher[int8(len(c.ChannelPoolPublisher))+1] = &ItemChannelPool{
		Channel: ch,
		isUsed:  false,
	}
	log.Printf("%s", "Channel Open And Add to map")
	return nil
}

func (c *Connection) Close(_ context.Context) error {
	c.Mu.Lock()
	defer c.Mu.Unlock()
	c.isClosed = true
	for _, ch := range c.ChannelPoolPublisher {
		if err := ch.Channel.Close(); err != nil {
			return errors.Wrap(err, "close rabbitMQ channel")
		}
		log.Printf("close Channel Ok - %s", "publisher")
	}
	for key, ch := range c.ChannelPoolConsumer {
		if err := ch.Close(); err != nil {
			return errors.Wrap(err, "close rabbitMQ channel")
		}
		log.Printf("close Channel Ok - %s", key.Consumer)
	}

	if err := c.Conn.Close(); err != nil {
		return errors.Wrap(err, "close rabbitMQ connection")
	}
	log.Println("close rabbitMQ connection OK")
	return nil
}

func (c *Connection) IsClosed() bool {
	return c.isClosed
}

func (c *Connection) connect(ctx context.Context) error {
	var err error
	if c.Conn, err = amqp.Dial(c.Dsn); err != nil {
		return errors.Wrap(err, "connect to rabbitMQ")
	}
	log.Println("get Connection")
	c.isClosed = false
	c.ErrorClose = false
	if c.ServiceChannel, err = c.Conn.Channel(); err != nil {
		return errors.Wrap(err, "create service rabbitMQ channel")
	}
	err = c.MakePoolChannel()
	if err != nil {
		return err
	}
	c.ChannelPoolConsumer = make(map[ChannelPoolItemKey]*amqp.Channel)
	log.Println("get Channel")
	return nil
}

// Connect auto reconnect to rabbitmq when we lost connection.
func (c *Connection) Connect(ctx context.Context, errorGroup *errgroup.Group) error {
	var connErr error
	c.ReconnectRabbit = make(chan int)
	c.ErrorGroup = errorGroup
	c.ChanCtx = ctx
	if c.isClosed {
		if err := c.connect(ctx); err != nil {
			return errors.Wrap(err, "connect")
		}
	}
	c.ErrorGroup.Go(func() error {
		log.Println("starting connection watcher")
		for {
			select {
			case <-ctx.Done():
				log.Println("connection watcher stopped")
				return ctx.Err()
			default:
				_, ok := <-c.Conn.NotifyClose(make(chan *amqp.Error))
				if !ok {
					c.ErrorClose = true
					if c.isClosed {
						return nil
					}
					log.Println("rabbitMQ connection unexpected closed")
					c.Mu.Lock()
					for i := 0; i < c.BackoffPolicy["count"]; i++ {
						if connErr = c.connect(ctx); connErr != nil {
							log.Println("connection failed, trying to reconnect to rabbitMQ")
							time.Sleep(time.Duration(c.BackoffPolicy["time"]) * time.Second)
							continue
						}
						break
					}
					c.Mu.Unlock()
					c.ReconnectRabbit <- 1
				}
			}
		}
	})
	return nil
}

func (c *Connection) ExchangeDeclare(name, kind string, durable, autoDelete, internal, noWait bool, args amqp.Table) error {
	c.Mu.RLock()
	defer c.Mu.RUnlock()
	return c.ServiceChannel.ExchangeDeclare(name, kind, durable, autoDelete, internal, noWait, args)
}

func (c *Connection) QueueDeclare(name string, durable, autoDelete, exclusive, noWait bool, args amqp.Table) (amqp.Queue, error) {
	c.Mu.RLock()
	defer c.Mu.RUnlock()
	return c.ServiceChannel.QueueDeclare(name, durable, autoDelete, exclusive, noWait, args)
}

func (c *Connection) QueueBind(name, key, exchange string, noWait bool, args amqp.Table) error {
	c.Mu.RLock()
	defer c.Mu.RUnlock()
	return c.ServiceChannel.QueueBind(name, key, exchange, noWait, args)
}

func (c *Connection) Consume(
	queue, consumer string,
	autoAck, exclusive, noLocal, noWait bool,
	args amqp.Table) (<-chan amqp.Delivery, error) {
	c.Mu.RLock()
	defer c.Mu.RUnlock()
	ch, err := c.GetChannelFromPoolConsumer("", "", queue, consumer)
	if err != nil {
		return nil, errors.Wrap(err, "get channel from pool")
	}
	return ch.Consume(queue, consumer, autoAck, exclusive, noLocal, noWait, args)
}

// nolint:gocritic // pass msg without pointer as in original func in amqp
func (c *Connection) Publish(exchange, keyRouting string, mandatory, immediate bool, msg amqp.Publishing) (err error) {
	c.Mu.RLock()
	defer c.Mu.RUnlock()
	var ch *ItemChannelPool
	var key int8
	ch, key, err = c.GetChannelFromPoolPublischer()
	log.Printf("%v", err)
	if err != nil && err.Error() == "no free channels" {
		ch, key, err = c.WaitChannel()
	}
	if err != nil {
		return err
	}
	log.Printf("%v", *ch)
	err = ch.Channel.Publish(exchange, keyRouting, mandatory, immediate, msg)
	if err != nil {
		c.ReturnChannelInPoolPublischer(ch, key)
		return err
	}
	c.ReturnChannelInPoolPublischer(ch, key)
	log.Printf("%v", *ch)
	return nil
}

func (c *Connection) WaitChannel() (ch *ItemChannelPool, key int8, err error) {
	key = 0
	log.Println("start waiting")
	for {
		select {
		case <-c.ChanCtx.Done():
			log.Println("channel watcher stopped")
			return ch, key, c.ChanCtx.Err()
		default:
			number := len(c.ChannelPoolPublisher)
			if number > 0 {
				log.Println("we have Channel")
				if c.isClosed {
					return ch, key, errors.New("connection is closed")
				}
				ch, key, err = c.GetChannelFromPoolPublischer()
				return ch, key, err
			}
		}
	}
	return ch, key, err
}

func (c *Connection) GetChannelFromPoolPublischer() (channel *ItemChannelPool, key int8, err error) {
	c.ChannelPoolMu.Lock()
	defer c.ChannelPoolMu.Unlock()
	key = 0
	closeResult := c.Conn.IsClosed()
	if closeResult != false {
		log.Println("Connection is close")
		return nil, key, errors.New("Connection is close")
	}
	for key, ch := range c.ChannelPoolPublisher {
		if !ch.isUsed {
			ch.isUsed = true
			return ch, key, nil
		}
	}
	return nil, key, errors.New("no free channels")
}

func (c *Connection) ReturnChannelInPoolPublischer(channel *ItemChannelPool, key int8) {
	c.ChannelPoolMu.Lock()
	defer c.ChannelPoolMu.Unlock()
	c.CloseChannel(channel, key)
	c.OpenChannel()
}

func (c *Connection) GetChannelFromPoolConsumer(exchange, key, queue, consumer string) (*amqp.Channel, error) {
	c.ChannelPoolMu.Lock()
	defer c.ChannelPoolMu.Unlock()
	closeResult := c.Conn.IsClosed()
	if closeResult == true {
		log.Println("Connection is close")
		return nil, errors.New("Connection is close")
	}

	var err error
	poolKey := ChannelPoolItemKey{
		Exchange: exchange,
		Key:      key,
		Queue:    queue,
		Consumer: consumer,
	}
	ch, ok := c.ChannelPoolConsumer[poolKey]
	if !ok {
		ch, err = c.Conn.Channel()
		if err != nil {
			return nil, errors.Wrap(err, "create channel")
		}
		c.ChannelPoolConsumer[poolKey] = ch
		c.chanWatcherConsumer(poolKey)
	}

	return ch, nil
}

func (c *Connection) chanWatcherConsumer(poolKey ChannelPoolItemKey) {
	ch := c.ChannelPoolConsumer[poolKey]
	c.ErrorGroup.Go(func() error {
		log.Println("starting channel watcher")
		for {
			select {
			case <-c.ChanCtx.Done():
				log.Println("channel watcher stopped")
				return c.ChanCtx.Err()
			default:
				_, ok := <-ch.NotifyClose(make(chan *amqp.Error))
				if !ok {
					log.Println("Close Channel")
					if c.isClosed {
						return nil
					}
					log.Println("rabbitMQ channel unexpected closed")
					c.ChannelPoolMu.Lock()
					delete(c.ChannelPoolConsumer, poolKey)
					c.ChannelPoolMu.Unlock()
					return nil
				}
			}
		}
	})
}

func (c *Connection) NotifyPublish(channel *amqp.Channel) error {
	fmt.Printf("%+v", channel)
	for {
		select {
		case <-c.ChanCtx.Done():
			log.Println("channel NotifyPublish stopped")
			return c.ChanCtx.Err()
		default:
			con, ok := <-channel.NotifyPublish(make(chan amqp.Confirmation))
			fmt.Printf("$$$$$$%+v", con)
			if !ok {
				log.Println("Close Channel")
				log.Println("rabbitMQ channel unexpected closed")
				return nil
			}
		}
	}
}
