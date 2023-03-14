package connection

import (
	"context"
	"fmt"
	"log"
	"net/url"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

const (
	State_Idle uint8 = iota
	State_Work
	State_Warning
	State_Error
)

type Config struct {
	Address                  string
	Port                     string
	Login                    string
	ClientID                 string
	Password                 string
	ProtocolVersion          uint
	Disconnect               uint
	ExecuteTimeSeconds       time.Duration
	BackOffPolicySendMassage []time.Duration
	BackOffPolicyConnection  []time.Duration
}

type BrokerMqttConnection struct {
	Config           *Config
	State            uint8
	ConfigConnection *mqtt.ClientOptions
	Client           mqtt.Client
}

func New() (con *BrokerMqttConnection, err error) {
	con = &BrokerMqttConnection{
		Config:           &Config{},
		ConfigConnection: mqtt.NewClientOptions(),
		State:            State_Idle,
	}
	return con, err
}

func (bmc *BrokerMqttConnection) StartConnection() error {
	err := bmc.SetConfigConnection()
	if err != nil {
		return err
	}
	err = bmc.Connect()
	return err
}

func (bmc *BrokerMqttConnection) SetConfig(address, port, login, password, clientId string,
	protocolVersion, disconnect uint,
	backOffPolicySendMassage, backOffPolicyConnection []time.Duration,
	executeTimeSeconds int) {
	bmc.Config.Address = address
	bmc.Config.Port = port
	bmc.Config.Login = login
	bmc.Config.ClientID = clientId
	bmc.Config.Password = password
	bmc.Config.ProtocolVersion = protocolVersion
	bmc.Config.Disconnect = disconnect
	bmc.Config.ExecuteTimeSeconds = time.Duration(executeTimeSeconds)
	bmc.Config.BackOffPolicySendMassage = backOffPolicySendMassage
	bmc.Config.BackOffPolicyConnection = backOffPolicyConnection
}

func (bmc *BrokerMqttConnection) SetConfigConnection() (err error) {
	defer checkError()
	bmc.checkDefaultSettingAndSetDefault()
	var urlMqtt *url.URL
	urlMqtt, err = bmc.MakeUrl()
	if err != nil {
		return err
	}
	sliceUrls := bmc.MakeSliceUrl(urlMqtt)
	bmc.ConfigConnection.Servers = sliceUrls
	bmc.SetValueConfig()
	return err
}

func (bmc *BrokerMqttConnection) MakeUrl() (urlMqtt *url.URL, err error) {
	stringUrl := fmt.Sprintf("mqtt://%s:%s", bmc.Config.Address, bmc.Config.Port)
	urlMqtt, err = url.Parse(stringUrl)
	return urlMqtt, err
}

func (bmc *BrokerMqttConnection) MakeSliceUrl(urls ...*url.URL) (sliceUrls []*url.URL) {
	sliceUrls = append(sliceUrls, urls...)
	return sliceUrls
}

func (bmc *BrokerMqttConnection) checkDefaultSettingAndSetDefault() (err error) {
	if len(bmc.Config.Address) < 1 {
		bmc.Config.Address = "127.0.0.1"
	}
	if len(bmc.Config.Port) < 1 {
		bmc.Config.Port = "1883"
	}
	if len(bmc.Config.BackOffPolicySendMassage) < 1 {
		bmc.Config.BackOffPolicySendMassage = []time.Duration{1, 2, 4, 8, 16, 32, 50, 60, 70, 80, 90, 100, 110, 128, 256}
	}
	if len(bmc.Config.BackOffPolicyConnection) < 1 {
		bmc.Config.BackOffPolicySendMassage = []time.Duration{1, 2, 4, 8, 16, 32, 50, 60, 70, 80, 90, 100, 110, 128, 256}
	}
	if bmc.Config.ExecuteTimeSeconds == time.Duration(0) {
		bmc.Config.ExecuteTimeSeconds = time.Duration(60)
	}
	if bmc.Config.ProtocolVersion == uint(0) {
		bmc.Config.ProtocolVersion = uint(3)
	}
	if len(bmc.Config.ClientID) < 1 {
		bmc.Config.ClientID = "ephorservices"
	}
	if bmc.Config.Disconnect == uint(0) {
		bmc.Config.Disconnect = uint(60000)
	}
	return err
}

func (bmc *BrokerMqttConnection) SetValueConfig() {
	bmc.ConfigConnection.Username = bmc.Config.Login
	bmc.ConfigConnection.Password = bmc.Config.Password
	bmc.ConfigConnection.ProtocolVersion = bmc.Config.ProtocolVersion
	bmc.ConfigConnection.ClientID = bmc.Config.ClientID
	bmc.ConfigConnection.MaxReconnectInterval = 1
	bmc.ConfigConnection.AutoReconnect = true
}

func (bmc *BrokerMqttConnection) Connect() (err error) {
	defer checkError()
	client := mqtt.NewClient(bmc.ConfigConnection)
	bmc.Client = client
	err = bmc.OpenConnection()
	return err
}

func (bmc *BrokerMqttConnection) SetOnConnection(f mqtt.OnConnectHandler) {
	bmc.ConfigConnection.OnConnect = f
}

func (bmc *BrokerMqttConnection) SetConnLost(f mqtt.ConnectionLostHandler) {
	bmc.ConfigConnection.OnConnectionLost = f
}

func (bmc *BrokerMqttConnection) SetReconnect(f mqtt.ReconnectHandler) {
	bmc.ConfigConnection.OnReconnecting = f
}

func (bmc *BrokerMqttConnection) OpenConnection() (err error) {
	if token := bmc.Client.Connect(); token.Wait() && token.Error() != nil {
		err = token.Error()
		bmc.State = State_Error
		fmt.Println(token.Error())
		return err
	}
	fmt.Println("Connection to server")
	bmc.State = State_Work
	return err
}

func (bmc *BrokerMqttConnection) CloseConnection(ctx context.Context) {
	bmc.Client.Disconnect(bmc.Config.Disconnect)
	log.Println("Disconnect")
}

func (bmc *BrokerMqttConnection) Publish(topic string, qos byte, retained bool, payload interface{}) mqtt.Token {
	return bmc.Client.Publish(topic, qos, retained, payload)
}

func (bmc *BrokerMqttConnection) Subscribe(topic string, qos byte, callback mqtt.MessageHandler) mqtt.Token {
	return bmc.Client.Subscribe(topic, qos, callback)
}

func checkError() (err error) {
	if r := recover(); r != nil {
		err = r.(error)
		return err
	}
	return err
}
