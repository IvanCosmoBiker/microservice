package store

import (
	automat_model "ephorservices/internal/model/schema/account/automat/model"
	model_interface "ephorservices/pkg/orm/model"
	request "ephorservices/pkg/orm/request"
	store_parent "ephorservices/pkg/orm/store"
)

type StoreAutomat struct {
	*store_parent.Store
}

func New(accountNumber ...int) *StoreAutomat {
	AccountNumber := 0
	if len(accountNumber) != 0 {
		AccountNumber = accountNumber[0]
	}
	Model := automat_model.New()
	store := &StoreAutomat{
		store_parent.New(Model),
	}
	store.Store.SetAccountNumber(AccountNumber)
	return store
}

func (sa *StoreAutomat) GetStructModel(model model_interface.Model) *automat_model.AutomatModel {
	if model != nil {
		return model.(*automat_model.AutomatModel)
	}
	return &automat_model.AutomatModel{}
}

func (sa *StoreAutomat) AddByParams(params map[string]interface{}) (*automat_model.AutomatModel, error) {
	model, err := sa.Store.AddByParams(params)
	Model := sa.GetStructModel(model)
	return Model, err
}

func (sa *StoreAutomat) SetByParams(params map[string]interface{}) (*automat_model.AutomatModel, error) {
	model, err := sa.Store.SetByParams(params)
	Model := sa.GetStructModel(model)
	return Model, err
}

func (sa *StoreAutomat) Set(m model_interface.Model) (*automat_model.AutomatModel, error) {
	model, err := sa.Store.Set(m)
	Model := sa.GetStructModel(model)
	return Model, err
}

func (sa *StoreAutomat) GetOneById(id int) (*automat_model.AutomatModel, error) {
	model, err := sa.Store.GetOneById(id)
	Model := sa.GetStructModel(model)
	return Model, err
}

func (sa *StoreAutomat) GetOneBy(req *request.Request) (*automat_model.AutomatModel, error) {
	model, err := sa.Store.GetOneBy(req)
	Model := sa.GetStructModel(model)
	return Model, err
}
