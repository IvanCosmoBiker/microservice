package payment

import (
	transaction "ephorservices/ephorsale/transaction/transaction_struct"
)

var (
	TypeSber     uint8 = 1
	TypeVendPay  uint8 = 2
	TypeSperPay  uint8 = 3
	TypeLifePay  uint8 = 4
	TypeModulSbp uint8 = 5
	TypeSkbSbp   uint8 = 6
)

type Payment interface {
	GetName() string
	Hold(payment *transaction.Transaction) map[string]interface{}
	Status(payment *transaction.Transaction) map[string]interface{}
	Debit(payment *transaction.Transaction) map[string]interface{}
	Return(payment *transaction.Transaction) map[string]interface{}
	Timeout()
}
