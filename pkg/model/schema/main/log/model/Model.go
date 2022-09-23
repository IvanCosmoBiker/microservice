package model

import (
	"database/sql"
	"encoding/json"
	modelInterface "ephorservices/pkg/model/schema/interface/model"
	"log"

	pgx "github.com/jackc/pgx/v4"
)

var nameTable string = "log"
var nameSchema string = "main"

type LogModel struct {
	Id              int
	Address         sql.NullString
	Login           sql.NullString
	Date            sql.NullString
	Request_id      sql.NullString
	Request_uri     sql.NullString
	Request_data    sql.NullString
	Response        sql.NullString
	Runtime         sql.NullInt32
	Runtime_details sql.NullInt32
}

func Init() modelInterface.Model {
	return &LogModel{}
}

func (cm *LogModel) InitData(jsonData []byte) {
	json.Unmarshal(jsonData, &cm)
}

func (cm *LogModel) JsonSerialize() []byte {
	data, _ := json.Marshal(cm)
	return data
}

func (cm *LogModel) GetNameTable(account string) string {
	return nameTable
}

func (cm *LogModel) GetNameSchema(account int) string {
	return nameSchema
}

type ReturningStruct struct {
	Id              int
	Address         string
	Login           string
	Date            string
	Request_id      string
	Request_uri     string
	Request_data    string
	Response        string
	Runtime         int32
	Runtime_details int32
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
	Model := LogModel{}
	err := row.Scan(&Model.Id,
		&Model.Address,
		&Model.Login,
		&Model.Date,
		&Model.Request_id,
		&Model.Request_uri,
		&Model.Request_data,
		&Model.Response,
		&Model.Runtime,
		&Model.Runtime_details)
	if err != nil {
		log.Printf("%v\n", err)
		return &ReturningStruct{}, err
	}
	Result := MakeDataInSturct(&Model)
	return Result, nil
}

func MakeDataInSturct(modelData *LogModel) *ReturningStruct {
	stuctReturn := ReturningStruct{}
	stuctReturn.Id = modelData.Id
	stuctReturn.Address = modelData.Address.String
	stuctReturn.Login = modelData.Login.String
	stuctReturn.Date = modelData.Date.String
	stuctReturn.Request_id = modelData.Request_id.String
	stuctReturn.Request_uri = modelData.Request_uri.String
	stuctReturn.Request_data = modelData.Request_data.String
	stuctReturn.Response = modelData.Response.String
	stuctReturn.Runtime = modelData.Runtime.Int32
	stuctReturn.Runtime_details = modelData.Runtime_details.Int32
	return &stuctReturn
}
