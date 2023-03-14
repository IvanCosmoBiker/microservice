package store

import (
	modem_model "ephorservices/internal/model/schema/main/modem/model"
	model_interface "ephorservices/pkg/orm/model"
	request "ephorservices/pkg/orm/request"
	store_parent "ephorservices/pkg/orm/store"
)

type StoreModem struct {
	*store_parent.Store
}

func New(accountNumber ...int) *StoreModem {
	AccountNumber := 0
	if len(accountNumber) != 0 {
		AccountNumber = accountNumber[0]
	}
	Model := modem_model.New()
	store := &StoreModem{
		store_parent.New(Model),
	}
	store.Store.SetAccountNumber(AccountNumber)
	return store
}

func (sm *StoreModem) GetStructModel(model model_interface.Model) *modem_model.ModemModel {
	if model != nil {
		return model.(*modem_model.ModemModel)
	}
	return &modem_model.ModemModel{}
}

func (sm *StoreModem) AddByParams(params map[string]interface{}) (*modem_model.ModemModel, error) {
	model, err := sm.Store.AddByParams(params)
	Model := sm.GetStructModel(model)
	return Model, err
}

func (sm *StoreModem) SetByParams(params map[string]interface{}) (*modem_model.ModemModel, error) {
	model, err := sm.Store.SetByParams(params)
	Model := sm.GetStructModel(model)
	return Model, err
}

func (sm *StoreModem) GetOneById(id int) (*modem_model.ModemModel, error) {
	model, err := sm.Store.GetOneById(id)
	Model := sm.GetStructModel(model)
	return Model, err
}

func (sm *StoreModem) GetOneBy(req *request.Request) (*modem_model.ModemModel, error) {
	model, err := sm.Store.GetOneBy(req)
	Model := sm.GetStructModel(model)
	return Model, err
}
