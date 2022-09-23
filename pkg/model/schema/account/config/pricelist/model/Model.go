package model

import (
	"database/sql"
	"encoding/json"
	modelInterface "ephorservices/pkg/model/schema/interface/model"
	"log"

	pgx "github.com/jackc/pgx/v4"
)

var nameTable string = "config_price_list"
var nameSchema string = "account"

type ConfigPriceListModel struct {
	Id             int
	Config_id      sql.NullInt32
	Payment_device sql.NullString
	Number         sql.NullInt32
	Work_week      sql.NullInt32
	Work_time      sql.NullString
	Work_interval  sql.NullInt32
}

func Init() modelInterface.Model {
	return &ConfigPriceListModel{}
}

func (aem *ConfigPriceListModel) InitData(jsonData []byte) {
	json.Unmarshal(jsonData, &aem)
}

func (aem *ConfigPriceListModel) JsonSerialize() []byte {
	data, _ := json.Marshal(aem)
	return data
}

func (aem *ConfigPriceListModel) GetNameTable() string {
	return nameTable
}

func (aem *ConfigPriceListModel) GetNameSchema(accountId string) string {
	return nameSchema + "" + accountId
}

type ReturningStruct struct {
	Id             int
	Config_id      int32
	Payment_device string
	Number         int32
	Work_week      int32
	Work_time      string
	Work_interval  int32
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
	Model := ConfigPriceListModel{}
	err := row.Scan(&Model.Id,
		&Model.Config_id,
		&Model.Payment_device,
		&Model.Number,
		&Model.Work_week,
		&Model.Work_time,
		&Model.Work_interval)

	if err != nil {
		log.Printf("%v\n", err)
		return &ReturningStruct{}, err
	}
	Result := MakeDataInSturct(&Model)

	return Result, nil
}

func MakeDataInSturct(modelData *ConfigPriceListModel) *ReturningStruct {
	stuctReturn := ReturningStruct{}
	stuctReturn.Id = modelData.Id
	stuctReturn.Config_id = modelData.Config_id.Int32
	stuctReturn.Payment_device = modelData.Payment_device.String
	stuctReturn.Number = modelData.Number.Int32
	stuctReturn.Work_week = modelData.Work_week.Int32
	stuctReturn.Work_time = modelData.Work_time.String
	stuctReturn.Work_interval = modelData.Work_interval.Int32
	return &stuctReturn
}
