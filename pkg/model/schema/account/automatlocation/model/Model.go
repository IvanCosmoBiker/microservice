package model

import (
	"database/sql"
	"encoding/json"
	modelInterface "ephorservices/pkg/model/schema/interface/model"
	"log"

	pgx "github.com/jackc/pgx/v4"
)

var nameTable string = "automat_location"
var nameSchema string = "account"

type AutomatLocationModel struct {
	Id               int
	Automat_id       sql.NullInt32
	Company_point_id sql.NullInt32
	From_date        sql.NullString
	To_date          sql.NullString
}

func Init() modelInterface.Model {
	return &AutomatLocationModel{}
}

func (am *AutomatLocationModel) InitData(jsonData []byte) {
	json.Unmarshal(jsonData, &am)
}

func (am *AutomatLocationModel) JsonSerialize() []byte {
	data, _ := json.Marshal(am)
	return data
}

func (am *AutomatLocationModel) GetNameTable() string {
	return nameTable
}

func (am *AutomatLocationModel) GetNameSchema(accountId string) string {
	return nameSchema + "" + accountId
}

type ReturningStruct struct {
	Id               int
	Automat_id       int32
	Company_point_id int32
	From_date        string
	To_date          string
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
	Model := AutomatLocationModel{}
	err := row.Scan(&Model.Id,
		&Model.Automat_id,
		&Model.Company_point_id,
		&Model.From_date,
		&Model.To_date)
	if err != nil {
		log.Printf("%v\n", err)
		return &ReturningStruct{}, err
	}
	Result := MakeDataInSturct(&Model)

	return Result, nil
}

func MakeDataInSturct(modelData *AutomatLocationModel) *ReturningStruct {
	stuctReturn := ReturningStruct{}
	stuctReturn.Id = modelData.Id
	stuctReturn.Automat_id = modelData.Automat_id.Int32
	stuctReturn.Company_point_id = modelData.Company_point_id.Int32
	stuctReturn.From_date = modelData.From_date.String
	stuctReturn.To_date = modelData.To_date.String
	return &stuctReturn
}
