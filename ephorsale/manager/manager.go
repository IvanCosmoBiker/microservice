package ephorsale

import (
	"context"
	config "ephorservices/config"
	commandManager "ephorservices/ephorsale/command"
	"ephorservices/ephorsale/control"
	"ephorservices/ephorsale/sale"
	transactionDispetcher "ephorservices/ephorsale/transaction"
	Transport_init "ephorservices/ephorsale/transport"
	logger "ephorservices/pkg/logger"
	connectionDb "ephorservices/pkg/orm/db"
	"fmt"
	"os"
	"runtime"
	"runtime/debug"
	"time"

	"golang.org/x/sync/errgroup"
)

type Manager struct {
	Config                *config.Config
	ConnectionDb          *connectionDb.Manager
	Transport             *Transport_init.Transport
	TransactionDispetcher *transactionDispetcher.TransactionDispetcher
	Sale                  *sale.SaleManager
	Control               *control.Control
	Command               *commandManager.CommandManager
	ContextApp            context.Context
	ErrorGroup            *errgroup.Group
}

func New(conf *config.Config) *Manager {
	return &Manager{
		Config:     conf,
		ContextApp: context.Background(),
		ErrorGroup: new(errgroup.Group),
	}
}

func (m *Manager) InitService() {
	go m.StartGc()
	var err error
	fmt.Printf("%+v\n", m.Config)
	m.ConnectionDb, err = connectionDb.Init(m.Config.Db.Login,
		m.Config.Db.Password,
		m.Config.Db.Address,
		m.Config.Db.DatabaseName,
		m.Config.Db.Port,
		m.Config.Db.PgConnectionPool,
		m.Config.Db.PgConnectionMin,
		m.Config.Db.PgConnectionMax,
		m.Config.Db.ReconnectSecond,
		m.Config.Db.PreferSimpleProtocol,
		m.Config.Debug,
		m.ContextApp)
	if err != nil {
		logger.Log.Error(err.Error())
		os.Exit(1)
	}
	m.TransactionDispetcher = transactionDispetcher.New(m.Config, m.ContextApp)
	err = m.InitTransport()
	if err != nil {
		logger.Log.Error(err.Error())
		os.Exit(1)
	}
	Sale, err := sale.New(m.ContextApp, m.Config)
	if err != nil {
		logger.Log.Error(err.Error())
		os.Exit(1)
	}
	m.Sale = Sale
	m.Sale.InitApi(m.Config, m.ContextApp, m.ErrorGroup, m.Transport)
	m.Command = commandManager.New()
	m.Command.InitApi()
	m.Control = control.New()
	m.Control.InitApi()
	m.Transport.Listen()
}

func (m *Manager) InitTransport() error {
	m.Transport = Transport_init.New(m.ContextApp)
	m.Transport.InitHttp(m.Config.Services.Http.Address, m.Config.Services.Http.Port)
	err := m.Transport.InitMqtt(m.Config.Transport.Mqtt.Address,
		m.Config.Transport.Mqtt.Port,
		m.Config.Transport.Mqtt.Login,
		m.Config.Transport.Mqtt.Password,
		m.Config.Transport.Mqtt.ClientID,
		m.Config.Transport.Mqtt.ProtocolVersion,
		m.Config.Transport.Mqtt.Disconnect,
		m.Config.Transport.Mqtt.BackOffPolicySendMassage,
		m.Config.Transport.Mqtt.BackOffPolicyConnection,
		m.Config.Transport.Mqtt.ExecuteTimeSeconds)
	return err
}

func (m *Manager) StopService() error {
	m.Transport.CloseHttp()
	for {
		transactions := m.TransactionDispetcher.GetTransactions()
		if len(transactions) < 1 {
			break
		}
	}
	m.Transport.CloseMqtt()
	m.ConnectionDb.Close()
	return nil
}

func (m *Manager) StartGc() {
	for {
		time.Sleep(time.Duration(30) * time.Minute)
		var ms1, ms2 runtime.MemStats
		fmt.Println("Start GC AND free memory")
		runtime.GC()
		runtime.ReadMemStats(&ms1)
		debug.FreeOSMemory()
		runtime.ReadMemStats(&ms2)
		fmt.Println("Idle memory delta: ", (int64(ms2.HeapIdle)-int64(ms1.HeapIdle))/int64(1024))
	}
}
