package automatpostpaid

import (
	payment "ephorservices/ephorsale/payment"
	"ephorservices/ephorsale/sale/interfaceSale"
	transaction_dispetcher "ephorservices/ephorsale/transaction"
	transaction "ephorservices/ephorsale/transaction/transaction_struct"
	transport_manager "ephorservices/ephorsale/transport"
	logger "ephorservices/pkg/logger"
	parserTypes "ephorservices/pkg/parser/typeParse"
	"fmt"
)

type SaleAutomatPostPaid struct {
	Debug              bool
	KeyReplay          int
	ExecuteTimeSeconds int
	Fiscalization      interfaceSale.FiscalFunc
}

type NewSaleAutomatPostPaid struct {
	SaleAutomatPostPaid
}

func (newA *NewSaleAutomatPostPaid) New(executeTimeSeconds int, debug bool) interfaceSale.Sale {
	return &NewSaleAutomatPostPaid{
		SaleAutomatPostPaid: SaleAutomatPostPaid{
			ExecuteTimeSeconds: executeTimeSeconds,
			Debug:              debug,
		},
	}
}

func (sapp *SaleAutomatPostPaid) SetFiscalisation(function interfaceSale.FiscalFunc) {
	sapp.Fiscalization = function
}

func (sapp *SaleAutomatPostPaid) Sale(tran *transaction.Transaction) (result map[string]interface{}) {
	defer sapp.returnRoutine(&result, tran.Config.Tid)
	if tran.Payment.Sum == 0 {
		result["id"] = tran.Config.Tid
		result["error"] = "Не найдены товары для продолжения оплаты"
		result["status"] = transaction.TransactionState_Error
		return
	}
	sapp.KeyReplay = tran.Config.AutomatId + tran.Config.AccountId
	tran.KeyReplay = sapp.KeyReplay
	resultPayment := sapp.Payment(tran)
	if resultPayment["status"] == transaction.TransactionState_Error {
		result = sapp.clearMap(result)
		result["id"] = tran.Config.Tid
		result["ps_desc"] = resultPayment["ps_desc"]
		result["error"] = "Ошибка"
		result["status"] = resultPayment["status"]
		return
	}
	logger.Log.Infof("%+v", result)
	result = sapp.clearMap(result)
	resultPayment["id"] = tran.Config.Tid
	transaction_dispetcher.Dispetcher.StoreTransaction.SetByParams(resultPayment)
	resultPaymentStatus := payment.Payment.Satus(tran)
	logger.Log.Infof("%+v", resultPaymentStatus)
	if parserTypes.ParseTypeInterfaceToInt(resultPaymentStatus["status"]) != transaction.TransactionState_MoneyDebitWait {
		result = sapp.clearMap(result)
		result["id"] = tran.Config.Tid
		result["status"] = transaction.TransactionState_Error
		result["error"] = "Ошибка оплаты"
		return
	}
	result["id"] = tran.Config.Tid
	result["status"] = resultPaymentStatus["status"]
	transaction_dispetcher.Dispetcher.StoreTransaction.SetByParams(result)
	transaction_dispetcher.Dispetcher.AddReplayProtection(sapp.KeyReplay, tran.Config.AutomatId)
	resultSendMessage := sapp.SendMassage(tran)
	if parserTypes.ParseTypeInterfaceToInt(resultSendMessage["status"]) != transaction.TransactionState_MoneyHoldWait {
		result = sapp.clearMap(result)
		result["id"] = tran.Config.Tid
		result["status"] = transaction.TransactionState_Error
		result["error"] = resultSendMessage["error"]
		return
	}
	responseRabbit := sapp.WaitMassage(tran)
	fmt.Printf("___%+v_____\n", responseRabbit)
	if parserTypes.ParseTypeInterfaceToInt(responseRabbit["status"]) == transaction.TransactionState_ReturnMoney {
		result = sapp.clearMap(result)
		result["id"] = tran.Config.Tid
		result["ps_desc"] = responseRabbit["ps_desc"]
		result["error"] = responseRabbit["error"]
		result["status"] = transaction.TransactionState_Error
		return
	}
	logger.Log.Infof("%+v", resultPaymentStatus)
	if parserTypes.ParseTypeInterfaceToInt(responseRabbit["status"]) != transaction.VendState_Session {
		resultReturnMoney := payment.Payment.Return(tran)
		result = sapp.clearMap(result)
		result["id"] = tran.Config.Tid
		result["ps_desc"] = resultReturnMoney["ps_desc"]
		result["error"] = tran.GetDescriptionErr(parserTypes.ParseTypeInterfaceToInt(responseRabbit["error"]))
		result["status"] = transaction.TransactionState_Error
		return
	}
	result = sapp.clearMap(result)
	result["id"] = tran.Config.Tid
	result["ps_desc"] = tran.GetDescriptionCode(parserTypes.ParseTypeInterfaceToInt(responseRabbit["status"]))
	result["error"] = tran.GetDescriptionErr(parserTypes.ParseTypeInterfaceToInt(responseRabbit["error"]))
	result["status"] = tran.GetStatusServer(transaction.VendState_Session)
	logger.Log.Infof("%+v", result)
	transaction_dispetcher.Dispetcher.StoreTransaction.SetByParams(result)
	responseRabbit = sapp.WaitMassage(tran)
	if parserTypes.ParseTypeInterfaceToInt(responseRabbit["status"]) == transaction.TransactionState_ReturnMoney {
		result = sapp.clearMap(result)
		result["id"] = tran.Config.Tid
		result["ps_desc"] = responseRabbit["ps_desc"]
		result["error"] = responseRabbit["error"]
		result["status"] = transaction.TransactionState_Error
		return
	}
	if parserTypes.ParseTypeInterfaceToInt(responseRabbit["status"]) != transaction.VendState_Approving && parserTypes.ParseTypeInterfaceToInt(responseRabbit["status"]) != transaction.VendState_Vending {
		resultReturnMoney := payment.Payment.Return(tran)
		result = sapp.clearMap(result)
		result["id"] = tran.Config.Tid
		result["ps_desc"] = resultReturnMoney["ps_desc"]
		result["error"] = tran.GetDescriptionErr(parserTypes.ParseTypeInterfaceToInt(responseRabbit["error"]))
		result["status"] = transaction.TransactionState_Error
		return
	}
	result["id"] = tran.Config.Tid
	result["ps_desc"] = tran.GetDescriptionCode(parserTypes.ParseTypeInterfaceToInt(responseRabbit["status"]))
	result["error"] = tran.GetDescriptionErr(parserTypes.ParseTypeInterfaceToInt(responseRabbit["error"]))
	result["status"] = tran.GetStatusServer(parserTypes.ParseTypeInterfaceToInt(responseRabbit["status"]))
	transaction_dispetcher.Dispetcher.StoreTransaction.SetByParams(result)
	logger.Log.Infof("%+v", result)
	result = sapp.setProduct(tran, responseRabbit)
	if _, ok := result["status"]; ok {
		transaction_dispetcher.Dispetcher.StoreTransaction.SetByParams(result)
	}
	result = sapp.clearMap(result)
	tran.Payment.DebitSum = parserTypes.ParseTypeInterfaceToInt(responseRabbit["sum"])
	resultDebitMoney := payment.Payment.Debit(tran)
	if parserTypes.ParseTypeInterfaceToInt(resultDebitMoney["status"]) != transaction.TransactionState_MoneyDebitOk {
		resultReturnMoney := payment.Payment.Return(tran)
		result = sapp.clearMap(result)
		result["id"] = tran.Config.Tid
		result["ps_desc"] = resultReturnMoney["ps_desc"]
		result["error"] = "не удалось списать деньги"
		result["status"] = transaction.TransactionState_Error
		return
	}
	result = sapp.clearMap(result)
	result["id"] = tran.Config.Tid
	result["ps_desc"] = resultDebitMoney["ps_desc"]
	result["error"] = tran.GetDescriptionErr(parserTypes.ParseTypeInterfaceToInt(responseRabbit["error"]))
	result["status"] = tran.GetStatusServer(transaction.VendState_Session)
	transaction_dispetcher.Dispetcher.StoreTransaction.SetByParams(result)
	responseRabbit = sapp.WaitMassage(tran)
	if parserTypes.ParseTypeInterfaceToInt(responseRabbit["status"]) == transaction.TransactionState_ReturnMoney {
		result = sapp.clearMap(result)
		result["id"] = tran.Config.Tid
		result["ps_desc"] = responseRabbit["ps_desc"]
		result["error"] = responseRabbit["error"]
		result["status"] = transaction.TransactionState_Error
		return
	}
	if parserTypes.ParseTypeInterfaceToInt(responseRabbit["status"]) != transaction.VendState_VendOk {
		result = sapp.clearMap(result)
		resultReturnMoney := payment.Payment.Return(tran)
		result["id"] = tran.Config.Tid
		result["ps_desc"] = resultReturnMoney["ps_desc"]
		result["error"] = tran.GetDescriptionErr(parserTypes.ParseTypeInterfaceToInt(responseRabbit["error"]))
		result["status"] = transaction.TransactionState_Error
		return
	}
	result = sapp.clearMap(result)
	result["id"] = tran.Config.Tid
	result["ps_desc"] = tran.GetDescriptionCode(parserTypes.ParseTypeInterfaceToInt(responseRabbit["status"]))
	result["error"] = tran.GetDescriptionErr(parserTypes.ParseTypeInterfaceToInt(responseRabbit["error"]))
	result["status"] = tran.GetStatusServer(transaction.VendState_VendOk)
	return
}

