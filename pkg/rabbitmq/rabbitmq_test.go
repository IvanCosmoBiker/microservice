package rabbitmq

import (
	"context"
	config "ephorservices/config"
	"fmt"
	"golang.org/x/sync/errgroup"
	"testing"
	"time"
)

var Config *config.Config
var ctx = context.Background()
var errGroup = new(errgroup.Group)

func init() {
	newConfig := config.Config{}
	newConfig.RabbitMq.Login = "device"
	newConfig.RabbitMq.Address = "188.225.18.140"
	newConfig.RabbitMq.Password = "ephor2021"
	newConfig.RabbitMq.Port = "5672"
	newConfig.RabbitMq.MaxAttempts = 10
	newConfig.RabbitMq.ExecuteTimeSeconds = 30
	newConfig.RabbitMq.PoolChannel = 1
	newConfig.RabbitMq.BackOffPolicySendMassage = make([]time.Duration, 3)
	newConfig.RabbitMq.BackOffPolicyConnection = make([]time.Duration, 3)
	Config = &newConfig
}

func PublishTest(pub *Publisher, massage map[string]interface{}) {
	err := pub.SendMessage(ctx, "1234", massage)
	if err != nil {
		fmt.Println(err.Error())
	}
}

func TestGetConnection(t *testing.T) {
	_, err := NewManager(ctx, Config)
	if err != nil {
		t.Error(err.Error())
	}
}

func TestOnePublish(t *testing.T) {
	manager, err := NewManager(ctx, Config)
	if err != nil {
		t.Error(err.Error())
	}
	err = manager.Start(ctx, errGroup)
	if err != nil {
		t.Error(err.Error())
	}
	send := make(map[string]interface{})
	send["one"] = "test"
	sendTwo := make(map[string]interface{})
	sendTwo["two"] = "test"
	go PublishTest(manager.Publisher, send)
	go PublishTest(manager.Publisher, sendTwo)
	time.Sleep(10 * time.Second)
}
