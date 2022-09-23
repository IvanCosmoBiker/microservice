package model

import (
	"database/sql"
	modelInterface "ephorservices/pkg/model/schema/interface/model"
	"log"

	pgx "github.com/jackc/pgx/v4"
)

var nameTable string = "company_point"
var nameSchema string = "account"

type CompanyPointModel struct {
	Id              int
	Company_id      sql.NullInt32
	Name            sql.NullString
	Address         sql.NullString
	Comment         sql.NullString
	State           sql.NullInt32
	Notification_id sql.NullInt32
}

func Init() modelInterface.Model {
	return &CompanyPointModel{}
}

func (am *CompanyPointModel) GetNameTable() string {
	return nameTable
}

func (am *CompanyPointModel) GetNameSchema(accountId string) string {
	return nameSchema + "" + accountId
}

type ReturningStruct struct {
	Id              int
	Company_id      int32
	Name            string
	Address         string
	Comment         string
	State           int32
	Notification_id int32
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
	Model := CompanyPointModel{}
	err := row.Scan(&Model.Id,
		&Model.Company_id,
		&Model.Name,
		&Model.Address,
		&Model.Comment,
		&Model.State,
		&Model.Notification_id)
	if err != nil {
		log.Printf("%v\n", err)
		return &ReturningStruct{}, err
	}
	Result := MakeDataInSturct(&Model)
	return Result, nil
}

func MakeDataInSturct(modelData *CompanyPointModel) *ReturningStruct {
	stuctReturn := ReturningStruct{}
	stuctReturn.Id = modelData.Id
	stuctReturn.Company_id = modelData.Company_id.Int32
	stuctReturn.Name = modelData.Name.String
	stuctReturn.Address = modelData.Address.String
	stuctReturn.Comment = modelData.Comment.String
	stuctReturn.State = modelData.State.Int32
	stuctReturn.Notification_id = modelData.Notification_id.Int32
	return &stuctReturn
}
