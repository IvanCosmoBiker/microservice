package service_http

import "encoding/json"

type RequestSale struct {
	Config struct {
		AccountId    int
		AutomatId    int
		DeviceType   int
		QrFormat     int
		SbpPoint     string
		Login        string
		Password     string
		UserPhone    string
		ReturnUrl    string
		DeepLink     string
		TokenType    int
		BankType     int
		PayType      int
		CurrensyCode int
		Language     int
		TaxSystem    int
		Description  string
	}
	Service_id    string
	KeyPayment    string
	SecretKey     string
	IdTransaction string
	Date          string
	MerchantId    string
	GateWay       string
	PaymentToken  string
	Sum           int
	Imei          string
	Products      []map[string]interface{}
}

type RequestFiscalServer struct {
	Events             []string
	Imei, CheckId, Inn string
	ConfigFR           struct {
		Host, Cert, Key, Sign, Port, Login, Password string
		Fiscalization                                int
	}
	InQueue int
	TypeFr  int
	Fields  struct {
		Request json.RawMessage
	}
}
