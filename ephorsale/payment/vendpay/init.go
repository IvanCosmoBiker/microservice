package vendpay

import (
	"ephorservices/ephorsale/payment/interface/payment"
	transaction "ephorservices/ephorsale/transaction/transaction_struct"
	"fmt"
	"math"
)

type VendPay struct {
	State   int
	Name    string
	Counter int
	Debug   bool
}

type NewVendStruct struct {
	VendPay
}

func (vend NewVendStruct) New(debug bool) payment.Payment /* тип interfaceBank.Bank*/ {
	return &NewVendStruct{
		VendPay: VendPay{
			Name:    "VendPay",
			Counter: 0,
			State:   0,
			Debug:   debug,
		},
	}
}
func (v *VendPay) GetName() string {
	return v.Name
}

func (v *VendPay) Hold(tran *transaction.Transaction) map[string]interface{} {
	ResponsePaymentSystem := make(map[string]interface{})
	sum := float64(tran.Payment.Sum)
	ResponsePaymentSystem["status"] = true
	ResponsePaymentSystem["message"] = fmt.Sprintf("Сумма %.2f удержана, ожидайте завершения транзакции", math.Round((sum / 100)))
	ResponsePaymentSystem["description"] = fmt.Sprintf("Сумма %.2f удержана, ожидайте завершения транзакции", math.Round((sum / 100)))
	ResponsePaymentSystem["code"] = transaction.TransactionState_MoneyHoldWait
	ResponsePaymentSystem["orderId"] = "vendpay"
	return ResponsePaymentSystem
}

func (v *VendPay) Debit(tran *transaction.Transaction) map[string]interface{} {
	ResponsePaymentSystem := make(map[string]interface{})
	ResponsePaymentSystem["status"] = true
	ResponsePaymentSystem["message"] = "Деньги списаны"
	ResponsePaymentSystem["description"] = "Деньги списаны"
	ResponsePaymentSystem["code"] = transaction.TransactionState_MoneyDebitOk
	return ResponsePaymentSystem
}

func (v *VendPay) Status(tran *transaction.Transaction) map[string]interface{} {
	ResponsePaymentSystem := make(map[string]interface{})
	ResponsePaymentSystem["status"] = true
	ResponsePaymentSystem["message"] = "Оплата подтверждена"
	ResponsePaymentSystem["description"] = "Оплата подтверждена"
	ResponsePaymentSystem["code"] = transaction.TransactionState_MoneyDebitOk
	return ResponsePaymentSystem
}

func (v *VendPay) Return(tran *transaction.Transaction) map[string]interface{} {
	ResponsePaymentSystem := make(map[string]interface{})
	ResponsePaymentSystem["status"] = true
	ResponsePaymentSystem["message"] = "Деньги возвращены"
	ResponsePaymentSystem["description"] = "Деньги возвращены"
	ResponsePaymentSystem["code"] = transaction.TransactionState_ReturnMoney
	return ResponsePaymentSystem
}

func (v *VendPay) Timeout() {
}
