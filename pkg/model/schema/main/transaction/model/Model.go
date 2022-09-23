package model

import (
	"database/sql"
	"encoding/json"
	modelInterface "ephorservices/pkg/model/schema/interface/model"
	"log"

	pgx "github.com/jackc/pgx/v4"
)

var nameTable string = "transaction"
var nameSchema string = "main"

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

type ReturningStruct struct {
	Id            int
	Noise         string
	Token_id      string
	Token_type    int64
	Account_id    int64
	Customer_id   int64
	Automat_id    int64
	Date          string
	Sum           int64
	Ps_type       int64
	Ps_order      string
	Ps_tid        string
	Ps_code       string
	Ps_desc       string
	Ps_invoice_id string
	Pay_type      int64
	Fn            int64
	Fd            int64
	Fp            string
	F_type        int64
	F_receipt     string
	F_desc        string
	F_status      int64
	Qr_format     int64
	F_qr          string
	Status        int64
	Error         string
}

func Init() modelInterface.Model {
	return &TransactionModel{}
}

func (tm *TransactionModel) InitData(jsonData []byte) {
	json.Unmarshal(jsonData, &tm)
}

func (tm *TransactionModel) JsonSerialize() []byte {
	data, _ := json.Marshal(tm)
	return data
}

func (tm *TransactionModel) GetNameTable() string {
	return nameTable
}

func (tm *TransactionModel) GetNameSchema(account string) string {
	return nameSchema
}

func (tm *TransactionModel) CheckType(parametrs map[string]interface{}) {

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
	Transaction := TransactionModel{}
	err := row.Scan(&Transaction.Id,
		&Transaction.Noise,
		&Transaction.Token_id,
		&Transaction.Token_type,
		&Transaction.Account_id,
		&Transaction.Customer_id,
		&Transaction.Automat_id,
		&Transaction.Date,
		&Transaction.Sum,
		&Transaction.Ps_type,
		&Transaction.Ps_order,
		&Transaction.Ps_tid,
		&Transaction.Ps_code,
		&Transaction.Ps_desc,
		&Transaction.Ps_invoice_id,
		&Transaction.Pay_type,
		&Transaction.Fn,
		&Transaction.Fd,
		&Transaction.Fp,
		&Transaction.F_type,
		&Transaction.F_receipt,
		&Transaction.F_desc,
		&Transaction.F_status,
		&Transaction.Qr_format,
		&Transaction.F_qr,
		&Transaction.Status,
		&Transaction.Error)
	if err != nil {
		log.Printf("%v\n", err)
		return &ReturningStruct{}, err
	}
	Result := MakeDataInSturct(&Transaction)

	return Result, nil
}

func MakeDataInSturct(modelData *TransactionModel) *ReturningStruct {
	stuctReturn := ReturningStruct{}
	stuctReturn.Id = modelData.Id
	stuctReturn.Noise = modelData.Noise.String
	stuctReturn.Token_id = modelData.Token_id.String
	stuctReturn.Token_type = modelData.Token_type.Int64
	stuctReturn.Account_id = modelData.Account_id.Int64
	stuctReturn.Automat_id = modelData.Automat_id.Int64
	stuctReturn.Date = modelData.Date.String
	stuctReturn.Sum = modelData.Sum.Int64
	stuctReturn.Ps_type = modelData.Ps_type.Int64
	stuctReturn.Ps_order = modelData.Ps_order.String
	stuctReturn.Ps_code = modelData.Ps_code.String
	stuctReturn.Ps_desc = modelData.Ps_desc.String
	stuctReturn.Ps_invoice_id = modelData.Ps_invoice_id.String
	stuctReturn.Pay_type = modelData.Pay_type.Int64
	stuctReturn.Fn = modelData.Fn.Int64
	stuctReturn.Fd = modelData.Fd.Int64
	stuctReturn.Fp = modelData.Fp.String
	stuctReturn.F_type = modelData.F_type.Int64
	stuctReturn.F_receipt = modelData.F_receipt.String
	stuctReturn.F_desc = modelData.F_desc.String
	stuctReturn.F_status = modelData.F_status.Int64
	stuctReturn.Qr_format = modelData.Qr_format.Int64
	stuctReturn.F_qr = modelData.F_qr.String
	stuctReturn.Status = modelData.Status.Int64
	stuctReturn.Error = modelData.Error.String
	return &stuctReturn
}
