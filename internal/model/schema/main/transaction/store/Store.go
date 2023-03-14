package store

import (
	transaction_model "ephorservices/internal/model/schema/main/transaction/model"
	model_interface "ephorservices/pkg/orm/model"
	store_parent "ephorservices/pkg/orm/store"
)

type StoreTransaction struct {
	*store_parent.Store
}

func New(accountNumber ...int) *StoreTransaction {
	AccountNumber := 0
	if len(accountNumber) != 0 {
		AccountNumber = accountNumber[0]
	}
	Model := transaction_model.New()
	store := &StoreTransaction{
		store_parent.New(Model),
	}
	store.Store.SetAccountNumber(AccountNumber)
	return store
}

func (st *StoreTransaction) GetStructModel(model model_interface.Model) *transaction_model.TransactionModel {
	if model != nil {
		return model.(*transaction_model.TransactionModel)
	}
	return &transaction_model.TransactionModel{}
}

func (st *StoreTransaction) AddByParams(params map[string]interface{}) (*transaction_model.TransactionModel, error) {
	model, err := st.Store.AddByParams(params)
	Model := st.GetStructModel(model)
	return Model, err
}

func (st *StoreTransaction) SetByParams(params map[string]interface{}) (*transaction_model.TransactionModel, error) {
	model, err := st.Store.SetByParams(params)
	Model := st.GetStructModel(model)
	return Model, err
}

func (st *StoreTransaction) GetOneById(id int) (*transaction_model.TransactionModel, error) {
	model, err := st.Store.GetOneById(id)
	Model := st.GetStructModel(model)
	return Model, err
}
