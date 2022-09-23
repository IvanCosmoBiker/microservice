package sberpay

import (
	"encoding/json"
	config "ephorservices/config"
	transaction "ephorservices/ephorsale/transaction"
	randstring "ephorservices/pkg/randgeneratestring"
	"fmt"
	"log"
)

type Core struct {
	cfg *config.Config
}

func (c *Core) makeRequestCreateOrder(tran *transaction.Transaction) (string, error) {
	stringRequest := ""
	orderString := randstring.Init()
	orderString.RandStringRunes()
	orderNumber := orderString.String
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
	stringRequest += fmt.Sprintf("jsonParams=%s", data)
	log.Printf("%s", stringRequest)
	return stringRequest, nil
}

func (c *Core) makeRequestStatusOrder(tran *transaction.Transaction) (string, error) {
	stringRequest := ""
	stringRequest += fmt.Sprintf("userName=%s&", tran.Payment.Login)
	stringRequest += fmt.Sprintf("password=%s&", tran.Payment.Password)
	stringRequest += fmt.Sprintf("orderId=%s", tran.Payment.OrderId)
	return stringRequest, nil
}

func (c *Core) makeRequestDepositOrder(tran *transaction.Transaction) (string,error) {
	stringRequest := ""
    stringRequest += fmt.Sprintf("userName=%s&",tran.Payment.Login)
    stringRequest += fmt.Sprintf("password=%s&",tran.Payment.Password)
    stringRequest += fmt.Sprintf("orderId=%s&",tran.Payment.OrderId)
    stringRequest += fmt.Sprintf("amount=%v",tran.Payment.DebitSum)
    return stringRequest ,nil

}

func (c *Core) makeRequestReturnMoney(tran *transaction.Transaction) (string, error) {
	stringRequest := ""
	stringRequest += fmt.Sprintf("userName=%s&", tran.Payment.Login)
	stringRequest += fmt.Sprintf("password=%s&", tran.Payment.Password)
	stringRequest += fmt.Sprintf("orderId=%s&", tran.Payment.OrderId)
	stringRequest += fmt.Sprintf("amount=%v", tran.Payment.Sum)
	return stringRequest, nil
}

func InitCore(conf *config.Config) *Core {
	return &Core{
		cfg: conf,
	}
}
