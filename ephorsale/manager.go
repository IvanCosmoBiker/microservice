package ephorsale

import (
	"context"
	config "ephorservices/config"
	"ephorservices/ephorsale/fiscal"
	"ephorservices/ephorsale/payment"
	"ephorservices/ephorsale/sale"
	transactionDispetcher "ephorservices/ephorsale/transaction"
	connectionDb "ephorservices/pkg/db"
	RabbitMQ "ephorservices/pkg/rabbitmq"
	transportHttp "ephorservices/pkg/transportprotocol/http/v1"
	"log"

	"golang.org/x/sync/errgroup"
)

type SaleServiceManager struct {
	Config                *config.Config
	Service               *SaleService
	ConnectionDb          *connectionDb.Manager
	QueueManager          *RabbitMQ.Manager
	RequestManager        *transportHttp.ServerHttp
	Fiscal                *fiscal.FiscalManager
	Payment               *payment.PaymentManager
	Sale                  *sale.SaleManager
	TransactionDispetcher *transactionDispetcher.TransactionDispetcher
	connectRabbit         chan bool
	ContextApp            context.Context
	ErrorGroup            *errgroup.Group
}

func (m *SaleServiceManager) InitService(s *SaleService) {
	var err error
	m.Service = s
	m.Config = m.Service.ConfigFile
	m.ContextApp = context.Background()
	m.ErrorGroup = new(errgroup.Group)
	m.ConnectionDb, err = connectionDb.Init(m.Service.ConfigFile, m.ContextApp)
	if err != nil {
		log.Fatal(err)
	}
	m.TransactionDispetcher = transactionDispetcher.New(m.ConnectionDb)
	m.QueueManager, err = RabbitMQ.Init(m.ContextApp, m.ErrorGroup, m.Service.ConfigFile)
	if err != nil {
		log.Fatal(err)
	}
	m.RequestManager = transportHttp.Init(m.Service.ConfigFile)
	m.Fiscal, err = fiscal.Init(m.ConnectionDb, m.Config, m.ContextApp)
	if err != nil {
		log.Fatal(err)
	}
	m.Payment, err = payment.Init(m.Config, m.ContextApp)
	if err != nil {
		log.Fatal(err)
	}
	m.Sale, err = sale.Init(m.ContextApp, m.Config, m.QueueManager, m.ErrorGroup, m.TransactionDispetcher, m.Fiscal, m.Payment, m.RequestManager)
	if err != nil {
		log.Fatal(err)
	}
	go m.RequestManager.StartListener()
}

func (m *SaleServiceManager) StopService() error {
	m.RequestManager.CloseListener()
	log.Println("1")
	for {
		transactions := m.TransactionDispetcher.GetTransactions()
		if len(transactions) < 1 {
			break
		}
	}
	log.Println("2")
	m.QueueManager.Shutdown(m.ContextApp)
	log.Println("3")
	m.ConnectionDb.Close()
	log.Println("4")
	return nil
}
