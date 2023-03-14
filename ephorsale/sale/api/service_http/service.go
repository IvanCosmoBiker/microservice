package service_http

import (
	"encoding/json"
	transaction_dispetcher "ephorservices/ephorsale/transaction"
	transaction "ephorservices/ephorsale/transaction/transaction_struct"
	parserTypes "ephorservices/pkg/parser/typeParse"
	transportHttp "ephorservices/pkg/transportprotocol/http/v1"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
)

type FUNC func(tran *transaction.Transaction)

type ServiceSale struct {
	SaleHandler    FUNC
	FiscalHandler  FUNC
	PaymentHandler FUNC
	Address        map[string]string
}

func New() *ServiceSale {
	address := make(map[string]string)
	address["pay"] = "/pay"
	address["fiscal"] = "/fiscal"
	return &ServiceSale{
		Address: address,
	}
}

func (ss *ServiceSale) HandlerSale(w http.ResponseWriter, req *http.Request) {
	response := make(map[string]interface{})
	defer req.Body.Close()
	saleRequest := &RequestSale{}
	json_data, _ := ioutil.ReadAll(req.Body)
	json.Unmarshal(json_data, saleRequest)
	if saleRequest.Sum == 0 && len(saleRequest.Products) < 1 {
		response["message"] = "Нельзя провети оплату с нулевой суммой"
		ss.JSON(w, response)
		return
	}
	check := transaction_dispetcher.Dispetcher.CheckDuplicate(saleRequest.Config.AutomatId, saleRequest.Config.AccountId)
	fmt.Printf("CHECK_DUPLICATE %v_______________\n", check)
	if check == true {
		response["message"] = "Действие над автоматом производится другим пользователем, пожалуйста подождите"
		ss.JSON(w, response)
		return
	}
	tran := ss.SetTransactionSale(saleRequest)
	tran.InitStores(tran.Config.AccountId)
	err := transaction_dispetcher.Dispetcher.StartTransaction(tran)
	fmt.Printf(" KYKKY %s\n", err)
	if err != nil {
		response["message"] = err.Error()
		ss.JSON(w, response)
		return
	}
	fmt.Println(" KYKKY\n")
	go ss.SaleHandler(tran)
	response["message"] = "ok"
	response["tid"] = tran.Config.Noise
	fmt.Printf("%s\n", response)
	ss.JSON(w, response)
}

func (ss *ServiceSale) HandlerFiscal(w http.ResponseWriter, req *http.Request) {
	fmt.Println(" KYKKY")
	if req.Method == "GET" {
		fmt.Fprintf(w, "Sorry, only POST method is supported.")
		return
	}
	fmt.Println(" KYKKY2")
	defer req.Body.Close()
	fiscalRequest := RequestFiscalServer{
		Events: make([]string, 0, 1),
	}
	fmt.Println(" KYKK3")
	json_data, _ := ioutil.ReadAll(req.Body)
	err := json.Unmarshal(json_data, &fiscalRequest)

	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(" KYKK3")
	tran := ss.SetTransactionFiscalServer(&fiscalRequest)
	fmt.Println(" KYKKY")
	fmt.Printf("TRANSACTION SET: %+v\n", tran)
	go ss.FiscalHandler(tran)
	return
}

func (ss *ServiceSale) HandlerPayment(w http.ResponseWriter, req *http.Request) {
	ss.JSON(w, "not emplement")
}

