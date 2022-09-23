package orange

import (
	"encoding/json"
	config "ephorservices/config"
	"ephorservices/ephorsale/fiscal/interfaceFiscal"
	transaction "ephorservices/ephorsale/transaction"
	randString "ephorservices/pkg/randgeneratestring"
	"fmt"
	"log"
	"strings"
)

type Orange struct {
	Name string
	cfg  *config.Config
}

type NewOrangeStruct struct {
	Orange
}

func New(conf *config.Config) interfaceFiscal.Fiscal {
	return &NewOrangeStruct{
		Orange: Orange{
			Name: "Orange",
			cfg:  conf,
		},
	}
}

func (o *Orange) SendCheck(data *transaction.Transaction) map[string]interface{} {
	result := make(map[string]interface{})
	key := ""
	TypeFr := int(int64(FrModel["type"].(int64)))
	var dataCheck = make(map[string]interface{})
	var orderString randString.GenerateString
	orderString.RandStringRunes()
	resiptId := orderString.String
	log.Printf("%+v", TransactionData)
	payments, positions := o.GenerateDataForCheck(TransactionData)
	content := make(map[string]interface{})
	content["type"] = 1
	content["automatNumber"] = data.Fiscal.Config.AutomatNumber
	content["SettlementAddress"] = TransactionData.Address
	content["SettlementPlace"] = TransactionData.PointName
	checkClose := make(map[string]interface{})
	checkClose["payments"] = payments
	checkClose["taxationSystem"] = o.ConvertTaxationSystem(o.Config.TaxSystem)
	content["checkClose"] = checkClose
	content["positions"] = positions
	if TypeFr == interfaceFiscal.Fr_EphorServerOrangeData || TypeFr == interfaceFiscal.Fr_EphorOrangeData {
		key = "4010004"
	} else {
		key = FrModel["inn"].(string)
	}
	dataCheck["id"] = resiptId
	dataCheck["group"] = o.Config.Group
	dataCheck["Inn"] = FrModel["inn"]
	dataCheck["key"] = key
	dataCheck["content"] = content
	jsonDataCheck, _ := json.Marshal(dataCheck)
	if TypeFr == interfaceFiscal.Fr_EphorServerOrangeData || TypeFr == interfaceFiscal.Fr_EphorOrangeData {
		certFile, keyFile, errFile := o.ReadFileCertificate()
		if errFile != nil {

		}
		o.Config.Cert = certFile
		o.Config.Key = keyFile
		sign, err := o.ComputeSignature(string(jsonDataCheck), privateKeySign)
		if err != nil {

		}
		o.Config.Sign = sign
	} else {
		sign, err := o.ComputeSignature(string(jsonDataCheck), []byte(string(FrModel["sign_private_key"].(string))))
		if err != nil {

		}
		o.Config.Sign = sign
	}
	url := fmt.Sprintf("https://%s:%s/api/v2/documents/", o.Config.Host, o.Config.Port)
	Response := o.Call("POST", url, jsonDataCheck)
	if Response.Code == 409 {
		return o.GetStatus(resiptId)
	}

	if Response.Code != 201 {
		result["code"] = Response.Code
		result["fr_id"] = resiptId
		result["status"] = "unsuccess"
		result["fr_status"] = interfaceFiscal.Status_Overflow
		result["message"] = strings.Join(Response.Errors[:], "\n")
		return result
	}
	return o.GetStatus(resiptId)
}

func (o *Orange) GetStatus(data *transaction.Transaction) map[string]interface{} {

}

func (o *Orange) GetQrPicture(data *transaction.Transaction) string {

}

func (o *Orange) GetQrUrl(data *transaction.Transaction) string {

}
