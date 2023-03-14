package fr

import (
	fr_model "ephorservices/internal/model/schema/account/fr/model"
	model_interface "ephorservices/pkg/orm/model"
	store_parent "ephorservices/pkg/orm/store"
)

type StoreFr struct {
	*store_parent.Store
}

func New(accountNumber ...int) *StoreFr {
	AccountNumber := 0
	if len(accountNumber) != 0 {
		AccountNumber = accountNumber[0]
	}
	Model := fr_model.New()
	store := &StoreFr{
		store_parent.New(Model),
	}
	store.Store.SetAccountNumber(AccountNumber)
	return store
}

func (sf *StoreFr) GetStructModel(model model_interface.Model) *fr_model.FrModel {
	if model != nil {
		return model.(*fr_model.FrModel)
	}
	return &fr_model.FrModel{}
}

func (sf *StoreFr) AddByParams(params map[string]interface{}) (*fr_model.FrModel, error) {
	model, err := sf.Store.AddByParams(params)
	Model := sf.GetStructModel(model)
	return Model, err
}

func (sf *StoreFr) SetByParams(params map[string]interface{}) (*fr_model.FrModel, error) {
	model, err := sf.Store.SetByParams(params)
	Model := sf.GetStructModel(model)
	return Model, err
}

func (sf *StoreFr) GetOneById(id int) (*fr_model.FrModel, error) {
	model, err := sf.Store.GetOneById(id)
	Model := sf.GetStructModel(model)
	return Model, err
}
