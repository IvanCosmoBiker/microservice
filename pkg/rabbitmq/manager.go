package rabbitmq

import (
	"amqp"
	"context"
	config "ephorservices/config"
	"log"
	"sync"

	"github.com/pkg/errors"
	"golang.org/x/sync/errgroup"
)

const ProviderName = "rabbitmq"

// Manager is a connection controlling structure. It controls
// connection, asynchronous queue and everything that related to
// specified connection.
type Manager struct {
	cfg           *config.Config
	conn          *Connection
	ContextRabbit context.Context
	consumers     map[string]*Consumer
	Publisher     *Publisher
	Mu            sync.Mutex
}

// NewManager create new Manager and init Connection to Rabbit.
func NewManager(ctx context.Context, cfg *config.Config) (*Manager, error) {
	var manager = &Manager{
		cfg:       cfg,
		consumers: make(map[string]*Consumer),
	}

	var err error
	manager.conn, err = NewConnection(manager.cfg)
	manager.ContextRabbit = ctx
	if err != nil {
		return nil, errors.Wrap(err, "create connection")
	}
	manager.conn.SetConfig()
	return manager, nil
}

func (m *Manager) GetConn() *Connection {
	return m.conn
}

func (m *Manager) AddConsumer(ctx context.Context, name, consumerName, exchangeName, routingKey string) (*Consumer, error) {
	if _, ok := m.consumers[name]; ok {
		return m.consumers[name], nil
	}
	consumer, err := NewConsumer(ctx, name, consumerName, exchangeName, routingKey, m.conn)
	if err != nil {
		log.Printf("%v", err)
		return nil, errors.Wrap(err, "new consumer")
	}
	m.consumers[name] = consumer
	return consumer, nil
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
		case result, ok := <-m.conn.ReconnectRabbit:
			{
				if !ok {
					return
				}
				if result == 1 {
					m.Mu.Lock()
					for _, item := range m.consumers {
						item.Reconnect <- 1
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
	mainPublisher, err := NewPublisher(ctx, "publisher", "", m.conn)
	mainPublisher.isConnected = true
	if err != nil {
		return err
	}
	m.Publisher = mainPublisher
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
	m.conn = nil
	return nil
}

func Init(ctx context.Context, errGroup *errgroup.Group, conf *config.Config) (*Manager, error) {
	manager, err := NewManager(ctx, conf)
	if err != nil {
		return manager, err
	}
	err = manager.Start(ctx, errGroup)
	if err != nil {
		return manager, err
	}
	err = manager.Ping(ctx)
	if err != nil {
		return manager, err
	}
	return manager, nil
}
