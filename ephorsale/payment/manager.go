package payment

import (
	"context"
	config "ephorservices/config"
	"ephorservices/ephorsale/payment/factory"
	"ephorservices/ephorsale/payment/interfacePayment"
	transaction "ephorservices/ephorsale/transaction"
	"fmt"
	"log"
	"sync"
)

var (
	SBER    uint8 = 1
	VENDPAY uint8 = 2
)

var typePayment = [...]uint8{SBER, VENDPAY}

type PaymentManager struct {
	Status   int
	Payments map[uint8]interfacePayment.Payment
	cfg      *config.Config
	ctx 	 context.Context
	rmutex   sync.RWMutex
}

func Init(cfg *config.Config,ctx context.Context) (*PaymentManager,error) {
	payments := make(map[uint8]interfacePayment.Payment)
	var newRmutex sync.RWMutex
	manager := PaymentManager{
		cfg:      cfg,
		Payments: payments,
		rmutex:   newRmutex,
	}
	manager.InitPayments(cfg)
	return &manager,nil
}

func (pm *PaymentManager) InitPayments(conf *config.Config) {
	for _, item := range typePayment {
		pm.Payments[item] = factory.NewPayment(item, conf)
	}
}

func (pm *PaymentManager) GetPaymentOfType(tp uint8) interfacePayment.Payment {
	payment, ok := pm.Payments[tp]
	if !ok {
		return nil
	}
	return payment
}

func (pm *PaymentManager) StartPrepayment(tran *transaction.Transaction, payment interfacePayment.Payment) map[string]interface{} {
	result := make(map[string]interface{})
	resultHold := payment.HoldMoney(tran)
	if resultHold["status"] == false {
		result["status"] = transaction.TransactionState_Error
		result["ps_order"] = fmt.Sprintf("%v", resultHold["orderId"])
		result["ps_invoice_id"] = fmt.Sprintf("%v", resultHold["invoiceId"])
		result["ps_desc"] = resultHold["description"]
		result["error"] = resultHold["message"]
		return result
	}
	orderId := fmt.Sprintf("%v", resultHold["orderId"])
	invoiceId := fmt.Sprintf("%v", resultHold["invoiceId"])
	log.Println(orderId)
	log.Printf("%+v", tran)
	resultHold = payment.DebitHoldMoney(tran)
	if resultHold["status"] == false {
		result["status"] = transaction.TransactionState_Error
		result["ps_order"] = orderId
		result["ps_invoice_id"] = invoiceId
		result["ps_desc"] = resultHold["description"]
		result["error"] = resultHold["message"]
		return result
	}
	if tid, exist := resultHold["tid"]; exist == true {
		result["ps_tid"] = tid
	}
	result["ps_desc"] = resultHold["description"]
	result["ps_order"] = orderId
	result["ps_invoice_id"] = invoiceId
	result["status"] = transaction.TransactionState_MoneyDebitOk
	log.Printf("%+v", result)
	return result
}

func (pm *PaymentManager) StartPostpaid(tran *transaction.Transaction, payment interfacePayment.Payment) map[string]interface{} {
	result := make(map[string]interface{})
	resultHold := payment.HoldMoney(tran)
	if resultHold["status"] == false {
		result["status"] = transaction.TransactionState_Error
		result["ps_order"] = "none"
		result["ps_desc"] = resultHold["description"]
		result["error"] = resultHold["message"]
		return result
	}
	log.Printf("%+v", resultHold)
	result["ps_desc"] = resultHold["description"]
	result["ps_order"] = resultHold["orderId"]
	result["ps_invoice_id"] = resultHold["invoiceId"]
	result["status"] = transaction.TransactionState_MoneyDebitOk
	return result
}

func (pm *PaymentManager) StartDebitHoldMoney(tran *transaction.Transaction, payment interfacePayment.Payment) map[string]interface{} {
	result := make(map[string]interface{})
	resultHold := payment.DebitHoldMoney(tran)
	if resultHold["status"] == false {
		result["status"] = transaction.TransactionState_Error
		result["ps_desc"] = resultHold["message"]
		return result
	}
	result["ps_desc"] = resultHold["description"]
	result["status"] = transaction.TransactionState_MoneyDebitOk
	return result
}

func (pm *PaymentManager) SatusHoldMoney(tran *transaction.Transaction, payment interfacePayment.Payment) map[string]interface{} {
	result := make(map[string]interface{})
	resultHold := payment.GetStatusHoldMoney(tran)
	if resultHold["status"] == false {
		result["status"] = transaction.TransactionState_Error
		result["ps_desc"] = resultHold["message"]
		return result
	}
	result["ps_desc"] = resultHold["description"]
	result["status"] = transaction.TransactionState_MoneyDebitOk
	return result
}

func (pm *PaymentManager) StartPayment(tran *transaction.Transaction) (map[string]interface{}, interfacePayment.Payment) {
	result := make(map[string]interface{})
	payment := pm.GetPaymentOfType(tran.Payment.Type)
	if payment == nil {
		result["status"] = transaction.TransactionState_Error
		result["ps_desc"] = "no available bank"
		result["error"] = "no available bank"
		return result, payment
	}
	result["status"] = transaction.TransactionState_MoneyHoldStart
	return result, payment
}

func (pm *PaymentManager) ReturnMoney(tran *transaction.Transaction, payment interfacePayment.Payment) map[string]interface{} {
	resultReturnMoney := payment.ReturnMoney(tran)
	if resultReturnMoney["status"] != transaction.TransactionState_ReturnMoney {
		resultReturnMoney["status"] = transaction.TransactionState_Error
		resultReturnMoney["ps_desc"] = resultReturnMoney["message"]
		return resultReturnMoney
	}
	resultReturnMoney["ps_desc"] = resultReturnMoney["description"]
	return resultReturnMoney
}
