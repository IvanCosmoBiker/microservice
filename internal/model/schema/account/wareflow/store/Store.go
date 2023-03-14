package store

import (
	wareflow_model "ephorservices/internal/model/schema/account/wareflow/model"
	model_interface "ephorservices/pkg/orm/model"
	request "ephorservices/pkg/orm/request"
	store_parent "ephorservices/pkg/orm/store"
)

const (
	State_Start    = 0
	State_Complete = 1
)

const (
	Type_Load   = 0
	Type_UnLoad = 1
)

type StoreWareFlow struct {
	*store_parent.Store
}

func New(accountNumber ...int) *StoreWareFlow {
	AccountNumber := 0
	if len(accountNumber) != 0 {
		AccountNumber = accountNumber[0]
	}
	Model := wareflow_model.New()
	store := &StoreWareFlow{
		store_parent.New(Model),
	}
	store.Store.SetAccountNumber(AccountNumber)
	return store
}

func (swf *StoreWareFlow) GetStructModel(model model_interface.Model) *wareflow_model.WareFlowModel {
	if model != nil {
		return model.(*wareflow_model.WareFlowModel)
	}
	return &wareflow_model.WareFlowModel{}
}

func (swf *StoreWareFlow) AddByParams(params map[string]interface{}) (*wareflow_model.WareFlowModel, error) {
	model, err := swf.Store.AddByParams(params)
	Model := swf.GetStructModel(model)
	return Model, err
}

func (swf *StoreWareFlow) SetByParams(params map[string]interface{}) (*wareflow_model.WareFlowModel, error) {
	model, err := swf.Store.SetByParams(params)
	Model := swf.GetStructModel(model)
	return Model, err
}

func (swf *StoreWareFlow) GetOneById(id int) (*wareflow_model.WareFlowModel, error) {
	model, err := swf.Store.GetOneById(id)
	Model := swf.GetStructModel(model)
	return Model, err
}

func (swf *StoreWareFlow) GetOneBy(req *request.Request) (*wareflow_model.WareFlowModel, error) {
	model, err := swf.Store.GetOneBy(req)
	Model := swf.GetStructModel(model)
	return Model, err
}

func (swf *StoreWareFlow) GetLatest(automatId int) (*wareflow_model.WareFlowModel, error) {
	req := request.New()
	req.AddFilterParam("automat_id", req.Operator.OperatorEqual, true, automatId)
	req.SetSorter("date", req.Operator.Desc)
	return swf.GetOneBy(req)
}
