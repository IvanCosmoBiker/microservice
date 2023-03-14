package service_http

import (
	"encoding/json"
	model "ephorservices/ephorsale/control/api/state/transaction"
	transaction_dispetcher "ephorservices/ephorsale/transaction"
	"ephorservices/pkg/orm/request"
	transportHttp "ephorservices/pkg/transportprotocol/http/v1"
	"fmt"
	"io/ioutil"
	"net/http"
)

type Func func(req *request.Request) []map[string]interface{}

type ServiceStateTransaction struct {
	Address   map[string]string
	Functions map[string]Func
}

func New() *ServiceStateTransaction {
	serviceApi := &ServiceStateTransaction{
		Address:   make(map[string]string),
		Functions: make(map[string]Func),
	}
	serviceApi.initAddress()
	return serviceApi
}

func (sst *ServiceStateTransaction) initAddress() {
	sst.Address["transactionGet"] = "/transactionGet"
	sst.Address["transactionEnd"] = "/transactionEnd"
}

func (sst *ServiceStateTransaction) SetFunctionHandler(key string, function Func) {
	sst.Functions[key] = function
}

func (sst *ServiceStateTransaction) JSON(w http.ResponseWriter, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(data)
}

func (sst *ServiceStateTransaction) InitApi(HttpManager *transportHttp.ServerHttp) {
	HttpManager.SetHandlerListener(sst.Address["transactionGet"], sst.HandlerTransactionGet, "GET", "POST")
	HttpManager.SetHandlerListener(sst.Address["transactionEnd"], sst.HandlerTransactionEnd, "GET", "POST")
}

func (sst *ServiceStateTransaction) HandlerTransactionGet(w http.ResponseWriter, req *http.Request) {
	result := make([]map[string]interface{}, 0, 1)
	transactions := transaction_dispetcher.Dispetcher.GetTransactions()
	fmt.Printf("%+v", transactions)
	for _, tran := range transactions {
		entry := make(map[string]interface{})
		entry["id"] = tran.Config.Tid
		entry["automat_id"] = tran.Config.AutomatId
		entry["date"] = tran.Date
		entry["ps_type"] = tran.Payment.Type
		entry["ps_order"] = tran.Payment.OrderId
		entry["sum"] = tran.Sum
		entry["status"] = tran.Status
		entry["f_status"] = tran.Fiscal.Status
		entry["error"] = tran.Error
		entry["ps_desc"] = tran.Payment.Message
		result = append(result, entry)
	}
	fmt.Printf("%+v", result)
	sst.JSON(w, result)
}

func (sst *ServiceStateTransaction) HandlerTransactionEnd(w http.ResponseWriter, req *http.Request) {
	result := make([]map[string]interface{}, 0, 1)
	json_data, _ := ioutil.ReadAll(req.Body)
	defer req.Body.Close()
	RequestEnd := model.RequestEndTransaction{}
	json.Unmarshal(json_data, &RequestEnd)
	transaction, exist := transaction_dispetcher.Dispetcher.GetOneTransaction(RequestEnd.Tid)
	if exist {
		transaction.Close <- []byte("Stop transaction of administrator")
	}
	transaction_dispetcher.Dispetcher.RemoveTransaction(RequestEnd.Tid)
	sst.JSON(w, result)
}
