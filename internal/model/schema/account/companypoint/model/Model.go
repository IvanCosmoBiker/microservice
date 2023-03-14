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
var table = "company_point"

type CompanyPointModel struct {
	Id              int
	Company_id      sql.NullInt32
	Name            sql.NullString
	Address         sql.NullString
	Comment         sql.NullString
	State           sql.NullInt32
	Notification_id sql.NullInt32
}

func New() model_interface.Model {
	return &CompanyPointModel{}
}

func (cpm *CompanyPointModel) New() model_interface.Model {
	return &CompanyPointModel{}
}

func (cpm *CompanyPointModel) GetNameSchema(accountNumber int) string {
	return fmt.Sprintf("%s%v", schema, accountNumber)
}

func (cpm *CompanyPointModel) GetNameTable() string {
	return table
}

func (cpm *CompanyPointModel) GetFields() []string {
	fields := make([]string, 0, 1)
	v := reflect.ValueOf(cpm).Elem()
	typeOfS := v.Type()
	for i := 0; i < v.NumField(); i++ {
		field := typeOfS.Field(i).Name
		fields = append(fields, field)
	}
	return fields
}

func (cpm *CompanyPointModel) GetIdKey() int64 {
	return int64(cpm.Id)
}

func (cpm *CompanyPointModel) ScanModelRow(row pgx.Row) (model_interface.Model, error) {
	model := CompanyPointModel{}
	err := row.Scan(&model.Id,
		&model.Company_id,
		&model.Name,
		&model.Address,
		&model.Comment,
		&model.State,
		&model.Notification_id)
	if err != nil {
		return &model, err
	}
	return &model, nil
}

func (cpm *CompanyPointModel) ScanModelRows(rows pgx.Rows) ([]model_interface.Model, error) {
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

func (cpm *CompanyPointModel) GetName() string {
	return "CompanyPointModel"
}

func (cpm *CompanyPointModel) Get() map[string]interface{} {
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

func (cpm *CompanyPointModel) GetNameFieldOfTag(field reflect.StructField) (string, bool) {
	if value, ok := field.Tag.Lookup("json"); ok {
		return value, ok
	} else {
		return "", false
	}

}

func (cpm *CompanyPointModel) GetValueField(f *reflect.Value) interface{} {
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
