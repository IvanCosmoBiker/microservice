package coninit

import (
	"context"
	"database/sql"
	logger "ephorservices/pkg/logger"
	"fmt"
	"time"

	pgxpool "github.com/jackc/pgx/v5/pgxpool"
)

type Connection struct {
	CloseConn         bool
	HealthCheckPeriod time.Duration
	configPool        *pgxpool.Config
	Connection        []*sql.Conn
	conn              *pgxpool.Pool
	debug             bool
	ctxApp            context.Context
}

func NewConnectionPool(ctx context.Context) *Connection {
	return &Connection{
		debug:  false,
		ctxApp: ctx,
	}
}

func (Con *Connection) Init(login, password, address, database string, port, maxConnection, minConnection uint16, healthCheckPeriod int, protocolBinary, debug bool) error {
	errConfig := Con.SetConfig(login, password, address, database, port, maxConnection, minConnection, healthCheckPeriod, protocolBinary, debug)
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

func (Con *Connection) SetConfig(login, password, address, database string, port, maxConnection, minConnection uint16, healthCheckPeriod int, protocolBinary, debug bool) error {
	Con.debug = debug
	stringConn := Con.GetStringConnect(login, password, address, database, port, maxConnection)
	ConfigPg, err := pgxpool.ParseConfig(stringConn)
	if err != nil {
		return err
	}
	ConfigPg.MaxConns = int32(maxConnection)
	ConfigPg.MinConns = int32(minConnection)
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
	Con.conn, err = pgxpool.NewWithConfig(Con.ctxApp, Con.configPool)
	if err != nil {
		logger.Log.Error(err.Error())
		return err
	}
	logger.Log.Infof("%+v", Con.conn)
	Con.conn.Exec(Con.ctxApp, "SET timezone TO UTC")
	return nil
}

func (Con *Connection) Acquire(ctx context.Context) (*pgxpool.Conn, error) {
	Con.Ping()
	connection, err := Con.conn.Acquire(ctx)
	if err != nil {
		return connection, err
	}
	return connection, nil
}

func (Con *Connection) Ping() error {
	if Con.CloseConn == true {
		return nil
	}
	err := Con.conn.Ping(Con.ctxApp)
	if err != nil {
		for {
			select {
			case <-Con.ctxApp.Done():
				return nil
			case <-time.After(Con.HealthCheckPeriod * time.Second):
				logger.Log.Errorf("Reconnect to db after %v seconds", Con.HealthCheckPeriod)
				err = Con.reconnect()
				if err == nil {
					return err
				}
			}
		}
	}
	return err
}

func (Con *Connection) reconnect() error {
	var err error
	Con.conn, err = pgxpool.NewWithConfig(Con.ctxApp, Con.configPool)
	if err != nil {
		logger.Log.Error(err.Error())
		return err
	}
	if err = Con.conn.Ping(context.Background()); err != nil {
		logger.Log.Error(err.Error())
		return err
	}
	return nil
}

func (Con *Connection) Close() {
	Con.CloseConn = true
	Con.conn.Close()
}
