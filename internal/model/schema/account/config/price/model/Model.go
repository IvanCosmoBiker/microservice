package model

import (
	"database/sql"
	"encoding/json"
	modelInterface "ephorservices/pkg/model/schema/interface/model"
	"log"

	pgx "github.com/jackc/pgx/v5"
)

var nameTable string = "config_product_price"
var nameSchema string = "account"

type ConfigProductPriceModel struct {
	Id            int
	Product_id    sql.NullInt32
	Price_list_id sql.NullInt32
	Value         sql.NullInt64
}

func Init() modelInterface.Model {
	return &ConfigProductPriceModel{}
}

func (aem *ConfigProductPriceModel) InitData(jsonData []byte) {
	json.Unmarshal(jsonData, &aem)
}

func (aem *ConfigProductPriceModel) JsonSerialize() []byte {
	data, _ := json.Marshal(aem)
	return data
}

func (aem *ConfigProductPriceModel) GetNameTable() string {
	return nameTable
}

func (aem *ConfigProductPriceModel) GetNameSchema(accountId string) string {
	return nameSchema + "" + accountId
}

type ReturningStruct struct {
	Id            int
	Product_id    int32
	Price_list_id int32
	Value         int64
}

func ScanModelRows(rows pgx.Rows) ([]*ReturningStruct, error) {
	defer rows.Close()
	Result := []*ReturningStruct{}
	for rows.Next() {
		resultStruct, err := ScanModelRow(rows)
		if err != nil {
			log.Printf("%v", err)
			continue
		}
		Result = append(Result, resultStruct)
	}
	return Result, nil
}

func ScanModelRow(row pgx.Row) (*ReturningStruct, error) {
	Model := ConfigProductPriceModel{}
	err := row.Scan(&Model.Id,
		&Model.Product_id,
		&Model.Price_list_id,
		&Model.Value)

	if err != nil {
		log.Printf("%v\n", err)
		return &ReturningStruct{}, err
	}
	Result := MakeDataInSturct(&Model)

	return Result, nil
}

func MakeDataInSturct(modelData *ConfigProductPriceModel) *ReturningStruct {
	stuctReturn := ReturningStruct{}
	stuctReturn.Id = modelData.Id
	stuctReturn.Product_id = modelData.Product_id.Int32
	stuctReturn.Price_list_id = modelData.Price_list_id.Int32
	stuctReturn.Value = modelData.Value.Int64
	return &stuctReturn
}
