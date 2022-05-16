package bankListener

import(
	"encoding/json"
    "fmt"
    "time"
    "net/http"
    "io/ioutil"
    "log"
    "runtime"
    "strconv"
    "math/rand"
    transactionDispetcher "transactionDispetcher"
    configEphor "configEphor"
    deviceFactory "factory/device"
    interfaseBank "interface/bankinterface"
    deviceInterfase "interface/deviceinterface"
    connectionPostgresql "connectionDB/connect"
    ConnectionRabbitMQ "lib-rabbitmq"
    listener "listeners/v1"
    transactionStruct "data/transaction"
    paymentmanager "paymentmanager"
    transactionProduct "internal/transactionproduct"
    automatEvent "internal/automatevent" 
    transactionStore "internal/transaction"
)

type RequestRabbit struct {
	Tid     int
	St 		int
	D       string
	Err 	int
	A 		int
	Wid 	int
	Sum		int
}

func Finish(){
    runtime.Goexit()
}

func GetDateTime(stringTime string,seconds int) string {
    t,_ := time.Parse("2006-01-02 15:04:05",stringTime)
    newTime := t.Add(time.Duration(seconds) * time.Second)
    resultTime := newTime.Format("2006-01-02 15:04:05")
    return resultTime
}

func InitDevice(device int) deviceInterfase.Device {
    return deviceFactory.GetDevice(device)
}

func SetDataTransactionProducts(transactionProduct map[string]interface{}){
    TransactionProductStore.SetByParams(transactionProduct)
}

func AddDataTransactionProducts(transaction transactionStruct.Transaction){
    params := make(map[string]interface{})
    for _, product := range transaction.Products {
        //for i := 0; i < product["quantity"].(int); i++{
            params["transaction_id"] = transaction.Tid
            params["name"] = product["name"]
            params["select_id"] = product["select_id"]
            params["ware_id"] = product["ware_id"]
            params["value"] = product["price"]
            params["tax_rate"] = product["tax_rate"]
            params["quantity"] = product["quantity"]
            TransactionProductStore.AddByParams(params)
        //} 
    }
}

func Random(min int, max int) int {
    return rand.Intn(max-min) + min
}

func AddDataAutomatEvent(transaction transactionStruct.Transaction){
    params := make(map[string]interface{})
    addSeconds := 1
    for _, product := range transaction.Products {
        fmt.Println(product["quantity"].(int))
        for i := 0; i < product["quantity"].(int); i++ {
            date := GetDateTime(transaction.Date,addSeconds)
            params["account_id"] = transaction.AccountId
            params["automat_id"] = transaction.AutomatId
            params["type"] = 3
            params["date"] = date
            params["category"] = 1
            params["name"] = product["name"]
            params["credit"] = product["price"]
            params["select_id"] = product["select_id"]
            params["ware_id"] = product["ware_id"]
            params["value"] = product["price"]
            params["status"] = 0
            //params["tax_rate"] = product["tax_rate"]
            AutomatEventStore.AddByParams(params)
            addSeconds +=1
        }
        
    }
}

func SetDataTransaction(parametrs,where map[string]interface{})  {
    connectDb.Set("transaction",parametrs,where)
    connectDb.SetLog(fmt.Sprintf("%q",parametrs["error"]))
}

func initDataTransaction(request interfaseBank.Request) transactionStruct.Transaction {
    products := request.Products
    transaction := transactionStruct.Transaction{}
    transaction.Tid = request.IdTransaction
    transaction.Date = request.Date
    transaction.Sum = request.Sum
    transaction.Token = request.PaymentToken
    transaction.PayType = request.Config.PayType
    transaction.AccountId = request.Config.AccountId
    transaction.AutomatId = request.Config.AutomatId
    transaction.DeviceType = request.Config.DeviceType
    transaction.SumMax = request.SumMax
    for _, product := range products {
        transaction.Products = append(transaction.Products,product)
    }
    return transaction
}

func CheckTransactionOfAutomat(request interfaseBank.Request) bool {
    keyFound := request.Config.AutomatId + request.Config.AccountId
    resultFound := Transactions.GetReplayProtection(keyFound)
    if resultFound == false {
        return false
    }
    return true
}

