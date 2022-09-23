package automatprepaid

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

type SaleAutomatPrePaid struct {
	cfg          *config.Config
	QueueManager *rabbit.Manager
	FiscalM      *fiscal.FiscalManager
	PaymentM     *payment.PaymentManager
	Dispetcher   *transaction.TransactionDispetcher
}

type NewSaleAutomatPrePaid struct {
	SaleAutomatPrePaid
}

func (newA NewSaleAutomatPrePaid) New(conf *config.Config, rabbitMq *rabbit.Manager, fiscalM *fiscal.FiscalManager, paymentM *payment.PaymentManager, dispether *transaction.TransactionDispetcher) interfaceSale.Sale {
	return &NewSaleAutomatPrePaid{
		SaleAutomatPrePaid: SaleAutomatPrePaid{
			cfg:          conf,
			QueueManager: rabbitMq,
			FiscalM:      fiscalM,
			PaymentM:     paymentM,
			Dispetcher:   dispether,
		},
	}
}

// создание goroutine
func (sapp *SaleAutomatPrePaid) Start(tran *transaction.Transaction) {
	resultDb := make(map[string]interface{})
	if len(tran.Products) < 1 {
		resultDb["id"] = tran.Config.Tid
		resultDb["ps_desc"] = "Не найдены товары для продолжения оплаты"
		resultDb["error"] = "Не найдены товары для продолжения оплаты"
		resultDb["status"] = transaction.TransactionState_Error
		sapp.Dispetcher.StoreTransaction.SetByParams(resultDb)
		sapp.Dispetcher.RemoveTransaction(tran.Config.Tid)
		sapp.returnRoutine()
	}
	keyReplayProtection := tran.Config.AutomatId + tran.Config.AccountId
	resultReturnMoney := make(map[string]interface{})
	resultPayment, PaymentSystem := sapp.Payment(tran)
	if resultPayment["status"] == transaction.TransactionState_Error {
		resultDb["id"] = tran.Config.Tid
		resultDb["ps_desc"] = resultPayment["ps_desc"]
		resultDb["error"] = "Ошибка"
		resultDb["status"] = resultPayment["status"]
		sapp.Dispetcher.StoreTransaction.SetByParams(resultDb)
		sapp.Dispetcher.RemoveTransaction(tran.Config.Tid)
		sapp.returnRoutine()
	}
	resultPayment["id"] = tran.Config.Tid
	sapp.Dispetcher.StoreTransaction.SetByParams(resultPayment)
	sapp.Dispetcher.AddReplayProtection(keyReplayProtection, tran.Config.AutomatId)
	resultMassage := sapp.SendMassage(tran, PaymentSystem)
	if parserTypes.ParseTypeInterfaceToInt(resultMassage["status"]) != transaction.TransactionState_MoneyHoldWait {
		resultMassage["id"] = tran.Config.Tid
		sapp.Dispetcher.StoreTransaction.SetByParams(resultMassage)
		sapp.Dispetcher.RemoveTransaction(tran.Config.Tid)
		sapp.Dispetcher.RemoveReplayProtection(keyReplayProtection)
		sapp.returnRoutine()
	}
	responseRabbit := sapp.WaitMassage(tran, PaymentSystem)
	if parserTypes.ParseTypeInterfaceToInt(responseRabbit["status"]) == transaction.TransactionState_ReturnMoney {
		resultDb["id"] = tran.Config.Tid
		resultDb["ps_desc"] = responseRabbit["ps_desc"]
		resultDb["error"] = responseRabbit["error"]
		resultDb["status"] = transaction.TransactionState_Error
		sapp.Dispetcher.StoreTransaction.SetByParams(resultDb)
		sapp.Dispetcher.RemoveTransaction(tran.Config.Tid)
		sapp.Dispetcher.RemoveReplayProtection(keyReplayProtection)
		sapp.returnRoutine()
	}
	if parserTypes.ParseTypeInterfaceToInt(responseRabbit["status"]) != transaction.VendState_Session {
		resultReturnMoney = sapp.returnMoney(tran, PaymentSystem)
		resultDb["id"] = tran.Config.Tid
		resultDb["ps_desc"] = resultReturnMoney["ps_desc"]
		resultDb["error"] = tran.GetDescriptionErr(parserTypes.ParseTypeInterfaceToInt(responseRabbit["error"]))
		resultDb["status"] = transaction.TransactionState_Error
		sapp.Dispetcher.StoreTransaction.SetByParams(resultDb)
		sapp.Dispetcher.RemoveTransaction(tran.Config.Tid)
		sapp.Dispetcher.RemoveReplayProtection(keyReplayProtection)
		sapp.returnRoutine()
	}
	resultDb["id"] = tran.Config.Tid
	resultDb["ps_desc"] = tran.GetDescriptionCode(parserTypes.ParseTypeInterfaceToInt(responseRabbit["status"]))
	resultDb["error"] = tran.GetDescriptionErr(parserTypes.ParseTypeInterfaceToInt(responseRabbit["error"]))
	resultDb["status"] = tran.GetStatusServer(transaction.VendState_Session)
	sapp.Dispetcher.StoreTransaction.SetByParams(resultDb)
	responseRabbit = sapp.WaitMassage(tran, PaymentSystem)
	if parserTypes.ParseTypeInterfaceToInt(responseRabbit["status"]) == transaction.TransactionState_ReturnMoney {
		resultDb["id"] = tran.Config.Tid
		resultDb["ps_desc"] = responseRabbit["ps_desc"]
		resultDb["error"] = responseRabbit["error"]
		resultDb["status"] = transaction.TransactionState_Error
		sapp.Dispetcher.StoreTransaction.SetByParams(resultDb)
		sapp.Dispetcher.RemoveTransaction(tran.Config.Tid)
		sapp.Dispetcher.RemoveReplayProtection(keyReplayProtection)
		sapp.returnRoutine()
	}
	if parserTypes.ParseTypeInterfaceToInt(responseRabbit["status"]) != transaction.VendState_Vending {
		resultReturnMoney = sapp.returnMoney(tran, PaymentSystem)
		resultDb["id"] = tran.Config.Tid
		resultDb["ps_desc"] = resultReturnMoney["ps_desc"]
		resultDb["error"] = tran.GetDescriptionErr(parserTypes.ParseTypeInterfaceToInt(responseRabbit["error"]))
		resultDb["status"] = transaction.TransactionState_Error
		sapp.Dispetcher.StoreTransaction.SetByParams(resultDb)
		sapp.Dispetcher.RemoveTransaction(tran.Config.Tid)
		sapp.Dispetcher.RemoveReplayProtection(keyReplayProtection)
		sapp.returnRoutine()
	}
	resultDb["id"] = tran.Config.Tid
	resultDb["ps_desc"] = tran.GetDescriptionCode(parserTypes.ParseTypeInterfaceToInt(responseRabbit["status"]))
	resultDb["error"] = tran.GetDescriptionErr(parserTypes.ParseTypeInterfaceToInt(responseRabbit["error"]))
	resultDb["status"] = tran.GetStatusServer(transaction.VendState_Session)
	sapp.Dispetcher.StoreTransaction.SetByParams(resultDb)
	responseRabbit = sapp.WaitMassage(tran, PaymentSystem)
	if parserTypes.ParseTypeInterfaceToInt(responseRabbit["status"]) == transaction.TransactionState_ReturnMoney {
		resultDb["id"] = tran.Config.Tid
		resultDb["ps_desc"] = responseRabbit["ps_desc"]
		resultDb["error"] = responseRabbit["error"]
		resultDb["status"] = transaction.TransactionState_Error
		sapp.Dispetcher.StoreTransaction.SetByParams(resultDb)
		sapp.Dispetcher.RemoveTransaction(tran.Config.Tid)
		sapp.Dispetcher.RemoveReplayProtection(keyReplayProtection)
		sapp.returnRoutine()
	}
	if parserTypes.ParseTypeInterfaceToInt(responseRabbit["status"]) != transaction.VendState_VendOk {
		resultReturnMoney = sapp.returnMoney(tran, PaymentSystem)
		resultDb["id"] = tran.Config.Tid
		resultDb["ps_desc"] = resultReturnMoney["ps_desc"]
		resultDb["error"] = tran.GetDescriptionErr(parserTypes.ParseTypeInterfaceToInt(responseRabbit["error"]))
		resultDb["status"] = transaction.TransactionState_Error
		sapp.Dispetcher.StoreTransaction.SetByParams(resultDb)
		sapp.Dispetcher.RemoveTransaction(tran.Config.Tid)
		sapp.Dispetcher.RemoveReplayProtection(keyReplayProtection)
		sapp.returnRoutine()
		sapp.Dispetcher.RemoveTransaction(tran.Config.Tid)
		sapp.Dispetcher.RemoveReplayProtection(keyReplayProtection)
		sapp.returnRoutine()
	}
	resultDb["id"] = tran.Config.Tid
	resultDb["ps_desc"] = tran.GetDescriptionCode(parserTypes.ParseTypeInterfaceToInt(responseRabbit["status"]))
	resultDb["error"] = tran.GetDescriptionErr(parserTypes.ParseTypeInterfaceToInt(responseRabbit["error"]))
	resultDb["status"] = tran.GetStatusServer(transaction.VendState_VendOk)
	sapp.Dispetcher.StoreTransaction.SetByParams(resultDb)
	sapp.Dispetcher.RemoveTransaction(tran.Config.Tid)
	sapp.Dispetcher.RemoveReplayProtection(keyReplayProtection)
	sapp.returnRoutine()
}

