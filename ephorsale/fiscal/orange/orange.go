package orange1

import (
	"bytes"
	"crypto/tls"
	"encoding/base64"
	"encoding/json"
	config "ephorservices/config"
	interfaceFiscal "ephorservices/ephorsale/fiscal"
	transactionStruct "ephorservices/ephorsale/transaction"
	parserTypes "ephorservices/pkg/parser/typeParse"
	randString "ephorservices/pkg/randgeneratestring"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"
)

var (
	layoutISO = "2006-01-02 15:04:05"
	layoutQr  = "20060102T150405"
)

// type ResponseOrange struct {
// 	Code    int
// 	Status  string
// 	Message string
// 	Errors  []string
// 	Data    map[string]interface{}
// }

// type ConfigOrange struct {
// 	Cert          string
// 	Host          string
// 	Port          string
// 	Key           string
// 	Group         string
// 	TaxSystem     int
// 	Sign          string
// 	AutomatNumber int
// 	Inn           string
// }

type Orange struct {
	Name string
	cfg  *config.Config
}

type NewOrangeStruct struct {
	Orange
}

func (ofd *Orange) GetQrPicture(date string, summ int, frResponse map[string]interface{}) string {
	result := ofd.MakeUrlQr(date, summ, frResponse)
	str := base64.StdEncoding.EncodeToString([]byte(result))
	return str
}

func (ofd *Orange) GetQrUrl(date string, summ int, frResponse map[string]interface{}) string {
	return ofd.MakeUrlQr(date, summ, frResponse)
}

func (ofd *Orange) SendCheckApi(data requestFiscal.Data) (map[string]interface{}, requestFiscal.Data) {
	ofd.Config.Cert = data.ConfigFR.Cert
	ofd.Config.Key = data.ConfigFR.Key
	ofd.Config.Host = data.ConfigFR.Host
	ofd.Config.Port = data.ConfigFR.Port
	ofd.Config.Inn = data.Inn
	ofd.Config.Sign = data.ConfigFR.Sign
	result := make(map[string]interface{})
	url := fmt.Sprintf("https://%s:%s/api/v2/documents/", data.ConfigFR.Host, data.ConfigFR.Port)
	Response := ofd.Call("POST", url, data.Fields.Request)
	if Response.Code == 409 {
		result["code"] = Response.Code
		result["fr_id"] = data.CheckId
		result["status"] = "success"
		result["fr_status"] = interfaceFiscal.Status_InQueue
		result["message"] = Response.Message
		return result, data
	}
	if Response.Code != 201 {
		result["code"] = Response.Code
		result["fr_id"] = nil
		result["fp_string"] = nil
		result["status"] = "unsuccess"
		result["fr_status"] = interfaceFiscal.Status_Error
		result["message"] = Response.Message
		return result, data
	}
	data.DataResponse = Response.Data
	result["code"] = Response.Code
	result["fr_id"] = data.CheckId
	result["status"] = "success"
	result["fr_status"] = interfaceFiscal.Status_InQueue
	result["message"] = "Нет ошибок"
	return result, data
}

func (ofd Orange) GetStatusApi(data requestFiscal.Data) map[string]interface{} {
	result := make(map[string]interface{})
	url := fmt.Sprintf("https://%s:%s/api/v2/documents/%s/status/%s", ofd.Config.Host, ofd.Config.Port, ofd.Config.Inn, data.CheckId)
	Response := ofd.Call("GET", url, []byte(""))
	if Response.Code > 299 {
		result["code"] = Response.Code
		result["fr_id"] = nil
		result["status"] = "unsuccess"
		result["fr_status"] = interfaceFiscal.Status_Error
		result["message"] = Response.Message
		return result
	}
	result["code"] = Response.Code
	result["fr_id"] = nil
	result["status"] = "success"
	result["fp"] = parserTypes.ParseTypeInString(Response.Data["fp"])
	result["fd"] = parserTypes.ParseTypeInFloat64(Response.Data["documentNumber"])
	result["fn"] = parserTypes.ParseTypeInString(Response.Data["fsNumber"])
	result["fr_status"] = interfaceFiscal.Status_Error
	result["message"] = Response.Message
	return result
}

func (ofd *Orange) SendCheck() map[string]interface{} {
	result := make(map[string]interface{})
	key := ""
	TypeFr := int(int64(FrModel["type"].(int64)))
	var dataCheck = make(map[string]interface{})
	var orderString randString.GenerateString
	orderString.RandStringRunes()
	resiptId := orderString.String
	log.Printf("%+v", TransactionData)
	payments, positions := ofd.GenerateDataForCheck(TransactionData)
	content := make(map[string]interface{})
	content["type"] = 1
	content["automatNumber"] = ofd.Config.AutomatNumber
	content["SettlementAddress"] = TransactionData.Address
	content["SettlementPlace"] = TransactionData.PointName
	checkClose := make(map[string]interface{})
	checkClose["payments"] = payments
	checkClose["taxationSystem"] = ofd.ConvertTaxationSystem(ofd.Config.TaxSystem)
	content["checkClose"] = checkClose
	content["positions"] = positions
	if TypeFr == interfaceFiscal.Fr_EphorServerOrangeData || TypeFr == interfaceFiscal.Fr_EphorOrangeData {
		key = "4010004"
	} else {
		key = FrModel["inn"].(string)
	}
	dataCheck["id"] = resiptId
	dataCheck["group"] = ofd.Config.Group
	dataCheck["Inn"] = FrModel["inn"]
	dataCheck["key"] = key
	dataCheck["content"] = content
	jsonDataCheck, _ := json.Marshal(dataCheck)
	if TypeFr == interfaceFiscal.Fr_EphorServerOrangeData || TypeFr == interfaceFiscal.Fr_EphorOrangeData {
		certFile, keyFile, errFile := ofd.ReadFileCertificate()
		if errFile != nil {

		}
		ofd.Config.Cert = certFile
		ofd.Config.Key = keyFile
		sign, err := ofd.ComputeSignature(string(jsonDataCheck), privateKeySign)
		if err != nil {

		}
		ofd.Config.Sign = sign
	} else {
		sign, err := ofd.ComputeSignature(string(jsonDataCheck), []byte(string(FrModel["sign_private_key"].(string))))
		if err != nil {

		}
		ofd.Config.Sign = sign
	}
	url := fmt.Sprintf("https://%s:%s/api/v2/documents/", ofd.Config.Host, ofd.Config.Port)
	Response := ofd.Call("POST", url, jsonDataCheck)
	if Response.Code == 409 {
		return ofd.GetStatus(resiptId)
	}

	if Response.Code != 201 {
		result["code"] = Response.Code
		result["fr_id"] = resiptId
		result["status"] = "unsuccess"
		result["fr_status"] = interfaceFiscal.Status_Overflow
		result["message"] = strings.Join(Response.Errors[:], "\n")
		return result
	}
	return ofd.GetStatus(resiptId)
}

