package ware

import (
	"database/sql"
	"encoding/json"
	modelInterface "ephorservices/pkg/model/schema/interface/model"
	"log"

	pgx "github.com/jackc/pgx/v4"
)

var nameTable string = "ware"
var nameSchema string = "account"

type WareModel struct {
	Id             int
	Code           sql.NullString
	Name           sql.NullString
	State          sql.NullInt32
	Type           sql.NullInt32
	Description    sql.NullString
	Img            sql.NullString
	Img_hash       sql.NullString
	Purchase_price sql.NullInt64
	Img_path       sql.NullString
}

func Init() modelInterface.Model {
	return &WareModel{}
}

func (wm *WareModel) InitData(jsonData []byte) {
	json.Unmarshal(jsonData, &wm)
}

func (wm *WareModel) JsonSerialize() []byte {
	data, _ := json.Marshal(wm)
	return data
}

func (wm *WareModel) GetNameTable() string {
	return nameTable
}

func (wm *WareModel) GetNameSchema(accountId string) string {
	return nameSchema + "" + accountId
}

type ReturningStruct struct {
	Id             int
	Code           string
	Name           string
	State          int32
	Type           int32
	Description    string
	Img            string
	Img_hash       string
	Purchase_price int64
	Img_path       string
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
	Model := WareModel{}
	err := row.Scan(&Model.Id,
		&Model.Code,
		&Model.Name,
		&Model.State,
		&Model.Type,
		&Model.Description,
		&Model.Img,
		&Model.Img_hash,
		&Model.Purchase_price,
		&Model.Img_path)
	if err != nil {
		log.Printf("%v\n", err)
		return &ReturningStruct{}, err
	}
	Result := MakeDataInSturct(&Model)
	return Result, nil
}

func MakeDataInSturct(modelData *WareModel) *ReturningStruct {
	stuctReturn := ReturningStruct{}
	stuctReturn.Id = modelData.Id
	stuctReturn.Code = modelData.Code.String
	stuctReturn.Name = modelData.Name.String
	stuctReturn.State = modelData.State.Int32
	stuctReturn.Type = modelData.Type.Int32
	stuctReturn.Description = modelData.Description.String
	stuctReturn.Img = modelData.Img.String
	stuctReturn.Img_hash = modelData.Img_hash.String
	stuctReturn.Purchase_price = modelData.Purchase_price.Int64
	stuctReturn.Img_path = modelData.Img_path.String
	return &stuctReturn
}
