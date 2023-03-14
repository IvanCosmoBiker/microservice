package store

import (
	warehouse_model "ephorservices/pkg/model/schema/account/warehouse/model"
	model_interface "ephorservices/pkg/orm/model"
	request "ephorservices/pkg/orm/request"
	store_parent "ephorservices/pkg/orm/store"
)

var (
	StateObsolete = 0x00
	StateActive   = 0x01
)

type WareHouseStore struct {
	*store_parent.Store
}

func New(accountNumber ...int) *WareHouseStore {
	AccountNumber := 0
	if len(accountNumber) != 0 {
		AccountNumber = accountNumber[0]
	}
	Model := warehouse_model.New()
	store := &WareHouseStore{
		store_parent.New(Model),
	}
	store.Store.SetAccountNumber(AccountNumber)
	return store
}

func (sl *WareHouseStore) GetStructModel(model model_interface.Model) *warehouse_model.WareHouseModel {
	if model != nil {
		return model.(*warehouse_model.WareHouseModel)
	}
	return &warehouse_model.WareHouseModel{}
}

func (sl *WareHouseStore) AddByParams(params map[string]interface{}) (*warehouse_model.WareHouseModel, error) {
	model, err := sl.Store.AddByParams(params)
	Model := sl.GetStructModel(model)
	return Model, err
}

func (sl *WareHouseStore) SetByParams(params map[string]interface{}) (*warehouse_model.WareHouseModel, error) {
	model, err := sl.Store.SetByParams(params)
	Model := sl.GetStructModel(model)
	return Model, err
}

func (sl *WareHouseStore) GetOneById(id int) (*warehouse_model.WareHouseModel, error) {
	model, err := sl.Store.GetOneById(id)
	Model := sl.GetStructModel(model)
	return Model, err
}

func (sl *WareHouseStore) GetOneBy(req *request.Request) (*warehouse_model.WareHouseModel, error) {
	model, err := sl.Store.GetOneBy(req)
	Model := sl.GetStructModel(model)
	return Model, err
}
