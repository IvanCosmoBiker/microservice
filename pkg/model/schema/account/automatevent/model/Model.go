package model

import (
	"database/sql"
	"encoding/json"
	modelInterface "ephorservices/pkg/model/schema/interface/model"
	"log"

	pgx "github.com/jackc/pgx/v4"
)

var nameTable string = "automat_event"
var nameSchema string = "account"

type AutomatEventModel struct {
	Id             int
	Automat_id     sql.NullInt32
	Operator_id    sql.NullInt32
	Date           sql.NullString
	Modem_date     sql.NullString
	Fiscal_date    sql.NullString
	Update_date    sql.NullString
	Type           sql.NullInt32
	Category       sql.NullInt32
	Select_id      sql.NullString
	Ware_id        sql.NullInt32
	Name           sql.NullString
	Payment_device sql.NullString
	Price_list     sql.NullInt32
	Value          sql.NullInt32
	Credit         sql.NullInt32
	Tax_system     sql.NullInt32
	Tax_rate       sql.NullInt32
	Tax_value      sql.NullInt32
	Fn             sql.NullInt64
	Fd             sql.NullInt32
	Fp             sql.NullInt64
	Fp_string      sql.NullString
	Id_fr          sql.NullString
	Status         sql.NullInt32
	Point_id       sql.NullInt32
	Loyality_type  sql.NullInt32
	Loyality_code  sql.NullString
	Error_detail   sql.NullString
	Warehouse_id   sql.NullString
	Type_fr        sql.NullInt32
}

func Init() modelInterface.Model {
	return &AutomatEventModel{}
}

func (aem *AutomatEventModel) InitData(jsonData []byte) {
	json.Unmarshal(jsonData, &aem)
}

func (aem *AutomatEventModel) JsonSerialize() []byte {
	data, _ := json.Marshal(aem)
	return data
}

func (aem *AutomatEventModel) GetNameTable() string {
	return nameTable
}

func (aem *AutomatEventModel) GetNameSchema(accountId string) string {
	return nameSchema + "" + accountId
}

type ReturningStruct struct {
	Id             int
	Automat_id     int32
	Operator_id    int32
	Date           string
	Modem_date     string
	Fiscal_date    string
	Update_date    string
	Type           int32
	Category       int32
	Select_id      string
	Ware_id        int32
	Name           string
	Payment_device string
	Price_list     int32
	Value          int32
	Credit         int32
	Tax_system     int32
	Tax_rate       int32
	Tax_value      int32
	Fn             int64
	Fd             int32
	Fp             int64
	Fp_string      string
	Id_fr          string
	Status         int32
	Point_id       int32
	Loyality_type  int32
	Loyality_code  string
	Error_detail   string
	Warehouse_id   string
	Type_fr        int32
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
	Model := AutomatEventModel{}
	err := row.Scan(&Model.Id,
		&Model.Automat_id,
		&Model.Operator_id,
		&Model.Date,
		&Model.Modem_date,
		&Model.Fiscal_date,
		&Model.Update_date,
		&Model.Type,
		&Model.Category,
		&Model.Select_id,
		&Model.Ware_id,
		&Model.Name,
		&Model.Payment_device,
		&Model.Price_list,
		&Model.Value,
		&Model.Credit,
		&Model.Tax_system,
		&Model.Tax_rate,
		&Model.Tax_value,
		&Model.Fn,
		&Model.Fd,
		&Model.Fp,
		&Model.Fp_string,
		&Model.Id_fr,
		&Model.Status,
		&Model.Point_id,
		&Model.Loyality_type,
		&Model.Loyality_code,
		&Model.Error_detail,
		&Model.Warehouse_id,
		&Model.Type_fr)

	if err != nil {
		log.Printf("%v\n", err)
		return &ReturningStruct{}, err
	}
	Result := MakeDataInSturct(&Model)

	return Result, nil
}

func MakeDataInSturct(modelData *AutomatEventModel) *ReturningStruct {
	stuctReturn := ReturningStruct{}
	stuctReturn.Id = modelData.Id
	stuctReturn.Automat_id = modelData.Automat_id.Int32
	stuctReturn.Operator_id = modelData.Operator_id.Int32
	stuctReturn.Date = modelData.Date.String
	stuctReturn.Modem_date = modelData.Modem_date.String
	stuctReturn.Fiscal_date = modelData.Fiscal_date.String
	stuctReturn.Update_date = modelData.Update_date.String
	stuctReturn.Type = modelData.Type.Int32
	stuctReturn.Category = modelData.Category.Int32
	stuctReturn.Select_id = modelData.Select_id.String
	stuctReturn.Ware_id = modelData.Ware_id.Int32
	stuctReturn.Name = modelData.Name.String
	stuctReturn.Payment_device = modelData.Payment_device.String
	stuctReturn.Price_list = modelData.Price_list.Int32
	stuctReturn.Value = modelData.Value.Int32
	stuctReturn.Credit = modelData.Credit.Int32
	stuctReturn.Tax_system = modelData.Tax_system.Int32
	stuctReturn.Tax_rate = modelData.Tax_rate.Int32
	stuctReturn.Tax_value = modelData.Tax_value.Int32
	stuctReturn.Fn = modelData.Fn.Int64
	stuctReturn.Fd = modelData.Fd.Int32
	stuctReturn.Fp = modelData.Fp.Int64
	stuctReturn.Fp_string = modelData.Fp_string.String
	stuctReturn.Id_fr = modelData.Id_fr.String
	stuctReturn.Status = modelData.Status.Int32
	stuctReturn.Point_id = modelData.Point_id.Int32
	stuctReturn.Loyality_type = modelData.Loyality_type.Int32
	stuctReturn.Loyality_code = modelData.Loyality_code.String
	stuctReturn.Error_detail = modelData.Error_detail.String
	stuctReturn.Warehouse_id = modelData.Warehouse_id.String
	stuctReturn.Type_fr = modelData.Type_fr.Int32
	return &stuctReturn
}
