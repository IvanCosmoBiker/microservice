package interfacePayment

import (
	transaction "ephorservices/ephorsale/transaction"
)

var (
	TypeSber     uint8 = 1
	TypeVendPay  uint8 = 2
	TypeLifePay  uint8 = 4
	TypeModulSbp uint8 = 5
	TypeSkbSbp   uint8 = 6
)

type Payment interface {
	HoldMoney(payment *transaction.Transaction) map[string]interface{}
	GetStatusHoldMoney(payment *transaction.Transaction) map[string]interface{}
	DebitHoldMoney(payment *transaction.Transaction) map[string]interface{}
	ReturnMoney(payment *transaction.Transaction) map[string]interface{}
	Timeout()
}
