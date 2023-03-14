package sberpay

import (
	"encoding/json"
	config "ephorservices/config"
	pb "ephorservices/ephorpayment/service"
	transaction "ephorservices/ephorsale/transaction/transaction_struct"
	randstring "ephorservices/pkg/randgeneratestring"
	"fmt"
	"log"
)

type Core struct {
	cfg *config.Config
}

func InitCore(conf *config.Config) *Core {
	return &Core{
		cfg: conf,
	}
}

func (c *Core) MakeRequestCreateOrder(tran *pb.Request) (string, error) {
	stringRequest := ""
	orderString := randstring.Init()
	orderString.RandStringRunes()
	orderNumber := orderString.String
	if tran.ReturnUrl == "" {
		tran.ReturnUrl = "https://paytest.ephor.online"
	}
	stringRequest += fmt.Sprintf("userName=%s&", tran.Login)
	stringRequest += fmt.Sprintf("password=%s&", tran.Password)
	stringRequest += fmt.Sprintf("orderNumber=%s&", orderNumber)
	stringRequest += fmt.Sprintf("amount=%v&", tran.Sum)
	stringRequest += fmt.Sprintf("returnUrl=%s&", tran.ReturnUrl)
	stringRequest += fmt.Sprintf("description=%s&", tran.Description)
	jsonParams := make(map[string]interface{})
	if tran.TokenType == int32(transaction.TypeTokenSberPayAndroid) || tran.TokenType == int32(transaction.TypeTokenSberPayiOS) {
		jsonParams["app2app"] = true
		jsonParams["app.deepLink"] = tran.DeepLink
		if tran.TokenType == int32(transaction.TypeTokenSberPayAndroid) {
			jsonParams["app.osType"] = "android"
		}
		if tran.TokenType == int32(transaction.TypeTokenSberPayiOS) {
			jsonParams["app.osType"] = "ios"
		}
	} else if tran.TokenType == int32(transaction.TypeTokenSberPayWeb) {
		jsonParams["back2app"] = true
		stringRequest += fmt.Sprintf("phone=%s&", tran.UserPhone)
	}
	data, _ := json.Marshal(jsonParams)
	stringRequest += fmt.Sprintf("jsonParams=%s", data)
	log.Printf("%s", stringRequest)
	return stringRequest, nil
}

func (c *Core) MakeRequestStatusOrder(tran *pb.Request) (string, error) {
	stringRequest := ""
	stringRequest += fmt.Sprintf("userName=%s&", tran.Login)
	stringRequest += fmt.Sprintf("password=%s&", tran.Password)
	stringRequest += fmt.Sprintf("orderId=%s", tran.OrderId)
	return stringRequest, nil
}

func (c *Core) MakeRequestDepositOrder(tran *pb.Request) (string, error) {
	stringRequest := ""
	stringRequest += fmt.Sprintf("userName=%s&", tran.Login)
	stringRequest += fmt.Sprintf("password=%s&", tran.Password)
	stringRequest += fmt.Sprintf("orderId=%s&", tran.OrderId)
	stringRequest += fmt.Sprintf("amount=%v", tran.DebitSum)
	return stringRequest, nil

}

func (c *Core) MakeRequestReturnMoney(tran *pb.Request) (string, error) {
	stringRequest := ""
	stringRequest += fmt.Sprintf("userName=%s&", tran.Login)
	stringRequest += fmt.Sprintf("password=%s&", tran.Password)
	stringRequest += fmt.Sprintf("orderId=%s&", tran.OrderId)
	stringRequest += fmt.Sprintf("amount=%v", tran.Sum)
	return stringRequest, nil
}
