package model

import (
	"database/sql"
	model_interface "ephorservices/pkg/orm/model"
	"fmt"
	"reflect"
	"strings"

	pgx "github.com/jackc/pgx/v5"
)

var Account = "account"

type WareHouseModel struct {
	Id        int            `json:"id"`
	W_1c_id   sql.NullString `json:"w_1c_id"`
	Date_sync sql.NullString `json:"date_sync"`
	Name      sql.NullString `json:"name"`
	Address   sql.NullString `json:"address"`
	Comment   sql.NullString `json:"comment"`
	State     sql.NullInt32  `json:"state"`
}

func New() model_interface.Model {
	return &WareHouseModel{}
}

func (whm *WareHouseModel) New() model_interface.Model {
	return &WareHouseModel{}
}

func (whm *WareHouseModel) GetNameSchema(accountNumber int) string {
	return fmt.Sprintf("%s%v", Account, accountNumber)
}

func (whm *WareHouseModel) GetNameTable() string {
	return "warehouse"
}

func (whm *WareHouseModel) GetFields() []string {
	fields := make([]string, 0, 1)
	v := reflect.ValueOf(whm).Elem()
	typeOfS := v.Type()
	for i := 0; i < v.NumField(); i++ {
		field := typeOfS.Field(i).Name
		fields = append(fields, field)
	}
	return fields
}

func (whm *WareHouseModel) GetIdKey() int64 {
	return int64(whm.Id)
}

func (whm *WareHouseModel) ScanModelRow(row pgx.Row) (model_interface.Model, error) {
	model := WareHouseModel{}
	err := row.Scan(&model.Id,
		&model.W_1c_id,
		&model.Date_sync,
		&model.Name,
		&model.Address,
		&model.Comment,
		&model.State)
	if err != nil {
		return &model, err
	}
	return &model, nil
}

func (whm *WareHouseModel) ScanModelRows(rows pgx.Rows) ([]model_interface.Model, error) {
	defer rows.Close()
	Result := make([]model_interface.Model, 0, 1)
	for rows.Next() {
		resultStruct, err := whm.ScanModelRow(rows)
		if err != nil {
			continue
		}
		Result = append(Result, resultStruct)
	}
	return Result, nil
}

func (whm *WareHouseModel) GetName() string {
	return "WareHouseModel"
}

func (whm *WareHouseModel) Get() map[string]interface{} {
	model := make(map[string]interface{})
	v := reflect.ValueOf(whm).Elem()
	typeOfS := v.Type()
	for i := 0; i < v.NumField(); i++ {
		field, ok := whm.GetNameFieldOfTag(typeOfS.Field(i))
		if !ok {
			field = strings.ToLower(typeOfS.Field(i).Name)
		}
		FieldValue := v.FieldByName(field)
		model[field] = whm.GetValueField(&FieldValue)
	}
	return model
}

func (whm *WareHouseModel) GetNameFieldOfTag(field reflect.StructField) (string, bool) {
	if value, ok := field.Tag.Lookup("json"); ok {
		return value, ok
	} else {
		return "", false
	}

}

func (whm *WareHouseModel) GetValueField(f *reflect.Value) interface{} {
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
