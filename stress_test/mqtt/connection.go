package main

import (
	mqtt "ephorservices/pkg/mqttmanager/connection"
	"fmt"
	"sync"
	"time"
)

var connectionMax = 200000000

func Init(clientId string) *mqtt.BrokerMqttConnection {
	// address, port, login, password, clientId string,
	// protocolVersion, disconnect uint,
	// backOffPolicySendMassage, backOffPolicyConnection []time.Duration,
	// executeTimeSeconds int
	connectionMqtt, _ := mqtt.New()
	BackOffPolicySendMassage := make([]time.Duration, 3)
	BackOffPolicyConnection := make([]time.Duration, 3)
	connectionMqtt.SetConfig("188.225.18.140", "1883", "device", "ephor2021", clientId, uint(0), uint(255), BackOffPolicySendMassage, BackOffPolicyConnection, 60)
	return connectionMqtt
}

func main() {
	Channel := make(chan bool)
	var wg sync.WaitGroup
	for i := 0; i < connectionMax; i++ {
		mqttConn := Init(fmt.Sprintf("Ephor%v", i))
		fmt.Printf("Ephor%v\n", i)
		wg.Add(1)
		go func(channel chan bool) {
			err := mqttConn.StartConnection()
			if err != nil {
				fmt.Printf("%s\n", err.Error())
			}
			mqttConn.Client.Connect()
			_, ok := <-channel
			fmt.Printf("%v\n", ok)
			if ok == false {
				wg.Done()
			}

		}(Channel)
	}
	time.Sleep(time.Duration(10) * time.Second)
	close(Channel)
	wg.Wait()
}
