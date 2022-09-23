package coninit

import (
	"context"
	"fmt"
	"log"

	pgxpool "github.com/jackc/pgx/v4/pgxpool"
)

type Connection struct {
	configPool *pgxpool.Config
	conn       *pgxpool.Pool
	debug      bool
	ctxApp     context.Context
}

func NewConnectionPool(ctx context.Context) *Connection {
	return &Connection{
		debug:  false,
		ctxApp: ctx,
	}
}

func (Con *Connection) Init(login, password, address, database string, port, maxConnection, minConnection uint16, protocolBinary, debug bool) error {
	errConfig := Con.SetConfig(login, password, address, database, port, maxConnection, minConnection, protocolBinary, debug)
	if errConfig != nil {
		return errConfig
	}
	errCon := Con.StartConnectionPool()
	if errCon != nil {
		return errCon
	}
	return nil
}

func (Con *Connection) GetStringConnect(login, password, address, database string, port, maxConnection uint16) string {
	stringConnection := fmt.Sprintf("user=%s password=%s host=%v port=%v dbname=%s sslmode=disable pool_max_conns=%v", login, password, address, port, database, maxConnection)
	return stringConnection
}

func (Con *Connection) SetConfig(login, password, address, database string, port, maxConnection, minConnection uint16, protocolBinary, debug bool) error {
	Con.debug = debug
	stringConn := Con.GetStringConnect(login, password, address, database, port, maxConnection)
	log.Printf("%s", stringConn)
	ConfigPg, err := pgxpool.ParseConfig(stringConn)
	log.Printf("%+v", ConfigPg)
	if err != nil {
		return err
	}
	ConfigPg.MinConns = int32(minConnection)
	ConfigPg.ConnConfig.PreferSimpleProtocol = protocolBinary
	log.Printf("%+v", ConfigPg.ConnConfig)
	Con.configPool = ConfigPg
	return nil
}

func (Con *Connection) GetConfig() *pgxpool.Config {
	return Con.configPool
}

func (Con *Connection) GetConn() *pgxpool.Pool {
	return Con.conn
}

func (Con *Connection) StartConnectionPool() error {
	var err error
	Con.conn, err = pgxpool.ConnectConfig(Con.ctxApp, Con.configPool)
	if err != nil {
		return err
	}
	Con.conn.Exec(Con.ctxApp, "SET timezone TO UTC")
	return nil
}

func (Con *Connection) Acquire() (*pgxpool.Conn, error) {
	connection, err := Con.conn.Acquire(Con.ctxApp)
	if err != nil {
		return connection, err
	}
	return connection, nil
}

func (Con *Connection) Ping() error {
	err := Con.conn.Ping(Con.ctxApp)
	return err
}

func (Con *Connection) Close() {
	Con.conn.Close()
}
