package mqttmanager

import (
	"context"
	config "ephorservices/config"
	//"errors"
	//"fmt"
	"log"
	"testing"
	"time"

	//mqtt "github.com/eclipse/paho.mqtt.golang"
	"golang.org/x/sync/errgroup"
)

var Config *config.Config
var ctx = context.Background()
var errGroup = new(errgroup.Group)
var Manager *BrokerManager

type CunsumerTestSubscriber struct {
}

func (cs *CunsumerTestSubscriber) Consume(ct context.Context, data []byte) error {
	log.Printf("%s", "Message")
	return nil
}
func (cs *CunsumerTestSubscriber) Shutdown(ct context.Context) error {
	log.Printf("%s", "Get Shutdown")
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

func AddConsumerTest(manager *BrokerManager) {
	consumer, err := manager.AddConsumer(ctx, "Pay", "pay.ephor")
	if err != nil {
		log.Printf("%v", err)
	}
	pay := &CunsumerTestSubscriber{}
	err = consumer.Subscribe(ctx, errGroup, pay)
	if err != nil {
		log.Printf("%v", err)
	}
}

func TestNewConsumer(t *testing.T) {
	manager, err := New(ctx)
	if err != nil {
		log.Printf("%v", err)
	}
	// address, port, login, password, clientId string,
	// protocolVersion, disconnect uint,
	// backOffPolicySendMassage, backOffPolicyConnection []time.Duration,
	// executeTimeSeconds int
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
	manager.
		Manager = manager
	AddConsumerTest(manager)

}

func TestPublish(t *testing.T) {
	manager, err := New(ctx)
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
	manager.Publisher.SendMessage(ctx, "pay.ephor", "Kyky")
	manager.Publisher.SendMessage(ctx, "pay.ephor", "Kyky")
	manager.Publisher.SendMessage(ctx, "pay.ephor", "Kyky")
	manager.Publisher.SendMessage(ctx, "pay.ephor", "Kyky")
	manager.Publisher.SendMessage(ctx, "pay.ephor", "Kyky")
	manager.Publisher.SendMessage(ctx, "pay.ephor", "Kyky")
	manager.Publisher.SendMessage(ctx, "pay.ephor", "Kyky")
	manager.Publisher.SendMessage(ctx, "pay.ephor", "Kyky")
}
