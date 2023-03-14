package interfaceSale

import (
	transaction "ephorservices/ephorsale/transaction/transaction_struct"
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

type FiscalFunc func(tran *transaction.Transaction)
type Sale interface {
	SetFiscalisation(function FiscalFunc)
	Sale(tran *transaction.Transaction) map[string]interface{}
	Payment(tran *transaction.Transaction) map[string]interface{}
	SendMassage(tran *transaction.Transaction) map[string]interface{}
	WaitMassage(tran *transaction.Transaction) map[string]interface{}
	Fiscal(tran *transaction.Transaction) map[string]interface{}
}
