package store

import (
	ingredient_model "ephorservices/internal/model/schema/account/ingredient/model"
	model_interface "ephorservices/pkg/orm/model"
	request "ephorservices/pkg/orm/request"
	store_parent "ephorservices/pkg/orm/store"
)

type StoreIngredient struct {
	*store_parent.Store
}

func New(accountNumber ...int) *StoreIngredient {
	AccountNumber := 0
	if len(accountNumber) != 0 {
		AccountNumber = accountNumber[0]
	}
	Model := ingredient_model.New()
	store := &StoreIngredient{
		store_parent.New(Model),
	}
	store.Store.SetAccountNumber(AccountNumber)
	return store
}

func (si *StoreIngredient) GetStructModel(model model_interface.Model) *ingredient_model.IngredientModel {
	if model != nil {
		return model.(*ingredient_model.IngredientModel)
	}
	return &ingredient_model.IngredientModel{}
}

func (si *StoreIngredient) AddByParams(params map[string]interface{}) (*ingredient_model.IngredientModel, error) {
	model, err := si.Store.AddByParams(params)
	Model := si.GetStructModel(model)
	return Model, err
}

func (si *StoreIngredient) SetByParams(params map[string]interface{}) (*ingredient_model.IngredientModel, error) {
	model, err := si.Store.SetByParams(params)
	Model := si.GetStructModel(model)
	return Model, err
}

func (si *StoreIngredient) GetOneById(id int) (*ingredient_model.IngredientModel, error) {
	model, err := si.Store.GetOneById(id)
	Model := si.GetStructModel(model)
	return Model, err
}

func (si *StoreIngredient) GetOneBy(req *request.Request) (*ingredient_model.IngredientModel, error) {
	model, err := si.Store.GetOneBy(req)
	Model := si.GetStructModel(model)
	return Model, err
}
