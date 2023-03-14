package core

import (
	"encoding/json"
	transaction "ephorservices/ephorsale/transaction/transaction_struct"
	parserTypes "ephorservices/pkg/parser/typeParse"
	randString "ephorservices/pkg/randgeneratestring"
	"log"
)

var (
	POST_INFO_COMPANY = "/v1/account-info"
	GET_INFO_POINTS   = "/v1/sbp/retail-points?companyId="
	POST_QR_CODE      = "/v1/sbp/qr-codes/dynamic"
	POST_QR_STATUS    = "/v1/sbp/qr-codes/"
	GET_QR_REFUND     = "/v1/sbp/qr-codes/refund"
	CONTENT_TYPE      = "application/json;charset=UTF-8"
	HOST              = "api.modulbank.ru"
)

type Core struct {
	Status int
}

func InitCore() *Core {
	return &Core{}
}

func (cm *Core) MakeRequestGetQrCode(pointId string, tran *transaction.Transaction) (string, []byte, error) {
	orderString := randString.Init()
	var sum float64
	orderString.RandStringRunes()
	orderNumber := orderString.String
	orderString = nil
	result := make(map[string]interface{})
	result["retailPointId"] = pointId
	result["extraInfo"] = "Заказ № " + orderNumber
	for _, product := range tran.Products {
		var value float64 = 0
		if product.Value == value {
			value = product.Price
		} else {
			value = product.Value
		}
		floatSum := (value / float64(100)) * (parserTypes.ParseTypeInFloat64(product.Quantity) / 1000)
		log.Println(floatSum)
		sum += floatSum
	}
	result["sum"] = sum
	log.Printf("%+v", result)
	jsonStr, err := json.Marshal(result)
	if err != nil {
		return "", nil, err
	}
	url := "https://" + HOST + POST_QR_CODE
	return url, jsonStr, nil
}

func (cm *Core) MakeRequestGetStatusQrCode(params map[string]interface{}) (string, error) {
	url := "https://" + HOST + POST_QR_STATUS + parserTypes.ParseTypeInString(params["qrId"])
	return url, nil
}

func (cm *Core) MakeRequestReturnMoney(params map[string]interface{}) (string, []byte, error) {
	jsonStr, err := json.Marshal(params)
	if err != nil {
		return "", nil, err
	}
	url := "https://" + HOST + GET_QR_REFUND
	return url, jsonStr, nil
}
