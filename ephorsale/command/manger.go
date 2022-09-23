package command

import (
	config "ephorservices/config"
	rabbit "ephorserices/pkg/rabbitmq"
)

type CommandManager struct {
	cfg *config.Config
	queue *rabbit.Manager
}

func Init(conf *config.Config, queueManager *rabbit.Manager) *CommandManager {
	return &CommandManager{
		cfg: conf,
		queue: queueManager,
	}
}