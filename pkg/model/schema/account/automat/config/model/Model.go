package model

import (
	"database/sql"
	"encoding/json"
	modelInterface "ephorservices/pkg/model/schema/interface/model"
	"log"

	pgx "github.com/jackc/pgx/v4"
)

var nameTable string = "automat_config"
var nameSchema string = "account"

type AutomatConfigModel struct {
	Id         int
	Automat_id sql.NullInt32
	Config_id  sql.NullInt32
	From_date  sql.NullString
	To_date    sql.NullString
}

func Init() modelInterface.Model {
	return &AutomatConfigModel{}
}

func (aem *AutomatConfigModel) InitData(jsonData []byte) {
	json.Unmarshal(jsonData, &aem)
}

func (aem *AutomatConfigModel) JsonSerialize() []byte {
	data, _ := json.Marshal(aem)
	return data
}

func (aem *AutomatConfigModel) GetNameTable() string {
	return nameTable
}

func (aem *AutomatConfigModel) GetNameSchema(accountId string) string {
	return nameSchema + "" + accountId
}

type ReturningStruct struct {
	Id         int
	Automat_id int32
	Config_id  int32
	From_date  string
	To_date    string
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
	Model := AutomatConfigModel{}
	err := row.Scan(&Model.Id,
		&Model.Automat_id,
		&Model.Config_id,
		&Model.From_date,
		&Model.To_date)

	if err != nil {
		log.Printf("%v\n", err)
		return &ReturningStruct{}, err
	}
	Result := MakeDataInSturct(&Model)

	return Result, nil
}

func MakeDataInSturct(modelData *AutomatConfigModel) *ReturningStruct {
	stuctReturn := ReturningStruct{}
	stuctReturn.Id = modelData.Id
	stuctReturn.Automat_id = modelData.Automat_id.Int32
	stuctReturn.Config_id = modelData.Config_id.Int32
	stuctReturn.From_date = modelData.From_date.String
	stuctReturn.To_date = modelData.To_date.String
	return &stuctReturn
}
