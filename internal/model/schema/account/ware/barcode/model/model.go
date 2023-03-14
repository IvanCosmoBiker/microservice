package barcode

import (
	"database/sql"
	model_interface "ephorservices/pkg/orm/model"
	"fmt"
	"reflect"
	"strings"

	pgx "github.com/jackc/pgx/v5"
)

var schema = "account"
var table = "ware_barcode"

type WareBarcodeModel struct {
	Id      int
	Ware_id sql.NullInt32
	Barcode sql.NullString
}

func New() model_interface.Model {
	return &WareBarcodeModel{}
}

func (wm *WareBarcodeModel) New() model_interface.Model {
	return &WareBarcodeModel{}
}

func (wm *WareBarcodeModel) GetNameSchema(accountNumber int) string {
	return fmt.Sprintf("%s%v", schema, accountNumber)
}

func (wm *WareBarcodeModel) GetNameTable() string {
	return table
}

func (wm *WareBarcodeModel) GetFields() []string {
	fields := make([]string, 0, 1)
	v := reflect.ValueOf(wm).Elem()
	typeOfS := v.Type()
	for i := 0; i < v.NumField(); i++ {
		field := typeOfS.Field(i).Name
		fields = append(fields, field)
	}
	return fields
}

func (wm *WareBarcodeModel) GetIdKey() int64 {
	return int64(wm.Id)
}

func (wm *WareBarcodeModel) ScanModelRow(row pgx.Row) (model_interface.Model, error) {
	model := WareBarcodeModel{}
	err := row.Scan(&model.Id,
		&model.Ware_id,
		&model.Barcode)
	if err != nil {
		return &model, err
	}
	return &model, nil
}

func (wm *WareBarcodeModel) ScanModelRows(rows pgx.Rows) ([]model_interface.Model, error) {
	defer rows.Close()
	Result := make([]model_interface.Model, 0, 1)
	for rows.Next() {
		resultStruct, err := wm.ScanModelRow(rows)
		if err != nil {
			continue
		}
		Result = append(Result, resultStruct)
	}
	return Result, nil
}

func (wm *WareBarcodeModel) GetName() string {
	return "WareBarcodeModel"
}

func (wm *WareBarcodeModel) Get() map[string]interface{} {
	model := make(map[string]interface{})
	v := reflect.ValueOf(wm).Elem()
	typeOfS := v.Type()
	for i := 0; i < v.NumField(); i++ {
		field, ok := wm.GetNameFieldOfTag(typeOfS.Field(i))
		if !ok {
			field = strings.ToLower(typeOfS.Field(i).Name)
		}
		FieldValue := v.FieldByName(field)
		model[field] = wm.GetValueField(&FieldValue)
	}
	return model
}

func (wm *WareBarcodeModel) GetNameFieldOfTag(field reflect.StructField) (string, bool) {
	if value, ok := field.Tag.Lookup("json"); ok {
		return value, ok
	} else {
		return "", false
	}

}

func (wm *WareBarcodeModel) GetValueField(f *reflect.Value) interface{} {
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
