package fiscalmanager

import (
    "fmt"
    "log"
    "time"
    "net/http"
    factoryFiscal "factory/fiscal"
    interfaseFiscal "interface/fiscalinterface"
    transactionStruct "data/transaction"
    fiscal "internal/fr"
    requestFiscal "data/requestApi"
   "encoding/json"
   "bytes"
    automat "internal/automat"
    transactionDispetcher "transactionDispetcher"
    configEphor "configEphor"
    connectionPostgresql "connectionDB/connect"
    ConnectionRabbitMQ "lib-rabbitmq"
    listener "listeners/v1"
    transactionProduct "internal/transactionproduct"
    automatEvent "internal/automatevent" 
    transactionStore "internal/transaction"
    "runtime"
    "join"
    "math"
    "strconv"
    "io/ioutil"
)
var (
    QrPicture = 1
    QrUrl = 2
)

var ( 
    TaxRate_NDSNone = 0;
    TaxRate_NDS0 = 1;
    TaxRate_NDS10 = 2;
    TaxRate_NDS20 = 3;
    TaxRate_NDS12 = 4;
)

type FiscalManager struct {
    Status int
}

type DataFiscal struct {
    Point_name string
    Point_address string
    Name_product string
    Ware_id int
    Select_id string
    Tax_rate int
    Tax_system int
}

func (fm *FiscalManager) getNds(tax_rate,value int) int {
    var valueResult int
    valueInt := float64(value/100)
    switch tax_rate {
        case TaxRate_NDS10:
        floatValue := math.Floor(((valueInt*10/110)*100))
        valueResult = int(float64(floatValue)*100)
        fallthrough
        case TaxRate_NDS12:
        floatValue := math.Floor(((valueInt*12/112)*100))
        valueResult = int(float64(floatValue)*100)
        fallthrough
        case TaxRate_NDS20:
        floatValue := math.Floor(((valueInt*20/120)*100))
        valueResult = int(float64(floatValue)*100)
    }
    return valueResult

}

func (fm *FiscalManager) SetFiscalResultData(transaction transactionStruct.Transaction,Fr interfaseFiscal.Fiscal,frResponse map[string]interface{}){
    params := make(map[string]interface{})
    if frResponse["status"] == "unsuccess" {
        params["id"] = transaction.Tid
        params["qr_format"] = 2
        params["f_type"] = frResponse["f_type"]
        params["f_receipt"] = frResponse["fr_id"].(string)
        params["f_desc"] = frResponse["message"].(string)
        params["f_status"] = frResponse["fr_status"]
        TransactionStore.SetByParams(params)
        return
    }
    var Qr string
    params["id"] = transaction.Tid
    params["qr_format"] = 2
    params["fn"] = frResponse["fn"]
    params["fd"] = frResponse["fd"]
    params["fp"] = frResponse["fp"]
    params["f_type"] = frResponse["f_type"]
    params["f_receipt"] = frResponse["fr_id"]
    params["f_desc"] = frResponse["message"]
    params["f_status"] = frResponse["fr_status"]
    Qr = Fr.GetQrUrl(transaction.Date,transaction.Sum,frResponse)
    params["f_qr"] = Qr
    TransactionStore.SetByParams(params)
}

func (fm *FiscalManager) GetStringSelectIdAndWareId(transaction transactionStruct.Transaction) (string,string) {
    stringSelectId := ""
    stringWareId := ""
     for _, product := range transaction.Products {
         if stringSelectId == ""{
             stringSelectId += "confp.select_id IN ("
             stringSelectId += fmt.Sprintf("'%s'",product["select_id"]) 
         }
         if stringWareId == ""{
             stringWareId += "confp.ware_id IN ("
             stringWareId += fmt.Sprintf("%v",product["ware_id"])
         }
        stringSelectId += fmt.Sprintf(",'%s'",product["select_id"]) 
        stringWareId += fmt.Sprintf(",%v",product["ware_id"])
     }
     stringSelectId += ")"
     stringWareId += ")"
     return stringSelectId,stringWareId
}

