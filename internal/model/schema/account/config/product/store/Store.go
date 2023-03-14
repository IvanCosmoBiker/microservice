package store

import (
	config_product_model "ephorservices/internal/model/schema/account/config/product/model"
	model_interface "ephorservices/pkg/orm/model"
	request "ephorservices/pkg/orm/request"
	store_parent "ephorservices/pkg/orm/store"
)

type StoreConfigProduct struct {
	*store_parent.Store
}

func New(accountNumber ...int) *StoreConfigProduct {
	AccountNumber := 0
	if len(accountNumber) != 0 {
		AccountNumber = accountNumber[0]
	}
	Model := config_product_model.New()
	store := &StoreConfigProduct{
		store_parent.New(Model),
	}
	store.Store.SetAccountNumber(AccountNumber)
	return store
}

func (scp *StoreConfigProduct) GetStructModel(model model_interface.Model) *config_product_model.ConfigProductModel {
	if model != nil {
		return model.(*config_product_model.ConfigProductModel)
	}
	return &config_product_model.ConfigProductModel{}
}

func (scp *StoreConfigProduct) AddByParams(params map[string]interface{}) (*config_product_model.ConfigProductModel, error) {
	model, err := scp.Store.AddByParams(params)
	Model := scp.GetStructModel(model)
	return Model, err
}

func (scp *StoreConfigProduct) SetByParams(params map[string]interface{}) (*config_product_model.ConfigProductModel, error) {
	model, err := scp.Store.SetByParams(params)
	Model := scp.GetStructModel(model)
	return Model, err
}

func (scp *StoreConfigProduct) Set(m model_interface.Model) (*config_product_model.ConfigProductModel, error) {
	model, err := scp.Store.Set(m)
	Model := scp.GetStructModel(model)
	return Model, err
}

func (scp *StoreConfigProduct) GetOneById(id int) (*config_product_model.ConfigProductModel, error) {
	model, err := scp.Store.GetOneById(id)
	Model := scp.GetStructModel(model)
	return Model, err
}

func (scp *StoreConfigProduct) GetOneBy(req *request.Request) (*config_product_model.ConfigProductModel, error) {
	model, err := scp.Store.GetOneBy(req)
	Model := scp.GetStructModel(model)
	return Model, err
}
