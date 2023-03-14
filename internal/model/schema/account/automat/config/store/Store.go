package store

import (
	automat_config_model "ephorservices/internal/model/schema/account/automat/config/model"
	model_interface "ephorservices/pkg/orm/model"
	request "ephorservices/pkg/orm/request"
	store_parent "ephorservices/pkg/orm/store"
	"fmt"
)

type StoreAutomatConfig struct {
	*store_parent.Store
}

func New(accountNumber ...int) *StoreAutomatConfig {
	AccountNumber := 0
	if len(accountNumber) != 0 {
		AccountNumber = accountNumber[0]
	}
	Model := automat_config_model.New()
	store := &StoreAutomatConfig{
		store_parent.New(Model),
	}
	store.Store.SetAccountNumber(AccountNumber)
	return store
}

func (sac *StoreAutomatConfig) GetStructModel(model model_interface.Model) *automat_config_model.AutomatConfigModel {
	if model != nil {
		return model.(*automat_config_model.AutomatConfigModel)
	}
	return &automat_config_model.AutomatConfigModel{}
}

func (sac *StoreAutomatConfig) AddByParams(params map[string]interface{}) (*automat_config_model.AutomatConfigModel, error) {
	model, err := sac.Store.AddByParams(params)
	Model := sac.GetStructModel(model)
	return Model, err
}

func (sac *StoreAutomatConfig) SetByParams(params map[string]interface{}) (*automat_config_model.AutomatConfigModel, error) {
	model, err := sac.Store.SetByParams(params)
	Model := sac.GetStructModel(model)
	return Model, err
}

func (sac *StoreAutomatConfig) GetOneById(id int) (*automat_config_model.AutomatConfigModel, error) {
	model, err := sac.Store.GetOneById(id)
	Model := sac.GetStructModel(model)
	return Model, err
}

func (sac *StoreAutomatConfig) GetOneBy(req *request.Request) (*automat_config_model.AutomatConfigModel, error) {
	model, err := sac.Store.GetOneBy(req)
	fmt.Printf("%+v___!!!!\n", model)
	fmt.Printf("%+v___!!!!\n", err)
	Model := sac.GetStructModel(model)
	fmt.Printf("%+v___!!!!\n", Model)
	return Model, err
}
