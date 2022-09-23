package model

import (
	"database/sql"
	"encoding/json"
	modelInterface "ephorservices/pkg/model/schema/interface/model"
	"log"

	pgx "github.com/jackc/pgx/v4"
)

const (
	SendUnSuccess = 0
	SendSuccess   = 1
)

var nameTable string = "modem_command"
var nameSchema string = "main"

type CommandModel struct {
	Id             int
	Modem_id       sql.NullInt32
	Command        sql.NullInt32
	Command_param1 sql.NullInt32
	Date           sql.NullString
	Sended         sql.NullInt32
}

func Init() modelInterface.Model {
	return &CommandModel{}
}

func (cm *CommandModel) InitData(jsonData []byte) {
	json.Unmarshal(jsonData, &cm)
}

func (cm *CommandModel) JsonSerialize() []byte {
	data, _ := json.Marshal(cm)
	return data
}

func (cm *CommandModel) GetNameTable(account string) string {
	return nameTable
}

func (cm *CommandModel) GetNameSchema(account int) string {
	return nameSchema
}

type ReturningStruct struct {
	Id             int
	Modem_id       int32
	Command        int32
	Command_param1 int32
	Date           string
	Sended         int32
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
	Model := CommandModel{}
	err := row.Scan(&Model.Id,
		&Model.Modem_id,
		&Model.Command,
		&Model.Command_param1,
		&Model.Date,
		&Model.Sended)
	if err != nil {
		log.Printf("%v\n", err)
		return &ReturningStruct{}, err
	}
	Result := MakeDataInSturct(&Model)
	return Result, nil
}

func MakeDataInSturct(modelData *CommandModel) *ReturningStruct {
	stuctReturn := ReturningStruct{}
	stuctReturn.Id = modelData.Id
	stuctReturn.Modem_id = modelData.Modem_id.Int32
	stuctReturn.Command = modelData.Command.Int32
	stuctReturn.Command_param1 = modelData.Command_param1.Int32
	stuctReturn.Date = modelData.Date.String
	stuctReturn.Sended = modelData.Sended.Int32
	return &stuctReturn
}
