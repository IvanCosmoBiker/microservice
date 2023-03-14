package model

import (
	"database/sql"
	"encoding/json"
	modelInterface "ephorservices/pkg/model/schema/interface/model"
	"log"

	pgx "github.com/jackc/pgx/v5"
)

const (
	Category_Error   = 1
	Category_Warning = 2
	Category_Info    = 3
)

const (
	StateInactive = 0
	StateActive   = 1
)

var nameTable string = "automat_device_event"
var nameSchema string = "account"

type AutomatDeviceEventModel struct {
	Id        int
	Device_id sql.NullInt32
	Category  sql.NullInt32
	Type      sql.NullInt32
	Date      sql.NullString
	Value     sql.NullString
	Count     sql.NullInt32
	State     sql.NullInt32
}

func Init() modelInterface.Model {
	return &AutomatDeviceEventModel{}
}

func (aem *AutomatDeviceEventModel) InitData(jsonData []byte) {
	json.Unmarshal(jsonData, &aem)
}

func (aem *AutomatDeviceEventModel) JsonSerialize() []byte {
	data, _ := json.Marshal(aem)
	return data
}

func (aem *AutomatDeviceEventModel) GetNameTable() string {
	return nameTable
}

func (aem *AutomatDeviceEventModel) GetNameSchema(accountId string) string {
	return nameSchema + "" + accountId
}

type ReturningStruct struct {
	Id        int
	Device_id int32
	Category  int32
	Type      int32
	Date      string
	Value     string
	Count     int32
	State     int32
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
	Model := AutomatDeviceEventModel{}
	err := row.Scan(&Model.Id,
		&Model.Device_id,
		&Model.Category,
		&Model.Type,
		&Model.Date,
		&Model.Value,
		&Model.Count,
		&Model.State)
	if err != nil {
		log.Printf("%v\n", err)
		return &ReturningStruct{}, err
	}
	Result := MakeDataInSturct(&Model)

	return Result, nil
}

func MakeDataInSturct(modelData *AutomatDeviceEventModel) *ReturningStruct {
	stuctReturn := ReturningStruct{}
	stuctReturn.Id = modelData.Id
	stuctReturn.Device_id = modelData.Device_id.Int32
	stuctReturn.Category = modelData.Category.Int32
	stuctReturn.Type = modelData.Type.Int32
	stuctReturn.Date = modelData.Date.String
	stuctReturn.Value = modelData.Value.String
	stuctReturn.Count = modelData.Count.Int32
	stuctReturn.State = modelData.State.Int32
	return &stuctReturn
}
