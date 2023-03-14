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
var table = "config_price_list"

type ConfigPriceListModel struct {
	Id        int
	Config_id sql.NullInt32
	Number    sql.NullInt32
}

func New() model_interface.Model {
	return &ConfigPriceListModel{}
}

func (cplm *ConfigPriceListModel) New() model_interface.Model {
	return &ConfigPriceListModel{}
}

func (cplm *ConfigPriceListModel) GetNameSchema(accountNumber int) string {
	return fmt.Sprintf("%s%v", schema, accountNumber)
}

func (cplm *ConfigPriceListModel) GetNameTable() string {
	return table
}

func (cplm *ConfigPriceListModel) GetFields() []string {
	fields := make([]string, 0, 1)
	v := reflect.ValueOf(cplm).Elem()
	typeOfS := v.Type()
	for i := 0; i < v.NumField(); i++ {
		field := typeOfS.Field(i).Name
		fields = append(fields, field)
	}
	return fields
}

func (cplm *ConfigPriceListModel) GetIdKey() int64 {
	return int64(cplm.Id)
}

func (cplm *ConfigPriceListModel) ScanModelRow(row pgx.Row) (model_interface.Model, error) {
	model := ConfigPriceListModel{}
	err := row.Scan(&model.Id,
		&model.Config_id,
		&model.Number)
	if err != nil {
		return &model, err
	}
	return &model, nil
}

func (cplm *ConfigPriceListModel) ScanModelRows(rows pgx.Rows) ([]model_interface.Model, error) {
	defer rows.Close()
	Result := make([]model_interface.Model, 0, 1)
	for rows.Next() {
		resultStruct, err := cplm.ScanModelRow(rows)
		if err != nil {
			continue
		}
		Result = append(Result, resultStruct)
	}
	return Result, nil
}

func (cplm *ConfigPriceListModel) GetName() string {
	return "ConfigPriceListModel"
}

func (cplm *ConfigPriceListModel) Get() map[string]interface{} {
	model := make(map[string]interface{})
	v := reflect.ValueOf(cplm).Elem()
	typeOfS := v.Type()
	for i := 0; i < v.NumField(); i++ {
		field, ok := cplm.GetNameFieldOfTag(typeOfS.Field(i))
		if !ok {
			field = strings.ToLower(typeOfS.Field(i).Name)
		}
		FieldValue := v.FieldByName(field)
		model[field] = cplm.GetValueField(&FieldValue)
	}
	return model
}

func (cplm *ConfigPriceListModel) GetNameFieldOfTag(field reflect.StructField) (string, bool) {
	if value, ok := field.Tag.Lookup("json"); ok {
		return value, ok
	} else {
		return "", false
	}

}

func (cplm *ConfigPriceListModel) GetValueField(f *reflect.Value) interface{} {
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
