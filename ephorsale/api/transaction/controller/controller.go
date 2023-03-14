package controller

import (
	transactionDispetcher "ephorservices/ephorsale/transaction"
	modelTransaction "ephorservices/pkg/model/schema/main/transaction/model"
	parserTypes "ephorservices/pkg/parser/typeParse"
	"log"
)

type ControllerTransaction struct {
	Dispetcher *transactionDispetcher.TransactionDispetcher
}

func Init(dispetcher *transactionDispetcher.TransactionDispetcher) *ControllerTransaction {
	return &ControllerTransaction{
		Dispetcher: dispetcher,
	}
}

func (ct *ControllerTransaction) GetTransactionActive(filters map[string]interface{}) []*modelTransaction.ReturningStruct {
	result := make([]*modelTransaction.ReturningStruct, 1, 10)
	transactionsActive := ct.Dispetcher.GetTransactions()
	if len(transactionsActive) < 1 {
		return result
	}
	for k, _ := range transactionsActive {
		transactionModel, err := ct.Dispetcher.StoreTransaction.GetOneById(k)
		if err != nil {
			continue
		}
		result = append(result, transactionModel)
	}
	return result
}

func (ct *ControllerTransaction) EndTransactionActive(RequestEnd RequestEndTransaction) (int, bool) {
	transactionsActive := ct.Dispetcher.GetOneTransaction(RequestEnd.Tid)
	log.Println(transactionsActive)
	if transactionsActive == false {
		return 0, false
	}
	transactionModel, err := ct.Dispetcher.StoreTransaction.GetOneById(RequestEnd.Tid)
	if err != nil {
		return 0, false
	}
	log.Printf("%v", transactionModel)
	resultRemoveChannek := ct.Dispetcher.RemoveTransaction(RequestEnd.Tid)
	log.Printf("%v", resultRemoveChannek)
	if resultRemoveChannek == false {
		return 0, false
	}
	keyAutomat := parserTypes.ParseTypeStringInInt(transactionModel.Automat_id) + parserTypes.ParseTypeStringInInt(transactionModel.Account_id)
	resultRemoveProtection := ct.Dispetcher.RemoveReplayProtection(keyAutomat)
	log.Printf("%v", resultRemoveProtection)
	if resultRemoveProtection == false {
		return 0, false
	}
	parametrs := make(map[string]interface{})
	parametrs["id"] = RequestEnd.Tid
	parametrs["ps_desc"] = "Трананзакция завершена в ручную"
	parametrs["error"] = "Трананзакция завершена в ручную"
	parametrs["status"] = transactionDispetcher.TransactionState_EndClient
	ct.Dispetcher.StoreTransaction.SetByParams(parametrs)
	return RequestEnd.Tid, true
}
