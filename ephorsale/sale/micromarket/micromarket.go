package micromarket

import (
	config "ephorservices/config"
	"ephorservices/ephorsale/fiscal"
	"ephorservices/ephorsale/fiscal/interfaceFiscal"
	"ephorservices/ephorsale/payment"
	"ephorservices/ephorsale/payment/interfacePayment"
	"ephorservices/ephorsale/sale/interfaceSale"
	transaction "ephorservices/ephorsale/transaction/transaction_struct"
	parserTypes "ephorservices/pkg/parser/typeParse"
	"fmt"
	"log"
	"runtime"
)

type SaleCoolerPrePaidUnLock struct {
	cfg        *config.Config
	FiscalM    *fiscal.FiscalManager
	PaymentM   *payment.Manager
	Dispetcher *transaction.TransactionDispetcher
}

type NewSaleCoolerPrePaidUnLock struct {
	SaleCoolerPrePaidUnLock
}

func (newA *NewSaleCoolerPrePaidUnLock) New(conf *config.Config, fiscalM *fiscal.FiscalManager, paymentM *payment.Manager, dispether *transaction.TransactionDispetcher) interfaceSale.Sale {
	return &NewSaleCoolerPrePaidUnLock{
		SaleCoolerPrePaidUnLock: SaleCoolerPrePaidUnLock{
			cfg:        conf,
			FiscalM:    fiscalM,
			PaymentM:   paymentM,
			Dispetcher: dispether,
		},
	}
}

func (scpp *SaleCoolerPrePaidUnLock) Start(tran *transaction.Transaction) {
	keyReplayProtection := tran.Config.AutomatId + tran.Config.AccountId
	defer scpp.returnRoutine(tran, keyReplayProtection)
	resultDb := make(map[string]interface{})
	if len(tran.Products) < 1 {
		resultDb["id"] = tran.Config.Tid
		resultDb["ps_desc"] = "Не найдены товары для продолжения оплаты"
		resultDb["error"] = "Не найдены товары для продолжения оплаты"
		resultDb["status"] = transaction.TransactionState_Error
		scpp.Dispetcher.StoreTransaction.SetByParams(resultDb)
		scpp.Dispetcher.RemoveTransaction(tran.Config.Tid)
		return
	}
	resultPayment, PaymentSystem := scpp.Payment(tran)
	if resultPayment["status"] == transaction.TransactionState_Error {
		resultDb["id"] = tran.Config.Tid
		resultDb["ps_desc"] = resultPayment["ps_desc"]
		resultDb["error"] = "Ошибка"
		resultDb["status"] = resultPayment["status"]
		scpp.Dispetcher.StoreTransaction.SetByParams(resultDb)
		scpp.Dispetcher.RemoveTransaction(tran.Config.Tid)
		return
	}
	resultPayment["id"] = tran.Config.Tid
	scpp.Dispetcher.StoreTransaction.SetByParams(resultPayment)
	resultStatusHoldMoney := scpp.PaymentM.SatusHoldMoney(tran, PaymentSystem)
	if parserTypes.ParseTypeInterfaceToInt(resultStatusHoldMoney["status"]) != transaction.TransactionState_MoneyDebitOk {
		resultStatusHoldMoney["id"] = tran.Config.Tid
		resultStatusHoldMoney["error"] = "Ошибка оплаты"
		scpp.Dispetcher.StoreTransaction.SetByParams(resultStatusHoldMoney)
		scpp.Dispetcher.RemoveTransaction(tran.Config.Tid)
		return
	}
	resultStatusHoldMoney["id"] = tran.Config.Tid
	scpp.Dispetcher.StoreTransaction.SetByParams(resultPayment)
	resultStatusHoldMoney = scpp.PaymentM.StartDebitHoldMoney(tran, PaymentSystem)
	if parserTypes.ParseTypeInterfaceToInt(resultStatusHoldMoney["status"]) != transaction.TransactionState_MoneyDebitOk {
		resultStatusHoldMoney["id"] = tran.Config.Tid
		resultStatusHoldMoney["error"] = "Ошибка оплаты"
		scpp.Dispetcher.StoreTransaction.SetByParams(resultStatusHoldMoney)
		scpp.Dispetcher.RemoveTransaction(tran.Config.Tid)
		return
	}
	resultStatusHoldMoney["id"] = tran.Config.Tid
	scpp.Dispetcher.StoreTransaction.SetByParams(resultPayment)
	resultDb["id"] = tran.Config.Tid
	resultDb["ps_desc"] = tran.GetDescriptionCodeCooler(transaction.TransactionState_WaitFiscal)
	resultDb["error"] = "Нет"
	resultDb["status"] = tran.GetStatusServer(transaction.TransactionState_WaitFiscal)
	scpp.Dispetcher.StoreTransaction.SetByParams(resultDb)

	resultFiscal := scpp.Fiscal(tran)
	if resultFiscal["status"].(uint8) != interfaceFiscal.Status_Complete {
		resultDb["id"] = tran.Config.Tid
		resultDb["f_desc"] = resultFiscal["error"]
		resultDb["f_status"] = resultFiscal["status"]
		scpp.Dispetcher.StoreTransaction.SetByParams(resultDb)
	} else {
		resultDb["id"] = tran.Config.Tid
		resultDb["f_desc"] = "Ошибок нет"
		resultDb["f_status"] = resultFiscal["status"]
		resultDb["f_qr"] = resultFiscal["f_qr"]
		resultDb["fn"] = tran.Fiscal.Fields.Fn
		resultDb["fd"] = tran.Fiscal.Fields.Fd
		resultDb["fp"] = tran.Fiscal.Fields.Fp
		resultDb["f_type"] = tran.Fiscal.Config.Type
		resultDb["f_receipt"] = tran.Fiscal.ResiptId
		scpp.Dispetcher.StoreTransaction.SetByParams(resultDb)
	}
	tran.Fiscal.Status = int(resultFiscal["status"].(uint8))
	err := scpp.Dispetcher.AddSales(tran)
	if err != nil {
		log.Printf("%v", err)
	}
	scpp.Dispetcher.RemoveTransaction(tran.Config.Tid)
	scpp.Dispetcher.RemoveReplayProtection(keyReplayProtection)
	return

}