func (sapp *SaleAutomatPostPaid) Payment(tran *transaction.Transaction) map[string]interface{} {
	result := make(map[string]interface{})
	resultPayment := payment.Payment.Payment(tran)
	if resultPayment["status"] != transaction.TransactionState_MoneyHoldStart {
		return resultPayment
	}
	result["id"] = tran.Config.Tid
	result["status"] = transaction.TransactionState_MoneyHoldWait
	transaction_dispetcher.Dispetcher.StoreTransaction.SetByParams(result)
	resultPay := payment.Payment.Hold(tran)
	if resultPay["status"] == transaction.TransactionState_Error {
		return resultPay
	}
	return resultPay
}

func (sapp *SaleAutomatPostPaid) SendMassage(tran *transaction.Transaction) map[string]interface{} {
	result := make(map[string]interface{})
	resultMessage := make(map[string]interface{})
	resultMessage["tid"] = tran.Config.Tid
	resultMessage["sum"] = tran.Payment.Sum
	resultMessage["m"] = 1
	resultMessage["a"] = 2
	err := transport_manager.TransportManager.QueueManager.SendMessage(resultMessage, tran.Config.Imei)
	if err != nil {
		return payment.Payment.Return(tran)
	}
	result["status"] = transaction.TransactionState_MoneyHoldWait
	return result
}

