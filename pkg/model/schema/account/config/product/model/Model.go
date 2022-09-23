package model

import (
	"database/sql"
	"encoding/json"
	modelInterface "ephorservices/pkg/model/schema/interface/model"
	"log"

	pgx "github.com/jackc/pgx/v4"
)

var nameTable string = "config_product"
var nameSchema string = "account"

type ConfigProductModel struct {
	Id        int
	Config_id sql.NullInt32
	Select_id sql.NullString
	Cid       sql.NullInt32
	Ware_id   sql.NullInt32
	Name      sql.NullString
	Tax_rate  sql.NullInt32
	Type      sql.NullInt32
	Qt_recom  sql.NullInt32
	Qt_max    sql.NullInt32
	Recipe_id sql.NullInt32
}

func Init() modelInterface.Model {
	return &ConfigProductModel{}
}

func (aem *ConfigProductModel) InitData(jsonData []byte) {
	json.Unmarshal(jsonData, &aem)
}

func (aem *ConfigProductModel) JsonSerialize() []byte {
	data, _ := json.Marshal(aem)
	return data
}

func (aem *ConfigProductModel) GetNameTable() string {
	return nameTable
}

func (aem *ConfigProductModel) GetNameSchema(accountId string) string {
	return nameSchema + "" + accountId
}

type ReturningStruct struct {
	Id        int
	Config_id int32
	Select_id string
	Cid       int32
	Ware_id   int32
	Name      string
	Tax_rate  int32
	Type      int32
	Qt_recom  int32
	Qt_max    int32
	Recipe_id int32
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
	Model := ConfigProductModel{}
	err := row.Scan(&Model.Id,
		&Model.Name,
		&Model.Config_id,
		&Model.Select_id,
		&Model.Cid,
		&Model.Ware_id,
		&Model.Name,
		&Model.Tax_rate,
		&Model.Type,
		&Model.Qt_recom,
		&Model.Qt_max,
		&Model.Recipe_id)

	if err != nil {
		log.Printf("%v\n", err)
		return &ReturningStruct{}, err
	}
	Result := MakeDataInSturct(&Model)

	return Result, nil
}

func MakeDataInSturct(modelData *ConfigProductModel) *ReturningStruct {
	stuctReturn := ReturningStruct{}
	stuctReturn.Id = modelData.Id
	stuctReturn.Name = modelData.Name.String
	stuctReturn.Config_id = modelData.Config_id.Int32
	stuctReturn.Select_id = modelData.Select_id.String
	stuctReturn.Cid = modelData.Cid.Int32
	stuctReturn.Ware_id = modelData.Ware_id.Int32
	stuctReturn.Name = modelData.Name.String
	stuctReturn.Tax_rate = modelData.Tax_rate.Int32
	stuctReturn.Type = modelData.Type.Int32
	stuctReturn.Qt_recom = modelData.Qt_recom.Int32
	stuctReturn.Qt_max = modelData.Qt_max.Int32
	stuctReturn.Recipe_id = modelData.Recipe_id.Int32
	return &stuctReturn
}
