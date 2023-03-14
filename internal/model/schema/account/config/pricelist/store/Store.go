package store

import (
	config_price_list_model "ephorservices/internal/model/schema/account/config/pricelist/model"
	model_interface "ephorservices/pkg/orm/model"
	request "ephorservices/pkg/orm/request"
	store_parent "ephorservices/pkg/orm/store"
)

const Price_List0 int32 = 0x0000 // basic price list used to be CA0
const Price_List1 int32 = 0x0001 // price list used to be DA1
const Price_List2 int32 = 0x0002 // price list. New price list, which will use to discount
type StoreConfigPriceList struct {
	*store_parent.Store
}

func New(accountNumber ...int) *StoreConfigPriceList {
	AccountNumber := 0
	if len(accountNumber) != 0 {
		AccountNumber = accountNumber[0]
	}
	Model := config_price_list_model.New()
	store := &StoreConfigPriceList{
		store_parent.New(Model),
	}
	store.Store.SetAccountNumber(AccountNumber)
	return store
}

func (scpl *StoreConfigPriceList) GetStructModel(model model_interface.Model) *config_price_list_model.ConfigPriceListModel {
	if model != nil {
		return model.(*config_price_list_model.ConfigPriceListModel)
	}
	return &config_price_list_model.ConfigPriceListModel{}
}

func (scpl *StoreConfigPriceList) AddByParams(params map[string]interface{}) (*config_price_list_model.ConfigPriceListModel, error) {
	model, err := scpl.Store.AddByParams(params)
	Model := scpl.GetStructModel(model)
	return Model, err
}

func (scpl *StoreConfigPriceList) SetByParams(params map[string]interface{}) (*config_price_list_model.ConfigPriceListModel, error) {
	model, err := scpl.Store.SetByParams(params)
	Model := scpl.GetStructModel(model)
	return Model, err
}

func (scpl *StoreConfigPriceList) Set(m model_interface.Model) (*config_price_list_model.ConfigPriceListModel, error) {
	model, err := scpl.Store.Set(m)
	Model := scpl.GetStructModel(model)
	return Model, err
}

func (scpl *StoreConfigPriceList) GetOneById(id int) (*config_price_list_model.ConfigPriceListModel, error) {
	model, err := scpl.Store.GetOneById(id)
	Model := scpl.GetStructModel(model)
	return Model, err
}

func (scpl *StoreConfigPriceList) GetOneBy(req *request.Request) (*config_price_list_model.ConfigPriceListModel, error) {
	model, err := scpl.Store.GetOneBy(req)
	Model := scpl.GetStructModel(model)
	return Model, err
}

func (scpl *StoreConfigPriceList) GetPaymentDevice(number int32) string {
	switch number {
	case Price_List0:
		return "CA"
	case Price_List1:
		return "DA"
	}
	return "CA"
}
