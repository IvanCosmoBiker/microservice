package deviceinterface 

import (
    transactionStruct "data/transaction"
    ConnectionRabbitMQ "lib-rabbitmq"
    connectionPostgresql "connectionDB/connect"
    interfaseBank "interface/bankinterface"
	requestApi "data/requestApi"
    transactionDispetcher "transactionDispetcher"
    configEphor "configEphor"
)

var (
    TypeCoffee = 0
	TypeSnack = 1
	TypeHoreca = 2
	TypeSodaWater = 3
	TypeMechanical = 4
	TypeComb = 5
	TypeMicromarket = 6
	TypeCooler = 7
)

type Device interface {
    InitDeviceData(transactionStruct.Transaction)
    SendMessage(trasactionStruct transactionStruct.Transaction, 
	conn connectionPostgresql.DatabaseInstance, 
	rebbit *ConnectionRabbitMQ.ChannelMQ, 
	conf *configEphor.Config, 
	req requestApi.Request, 
	bank interfaseBank.Bank,
	resultBank map[string]interface{},
	transactionDispetcer transactionDispetcher.TransactionDispetcher,
	channel chan []byte) map[string]interface{}
}