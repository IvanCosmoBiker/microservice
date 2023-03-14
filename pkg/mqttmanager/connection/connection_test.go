package connection

import (
	"context"
	"errors"
	"fmt"
	mqtt "github.com/eclipse/paho.mqtt.golang"
	"testing"
	"time"
)

var TestConfig *Config
var ctx = context.Background()

func init() {
	newConfig := Config{}
	newConfig.Login = "device"
	newConfig.Address = "188.225.18.140"
	newConfig.Password = "ephor2021"
	newConfig.Port = "1883"
	newConfig.ExecuteTimeSeconds = 30
	newConfig.BackOffPolicySendMassage = make([]time.Duration, 3)
	newConfig.BackOffPolicyConnection = make([]time.Duration, 3)
	TestConfig = &newConfig
}

var connectionOn mqtt.OnConnectHandler = func(client mqtt.Client) {
	fmt.Println("Connected")
}

func TestInit(t *testing.T) {
	connection, err := New()
	if err != nil {
		t.Error(err.Error())
	}
	connection.Config = TestConfig
	err = connection.StartConnection()
	if err != nil {
		t.Error(err.Error())
	}
	if connection.ConfigConnection.Username == "" || connection.ConfigConnection.Password == "" || len(connection.ConfigConnection.Servers) < 1 {
		t.Error(errors.New("Empty config parametrs"))
	}
}

func TestConnection(t *testing.T) {
	connection, err := New()
	if err != nil {
		t.Error(err.Error())
	}
	connection.Config = TestConfig
	err = connection.StartConnection()
	if err != nil {
		t.Error(err.Error())
	}
	if connection.ConfigConnection.Username == "" || connection.ConfigConnection.Password == "" || len(connection.ConfigConnection.Servers) < 1 {
		t.Error(errors.New("Empty config parametrs"))
	}
	if connection.Client == nil {
		t.Error("Client Mqtt is nil")
	}
}

func TestStartConnection(t *testing.T) {
	connection, err := New()
	if err != nil {
		t.Error(err.Error())
	}
	connection.Config = TestConfig
	err = connection.StartConnection()
	if err != nil {
		t.Error(err.Error())
	}
	if connection.Client == nil {
		t.Error("Client Mqtt is nil")
	}
	connection.SetOnConnection(connectionOn)
	connection.OpenConnection()
	connection.CloseConnection(ctx)
}
