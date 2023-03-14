package store

import (
	recipe_model "ephorservices/internal/model/schema/account/ingredient/recipe/model"
	model_interface "ephorservices/pkg/orm/model"
	store_parent "ephorservices/pkg/orm/store"
)

type StoreRecipe struct {
	*store_parent.Store
}

func New(accountNumber ...int) *StoreRecipe {
	AccountNumber := 0
	if len(accountNumber) != 0 {
		AccountNumber = accountNumber[0]
	}
	Model := recipe_model.New()
	store := &StoreRecipe{
		store_parent.New(Model),
	}
	store.Store.SetAccountNumber(AccountNumber)
	return store
}

func (sw *StoreRecipe) GetStructModel(model model_interface.Model) *recipe_model.RecipeModel {
	if model != nil {
		return model.(*recipe_model.RecipeModel)
	}
	return &recipe_model.RecipeModel{}
}

func (sw *StoreRecipe) AddByParams(params map[string]interface{}) (*recipe_model.RecipeModel, error) {
	model, err := sw.Store.AddByParams(params)
	Model := sw.GetStructModel(model)
	return Model, err
}

func (sw *StoreRecipe) SetByParams(params map[string]interface{}) (*recipe_model.RecipeModel, error) {
	model, err := sw.Store.SetByParams(params)
	Model := sw.GetStructModel(model)
	return Model, err
}

func (sw *StoreRecipe) GetOneById(id int) (*recipe_model.RecipeModel, error) {
	model, err := sw.Store.GetOneById(id)
	Model := sw.GetStructModel(model)
	return Model, err
}
