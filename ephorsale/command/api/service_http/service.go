package service_http

import (
	"encoding/json"
	transaction_dispetcher "ephorservices/ephorsale/transaction"
	transaction "ephorservices/ephorsale/transaction/transaction_struct"
	logger "ephorservices/pkg/logger"
	transportHttp "ephorservices/pkg/transportprotocol/http/v1"
	"io/ioutil"
	"log"
	"net/http"
)

type FUNC func(tran *transaction.Transaction)

type ServiceCommand struct {
	CommandDeviceHandler FUNC
	Address              map[string]string
}

func New() *ServiceCommand {
	address := make(map[string]string)
	address["command"] = "/command"
	return &ServiceCommand{
		Address: address,
	}
}

func (sc *ServiceCommand) JSON(w http.ResponseWriter, data interface{}) {
	json.NewEncoder(w).Encode(data)
}

func (sc *ServiceCommand) HandlerCommandDevice(w http.ResponseWriter, req *http.Request) {
	defer req.Body.Close()
	response := make(map[string]interface{})
	json_data, err := ioutil.ReadAll(req.Body)
	logger.Log.Infof("%+v", json_data)
	if err != nil {
		response["status"] = false
		response["message"] = err.Error()
		sc.JSON(w, response)
		logger.Log.Errorf("%+v", response)
		return
	}
	CommandRequest := CommandServerRequest{}
	json.Unmarshal(json_data, &CommandRequest)
	log.Println(CommandRequest.Id)
	tran := transaction.InitTransaction()
	tran.Config.Imei = CommandRequest.Imei
	reqModem := tran.NewRequest()
	reqModem.AddFilterParam("imei", reqModem.Operator.OperatorEqual, true, CommandRequest.Imei)
	modem, err := transaction_dispetcher.Dispetcher.StoreModem.GetOneBy(reqModem)
	if err != nil {
		response["status"] = false
		response["message"] = err.Error()
		logger.Log.Errorf("%+v", response)
		sc.JSON(w, response)
		return
	}
	tran.Config.Command_id = CommandRequest.Id
	tran.Config.AccountId = int(modem.Account_id.Int32)
	tran.InitStores(tran.Config.AccountId)
	go sc.CommandDeviceHandler(tran)
	response["status"] = true
	sc.JSON(w, response)
	return
}

func (sc *ServiceCommand) InitApi(HttpManager *transportHttp.ServerHttp) {
	HttpManager.SetHandlerListener(sc.Address["command"], sc.HandlerCommandDevice, "POST")
}