func (sapp *SaleAutomatPrePaid) Payment(tran *transaction.Transaction) (map[string]interface{}, interfacePayment.Payment) {
	resultDb := make(map[string]interface{})
	resultPayment, PaymentSystem := sapp.PaymentM.StartPayment(tran)
	if resultPayment["status"] != transaction.TransactionState_MoneyHoldStart {
		return resultPayment, nil
	}
	resultDb["id"] = tran.Config.Tid
	resultDb["status"] = transaction.TransactionState_MoneyHoldWait
	sapp.Dispetcher.StoreTransaction.SetByParams(resultDb)
	resultPay := sapp.PaymentM.StartPrepayment(tran, PaymentSystem)
	if resultPay["status"] == transaction.TransactionState_Error {
		return resultPay, nil
	}
	return resultPay, PaymentSystem
}

func (sapp *SaleAutomatPrePaid) SendMassage(tran *transaction.Transaction, PaymentSystem interfacePayment.Payment) map[string]interface{} {
	result := make(map[string]interface{})
	resultMessage := make(map[string]interface{})
	var errSend error
	resultMessage["tid"] = tran.Config.Tid
	resultMessage["sum"] = tran.Payment.Sum
	resultMessage["wid"] = tran.Products[0]["ware_id"]
	resultMessage["m"] = 1
	resultMessage["a"] = 1
	publicher, err := sapp.QueueManager.AddPublisher(sapp.QueueManager.ContextRabbit, fmt.Sprintf("ephor.1.dev.%v", tran.Config.Imei), "", fmt.Sprintf("ephor.1.dev.%v", tran.Config.Imei))
	if err != nil {
		return sapp.returnMoney(tran, PaymentSystem)
	}
	err = publicher.SendMessage(sapp.QueueManager.ContextRabbit, resultMessage)
	if err != nil {
		for _, v := range sapp.cfg.RabbitMq.BackOffPolicySendMassage {
			time.Sleep(v * time.Second)
			if errSend = publicher.SendMessage(sapp.QueueManager.ContextRabbit, resultMessage); errSend != nil {
				log.Printf("Fail send massage to device with imei: %v, wait time to retry send is %v", tran.Config.Imei, v)
				time.Sleep(v * time.Second)
				continue
			}
			break
		}
		if errSend != nil {
			return sapp.returnMoney(tran, PaymentSystem)
		}
	}
	result["status"] = transaction.TransactionState_MoneyHoldWait
	return result
}

