package store

import (
	transaction_product_model "ephorservices/internal/model/schema/main/transactionproduct/model"
	model_interface "ephorservices/pkg/orm/model"
	store_parent "ephorservices/pkg/orm/store"
)

type StoreTransactionProduct struct {
	*store_parent.Store
}

func New(accountNumber ...int) *StoreTransactionProduct {
	AccountNumber := 0
	if len(accountNumber) != 0 {
		AccountNumber = accountNumber[0]
	}
	Model := transaction_product_model.New()
	store := &StoreTransactionProduct{
		store_parent.New(Model),
	}
	store.Store.SetAccountNumber(AccountNumber)
	return store
}

func (stp *StoreTransactionProduct) GetStructModel(model model_interface.Model) *transaction_product_model.TransactionProductModel {
	if model != nil {
		return model.(*transaction_product_model.TransactionProductModel)
	}
	return &transaction_product_model.TransactionProductModel{}
}

func (stp *StoreTransactionProduct) AddByParams(params map[string]interface{}) (*transaction_product_model.TransactionProductModel, error) {
	model, err := stp.Store.AddByParams(params)
	Model := stp.GetStructModel(model)
	return Model, err
}

func (stp *StoreTransactionProduct) SetByParams(params map[string]interface{}) (*transaction_product_model.TransactionProductModel, error) {
	model, err := stp.Store.SetByParams(params)
	Model := stp.GetStructModel(model)
	return Model, err
}

func (stp *StoreTransactionProduct) GetOneById(id int) (*transaction_product_model.TransactionProductModel, error) {
	model, err := stp.Store.GetOneById(id)
	Model := stp.GetStructModel(model)
	return Model, err
}
