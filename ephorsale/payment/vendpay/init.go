package vendpay

import (
	config "ephorservices/config"
	"ephorservices/ephorsale/payment/interfacePayment"
	"ephorservices/ephorsale/transaction"
	"fmt"
	"math"
)

type VendPay struct {
	Status  int
	Name    string
	Counter int
	cfg     *config.Config
}

type NewVendStruct struct {
	VendPay
}

func (vend NewVendStruct) New(conf *config.Config) interfacePayment.Payment /* тип interfaceBank.Bank*/ {
	return &NewVendStruct{
		VendPay: VendPay{
			Name:    "VendPay",
			Counter: 0,
			Status:  0,
			cfg:     conf,
		},
	}
}

func (v *VendPay) HoldMoney(tran *transaction.Transaction) map[string]interface{} {
	ResponsePaymentSystem := make(map[string]interface{})
	sum := float64(tran.Payment.Sum)
	ResponsePaymentSystem["status"] = true
	ResponsePaymentSystem["message"] = fmt.Sprintf("Сумма %.2f удержана, ожидайте завершения транзакции", math.Round((sum / 100)))
	ResponsePaymentSystem["description"] = fmt.Sprintf("Сумма %.2f удержана, ожидайте завершения транзакции", math.Round((sum / 100)))
	ResponsePaymentSystem["code"] = transaction.TransactionState_MoneyHoldWait
	ResponsePaymentSystem["orderId"] = "vendpay"
	return ResponsePaymentSystem
}

func (v *VendPay) DebitHoldMoney(tran *transaction.Transaction) map[string]interface{} {
	ResponsePaymentSystem := make(map[string]interface{})
	ResponsePaymentSystem["status"] = true
	ResponsePaymentSystem["message"] = "Деньги списаны"
	ResponsePaymentSystem["description"] = "Деньги списаны"
	ResponsePaymentSystem["code"] = transaction.TransactionState_MoneyDebitOk
	return ResponsePaymentSystem
}

func (v *VendPay) GetStatusHoldMoney(tran *transaction.Transaction) map[string]interface{} {
	ResponsePaymentSystem := make(map[string]interface{})
	ResponsePaymentSystem["status"] = true
	ResponsePaymentSystem["message"] = "Оплата подтверждена"
	ResponsePaymentSystem["description"] = "Оплата подтверждена"
	ResponsePaymentSystem["code"] = transaction.TransactionState_MoneyDebitOk
	return ResponsePaymentSystem
}

func (v *VendPay) ReturnMoney(tran *transaction.Transaction) map[string]interface{} {
	ResponsePaymentSystem := make(map[string]interface{})
	ResponsePaymentSystem["status"] = true
	ResponsePaymentSystem["message"] = "Деньги возвращены"
	ResponsePaymentSystem["description"] = "Деньги возвращены"
	ResponsePaymentSystem["code"] = transaction.TransactionState_ReturnMoney
	return ResponsePaymentSystem
}

func (v *VendPay) Timeout() {
}



