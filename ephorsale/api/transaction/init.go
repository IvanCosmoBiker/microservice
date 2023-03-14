package transaction

import (
	"encoding/json"
	controller "ephorservices/ephorsale/api/transaction/controller"
	transactionDispetcher "ephorservices/ephorsale/transaction"
	transportHttp "ephorservices/pkg/transportprotocol/http/v1"
	"fmt"
	"io/ioutil"
	"net/http"
)

var ControllerTransaction *controller.ControllerTransaction

func Init(transport *transportHttp.ServerHttp, dispecther *transactionDispetcher.TransactionDispetcher) {
	getTransactions := transport.SetHandlerListener("/transactionGet", GetTransactions)
	deleteTransaction := transport.SetHandlerListener("/transactionEnd", DeleteTransaction)
	getTransactions.Methods("GET", "POST")
	deleteTransaction.Methods("GET", "POST")
	ControllerTransaction = controller.Init(dispecther)
}

func GetTransactions(w http.ResponseWriter, req *http.Request) {
	switch req.Method {
	case "GET":
		filters := make(map[string]interface{})
		result := ControllerTransaction.GetTransactionActive(filters)
		body, err := json.Marshal(result)
		if err != nil {
			return
		}
		w.Write(body)
	default:
		fmt.Fprintf(w, "Sorry, only GET method is supported.")
	}
}

func DeleteTransaction(w http.ResponseWriter, req *http.Request) {
	switch req.Method {
	case "POST":
		json_data, _ := ioutil.ReadAll(req.Body)
		defer req.Body.Close()
		RequestEnd := controller.RequestEndTransaction{}
		json.Unmarshal(json_data, &RequestEnd)
		result, resulterr := ControllerTransaction.EndTransactionActive(RequestEnd)
		if resulterr != true {
			return
		}
		body, err := json.Marshal(result)
		if err != nil {
			return
		}
		w.Write(body)
	case "GET":
		fmt.Fprintf(w, "Sorry, only GET method is supported.")
	default:
		fmt.Fprintf(w, "Sorry, only POST method is supported.")
	}

}
