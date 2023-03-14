package barcode

import (
	ware_barcode_model "ephorservices/internal/model/schema/account/ware/barcode/model"
	model_interface "ephorservices/pkg/orm/model"
	request "ephorservices/pkg/orm/request"
	store_parent "ephorservices/pkg/orm/store"
)

type StoreWareBarcode struct {
	*store_parent.Store
}

func New(accountNumber ...int) *StoreWareBarcode {
	AccountNumber := 0
	if len(accountNumber) != 0 {
		AccountNumber = accountNumber[0]
	}
	Model := ware_barcode_model.New()
	store := &StoreWareBarcode{
		store_parent.New(Model),
	}
	store.Store.SetAccountNumber(AccountNumber)
	return store
}

func (swb *StoreWareBarcode) GetStructModel(model model_interface.Model) *ware_barcode_model.WareBarcodeModel {
	if model != nil {
		return model.(*ware_barcode_model.WareBarcodeModel)
	}
	return &ware_barcode_model.WareBarcodeModel{}
}

func (swb *StoreWareBarcode) AddByParams(params map[string]interface{}) (*ware_barcode_model.WareBarcodeModel, error) {
	model, err := swb.Store.AddByParams(params)
	Model := swb.GetStructModel(model)
	return Model, err
}

func (swb *StoreWareBarcode) SetByParams(params map[string]interface{}) (*ware_barcode_model.WareBarcodeModel, error) {
	model, err := swb.Store.SetByParams(params)
	Model := swb.GetStructModel(model)
	return Model, err
}

func (swb *StoreWareBarcode) GetOneById(id int) (*ware_barcode_model.WareBarcodeModel, error) {
	model, err := swb.Store.GetOneById(id)
	Model := swb.GetStructModel(model)
	return Model, err
}

func (swb *StoreWareBarcode) GetOneBy(req *request.Request) (*ware_barcode_model.WareBarcodeModel, error) {
	model, err := swb.Store.GetOneBy(req)
	Model := swb.GetStructModel(model)
	return Model, err
}
