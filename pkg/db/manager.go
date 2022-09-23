package db

import (
	"context"
	config "ephorservices/config"
	coninit "ephorservices/pkg/db/coninit"
	core "ephorservices/pkg/db/core"
)

type Manager struct {
	Conn      *core.CoreConnection
	Debug     bool
	reconnect chan bool
	Context   context.Context
}

func NewManager(ctx context.Context, debug bool) (*Manager, error) {
	return &Manager{
		Debug:     debug,
		Context:   ctx,
		reconnect: make(chan bool),
	}, nil
}

func (m *Manager) Init(conf *config.Config) error {
	Conn := coninit.NewConnectionPool(m.Context)
	errInitConn := Conn.Init(conf.Db.Login, conf.Db.Password, conf.Db.Address, conf.Db.DatabaseName, conf.Db.Port, conf.Db.PgConnectionPool, conf.Db.PgConnectionMin, conf.Db.PreferSimpleProtocol, conf.Debug)
	if errInitConn != nil {
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

func Init(conf *config.Config, ctx context.Context) (*Manager, error) {
	var err error
	ConnectionDb, _ := NewManager(ctx, true)
	err = ConnectionDb.Init(conf)
	if err != nil {
		return ConnectionDb, err
	}
	return ConnectionDb, err
}
