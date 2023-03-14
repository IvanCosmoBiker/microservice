package store

import (
	wareflow_ingredient_model "ephorservices/internal/model/schema/account/wareflow/ingredient/model"
	model_interface "ephorservices/pkg/orm/model"
	request "ephorservices/pkg/orm/request"
	store_parent "ephorservices/pkg/orm/store"
)

var (
	Warning_Low_Product = 0x01
)

type StoreWareFlowIngredient struct {
	*store_parent.Store
}

func New(accountNumber ...int) *StoreWareFlowIngredient {
	AccountNumber := 0
	if len(accountNumber) != 0 {
		AccountNumber = accountNumber[0]
	}
	Model := wareflow_ingredient_model.New()
	store := &StoreWareFlowIngredient{
		store_parent.New(Model),
	}
	store.Store.SetAccountNumber(AccountNumber)
	return store
}

func (swfi *StoreWareFlowIngredient) GetStructModel(model model_interface.Model) *wareflow_ingredient_model.WareFlowIngredientModel {
	if model != nil {
		return model.(*wareflow_ingredient_model.WareFlowIngredientModel)
	}
	return &wareflow_ingredient_model.WareFlowIngredientModel{}
}

func (swfi *StoreWareFlowIngredient) AddByParams(params map[string]interface{}) (*wareflow_ingredient_model.WareFlowIngredientModel, error) {
	model, err := swfi.Store.AddByParams(params)
	Model := swfi.GetStructModel(model)
	return Model, err
}

func (swfi *StoreWareFlowIngredient) SetByParams(params map[string]interface{}) (*wareflow_ingredient_model.WareFlowIngredientModel, error) {
	model, err := swfi.Store.SetByParams(params)
	Model := swfi.GetStructModel(model)
	return Model, err
}

func (swfi *StoreWareFlowIngredient) GetOneById(id int) (*wareflow_ingredient_model.WareFlowIngredientModel, error) {
	model, err := swfi.Store.GetOneById(id)
	Model := swfi.GetStructModel(model)
	return Model, err
}

func (swfi *StoreWareFlowIngredient) GetOneBy(req *request.Request) (*wareflow_ingredient_model.WareFlowIngredientModel, error) {
	model, err := swfi.Store.GetOneBy(req)
	Model := swfi.GetStructModel(model)
	return Model, err
}

func (swfi *StoreWareFlowIngredient) Deduction(wareflow_id, ingredient_id int, quantity int32) (int, error) {
	req := request.New()
	req.AddFilterParam("wareflow_id", req.Operator.OperatorEqual, true, wareflow_id)
	req.AddFilterParam("ingredient_id", req.Operator.OperatorEqual, true, ingredient_id)
	model, err := swfi.GetOneBy(req)
	if err != nil {
		return 0, err
	}
	if model.Id == 0 {
		return 0, store_parent.Error_Empty_Data
	}
	model.Qt_counter.Int32 -= quantity
	if model.Qt_counter.Int32 <= int32(0) {
		model.Qt_counter.Int32 = int32(0)
		_, err = swfi.Store.Set(model)
		return Warning_Low_Product, err
	}
	_, err = swfi.Store.Set(model)
	return 0, err
}
