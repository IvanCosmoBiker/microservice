package model

import (
	"database/sql"
	model_interface "ephorservices/pkg/orm/model"
	"reflect"
	"strings"

	pgx "github.com/jackc/pgx/v5"
)

var Schema string = "main"
var nameTable string = "modem_command"

type CommandModel struct {
	Id             int
	Modem_id       sql.NullInt32
	Command        sql.NullInt32
	Command_param1 sql.NullInt32
	Date           sql.NullString
	Sended         sql.NullInt32
}

func New() model_interface.Model {
	return &CommandModel{}
}

func (cm *CommandModel) New() model_interface.Model {
	return &CommandModel{}
}

func (cm *CommandModel) GetNameSchema(accountNumber int) string {
	return Schema
}

func (cm *CommandModel) GetNameTable() string {
	return nameTable
}

func (cm *CommandModel) GetFields() []string {
	fields := make([]string, 0, 1)
	v := reflect.ValueOf(cm).Elem()
	typeOfS := v.Type()
	for i := 0; i < v.NumField(); i++ {
		field := typeOfS.Field(i).Name
		fields = append(fields, field)
	}
	return fields
}

func (cm *CommandModel) GetIdKey() int64 {
	return int64(cm.Id)
}

func (cm *CommandModel) ScanModelRow(row pgx.Row) (model_interface.Model, error) {
	model := CommandModel{}
	err := row.Scan(&model.Id,
		&model.Modem_id,
		&model.Command,
		&model.Command_param1,
		&model.Date,
		&model.Sended)
	if err != nil {
		return &model, err
	}
	return &model, nil
}

func (cm *CommandModel) ScanModelRows(rows pgx.Rows) ([]model_interface.Model, error) {
	defer rows.Close()
	Result := make([]model_interface.Model, 0, 1)
	for rows.Next() {
		resultStruct, err := cm.ScanModelRow(rows)
		if err != nil {
			continue
		}
		Result = append(Result, resultStruct)
	}
	return Result, nil
}

func (cm *CommandModel) GetName() string {
	return "CommandModel"
}

func (cm *CommandModel) Get() map[string]interface{} {
	model := make(map[string]interface{})
	v := reflect.ValueOf(cm).Elem()
	typeOfS := v.Type()
	for i := 0; i < v.NumField(); i++ {
		field, ok := cm.GetNameFieldOfTag(typeOfS.Field(i))
		if !ok {
			field = strings.ToLower(typeOfS.Field(i).Name)
		}
		FieldValue := v.FieldByName(field)
		model[field] = cm.GetValueField(&FieldValue)
	}
	return model
}

func (cm *CommandModel) GetNameFieldOfTag(field reflect.StructField) (string, bool) {
	if value, ok := field.Tag.Lookup("json"); ok {
		return value, ok
	} else {
		return "", false
	}

}

func (cm *CommandModel) GetValueField(f *reflect.Value) interface{} {
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
