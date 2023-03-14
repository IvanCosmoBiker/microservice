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
var table = "automat_config"

type AutomatConfigModel struct {
	Id         int
	Automat_id sql.NullInt32
	Config_id  sql.NullInt32
	From_date  sql.NullString
	To_date    sql.NullString
}

func New() model_interface.Model {
	return &AutomatConfigModel{}
}

func (acm *AutomatConfigModel) New() model_interface.Model {
	return &AutomatConfigModel{}
}

func (acm *AutomatConfigModel) GetNameSchema(accountNumber int) string {
	return fmt.Sprintf("%s%v", schema, accountNumber)
}

func (acm *AutomatConfigModel) GetNameTable() string {
	return table
}

func (acm *AutomatConfigModel) GetFields() []string {
	fields := make([]string, 0, 1)
	v := reflect.ValueOf(acm).Elem()
	typeOfS := v.Type()
	for i := 0; i < v.NumField(); i++ {
		field := typeOfS.Field(i).Name
		fields = append(fields, field)
	}
	return fields
}

func (acm *AutomatConfigModel) GetIdKey() int64 {
	return int64(acm.Id)
}

func (acm *AutomatConfigModel) ScanModelRow(row pgx.Row) (model_interface.Model, error) {
	model := AutomatConfigModel{}
	err := row.Scan(&model.Id,
		&model.Automat_id,
		&model.Config_id,
		&model.From_date,
		&model.To_date)
	if err != nil {
		return &model, err
	}
	return &model, nil
}

func (acm *AutomatConfigModel) ScanModelRows(rows pgx.Rows) ([]model_interface.Model, error) {
	defer rows.Close()
	Result := make([]model_interface.Model, 0, 1)
	for rows.Next() {
		resultStruct, err := acm.ScanModelRow(rows)
		if err != nil {
			continue
		}
		Result = append(Result, resultStruct)
	}
	return Result, nil
}

func (acm *AutomatConfigModel) GetName() string {
	return "AutomatConfigModel"
}

func (acm *AutomatConfigModel) Get() map[string]interface{} {
	model := make(map[string]interface{})
	v := reflect.ValueOf(acm).Elem()
	typeOfS := v.Type()
	for i := 0; i < v.NumField(); i++ {
		field, ok := acm.GetNameFieldOfTag(typeOfS.Field(i))
		if !ok {
			field = strings.ToLower(typeOfS.Field(i).Name)
		}
		FieldValue := v.FieldByName(field)
		model[field] = acm.GetValueField(&FieldValue)
	}
	return model
}

func (acm *AutomatConfigModel) GetNameFieldOfTag(field reflect.StructField) (string, bool) {
	if value, ok := field.Tag.Lookup("json"); ok {
		return value, ok
	} else {
		return "", false
	}
}

func (acm *AutomatConfigModel) GetValueField(f *reflect.Value) interface{} {
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
