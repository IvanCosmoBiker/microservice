package model

import (
	"database/sql"
	model_interface "ephorservices/pkg/orm/model"
	"reflect"
	"strings"

	pgx "github.com/jackc/pgx/v5"
)

var Schema string = "main"
var nameTable string = "transaction_product"

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

func New() model_interface.Model {
	return &TransactionProductModel{}
}

func (tpm *TransactionProductModel) New() model_interface.Model {
	return &TransactionProductModel{}
}

func (tpm *TransactionProductModel) GetNameSchema(accountNumber int) string {
	return Schema
}

func (tpm *TransactionProductModel) GetNameTable() string {
	return nameTable
}

func (tpm *TransactionProductModel) GetFields() []string {
	fields := make([]string, 0, 1)
	v := reflect.ValueOf(tpm).Elem()
	typeOfS := v.Type()
	for i := 0; i < v.NumField(); i++ {
		field := typeOfS.Field(i).Name
		fields = append(fields, field)
	}
	return fields
}

func (tpm *TransactionProductModel) GetIdKey() int64 {
	return int64(tpm.Id)
}

func (tpm *TransactionProductModel) ScanModelRow(row pgx.Row) (model_interface.Model, error) {
	model := TransactionProductModel{}
	err := row.Scan(&model.Id,
		&model.Transaction_id,
		&model.Name,
		&model.Select_id,
		&model.Ware_id,
		&model.Value,
		&model.Tax_rate,
		&model.Quantity)
	if err != nil {
		return &model, err
	}
	return &model, nil
}

func (tpm *TransactionProductModel) ScanModelRows(rows pgx.Rows) ([]model_interface.Model, error) {
	defer rows.Close()
	Result := make([]model_interface.Model, 0, 1)
	for rows.Next() {
		resultStruct, err := tpm.ScanModelRow(rows)
		if err != nil {
			continue
		}
		Result = append(Result, resultStruct)
	}
	return Result, nil
}

func (tpm *TransactionProductModel) GetName() string {
	return "TransactionProductModel"
}

func (tpm *TransactionProductModel) Get() map[string]interface{} {
	model := make(map[string]interface{})
	v := reflect.ValueOf(tpm).Elem()
	typeOfS := v.Type()
	for i := 0; i < v.NumField(); i++ {
		field, ok := tpm.GetNameFieldOfTag(typeOfS.Field(i))
		if !ok {
			field = strings.ToLower(typeOfS.Field(i).Name)
		}
		FieldValue := v.FieldByName(field)
		model[field] = tpm.GetValueField(&FieldValue)
	}
	return model
}

func (tpm *TransactionProductModel) GetNameFieldOfTag(field reflect.StructField) (string, bool) {
	if value, ok := field.Tag.Lookup("json"); ok {
		return value, ok
	} else {
		return "", false
	}

}

func (tpm *TransactionProductModel) GetValueField(f *reflect.Value) interface{} {
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
