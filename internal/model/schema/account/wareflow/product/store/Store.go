package ware

import (
	wareflow_product_model "ephorservices/internal/model/schema/account/wareflow/product/model"
	model_interface "ephorservices/pkg/orm/model"
	request "ephorservices/pkg/orm/request"
	store_parent "ephorservices/pkg/orm/store"
)

var (
	Warning_Low_Product = 0x01
)

type StoreWareFlowProduct struct {
	*store_parent.Store
}

func New(accountNumber ...int) *StoreWareFlowProduct {
	AccountNumber := 0
	if len(accountNumber) != 0 {
		AccountNumber = accountNumber[0]
	}
	Model := wareflow_product_model.New()
	store := &StoreWareFlowProduct{
		store_parent.New(Model),
	}
	store.Store.SetAccountNumber(AccountNumber)
	return store
}

func (swf *StoreWareFlowProduct) GetStructModel(model model_interface.Model) *wareflow_product_model.WareFlowProductModel {
	if model != nil {
		return model.(*wareflow_product_model.WareFlowProductModel)
	}
	return &wareflow_product_model.WareFlowProductModel{}
}

func (swf *StoreWareFlowProduct) AddByParams(params map[string]interface{}) (*wareflow_product_model.WareFlowProductModel, error) {
	model, err := swf.Store.AddByParams(params)
	Model := swf.GetStructModel(model)
	return Model, err
}

func (swf *StoreWareFlowProduct) SetByParams(params map[string]interface{}) (*wareflow_product_model.WareFlowProductModel, error) {
	model, err := swf.Store.SetByParams(params)
	Model := swf.GetStructModel(model)
	return Model, err
}

func (swf *StoreWareFlowProduct) GetOneById(id int) (*wareflow_product_model.WareFlowProductModel, error) {
	model, err := swf.Store.GetOneById(id)
	Model := swf.GetStructModel(model)
	return Model, err
}

func (swf *StoreWareFlowProduct) GetOneBy(req *request.Request) (*wareflow_product_model.WareFlowProductModel, error) {
	model, err := swf.Store.GetOneBy(req)
	Model := swf.GetStructModel(model)
	return Model, err
}

func (swf *StoreWareFlowProduct) Deduction(wareFlowId int, select_id string, ware_id int, quantity int32) (int, error) {
	req := request.New()
	req.AddFilterParam("select_id", req.Operator.OperatorEqual, true, select_id)
	req.AddFilterParam("wareflow_id", req.Operator.OperatorEqual, true, wareFlowId)
	req.AddFilterParam("ware_id", req.Operator.OperatorEqual, true, ware_id)
	model, err := swf.GetOneBy(req)
	if err != nil {
		return 0, err
	}
	if model.Id == 0 {
		return 0, store_parent.Error_Empty_Data
	}
	model.Qt_counter.Int32 -= quantity
	if model.Qt_counter.Int32 <= int32(0) {
		model.Qt_counter.Int32 = int32(0)
		_, err = swf.Store.Set(model)
		return Warning_Low_Product, err
	}
	_, err = swf.Store.Set(model)
	return 0, err
}