func StartTransaction(bankChannel chan bool, json_data interfaseBank.Request){
    where = make(map[string]interface{})
    parametrs = make(map[string]interface{})
    request := json_data
    transactionData := initDataTransaction(request)
    result,Bank := PaymentManager.StartСommunicationBank(request,transactionData)
    if result["success"] == false {
        where["id"] = request.IdTransaction
        parametrs["ps_order"] = result["ps_order"]
        parametrs["ps_desc"] = result["ps_desc"]
        parametrs["error"] = result["error"]
        parametrs["status"] = transactionStruct.TransactionState_Error
        SetDataTransaction(parametrs,where)
        AddDataTransactionProducts(transactionData)
        return 
    }
    where["id"] = request.IdTransaction
    parametrs["ps_order"] = result["ps_order"]
    parametrs["ps_desc"] = result["ps_desc"]
    parametrs["error"] = result["error"]
    parametrs["status"] = transactionStruct.TransactionState_MoneyDebitOk
    SetDataTransaction(parametrs,where)
    device := InitDevice(transactionData.DeviceType)
    if device == nil {
        where["id"] = request.IdTransaction 
        parametrs["error"] = "no available device type"
        parametrs["status"] = transactionStruct.TransactionState_Error
        SetDataTransaction(parametrs,where)
        AddDataTransactionProducts(transactionData)
        Bank.ReturnMoney(request.IdTransaction)
        return 
    }
    tidInt,_ := strconv.Atoi(request.IdTransaction)
    keyAutomat := request.Config.AutomatId + request.Config.AccountId
    channel := Transactions.AddChannel(tidInt)
    Transactions.AddReplayProtection(keyAutomat,request.Config.AutomatId)
    device.SendMessage(transactionData,connectDb,Rabbit,conf,request,Bank,result,Transactions,channel)
    return 
}

func GetDateNow() string {
    date := time.Now()
    return date.Format("2006-01-02 15:04:05")
}

func AddTransaction(request *interfaseBank.Request) string {
    rand.Seed(time.Now().UnixNano())
    randNoise := Random(10000000,20000000)
    parametrInsert := make(map[string]interface{})
    parametrInsert["automat_id"] = request.Config.AutomatId
    parametrInsert["account_id"] = request.Config.AccountId
    parametrInsert["token_id"] = request.PaymentToken
    parametrInsert["status"] = 1
    parametrInsert["date"] = request.Date
    parametrInsert["pay_type"] = request.Config.PayType
    parametrInsert["sum"] = request.Sum
    parametrInsert["ps_type"] = request.Config.BankType
    Id := TransactionStore.AddByParams(parametrInsert)
    request.IdTransaction = strconv.Itoa(Id)
    noise := randNoise + Id
    updateNoise := make(map[string]interface{})
    updateNoise["id"] = Id
    updateNoise["noise"] = strconv.Itoa(noise)
    TransactionStore.SetByParams(updateNoise)
    return strconv.Itoa(noise)
}

func handler(w http.ResponseWriter, req *http.Request) {
    switch req.Method {
    case "POST":
        json_data, _ := ioutil.ReadAll(req.Body)
        defer req.Body.Close()
        bankChannel := make(chan bool)
        var requestCheck interfaseBank.Request
        json.Unmarshal(json_data, &requestCheck)
        connectDb.AddLog(fmt.Sprintf("%+v",requestCheck),"PaymentSystem" ,fmt.Sprintf("%+v",requestCheck),"EphorErp")
        resultFoundActiveTid := CheckTransactionOfAutomat(requestCheck)
        response := make(map[string]interface{})
        fmt.Println(resultFoundActiveTid)
        if resultFoundActiveTid == true {
            response["message"] = "the machine is busy"
            body, err := json.Marshal(response)
            if err != nil {
                return
            }
            w.Write(body)
            return 
        }else {
            noise := AddTransaction(&requestCheck)
            fmt.Println(requestCheck.IdTransaction)
            response["message"] = "ok"
            response["tid"] = noise
            body, err := json.Marshal(response)
            if err != nil {
                return
            }
            w.Write(body)
            go StartTransaction(bankChannel, requestCheck)
            return 
        }
        return 
    case "GET":
        fmt.Fprintf(w, "%s: Running\n", "Servirce Payment")
        log.Println("Running")
    default:
        fmt.Fprintf(w, "Sorry, only POST and GET method is supported.")
    }

}

var connectDb *connectionPostgresql.DatabaseInstance
var conf *configEphor.Config
var Rabbit *ConnectionRabbitMQ.ChannelMQ
var ChannelParentClose chan bool
var Transactions transactionDispetcher.TransactionDispetcher
var TransactionStruct transactionStruct.Transaction
var PaymentManager paymentmanager.PaymentManager
var TransactionProductStore transactionProduct.StoreTransactionProduct
var AutomatEventStore automatEvent.StoreAutomatEvent
var TransactionStore  transactionStore.StoreTransaction
var (
    where  map[string]interface{}
    parametrs map[string]interface{}
)

func StartBank(cfg *configEphor.Config,connectPg *connectionPostgresql.DatabaseInstance,rabbitmq *ConnectionRabbitMQ.ChannelMQ,channelParent chan bool,transactions transactionDispetcher.TransactionDispetcher){
    TransactionProductStore.Connection = connectPg
    AutomatEventStore.Connection = connectPg
    TransactionStore.Connection = connectPg
    ChannelParentClose = make(chan bool)
    connectDb = connectPg
    conf = cfg
    Rabbit = rabbitmq
    Transactions = transactions
    point := fmt.Sprintf("%s:%s",cfg.Services.EphorPay.Bank.Address,cfg.Services.EphorPay.Bank.Port)
    listener.StartListener("/pay",point,handler)
    log.Println("Start Bank..")
}