func (fm *FiscalManager) GetDataForFiscal(transaction transactionStruct.Transaction,FrModel map[string]interface{}) ([]DataFiscal,bool) {
    stringSelectId,stringWareId := fm.GetStringSelectIdAndWareId(transaction)
    sql := fmt.Sprintf("SELECT cp.name AS point_name,cp.address AS point_address,confp.name AS name_product,confp.ware_id AS ware_id,confp.select_id AS select_id,confp.tax_rate AS tax_rate,conf.tax_system AS tax_system FROM /schema/.config_product AS confp INNER JOIN /schema/.automat_config AS ac ON ac.automat_id = %v AND ac.to_date IS NULL INNER JOIN /schema/.config AS conf ON confp.config_id = ac.config_id INNER JOIN /schema/.automat_location AS al ON al.automat_id = %v INNER JOIN /schema/.company_point AS cp ON cp.id = al.company_point_id WHERE %s AND %s AND confp.config_id = ac.config_id GROUP BY confp.name,cp.name,cp.address,conf.tax_system,confp.tax_rate,confp.ware_id,confp.select_id",transaction.AutomatId,transaction.AutomatId,stringSelectId,stringWareId)
    JoinSql.GetJoin(transaction.AccountId,sql)
    FrDatas := []DataFiscal{}
    for JoinSql.RowsData.Next(){
        FrData := DataFiscal{}
        err := JoinSql.RowsData.Scan(&FrData.Point_name, 
            &FrData.Point_address,
            &FrData.Name_product,
            &FrData.Ware_id,  
            &FrData.Select_id,
            &FrData.Tax_rate,
            &FrData.Tax_system)
        if err != nil{
            fmt.Println(err)
            continue
        }
        FrDatas = append(FrDatas,FrData)
    }
    if len(FrDatas) < 1 {
        fm.AddDataAutomatEventErr(transaction,FrModel,"нет данных для фискализации: [точка,конфигурация]. Проверьте настройки конфигурации автомата",interfaseFiscal.Status_Error)
        return nil,false
    }
    return FrDatas,true
}

func GetDateTime(stringTime string,seconds int) string {
    t,_ := time.Parse("2006-01-02 15:04:05",stringTime)
    t2 := t.Add(-180 * time.Minute)
    newTime := t2.Add(time.Duration(seconds) * time.Second)
    resultTime := newTime.Format("2006-01-02 15:04:05")
    return resultTime
}

func (fm *FiscalManager) GetAutomat(id,accountId interface{}) ([]map[string]interface{},bool) {
   return AutomatStore.GetOneById(id,accountId);
}

func (fm *FiscalManager) GetFr(id,accountId interface{}) ([]map[string]interface{},bool) {
    return FrStore.GetOneById(id,accountId)
}

func (fm *FiscalManager) InitFiscal(typeFiscal int) interfaseFiscal.Fiscal {
    return factoryFiscal.GetFiscal(typeFiscal)
}

func (fm *FiscalManager) AddDataAutomatEventErr(transaction transactionStruct.Transaction,frModel map[string]interface{},err string,status int){
    addSeconds := 1
    for _, product := range transaction.Products {
        for i := 0; i < int(product["quantity"].(float64)); i++ {
            params := make(map[string]interface{})
            _,exist := product["tax_rate"]
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
            params["status"] = status
            params["error_detail"] = err
            if len(frModel) > 1 {
                params["type_fr"] = frModel["type"]
            }
            params["payment_device"] = "DA"
            params["tax_system"] = transaction.Tax_system
            if exist != false {
                params["tax_rate"] = product["tax_rate"]
                taxValue := fm.getNds(int(product["tax_rate"].(float64)),int(product["price"].(float64)))
                params["tax_value"] = taxValue
            }
            if transaction.PointId != 0 {
                params["point_id"] = transaction.PointId
            }else {
                params["point_id"] = nil
            }
            log.Printf("%+v",params)
            AutomatEventStore.AddByParams(params)
            addSeconds +=1
        }
    }
}