func (sapp *SaleAutomatPostPaid) WaitMassage(tran *transaction.Transaction) map[string]interface{} {
	response := make(map[string]interface{})
	result, err := transport_manager.TransportManager.QueueManager.WaitMessage(tran.ChannelMessage, sapp.ExecuteTimeSeconds)
	if err != nil {
		request := payment.Payment.Return(tran)
		request["status"] = transaction.TransactionState_ReturnMoney
		request["error"] = err.Error()
	}
	response["status"] = result["st"]
	response["error"] = result["err"]
	response["tid"] = result["tid"]
	response["ware_id"] = result["wid"]
	response["sum"] = result["sum"]
	response["a"] = result["a"]
	response["d"] = result["d"]
	return response
}

func (sapp *SaleAutomatPostPaid) clearMap(Map map[string]interface{}) map[string]interface{} {
	Map = nil
	Map = make(map[string]interface{})
	return Map
}

func (sapp *SaleAutomatPostPaid) Fiscal(tran *transaction.Transaction) map[string]interface{} {
	return make(map[string]interface{})
}

func (sapp *SaleAutomatPostPaid) setProduct(tran *transaction.Transaction, responseDevice map[string]interface{}) (result map[string]interface{}) {
	result = make(map[string]interface{})
	reqAutomatConfig := tran.NewRequest()
	reqAutomatConfig.AddFilterParam("automat_id", reqAutomatConfig.Operator.OperatorEqual, true, tran.Config.AutomatId)
	reqAutomatConfig.AddFilterParam("to_date", reqAutomatConfig.Operator.OperatorEqual, true)
	automatConfig, err := tran.Stores.StoreAutomatConfig.GetOneBy(reqAutomatConfig)
	if err != nil {
		logger.Log.Errorf("нет конфигурации")
		resultReturnMoney := payment.Payment.Return(tran)
		result["id"] = tran.Config.Tid
		result["ps_desc"] = resultReturnMoney["ps_desc"]
		result["error"] = "нет конфигурации"
		result["status"] = transaction.TransactionState_Error
		return
	}
	if automatConfig.Id == 0 {
		logger.Log.Errorf("нет конфигурации")
		resultReturnMoney := payment.Payment.Return(tran)
		result["id"] = tran.Config.Tid
		result["ps_desc"] = resultReturnMoney["ps_desc"]
		result["error"] = "нет конфигурации"
		result["status"] = transaction.TransactionState_Error
		return
	}
	reqConfigProduct := tran.NewRequest()
	reqConfigProduct.AddFilterParam("config_id", reqConfigProduct.Operator.OperatorEqual, true, automatConfig.Config_id.Int32)
	reqConfigProduct.AddFilterParam("ware_id", reqConfigProduct.Operator.OperatorEqual, true, responseDevice["ware_id"])
	configProduct, er := tran.Stores.StoreConfigProduct.GetOneBy(reqConfigProduct)
	if er != nil {
		logger.Log.Errorf("нет такого продукта в ассортименте")
		resultReturnMoney := payment.Payment.Return(tran)
		result["id"] = tran.Config.Tid
		result["ps_desc"] = resultReturnMoney["ps_desc"]
		result["error"] = "нет такого продукта в ассортименте"
		result["status"] = transaction.TransactionState_Error
		return
	}
	if configProduct.Id == 0 {
		logger.Log.Errorf("нет такого продукта в ассортименте")
		resultReturnMoney := payment.Payment.Return(tran)
		result["id"] = tran.Config.Tid
		result["ps_desc"] = resultReturnMoney["ps_desc"]
		result["error"] = "нет такого продукта в ассортименте"
		result["status"] = transaction.TransactionState_Error
		return
	}
	ware, errW := tran.Stores.StoreWare.GetOneById(parserTypes.ParseTypeInterfaceToInt(responseDevice["ware_id"]))
	if errW != nil {
		logger.Log.Errorf("нет такого товара в номенклатуре")
		resultReturnMoney := payment.Payment.Return(tran)
		result["id"] = tran.Config.Tid
		result["ps_desc"] = resultReturnMoney["ps_desc"]
		result["error"] = "нет такого товара в номенклатуре"
		result["status"] = transaction.TransactionState_Error
		return
	}
	if ware.Id == 0 {
		logger.Log.Errorf("нет такого товара в номенклатуре")
		resultReturnMoney := payment.Payment.Return(tran)
		result["id"] = tran.Config.Tid
		result["ps_desc"] = resultReturnMoney["ps_desc"]
		result["error"] = "нет такого товара в номенклатуре"
		result["status"] = transaction.TransactionState_Error
		return
	}
	logger.Log.Infof("%+v", ware)
	logger.Log.Infof("%+v", configProduct)
	product := &transaction.Product{
		Name:           ware.Name.String,
		Payment_device: "DA",
		Type:           ware.Type.Int32,
		Select_id:      configProduct.Select_id.String,
		Ware_id:        configProduct.Ware_id.Int32,
		Tax_rate:       configProduct.Tax_rate.Int32,
		Quantity:       int64(1000),
		Value:          parserTypes.ParseTypeInFloat64(responseDevice["sum"]),
		Price:          parserTypes.ParseTypeInFloat64(responseDevice["sum"]),
	}
	tran.Products = append(tran.Products, product)
	sapp.AddProductTransaction(tran)
	logger.Log.Infof("%+v", tran.Products)
	return
}

func (sapp *SaleAutomatPostPaid) AddProductTransaction(tran *transaction.Transaction) {
	transaction_dispetcher.Dispetcher.AddTransactionProduct(tran)
}

func (sapp *SaleAutomatPostPaid) returnRoutine(result *map[string]interface{}, tid int) {
	status := *result
	if r := recover(); r != nil {
		status["status"] = transaction.TransactionState_Error
		status["error"] = fmt.Sprintf("%v", r)
		logger.Log.Errorf("recovered from %v", r)
	}
	transaction_dispetcher.Dispetcher.StoreTransaction.SetByParams(*result)
	if parserTypes.ParseTypeInterfaceToInt(status["status"]) == transaction.TransactionState_Error {
		transaction_dispetcher.Dispetcher.RemoveTransaction(tid)
		transaction_dispetcher.Dispetcher.RemoveReplayProtection(sapp.KeyReplay)
	}
}
