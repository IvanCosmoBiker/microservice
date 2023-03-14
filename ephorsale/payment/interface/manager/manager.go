package manager

import (
	"ephorservices/ephorsale/payment/interface/payment"
	transaction "ephorservices/ephorsale/transaction/transaction_struct"
)

type ManagerPayment interface {
	Hold(tran *transaction.Transaction) map[string]interface{}
	Debit(tran *transaction.Transaction) map[string]interface{}
	Satus(tran *transaction.Transaction) map[string]interface{}
	Payment(tran *transaction.Transaction) map[string]interface{}
	Return(tran *transaction.Transaction) map[string]interface{}
	InitPayment(debug bool)
	SetPayment(mapPayments map[uint8]payment.Payment)
	GetPaymentOfType(tp uint8) (payment.Payment, error)
}
