package modul

import (
	"encoding/json"
	config "ephorservices/config"
	parserTypes "ephorservices/pkg/parser/typeParse"
	randString "ephorservices/pkg/randgeneratestring"
	"log"
	"math"
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
	cfg    *config.Config
}

func (cm *Core) makeRequestGetQrCode(pointId string, params []map[string]interface{}) (string, []byte, error) {
	orderString := randString.GenerateString{}
	var sum float64
	orderString.RandStringRunes()
	orderNumber := orderString.String
	result := make(map[string]interface{})
	result["retailPointId"] = pointId
	result["extraInfo"] = "Заказ № " + orderNumber
	for _, item := range params {
		sum += math.Floor((parserTypes.ParseTypeInFloat64(item["price"]) * parserTypes.ParseTypeInFloat64(item["quantity"])) / float64(100))
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

func (cm *Core) makeRequestGetStatusQrCode(params map[string]interface{}) (string, error) {
	url := "https://" + HOST + POST_QR_STATUS + parserTypes.ParseTypeInString(params["qrId"])
	return url, nil
}

func (cm *Core) makeRequestReturnMoney(params map[string]interface{}) (string, []byte, error) {
	jsonStr, err := json.Marshal(params)
	if err != nil {
		return "", nil, err
	}
	url := "https://" + HOST + GET_QR_REFUND
	return url, jsonStr, nil
}

func InitCore(conf *config.Config) *Core {
	return &Core{
		cfg: conf,
	}
}