func (fm *FiscalManager) CalcSumProducts(transaction transactionStruct.Transaction) int {
     sum := 0 
     for _, product := range transaction.Products {
         fmt.Println(int(product["price"].(float64)*product["quantity"].(float64)))
         sum += int(product["price"].(float64)*product["quantity"].(float64))
     }
     return sum
}

func (fm *FiscalManager) CheckMaxSumm(transaction transactionStruct.Transaction,automatModel map[string]interface{}) bool {
    maxSum := fm.CalcSumProducts(transaction)
    fmt.Println(maxSum)
    fmt.Println(int(automatModel["summ_max_fr"].(int64)))
    if int(automatModel["summ_max_fr"].(int64)) == 0 {
        return true
    }
    if maxSum > int(automatModel["summ_max_fr"].(int64)) {
        return false
    }
    return true
}


func (fm *FiscalManager) CheckFiscalisation(transaction transactionStruct.Transaction,frModel map[string]interface{}) ([]map[string]interface{}) {
    var productFiscal []map[string]interface{}
    addSeconds := 1
    if int(frModel["fr_disable_cashless"].(int64)) != 1 {
        fm.AddDataAutomatEventErr(transaction,frModel,"отключение фискализации со стороны клиента",interfaseFiscal.Status_None)
        return productFiscal
    }
     for _, product := range transaction.Products {
            params := make(map[string]interface{})
            date := GetDateTime(transaction.Date,addSeconds)
            params["account_id"] = transaction.AccountId
            params["automat_id"] = transaction.AutomatId
            params["type"] = 3
            params["date"] = date
            params["category"] = 1
            params["name"] = product["name"]
            params["credit"] = product["price"]
            params["price"] = product["price"]
            params["select_id"] = product["select_id"]
            params["ware_id"] = product["ware_id"]
            params["value"] = product["price"]
			params["type_fr"] = frModel["type"]
            params["quantity"] = product["quantity"]
            params["payment_device"] = "DA"
            params["tax_system"] = transaction.Tax_system
            params["tax_rate"] = product["tax_rate"]
            taxValue := fm.getNds(int(product["tax_rate"].(float64)),int(product["price"].(float64)))
            params["tax_value"] = taxValue
            params["point_id"] = transaction.PointId
            productFiscal = append(productFiscal,params)
            addSeconds +=1
        }
     return productFiscal
}