func (ss *ServiceSale) SetTransactionSale(request *RequestSale) *transaction.Transaction {
	tran := transaction_dispetcher.Dispetcher.NewTransaction()
	tran.Date = transaction_dispetcher.Dispetcher.Date.Now()
	tran.Config.AccountId = request.Config.AccountId
	tran.Config.AutomatId = request.Config.AutomatId
	tran.Config.TokenType = request.Config.TokenType
	tran.Config.DeviceType = request.Config.DeviceType
	tran.Config.CurrensyCode = request.Config.CurrensyCode
	tran.Config.Imei = request.Imei
	tran.TaxSystem.Type = request.Config.TaxSystem
	tran.Payment.PayType = request.Config.PayType
	tran.Payment.TokenType = request.Config.TokenType
	tran.Payment.Token = request.PaymentToken
	tran.Payment.UserPhone = request.Config.UserPhone
	tran.Payment.ReturnUrl = request.Config.ReturnUrl
	tran.Payment.DeepLink = request.Config.DeepLink
	tran.Payment.Login = request.Config.Login
	tran.Payment.Password = request.Config.Password
	tran.Payment.Type = uint8(request.Config.BankType)
	tran.Payment.CurrensyCode = request.Config.CurrensyCode
	tran.Payment.GateWay = request.GateWay
	tran.Payment.MerchantId = request.MerchantId
	tran.Payment.SbpPoint = request.Config.SbpPoint
	tran.Payment.Service_id = request.Service_id
	tran.Payment.SecretKey = request.SecretKey
	tran.Payment.KeyPayment = request.KeyPayment
	tran.Payment.Sum = request.Sum
	tran.Fiscal.Config.QrFormat = request.Config.QrFormat
	ss.setProductsTransaction(request.Products, tran)
	return tran
}

func (ss *ServiceSale) SetTransactionFiscalServer(fiscalRequest *RequestFiscalServer) *transaction.Transaction {
	tran := transaction_dispetcher.Dispetcher.NewTransaction()
	port, err := strconv.Atoi(fiscalRequest.ConfigFR.Port)
	if err != nil {
		tran.Fiscal.Config.Dev_port = 0
	} else {
		tran.Fiscal.Config.Dev_port = port
	}
	tran.Config.Imei = fiscalRequest.Imei
	tran.Fiscal.Config.Type = uint8(fiscalRequest.TypeFr)
	tran.Fiscal.Config.Inn = fiscalRequest.Inn
	tran.Fiscal.Config.Login = fiscalRequest.ConfigFR.Login
	tran.Fiscal.Config.Password = fiscalRequest.ConfigFR.Password
	tran.Fiscal.ResiptId = fiscalRequest.CheckId
	tran.Fiscal.Config.Dev_addr = fiscalRequest.ConfigFR.Host
	tran.Fiscal.Config.Auth_public_key = fiscalRequest.ConfigFR.Cert
	tran.Fiscal.Config.Auth_private_key = fiscalRequest.ConfigFR.Key
	tran.Fiscal.Signature = fiscalRequest.ConfigFR.Sign
	tran.Fiscal.FiscalRequest = fiscalRequest.Fields
	tran.Fiscal.Events = make([]string, 0, 1)
	tran.Fiscal.Events = append(tran.Fiscal.Events, fiscalRequest.Events...)
	tran.Fiscal.OnlyFiscal = true
	tran.Fiscal.Send = true
	tran.Fiscal.NeedFiscal = true
	return tran
}

func (ss *ServiceSale) SetTransactionPayment() {

}

func (ss *ServiceSale) setProductsTransaction(products []map[string]interface{}, tran *transaction.Transaction) {
	for _, product := range products {
		quantity := parserTypes.ParseTypeInterfaceToInt64(product["quantity"])
		if float64(quantity/1000) < float64(1) {
			quantity = quantity * 1000
		}
		Product := &transaction.Product{}
		Product.Name = parserTypes.ParseTypeInString(product["name"])
		Product.Ware_id = parserTypes.ParseTypeInterfaceToInt32(product["ware_id"])
		Product.Select_id = parserTypes.ParseTypeInString(product["select_id"])
		Product.Price = parserTypes.ParseTypeInFloat64(product["price"])
		Product.Value = parserTypes.ParseTypeInFloat64(product["value"])
		Product.Tax_rate = parserTypes.ParseTypeInterfaceToInt32(product["tax_rate"])
		Product.Quantity = quantity
		tran.Products = append(tran.Products, Product)
	}
}

func (ss *ServiceSale) JSON(w http.ResponseWriter, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(data)
}

func (ss *ServiceSale) InitApi(HttpManager *transportHttp.ServerHttp) {
	HttpManager.SetHandlerListener(ss.Address["pay"], ss.HandlerSale, "POST")
	HttpManager.SetHandlerListener(ss.Address["fiscal"], ss.HandlerFiscal, "POST", "GET")
}
