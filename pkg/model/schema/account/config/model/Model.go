package model

import (
	"database/sql"
	"encoding/json"
	modelInterface "ephorservices/pkg/model/schema/interface/model"
	"log"

	pgx "github.com/jackc/pgx/v4"
)

var nameTable string = "config"
var nameSchema string = "account"

type ConfigModel struct {
	Id               int
	Name             sql.NullString
	Automat_model_id sql.NullInt32
	Payment_type     sql.NullInt32
	Dex_port         sql.NullInt32
	Template_id      sql.NullInt32
	Currency         sql.NullInt32
	Decimal_point    sql.NullInt32
	Max_credit       sql.NullInt64
	Show_change      sql.NullInt32
	Cl_num           sql.NullInt32
	Cl_sf            sql.NullInt32
	Category_money   sql.NullInt32
	Tax_system       sql.NullInt32
	State            sql.NullInt32
	Price_holding    sql.NullInt32
	Multivend        sql.NullInt32
	Credit_holding   sql.NullInt32
	Cl_2_click       sql.NullInt32
	Comment          sql.NullString
	Bv_enabled       sql.NullInt32
	Vend_fail_ignore sql.NullInt32
	Cl_always_idle   sql.NullInt32
	Dex_stat         sql.NullInt32
}

func Init() modelInterface.Model {
	return &ConfigModel{}
}

func (aem *ConfigModel) InitData(jsonData []byte) {
	json.Unmarshal(jsonData, &aem)
}

func (aem *ConfigModel) JsonSerialize() []byte {
	data, _ := json.Marshal(aem)
	return data
}

func (aem *ConfigModel) GetNameTable() string {
	return nameTable
}

func (aem *ConfigModel) GetNameSchema(accountId string) string {
	return nameSchema + "" + accountId
}

type ReturningStruct struct {
	Id               int
	Name             string
	Automat_model_id int32
	Payment_type     int32
	Dex_port         int32
	Template_id      int32
	Currency         int32
	Decimal_point    int32
	Max_credit       int64
	Show_change      int32
	Cl_num           int32
	Cl_sf            int32
	Category_money   int32
	Tax_system       int32
	State            int32
	Price_holding    int32
	Multivend        int32
	Credit_holding   int32
	Cl_2_click       int32
	Comment          string
	Bv_enabled       int32
	Vend_fail_ignore int32
	Cl_always_idle   int32
	Dex_stat         int32
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
	Model := ConfigModel{}
	err := row.Scan(&Model.Id,
		&Model.Name,
		&Model.Automat_model_id,
		&Model.Payment_type,
		&Model.Dex_port,
		&Model.Template_id,
		&Model.Currency,
		&Model.Decimal_point,
		&Model.Max_credit,
		&Model.Show_change,
		&Model.Cl_num,
		&Model.Cl_sf,
		&Model.Category_money,
		&Model.Tax_system,
		&Model.State,
		&Model.Price_holding,
		&Model.Multivend,
		&Model.Credit_holding,
		&Model.Cl_2_click,
		&Model.Comment,
		&Model.Bv_enabled,
		&Model.Vend_fail_ignore,
		&Model.Cl_always_idle,
		&Model.Dex_stat)

	if err != nil {
		log.Printf("%v\n", err)
		return &ReturningStruct{}, err
	}
	Result := MakeDataInSturct(&Model)

	return Result, nil
}

func MakeDataInSturct(modelData *ConfigModel) *ReturningStruct {
	stuctReturn := ReturningStruct{}
	stuctReturn.Id = modelData.Id
	stuctReturn.Name = modelData.Name.String
	stuctReturn.Automat_model_id = modelData.Automat_model_id.Int32
	stuctReturn.Payment_type = modelData.Payment_type.Int32
	stuctReturn.Dex_port = modelData.Dex_port.Int32
	stuctReturn.Template_id = modelData.Template_id.Int32
	stuctReturn.Currency = modelData.Currency.Int32
	stuctReturn.Decimal_point = modelData.Decimal_point.Int32
	stuctReturn.Max_credit = modelData.Max_credit.Int64
	stuctReturn.Show_change = modelData.Show_change.Int32
	stuctReturn.Cl_num = modelData.Cl_num.Int32
	stuctReturn.Cl_sf = modelData.Cl_sf.Int32
	stuctReturn.Category_money = modelData.Category_money.Int32
	stuctReturn.Tax_system = modelData.Tax_system.Int32
	stuctReturn.State = modelData.State.Int32
	stuctReturn.Price_holding = modelData.Price_holding.Int32
	stuctReturn.Multivend = modelData.Multivend.Int32
	stuctReturn.Credit_holding = modelData.Credit_holding.Int32
	stuctReturn.Cl_2_click = modelData.Cl_2_click.Int32
	stuctReturn.Comment = modelData.Comment.String
	stuctReturn.Bv_enabled = modelData.Bv_enabled.Int32
	stuctReturn.Vend_fail_ignore = modelData.Vend_fail_ignore.Int32
	stuctReturn.Cl_always_idle = modelData.Cl_always_idle.Int32
	stuctReturn.Dex_stat = modelData.Dex_stat.Int32
	return &stuctReturn
}
