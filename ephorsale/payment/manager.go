package payment

import (
	"context"
	"ephorservices/ephorsale/payment/interface/manager"
	"ephorservices/ephorsale/payment/interface/payment"
	sberPay "ephorservices/ephorsale/payment/sberpay"
	modulSbp "ephorservices/ephorsale/payment/sbp/modul"
	vendpay "ephorservices/ephorsale/payment/vendpay"
	transaction "ephorservices/ephorsale/transaction/transaction_struct"
	logger "ephorservices/pkg/logger"
	"errors"
	"sync"
)

var (
	SBERPAY  uint8 = 3
	VENDPAY  uint8 = 2
	SBPMODUL uint8 = 5
)

var TypePayment = [3]uint8{SBERPAY, VENDPAY, SBPMODUL}

type Manager struct {
	Status   int
	Payments map[uint8]payment.Payment
	Debug    bool
	Ctx      context.Context
	rmutex   sync.RWMutex
}

var Payment manager.ManagerPayment

func New(debug bool, ctx context.Context) manager.ManagerPayment {
	payments := make(map[uint8]payment.Payment)
	var newRmutex sync.RWMutex
	manager := &Manager{
		Payments: payments,
		rmutex:   newRmutex,
		Debug:    debug,
	}
	manager.InitPayment(debug)
	Payment = manager
	return manager
}

func (pm *Manager) InitPayment(debug bool) {
	for _, item := range TypePayment {
		pm.Payments[item] = pm.NewPayment(item, debug)
	}
}

// instance of type banks
var SberPay sberPay.NewSberPayStruct
var vendPay vendpay.NewVendStruct
var sbpModul modulSbp.NewSbpModul

func (pm *Manager) NewPayment(TypePayment uint8, debug bool) payment.Payment {
	switch TypePayment {
	case payment.TypeSperPay:
		return SberPay.New(debug)
	case payment.TypeVendPay:
		return vendPay.New(debug)
	case payment.TypeModulSbp:
		return sbpModul.New(debug)
	}
	return nil
}

func (pm *Manager) GetPaymentOfType(tp uint8) (payment.Payment, error) {
	payment, ok := pm.Payments[tp]
	if !ok {
		return nil, errors.New("No Avalable Type Payment System")
	}
	return payment, nil
}

func (pm *Manager) SetPayment(mapPayments map[uint8]payment.Payment) {
	pm.Payments = mapPayments
}

func (pm *Manager) Hold(tran *transaction.Transaction) map[string]interface{} {
	result := make(map[string]interface{})
	payment, err := pm.GetPaymentOfType(tran.Payment.Type)
	if err != nil {
		result["status"] = transaction.TransactionState_Error
		result["ps_order"] = "none"
		result["ps_desc"] = err.Error()
		result["error"] = err.Error()
		logger.Log.Error(err.Error())
		return result
	}
	resultHold := payment.Hold(tran)
	if resultHold["status"] == false {
		result["status"] = transaction.TransactionState_Error
		result["ps_order"] = "none"
		result["ps_desc"] = resultHold["description"]
		result["error"] = resultHold["message"]
		logger.Log.Error(resultHold["message"])
		return result
	}
	if tid, exist := resultHold["tid"]; !exist {
		result["ps_tid"] = tid
	}
	result["ps_desc"] = resultHold["description"]
	result["ps_order"] = resultHold["orderId"]
	result["ps_invoice_id"] = resultHold["invoiceId"]
	result["error"] = ""
	result["status"] = transaction.TransactionState_MoneyHoldWait
	if pm.Debug {
		logger.Log.Debugf("Hold %+v", result)
	}
	return result
}

func (pm *Manager) Debit(tran *transaction.Transaction) map[string]interface{} {
	result := make(map[string]interface{})
	payment, err := pm.GetPaymentOfType(tran.Payment.Type)
	if err != nil {
		result["status"] = transaction.TransactionState_Error
		result["ps_order"] = "none"
		result["ps_desc"] = err.Error()
		result["error"] = err.Error()
		logger.Log.Error(err.Error())
		return result
	}
	resultHold := payment.Debit(tran)
	if resultHold["status"] == false {
		result["status"] = transaction.TransactionState_Error
		result["ps_desc"] = resultHold["message"]
		logger.Log.Error(resultHold["message"])
		return result
	}
	result["ps_desc"] = resultHold["description"]
	result["status"] = transaction.TransactionState_MoneyDebitOk
	if pm.Debug {
		logger.Log.Debugf("Debit %+v", result)
	}
	return result
}

func (pm *Manager) Satus(tran *transaction.Transaction) map[string]interface{} {
	result := make(map[string]interface{})
	payment, err := pm.GetPaymentOfType(tran.Payment.Type)
	if err != nil {
		result["status"] = transaction.TransactionState_Error
		result["ps_order"] = "none"
		result["ps_desc"] = err.Error()
		result["error"] = err.Error()
		logger.Log.Error(err.Error())
		return result
	}
	resultHold := payment.Status(tran)
	if resultHold["status"] == false {
		result["status"] = transaction.TransactionState_Error
		result["ps_desc"] = resultHold["message"]
		logger.Log.Error(resultHold["message"])
		return result
	}
	result["ps_desc"] = resultHold["message"]
	result["status"] = transaction.TransactionState_MoneyDebitWait
	if pm.Debug {
		logger.Log.Debugf("Get Status %+v", result)
	}
	return result
}

func (pm *Manager) Payment(tran *transaction.Transaction) map[string]interface{} {
	result := make(map[string]interface{})
	_, err := pm.GetPaymentOfType(tran.Payment.Type)
	if err != nil {
		result["status"] = transaction.TransactionState_Error
		result["ps_order"] = "none"
		result["ps_desc"] = err.Error()
		result["error"] = err.Error()
		logger.Log.Error(err.Error())
		return result
	}
	result["status"] = transaction.TransactionState_MoneyHoldStart
	if pm.Debug {
		logger.Log.Debugf("Payment %+v", result)
	}
	return result
}

func (pm *Manager) Return(tran *transaction.Transaction) map[string]interface{} {
	resultReturnMoney := make(map[string]interface{})
	payment, err := pm.GetPaymentOfType(tran.Payment.Type)
	if err != nil {
		resultReturnMoney["status"] = transaction.TransactionState_Error
		resultReturnMoney["ps_order"] = "none"
		resultReturnMoney["ps_desc"] = err.Error()
		resultReturnMoney["error"] = err.Error()
		logger.Log.Error(err.Error())
		return resultReturnMoney
	}
	resultReturnMoney = payment.Return(tran)
	if resultReturnMoney["status"] != transaction.TransactionState_ReturnMoney {
		resultReturnMoney["status"] = transaction.TransactionState_Error
		resultReturnMoney["ps_desc"] = resultReturnMoney["description"]
		logger.Log.Error(resultReturnMoney["description"])
		return resultReturnMoney
	}
	resultReturnMoney["ps_desc"] = resultReturnMoney["description"]
	if pm.Debug {
		logger.Log.Debugf("Return %+v", resultReturnMoney)
	}
	return resultReturnMoney
}
