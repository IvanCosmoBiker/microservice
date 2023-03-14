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
var table = "recipe_ingredient"

type RecipeIngredientModel struct {
	Id            int
	Recipe_id     sql.NullInt32
	Ingredient_id sql.NullInt32
	Count         sql.NullInt32
}

func New() model_interface.Model {
	return &RecipeIngredientModel{}
}

func (rim *RecipeIngredientModel) New() model_interface.Model {
	return &RecipeIngredientModel{}
}

func (rim *RecipeIngredientModel) GetNameSchema(accountNumber int) string {
	return fmt.Sprintf("%s%v", schema, accountNumber)
}

func (rim *RecipeIngredientModel) GetNameTable() string {
	return table
}

func (rim *RecipeIngredientModel) GetFields() []string {
	fields := make([]string, 0, 1)
	v := reflect.ValueOf(rim).Elem()
	typeOfS := v.Type()
	for i := 0; i < v.NumField(); i++ {
		field := typeOfS.Field(i).Name
		fields = append(fields, field)
	}
	return fields
}

func (rim *RecipeIngredientModel) GetIdKey() int64 {
	return int64(rim.Id)
}

func (rim *RecipeIngredientModel) ScanModelRow(row pgx.Row) (model_interface.Model, error) {
	model := RecipeIngredientModel{}
	err := row.Scan(&model.Id,
		&model.Recipe_id,
		&model.Ingredient_id,
		&model.Count)
	if err != nil {
		return &model, err
	}
	return &model, nil
}

func (rim *RecipeIngredientModel) ScanModelRows(rows pgx.Rows) ([]model_interface.Model, error) {
	defer rows.Close()
	Result := make([]model_interface.Model, 0, 1)
	for rows.Next() {
		resultStruct, err := rim.ScanModelRow(rows)
		if err != nil {
			continue
		}
		Result = append(Result, resultStruct)
	}
	return Result, nil
}

func (rim *RecipeIngredientModel) GetName() string {
	return "RecipeIngredientModel"
}

func (rim *RecipeIngredientModel) Get() map[string]interface{} {
	model := make(map[string]interface{})
	v := reflect.ValueOf(rim).Elem()
	typeOfS := v.Type()
	for i := 0; i < v.NumField(); i++ {
		field, ok := rim.GetNameFieldOfTag(typeOfS.Field(i))
		if !ok {
			field = strings.ToLower(typeOfS.Field(i).Name)
		}
		FieldValue := v.FieldByName(field)
		model[field] = rim.GetValueField(&FieldValue)
	}
	return model
}

func (rim *RecipeIngredientModel) GetNameFieldOfTag(field reflect.StructField) (string, bool) {
	if value, ok := field.Tag.Lookup("json"); ok {
		return value, ok
	} else {
		return "", false
	}
}

func (rim *RecipeIngredientModel) GetValueField(f *reflect.Value) interface{} {
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
