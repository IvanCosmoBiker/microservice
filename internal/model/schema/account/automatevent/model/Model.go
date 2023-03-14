package model

import (
	"database/sql"
	model_interface "ephorservices/pkg/orm/model"
	"fmt"
	"reflect"
	"strings"

	pgx "github.com/jackc/pgx/v5"
)

var schema = "account"
var table = "automat_event"

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

func New() model_interface.Model {
	return &AutomatEventModel{}
}

func (aem *AutomatEventModel) New() model_interface.Model {
	return &AutomatEventModel{}
}

func (aem *AutomatEventModel) GetNameSchema(accountNumber int) string {
	return fmt.Sprintf("%s%v", schema, accountNumber)
}

func (aem *AutomatEventModel) GetNameTable() string {
	return table
}

func (aem *AutomatEventModel) GetFields() []string {
	fields := make([]string, 0, 1)
	v := reflect.ValueOf(aem).Elem()
	typeOfS := v.Type()
	for i := 0; i < v.NumField(); i++ {
		field := typeOfS.Field(i).Name
		fields = append(fields, field)
	}
	return fields
}

func (aem *AutomatEventModel) GetIdKey() int64 {
	return int64(aem.Id)
}

func (aem *AutomatEventModel) ScanModelRow(row pgx.Row) (model_interface.Model, error) {
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
		return &Model, err
	}
	return &Model, nil
}

func (aem *AutomatEventModel) ScanModelRows(rows pgx.Rows) ([]model_interface.Model, error) {
	defer rows.Close()
	Result := make([]model_interface.Model, 0, 1)
	for rows.Next() {
		resultStruct, err := aem.ScanModelRow(rows)
		if err != nil {
			continue
		}
		Result = append(Result, resultStruct)
	}
	return Result, nil
}

func (aem *AutomatEventModel) GetName() string {
	return "AutomatEventModel"
}

func (aem *AutomatEventModel) Get() map[string]interface{} {
	model := make(map[string]interface{})
	v := reflect.ValueOf(aem).Elem()
	typeOfS := v.Type()
	for i := 0; i < v.NumField(); i++ {
		field, ok := aem.GetNameFieldOfTag(typeOfS.Field(i))
		if !ok {
			field = strings.ToLower(typeOfS.Field(i).Name)
		}
		FieldValue := v.FieldByName(field)
		model[field] = aem.GetValueField(&FieldValue)
	}
	return model
}

func (aem *AutomatEventModel) GetNameFieldOfTag(field reflect.StructField) (string, bool) {
	if value, ok := field.Tag.Lookup("json"); ok {
		return value, ok
	} else {
		return "", false
	}

}

func (aem *AutomatEventModel) GetValueField(f *reflect.Value) interface{} {
	typeOfS := f.Type()
	switch typeOfS.Name() {
	case "NullInt32":
		return f.Interface().(sql.NullInt32).Int32
	case "NullString":
		return f.Interface().(sql.NullString).String
	case "int", "int32", "int8", "int16", "int64":
		return f.Int()
	case "uint", "uint8", "uint16", "uint32", "uint64":
		return f.Uint()
	case "string":
		return f.String()
	}
	return false
}
