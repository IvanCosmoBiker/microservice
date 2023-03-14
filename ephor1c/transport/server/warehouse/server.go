package server

import (
	"encoding/json"
	warehouse_service "ephorservices/ephor1c/service/intergation_1c/warehouse"
	request "ephorservices/ephor1c/transport/request"
	transportHttp "ephorservices/pkg/transportprotocol/http/v1"
	"net/http"
)

type ServerHttpWareHouse struct {
	Address          map[string]string
	Method           []string
	WareHouseService *warehouse_service.ServiceWareHouse
}

func New(host string) *ServerHttpWareHouse {
	api := &ServerHttpWareHouse{
		WareHouseService: warehouse_service.New(),
	}
	api.InitAdress()
	return api
}

func (she *ServerHttpWareHouse) InitAdress() {
	address := make(map[string]string)
	address["WarehouseGet"] = "/WarehouseGet"
	address["WarehouseCreate"] = "/WarehouseCreate"
	address["WarehouseDelete"] = "/WarehouseDelete"
	address["WarehouseUpdate"] = "/WarehouseUpdate"
	address["WarehouseProductGet"] = "/WarehouseProductGet"
}

func (she *ServerHttpWareHouse) JSON(w http.ResponseWriter, data interface{}) {
	json.NewEncoder(w).Encode(data)
}

func (she *ServerHttpWareHouse) handlerSyncWareHouse(w http.ResponseWriter, req *http.Request) {
	requestWareHouse := request.New("http", params["host"].(string), request.Resource_Catalog, "Склады")
	requestWareHouse.SetBasicAuth(params["login"].(string), params["password"].(string))
	result := she.WareHouseService.GetWarehouse(requestWareHouse)
	she.JSON(w, result)
}

func (she *ServerHttpWareHouse) handlerSyncWareHouseProduct(w http.ResponseWriter, req *http.Request) {

}

func (she *ServerHttpWareHouse) InitUrl(HttpManager *transportHttp.ServerHttp) {
	HttpManager.SetHandlerListener(she.Address["WarehouseGet"], she.handlerSyncWareHouse)
	HttpManager.SetHandlerListener(she.Address["WarehouseProductGet"], she.handlerSyncWareHouseProduct)
}
