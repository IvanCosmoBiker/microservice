package store

import (
	"context"
	connectionPostgresql "ephorservices/pkg/db"
	model "ephorservices/pkg/model/schema/account/companypoint/model"
	modelInterface "ephorservices/pkg/model/schema/interface/model"
	parserTypes "ephorservices/pkg/parser/typeParse"
	"errors"
	"fmt"
	"log"
)

type StoreCompany struct {
	Connection      *connectionPostgresql.Manager
	model           modelInterface.Model
	fieldsReturning string
}

func NewStore(conn *connectionPostgresql.Manager) *StoreCompany {
	modelStore := model.Init()
	fields := conn.Conn.PrepareReturningFields(modelStore)
	return &StoreCompany{
		Connection:      conn,
		model:           modelStore,
		fieldsReturning: fields,
	}
}

func (sc *StoreCompany) Get(accountId interface{}) ([]*model.ReturningStruct, error) {
	ctx := context.Background()
	var options map[string]interface{}
	fieldSql, _, _ := sc.Connection.Conn.PrepareGet(options, sc.model)
	sql := fmt.Sprintf("SELECT %v FROM %v.%v", fieldSql, sc.model.GetNameSchema(parserTypes.ParseTypeInString(accountId), sc.model.GetNameTable()))
	return sc.GetDataOfMapStuct(ctx, sql)
}

func (sc *StoreCompany) AddByParams(parametrs map[string]interface{}) (*model.ReturningStruct, error) {
	if _, exist := parametrs["account_id"]; exist == false {
		return &model.ReturningStruct{}, errors.New("not found paramentr account_id")
	}
	accountId := parserTypes.ParseTypeInString(parametrs["account_id"])
	delete(parametrs, "account_id")
	ctx := context.Background()
	sqlField, sqlValues, values := sc.Connection.Conn.PrepareInsert(parametrs, sc.model)
	sql := fmt.Sprintf("INSERT INTO %v.%v (%v) VALUES (%s) RETURNING %s", sc.model.GetNameSchema(accountId), sc.model.GetNameTable(), sqlField, sqlValues, sc.fieldsReturning)
	result, err := sc.GetDataOfStruct(ctx, sql, values...)
	if err != nil {
		return result, err
	}
	return result, nil
}

func (sc *StoreCompany) SetByParams(parametrs map[string]interface{}) (*model.ReturningStruct, error) {
	if _, exist := parametrs["id"]; exist != true {
		return &model.ReturningStruct{}, errors.New("not found id in parametrs")
	}
	if _, exist := parametrs["account_id"]; exist == false {
		return &model.ReturningStruct{}, errors.New("not found paramentr account_id")
	}
	accountId := parserTypes.ParseTypeInString(parametrs["account_id"])
	delete(parametrs, "account_id")
	id := parametrs["id"]
	delete(parametrs, "id")
	options := make(map[string]interface{})
	options["id"] = id
	ctx := context.Background()
	sqlValues, sqlWhere, Values := sc.Connection.Conn.PrepareUpdate(parametrs, options, sc.model)
	sql := fmt.Sprintf("UPDATE %v.%v SET %s %v", sc.model.GetNameSchema(accountId), sc.model.GetNameTable(), sqlValues, sqlWhere)
	result, err := sc.GetDataOfStruct(ctx, sql, Values...)
	if err != nil {
		return result, err
	}
	return result, nil
}

func (sc *StoreCompany) GetOneById(id, accountId interface{}) (*model.ReturningStruct, error) {
	var result *model.ReturningStruct
	var err error
	options := make(map[string]interface{})
	options["id"] = id
	sqlField, sqlWhere, valuesWhere := sc.Connection.Conn.PrepareGet(options, sc.model)
	sql := fmt.Sprintf("SELECT %v FROM %v.%v %v", sqlField, sc.model.GetNameSchema(parserTypes.ParseTypeInString(accountId)), sc.model.GetNameTable(), sqlWhere)
	ctx := context.Background()
	result, err = sc.GetDataOfStruct(ctx, sql, valuesWhere...)
	if err != nil {
		return result, err
	}
	return result, nil
}

func (sc *StoreCompany) GetDataOfStruct(ctx context.Context, sql string, args ...interface{}) (*model.ReturningStruct, error) {
	var result *model.ReturningStruct
	row, err := sc.Connection.Conn.QueryRow(ctx, sql, args...)
	if err != nil {
		return result, err
	}
	return model.ScanModelRow(row)
}

func (sc *StoreCompany) GetDataOfMapStuct(ctx context.Context, sql string, args ...interface{}) ([]*model.ReturningStruct, error) {
	var result []*model.ReturningStruct
	rows, err := sc.Connection.Conn.Query(ctx, sql, args...)
	if err != nil {
		return result, err
	}
	return model.ScanModelRows(rows)
}

func (sc *StoreCompany) GetWithOptions(options map[string]interface{}) ([]*model.ReturningStruct, error) {
	var result []*model.ReturningStruct
	if _, exist := options["account_id"]; exist == false {
		return result, errors.New("not found paramentr account_id")
	}
	accountId := parserTypes.ParseTypeInString(options["account_id"])
	delete(options, "account_id")
	ctx := context.Background()
	sqlField, sqlWhere, valuesWhere := sc.Connection.Conn.PrepareGet(options, sc.model)
	sql := fmt.Sprintf("SELECT %v FROM %v.%v %v", sqlField, sc.model.GetNameSchema(accountId), sc.model.GetNameTable(), sqlWhere)
	log.Println(sql)
	return sc.GetDataOfMapStuct(ctx, sql, valuesWhere...)
}
