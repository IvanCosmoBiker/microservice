package model

import (
	"database/sql"
	"encoding/json"
	modelInterface "ephorservices/pkg/model/schema/interface/model"
	"log"

	pgx "github.com/jackc/pgx/v4"
)

var nameTable string = "automat"
var nameSchema string = "account"

type AutomatModel struct {
	Id                    int
	Automat_model_id      sql.NullInt32
	Modem_id              sql.NullInt32
	Fr_id                 sql.NullInt32
	Pay_type              sql.NullInt32
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
	Now_cash_val          sql.NullInt32
	Now_cashless_val      sql.NullInt32
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
	Cash_error            sql.NullInt32
	Cashless_error        sql.NullInt32
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

func Init() modelInterface.Model {
	return &AutomatModel{}
}

func (am *AutomatModel) InitData(jsonData []byte) {
	json.Unmarshal(jsonData, &am)
}

func (am *AutomatModel) JsonSerialize() []byte {
	data, _ := json.Marshal(am)
	return data
}

func (am *AutomatModel) GetNameTable() string {
	return nameTable
}

func (am AutomatModel) GetNameSchema(accountId string) string {
	return nameSchema + "" + accountId
}

type ReturningStruct struct {
	Id                    int
	Automat_model_id      int32
	Modem_id              int32
	Fr_id                 int32
	Pay_type              int32
	Sbp_id                int32
	Ext1                  int32
	Serial_number         string
	Key                   string
	Production_date       string
	From_date             string
	To_date               string
	Update_date           string
	Last_sale             string
	Last_audit            string
	Last_encash           string
	Type_nosale           int32
	Type_service          int32
	Type_encashment       int32
	Now_cash_val          int32
	Now_cashless_val      int32
	Tube_val              int64
	Now_tube_val          int64
	Control_billvalidator int32
	Control_coinchanger   int32
	Control_cashless      int32
	Last_coin             string
	Last_bill             string
	Last_cashless         string
	Load_date             string
	Update_config_id      string
	Now_cash_num          int32
	Now_cashless_num      int32
	Cash_error            int32
	Cashless_error        int32
	Qr                    int32
	Qr_type               int32
	Ext2                  int32
	Usb1                  int32
	Internet              int32
	Ethernet_mac          string
	Ethernet_ip           string
	Ethernet_netmask      string
	Ethernet_gateway      string
	Faceid_type           int32
	Faceid_id             string
	Faceid_addr           string
	Faceid_port           int32
	Faceid_username       string
	Faceid_password       string
	Summ_max_fr           int64
	Last_login            string
	Warehouse_id          string
	Screen_text           string
}

func ScanModelRows(rows pgx.Rows) ([]*ReturningStruct, error) {
	defer rows.Close()
	Result := []*ReturningStruct{}
	for rows.Next() {
		resultStruct, err := ScanModelRow(rows)
		if err != nil {
			log.Printf("%v", err)
			continue
		}
		Result = append(Result, resultStruct)
	}
	return Result, nil
}

func ScanModelRow(row pgx.Row) (*ReturningStruct, error) {
	Model := AutomatModel{}
	err := row.Scan(&Model.Id,
		&Model.Automat_model_id,
		&Model.Modem_id,
		&Model.Fr_id,
		&Model.Pay_type,
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
		&Model.Cash_error,
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
		log.Printf("%v\n", err)
		return &ReturningStruct{}, err
	}
	Result := MakeDataInSturct(&Model)

	return Result, nil
}

func MakeDataInSturct(modelData *AutomatModel) *ReturningStruct {
	stuctReturn := ReturningStruct{}
	stuctReturn.Id = modelData.Id
	stuctReturn.Automat_model_id = modelData.Automat_model_id.Int32
	stuctReturn.Modem_id = modelData.Modem_id.Int32
	stuctReturn.Fr_id = modelData.Fr_id.Int32
	stuctReturn.Pay_type = modelData.Pay_type.Int32
	stuctReturn.Sbp_id = modelData.Sbp_id.Int32
	stuctReturn.Ext1 = modelData.Ext1.Int32
	stuctReturn.Serial_number = modelData.Serial_number.String
	stuctReturn.Key = modelData.Key.String
	stuctReturn.Production_date = modelData.Production_date.String
	stuctReturn.From_date = modelData.From_date.String
	stuctReturn.To_date = modelData.To_date.String
	stuctReturn.Update_date = modelData.Update_date.String
	stuctReturn.Last_sale = modelData.Last_sale.String
	stuctReturn.Last_audit = modelData.Last_audit.String
	stuctReturn.Last_encash = modelData.Last_encash.String
	stuctReturn.Type_nosale = modelData.Type_nosale.Int32
	stuctReturn.Type_service = modelData.Type_service.Int32
	stuctReturn.Type_encashment = modelData.Type_encashment.Int32
	stuctReturn.Now_cash_val = modelData.Now_cash_val.Int32
	stuctReturn.Now_cashless_val = modelData.Now_cashless_val.Int32
	stuctReturn.Tube_val = modelData.Tube_val.Int64
	stuctReturn.Now_tube_val = modelData.Now_tube_val.Int64
	stuctReturn.Control_billvalidator = modelData.Control_billvalidator.Int32
	stuctReturn.Control_coinchanger = modelData.Control_coinchanger.Int32
	stuctReturn.Control_cashless = modelData.Control_cashless.Int32
	stuctReturn.Last_coin = modelData.Last_coin.String
	stuctReturn.Last_bill = modelData.Last_bill.String
	stuctReturn.Last_cashless = modelData.Last_cashless.String
	stuctReturn.Load_date = modelData.Load_date.String
	stuctReturn.Update_config_id = modelData.Update_config_id.String
	stuctReturn.Now_cash_num = modelData.Now_cash_num.Int32
	stuctReturn.Now_cashless_num = modelData.Now_cashless_num.Int32
	stuctReturn.Cash_error = modelData.Cash_error.Int32
	stuctReturn.Qr = modelData.Qr.Int32
	stuctReturn.Qr_type = modelData.Qr_type.Int32
	stuctReturn.Ext2 = modelData.Ext2.Int32
	stuctReturn.Usb1 = modelData.Usb1.Int32
	stuctReturn.Internet = modelData.Internet.Int32
	stuctReturn.Ethernet_mac = modelData.Ethernet_mac.String
	stuctReturn.Ethernet_ip = modelData.Ethernet_ip.String
	stuctReturn.Ethernet_netmask = modelData.Ethernet_netmask.String
	stuctReturn.Ethernet_gateway = modelData.Ethernet_gateway.String
	stuctReturn.Faceid_type = modelData.Faceid_type.Int32
	stuctReturn.Faceid_id = modelData.Faceid_id.String
	stuctReturn.Faceid_addr = modelData.Faceid_addr.String
	stuctReturn.Faceid_port = modelData.Faceid_port.Int32
	stuctReturn.Faceid_username = modelData.Faceid_username.String
	stuctReturn.Faceid_password = modelData.Faceid_password.String
	stuctReturn.Summ_max_fr = modelData.Summ_max_fr.Int64
	stuctReturn.Last_login = modelData.Last_login.String
	stuctReturn.Warehouse_id = modelData.Warehouse_id.String
	stuctReturn.Screen_text = modelData.Screen_text.String
	return &stuctReturn
}
