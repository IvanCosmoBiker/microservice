package interfaceSale

import (
	transaction "ephorservices/ephorsale/transaction"
	"ephorservices/ephorsale/payment/interfacePayment"
)

const (
	Type_Prepayment uint8 = 1
	Type_PostPaid   uint8 = 2
)

const (
	TypeCoffee      uint8 = 0
	TypeSnack       uint8 = 1
	TypeHoreca      uint8 = 2
	TypeSodaWater   uint8 = 3
	TypeMechanical  uint8 = 4
	TypeComb        uint8 = 5
	TypeMicromarket uint8 = 6
	TypeCooler      uint8 = 7
)

type Sale interface {
	Start(tran *transaction.Transaction)
	Payment(tran *transaction.Transaction) (map[string]interface{},interfacePayment.Payment)
	SendMassage(tran *transaction.Transaction,paymenySystem interfacePayment.Payment) map[string]interface{}
	WaitMassage(tran *transaction.Transaction,paymenySystem interfacePayment.Payment) map[string]interface{}
	Fiscal(tran *transaction.Transaction) map[string]interface{}
}
