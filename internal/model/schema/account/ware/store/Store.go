package ware

import (
	ware_model "ephorservices/internal/model/schema/account/ware/model"
	model_interface "ephorservices/pkg/orm/model"
	store_parent "ephorservices/pkg/orm/store"
)

var (
	Type_NotSet = 0
	Type_Recipe = 1
	Type_Snack  = 2
)

var (
	StateUncomplete = 0
	StateActive     = 1
	StateObsolete   = 2
)

type StoreWare struct {
	*store_parent.Store
}

func New(accountNumber ...int) *StoreWare {
	AccountNumber := 0
	if len(accountNumber) != 0 {
		AccountNumber = accountNumber[0]
	}
	Model := ware_model.New()
	store := &StoreWare{
		store_parent.New(Model),
	}
	store.Store.SetAccountNumber(AccountNumber)
	return store
}

func (sw *StoreWare) GetStructModel(model model_interface.Model) *ware_model.WareModel {
	if model != nil {
		return model.(*ware_model.WareModel)
	}
	return &ware_model.WareModel{}
}

func (sw *StoreWare) AddByParams(params map[string]interface{}) (*ware_model.WareModel, error) {
	model, err := sw.Store.AddByParams(params)
	Model := sw.GetStructModel(model)
	return Model, err
}

func (sw *StoreWare) SetByParams(params map[string]interface{}) (*ware_model.WareModel, error) {
	model, err := sw.Store.SetByParams(params)
	Model := sw.GetStructModel(model)
	return Model, err
}

func (sw *StoreWare) GetOneById(id int) (*ware_model.WareModel, error) {
	model, err := sw.Store.GetOneById(id)
	Model := sw.GetStructModel(model)
	return Model, err
}
