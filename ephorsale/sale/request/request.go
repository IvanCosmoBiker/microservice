package request

import (
	"encoding/json"
)

type RequestPay struct {
	Config struct {
		UserPhone    string
		ReturnUrl    string
		DeepLink     string
		Login        string
		Password     string
		TokenType    int
		BankType     int
		PayType      int
		CurrensyCode int
		Language     string
		Description  string
		AccountId    int
		AutomatId    int
		DeviceType   int
		SbpPoint     string
		QrFormat     int
	}
	SecretKey         string
	KeyPayment        string
	Service_id        string
	TidPaymentSystem  string
	HostPaymentSystem string
	Products          []map[string]interface{}
	Date              string
	OrderId           string
	OperationId       string
	SbolBankInvoiceId string
	IdTransaction     string
	MerchantId        string
	GateWay           string
	PaymentToken      string
	WareId            string
	Sum               int
	SumOneProduct     int
	SumMax            int
	Imei              string
}

func (rp *RequestPay) Init(json_param []byte) error {
	err := json.Unmarshal(json_param, rp)
	if err != nil {
		return err
	}
	return nil
}
