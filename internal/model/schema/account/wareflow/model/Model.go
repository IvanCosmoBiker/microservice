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
var table = "wareflow"

type WareFlowModel struct {
	Id           int
	Date         sql.NullString
	Load_date    sql.NullString
	Automat_id   sql.NullInt32
	Collector_id sql.NullInt32
	Operator_id  sql.NullInt32
	State        sql.NullInt32
	Type         sql.NullInt32
}

func New() model_interface.Model {
	return &WareFlowModel{}
}

func (wfm *WareFlowModel) New() model_interface.Model {
	return &WareFlowModel{}
}

func (wfm *WareFlowModel) GetNameSchema(accountNumber int) string {
	return fmt.Sprintf("%s%v", schema, accountNumber)
}

func (wfm *WareFlowModel) GetNameTable() string {
	return table
}

func (wfm *WareFlowModel) GetFields() []string {
	fields := make([]string, 0, 1)
	v := reflect.ValueOf(wfm).Elem()
	typeOfS := v.Type()
	for i := 0; i < v.NumField(); i++ {
		field := typeOfS.Field(i).Name
		fields = append(fields, field)
	}
	return fields
}

func (wfm *WareFlowModel) GetIdKey() int64 {
	return int64(wfm.Id)
}

func (wfm *WareFlowModel) ScanModelRow(row pgx.Row) (model_interface.Model, error) {
	model := WareFlowModel{}
	err := row.Scan(&model.Id,
		&model.Date,
		&model.Load_date,
		&model.Automat_id,
		&model.Collector_id,
		&model.Operator_id,
		&model.State,
		&model.Type)
	if err != nil {
		return &model, err
	}
	return &model, nil
}

func (wfm *WareFlowModel) ScanModelRows(rows pgx.Rows) ([]model_interface.Model, error) {
	defer rows.Close()
	Result := make([]model_interface.Model, 0, 1)
	for rows.Next() {
		resultStruct, err := wfm.ScanModelRow(rows)
		if err != nil {
			continue
		}
		Result = append(Result, resultStruct)
	}
	return Result, nil
}

func (wfm *WareFlowModel) GetName() string {
	return "WareFlowModel"
}

func (wfm *WareFlowModel) Get() map[string]interface{} {
	model := make(map[string]interface{})
	v := reflect.ValueOf(wfm).Elem()
	typeOfS := v.Type()
	for i := 0; i < v.NumField(); i++ {
		field, ok := wfm.GetNameFieldOfTag(typeOfS.Field(i))
		if !ok {
			field = strings.ToLower(typeOfS.Field(i).Name)
		}
		FieldValue := v.FieldByName(field)
		model[field] = wfm.GetValueField(&FieldValue)
	}
	return model
}

func (wfm *WareFlowModel) GetNameFieldOfTag(field reflect.StructField) (string, bool) {
	if value, ok := field.Tag.Lookup("json"); ok {
		return value, ok
	} else {
		return "", false
	}

}

func (wfm *WareFlowModel) GetValueField(f *reflect.Value) interface{} {
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
