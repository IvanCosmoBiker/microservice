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
var table = "wareflow_product"

type WareFlowProductModel struct {
	Id           int
	Wareflow_id  sql.NullInt32
	Select_id    sql.NullString
	Ware_id      sql.NullInt32
	Qt_max       sql.NullInt32
	Qt_recom     sql.NullInt32
	Qt_rest      sql.NullInt32
	Qt_rest_fact sql.NullInt32
	Qt_take      sql.NullInt32
	Qt_take_fact sql.NullInt32
	Qt_fill      sql.NullInt32
	Qt_fill_fact sql.NullInt32
	Qt_counter   sql.NullInt32
	Qt_pull      sql.NullInt32
}

func New() model_interface.Model {
	return &WareFlowProductModel{}
}

func (wfpm *WareFlowProductModel) New() model_interface.Model {
	return &WareFlowProductModel{}
}

func (wfpm *WareFlowProductModel) GetNameSchema(accountNumber int) string {
	return fmt.Sprintf("%s%v", schema, accountNumber)
}

func (wfpm *WareFlowProductModel) GetNameTable() string {
	return table
}

func (wfpm *WareFlowProductModel) GetFields() []string {
	fields := make([]string, 0, 1)
	v := reflect.ValueOf(wfpm).Elem()
	typeOfS := v.Type()
	for i := 0; i < v.NumField(); i++ {
		field := typeOfS.Field(i).Name
		fields = append(fields, field)
	}
	return fields
}

func (wfpm *WareFlowProductModel) GetIdKey() int64 {
	return int64(wfpm.Id)
}

func (wfpm *WareFlowProductModel) ScanModelRow(row pgx.Row) (model_interface.Model, error) {
	model := WareFlowProductModel{}
	err := row.Scan(&model.Id,
		&model.Wareflow_id,
		&model.Select_id,
		&model.Ware_id,
		&model.Qt_max,
		&model.Qt_recom,
		&model.Qt_rest,
		&model.Qt_rest_fact,
		&model.Qt_take,
		&model.Qt_take_fact,
		&model.Qt_fill,
		&model.Qt_fill_fact,
		&model.Qt_counter,
		&model.Qt_pull)
	if err != nil {
		return &model, err
	}
	return &model, nil
}

func (wfpm *WareFlowProductModel) ScanModelRows(rows pgx.Rows) ([]model_interface.Model, error) {
	defer rows.Close()
	Result := make([]model_interface.Model, 0, 1)
	for rows.Next() {
		resultStruct, err := wfpm.ScanModelRow(rows)
		if err != nil {
			continue
		}
		Result = append(Result, resultStruct)
	}
	return Result, nil
}

func (wfpm *WareFlowProductModel) GetName() string {
	return "WareFlowProductModel"
}

func (wfpm *WareFlowProductModel) Get() map[string]interface{} {
	model := make(map[string]interface{})
	v := reflect.ValueOf(wfpm).Elem()
	typeOfS := v.Type()
	for i := 0; i < v.NumField(); i++ {
		field, ok := wfpm.GetNameFieldOfTag(typeOfS.Field(i))
		if !ok {
			field = strings.ToLower(typeOfS.Field(i).Name)
		}
		FieldValue := v.FieldByName(field)
		model[field] = wfpm.GetValueField(&FieldValue)
	}
	return model
}

func (wfpm *WareFlowProductModel) GetNameFieldOfTag(field reflect.StructField) (string, bool) {
	if value, ok := field.Tag.Lookup("json"); ok {
		return value, ok
	} else {
		return "", false
	}

}

func (wfpm *WareFlowProductModel) GetValueField(f *reflect.Value) interface{} {
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
	return false
}
