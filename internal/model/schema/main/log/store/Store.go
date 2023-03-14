package store

import (
	log_model "ephorservices/internal/model/schema/main/log/model"
	model_interface "ephorservices/pkg/orm/model"
	store_parent "ephorservices/pkg/orm/store"
)

type StoreLog struct {
	*store_parent.Store
}

func New(accountNumber ...int) *StoreLog {
	AccountNumber := 0
	if len(accountNumber) != 0 {
		AccountNumber = accountNumber[0]
	}
	Model := log_model.New()
	store := &StoreLog{
		store_parent.New(Model),
	}
	store.Store.SetAccountNumber(AccountNumber)
	return store
}

func (sl *StoreLog) GetStructModel(model model_interface.Model) *log_model.LogModel {
	if model != nil {
		return model.(*log_model.LogModel)
	}
	return &log_model.LogModel{}
}

func (sl *StoreLog) AddByParams(params map[string]interface{}) (*log_model.LogModel, error) {
	model, err := sl.Store.AddByParams(params)
	Model := sl.GetStructModel(model)
	return Model, err
}

func (sl *StoreLog) SetByParams(params map[string]interface{}) (*log_model.LogModel, error) {
	model, err := sl.Store.SetByParams(params)
	Model := sl.GetStructModel(model)
	return Model, err
}

func (sl *StoreLog) GetOneById(id int) (*log_model.LogModel, error) {
	model, err := sl.Store.GetOneById(id)
	Model := sl.GetStructModel(model)
	return Model, err
}
