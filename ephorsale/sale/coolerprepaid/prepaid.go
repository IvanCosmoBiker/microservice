package coolerprepaid

import (
	"encoding/json"
	config "ephorservices/config"
	"ephorservices/ephorsale/fiscal"
	"ephorservices/ephorsale/payment"
	"ephorservices/ephorsale/payment/interfacePayment"
	"ephorservices/ephorsale/sale/interfaceSale"
	responseQueueManager "ephorservices/ephorsale/sale/responseQueueManager"
	transaction "ephorservices/ephorsale/transaction"
	parserTypes "ephorservices/pkg/parser/typeParse"
	rabbit "ephorservices/pkg/rabbitmq"
	"fmt"
	"log"
	"runtime"
	"time"
)

type SaleCoolerPrePaid struct {
	cfg          *config.Config
	QueueManager *rabbit.Manager
	FiscalM      *fiscal.FiscalManager
	PaymentM     *payment.PaymentManager
	Dispetcher   *transaction.TransactionDispetcher
}

type NewSaleCoolerPrePaid struct {
	SaleCoolerPrePaid
}

func (newA *NewSaleCoolerPrePaid) New(conf *config.Config, rabbitMq *rabbit.Manager, fiscalM *fiscal.FiscalManager, paymentM *payment.PaymentManager, dispether *transaction.TransactionDispetcher) interfaceSale.Sale {
	return &NewSaleCoolerPrePaid{
		SaleCoolerPrePaid: SaleCoolerPrePaid{
			cfg:          conf,
			QueueManager: rabbitMq,
			FiscalM:      fiscalM,
			PaymentM:     paymentM,
			Dispetcher:   dispether,
		},
	}
}

