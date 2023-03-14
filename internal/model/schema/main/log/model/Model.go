package model

import (
	"database/sql"
	model_interface "ephorservices/pkg/orm/model"
	"reflect"
	"strings"

	pgx "github.com/jackc/pgx/v5"
)

var Schema string = "main"
var nameTable string = "log"

type LogModel struct {
	Id              int
	Address         sql.NullString
	Login           sql.NullString
	Date            sql.NullString
	Request_id      sql.NullString
	Request_uri     sql.NullString
	Request_data    sql.NullString
	Response        sql.NullString
	Runtime         sql.NullInt32
	Runtime_details sql.NullInt32
}

func New() model_interface.Model {
	return &LogModel{}
}

func (lm *LogModel) New() model_interface.Model {
	return &LogModel{}
}

func (lm *LogModel) GetNameSchema(accountNumber int) string {
	return Schema
}

func (lm *LogModel) GetNameTable() string {
	return nameTable
}

func (lm *LogModel) GetFields() []string {
	fields := make([]string, 0, 1)
	v := reflect.ValueOf(lm).Elem()
	typeOfS := v.Type()
	for i := 0; i < v.NumField(); i++ {
		field := typeOfS.Field(i).Name
		fields = append(fields, field)
	}
	return fields
}

func (lm *LogModel) GetIdKey() int64 {
	return int64(lm.Id)
}

func (lm *LogModel) ScanModelRow(row pgx.Row) (model_interface.Model, error) {
	model := LogModel{}
	err := row.Scan(&model.Id,
		&model.Address,
		&model.Login,
		&model.Date,
		&model.Request_id,
		&model.Request_uri,
		&model.Request_data,
		&model.Response,
		&model.Runtime,
		&model.Runtime_details)
	if err != nil {
		return &model, err
	}
	return &model, nil
}

func (lm *LogModel) ScanModelRows(rows pgx.Rows) ([]model_interface.Model, error) {
	defer rows.Close()
	Result := make([]model_interface.Model, 0, 1)
	for rows.Next() {
		resultStruct, err := lm.ScanModelRow(rows)
		if err != nil {
			continue
		}
		Result = append(Result, resultStruct)
	}
	return Result, nil
}

func (lm *LogModel) GetName() string {
	return "LogModel"
}

func (lm *LogModel) Get() map[string]interface{} {
	model := make(map[string]interface{})
	v := reflect.ValueOf(lm).Elem()
	typeOfS := v.Type()
	for i := 0; i < v.NumField(); i++ {
		field, ok := lm.GetNameFieldOfTag(typeOfS.Field(i))
		if !ok {
			field = strings.ToLower(typeOfS.Field(i).Name)
		}
		FieldValue := v.FieldByName(field)
		model[field] = lm.GetValueField(&FieldValue)
	}
	return model
}

func (lm *LogModel) GetNameFieldOfTag(field reflect.StructField) (string, bool) {
	if value, ok := field.Tag.Lookup("json"); ok {
		return value, ok
	} else {
		return "", false
	}

}

func (lm *LogModel) GetValueField(f *reflect.Value) interface{} {
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
