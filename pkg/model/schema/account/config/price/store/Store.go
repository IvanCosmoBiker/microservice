package store

import (
	"context"
	connectionPostgresql "ephorservices/pkg/db"
	model "ephorservices/pkg/model/schema/account/config/product/model"
	modelInterface "ephorservices/pkg/model/schema/interface/model"
	parserTypes "ephorservices/pkg/parser/typeParse"
	"errors"
	"fmt"
	"log"
)

type StoreConfigProductPrice struct {
	Connection      *connectionPostgresql.Manager
	model           modelInterface.Model
	fieldsReturning string
}

func NewStore(conn *connectionPostgresql.Manager) *StoreConfigProductPrice {
	modelStore := model.Init()
	fields := conn.Conn.PrepareReturningFields(modelStore)
	return &StoreConfigProductPrice{
		Connection:      conn,
		model:           modelStore,
		fieldsReturning: fields,
	}
}

func (scpp *StoreConfigProductPrice) Get(accountId interface{}) ([]*model.ReturningStruct, error) {
	ctx := context.Background()
	var options map[string]interface{}
	fieldSql, _, _ := scpp.Connection.Conn.PrepareGet(options, scpp.model)
	sql := fmt.Sprintf("SELECT %v FROM %v.%v", fieldSql, scpp.model.GetNameSchema(parserTypes.ParseTypeInString(accountId)), scpp.model.GetNameTable())
	return scpp.GetDataOfMapStuct(ctx, sql)
}

func (scpp *StoreConfigProductPrice) AddByParams(parametrs map[string]interface{}) (*model.ReturningStruct, error) {
	if _, exist := parametrs["account_id"]; !exist {
		return &model.ReturningStruct{}, errors.New("not found paramentr account_id")
	}
	accountId := parserTypes.ParseTypeInString(parametrs["account_id"])
	delete(parametrs, "account_id")
	ctx := context.Background()
	sqlField, sqlValues, values := scpp.Connection.Conn.PrepareInsert(parametrs, scpp.model)
	sql := fmt.Sprintf("INSERT INTO %v.%v (%v) VALUES (%s) RETURNING %s", scpp.model.GetNameSchema(accountId), scpp.model.GetNameTable(), sqlField, sqlValues, scpp.fieldsReturning)
	result, err := scpp.GetDataOfStruct(ctx, sql, values...)
	if err != nil {
		return result, err
	}
	return result, nil
}

func (scpp *StoreConfigProductPrice) SetByParams(parametrs map[string]interface{}) (*model.ReturningStruct, error) {
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
	sqlValues, sqlWhere, Values := scpp.Connection.Conn.PrepareUpdate(parametrs, options, scpp.model)
	sql := fmt.Sprintf("UPDATE %v.%v SET %s %v", scpp.model.GetNameSchema(accountId), scpp.model.GetNameTable(), sqlValues, sqlWhere)
	result, err := scpp.GetDataOfStruct(ctx, sql, Values...)
	if err != nil {
		return result, err
	}
	return result, nil
}

func (scpp *StoreConfigProductPrice) GetOneById(id, accountId interface{}) (*model.ReturningStruct, error) {
	var result *model.ReturningStruct
	var err error
	options := make(map[string]interface{})
	options["id"] = id
	sqlField, sqlWhere, valuesWhere := scpp.Connection.Conn.PrepareGet(options, scpp.model)
	sql := fmt.Sprintf("SELECT %v FROM %v.%v %v", sqlField, scpp.model.GetNameSchema(parserTypes.ParseTypeInString(accountId)), scpp.model.GetNameTable(), sqlWhere)
	ctx := context.Background()
	result, err = scpp.GetDataOfStruct(ctx, sql, valuesWhere...)
	if err != nil {
		return result, err
	}
	return result, nil
}

func (scpp *StoreConfigProductPrice) GetDataOfStruct(ctx context.Context, sql string, args ...interface{}) (*model.ReturningStruct, error) {
	var result *model.ReturningStruct
	row, err := scpp.Connection.Conn.QueryRow(ctx, sql, args...)
	if err != nil {
		return result, err
	}
	return model.ScanModelRow(row)
}

func (scpp *StoreConfigProductPrice) GetDataOfMapStuct(ctx context.Context, sql string, args ...interface{}) ([]*model.ReturningStruct, error) {
	var result []*model.ReturningStruct
	rows, err := scpp.Connection.Conn.Query(ctx, sql, args...)
	if err != nil {
		return result, err
	}
	return model.ScanModelRows(rows)
}

func (scpp *StoreConfigProductPrice) GetWithOptions(options map[string]interface{}) ([]*model.ReturningStruct, error) {
	var result []*model.ReturningStruct
	if _, exist := options["account_id"]; !exist {
		return result, errors.New("not found paramentr account_id")
	}
	accountId := parserTypes.ParseTypeInString(options["account_id"])
	delete(options, "account_id")
	ctx := context.Background()
	sqlField, sqlWhere, valuesWhere := scpp.Connection.Conn.PrepareGet(options, scpp.model)
	sql := fmt.Sprintf("SELECT %v FROM %v.%v %v", sqlField, scpp.model.GetNameSchema(accountId), scpp.model.GetNameTable(), sqlWhere)
	log.Println(sql)
	return scpp.GetDataOfMapStuct(ctx, sql, valuesWhere...)
}
