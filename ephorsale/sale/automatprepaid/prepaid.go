package automatprepaid

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

type SaleAutomatPrePaid struct {
	Debug              bool
	KeyReplay          int
	ExecuteTimeSeconds int
	Fiscalization      interfaceSale.FiscalFunc
}

type NewSaleAutomatPrePaid struct {
	SaleAutomatPrePaid
}

func (newA NewSaleAutomatPrePaid) New(executeTimeSeconds int, debug bool) interfaceSale.Sale {
	return &NewSaleAutomatPrePaid{
		SaleAutomatPrePaid: SaleAutomatPrePaid{
			ExecuteTimeSeconds: executeTimeSeconds,
			Debug:              debug,
		},
	}
}

func (sapp *SaleAutomatPrePaid) SetFiscalisation(function interfaceSale.FiscalFunc) {
	sapp.Fiscalization = function
}

// создание goroutine
func (sapp *SaleAutomatPrePaid) Sale(tran *transaction.Transaction) (result map[string]interface{}) {
	defer sapp.returnRoutine(&result, tran.Config.Tid)
	if len(tran.Products) < 1 {
		result["id"] = tran.Config.Tid
		result["ps_desc"] = "Не найдены товары для продолжения оплаты"
		result["error"] = "Не найдены товары для продолжения оплаты"
		result["status"] = transaction.TransactionState_Error
		return
	}
	sapp.KeyReplay = tran.Config.AutomatId + tran.Config.AccountId
	tran.KeyReplay = sapp.KeyReplay
	resultPayment := sapp.Payment(tran)
	if resultPayment["status"] == transaction.TransactionState_Error {
		result["id"] = tran.Config.Tid
		result["ps_desc"] = resultPayment["ps_desc"]
		result["error"] = "Ошибка"
		result["status"] = resultPayment["status"]
		return
	}
	resultPayment["id"] = tran.Config.Tid
	transaction_dispetcher.Dispetcher.StoreTransaction.SetByParams(resultPayment)
	result = payment.Payment.Satus(tran)
	if parserTypes.ParseTypeInterfaceToInt(result["status"]) != transaction.TransactionState_MoneyDebitWait {
		result["id"] = tran.Config.Tid
		result["status"] = transaction.TransactionState_Error
		result["error"] = "Ошибка оплаты"
		return
	}
	result["id"] = tran.Config.Tid
	transaction_dispetcher.Dispetcher.StoreTransaction.SetByParams(result)
	result = payment.Payment.Debit(tran)
	if parserTypes.ParseTypeInterfaceToInt(result["status"]) != transaction.TransactionState_MoneyDebitOk {
		result["id"] = tran.Config.Tid
		result["status"] = transaction.TransactionState_Error
		result["error"] = "Ошибка оплаты"
		return
	}
	result["id"] = tran.Config.Tid
	transaction_dispetcher.Dispetcher.StoreTransaction.SetByParams(result)
	transaction_dispetcher.Dispetcher.AddReplayProtection(sapp.KeyReplay, tran.Config.AutomatId)
	result = sapp.SendMassage(tran)
	if parserTypes.ParseTypeInterfaceToInt(result["status"]) != transaction.TransactionState_MoneyHoldWait {
		result["id"] = tran.Config.Tid
		result["status"] = transaction.TransactionState_Error
		return
	}
	responseRabbit := sapp.WaitMassage(tran)
	if parserTypes.ParseTypeInterfaceToInt(responseRabbit["status"]) == transaction.TransactionState_ReturnMoney {
		result["id"] = tran.Config.Tid
		result["ps_desc"] = responseRabbit["ps_desc"]
		result["error"] = responseRabbit["error"]
		result["status"] = transaction.TransactionState_Error
		return
	}
	if parserTypes.ParseTypeInterfaceToInt(responseRabbit["status"]) != transaction.VendState_Session {
		resultReturnMoney := payment.Payment.Return(tran)
		result["id"] = tran.Config.Tid
		result["ps_desc"] = resultReturnMoney["description"]
		result["error"] = tran.GetDescriptionErr(parserTypes.ParseTypeInterfaceToInt(responseRabbit["error"]))
		result["status"] = transaction.TransactionState_Error
		return
	}
	result["id"] = tran.Config.Tid
	result["ps_desc"] = tran.GetDescriptionCode(parserTypes.ParseTypeInterfaceToInt(responseRabbit["status"]))
	result["error"] = tran.GetDescriptionErr(parserTypes.ParseTypeInterfaceToInt(responseRabbit["error"]))
	result["status"] = tran.GetStatusServer(transaction.VendState_Session)
	transaction_dispetcher.Dispetcher.StoreTransaction.SetByParams(result)
	responseRabbit = sapp.WaitMassage(tran)
	if parserTypes.ParseTypeInterfaceToInt(responseRabbit["status"]) == transaction.TransactionState_ReturnMoney {
		result["id"] = tran.Config.Tid
		result["ps_desc"] = responseRabbit["ps_desc"]
		result["error"] = responseRabbit["error"]
		result["status"] = transaction.TransactionState_Error
		return
	}
	if parserTypes.ParseTypeInterfaceToInt(responseRabbit["status"]) != transaction.VendState_Vending {
		resultReturnMoney := payment.Payment.Return(tran)
		result["id"] = tran.Config.Tid
		result["ps_desc"] = resultReturnMoney["description"]
		result["error"] = tran.GetDescriptionErr(parserTypes.ParseTypeInterfaceToInt(responseRabbit["error"]))
		result["status"] = transaction.TransactionState_Error
		return
	}
	result["id"] = tran.Config.Tid
	result["ps_desc"] = tran.GetDescriptionCode(parserTypes.ParseTypeInterfaceToInt(responseRabbit["status"]))
	result["error"] = tran.GetDescriptionErr(parserTypes.ParseTypeInterfaceToInt(responseRabbit["error"]))
	result["status"] = tran.GetStatusServer(transaction.VendState_Session)
	transaction_dispetcher.Dispetcher.StoreTransaction.SetByParams(result)
	responseRabbit = sapp.WaitMassage(tran)
	if parserTypes.ParseTypeInterfaceToInt(responseRabbit["status"]) == transaction.TransactionState_ReturnMoney {
		result["id"] = tran.Config.Tid
		result["error"] = tran.GetDescriptionErr(parserTypes.ParseTypeInterfaceToInt(responseRabbit["error"]))
		result["status"] = transaction.TransactionState_Error
		return
	}
	if parserTypes.ParseTypeInterfaceToInt(responseRabbit["status"]) != transaction.VendState_VendOk {
		resultReturnMoney := payment.Payment.Return(tran)
		result["id"] = tran.Config.Tid
		result["ps_desc"] = resultReturnMoney["description"]
		result["error"] = tran.GetDescriptionErr(parserTypes.ParseTypeInterfaceToInt(responseRabbit["error"]))
		result["status"] = transaction.TransactionState_Error
		return
	}
	result["id"] = tran.Config.Tid
	result["ps_desc"] = tran.GetDescriptionCode(parserTypes.ParseTypeInterfaceToInt(responseRabbit["status"]))
	result["error"] = tran.GetDescriptionErr(parserTypes.ParseTypeInterfaceToInt(responseRabbit["error"]))
	result["status"] = tran.GetStatusServer(transaction.VendState_VendOk)
	return
}

