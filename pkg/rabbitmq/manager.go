package rabbitmq

import (
	"log"
	"context"
	"github.com/pkg/errors"
	"amqp"
	"golang.org/x/sync/errgroup"
	config "ephorservices/config"
	"sync"
)

const ProviderName = "rabbitmq"

type Config struct {
	Login string
	Password string
	Address string
	Port string
	MaxAttempts int
	Timeout int
}


// Manager is a connection controlling structure. It controls
// connection, asynchronous queue and everything that related to
// specified connection.
type Manager struct {
	ConfigRabbit     *Config
	conn       		 *Connection
	ContextRabbit 	 context.Context
	consumers  		 map[string]*Consumer
	publishers 		 map[string]*Publisher
	Mu             	 sync.Mutex
}

// NewManager create new Manager and init Connection to Rabbit.
func NewManager(ctx context.Context,login, password, address, port string, maxAttempts int) (*Manager, error) {
	var manager = &Manager{
		consumers:         make(map[string]*Consumer),
		publishers:        make(map[string]*Publisher),
	}
	var config = &Config{}
	var err error
	config.Login = login
	config.Password = password
	config.Address = address
	config.Port = port
	config.MaxAttempts = maxAttempts
	manager.ConfigRabbit = config
	manager.conn, err = NewConnection()
	manager.ContextRabbit = ctx
	if err != nil {
		return nil, errors.Wrap(err, "create connection")
	}
	manager.conn.SetConfig(login, password, address, port, maxAttempts)
	return manager, nil
}

func (m *Manager) GetConn() *Connection {
	return m.conn
}

func (m *Manager) GetConfig() *Config {
	return m.ConfigRabbit
}

func (m *Manager) AddConsumer(ctx context.Context,name,consumerName,exchangeName,routingKey string) (*Consumer, error) {
	if _, ok := m.consumers[name]; ok {
		return m.consumers[name],nil
	}
	consumer, err := NewConsumer(ctx, name,consumerName,exchangeName,routingKey, m.conn)
	if err != nil {
		log.Printf("%v",err)
		return nil, errors.Wrap(err, "new consumer")
	}
	m.consumers[name] = consumer
	return consumer, nil
}

func (m *Manager) AddPublisher(ctx context.Context,name,exchangeName,routingKey string) (*Publisher, error) {
	if _, ok := m.publishers[name]; ok {
		return m.publishers[name],nil
	}
	publisher, err := NewPublisher(ctx, name,exchangeName,routingKey, m.conn)
	if err != nil {
		return nil, errors.Wrap(err, "new publicher")
	}
	m.publishers[name] = publisher
	return publisher, nil
}

// Ping checks that rabbitMQ connections is live
func (m *Manager) Ping(ctx context.Context) error {
	client, err := amqp.Dial(m.conn.Dsn)
	if err != nil {
		return errors.Wrap(err, "ampq dial")
	}

	if err := client.Close(); err != nil {
		return errors.Wrap(err, "close client")
	}
	return nil
}

func (m *Manager) Reconnect() {
	for {
		select {
			case result,ok := <-m.conn.ReconnectRabbit: {
				if !ok {
					return 
				}
				if result == 1 {
					m.Mu.Lock()
					for _,item := range m.consumers {
						item.Reconnect<-1
					}
					m.Mu.Unlock()
				}
			}
		}
	}
}

func (m *Manager) Start(ctx context.Context, errorGroup *errgroup.Group) error {
	if m.conn.OriConn() != nil {
		return nil
	}
	log.Println("establishing connection...")
	if err := m.conn.Connect(ctx, errorGroup); err != nil {
		return errors.Wrap(err, "connect")
	}
	go m.Reconnect()
	return nil
}

// Shutdown shutdowns queue worker. Later will also
func (m *Manager) Shutdown(ctx context.Context) error {
	log.Println("shutting down")
	if err := m.shutdown(ctx); err != nil {
		return errors.Wrapf(err, "shutdown %q", ProviderName)
	}
	log.Println("shutted down")
	return nil
}

func (m *Manager) shutdown(ctx context.Context) error {
	if m.conn == nil {
		return nil
	}
	log.Println("closing connection...")
	if err := m.conn.Close(ctx); err != nil {
		return errors.Wrap(err, "close connection")
	}
	m.consumers = nil  		 
	m.publishers = nil
	m.conn = nil
	return nil
}

func Init(ctx context.Context,errGroup *errgroup.Group, conf *config.Config) (*Manager,error) {
    manager,err := NewManager(ctx, conf.RabbitMq.Login, conf.RabbitMq.Password, conf.RabbitMq.Address, conf.RabbitMq.Port,conf.RabbitMq.MaxAttempts)
    if err != nil {
       return manager,err
    }
    err = manager.Start(ctx,errGroup)
    if err != nil {
       return manager,err
    }
    err = manager.Ping(ctx)
    if err != nil {
       return manager,err
    }
    return manager,nil
}