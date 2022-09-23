package automatpostpaid

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

type SaleAutomatPostPaid struct {
	cfg          *config.Config
	QueueManager *rabbit.Manager
	FiscalM      *fiscal.FiscalManager
	PaymentM     *payment.PaymentManager
	Dispetcher   *transaction.TransactionDispetcher
}

type NewSaleAutomatPostPaid struct {
	SaleAutomatPostPaid
}

func (newA *NewSaleAutomatPostPaid) New(conf *config.Config, rabbitMq *rabbit.Manager, fiscalM *fiscal.FiscalManager, paymentM *payment.PaymentManager, dispether *transaction.TransactionDispetcher) interfaceSale.Sale {
	return &NewSaleAutomatPostPaid{
		SaleAutomatPostPaid: SaleAutomatPostPaid{
			cfg:          conf,
			QueueManager: rabbitMq,
			FiscalM:      fiscalM,
			PaymentM:     paymentM,
			Dispetcher:   dispether,
		},
	}
}

func (sapp *SaleAutomatPostPaid) Start(tran *transaction.Transaction) {
	resultDb := make(map[string]interface{})
	if tran.Payment.Sum == 0 {
		resultDb["id"] = tran.Config.Tid
		resultDb["error"] = "Не найдены товары для продолжения оплаты"
		resultDb["status"] = transaction.TransactionState_Error
		sapp.Dispetcher.StoreTransaction.SetByParams(resultDb)
		sapp.Dispetcher.RemoveTransaction(tran.Config.Tid)
		sapp.returnRoutine()
	}
	keyReplayProtection := tran.Config.AutomatId + tran.Config.AccountId
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
	resultStatusHoldMoney := sapp.PaymentM.SatusHoldMoney(tran, PaymentSystem)
	if parserTypes.ParseTypeInterfaceToInt(resultStatusHoldMoney["status"]) != transaction.TransactionState_MoneyDebitOk {
		resultStatusHoldMoney["id"] = tran.Config.Tid
		resultStatusHoldMoney["error"] = "Ошибка оплаты"
		sapp.Dispetcher.StoreTransaction.SetByParams(resultStatusHoldMoney)
		sapp.Dispetcher.RemoveTransaction(tran.Config.Tid)
		sapp.returnRoutine()
	}
	resultStatusHoldMoney["id"] = tran.Config.Tid
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
		resultReturnMoney := sapp.returnMoney(tran, PaymentSystem)
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
	if parserTypes.ParseTypeInterfaceToInt(responseRabbit["status"]) != transaction.VendState_Approving && parserTypes.ParseTypeInterfaceToInt(responseRabbit["status"]) != transaction.VendState_Vending {
		resultReturnMoney := sapp.returnMoney(tran, PaymentSystem)
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
	resultDb["status"] = tran.GetStatusServer(parserTypes.ParseTypeInterfaceToInt(responseRabbit["status"]))
	sapp.Dispetcher.StoreTransaction.SetByParams(resultDb)
	reqAutomatConfig := make(map[string]interface{})
	reqAutomatConfig["to_date"] = nil
	reqAutomatConfig["account_id"] = tran.Config.AccountId
	reqAutomatConfig["automat_id"] = tran.Config.AutomatId
	automatConfig, errC := sapp.Dispetcher.StoreAutomatConfig.GetOneWithOptions(reqAutomatConfig)
	log.Printf("%+v", automatConfig)
	if errC != nil {
		resultReturnMoney := sapp.returnMoney(tran, PaymentSystem)
		resultDb["id"] = tran.Config.Tid
		resultDb["ps_desc"] = resultReturnMoney["ps_desc"]
		resultDb["error"] = "нет конфигурации"
		resultDb["status"] = transaction.TransactionState_Error
		sapp.Dispetcher.StoreTransaction.SetByParams(resultDb)
		sapp.Dispetcher.RemoveTransaction(tran.Config.Tid)
		sapp.Dispetcher.RemoveReplayProtection(keyReplayProtection)
		sapp.returnRoutine()
	}
	reqConfigProduct := make(map[string]interface{})
	reqConfigProduct["config_id"] = automatConfig.Config_id
	reqConfigProduct["account_id"] = tran.Config.AccountId
	reqConfigProduct["ware_id"] = responseRabbit["ware_id"]
	log.Printf("%v", reqConfigProduct)
	configProduct, errP := sapp.Dispetcher.StoreConfigProduct.GetOneWithOptions(reqConfigProduct)
	log.Printf("%+v", configProduct)
	if errP != nil {
		log.Printf("%+v", errP)
		resultReturnMoney := sapp.returnMoney(tran, PaymentSystem)
		resultDb["id"] = tran.Config.Tid
		resultDb["ps_desc"] = resultReturnMoney["ps_desc"]
		resultDb["error"] = "нет такого продукта в ассортименте"
		resultDb["status"] = transaction.TransactionState_Error
		sapp.Dispetcher.StoreTransaction.SetByParams(resultDb)
		sapp.Dispetcher.RemoveTransaction(tran.Config.Tid)
		sapp.Dispetcher.RemoveReplayProtection(keyReplayProtection)
		sapp.returnRoutine()
	}
	ware, errW := sapp.Dispetcher.StoreWare.GetOneById(configProduct.Ware_id, tran.Config.AccountId)
	log.Printf("%+v", ware)
	if errW != nil {
		log.Printf("%+v", errW)
		resultReturnMoney := sapp.returnMoney(tran, PaymentSystem)
		resultDb["id"] = tran.Config.Tid
		resultDb["ps_desc"] = resultReturnMoney["ps_desc"]
		resultDb["error"] = "нет такого товара в номенклатуре"
		resultDb["status"] = transaction.TransactionState_Error
		sapp.Dispetcher.StoreTransaction.SetByParams(resultDb)
		sapp.Dispetcher.RemoveTransaction(tran.Config.Tid)
		sapp.Dispetcher.RemoveReplayProtection(keyReplayProtection)
		sapp.returnRoutine()
	}
	tran.Payment.DebitSum = parserTypes.ParseTypeInterfaceToInt(responseRabbit["sum"])
	resultDebitMoney := sapp.PaymentM.StartDebitHoldMoney(tran, PaymentSystem)
	if parserTypes.ParseTypeInterfaceToInt(resultDebitMoney["status"]) != transaction.TransactionState_MoneyDebitOk {
		resultReturnMoney := sapp.returnMoney(tran, PaymentSystem)
		resultDb["id"] = tran.Config.Tid
		resultDb["ps_desc"] = resultReturnMoney["ps_desc"]
		resultDb["error"] = "не удалось списать деньги"
		resultDb["status"] = transaction.TransactionState_Error
		sapp.Dispetcher.StoreTransaction.SetByParams(resultDb)
		sapp.Dispetcher.RemoveTransaction(tran.Config.Tid)
		sapp.Dispetcher.RemoveReplayProtection(keyReplayProtection)
		sapp.returnRoutine()
	}
	resultDb["id"] = tran.Config.Tid
	resultDb["ps_desc"] = resultDebitMoney["ps_desc"]
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
		resultReturnMoney := sapp.returnMoney(tran, PaymentSystem)
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
	resultDb["status"] = tran.GetStatusServer(transaction.VendState_VendOk)
	sapp.Dispetcher.StoreTransaction.SetByParams(resultDb)
	sapp.Dispetcher.RemoveTransaction(tran.Config.Tid)
	sapp.Dispetcher.RemoveReplayProtection(keyReplayProtection)
	sapp.returnRoutine()
}

func (sapp *SaleAutomatPostPaid) Payment(tran *transaction.Transaction) (map[string]interface{}, interfacePayment.Payment) {
	resultDb := make(map[string]interface{})
	resultPayment, PaymentSystem := sapp.PaymentM.StartPayment(tran)
	if resultPayment["status"] != transaction.TransactionState_MoneyHoldStart {
		return resultPayment, nil
	}
	resultDb["id"] = tran.Config.Tid
	resultDb["status"] = transaction.TransactionState_MoneyHoldWait
	sapp.Dispetcher.StoreTransaction.SetByParams(resultDb)
	resultPay := sapp.PaymentM.StartPostpaid(tran, PaymentSystem)
	if resultPay["status"] == transaction.TransactionState_Error {
		return resultPay, nil
	}
	return resultPay, PaymentSystem
}

func (sapp *SaleAutomatPostPaid) SendMassage(tran *transaction.Transaction, PaymentSystem interfacePayment.Payment) map[string]interface{} {
	result := make(map[string]interface{})
	resultMessage := make(map[string]interface{})
	var errSend error
	resultMessage["tid"] = tran.Config.Tid
	resultMessage["sum"] = tran.Payment.Sum
	resultMessage["m"] = 1
	resultMessage["a"] = 2
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

func (sapp *SaleAutomatPostPaid) WaitMassage(tran *transaction.Transaction, PaymentSystem interfacePayment.Payment) map[string]interface{} {
	timer := time.NewTimer(sapp.cfg.RabbitMq.ExecuteTimeSeconds * time.Second)
	log.Printf("\n [x] %s", "Timer")
	Tran := *tran
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
	case result, ok := <-Tran.ChannelMessage:
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

func (sapp *SaleAutomatPostPaid) Fiscal(tran *transaction.Transaction) map[string]interface{} {
	return make(map[string]interface{})
}

func (sapp *SaleAutomatPostPaid) returnMoney(tran *transaction.Transaction, PaymentSystem interfacePayment.Payment) map[string]interface{} {
	return sapp.PaymentM.ReturnMoney(tran, PaymentSystem)
}

func (sapp *SaleAutomatPostPaid) returnRoutine() {
	runtime.Goexit()
}
