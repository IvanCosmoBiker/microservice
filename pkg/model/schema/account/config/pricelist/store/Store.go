package store

import (
	"context"
	connectionPostgresql "ephorservices/pkg/db"
	model "ephorservices/pkg/model/schema/account/config/pricelist/model"
	modelInterface "ephorservices/pkg/model/schema/interface/model"
	parserTypes "ephorservices/pkg/parser/typeParse"
	"errors"
	"fmt"
	"log"
)

type StorePriceList struct {
	Connection      *connectionPostgresql.Manager
	model           modelInterface.Model
	fieldsReturning string
}

func NewStore(conn *connectionPostgresql.Manager) *StorePriceList {
	modelStore := model.Init()
	fields := conn.Conn.PrepareReturningFields(modelStore)
	return &StorePriceList{
		Connection:      conn,
		model:           modelStore,
		fieldsReturning: fields,
	}
}

func (spl *StorePriceList) Get(accountId interface{}) ([]*model.ReturningStruct, error) {
	ctx := context.Background()
	var options map[string]interface{}
	fieldSql, _, _ := spl.Connection.Conn.PrepareGet(options, spl.model)
	sql := fmt.Sprintf("SELECT %v FROM %v.%v", fieldSql, spl.model.GetNameSchema(parserTypes.ParseTypeInString(accountId)), spl.model.GetNameTable())
	return spl.GetDataOfMapStuct(ctx, sql)
}

func (spl *StorePriceList) AddByParams(parametrs map[string]interface{}) (*model.ReturningStruct, error) {
	if _, exist := parametrs["account_id"]; !exist {
		return &model.ReturningStruct{}, errors.New("not found paramentr account_id")
	}
	accountId := parserTypes.ParseTypeInString(parametrs["account_id"])
	delete(parametrs, "account_id")
	ctx := context.Background()
	sqlField, sqlValues, values := spl.Connection.Conn.PrepareInsert(parametrs, spl.model)
	sql := fmt.Sprintf("INSERT INTO %v.%v (%v) VALUES (%s) RETURNING %s", spl.model.GetNameSchema(accountId), spl.model.GetNameTable(), sqlField, sqlValues, spl.fieldsReturning)
	result, err := spl.GetDataOfStruct(ctx, sql, values...)
	if err != nil {
		return result, err
	}
	return result, nil
}

func (spl *StorePriceList) SetByParams(parametrs map[string]interface{}) (*model.ReturningStruct, error) {
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
	sqlValues, sqlWhere, Values := spl.Connection.Conn.PrepareUpdate(parametrs, options, spl.model)
	sql := fmt.Sprintf("UPDATE %v.%v SET %s %v", spl.model.GetNameSchema(accountId), spl.model.GetNameTable(), sqlValues, sqlWhere)
	result, err := spl.GetDataOfStruct(ctx, sql, Values...)
	if err != nil {
		return result, err
	}
	return result, nil
}

func (spl *StorePriceList) GetOneById(id, accountId interface{}) (*model.ReturningStruct, error) {
	var result *model.ReturningStruct
	var err error
	options := make(map[string]interface{})
	options["id"] = id
	sqlField, sqlWhere, valuesWhere := spl.Connection.Conn.PrepareGet(options, spl.model)
	sql := fmt.Sprintf("SELECT %v FROM %v.%v %v", sqlField, spl.model.GetNameSchema(parserTypes.ParseTypeInString(accountId)), spl.model.GetNameTable(), sqlWhere)
	ctx := context.Background()
	result, err = spl.GetDataOfStruct(ctx, sql, valuesWhere...)
	if err != nil {
		return result, err
	}
	return result, nil
}

func (spl *StorePriceList) GetDataOfStruct(ctx context.Context, sql string, args ...interface{}) (*model.ReturningStruct, error) {
	var result *model.ReturningStruct
	row, err := spl.Connection.Conn.QueryRow(ctx, sql, args...)
	if err != nil {
		return result, err
	}
	return model.ScanModelRow(row)
}

func (spl *StorePriceList) GetDataOfMapStuct(ctx context.Context, sql string, args ...interface{}) ([]*model.ReturningStruct, error) {
	var result []*model.ReturningStruct
	rows, err := spl.Connection.Conn.Query(ctx, sql, args...)
	if err != nil {
		return result, err
	}
	return model.ScanModelRows(rows)
}

func (spl *StorePriceList) GetWithOptions(options map[string]interface{}) ([]*model.ReturningStruct, error) {
	var result []*model.ReturningStruct
	if _, exist := options["account_id"]; !exist {
		return result, errors.New("not found paramentr account_id")
	}
	accountId := parserTypes.ParseTypeInString(options["account_id"])
	delete(options, "account_id")
	ctx := context.Background()
	sqlField, sqlWhere, valuesWhere := spl.Connection.Conn.PrepareGet(options, spl.model)
	sql := fmt.Sprintf("SELECT %v FROM %v.%v %v", sqlField, spl.model.GetNameSchema(accountId), spl.model.GetNameTable(), sqlWhere)
	log.Println(sql)
	return spl.GetDataOfMapStuct(ctx, sql, valuesWhere...)
}
