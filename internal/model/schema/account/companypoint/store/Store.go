package store

import (
	company_point_model "ephorservices/internal/model/schema/account/companypoint/model"
	model_interface "ephorservices/pkg/orm/model"
	request "ephorservices/pkg/orm/request"
	store_parent "ephorservices/pkg/orm/store"
)

type StoreCompanyPoint struct {
	*store_parent.Store
}

func New(accountNumber ...int) *StoreCompanyPoint {
	AccountNumber := 0
	if len(accountNumber) != 0 {
		AccountNumber = accountNumber[0]
	}
	Model := company_point_model.New()
	store := &StoreCompanyPoint{
		store_parent.New(Model),
	}
	store.Store.SetAccountNumber(AccountNumber)
	return store
}

func (scp *StoreCompanyPoint) GetStructModel(model model_interface.Model) *company_point_model.CompanyPointModel {
	if model != nil {
		return model.(*company_point_model.CompanyPointModel)
	}
	return &company_point_model.CompanyPointModel{}
}

func (scp *StoreCompanyPoint) AddByParams(params map[string]interface{}) (*company_point_model.CompanyPointModel, error) {
	model, err := scp.Store.AddByParams(params)
	Model := scp.GetStructModel(model)
	return Model, err
}

func (scp *StoreCompanyPoint) SetByParams(params map[string]interface{}) (*company_point_model.CompanyPointModel, error) {
	model, err := scp.Store.SetByParams(params)
	Model := scp.GetStructModel(model)
	return Model, err
}

func (scp *StoreCompanyPoint) GetOneById(id int) (*company_point_model.CompanyPointModel, error) {
	model, err := scp.Store.GetOneById(id)
	Model := scp.GetStructModel(model)
	return Model, err
}

func (scp *StoreCompanyPoint) GetOneBy(req *request.Request) (*company_point_model.CompanyPointModel, error) {
	model, err := scp.Store.GetOneBy(req)
	Model := scp.GetStructModel(model)
	return Model, err
}
