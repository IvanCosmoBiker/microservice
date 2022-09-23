package model

import (
	"database/sql"
	"encoding/json"
	"ephorservices/pkg/model/schema/interface/model"
	"log"

	pgx "github.com/jackc/pgx/v4"
)

var nameTable string = "transaction_product"
var nameSchema string = "main"

type TransactionProductModel struct {
	Id             int
	Transaction_id sql.NullInt32
	Name           sql.NullString
	Select_id      sql.NullString
	Ware_id        sql.NullInt32
	Value          sql.NullInt32
	Tax_rate       sql.NullInt32
	Quantity       sql.NullInt32
}

func Init() model.Model {
	return &TransactionProductModel{}
}

func (cm *TransactionProductModel) InitData(jsonData []byte) {
	json.Unmarshal(jsonData, &cm)
}

func (cm *TransactionProductModel) JsonSerialize() []byte {
	data, _ := json.Marshal(cm)
	return data
}

func (cm *TransactionProductModel) GetNameTable() string {
	return nameTable
}

func (cm *TransactionProductModel) GetNameSchema(account string) string {
	return nameSchema
}

type ReturningStruct struct {
	Id             int
	Transaction_id int32
	Name           string
	Select_id      string
	Ware_id        int32
	Value          int32
	Quantity       int32
}

func ScanModelRows(rows pgx.Rows) ([]*ReturningStruct, error) {
	defer rows.Close()
	Result := []*ReturningStruct{}
	for rows.Next() {
		resultStruct, err := ScanModelRow(rows)
		if err != nil {
			log.Println(err)
			continue
		}
		Result = append(Result, resultStruct)
	}
	return Result, nil
}

func ScanModelRow(row pgx.Row) (*ReturningStruct, error) {
	Transaction := TransactionProductModel{}
	err := row.Scan(&Transaction.Id,
		&Transaction.Transaction_id,
		&Transaction.Name,
		&Transaction.Select_id,
		&Transaction.Ware_id,
		&Transaction.Value,
		&Transaction.Quantity)
	if err != nil {
		log.Printf("%v", err)
		return &ReturningStruct{}, err
	}
	Result := MakeDataInSturct(&Transaction)
	return Result, nil
}

func MakeDataInSturct(modelData *TransactionProductModel) *ReturningStruct {
	stuctReturn := ReturningStruct{}
	stuctReturn.Id = modelData.Id
	stuctReturn.Transaction_id = modelData.Transaction_id.Int32
	stuctReturn.Name = modelData.Name.String
	stuctReturn.Select_id = modelData.Select_id.String
	stuctReturn.Ware_id = modelData.Ware_id.Int32
	stuctReturn.Value = modelData.Value.Int32
	stuctReturn.Quantity = modelData.Quantity.Int32
	return &stuctReturn
}
