package store

import (
	automat_location_model "ephorservices/internal/model/schema/account/automatlocation/model"
	model_interface "ephorservices/pkg/orm/model"
	request "ephorservices/pkg/orm/request"
	store_parent "ephorservices/pkg/orm/store"
)

type StoreAutomatLocation struct {
	*store_parent.Store
}

func New(accountNumber ...int) *StoreAutomatLocation {
	AccountNumber := 0
	if len(accountNumber) != 0 {
		AccountNumber = accountNumber[0]
	}
	Model := automat_location_model.New()
	store := &StoreAutomatLocation{
		store_parent.New(Model),
	}
	store.Store.SetAccountNumber(AccountNumber)
	return store
}

func (sal *StoreAutomatLocation) GetStructModel(model model_interface.Model) *automat_location_model.AutomatLocationModel {
	if model != nil {
		return model.(*automat_location_model.AutomatLocationModel)
	}
	return &automat_location_model.AutomatLocationModel{}
}

func (sal *StoreAutomatLocation) AddByParams(params map[string]interface{}) (*automat_location_model.AutomatLocationModel, error) {
	model, err := sal.Store.AddByParams(params)
	Model := sal.GetStructModel(model)
	return Model, err
}

func (sal *StoreAutomatLocation) SetByParams(params map[string]interface{}) (*automat_location_model.AutomatLocationModel, error) {
	model, err := sal.Store.SetByParams(params)
	Model := sal.GetStructModel(model)
	return Model, err
}

func (sal *StoreAutomatLocation) GetOneById(id int) (*automat_location_model.AutomatLocationModel, error) {
	model, err := sal.Store.GetOneById(id)
	Model := sal.GetStructModel(model)
	return Model, err
}

func (sal *StoreAutomatLocation) GetOneBy(req *request.Request) (*automat_location_model.AutomatLocationModel, error) {
	model, err := sal.Store.GetOneBy(req)
	Model := sal.GetStructModel(model)
	return Model, err
}
