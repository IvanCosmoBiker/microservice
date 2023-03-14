package coolerprepaid

import (
	"ephorservices/ephorsale/fiscal/interface/fr"
	payment "ephorservices/ephorsale/payment"
	"ephorservices/ephorsale/sale/interfaceSale"
	transaction_dispetcher "ephorservices/ephorsale/transaction"
	transaction "ephorservices/ephorsale/transaction/transaction_struct"
	transport_manager "ephorservices/ephorsale/transport"
	logger "ephorservices/pkg/logger"
	parserTypes "ephorservices/pkg/parser/typeParse"
	"fmt"
)

type SaleCoolerPrePaid struct {
	Debug              bool
	KeyReplay          int
	ExecuteTimeSeconds int
	Fiscalization      interfaceSale.FiscalFunc
}

type NewSaleCoolerPrePaid struct {
	SaleCoolerPrePaid
}

func (newA *NewSaleCoolerPrePaid) New(executeTimeSeconds int, debug bool) interfaceSale.Sale {
	return &NewSaleCoolerPrePaid{
		SaleCoolerPrePaid: SaleCoolerPrePaid{
			ExecuteTimeSeconds: executeTimeSeconds,
			Debug:              debug,
		},
	}
}

func (scpp *SaleCoolerPrePaid) SetFiscalisation(function interfaceSale.FiscalFunc) {
	scpp.Fiscalization = function
}

