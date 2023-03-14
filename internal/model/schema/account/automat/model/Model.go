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
var table = "automat"

type AutomatModel struct {
	Id                    int
	Automat_model_id      sql.NullInt32
	Modem_id              sql.NullInt32
	Fr_id                 sql.NullInt32
	Pay_type              sql.NullInt32
	Work_mode             sql.NullInt32
	Sbp_id                sql.NullInt32
	Ext1                  sql.NullInt32
	Serial_number         sql.NullString
	Key                   sql.NullString
	Production_date       sql.NullString
	From_date             sql.NullString
	To_date               sql.NullString
	Update_date           sql.NullString
	Last_sale             sql.NullString
	Last_audit            sql.NullString
	Last_encash           sql.NullString
	Type_nosale           sql.NullInt32
	Type_service          sql.NullInt32
	Type_encashment       sql.NullInt32
	Now_cash_val          sql.NullInt64
	Now_cashless_val      sql.NullInt64
	Now_token_val         sql.NullInt64
	Tube_val              sql.NullInt64
	Now_tube_val          sql.NullInt64
	Control_billvalidator sql.NullInt32
	Control_coinchanger   sql.NullInt32
	Control_cashless      sql.NullInt32
	Last_coin             sql.NullString
	Last_bill             sql.NullString
	Last_cashless         sql.NullString
	Load_date             sql.NullString
	Update_config_id      sql.NullString
	Now_cash_num          sql.NullInt32
	Now_cashless_num      sql.NullInt32
	Now_token_num         sql.NullInt32
	Cash_error            sql.NullInt32
	Cashless_error        sql.NullInt32
	Token_error           sql.NullInt32
	Qr                    sql.NullInt32
	Qr_type               sql.NullInt32
	Ext2                  sql.NullInt32
	Usb1                  sql.NullInt32
	Internet              sql.NullInt32
	Ethernet_mac          sql.NullString
	Ethernet_ip           sql.NullString
	Ethernet_netmask      sql.NullString
	Ethernet_gateway      sql.NullString
	Faceid_type           sql.NullInt32
	Faceid_id             sql.NullString
	Faceid_addr           sql.NullString
	Faceid_port           sql.NullInt32
	Faceid_username       sql.NullString
	Faceid_password       sql.NullString
	Summ_max_fr           sql.NullInt64
	Last_login            sql.NullString
	Warehouse_id          sql.NullString
	Screen_text           sql.NullString
}

func New() model_interface.Model {
	return &AutomatModel{}
}

func (am *AutomatModel) New() model_interface.Model {
	return &AutomatModel{}
}

func (am *AutomatModel) GetNameSchema(accountNumber int) string {
	return fmt.Sprintf("%s%v", Account, accountNumber)
}

func (am *AutomatModel) GetNameTable() string {
	return table
}

func (am *AutomatModel) GetFields() []string {
	fields := make([]string, 0, 1)
	v := reflect.ValueOf(am).Elem()
	typeOfS := v.Type()
	for i := 0; i < v.NumField(); i++ {
		field := typeOfS.Field(i).Name
		fields = append(fields, field)
	}
	return fields
}

func (am *AutomatModel) GetIdKey() int64 {
	return int64(am.Id)
}

func (am *AutomatModel) ScanModelRow(row pgx.Row) (model_interface.Model, error) {
	Model := AutomatModel{}
	err := row.Scan(&Model.Id,
		&Model.Automat_model_id,
		&Model.Modem_id,
		&Model.Fr_id,
		&Model.Pay_type,
		&Model.Work_mode,
		&Model.Sbp_id,
		&Model.Ext1,
		&Model.Serial_number,
		&Model.Key,
		&Model.Production_date,
		&Model.From_date,
		&Model.To_date,
		&Model.Update_date,
		&Model.Last_sale,
		&Model.Last_audit,
		&Model.Last_encash,
		&Model.Type_nosale,
		&Model.Type_service,
		&Model.Type_encashment,
		&Model.Now_cash_val,
		&Model.Now_cashless_val,
		&Model.Now_token_val,
		&Model.Tube_val,
		&Model.Now_tube_val,
		&Model.Control_billvalidator,
		&Model.Control_coinchanger,
		&Model.Control_cashless,
		&Model.Last_coin,
		&Model.Last_bill,
		&Model.Last_cashless,
		&Model.Load_date,
		&Model.Update_config_id,
		&Model.Now_cash_num,
		&Model.Now_cashless_num,
		&Model.Now_token_num,
		&Model.Cash_error,
		&Model.Cashless_error,
		&Model.Token_error,
		&Model.Qr,
		&Model.Qr_type,
		&Model.Ext2,
		&Model.Usb1,
		&Model.Internet,
		&Model.Ethernet_mac,
		&Model.Ethernet_ip,
		&Model.Ethernet_netmask,
		&Model.Ethernet_gateway,
		&Model.Faceid_type,
		&Model.Faceid_id,
		&Model.Faceid_addr,
		&Model.Faceid_port,
		&Model.Faceid_username,
		&Model.Faceid_password,
		&Model.Summ_max_fr,
		&Model.Last_login,
		&Model.Warehouse_id,
		&Model.Screen_text)
	if err != nil {
		return &Model, err
	}
	return &Model, nil
}

func (am *AutomatModel) ScanModelRows(rows pgx.Rows) ([]model_interface.Model, error) {
	defer rows.Close()
	Result := make([]model_interface.Model, 0, 1)
	for rows.Next() {
		resultStruct, err := am.ScanModelRow(rows)
		if err != nil {
			continue
		}
		Result = append(Result, resultStruct)
	}
	return Result, nil
}

func (am *AutomatModel) GetName() string {
	return "AutomatModel"
}

func (am *AutomatModel) Get() map[string]interface{} {
	model := make(map[string]interface{})
	v := reflect.ValueOf(am).Elem()
	typeOfS := v.Type()
	for i := 0; i < v.NumField(); i++ {
		field, ok := am.GetNameFieldOfTag(typeOfS.Field(i))
		if !ok {
			field = strings.ToLower(typeOfS.Field(i).Name)
		}
		FieldValue := v.FieldByName(field)
		model[field] = am.GetValueField(&FieldValue)
	}
	return model
}

func (am *AutomatModel) GetNameFieldOfTag(field reflect.StructField) (string, bool) {
	if value, ok := field.Tag.Lookup("json"); ok {
		return value, ok
	} else {
		return "", false
	}

}

func (am *AutomatModel) GetValueField(f *reflect.Value) interface{} {
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
