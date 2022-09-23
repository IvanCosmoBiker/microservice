package fr

import (
	"database/sql"
	"encoding/json"
	modelInterface "ephorservices/pkg/model/schema/interface/model"
	"log"

	pgx "github.com/jackc/pgx/v4"
)

var nameTable string = "fr"
var nameSchema string = "account"

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

func Init() modelInterface.Model {
	return &FrModel{}
}

func (fm *FrModel) InitData(jsonData []byte) {
	json.Unmarshal(jsonData, &fm)
}

func (fm *FrModel) JsonSerialize() []byte {
	data, _ := json.Marshal(fm)
	return data
}

func (fm *FrModel) GetNameTable() string {
	return nameTable
}

func (fm *FrModel) GetNameSchema(accountId string) string {
	return nameSchema + "" + accountId
}

type ReturningStruct struct {
	Id                  int
	Name                string
	Type                int32
	Dev_interface       int32
	Login               string
	Password            string
	Phone               string
	Email               string
	Dev_addr            string
	Dev_port            int32
	Ofd_addr            string
	Ofd_port            int32
	Inn                 string
	Auth_public_key     string
	Auth_private_key    string
	Sign_private_key    string
	Param1              string
	Use_sn              int32
	Add_fiscal          int32
	Id_shift            string
	Fr_disable_cash     int32
	Fr_disable_cashless int32
	Ffd_version         int32
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
	Model := FrModel{}
	err := row.Scan(&Model.Id,
		&Model.Name,
		&Model.Type,
		&Model.Dev_interface,
		&Model.Login,
		&Model.Password,
		&Model.Phone,
		&Model.Email,
		&Model.Dev_addr,
		&Model.Dev_port,
		&Model.Ofd_addr,
		&Model.Ofd_port,
		&Model.Inn,
		&Model.Auth_public_key,
		&Model.Auth_private_key,
		&Model.Sign_private_key,
		&Model.Param1,
		&Model.Use_sn,
		&Model.Add_fiscal,
		&Model.Id_shift,
		&Model.Fr_disable_cash,
		&Model.Fr_disable_cashless)
	if err != nil {
		log.Printf("%v\n", err)
		return &ReturningStruct{}, err
	}
	Result := MakeDataInSturct(&Model)
	return Result, nil
}

func MakeDataInSturct(modelData *FrModel) *ReturningStruct {
	stuctReturn := ReturningStruct{}
	stuctReturn.Id = modelData.Id
	stuctReturn.Name = modelData.Name.String
	stuctReturn.Type = modelData.Type.Int32
	stuctReturn.Dev_interface = modelData.Dev_interface.Int32
	stuctReturn.Login = modelData.Login.String
	stuctReturn.Password = modelData.Password.String
	stuctReturn.Phone = modelData.Phone.String
	stuctReturn.Email = modelData.Email.String
	stuctReturn.Dev_addr = modelData.Dev_addr.String
	stuctReturn.Dev_port = modelData.Dev_port.Int32
	stuctReturn.Ofd_addr = modelData.Ofd_addr.String
	stuctReturn.Ofd_port = modelData.Ofd_port.Int32
	stuctReturn.Inn = modelData.Inn.String
	stuctReturn.Auth_public_key = modelData.Auth_public_key.String
	stuctReturn.Auth_private_key = modelData.Auth_private_key.String
	stuctReturn.Sign_private_key = modelData.Sign_private_key.String
	stuctReturn.Param1 = modelData.Param1.String
	stuctReturn.Use_sn = modelData.Use_sn.Int32
	stuctReturn.Add_fiscal = modelData.Add_fiscal.Int32
	stuctReturn.Id_shift = modelData.Id_shift.String
	stuctReturn.Fr_disable_cash = modelData.Fr_disable_cash.Int32
	stuctReturn.Fr_disable_cashless = modelData.Fr_disable_cashless.Int32
	stuctReturn.Ffd_version = modelData.Ffd_version.Int32
	return &stuctReturn
}
