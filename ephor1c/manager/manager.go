package manager

import (
	"ephorservices/ephor1c/config"
	connectionDb "ephorservices/pkg/db"
	logger "ephorservices/pkg/logger"
)

type Manager struct {
	ConnectionDb *connectionDb.Manager
	Logger       *logger.Logger
	Config       *config.Config
}

func Init() *Manager {
	manager := &Manager{}
	manager.ReadConfig()
	manager.initLogger()
	return manager
}

func (m *Manager) initLogger() error {
	var err error
	m.Logger, err = logger.Init(m.Config.LogFile, m.Config.Name, m.Config.LoggingType, m.Config.LogFileEnable)
	return err
}

func (m *Manager) ReadConfig() {
	m.Config = &config.Config{}
	m.Config.Load()
}
