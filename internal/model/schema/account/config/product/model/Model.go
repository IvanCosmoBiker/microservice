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
var table = "config_product"

type ConfigProductModel struct {
	Id        int
	Name      sql.NullString
	Config_id sql.NullInt32
	Select_id sql.NullString
	Cid       sql.NullInt32
	Ware_id   sql.NullInt32
	Tax_rate  sql.NullInt32
	Type      sql.NullInt32
	Qt_recom  sql.NullInt32
	Qt_max    sql.NullInt32
	Recipe_id sql.NullInt32
}

func New() model_interface.Model {
	return &ConfigProductModel{}
}

func (cpm *ConfigProductModel) New() model_interface.Model {
	return &ConfigProductModel{}
}

func (cpm *ConfigProductModel) GetNameSchema(accountNumber int) string {
	return fmt.Sprintf("%s%v", schema, accountNumber)
}

func (cpm *ConfigProductModel) GetNameTable() string {
	return table
}

func (cpm *ConfigProductModel) GetFields() []string {
	fields := make([]string, 0, 1)
	v := reflect.ValueOf(cpm).Elem()
	typeOfS := v.Type()
	for i := 0; i < v.NumField(); i++ {
		field := typeOfS.Field(i).Name
		fields = append(fields, field)
	}
	return fields
}

func (cpm *ConfigProductModel) GetIdKey() int64 {
	return int64(cpm.Id)
}

func (cpm *ConfigProductModel) ScanModelRow(row pgx.Row) (model_interface.Model, error) {
	model := ConfigProductModel{}
	err := row.Scan(&model.Id,
		&model.Name,
		&model.Config_id,
		&model.Select_id,
		&model.Cid,
		&model.Ware_id,
		&model.Tax_rate,
		&model.Type,
		&model.Qt_recom,
		&model.Qt_max,
		&model.Recipe_id)
	if err != nil {
		return &model, err
	}
	return &model, nil
}

func (cpm *ConfigProductModel) ScanModelRows(rows pgx.Rows) ([]model_interface.Model, error) {
	defer rows.Close()
	Result := make([]model_interface.Model, 0, 1)
	for rows.Next() {
		resultStruct, err := cpm.ScanModelRow(rows)
		if err != nil {
			continue
		}
		Result = append(Result, resultStruct)
	}
	return Result, nil
}

func (cpm *ConfigProductModel) GetName() string {
	return "ConfigProductModel"
}

func (cpm *ConfigProductModel) Get() map[string]interface{} {
	model := make(map[string]interface{})
	v := reflect.ValueOf(cpm).Elem()
	typeOfS := v.Type()
	for i := 0; i < v.NumField(); i++ {
		field, ok := cpm.GetNameFieldOfTag(typeOfS.Field(i))
		if !ok {
			field = strings.ToLower(typeOfS.Field(i).Name)
		}
		FieldValue := v.FieldByName(field)
		model[field] = cpm.GetValueField(&FieldValue)
	}
	return model
}

func (cpm *ConfigProductModel) GetNameFieldOfTag(field reflect.StructField) (string, bool) {
	if value, ok := field.Tag.Lookup("json"); ok {
		return value, ok
	} else {
		return "", false
	}

}

func (cpm *ConfigProductModel) GetValueField(f *reflect.Value) interface{} {
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
