package model

import (
	"database/sql"
	model_interface "ephorservices/pkg/orm/model"
	"reflect"
	"strings"

	pgx "github.com/jackc/pgx/v5"
)

var Schema string = "main"
var nameTable string = "modem"

type ModemModel struct {
	Id               int
	Account_id       sql.NullInt32
	Imei             sql.NullString
	Hash             sql.NullString
	Nonce            sql.NullString
	Hardware_version sql.NullInt32
	Software_version sql.NullInt32
	Software_release sql.NullInt32
	Phone            sql.NullString
	Signal_quality   sql.NullInt32
	Last_login       sql.NullString
	Last_ex_id       sql.NullInt32
	Ipaddr           sql.NullString
	On_log           sql.NullInt32
	Gsm_apn          sql.NullString
	Gsm_username     sql.NullString
	Gsm_password     sql.NullString
	Dns1             sql.NullString
	Dns2             sql.NullString
	Add_date         sql.NullString
	Type             sql.NullInt32
}

func New() model_interface.Model {
	return &ModemModel{}
}

func (mm *ModemModel) New() model_interface.Model {
	return &ModemModel{}
}

func (mm *ModemModel) GetNameSchema(accountNumber int) string {
	return Schema
}

func (mm *ModemModel) GetNameTable() string {
	return nameTable
}

func (mm *ModemModel) GetFields() []string {
	fields := make([]string, 0, 1)
	v := reflect.ValueOf(mm).Elem()
	typeOfS := v.Type()
	for i := 0; i < v.NumField(); i++ {
		field := typeOfS.Field(i).Name
		fields = append(fields, field)
	}
	return fields
}

func (mm *ModemModel) GetIdKey() int64 {
	return int64(mm.Id)
}

func (mm *ModemModel) ScanModelRow(row pgx.Row) (model_interface.Model, error) {
	model := ModemModel{}
	err := row.Scan(&model.Id,
		&model.Account_id,
		&model.Imei,
		&model.Hash,
		&model.Nonce,
		&model.Hardware_version,
		&model.Software_version,
		&model.Software_release,
		&model.Phone,
		&model.Signal_quality,
		&model.Last_login,
		&model.Last_ex_id,
		&model.Ipaddr,
		&model.On_log,
		&model.Gsm_apn,
		&model.Gsm_username,
		&model.Gsm_password,
		&model.Dns1,
		&model.Dns2,
		&model.Add_date,
		&model.Type)
	if err != nil {
		return &model, err
	}
	return &model, nil
}

func (mm *ModemModel) ScanModelRows(rows pgx.Rows) ([]model_interface.Model, error) {
	defer rows.Close()
	Result := make([]model_interface.Model, 0, 1)
	for rows.Next() {
		resultStruct, err := mm.ScanModelRow(rows)
		if err != nil {
			continue
		}
		Result = append(Result, resultStruct)
	}
	return Result, nil
}

func (mm *ModemModel) GetName() string {
	return "ModemModel"
}

func (mm *ModemModel) Get() map[string]interface{} {
	model := make(map[string]interface{})
	v := reflect.ValueOf(mm).Elem()
	typeOfS := v.Type()
	for i := 0; i < v.NumField(); i++ {
		field, ok := mm.GetNameFieldOfTag(typeOfS.Field(i))
		if !ok {
			field = strings.ToLower(typeOfS.Field(i).Name)
		}
		FieldValue := v.FieldByName(field)
		model[field] = mm.GetValueField(&FieldValue)
	}
	return model
}

func (mm *ModemModel) GetNameFieldOfTag(field reflect.StructField) (string, bool) {
	if value, ok := field.Tag.Lookup("json"); ok {
		return value, ok
	} else {
		return "", false
	}

}

func (mm *ModemModel) GetValueField(f *reflect.Value) interface{} {
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
