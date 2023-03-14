package db

import (
	"context"
	coninit "ephorservices/pkg/orm/db/coninit"
	core "ephorservices/pkg/orm/db/core"
	"log"
)

type Manager struct {
	Conn      *core.CoreConnection
	Debug     bool
	reconnect chan bool
	Context   context.Context
}

var ConnectionManager *Manager

func NewManager(ctx context.Context, debug bool) (*Manager, error) {
	manager := &Manager{
		Debug:     debug,
		Context:   ctx,
		reconnect: make(chan bool),
	}
	ConnectionManager = manager
	return manager, nil
}

func (m *Manager) Init(login, password, address, databaseName string, port, pgConnectionPool, pgConnectionMin, pgConnectionMax uint16, healthCheckPeriod int, preferSimpleProtocol bool, debug bool) error {
	Conn := coninit.NewConnectionPool(m.Context)
	errInitConn := Conn.Init(login, password, address, databaseName, port, pgConnectionPool, pgConnectionMin, healthCheckPeriod, preferSimpleProtocol, debug)
	if errInitConn != nil {
		log.Println(errInitConn)
		return errInitConn
	}
	m.Conn = core.NewCoreConnection(m.Context, Conn)
	return nil
}

func (m *Manager) Close() {
	m.Conn.Close()
}

func (m *Manager) Reconnect(ctx context.Context) {

}

func Init(login, password, address, databaseName string, port, pgConnectionPool, pgConnectionMin, pgConnectionMax uint16, healthCheckPeriod int, preferSimpleProtocol bool, debug bool, ctx context.Context) (*Manager, error) {
	var err error
	ConnectionDb, _ := NewManager(ctx, true)
	err = ConnectionDb.Init(login, password, address, databaseName, port, pgConnectionPool, pgConnectionMin, pgConnectionMax, healthCheckPeriod, preferSimpleProtocol, debug)
	if err != nil {
		return ConnectionDb, err
	}
	return ConnectionDb, err
}
