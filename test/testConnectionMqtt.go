package main

import (
	config "ephorservices/config"
	mqttConnection "ephorservices/pkg/mqttmanager/connection"
	"fmt"
	"log"
	"os"
	"os/signal"

	"golang.org/x/sync/errgroup"
)

func main() {
	config, _ := ReadConfig()
	ErrorGroup := new(errgroup.Group)
	InitConnections(config, ErrorGroup)
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	func() {
		<-c
		os.Exit(1)
	}()
}

func ReadConfig() (*config.Config, bool) {
	var config = config.Config{}
	config.Load()
	if config.LogFile != "" {
		file, err := os.OpenFile(config.LogFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
		if err != nil {
			log.Panic(err)
		}
		fmt.Println(config.Db.Login)
		log.SetOutput(file)
		return &config, true
	}
	return &config, true
}

func InitConnections(cfg *config.Config, errGroup *errgroup.Group) {
	for i := 0; i < 100; i++ {
		id := fmt.Sprintf("%s%v", cfg.Transport.Mqtt.ClientID, i)
		cfg.Transport.Mqtt.ClientID = id
		client, _ := mqttConnection.Init(cfg, errGroup)
		client.OpenConnection()
	}
}
