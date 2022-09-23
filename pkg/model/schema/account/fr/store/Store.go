package fr

import (
	"context"
	connectionPostgresql "ephorservices/pkg/db"
	model "ephorservices/pkg/model/schema/account/fr/model"
	modelInterface "ephorservices/pkg/model/schema/interface/model"
	parserTypes "ephorservices/pkg/parser/typeParse"
	"errors"
	"fmt"
	"log"
)

type StoreFr struct {
	Connection      *connectionPostgresql.Manager
	model           modelInterface.Model
	fieldsReturning string
}

func NewStore(conn *connectionPostgresql.Manager) *StoreFr {
	modelStore := model.Init()
	fields := conn.Conn.PrepareReturningFields(modelStore)
	return &StoreFr{
		Connection:      conn,
		model:           modelStore,
		fieldsReturning: fields,
	}
}

func (sf *StoreFr) Get(accountId interface{}) ([]*model.ReturningStruct, error) {
	ctx := context.Background()
	var options map[string]interface{}
	fieldSql, _, _ := sf.Connection.Conn.PrepareGet(options, sf.model)
	sql := fmt.Sprintf("SELECT %v FROM %v.%v", fieldSql, sf.model.GetNameSchema(parserTypes.ParseTypeInString(accountId)), sf.model.GetNameTable())
	return sf.GetDataOfMapStuct(ctx, sql)
}

func (sf *StoreFr) AddByParams(parametrs map[string]interface{}) (*model.ReturningStruct, error) {
	if _, exist := parametrs["account_id"]; !exist {
		return &model.ReturningStruct{}, errors.New("not found paramentr account_id")
	}
	accountId := parserTypes.ParseTypeInString(parametrs["account_id"])
	delete(parametrs, "account_id")
	ctx := context.Background()
	sqlField, sqlValues, values := sf.Connection.Conn.PrepareInsert(parametrs, sf.model)
	sql := fmt.Sprintf("INSERT INTO %v.%v (%v) VALUES (%s) RETURNING %s", sf.model.GetNameSchema(accountId), sf.model.GetNameTable(), sqlField, sqlValues, sf.fieldsReturning)
	result, err := sf.GetDataOfStruct(ctx, sql, values...)
	if err != nil {
		return result, err
	}
	return result, nil
}

func (sf *StoreFr) SetByParams(parametrs map[string]interface{}) (*model.ReturningStruct, error) {
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
	sqlValues, sqlWhere, Values := sf.Connection.Conn.PrepareUpdate(parametrs, options, sf.model)
	sql := fmt.Sprintf("UPDATE %v.%v SET %s %v", sf.model.GetNameSchema(accountId), sf.model.GetNameTable(), sqlValues, sqlWhere)
	result, err := sf.GetDataOfStruct(ctx, sql, Values...)
	if err != nil {
		return result, err
	}
	return result, nil
}

func (sf *StoreFr) GetOneById(id, accountId interface{}) (*model.ReturningStruct, error) {
	var result *model.ReturningStruct
	var err error
	options := make(map[string]interface{})
	options["id"] = id
	sqlField, sqlWhere, valuesWhere := sf.Connection.Conn.PrepareGet(options, sf.model)
	sql := fmt.Sprintf("SELECT %v FROM %v.%v %v", sqlField, sf.model.GetNameSchema(parserTypes.ParseTypeInString(accountId)), sf.model.GetNameTable(), sqlWhere)
	ctx := context.Background()
	result, err = sf.GetDataOfStruct(ctx, sql, valuesWhere...)
	if err != nil {
		return result, err
	}
	return result, nil
}

func (sf *StoreFr) GetDataOfStruct(ctx context.Context, sql string, args ...interface{}) (*model.ReturningStruct, error) {
	var result *model.ReturningStruct
	row, err := sf.Connection.Conn.QueryRow(ctx, sql, args...)
	if err != nil {
		return result, err
	}
	return model.ScanModelRow(row)
}

func (sf *StoreFr) GetDataOfMapStuct(ctx context.Context, sql string, args ...interface{}) ([]*model.ReturningStruct, error) {
	var result []*model.ReturningStruct
	rows, err := sf.Connection.Conn.Query(ctx, sql, args...)
	if err != nil {
		return result, err
	}
	return model.ScanModelRows(rows)
}

func (sf *StoreFr) GetWithOptions(options map[string]interface{}) ([]*model.ReturningStruct, error) {
	var result []*model.ReturningStruct
	if _, exist := options["account_id"]; !exist {
		return result, errors.New("not found paramentr account_id")
	}
	accountId := parserTypes.ParseTypeInString(options["account_id"])
	delete(options, "account_id")
	ctx := context.Background()
	sqlField, sqlWhere, valuesWhere := sf.Connection.Conn.PrepareGet(options, sf.model)
	sql := fmt.Sprintf("SELECT %v FROM %v.%v %v", sqlField, sf.model.GetNameSchema(accountId), sf.model.GetNameTable(), sqlWhere)
	log.Println(sql)
	return sf.GetDataOfMapStuct(ctx, sql, valuesWhere...)
}
