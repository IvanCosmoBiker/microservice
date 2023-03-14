package model

import (
	"database/sql"
	"encoding/json"
	modelInterface "ephorservices/pkg/model/schema/interface/model"
	"log"

	pgx "github.com/jackc/pgx/v5"
)

const (
	Type_Modem           = 0x01
	Type_Automat         = 0x02
	Type_FiscalRegistrar = 0x03
	Type_BillValidator   = 0x04
	Type_CoinChanger     = 0x05
	Type_Cashless1       = 0x06
	Type_Cashless2       = 0x07
	Type_ExternCashless  = 0x08
	Type_Screen          = 0x09
	Type_NFC             = 0x0A
	Type_USB1            = 0x0B
)

// Eversyst
const (
	Type_CoffeeGroupLeft     = 0x0C
	Type_CoffeeGroupRigth    = 0x0D
	Type_SteamGeneratorLeft  = 0x0E
	Type_SteamGeneratorRigth = 0x0F
	Type_Boiler              = 0x010
)

// IM30
const Type_Im30 = 0x011

const (
	State_NotFound = 0
	State_Init     = 1
	State_Work     = 2
	State_Error    = 3
	State_Disabled = 4
	State_Enabled  = 5
)

var nameTable string = "automat_device"
var nameSchema string = "account"

type AutomatDeviceModel struct {
	Id         int
	Automat_id sql.NullInt32
	Type       sql.NullInt32
	State      sql.NullInt32
}

func Init() modelInterface.Model {
	return &AutomatDeviceModel{}
}

func (aem *AutomatDeviceModel) InitData(jsonData []byte) {
	json.Unmarshal(jsonData, &aem)
}

func (aem *AutomatDeviceModel) JsonSerialize() []byte {
	data, _ := json.Marshal(aem)
	return data
}

func (aem *AutomatDeviceModel) GetNameTable() string {
	return nameTable
}

func (aem *AutomatDeviceModel) GetNameSchema(accountId string) string {
	return nameSchema + "" + accountId
}

type ReturningStruct struct {
	Id         int
	Automat_id int32
	Type       int32
	State      int32
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
	Model := AutomatDeviceModel{}
	err := row.Scan(&Model.Id,
		&Model.Automat_id,
		&Model.Type,
		&Model.State)
	if err != nil {
		log.Printf("%v\n", err)
		return &ReturningStruct{}, err
	}
	Result := MakeDataInSturct(&Model)

	return Result, nil
}

func MakeDataInSturct(modelData *AutomatDeviceModel) *ReturningStruct {
	stuctReturn := ReturningStruct{}
	stuctReturn.Id = modelData.Id
	stuctReturn.Automat_id = modelData.Automat_id.Int32
	stuctReturn.Type = modelData.Type.Int32
	stuctReturn.State = modelData.State.Int32
	return &stuctReturn
}
