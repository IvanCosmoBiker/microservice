package mqttmanager

import (
	"context"
	"encoding/json"
	connection "ephorservices/pkg/mqttmanager/connection"
	consumer "ephorservices/pkg/mqttmanager/consumer"
	publisher "ephorservices/pkg/mqttmanager/publisher"
	"errors"
	"fmt"
	"log"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

var (
	Time_Is_End = errors.New("Time is end")
)

var (
	TopicTemplateDevice = "ephor/1/dev/"
)

type BrokerManager struct {
	State              uint8
	ContextMqtt        context.Context
	Connection         *connection.BrokerMqttConnection
	Consumers          map[string]*consumer.Consumer
	Publisher          *publisher.Publisher
	OpenConnectionChan chan bool
}

var Broker *BrokerManager

func New(contex context.Context) (*BrokerManager, error) {
	broker := &BrokerManager{
		ContextMqtt:        contex,
		Consumers:          make(map[string]*consumer.Consumer),
		OpenConnectionChan: make(chan bool),
	}
	Broker = broker
	err := broker.NewConnect()
	if err != nil {
		return broker, err
	}

	return broker, err
}

func (bm *BrokerManager) SetConfig(address, port, login, password, clientId string,
	protocolVersion, disconnect uint,
	backOffPolicySendMassage, backOffPolicyConnection []time.Duration,
	executeTimeSeconds int) {
	bm.Connection.Config.Address = address
	bm.Connection.Config.Port = port
	bm.Connection.Config.Login = login
	bm.Connection.Config.ClientID = clientId
	bm.Connection.Config.Password = password
	bm.Connection.Config.ProtocolVersion = protocolVersion
	bm.Connection.Config.Disconnect = disconnect
	bm.Connection.Config.ExecuteTimeSeconds = time.Duration(executeTimeSeconds)
	bm.Connection.Config.BackOffPolicySendMassage = backOffPolicySendMassage
	bm.Connection.Config.BackOffPolicyConnection = backOffPolicyConnection
}

func (bm *BrokerManager) InitPublisher() {
	bm.Publisher, _ = publisher.NewPublisher(bm.ContextMqtt, "Publisher", bm.Connection)
}

func (bm *BrokerManager) NewConnect() error {
	con, err := connection.New()
	if err != nil {
		return err
	}
	bm.Connection = con
	bm.Connection.SetOnConnection(bm.OnConnect)
	bm.Connection.SetReconnect(bm.OnReconnect)
	bm.Connection.SetConnLost(bm.OnConnectionLost)
	return err
}

func (bm *BrokerManager) Start() error {
	err := bm.Connection.StartConnection()
	if err != nil {
		return err
	}
	bm.InitPublisher()
	return err
}

func (bm *BrokerManager) WaitConnection() {
	// openConnectionMqtt, _ := <-bm.OpenConnectionChan
	// if openConnectionMqtt {
	// 	return
	// }
}

func (bm *BrokerManager) Shutdown(ctx context.Context) (err error) {
	bm.Connection.CloseConnection(ctx)
	return err
}

func (bm *BrokerManager) AddConsumer(ctx context.Context, name, topic string) (*consumer.Consumer, error) {
	if _, ok := bm.Consumers[name]; ok {
		return bm.Consumers[name], nil
	}
	Consumer, err := consumer.NewConsumer(ctx, name, topic, bm.Connection)
	if err != nil {
		log.Printf("%v", err)
		return Consumer, err
	}
	bm.Consumers[name] = Consumer
	return Consumer, nil
}

func (bm *BrokerManager) OnConnect(client mqtt.Client) {
	log.Println("Connection open")
}

func (bm *BrokerManager) OnReconnect(client mqtt.Client, confMqtt *mqtt.ClientOptions) {
	log.Println("Connection reconnect")
}

func (bm *BrokerManager) OnConnectionLost(client mqtt.Client, err error) {
	log.Println("Connection Close")
}

func (bm *BrokerManager) SendMessage(message map[string]interface{}, topic string) (err error) {
	err = bm.Publisher.SendMessage(bm.ContextMqtt, fmt.Sprintf("%s%v", TopicTemplateDevice, topic), message)
	if err != nil {
		for _, v := range bm.Connection.Config.BackOffPolicySendMassage {
			time.Sleep(v * time.Second)
			if err = bm.Publisher.SendMessage(bm.ContextMqtt, fmt.Sprintf("%s%v", TopicTemplateDevice, topic), message); err != nil {
				log.Printf("Fail send massage to device with imei: %s, wait time to retry send is %v", topic, v)
				time.Sleep(v * time.Second)
				continue
			}
			break
		}
	}
	return
}

func (bm *BrokerManager) WaitMessage(ChannelMessage chan []byte, waitSeconds int) (result map[string]interface{}, err error) {
	fmt.Println("[X]Timer")
	result = make(map[string]interface{})
	timer := time.NewTimer(time.Duration(waitSeconds) * time.Second)
	select {
	case <-timer.C:
		{
			timer.Stop()
			err = Time_Is_End
			return
		}
	case message, ok := <-ChannelMessage:
		{
			if !ok {
				err = errors.New("Close channel Transaction")
				timer.Stop()
				return
			}
			json.Unmarshal(message, &result)
			timer.Stop()
			fmt.Printf("%+v", result)
			err = nil
			return
		}
	}
}
