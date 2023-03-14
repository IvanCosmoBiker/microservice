package core

import (
	"encoding/json"
	transaction "ephorservices/ephorsale/transaction/transaction_struct"
	randstring "ephorservices/pkg/randgeneratestring"
	"fmt"
	"log"
)

type Core struct{}

func InitCore() *Core {
	return &Core{}
}

func (c *Core) MakeRequestCreateOrder(tran *transaction.Transaction) (string, error) {
	stringRequest := ""
	orderString := randstring.Init()
	orderString.RandStringRunes()
	orderNumber := orderString.String
	orderString = nil
	if tran.Payment.ReturnUrl == "" {
		tran.Payment.ReturnUrl = "https://paytest.ephor.online"
	}
	stringRequest += fmt.Sprintf("userName=%s&", tran.Payment.Login)
	stringRequest += fmt.Sprintf("password=%s&", tran.Payment.Password)
	stringRequest += fmt.Sprintf("orderNumber=%s&", orderNumber)
	stringRequest += fmt.Sprintf("amount=%v&", tran.Payment.Sum)
	stringRequest += fmt.Sprintf("returnUrl=%s&", tran.Payment.ReturnUrl)
	stringRequest += fmt.Sprintf("description=%s&", tran.Payment.Description)
	jsonParams := make(map[string]interface{})
	if tran.Payment.TokenType == transaction.TypeTokenSberPayAndroid || tran.Payment.TokenType == transaction.TypeTokenSberPayiOS {
		jsonParams["app2app"] = true
		jsonParams["app.deepLink"] = tran.Payment.DeepLink
		if tran.Payment.TokenType == transaction.TypeTokenSberPayAndroid {
			jsonParams["app.osType"] = "android"
		}
		if tran.Payment.TokenType == transaction.TypeTokenSberPayiOS {
			jsonParams["app.osType"] = "ios"
		}
	} else if tran.Payment.TokenType == transaction.TypeTokenSberPayWeb {
		jsonParams["back2app"] = true
		stringRequest += fmt.Sprintf("phone=%s&", tran.Payment.UserPhone)
	}
	data, _ := json.Marshal(jsonParams)
	jsonParams = nil
	stringRequest += fmt.Sprintf("jsonParams=%s", data)
	log.Printf("%s", stringRequest)
	return stringRequest, nil
}

func (c *Core) MakeRequestStatusOrder(tran *transaction.Transaction) (string, error) {
	stringRequest := ""
	stringRequest += fmt.Sprintf("userName=%s&", tran.Payment.Login)
	stringRequest += fmt.Sprintf("password=%s&", tran.Payment.Password)
	stringRequest += fmt.Sprintf("orderId=%s", tran.Payment.OrderId)
	return stringRequest, nil
}

func (c *Core) MakeRequestDepositOrder(tran *transaction.Transaction) (string, error) {
	stringRequest := ""
	stringRequest += fmt.Sprintf("userName=%s&", tran.Payment.Login)
	stringRequest += fmt.Sprintf("password=%s&", tran.Payment.Password)
	stringRequest += fmt.Sprintf("orderId=%s&", tran.Payment.OrderId)
	stringRequest += fmt.Sprintf("amount=%v", tran.Payment.DebitSum)
	return stringRequest, nil

}

func (c *Core) MakeRequestReturnMoney(tran *transaction.Transaction) (string, error) {
	stringRequest := ""
	stringRequest += fmt.Sprintf("userName=%s&", tran.Payment.Login)
	stringRequest += fmt.Sprintf("password=%s&", tran.Payment.Password)
	stringRequest += fmt.Sprintf("orderId=%s&", tran.Payment.OrderId)
	stringRequest += fmt.Sprintf("amount=%v", tran.Payment.Sum)
	return stringRequest, nil
}