func (scpp *SaleCoolerPrePaidUnLock) Payment(tran *transaction.Transaction) (map[string]interface{}, interfacePayment.Payment) {
	resultDb := make(map[string]interface{})
	resultPayment, PaymentSystem := scpp.PaymentM.StartPayment(tran)
	if resultPayment["status"] != transaction.TransactionState_MoneyHoldStart {
		return resultPayment, nil
	}
	resultDb["id"] = tran.Config.Tid
	resultDb["status"] = transaction.TransactionState_MoneyHoldWait
	scpp.Dispetcher.StoreTransaction.SetByParams(resultDb)
	resultPay := scpp.PaymentM.HoldMoney(tran, PaymentSystem)
	if resultPay["status"] == transaction.TransactionState_Error {
		return resultPay, nil
	}
	return resultPay, PaymentSystem
}

func (scpp *SaleCoolerPrePaidUnLock) SendMassage(tran *transaction.Transaction, PaymentSystem interfacePayment.Payment) (result map[string]interface{}) {
	return result
}

func (scpp *SaleCoolerPrePaidUnLock) WaitMassage(tran *transaction.Transaction, PaymentSystem interfacePayment.Payment) (result map[string]interface{}) {
	return result
}

func (scpp *SaleCoolerPrePaidUnLock) Fiscal(tran *transaction.Transaction) map[string]interface{} {
	resultFiscal := make(map[string]interface{})
	automat, err := scpp.Dispetcher.StoreAutomat.GetOneById(tran.Config.AutomatId, tran.Config.AccountId)
	if err != nil {
		log.Printf("%v", err)
		tran.Fiscal.Message = "невозможно взять автомат"
		resultFiscal["status"] = interfaceFiscal.Status_Error
		resultFiscal["error"] = "невозможно взять автомат"
		return resultFiscal
	}
	if automat.Id == 0 {
		log.Printf("%+v", automat)
		tran.Fiscal.Message = "невозможно взять автомат"
		resultFiscal["status"] = interfaceFiscal.Status_Error
		resultFiscal["error"] = "невозможно взять автомат"
		return resultFiscal
	}
	options := make(map[string]interface{})
	options["account_id"] = tran.Config.AccountId
	options["automat_id"] = automat.Id
	options["to_date"] = nil
	locationAutomat, errLoc := scpp.Dispetcher.StoreLocation.GetOneWithOptions(options)
	if errLoc != nil {
		log.Printf("%v", errLoc)
		resultFiscal["status"] = interfaceFiscal.Status_Error
		resultFiscal["error"] = fmt.Sprintf("%v", errLoc)
		tran.Fiscal.Message = string(resultFiscal["error"].(string))
		return resultFiscal
	}
	if locationAutomat.Id == 0 {
		log.Printf("%+v", locationAutomat)
		resultFiscal["status"] = interfaceFiscal.Status_Error
		resultFiscal["error"] = "автомат не стоит на точке"
		tran.Fiscal.Message = string(resultFiscal["error"].(string))
		return resultFiscal
	}
	tran.Fiscal.Config.Id = int64(automat.Fr_id)
	tran.Point.Id = parserTypes.ParseTypeInterfaceToInt(locationAutomat.Company_point_id)
	entryPoint, errPoint := scpp.Dispetcher.StorePoint.GetOneById(locationAutomat.Company_point_id, tran.Config.AccountId)
	if errPoint != nil {
		log.Printf("%v", errPoint)
		resultFiscal["status"] = interfaceFiscal.Status_Error
		resultFiscal["error"] = fmt.Sprintf("%v", errPoint)
		tran.Fiscal.Message = string(resultFiscal["error"].(string))
		return resultFiscal
	}
	if entryPoint.Id == 0 {
		log.Printf("%+v", entryPoint)
		resultFiscal["status"] = interfaceFiscal.Status_Error
		resultFiscal["error"] = "нет точки"
		tran.Fiscal.Message = string(resultFiscal["error"].(string))
		return resultFiscal
	}
	tran.Point.Address = entryPoint.Address
	tran.Point.PointName = entryPoint.Name
	resultStartFiscal, frKass := scpp.FiscalM.DataVerification(tran)
	if !resultStartFiscal["status"].(bool) {
		resultFiscal["status"] = interfaceFiscal.Status_Error
		resultFiscal["error"] = resultStartFiscal["message"]
		tran.Fiscal.Message = string(resultFiscal["error"].(string))
		return resultFiscal
	}
	statusSendCheck := frKass.SendCheck(tran)
	if statusSendCheck["status"].(uint8) != interfaceFiscal.Status_InQueue && statusSendCheck["status"].(uint8) != interfaceFiscal.Status_Retry_Status {
		resultFiscal["status"] = interfaceFiscal.Status_Error
		resultFiscal["error"] = resultStartFiscal["f_desc"]
		tran.Fiscal.Message = string(resultFiscal["error"].(string))
		return resultFiscal
	}
	resultStatusCheck := frKass.GetStatus(tran)
	if resultStatusCheck["status"].(uint8) != interfaceFiscal.Status_Complete {
		resultFiscal["status"] = interfaceFiscal.Status_Error
		resultFiscal["error"] = resultStartFiscal["f_desc"]
		tran.Fiscal.Message = string(resultFiscal["error"].(string))
		return resultFiscal
	}
	resultFiscal["status"] = interfaceFiscal.Status_Complete
	resultFiscal["f_qr"] = frKass.GetQrUrl(tran)
	resultFiscal["error"] = "Нет ошибок"
	tran.Fiscal.Message = string(resultFiscal["error"].(string))
	return resultFiscal
}

func (scpp *SaleCoolerPrePaidUnLock) returnRoutine(tran *transaction.Transaction, keyReplayProtection int) {
	if r := recover(); r != nil {
		log.Printf("recovered from %v", r)
		err := scpp.Dispetcher.AddSales(tran)
		if err != nil {
			log.Printf("%v", err)
		}
		scpp.Dispetcher.RemoveTransaction(tran.Config.Tid)
		scpp.Dispetcher.RemoveReplayProtection(keyReplayProtection)
	}
	runtime.Goexit()
}
