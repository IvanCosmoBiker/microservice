package model

// public static $fields = array(
// 	'id'    => array('type' => 'int', 'id' => true, 'autoIncrement' => true),
// 	'name'  => array('type' => 'string', 'max' => 50),
// 	'type'  => array('type' => 'int'),
// 	'pack'  => array('type' => 'int'),
// 	'code'  => array('type' => 'string', 'max' => 36),
// 	'state' => array('type' => 'int', 'default' => 1)
// );

import (
	"database/sql"
	model_interface "ephorservices/pkg/orm/model"
	"fmt"
	"reflect"
	"strings"

	pgx "github.com/jackc/pgx/v5"
)

var schema = "account"
var table = "ingredient"

type IngredientModel struct {
	Id    int
	Name  sql.NullString
	Type  sql.NullInt32
	Pack  sql.NullInt32
	Code  sql.NullString
	State sql.NullInt32
}

func New() model_interface.Model {
	return &IngredientModel{}
}

func (im *IngredientModel) New() model_interface.Model {
	return &IngredientModel{}
}

func (im *IngredientModel) GetNameSchema(accountNumber int) string {
	return fmt.Sprintf("%s%v", schema, accountNumber)
}

func (im *IngredientModel) GetNameTable() string {
	return table
}

func (im *IngredientModel) GetFields() []string {
	fields := make([]string, 0, 1)
	v := reflect.ValueOf(im).Elem()
	typeOfS := v.Type()
	for i := 0; i < v.NumField(); i++ {
		field := typeOfS.Field(i).Name
		fields = append(fields, field)
	}
	return fields
}

func (im *IngredientModel) GetIdKey() int64 {
	return int64(im.Id)
}

func (im *IngredientModel) ScanModelRow(row pgx.Row) (model_interface.Model, error) {
	model := IngredientModel{}
	err := row.Scan(&model.Id,
		&model.Name,
		&model.Type,
		&model.Pack,
		&model.Code,
		&model.State)
	if err != nil {
		return &model, err
	}
	return &model, nil
}

func (im *IngredientModel) ScanModelRows(rows pgx.Rows) ([]model_interface.Model, error) {
	defer rows.Close()
	Result := make([]model_interface.Model, 0, 1)
	for rows.Next() {
		resultStruct, err := im.ScanModelRow(rows)
		if err != nil {
			continue
		}
		Result = append(Result, resultStruct)
	}
	return Result, nil
}

func (im *IngredientModel) GetName() string {
	return "IngredientModel"
}

func (im *IngredientModel) Get() map[string]interface{} {
	model := make(map[string]interface{})
	v := reflect.ValueOf(im).Elem()
	typeOfS := v.Type()
	for i := 0; i < v.NumField(); i++ {
		field, ok := im.GetNameFieldOfTag(typeOfS.Field(i))
		if !ok {
			field = strings.ToLower(typeOfS.Field(i).Name)
		}
		FieldValue := v.FieldByName(field)
		model[field] = im.GetValueField(&FieldValue)
	}
	return model
}

func (im *IngredientModel) GetNameFieldOfTag(field reflect.StructField) (string, bool) {
	if value, ok := field.Tag.Lookup("json"); ok {
		return value, ok
	} else {
		return "", false
	}
}

func (im *IngredientModel) GetValueField(f *reflect.Value) interface{} {
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