func (sapp *SaleAutomatPrePaid) WaitMassage(tran *transaction.Transaction, PaymentSystem interfacePayment.Payment) map[string]interface{} {
	timer := time.NewTimer(sapp.cfg.RabbitMq.ExecuteTimeSeconds * time.Second)
	log.Printf("\n [x] %s", "Timer")
	request := make(map[string]interface{})
	select {
	case <-timer.C:
		{
			request = sapp.returnMoney(tran, PaymentSystem)
			request["status"] = transaction.TransactionState_ReturnMoney
			request["error"] = "Время ожидания ответа от автомата истекло"
			fmt.Println("Timer is end")
			timer.Stop()
		}
	case result, ok := <-tran.ChannelMessage:
		{
			if !ok {
				request = sapp.returnMoney(tran, PaymentSystem)
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
			if sapp.cfg.Debug {
				log.Printf("\n [x] %s", "Timer ok")
				log.Println("Timer ok")
			}
		}
	}
	return request
}

func (sapp *SaleAutomatPrePaid) Fiscal(tran *transaction.Transaction) map[string]interface{} {
	return make(map[string]interface{})
}

func (sapp *SaleAutomatPrePaid) returnMoney(tran *transaction.Transaction, PaymentSystem interfacePayment.Payment) map[string]interface{} {
	return sapp.PaymentM.ReturnMoney(tran, PaymentSystem)
}

func (sapp *SaleAutomatPrePaid) returnRoutine() {
	runtime.Goexit()
}
