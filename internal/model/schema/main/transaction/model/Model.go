package model

import (
	"database/sql"
	model_interface "ephorservices/pkg/orm/model"
	"reflect"
	"strings"

	pgx "github.com/jackc/pgx/v5"
)

var Schema string = "main"
var nameTable string = "transaction"

type TransactionModel struct {
	Id            int
	Noise         sql.NullString
	Token_id      sql.NullString
	Token_type    sql.NullInt64
	Account_id    sql.NullInt64
	Customer_id   sql.NullInt64
	Automat_id    sql.NullInt64
	Date          sql.NullString
	Sum           sql.NullInt64
	Ps_type       sql.NullInt64
	Ps_order      sql.NullString
	Ps_tid        sql.NullString
	Ps_code       sql.NullString
	Ps_desc       sql.NullString
	Ps_invoice_id sql.NullString
	Pay_type      sql.NullInt64
	Fn            sql.NullInt64
	Fd            sql.NullInt64
	Fp            sql.NullString
	F_type        sql.NullInt64
	F_receipt     sql.NullString
	F_desc        sql.NullString
	F_status      sql.NullInt64
	Qr_format     sql.NullInt64
	F_qr          sql.NullString
	Status        sql.NullInt64
	Error         sql.NullString
}

func New() model_interface.Model {
	return &TransactionModel{}
}

func (tm *TransactionModel) New() model_interface.Model {
	return &TransactionModel{}
}

func (tm *TransactionModel) GetNameSchema(accountNumber int) string {
	return Schema
}

func (tm *TransactionModel) GetNameTable() string {
	return nameTable
}

func (tm *TransactionModel) GetFields() []string {
	fields := make([]string, 0, 1)
	v := reflect.ValueOf(tm).Elem()
	typeOfS := v.Type()
	for i := 0; i < v.NumField(); i++ {
		field := typeOfS.Field(i).Name
		fields = append(fields, field)
	}
	return fields
}

func (tm *TransactionModel) GetIdKey() int64 {
	return int64(tm.Id)
}

func (tm *TransactionModel) ScanModelRow(row pgx.Row) (model_interface.Model, error) {
	model := TransactionModel{}
	err := row.Scan(&model.Id,
		&model.Noise,
		&model.Token_id,
		&model.Token_type,
		&model.Account_id,
		&model.Customer_id,
		&model.Automat_id,
		&model.Date,
		&model.Sum,
		&model.Ps_type,
		&model.Ps_order,
		&model.Ps_tid,
		&model.Ps_code,
		&model.Ps_desc,
		&model.Ps_invoice_id,
		&model.Pay_type,
		&model.Fn,
		&model.Fd,
		&model.Fp,
		&model.F_type,
		&model.F_receipt,
		&model.F_desc,
		&model.F_status,
		&model.Qr_format,
		&model.F_qr,
		&model.Status,
		&model.Error)
	if err != nil {
		return &model, err
	}
	return &model, nil
}

func (tm *TransactionModel) ScanModelRows(rows pgx.Rows) ([]model_interface.Model, error) {
	defer rows.Close()
	Result := make([]model_interface.Model, 0, 1)
	for rows.Next() {
		resultStruct, err := tm.ScanModelRow(rows)
		if err != nil {
			continue
		}
		Result = append(Result, resultStruct)
	}
	return Result, nil
}

func (tm *TransactionModel) GetName() string {
	return "TransactionModel"
}

func (tm *TransactionModel) Get() map[string]interface{} {
	model := make(map[string]interface{})
	v := reflect.ValueOf(tm).Elem()
	typeOfS := v.Type()
	for i := 0; i < v.NumField(); i++ {
		field, ok := tm.GetNameFieldOfTag(typeOfS.Field(i))
		if !ok {
			field = strings.ToLower(typeOfS.Field(i).Name)
		}
		FieldValue := v.FieldByName(field)
		model[field] = tm.GetValueField(&FieldValue)
	}
	return model
}

func (tm *TransactionModel) GetNameFieldOfTag(field reflect.StructField) (string, bool) {
	if value, ok := field.Tag.Lookup("json"); ok {
		return value, ok
	} else {
		return "", false
	}

}

func (tm *TransactionModel) GetValueField(f *reflect.Value) interface{} {
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