func (scpp *SaleCoolerPrePaid) Sale(tran *transaction.Transaction) (result map[string]interface{}) {
	defer scpp.returnRoutine(&result, scpp.KeyReplay)
	if len(tran.Products) < 1 {
		result["id"] = tran.Config.Tid
		result["ps_desc"] = "Не найдены товары для продолжения оплаты"
		result["error"] = "Не найдены товары для продолжения оплаты"
		result["status"] = transaction.TransactionState_Error
		return
	}
	scpp.KeyReplay = tran.Config.AutomatId + tran.Config.AccountId
	tran.KeyReplay = scpp.KeyReplay
	result = scpp.Payment(tran)
	if result["status"] == transaction.TransactionState_Error {
		result["id"] = tran.Config.Tid
		result["status"] = transaction.TransactionState_Error
		result["error"] = "Ошибка"
		return
	}
	result["id"] = tran.Config.Tid
	transaction_dispetcher.Dispetcher.StoreTransaction.SetByParams(result)
	result = scpp.clearMap(result)
	result = payment.Payment.Satus(tran)
	if parserTypes.ParseTypeInterfaceToInt(result["status"]) != transaction.TransactionState_MoneyDebitWait {
		result["id"] = tran.Config.Tid
		result["status"] = transaction.TransactionState_Error
		result["error"] = "Ошибка оплаты"
		return
	}
	tran.Payment.DebitSum = tran.Payment.Sum
	result["id"] = tran.Config.Tid
	transaction_dispetcher.Dispetcher.StoreTransaction.SetByParams(result)
	result = scpp.clearMap(result)
	result = payment.Payment.Debit(tran)
	if parserTypes.ParseTypeInterfaceToInt(result["status"]) != transaction.TransactionState_MoneyDebitOk {
		result["id"] = tran.Config.Tid
		result["status"] = transaction.TransactionState_Error
		result["error"] = "Ошибка оплаты"
		return
	}
	result["id"] = tran.Config.Tid
	transaction_dispetcher.Dispetcher.StoreTransaction.SetByParams(result)
	transaction_dispetcher.Dispetcher.AddReplayProtection(scpp.KeyReplay, tran.Config.AutomatId)
	result = scpp.clearMap(result)
	result = scpp.SendMassage(tran)
	if parserTypes.ParseTypeInterfaceToInt(result["status"]) != transaction.TransactionState_MoneyHoldWait {
		result["id"] = tran.Config.Tid
		result["status"] = transaction.TransactionState_Error
		return
	}
	result = scpp.clearMap(result)
	responseRabbit := scpp.WaitMassage(tran)
	if parserTypes.ParseTypeInterfaceToInt(responseRabbit["status"]) == transaction.TransactionState_ReturnMoney {
		result["id"] = tran.Config.Tid
		result["ps_desc"] = responseRabbit["ps_desc"]
		result["error"] = responseRabbit["error"]
		result["status"] = transaction.TransactionState_Error
		return
	}
	if parserTypes.ParseTypeInterfaceToInt(responseRabbit["status"]) != transaction.IceboxStatus_Drink && parserTypes.ParseTypeInterfaceToInt(responseRabbit["status"]) != transaction.IceboxStatus_Icebox {
		resultReturnMoney := payment.Payment.Return(tran)
		result["id"] = tran.Config.Tid
		result["ps_desc"] = resultReturnMoney["ps_desc"]
		result["error"] = tran.GetDescriptionErr(parserTypes.ParseTypeInterfaceToInt(responseRabbit["error"]))
		result["status"] = transaction.TransactionState_Error
		resultReturnMoney = nil
		return
	}
	result["id"] = tran.Config.Tid
	result["ps_desc"] = tran.GetDescriptionCodeCooler(parserTypes.ParseTypeInterfaceToInt(responseRabbit["status"]))
	result["error"] = tran.GetDescriptionErr(parserTypes.ParseTypeInterfaceToInt(responseRabbit["error"]))
	result["status"] = tran.GetStatusServer(transaction.VendState_Session)
	transaction_dispetcher.Dispetcher.StoreTransaction.SetByParams(result)
	result = scpp.clearMap(result)
	responseRabbit = scpp.WaitMassage(tran)
	if parserTypes.ParseTypeInterfaceToInt(responseRabbit["status"]) == transaction.TransactionState_ReturnMoney {
		result["id"] = tran.Config.Tid
		result["ps_desc"] = responseRabbit["ps_desc"]
		result["error"] = responseRabbit["error"]
		result["status"] = transaction.TransactionState_Error
		return
	}
	if parserTypes.ParseTypeInterfaceToInt(responseRabbit["status"]) != transaction.IceboxStatus_End {
		resultReturnMoney := payment.Payment.Return(tran)
		result["id"] = tran.Config.Tid
		result["ps_desc"] = resultReturnMoney["ps_desc"]
		result["error"] = tran.GetDescriptionErr(parserTypes.ParseTypeInterfaceToInt(responseRabbit["error"]))
		result["status"] = transaction.TransactionState_Error
		resultReturnMoney = nil
		return
	}
	result["id"] = tran.Config.Tid
	result["ps_desc"] = tran.GetDescriptionCodeCooler(parserTypes.ParseTypeInterfaceToInt(responseRabbit["status"]))
	result["error"] = tran.GetDescriptionErr(parserTypes.ParseTypeInterfaceToInt(responseRabbit["error"]))
	result["status"] = tran.GetStatusServer(transaction.VendState_Session)
	transaction_dispetcher.Dispetcher.StoreTransaction.SetByParams(result)
	result["id"] = tran.Config.Tid
	result["error"] = "Нет"
	result["ps_desc"] = "Подождите. Получаем Ваш чек"
	result["status"] = transaction.TransactionState_WaitFiscal
	transaction_dispetcher.Dispetcher.StoreTransaction.SetByParams(result)
	resaltFiscal := scpp.Fiscal(tran)
	transaction_dispetcher.Dispetcher.StoreTransaction.SetByParams(resaltFiscal)
	resaltFiscal = nil
	result["id"] = tran.Config.Tid
	result["ps_desc"] = tran.GetDescriptionCodeCooler(parserTypes.ParseTypeInterfaceToInt(responseRabbit["status"]))
	result["error"] = tran.GetDescriptionErr(parserTypes.ParseTypeInterfaceToInt(responseRabbit["error"]))
	result["status"] = tran.GetStatusServerCooler(transaction.IceboxStatus_End)
	transaction_dispetcher.Dispetcher.StoreTransaction.SetByParams(result)
	return
}

func (scpp *SaleCoolerPrePaid) Payment(tran *transaction.Transaction) map[string]interface{} {
	result := make(map[string]interface{})
	result = payment.Payment.Payment(tran)
	if parserTypes.ParseTypeInterfaceToInt(result["status"]) != transaction.TransactionState_MoneyHoldStart {
		return result
	}
	result["id"] = tran.Config.Tid
	result["status"] = transaction.TransactionState_MoneyHoldWait
	transaction_dispetcher.Dispetcher.StoreTransaction.SetByParams(result)
	result = scpp.clearMap(result)
	result = payment.Payment.Hold(tran)
	if parserTypes.ParseTypeInterfaceToInt(result["status"]) == transaction.TransactionState_Error {
		return result
	}
	return result
}

