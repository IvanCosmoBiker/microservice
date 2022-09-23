package model

import (
	"database/sql"
	modelInterface "ephorservices/pkg/model/schema/interface/model"
	"log"

	pgx "github.com/jackc/pgx/v4"
)

var nameTable string = "company_point"
var nameSchema string = "account"

type CompanyModel struct {
	Id               int
	Name             sql.NullString
	Comment          sql.NullString
	State            sql.NullInt32
	Balance          sql.NullInt64
	Ephor_manager_id sql.NullInt32
}

func Init() modelInterface.Model {
	return &CompanyModel{}
}

func (am *CompanyModel) GetNameTable() string {
	return nameTable
}

func (am *CompanyModel) GetNameSchema(accountId string) string {
	return nameSchema + "" + accountId
}

type ReturningStruct struct {
	Id               int
	Name             string
	Comment          string
	State            int32
	Balance          int64
	Ephor_manager_id int32
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
	Model := CompanyModel{}
	err := row.Scan(&Model.Id,
		&Model.Name,
		&Model.Comment,
		&Model.State,
		&Model.Balance,
		&Model.Ephor_manager_id)
	if err != nil {
		log.Printf("%v\n", err)
		return &ReturningStruct{}, err
	}
	Result := MakeDataInSturct(&Model)

	return Result, nil
}

func MakeDataInSturct(modelData *CompanyModel) *ReturningStruct {
	stuctReturn := ReturningStruct{}
	stuctReturn.Id = modelData.Id
	stuctReturn.Name = modelData.Name.String
	stuctReturn.Comment = modelData.Comment.String
	stuctReturn.State = modelData.State.Int32
	stuctReturn.Balance = modelData.Balance.Int64
	stuctReturn.Ephor_manager_id = modelData.Ephor_manager_id.Int32
	return &stuctReturn
}
