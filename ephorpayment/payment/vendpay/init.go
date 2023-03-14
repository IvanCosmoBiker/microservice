package vendpay

import (
	config "ephorservices/config"
	"ephorservices/ephorpayment/payment/interface/payment"
	pb "ephorservices/ephorpayment/service"
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

func (vend NewVendStruct) New(conf *config.Config) payment.Payment /* тип interfaceBank.Bank*/ {
	return &NewVendStruct{
		VendPay: VendPay{
			Name:    "VendPay",
			Counter: 0,
			Status:  0,
			cfg:     conf,
		},
	}
}

func (v *VendPay) HoldMoney(tran *pb.Request) map[string]interface{} {
	ResponsePaymentSystem := make(map[string]interface{})
	sum := float64(tran.Sum)
	ResponsePaymentSystem["status"] = true
	ResponsePaymentSystem["message"] = fmt.Sprintf("Сумма %.2f удержана, ожидайте завершения транзакции", math.Round((sum / 100)))
	ResponsePaymentSystem["description"] = fmt.Sprintf("Сумма %.2f удержана, ожидайте завершения транзакции", math.Round((sum / 100)))
	ResponsePaymentSystem["code"] = transaction.TransactionState_MoneyHoldWait
	ResponsePaymentSystem["orderId"] = "vendpay"
	return ResponsePaymentSystem
}

func (v *VendPay) DebitHoldMoney(tran *pb.Request) map[string]interface{} {
	ResponsePaymentSystem := make(map[string]interface{})
	ResponsePaymentSystem["status"] = true
	ResponsePaymentSystem["message"] = "Деньги списаны"
	ResponsePaymentSystem["description"] = "Деньги списаны"
	ResponsePaymentSystem["code"] = transaction.TransactionState_MoneyDebitOk
	return ResponsePaymentSystem
}

func (v *VendPay) GetStatusHoldMoney(tran *pb.Request) map[string]interface{} {
	ResponsePaymentSystem := make(map[string]interface{})
	ResponsePaymentSystem["status"] = true
	ResponsePaymentSystem["message"] = "Оплата подтверждена"
	ResponsePaymentSystem["description"] = "Оплата подтверждена"
	ResponsePaymentSystem["code"] = transaction.TransactionState_MoneyDebitOk
	return ResponsePaymentSystem
}

func (v *VendPay) ReturnMoney(tran *pb.Request) map[string]interface{} {
	ResponsePaymentSystem := make(map[string]interface{})
	ResponsePaymentSystem["status"] = true
	ResponsePaymentSystem["message"] = "Деньги возвращены"
	ResponsePaymentSystem["description"] = "Деньги возвращены"
	ResponsePaymentSystem["code"] = transaction.TransactionState_ReturnMoney
	return ResponsePaymentSystem
}

func (v *VendPay) Timeout() {
}