func (sapp *SaleAutomatPrePaid) Payment(tran *transaction.Transaction) map[string]interface{} {
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

func (sapp *SaleAutomatPrePaid) SendMassage(tran *transaction.Transaction) map[string]interface{} {
	result := make(map[string]interface{})
	resultMessage := make(map[string]interface{})
	resultMessage["tid"] = tran.Config.Tid
	resultMessage["sum"] = tran.Payment.Sum
	resultMessage["wid"] = tran.Products[0].Ware_id
	resultMessage["m"] = 1
	resultMessage["a"] = 1
	err := transport_manager.TransportManager.QueueManager.SendMessage(resultMessage, tran.Config.Imei)
	if err != nil {
		return payment.Payment.Return(tran)
	}
	result["status"] = transaction.TransactionState_MoneyHoldWait
	return result
}

func (sapp *SaleAutomatPrePaid) WaitMassage(tran *transaction.Transaction) map[string]interface{} {
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

func (sapp *SaleAutomatPrePaid) Fiscal(tran *transaction.Transaction) map[string]interface{} {
	return make(map[string]interface{})
}

func (sapp *SaleAutomatPrePaid) returnRoutine(result *map[string]interface{}, tid int) {
	status := *result
	if r := recover(); r != nil {
		status["status"] = transaction.TransactionState_Error
		status["error"] = fmt.Sprintf("%v", r)
		logger.Log.Errorf("%+v", status)
	}
	transaction_dispetcher.Dispetcher.StoreTransaction.SetByParams(*result)
	if parserTypes.ParseTypeInterfaceToInt(status["status"]) == transaction.TransactionState_Error {
		transaction_dispetcher.Dispetcher.RemoveTransaction(tid)
		transaction_dispetcher.Dispetcher.RemoveReplayProtection(sapp.KeyReplay)
	}
}
