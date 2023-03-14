package manager

import (
	fr "ephorservice/ephorsale/fiscal/interface/fr"
	transaction "ephorservice/ephorsale/transaction/transaction_struct"
)

type ManagerFiscal interface {
	GetFr() fr.Fiscal
	Fiscal() error
	Status() error
	ValidateData(tran *transaction.Transaction) map[string]interface{}
	Send(protocol string, tran *transaction.Transaction) error
}
