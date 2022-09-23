package model

import (
	"database/sql"
	"encoding/json"
	modelInterface "ephorservices/pkg/model/schema/interface/model"
	"log"

	pgx "github.com/jackc/pgx/v4"
)

var nameTable string = "modem"
var nameSchema string = "main"

type ModemModel struct {
	Id               int
	Account_id       int
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
	Static           sql.NullInt32
	Gsm_apn          sql.NullString
	Gsm_username     sql.NullString
	Gsm_password     sql.NullString
	Dns1             sql.NullString
	Dns2             sql.NullString
	Add_date         sql.NullString
}

func Init() modelInterface.Model {
	return &ModemModel{}
}

func (cm *ModemModel) InitData(jsonData []byte) {
	json.Unmarshal(jsonData, &cm)
}

func (cm *ModemModel) JsonSerialize() []byte {
	data, _ := json.Marshal(cm)
	return data
}

func (cm *ModemModel) GetNameTable(account string) string {
	return nameTable
}

func (cm *ModemModel) GetNameSchema(account int) string {
	return nameSchema
}

type ReturningStruct struct {
	Id               int
	Account_id       int
	Imei             string
	Hash             string
	Nonce            string
	Hardware_version int32
	Software_version int32
	Software_release int32
	Phone            string
	Signal_quality   int32
	Last_login       string
	Last_ex_id       int32
	Ipaddr           string
	Static           int32
	Gsm_apn          string
	Gsm_username     string
	Gsm_password     string
	Dns1             string
	Dns2             string
	Add_date         string
}

func ScanModelRows(rows pgx.Rows) ([]*ReturningStruct, error) {
	defer rows.Close()
	Result := []*ReturningStruct{}
	for rows.Next() {
		resultStruct, err := ScanModelRow(rows)
		if err != nil {
			log.Println(err)
			continue
		}
		Result = append(Result, resultStruct)
	}
	return Result, nil
}

func ScanModelRow(row pgx.Row) (*ReturningStruct, error) {
	Model := ModemModel{}
	err := row.Scan(&Model.Id,
		&Model.Account_id,
		&Model.Imei,
		&Model.Hash,
		&Model.Nonce,
		&Model.Hardware_version,
		&Model.Software_version,
		&Model.Software_release,
		&Model.Phone,
		&Model.Signal_quality,
		&Model.Last_login,
		&Model.Last_ex_id,
		&Model.Ipaddr,
		&Model.Ipaddr,
		&Model.Static,
		&Model.Gsm_apn,
		&Model.Gsm_username,
		&Model.Gsm_password,
		&Model.Dns1,
		&Model.Dns2,
		&Model.Add_date)
	if err != nil {
		log.Printf("%v", err)
		return &ReturningStruct{}, err
	}
	Result := MakeDataInSturct(&Model)
	return Result, nil
}

func MakeDataInSturct(modelData *ModemModel) *ReturningStruct {
	stuctReturn := ReturningStruct{}
	stuctReturn.Id = modelData.Id
	stuctReturn.Account_id = modelData.Account_id
	stuctReturn.Imei = modelData.Imei.String
	stuctReturn.Hash = modelData.Hash.String
	stuctReturn.Nonce = modelData.Nonce.String
	stuctReturn.Hardware_version = modelData.Hardware_version.Int32
	stuctReturn.Software_version = modelData.Software_version.Int32
	stuctReturn.Software_release = modelData.Software_release.Int32
	stuctReturn.Phone = modelData.Phone.String
	stuctReturn.Signal_quality = modelData.Signal_quality.Int32
	stuctReturn.Last_login = modelData.Last_login.String
	stuctReturn.Last_ex_id = modelData.Last_ex_id.Int32
	stuctReturn.Ipaddr = modelData.Ipaddr.String
	stuctReturn.Static = modelData.Static.Int32
	stuctReturn.Gsm_apn = modelData.Gsm_apn.String
	stuctReturn.Gsm_username = modelData.Gsm_username.String
	stuctReturn.Gsm_password = modelData.Gsm_password.String
	stuctReturn.Dns1 = modelData.Dns1.String
	stuctReturn.Dns2 = modelData.Dns2.String
	stuctReturn.Add_date = modelData.Add_date.String
	return &stuctReturn
}