func (scpp *SaleCoolerPrePaid) Start(tran *transaction.Transaction) {
	resultDb := make(map[string]interface{})
	if len(tran.Products) < 1 {
		resultDb["id"] = tran.Config.Tid
		resultDb["ps_desc"] = "Не найдены товары для продолжения оплаты"
		resultDb["error"] = "Не найдены товары для продолжения оплаты"
		resultDb["status"] = transaction.TransactionState_Error
		scpp.Dispetcher.StoreTransaction.SetByParams(resultDb)
		scpp.Dispetcher.RemoveTransaction(tran.Config.Tid)
		scpp.returnRoutine()
	}
	keyReplayProtection := tran.Config.AutomatId + tran.Config.AccountId
	resultReturnMoney := make(map[string]interface{})
	resultPayment, PaymentSystem := scpp.Payment(tran)
	if resultPayment["status"] == transaction.TransactionState_Error {
		resultDb["id"] = tran.Config.Tid
		resultDb["ps_desc"] = resultPayment["ps_desc"]
		resultDb["error"] = "Ошибка"
		resultDb["status"] = resultPayment["status"]
		scpp.Dispetcher.StoreTransaction.SetByParams(resultDb)
		scpp.Dispetcher.RemoveTransaction(tran.Config.Tid)
		scpp.returnRoutine()
	}
	resultPayment["id"] = tran.Config.Tid
	scpp.Dispetcher.StoreTransaction.SetByParams(resultPayment)
	scpp.Dispetcher.AddReplayProtection(keyReplayProtection, tran.Config.AutomatId)
	resultMassage := scpp.SendMassage(tran, PaymentSystem)
	if parserTypes.ParseTypeInterfaceToInt(resultMassage["status"]) != transaction.TransactionState_MoneyHoldWait {
		resultMassage["id"] = tran.Config.Tid
		scpp.Dispetcher.StoreTransaction.SetByParams(resultMassage)
		scpp.Dispetcher.RemoveTransaction(tran.Config.Tid)
		scpp.Dispetcher.RemoveReplayProtection(keyReplayProtection)
		scpp.returnRoutine()
	}
	responseRabbit := scpp.WaitMassage(tran, PaymentSystem)
	if parserTypes.ParseTypeInterfaceToInt(responseRabbit["status"]) == transaction.TransactionState_ReturnMoney {
		resultDb["id"] = tran.Config.Tid
		resultDb["ps_desc"] = responseRabbit["ps_desc"]
		resultDb["error"] = responseRabbit["error"]
		resultDb["status"] = transaction.TransactionState_Error
		scpp.Dispetcher.StoreTransaction.SetByParams(resultDb)
		scpp.Dispetcher.RemoveTransaction(tran.Config.Tid)
		scpp.Dispetcher.RemoveReplayProtection(keyReplayProtection)
		scpp.returnRoutine()
	}
	if parserTypes.ParseTypeInterfaceToInt(responseRabbit["status"]) != transaction.IceboxStatus_Drink && parserTypes.ParseTypeInterfaceToInt(responseRabbit["status"]) != transaction.IceboxStatus_Icebox {
		resultReturnMoney = scpp.returnMoney(tran, PaymentSystem)
		resultDb["id"] = tran.Config.Tid
		resultDb["ps_desc"] = resultReturnMoney["ps_desc"]
		resultDb["error"] = tran.GetDescriptionErr(parserTypes.ParseTypeInterfaceToInt(responseRabbit["error"]))
		resultDb["status"] = transaction.TransactionState_Error
		scpp.Dispetcher.StoreTransaction.SetByParams(resultDb)
		scpp.Dispetcher.RemoveTransaction(tran.Config.Tid)
		scpp.Dispetcher.RemoveReplayProtection(keyReplayProtection)
		scpp.returnRoutine()
	}
	resultDb["id"] = tran.Config.Tid
	resultDb["ps_desc"] = tran.GetDescriptionCodeCooler(parserTypes.ParseTypeInterfaceToInt(responseRabbit["status"]))
	resultDb["error"] = tran.GetDescriptionErr(parserTypes.ParseTypeInterfaceToInt(responseRabbit["error"]))
	resultDb["status"] = tran.GetStatusServer(transaction.VendState_Session)
	scpp.Dispetcher.StoreTransaction.SetByParams(resultDb)
	responseRabbit = scpp.WaitMassage(tran, PaymentSystem)
	if parserTypes.ParseTypeInterfaceToInt(responseRabbit["status"]) == transaction.TransactionState_ReturnMoney {
		resultDb["id"] = tran.Config.Tid
		resultDb["ps_desc"] = responseRabbit["ps_desc"]
		resultDb["error"] = responseRabbit["error"]
		resultDb["status"] = transaction.TransactionState_Error
		scpp.Dispetcher.StoreTransaction.SetByParams(resultDb)
		scpp.Dispetcher.RemoveTransaction(tran.Config.Tid)
		scpp.Dispetcher.RemoveReplayProtection(keyReplayProtection)
		scpp.returnRoutine()
	}
	if parserTypes.ParseTypeInterfaceToInt(responseRabbit["status"]) != transaction.IceboxStatus_End {
		resultReturnMoney = scpp.returnMoney(tran, PaymentSystem)
		resultDb["id"] = tran.Config.Tid
		resultDb["ps_desc"] = resultReturnMoney["ps_desc"]
		resultDb["error"] = tran.GetDescriptionErr(parserTypes.ParseTypeInterfaceToInt(responseRabbit["error"]))
		resultDb["status"] = transaction.TransactionState_Error
		scpp.Dispetcher.StoreTransaction.SetByParams(resultDb)
		scpp.Dispetcher.RemoveTransaction(tran.Config.Tid)
		scpp.Dispetcher.RemoveReplayProtection(keyReplayProtection)
		scpp.returnRoutine()
	}
	resultDb["id"] = tran.Config.Tid
	resultDb["ps_desc"] = tran.GetDescriptionCodeCooler(parserTypes.ParseTypeInterfaceToInt(responseRabbit["status"]))
	resultDb["error"] = tran.GetDescriptionErr(parserTypes.ParseTypeInterfaceToInt(responseRabbit["error"]))
	resultDb["status"] = tran.GetStatusServer(transaction.VendState_Session)
	scpp.Dispetcher.StoreTransaction.SetByParams(resultDb)
	//resultFiscal := scpp.Fiscal(tran)
	scpp.Dispetcher.RemoveTransaction(tran.Config.Tid)
	scpp.Dispetcher.RemoveReplayProtection(keyReplayProtection)
	scpp.returnRoutine()

}

func (scpp *SaleCoolerPrePaid) Payment(tran *transaction.Transaction) (map[string]interface{}, interfacePayment.Payment) {
	resultDb := make(map[string]interface{})
	resultPayment, PaymentSystem := scpp.PaymentM.StartPayment(tran)
	if resultPayment["status"] != transaction.TransactionState_MoneyHoldStart {
		return resultPayment, nil
	}
	resultDb["id"] = tran.Config.Tid
	resultDb["status"] = transaction.TransactionState_MoneyHoldWait
	scpp.Dispetcher.StoreTransaction.SetByParams(resultDb)
	resultPay := scpp.PaymentM.StartPrepayment(tran, PaymentSystem)
	if resultPay["status"] == transaction.TransactionState_Error {
		return resultPay, nil
	}
	return resultPay, PaymentSystem
}

