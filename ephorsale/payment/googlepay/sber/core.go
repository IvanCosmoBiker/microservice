package sber

import (
	"encoding/json"
	config "ephorservices/config"
	transaction "ephorservices/ephorsale/transaction"
	randString "ephorservices/pkg/randgeneratestring"
	"log"
)

type Core struct {
	cfg *config.Config
}

func (c *Core) makeRequestDepositOrder(tran *transaction.Transaction) ([]byte, error) {
	requestOrder := make(map[string]interface{})
	requestOrder["amount"] = tran.Payment.Sum
	requestOrder["orderId"] = tran.Payment.OrderId
	data, err := json.Marshal(requestOrder)
	if err != nil {
		log.Printf("%+v", err)
		return nil, err
	}
	return data, nil
}

func (c *Core) makeRequestStatusOrder(tran *transaction.Transaction) ([]byte, error) {
	requestOrder := make(map[string]interface{})
	requestOrder["token"] = tran.Payment.Token
	requestOrder["orderId"] = tran.Payment.OrderId
	data, err := json.Marshal(requestOrder)
	if err != nil {
		log.Printf("%+v", err)
		return nil, err
	}
	return data, nil
}

func (c *Core) makeOrderRequestCreateOrder(tran *transaction.Transaction) ([]byte, error) {
	orderString := randString.GenerateString{}
	orderString.RandStringRunes()
	orderNumber := orderString.String
	requestOrder := make(map[string]interface{})
	requestOrder["merchant"] = tran.Payment.MerchantId
	requestOrder["orderNumber"] = orderNumber
	requestOrder["language"] = tran.Payment.Language
	requestOrder["preAuth"] = true
	requestOrder["description"] = tran.Payment.Description
	requestOrder["paymentToken"] = tran.Payment.Token
	requestOrder["amount"] = tran.Payment.Sum
	requestOrder["currencyCode"] = tran.Payment.CurrensyCode
	requestOrder["returnUrl"] = "https://test.ru"
	data, err := json.Marshal(requestOrder)
	if err != nil {
		log.Printf("%+v", err)
		return nil, err
	}
	if c.cfg.Debug {
		log.Printf("%+v", requestOrder)
	}
	return data, nil
}

func InitCore(conf *config.Config) *Core {
	return &Core{
		cfg: conf,
	}
}