func (scpp *SaleCoolerPrePaid) SendMassage(tran *transaction.Transaction) map[string]interface{} {
	result := make(map[string]interface{})
	resultMessage := make(map[string]interface{})
	resultMessage["tid"] = tran.Config.Tid
	resultMessage["sum"] = tran.Payment.Sum
	resultMessage["wid"] = tran.Products[0].Ware_id
	resultMessage["m"] = 4
	resultMessage["a"] = 1
	err := transport_manager.TransportManager.QueueManager.SendMessage(resultMessage, tran.Config.Imei)
	resultMessage = nil
	if err != nil {
		return payment.Payment.Return(tran)
	}
	result["status"] = transaction.TransactionState_MoneyHoldWait
	return result
}

func (scpp *SaleCoolerPrePaid) WaitMassage(tran *transaction.Transaction) map[string]interface{} {
	response := make(map[string]interface{})
	result, err := transport_manager.TransportManager.QueueManager.WaitMessage(tran.ChannelMessage, scpp.ExecuteTimeSeconds)
	if err != nil {
		request := payment.Payment.Return(tran)
		request["status"] = transaction.TransactionState_ReturnMoney
		request["error"] = err.Error()
		result, err = nil, nil
		return request
	}
	response["status"] = result["st"]
	response["error"] = result["err"]
	response["tid"] = result["tid"]
	response["ware_id"] = result["wid"]
	response["sum"] = result["sum"]
	response["a"] = result["a"]
	response["d"] = result["d"]
	result = nil
	return response
}

func (scpp *SaleCoolerPrePaid) Fiscal(tran *transaction.Transaction) map[string]interface{} {
	result := make(map[string]interface{})
	tran.Fiscal.NeedFiscal = true
	scpp.GetTypeProduct(tran)
	scpp.Fiscalization(tran)
	if parserTypes.ParseTypeInterfaceToUint8(tran.Status) != fr.Status_Complete {
		result["id"] = tran.Config.Tid
		result["f_desc"] = tran.Error
		result["f_status"] = tran.Status
		return result
	} else {
		result["id"] = tran.Config.Tid
		result["f_desc"] = "Ошибок нет"
		result["f_status"] = tran.Status
		result["f_qr"] = tran.Fiscal.QrCode
		result["fn"] = tran.Fiscal.Fields.Fn
		result["fd"] = tran.Fiscal.Fields.Fd
		result["fp"] = tran.Fiscal.Fields.Fp
		result["f_type"] = tran.Fiscal.Config.Type
		result["f_receipt"] = tran.Fiscal.ResiptId
		return result
	}
	tran.Fiscal.Status = int(parserTypes.ParseTypeInterfaceToUint8(tran.Status))
	return result
}

func (scpp *SaleCoolerPrePaid) GetTypeProduct(tran *transaction.Transaction) {
	for _, product := range tran.Products {
		ware, _ := tran.Stores.StoreWare.GetOneById(int(product.Ware_id))
		if ware.Id != 0 {
			product.Type = ware.Type.Int32
		}
	}
}

func (scpp *SaleCoolerPrePaid) clearMap(Map map[string]interface{}) map[string]interface{} {
	Map = nil
	Map = make(map[string]interface{})
	return Map
}

func (scpp *SaleCoolerPrePaid) returnRoutine(result *map[string]interface{}, tid int) {
	status := *result
	if r := recover(); r != nil {
		status["status"] = transaction.TransactionState_Error
		status["error"] = fmt.Sprintf("%v", r)
		logger.Log.Errorf("%+v", status)
	}
	transaction_dispetcher.Dispetcher.StoreTransaction.SetByParams(*result)
	if parserTypes.ParseTypeInterfaceToInt(status["status"]) == transaction.TransactionState_Error {
		transaction_dispetcher.Dispetcher.RemoveTransaction(tid)
		transaction_dispetcher.Dispetcher.RemoveReplayProtection(scpp.KeyReplay)
	}
	status = nil
}
