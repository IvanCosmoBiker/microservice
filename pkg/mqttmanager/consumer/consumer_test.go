package consumer

import (
	"context"
	config "ephorservices/config"
	connection "ephorservices/pkg/mqttmanager/connection"
	//"errors"
	"fmt"
	"golang.org/x/sync/errgroup"
	"log"
	"testing"
	"time"
)

var Config *config.Config
var ctx = context.Background()
var errGroup = new(errgroup.Group)

// var connectionOn mqtt.OnConnectHandler = func(client mqtt.Client) {
// 	fmt.Println("Connected")
// }

type CunsumerTestSubscriber struct {
}

func (cs *CunsumerTestSubscriber) Consume(ct context.Context, data []byte) error {
	fmt.Printf("%s", data)
	return nil
}
func (cs *CunsumerTestSubscriber) Shutdown(ct context.Context) error {
	fmt.Printf("%s", "Get Shutdown")
	return nil
}

func init() {
	newConfig := config.Config{}
	newConfig.Transport.Mqtt.Login = "device"
	newConfig.Transport.Mqtt.Address = "188.225.18.140"
	newConfig.Transport.Mqtt.Password = "ephor2021"
	newConfig.Transport.Mqtt.Port = "1883"
	newConfig.Transport.Mqtt.ExecuteTimeSeconds = 30
	newConfig.Transport.Mqtt.BackOffPolicySendMassage = make([]time.Duration, 3)
	newConfig.Transport.Mqtt.BackOffPolicyConnection = make([]time.Duration, 3)
	Config = &newConfig
}

func TestNewConsumer(t *testing.T) {
	manager, err := connection.New()
	if err != nil {
		log.Printf("%v", err)
	}
	manager.SetConfig(Config.Transport.Mqtt.Address,
		Config.Transport.Mqtt.Port,
		Config.Transport.Mqtt.Login,
		Config.Transport.Mqtt.Password,
		"test",
		uint(0),
		uint(200),
		Config.Transport.Mqtt.BackOffPolicySendMassage,
		Config.Transport.Mqtt.BackOffPolicyConnection,
		Config.Transport.Mqtt.ExecuteTimeSeconds)
	manager.StartConnection()
	ConsumerNew, err := NewConsumer(ctx, "Pay", "ephor/1/pay", manager)
	if err != nil {
		t.Error(err.Error())
	}
	sub := &CunsumerTestSubscriber{}
	//ctx context.Context, errorGroup *errgroup.Group, subscriber Subscriber
	ConsumerNew.Subscribe(ctx, errGroup, sub)
	//Publish(topic string, qos byte, retained bool, payload interface{})
	manager.Publish("ephor/1/pay", 2, false, "Kyky")
	time.Sleep(time.Duration(6) * time.Second)
}

func TestPublish(t *testing.T) {
	//ctx context.Context, name, consumerName, topic, routingKey string, ch ConnectorConsumer
	manager, err := connection.New()
	if err != nil {
		log.Printf("%v", err)
	}
	manager.SetConfig(Config.Transport.Mqtt.Address,
		Config.Transport.Mqtt.Port,
		Config.Transport.Mqtt.Login,
		Config.Transport.Mqtt.Password,
		"test",
		uint(0),
		uint(200),
		Config.Transport.Mqtt.BackOffPolicySendMassage,
		Config.Transport.Mqtt.BackOffPolicyConnection,
		Config.Transport.Mqtt.ExecuteTimeSeconds)
	manager.StartConnection()
	manager.Publish("ephor.pay", 2, false, "Kyky")
}
