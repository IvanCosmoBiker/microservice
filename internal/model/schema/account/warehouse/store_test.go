package warehouse

import (
	"context"

	warehouse_store "ephorservices/pkg/model/schema/account/warehouse/store"
	connectionPostgresql "ephorservices/pkg/orm/db"
	"ephorservices/pkg/orm/request"
	"fmt"
	"testing"
)

var ConnectionDb *connectionPostgresql.Manager

func init() {
	var err error
	ConnectionDb, err = connectionPostgresql.Init("postgres", "123", "127.0.0.1", "local", uint16(5432), uint16(10), uint16(10), uint16(10), false, true, context.Background())
	if err != nil {
		panic(err)
	}
}

func TestCreateWareHouse(t *testing.T) {
	store := warehouse_store.New(1)
	params := make(map[string]interface{})
	params["name"] = "test1"
	params["date_sync"] = "2022-09-11 00:00:00"
	params["state"] = warehouse_store.StateActive
	modelWareHouse, err := store.AddByParams(params)
	if err != nil {
		fmt.Println(err.Error())
	}
	req := request.New()
	req.AddFilterParam("id", request.OperatorEqual, true, modelWareHouse.Id)
	models, _ := store.GetOneBy(req)
	fmt.Printf("%+v", models[0])
	store.DeleteById(models[0].Id)
}
