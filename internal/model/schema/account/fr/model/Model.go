package fr

import (
	"database/sql"
	model_interface "ephorservices/pkg/orm/model"
	"fmt"
	"reflect"
	"strings"

	pgx "github.com/jackc/pgx/v5"
)

var schema = "account"
var table = "fr"

type FrModel struct {
	Id                  int
	Name                sql.NullString
	Type                sql.NullInt32
	Dev_interface       sql.NullInt32
	Login               sql.NullString
	Password            sql.NullString
	Phone               sql.NullString
	Email               sql.NullString
	Dev_addr            sql.NullString
	Dev_port            sql.NullInt32
	Ofd_addr            sql.NullString
	Ofd_port            sql.NullInt32
	Inn                 sql.NullString
	Auth_public_key     sql.NullString
	Auth_private_key    sql.NullString
	Sign_private_key    sql.NullString
	Param1              sql.NullString
	Use_sn              sql.NullInt32
	Add_fiscal          sql.NullInt32
	Id_shift            sql.NullString
	Fr_disable_cash     sql.NullInt32
	Fr_disable_cashless sql.NullInt32
	Ffd_version         sql.NullInt32
}

func New() model_interface.Model {
	return &FrModel{}
}

func (fm *FrModel) New() model_interface.Model {
	return &FrModel{}
}

func (fm *FrModel) GetNameSchema(accountNumber int) string {
	return fmt.Sprintf("%s%v", schema, accountNumber)
}

func (fm *FrModel) GetNameTable() string {
	return table
}

func (fm *FrModel) GetFields() []string {
	fields := make([]string, 0, 1)
	v := reflect.ValueOf(fm).Elem()
	typeOfS := v.Type()
	for i := 0; i < v.NumField(); i++ {
		field := typeOfS.Field(i).Name
		fields = append(fields, field)
	}
	return fields
}

func (fm *FrModel) GetIdKey() int64 {
	return int64(fm.Id)
}

func (fm *FrModel) ScanModelRow(row pgx.Row) (model_interface.Model, error) {
	model := FrModel{}
	err := row.Scan(&model.Id,
		&model.Name,
		&model.Type,
		&model.Dev_interface,
		&model.Login,
		&model.Password,
		&model.Phone,
		&model.Email,
		&model.Dev_addr,
		&model.Dev_port,
		&model.Ofd_addr,
		&model.Ofd_port,
		&model.Inn,
		&model.Auth_public_key,
		&model.Auth_private_key,
		&model.Sign_private_key,
		&model.Param1,
		&model.Use_sn,
		&model.Add_fiscal,
		&model.Id_shift,
		&model.Fr_disable_cash,
		&model.Fr_disable_cashless,
		&model.Ffd_version)
	if err != nil {
		return &model, err
	}
	return &model, nil
}

func (fm *FrModel) ScanModelRows(rows pgx.Rows) ([]model_interface.Model, error) {
	defer rows.Close()
	Result := make([]model_interface.Model, 0, 1)
	for rows.Next() {
		resultStruct, err := fm.ScanModelRow(rows)
		if err != nil {
			continue
		}
		Result = append(Result, resultStruct)
	}
	return Result, nil
}

func (fm *FrModel) GetName() string {
	return "FrModel"
}

func (fm *FrModel) Get() map[string]interface{} {
	model := make(map[string]interface{})
	v := reflect.ValueOf(fm).Elem()
	typeOfS := v.Type()
	for i := 0; i < v.NumField(); i++ {
		field, ok := fm.GetNameFieldOfTag(typeOfS.Field(i))
		if !ok {
			field = strings.ToLower(typeOfS.Field(i).Name)
		}
		FieldValue := v.FieldByName(field)
		model[field] = fm.GetValueField(&FieldValue)
	}
	return model
}

func (fm *FrModel) GetNameFieldOfTag(field reflect.StructField) (string, bool) {
	if value, ok := field.Tag.Lookup("json"); ok {
		return value, ok
	} else {
		return "", false
	}

}

func (fm *FrModel) GetValueField(f *reflect.Value) interface{} {
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