func (ofd *Orange) setTimeOut() chan bool {
	timeout := make(chan bool)
	go func() {
		select {
		case <-time.After(5 * time.Minute):
			timeout <- true
		}
	}()
	return timeout
}

func (ofd *Orange) SendRequestOfGetStatus(orderId string) map[string]interface{} {
	result := make(map[string]interface{})
	url := fmt.Sprintf("https://%s:%s/api/v2/documents/%s/status/%s", ofd.Config.Host, ofd.Config.Port, ofd.Config.Inn, orderId)
	Response := ofd.Call("GET", url, []byte(""))
	if Response.Code == 200 {
		result["code"] = Response.Code
		result["fr_id"] = orderId
		result["status"] = "success"
		result["fp"] = Response.Data["fp"].(string)
		result["fd"] = int(Response.Data["documentNumber"].(float64))
		result["fn"], _ = strconv.Atoi(Response.Data["fsNumber"].(string))
		result["date"] = Response.Data["processedAt"]
		result["message"] = "нет ошибок"
		result["fr_status"] = interfaceFiscal.Status_Complete
	}
	if Response.Code == 503 {
		result["code"] = Response.Code
		result["fr_id"] = orderId
		result["status"] = "unsuccess"
		result["fr_status"] = interfaceFiscal.Status_Overflow
		result["message"] = Response.Message
	}
	if Response.Code > 299 {
		result["code"] = Response.Code
		result["fr_id"] = orderId
		result["status"] = "unsuccess"
		result["fr_status"] = interfaceFiscal.Status_Error
		result["message"] = strings.Join(Response.Errors[:], "\n")
	}
	return result
}

func (ofd *Orange) GetStatus(parametrs ...string) map[string]interface{} {
	var orderId string
	orderId = parametrs[0]
	result := make(map[string]interface{})
	chanTimeOut := ofd.setTimeOut()
	for {
		select {
		case <-time.After(2 * time.Second):
			result = ofd.SendRequestOfGetStatus(orderId)
			if result["status"] == "unsuccess" {
				return result
			}
			if result["code"] == 200 {
				return result
			}
		case <-chanTimeOut:
			result["fr_id"] = orderId
			result["status"] = "unsuccess"
			result["code"] = 0
			result["message"] = fmt.Sprintf("Cancelled by a Timeout of %s", ofd.Name)
			return result
		}
	}
}

func (ofd Orange) Call(method string, url string, json_request []byte) ResponseOrange {
	Response := ResponseOrange{}
	req, _ := http.NewRequest(method, url, bytes.NewBuffer(json_request))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("x-Signature", ofd.Config.Sign)
	log.Printf("x-Signature %s", ofd.Config.Sign)
	req.Close = true
	log.Println(ofd.Config.Cert)
	log.Println(ofd.Config.Key)
	cert, err := tls.X509KeyPair([]byte(ofd.Config.Cert), []byte(ofd.Config.Key))
	if err != nil {
		Response.Code = 0
		Response.Status = "unsuccess"
		Response.Message = fmt.Sprintf("%v", err)
		return Response
	}
	client := &http.Client{}
	tlsConfig := &tls.Config{
		Certificates:       []tls.Certificate{cert},
		InsecureSkipVerify: true,
	}
	client.Transport = &http.Transport{
		TLSClientConfig: tlsConfig,
	}
	resp, err := client.Do(req)
	if err != nil {
		Response.Code = 0
		Response.Status = "unsuccess"
		Response.Message = fmt.Sprintf("%v", err)
		return Response
	}
	body, _ := ioutil.ReadAll(resp.Body)
	defer resp.Body.Close()
	json.Unmarshal([]byte(body), &Response.Data)
	log.Printf("%+v", Response)
	Response.Code = resp.StatusCode
	if resp.StatusCode > 299 {
		Response.Status = "unsuccess"
		ArrayInterface := Response.Data["errors"].([]interface{})
		errorStrings := parserTypes.ParseArrayInrefaceToArrayString(ArrayInterface)
		Response.Message = strings.Join(errorStrings[:], "\n")
		return Response
	}
	return Response
}

var TransactionData transactionStruct.Transaction
var FrModel map[string]interface{}

func (newf *NewOrangeStruct) NewFiscal() interfaceFiscal.Fiscal /* тип interfaceFiscal.Fiscal*/ {
	return &NewOrangeStruct{
		Orange: Orange{
			Name:   "Orange",
			Config: ConfigOrange{},
		},
	}
}
