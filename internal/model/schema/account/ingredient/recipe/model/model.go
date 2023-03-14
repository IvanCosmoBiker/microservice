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
var table = "recipe"

type RecipeModel struct {
	Id      int
	Ware_id sql.NullInt32
	Name    sql.NullString
}

func New() model_interface.Model {
	return &RecipeModel{}
}

func (rm *RecipeModel) New() model_interface.Model {
	return &RecipeModel{}
}

func (rm *RecipeModel) GetNameSchema(accountNumber int) string {
	return fmt.Sprintf("%s%v", schema, accountNumber)
}

func (rm *RecipeModel) GetNameTable() string {
	return table
}

func (rm *RecipeModel) GetFields() []string {
	fields := make([]string, 0, 1)
	v := reflect.ValueOf(rm).Elem()
	typeOfS := v.Type()
	for i := 0; i < v.NumField(); i++ {
		field := typeOfS.Field(i).Name
		fields = append(fields, field)
	}
	return fields
}

func (rm *RecipeModel) GetIdKey() int64 {
	return int64(rm.Id)
}

func (rm *RecipeModel) ScanModelRow(row pgx.Row) (model_interface.Model, error) {
	model := RecipeModel{}
	err := row.Scan(&model.Id,
		&model.Ware_id,
		&model.Name)
	if err != nil {
		return &model, err
	}
	return &model, nil
}

func (rm *RecipeModel) ScanModelRows(rows pgx.Rows) ([]model_interface.Model, error) {
	defer rows.Close()
	Result := make([]model_interface.Model, 0, 1)
	for rows.Next() {
		resultStruct, err := rm.ScanModelRow(rows)
		if err != nil {
			continue
		}
		Result = append(Result, resultStruct)
	}
	return Result, nil
}

func (rm *RecipeModel) GetName() string {
	return "RecipeModel"
}

func (rm *RecipeModel) Get() map[string]interface{} {
	model := make(map[string]interface{})
	v := reflect.ValueOf(rm).Elem()
	typeOfS := v.Type()
	for i := 0; i < v.NumField(); i++ {
		field, ok := rm.GetNameFieldOfTag(typeOfS.Field(i))
		if !ok {
			field = strings.ToLower(typeOfS.Field(i).Name)
		}
		FieldValue := v.FieldByName(field)
		model[field] = rm.GetValueField(&FieldValue)
	}
	return model
}

func (rm *RecipeModel) GetNameFieldOfTag(field reflect.StructField) (string, bool) {
	if value, ok := field.Tag.Lookup("json"); ok {
		return value, ok
	} else {
		return "", false
	}

}

func (rm *RecipeModel) GetValueField(f *reflect.Value) interface{} {
	typeOfS := f.Type()
	switch typeOfS.Name() {
	case "NullInt32":
		return f.Interface().(sql.NullInt32).Int32
	case "NullInt64":
		return f.Interface().(sql.NullInt64).Int64
	case "NullString":
		return f.Interface().(sql.NullString).String
	case "int", "int32", "int8", "int16", "int64":
		return f.Int()
	case "uint", "uint8", "uint16", "uint32", "uint64":
		return f.Uint()
	case "string":
		return f.String()
	}
	return nil
}
