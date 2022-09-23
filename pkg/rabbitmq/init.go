package rabbitmq

import(
    "context"
	"sync"
	"fmt"
	"log"
	"amqp"
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

type Connection struct {
	Dsn            string
	BackoffPolicy  map[string]int
	Conn           *amqp.Connection
	ServiceChannel *amqp.Channel
	Mu             sync.RWMutex
	ChannelPool    map[ChannelPoolItemKey]*amqp.Channel
	ChannelPoolMu  sync.RWMutex
	isClosed       bool
	ErrorClose 	   bool
	ErrorGroup     *errgroup.Group
	ChanCtx        context.Context
	ReconnectRabbit chan int
}

func (c *Connection) defaultBackoffPolicy() map[string]int {
	backoffPolicyDefault := make(map[string]int,2)
	backoffPolicyDefault["time"] = 10 
	backoffPolicyDefault["count"] = 400
	return backoffPolicyDefault
}

func (c *Connection) setStringDsn(login string, password string, address string, port string,maxAttempts int) string {
    return fmt.Sprintf("amqp://%s:%s@%s:%s",login,password,address,port)
}

func (c *Connection) SetConfig(login string, password string, address string, port string,maxAttempts int){
    c.Dsn = c.setStringDsn(login, password, address, port, maxAttempts)
    c.BackoffPolicy = c.defaultBackoffPolicy()
}

func NewConnection() (*Connection, error) {
	Conn := &Connection{
		isClosed: true,
		ErrorClose: false,
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

func (c *Connection) Close(_ context.Context) error {
	c.Mu.Lock()
	defer c.Mu.Unlock()
	c.isClosed = true
	for key, ch := range c.ChannelPool {
		if err := ch.Close(); err != nil {
			return errors.Wrap(err, "close rabbitMQ channel")
		}
		log.Printf("close Channel Ok - %s",key.Consumer)
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
	c.ChannelPool = make(map[ChannelPoolItemKey]*amqp.Channel)
	log.Println("get Channel")
	return nil
}


// Connect auto reconnect to rabbitmq when we lost connection.
func (c *Connection) Connect(ctx context.Context, errorGroup *errgroup.Group) error {
	var connErr error
	c.ReconnectRabbit = make(chan int)
	if c.isClosed {
		if err := c.connect(ctx); err != nil {
			return errors.Wrap(err, "connect")
		}
	}
	c.ErrorGroup = errorGroup
	c.ChanCtx = ctx
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
					for i:= 0; i < c.BackoffPolicy["count"]; i++ {
						if connErr = c.connect(ctx); connErr != nil {
							log.Println("connection failed, trying to reconnect to rabbitMQ")
							time.Sleep(time.Duration(c.BackoffPolicy["time"]) * time.Second)
							continue
						}
						break
					}
					c.Mu.Unlock()
					c.ReconnectRabbit<-1
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
	ch, err := c.GetChannelFromPool("", "", queue, consumer)
	if err != nil {
		return nil, errors.Wrap(err, "get channel from pool")
	}
	return ch.Consume(queue, consumer, autoAck, exclusive, noLocal, noWait, args)
}

// nolint:gocritic // pass msg without pointer as in original func in amqp
func (c *Connection) Publish(exchange, key string, mandatory, immediate bool, msg amqp.Publishing) error {
	c.Mu.RLock()
	defer c.Mu.RUnlock()

	ch, err := c.GetChannelFromPool(exchange, key, "", "")
	if err != nil {
		return errors.Wrap(err, "get channel from pool")
	}

	return ch.Publish(exchange, key, mandatory, immediate, msg)
}

func (c *Connection) GetChannelFromPool(exchange, key, queue, consumer string) (*amqp.Channel, error) {
	c.ChannelPoolMu.Lock()
	closeResult := c.Conn.IsClosed()
	if closeResult == true {
		log.Println("Connection is close")
		return nil,errors.New("Connection is close")
	}
	defer c.ChannelPoolMu.Unlock()
	var err error
	poolKey := ChannelPoolItemKey{
		Exchange: exchange,
		Key:      key,
		Queue:    queue,
		Consumer: consumer,
	}
	ch, ok := c.ChannelPool[poolKey]
	if !ok {
		ch, err = c.Conn.Channel()
		if err != nil {
			return nil, errors.Wrap(err, "create channel")
		}
		c.ChannelPool[poolKey] = ch
		c.chanWatcher(poolKey)
	}

	return ch, nil
}

func (c *Connection) chanWatcher(poolKey ChannelPoolItemKey) {
	ch := c.ChannelPool[poolKey]
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
					delete(c.ChannelPool, poolKey)
					c.ChannelPoolMu.Unlock()
					return nil
				}
			}
		}
	})
}