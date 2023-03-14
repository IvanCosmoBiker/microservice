package store

import (
	recipe_ingredient_model "ephorservices/internal/model/schema/account/ingredient/recipe/ingredient/model"
	model_interface "ephorservices/pkg/orm/model"
	request "ephorservices/pkg/orm/request"
	store_parent "ephorservices/pkg/orm/store"
)

type StoreRecipeIngredient struct {
	*store_parent.Store
}

func New(accountNumber ...int) *StoreRecipeIngredient {
	AccountNumber := 0
	if len(accountNumber) != 0 {
		AccountNumber = accountNumber[0]
	}
	Model := recipe_ingredient_model.New()
	store := &StoreRecipeIngredient{
		store_parent.New(Model),
	}
	store.Store.SetAccountNumber(AccountNumber)
	return store
}

func (sri *StoreRecipeIngredient) GetStructModel(model model_interface.Model) *recipe_ingredient_model.RecipeIngredientModel {
	return model.(*recipe_ingredient_model.RecipeIngredientModel)
}

func (sri *StoreRecipeIngredient) Get(req *request.Request) (Models []*recipe_ingredient_model.RecipeIngredientModel, err error) {
	models, err := sri.Store.Get(req)
	if len(models) > 0 {
		Models = make([]*recipe_ingredient_model.RecipeIngredientModel, 0, len(models))
		for _, model := range models {
			Models = append(Models, sri.GetStructModel(model))
		}
	}
	return
}

func (sri *StoreRecipeIngredient) AddByParams(params map[string]interface{}) (*recipe_ingredient_model.RecipeIngredientModel, error) {
	model, err := sri.Store.AddByParams(params)
	Model := sri.GetStructModel(model)
	return Model, err
}

func (sri *StoreRecipeIngredient) SetByParams(params map[string]interface{}) (*recipe_ingredient_model.RecipeIngredientModel, error) {
	model, err := sri.Store.SetByParams(params)
	Model := sri.GetStructModel(model)
	return Model, err
}

func (sri *StoreRecipeIngredient) GetOneById(id int) (*recipe_ingredient_model.RecipeIngredientModel, error) {
	model, err := sri.Store.GetOneById(id)
	Model := sri.GetStructModel(model)
	return Model, err
}
