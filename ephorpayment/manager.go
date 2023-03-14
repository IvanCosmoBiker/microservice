package ephorpayment

import (
	config "ephorservices/config"
	client "ephorservices/ephorpayment/client"
	server "ephorservices/ephorpayment/server"
)

type Manager struct {
	Config     *config.Config
	ServerGRPC *server.PaymentServer
	ClientGRPC *client.PaymentClient
}

func Init(cfg *config.Config) *Manager {
	manager := &Manager{
		Config: cfg,
	}
	manager.InitServer()
	manager.InitClient()
	return manager
}

func (m *Manager) InitServer() {
	serverPayment := server.Init(m.Config)
	m.ServerGRPC = serverPayment
	go m.ServerGRPC.Serve()
}

func (m *Manager) InitClient() {
	clientPayment := client.Init(m.Config)
	m.ClientGRPC = clientPayment
}

func (m *Manager) GetClient() *client.PaymentClient {
	return m.ClientGRPC
}
