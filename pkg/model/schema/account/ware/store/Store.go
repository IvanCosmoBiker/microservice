package ware

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

type StoreWare struct {
	Connection      *connectionPostgresql.Manager
	model           modelInterface.Model
	fieldsReturning string
}

func NewStore(conn *connectionPostgresql.Manager) *StoreWare {
	modelStore := model.Init()
	fields := conn.Conn.PrepareReturningFields(modelStore)
	return &StoreWare{
		Connection:      conn,
		model:           modelStore,
		fieldsReturning: fields,
	}
}

func (sw *StoreWare) Get(accountId interface{}) ([]*model.ReturningStruct, error) {
	ctx := context.Background()
	var options map[string]interface{}
	fieldSql, _, _ := sw.Connection.Conn.PrepareGet(options, sw.model)
	sql := fmt.Sprintf("SELECT %v FROM %v.%v", fieldSql, sw.model.GetNameSchema(parserTypes.ParseTypeInString(accountId)), sw.model.GetNameTable())
	return sw.GetDataOfMapStuct(ctx, sql)
}

func (sw *StoreWare) AddByParams(parametrs map[string]interface{}) (*model.ReturningStruct, error) {
	if _, exist := parametrs["account_id"]; exist == false {
		return &model.ReturningStruct{}, errors.New("not found paramentr account_id")
	}
	accountId := parserTypes.ParseTypeInString(parametrs["account_id"])
	delete(parametrs, "account_id")
	ctx := context.Background()
	sqlField, sqlValues, values := sw.Connection.Conn.PrepareInsert(parametrs, sw.model)
	sql := fmt.Sprintf("INSERT INTO %v.%v (%v) VALUES (%s) RETURNING %s", sw.model.GetNameSchema(accountId), sw.model.GetNameTable(), sqlField, sqlValues, sw.fieldsReturning)
	result, err := sw.GetDataOfStruct(ctx, sql, values...)
	if err != nil {
		return result, err
	}
	return result, nil
}

func (sw *StoreWare) SetByParams(parametrs map[string]interface{}) (*model.ReturningStruct, error) {
	if _, exist := parametrs["id"]; !exist {
		return &model.ReturningStruct{}, errors.New("not found id in parametrs")
	}
	if _, exist := parametrs["account_id"]; !exist {
		return &model.ReturningStruct{}, errors.New("not found paramentr account_id")
	}
	accountId := parserTypes.ParseTypeInString(parametrs["account_id"])
	delete(parametrs, "account_id")
	id := parametrs["id"]
	delete(parametrs, "id")
	options := make(map[string]interface{})
	options["id"] = id
	ctx := context.Background()
	sqlValues, sqlWhere, Values := sw.Connection.Conn.PrepareUpdate(parametrs, options, sw.model)
	sql := fmt.Sprintf("UPDATE %v.%v SET %s %v", sw.model.GetNameSchema(accountId), sw.model.GetNameTable(), sqlValues, sqlWhere)
	result, err := sw.GetDataOfStruct(ctx, sql, Values...)
	if err != nil {
		return result, err
	}
	return result, nil
}

func (sw *StoreWare) GetOneById(id, accountId interface{}) (*model.ReturningStruct, error) {
	var result *model.ReturningStruct
	var err error
	options := make(map[string]interface{})
	options["id"] = id
	sqlField, sqlWhere, valuesWhere := sw.Connection.Conn.PrepareGet(options, sw.model)
	sql := fmt.Sprintf("SELECT %v FROM %v.%v %v", sqlField, sw.model.GetNameSchema(parserTypes.ParseTypeInString(accountId)), sw.model.GetNameTable(), sqlWhere)
	ctx := context.Background()
	result, err = sw.GetDataOfStruct(ctx, sql, valuesWhere...)
	if err != nil {
		return result, err
	}
	return result, nil
}

func (sw *StoreWare) GetDataOfStruct(ctx context.Context, sql string, args ...interface{}) (*model.ReturningStruct, error) {
	var result *model.ReturningStruct
	row, err := sw.Connection.Conn.QueryRow(ctx, sql, args...)
	if err != nil {
		return result, err
	}
	return model.ScanModelRow(row)
}

func (sw *StoreWare) GetDataOfMapStuct(ctx context.Context, sql string, args ...interface{}) ([]*model.ReturningStruct, error) {
	var result []*model.ReturningStruct
	rows, err := sw.Connection.Conn.Query(ctx, sql, args...)
	if err != nil {
		return result, err
	}
	return model.ScanModelRows(rows)
}

func (sw *StoreWare) GetWithOptions(options map[string]interface{}) ([]*model.ReturningStruct, error) {
	var result []*model.ReturningStruct
	if _, exist := options["account_id"]; !exist {
		return result, errors.New("not found paramentr account_id")
	}
	accountId := parserTypes.ParseTypeInString(options["account_id"])
	delete(options, "account_id")
	ctx := context.Background()
	sqlField, sqlWhere, valuesWhere := sw.Connection.Conn.PrepareGet(options, sw.model)
	sql := fmt.Sprintf("SELECT %v FROM %v.%v %v", sqlField, sw.model.GetNameSchema(accountId), sw.model.GetNameTable(), sqlWhere)
	log.Println(sql)
	return sw.GetDataOfMapStuct(ctx, sql, valuesWhere...)
}