func (fm *FiscalManager) FiscalProductsRabbit()  {
}
// This function for fiscalisation products in Transaction
func (fm *FiscalManager) FiscalProducts(transaction transactionStruct.Transaction) (map[string]interface{},transactionStruct.Transaction) {
    var result = make(map[string]interface{})
    var FrModel = make(map[string]interface{})
    automatModel,err :=  fm.GetAutomat(transaction.AutomatId,transaction.AccountId)
    if err == false || len(automatModel) < 1 {
        fm.AddDataAutomatEventErr(transaction,FrModel,"не могли получить автомат",interfaseFiscal.Status_None)
        result["status"] = false
        result["message"] = "не могли получить автомат"
        result["fr_status"] = interfaseFiscal.Status_None
        return result,transaction
    }
    frId := automatModel[0]["fr_id"]
    if frId == nil {
        fm.AddDataAutomatEventErr(transaction,FrModel,"нет активной кассы",interfaseFiscal.Status_None)
        result["status"] = false
        result["fr_status"] = interfaseFiscal.Status_None
        result["message"] = "нет активной кассы"
        return result,transaction
    }
    fiscalModel,err :=  fm.GetFr(automatModel[0]["fr_id"],transaction.AccountId)
    if err == false {
        fm.AddDataAutomatEventErr(transaction,FrModel,"касса не привязана к автомату",interfaseFiscal.Status_None)
        result["status"] = false
        result["fr_status"] = interfaseFiscal.Status_None
        result["message"] = "касса не привязана к автомату"
        return result,transaction
    }
    FrModel = fiscalModel[0]
    if transaction.PointId == 0 {
        fm.AddDataAutomatEventErr(transaction,FrModel,"нет торговой точки",interfaseFiscal.Status_None)
        result["status"] = false
        result["message"] = "нет торговой точки"
        result["fr_status"] = interfaseFiscal.Status_None
        return result,transaction
    }
    checkSum := fm.CheckMaxSumm(transaction,automatModel[0])
    if checkSum == false {
        fm.AddDataAutomatEventErr(transaction,FrModel,"превышен лимит суммы по чеку",interfaseFiscal.Status_MAX_CHECK)
        result["status"] = false
        result["message"] = "превышен лимит суммы по чеку"
        result["fr_status"] = interfaseFiscal.Status_MAX_CHECK
        result["type_fr"] = FrModel["type"]
        return result,transaction
    }
    checkFiscal := fm.CheckFiscalisation(transaction,FrModel)
    log.Printf("[Fiscal]: %T",FrModel["fr_disable_cashless"])
    log.Printf("[Fiscal]: %+v",checkFiscal)
    log.Printf("[Fiscal]: %v",len(checkFiscal))
    if len(checkFiscal) < 1 {
        result["status"] = false
        result["message"] = "отключение фискализации со стороны клиента"
        result["fr_status"] = interfaseFiscal.Status_OFF_FR
        return result,transaction
    }
    transaction.Products = checkFiscal
    typeFr := int(FrModel["type"].(int64))
    fr := fm.InitFiscal(typeFr)
    Data,frResult := fm.GetDataForFiscal(transaction,FrModel)
    if frResult == false {
        result["status"] = false
        result["message"] = "нет данных для фискализации: [точка,конфигурация]. Проверьте настройки конфигурации автомата"
        result["fr_status"] = interfaseFiscal.Status_Error
        result["type_fr"] = FrModel["type"]
        return result,transaction
    }
    tax_system := Data[0].Tax_system
    addressAutomat := Data[0].Point_address
    pointName := Data[0].Point_name
    transaction.Tax_system = tax_system
    transaction.Address = addressAutomat
    transaction.PointName = pointName
    for _, product := range transaction.Products {
         for i := 0; i < len(Data); i++ {
             FrProduct := Data[i]
             if FrProduct.Select_id == product["select_id"] {
                 product["tax_rate"] = FrProduct.Tax_rate
             }
         }
    }
    fr.InitData(transaction,FrModel)
    resultFiscal := fr.SendCheck()
    resultFiscal["f_type"] = FrModel["type"]
    fm.SetFiscalResultData(transaction,fr,resultFiscal)
    result["status"] = true
    result["message"] = "нет ошибок"
    result["fn"] = resultFiscal["fn"]
    result["fd"] = resultFiscal["fd"]
    result["fp"] = resultFiscal["fp"]
    result["type_fr"] = resultFiscal["f_type"]
    result["id_fr"] = resultFiscal["fr_id"]
    result["fr_status"] = resultFiscal["fr_status"]
    return result,transaction
}

type Event struct {
	Id string
}

type Outcome struct {
	Imei string
	Data struct {
		Message, Status  string
		Events           []Event
		Code, StatusCode,Fiscalization int
		Fields           struct {
			Fp, Fn string
			Fd     float64
		}
	}
}

func (out *Outcome) MakeEvents(ev []string) {
	for _, s := range ev {
		out.Data.Events = append(out.Data.Events, Event{Id: s})
	}
}

func (out *Outcome) Finish() {
	runtime.Goexit()
}

func (out Outcome) Send() {

	defer out.Finish()
	json_request, _ := json.Marshal(out.Data)

	dc := strconv.FormatInt(time.Now().UnixNano()/1000000, 10)
	url := fmt.Sprintf("%s&login=%s&password=12345678&_dc=%s", conf.Services.EphorFiscal.ResponseUrl, out.Imei, dc)

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(json_request))
	req.Header.Set("Content-Type", "application/json")
	req.Close = true

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Panic(err)
	}
	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)
    if conf.Debug {
		log.Println(out.Data)
		log.Println(resp)
		log.Println(string(body))
	}
	return
}

func (out Outcome) Timeout() {
	out.Data.Status = "unsuccess"
	out.Data.Code = 0
	out.Data.Message = fmt.Sprintf("Cancelled by a Timeout of %s", "FiscalManager")
	out.Send()
}

