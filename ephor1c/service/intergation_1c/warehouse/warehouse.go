package warehouse

import (
	httpClient "ephorservices/ephor1c/service/intergation_1c"
	catalog_warehouse "ephorservices/ephor1c/service/intergation_1c/model/catalog/warehouse"
	request "ephorservices/ephor1c/transport/request"
)

type ServiceWareHouse struct {
	Http *httpClient.Http
}

func New() *ServiceWareHouse {
	return &ServiceWareHouse{
		Http: httpClient.New(),
	}
}

func (swh *ServiceWareHouse) GetWarehouse(req *request.Request1c) []*catalog_warehouse.Warehouse {
	result := make([]*catalog_warehouse.Warehouse, 0, 1)
	httpClient.New()
	return result
}

func (swh *ServiceWareHouse) GetWareHouseProduct(req *request.Request1c) {

}