func (scpp *SaleCoolerPrePaid) SendMassage(tran *transaction.Transaction, PaymentSystem interfacePayment.Payment) map[string]interface{} {
	result := make(map[string]interface{})
	resultMessage := make(map[string]interface{})
	var errSend error
	resultMessage["tid"] = tran.Config.Tid
	resultMessage["sum"] = tran.Payment.Sum
	resultMessage["wid"] = tran.Products[0]["ware_id"]
	resultMessage["m"] = 1
	resultMessage["a"] = 1
	publicher, err := scpp.QueueManager.AddPublisher(scpp.QueueManager.ContextRabbit, fmt.Sprintf("ephor.1.dev.%v", tran.Config.Imei), "", fmt.Sprintf("ephor.1.dev.%v", tran.Config.Imei))
	if err != nil {
		return scpp.returnMoney(tran, PaymentSystem)
	}
	err = publicher.SendMessage(scpp.QueueManager.ContextRabbit, resultMessage)
	if err != nil {
		for _, v := range scpp.cfg.RabbitMq.BackOffPolicySendMassage {
			time.Sleep(v * time.Second)
			if errSend = publicher.SendMessage(scpp.QueueManager.ContextRabbit, resultMessage); errSend != nil {
				log.Printf("Fail send massage to device with imei: %v, wait time to retry send is %v", tran.Config.Imei, v)
				time.Sleep(v * time.Second)
				continue
			}
			break
		}
		if errSend != nil {
			return scpp.returnMoney(tran, PaymentSystem)
		}
	}
	result["status"] = transaction.TransactionState_MoneyHoldWait
	return result
}

func (scpp *SaleCoolerPrePaid) WaitMassage(tran *transaction.Transaction, PaymentSystem interfacePayment.Payment) map[string]interface{} {
	timer := time.NewTimer(scpp.cfg.RabbitMq.ExecuteTimeSeconds * time.Second)
	log.Printf("\n [x] %s", "Timer")
	request := make(map[string]interface{})
	select {
	case <-timer.C:
		{
			request = scpp.returnMoney(tran, PaymentSystem)
			request["status"] = transaction.TransactionState_ReturnMoney
			request["error"] = "Время ожидания ответа от автомата истекло"
			fmt.Println("Timer is end")
			timer.Stop()
		}
	case result, ok := <-tran.ChannelMessage:
		{
			if !ok {
				request = scpp.returnMoney(tran, PaymentSystem)
				request["status"] = transaction.TransactionState_ReturnMoney
				request["error"] = "Трананзакция завершена в ручную"
				fmt.Println("Timer is end")
				timer.Stop()
			}
			response := responseQueueManager.ResponseQueue{}
			json.Unmarshal(result, &response)
			request["status"] = response.St
			request["error"] = response.Err
			request["ware_id"] = response.Wid
			request["sum"] = response.Sum
			request["action"] = response.A
			request["device_imei"] = response.D
			timer.Stop()
			if scpp.cfg.Debug {
				log.Printf("\n [x] %s", "Timer ok")
				log.Println("Timer ok")
			}
		}
	}
	return request
}

func (scpp *SaleCoolerPrePaid) Fiscal(tran *transaction.Transaction) map[string]interface{} {
	return make(map[string]interface{})
	// resultFiscal := make(map[string]interface{})
	// automat := scpp.Dispetcher.StoreAutomat.GetOneById(tran.Config.AutomatId, tran.Config.AccountId)
	// if len(automat) < 1 {
	// 	resultFiscal["status"] = false
	// 	resultFiscal["error"] = "невозможно взять автомат"
	// 	return resultFiscal
	// }
	// options := make(map[string]interface{})
	// options["account_id"] = tran.Config.AccountId
	// options["automat_id"] = automat["id"]
	// options["to_date"] = "null"
	// locationAutomat, err := scpp.Dispetcher.StoreAutomatLocation.GetWithOptions(options)
	// if err != nil {
	// 	resultFiscal["status"] = false
	// 	resultFiscal["error"] = fmt.Sprintf("%v", err)
	// 	return resultFiscal
	// }
	// if len(locationAutomat) < 1 {
	// 	resultFiscal["status"] = false
	// 	resultFiscal["error"] = "автомат не стоит на точке"
	// 	return resultFiscal
	// }
	// tran.Fiscal.Config.Id = int64(locationAutomat[0].Fr_id)
	// tran.Point.Id = parserTypes.ParseTypeInterfaceToInt(locationAutomat[0].Company_point_id)
	// entryPoint, errPoint := scpp.Dispetcher.StorePoint.GetOneById(locationAutomat[0].Company_point_id, tran.Config.Account_id)
	// if errPoint != nil {
	// 	resultFiscal["status"] = false
	// 	resultFiscal["error"] = fmt.Sprintf("%v", err)
	// 	return resultFiscal
	// }
	// if len(entryPoint) < 1 {
	// 	resultFiscal["status"] = false
	// 	resultFiscal["error"] = "нет точки"
	// 	return resultFiscal
	// }
	// tran.Point.Address = entryPoint[0].Address
	// tran.Point.PointName = entryPoint[0].Name
	// return resultFiscal
}

func (scpp *SaleCoolerPrePaid) returnMoney(tran *transaction.Transaction, PaymentSystem interfacePayment.Payment) map[string]interface{} {
	return scpp.PaymentM.ReturnMoney(tran, PaymentSystem)
}

func (scpp *SaleCoolerPrePaid) returnRoutine() {
	runtime.Goexit()
}