func (out *Outcome) initDataOutCome(data requestFiscal.Data) {
    out.Imei = data.Imei
	out.Data.Fiscalization = data.ConfigFR.Fiscalization
	out.MakeEvents(data.Events)
}

// This function for fiscalisation products which come from api 

func StartFiscal(timeout chan bool, json_data []byte) {
    var out Outcome
    var data requestFiscal.Data
    var fiscalManager FiscalManager
    json.Unmarshal(json_data, &data)
    log.Println(data)
    fmt.Sprintf("%+v",data)
    out.initDataOutCome(data)
    kassa := fiscalManager.InitFiscal(data.TypeFr)
    response,DataFiscal := kassa.SendCheckApi(data)
    if response["status"].(string) == "unsuccess" {
        out.Data.Status = response["status"].(string)
        out.Data.Code = response["code"].(int)
        out.Data.Message = response["message"].(string)
        out.Send()
        return 
    }
    for {
		select {
        case <-time.After(conf.Services.EphorFiscal.SleepMilliSec * time.Millisecond):
           resultStatus := kassa.GetStatusApi(DataFiscal)
           log.Println(resultStatus)
           fmt.Sprintf("%+v",resultStatus)
           if resultStatus["status"].(string) == "unsuccess" {
                out.Data.Status = resultStatus["status"].(string)
                out.Data.Code = resultStatus["code"].(int)
                out.Data.Message = resultStatus["message"].(string)
                out.Send()
                return 
           }
           if resultStatus["code"] == 200 {
                out.Data.Status = "success"
                out.Data.Code = resultStatus["code"].(int)
                out.Data.Fields.Fp = resultStatus["fp"].(string)
                out.Data.Fields.Fd = resultStatus["fd"].(float64)
                out.Data.Fields.Fn = resultStatus["fn"].(string)
                out.Send()
                return
           }
		case <-timeout:
			out.Timeout()
			return
		}
	}
}


func handler(w http.ResponseWriter, req *http.Request) {
    switch req.Method {
        case "POST":
            json_data, _ := ioutil.ReadAll(req.Body)
            defer req.Body.Close()
            timeout := make(chan bool)
            go StartFiscal(timeout, json_data)
            go func() {
                select {
                case <-time.After(conf.Services.EphorFiscal.ExecuteMinutes * time.Minute):
                    timeout <- true
                }
            }()
            return 
        case "GET":
            fmt.Fprintf(w, "%s: Running\n", "Servirce Fiscal")
        default:
        fmt.Fprintf(w, "Sorry, only POST and GET method is supported.")
    }

}

var connectDb connectionPostgresql.DatabaseInstance
var conf *configEphor.Config
var Rabbit *ConnectionRabbitMQ.ChannelMQ
var Transactions transactionDispetcher.TransactionDispetcher
var TransactionStruct transactionStruct.Transaction
var TransactionProductStore transactionProduct.StoreTransactionProduct
var AutomatEventStore automatEvent.StoreAutomatEvent
var TransactionStore  transactionStore.StoreTransaction
var FrStore fiscal.StoreFr
var AutomatStore automat.StoreAutomat
var JoinSql join.Join
var (
    where  map[string]interface{}
    parametrs map[string]interface{}
)

func StartListener(cfg *configEphor.Config,connectPg connectionPostgresql.DatabaseInstance,rabbitmq *ConnectionRabbitMQ.ChannelMQ,channelParent chan bool,transactions transactionDispetcher.TransactionDispetcher){
    TransactionProductStore.Connection = connectPg
    AutomatEventStore.Connection = connectPg
    TransactionStore.Connection = connectPg
    FrStore.Connection = connectPg
    AutomatStore.Connection = connectPg
    connectDb = connectPg
    conf = cfg
    Rabbit = rabbitmq
    Transactions = transactions
    JoinSql.ConnectionDb = &connectPg
    point := fmt.Sprintf("%s:%s",cfg.Services.Address,cfg.Services.Port)
    listener.StartListener("/fiscal",point,handler)
    log.Println("Start Module Fiscal..")
}
