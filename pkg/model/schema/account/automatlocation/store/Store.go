package store

import (
	"context"
	connectionPostgresql "ephorservices/pkg/db"
	model "ephorservices/pkg/model/schema/account/automatlocation/model"
	modelInterface "ephorservices/pkg/model/schema/interface/model"
	parserTypes "ephorservices/pkg/parser/typeParse"
	"errors"
	"fmt"
	"log"
)

type StoreAutomatLocation struct {
	Connection      *connectionPostgresql.Manager
	model           modelInterface.Model
	fieldsReturning string
}

func NewStore(conn *connectionPostgresql.Manager) *StoreAutomatLocation {
	modelStore := model.Init()
	fields := conn.Conn.PrepareReturningFields(modelStore)
	return &StoreAutomatLocation{
		Connection:      conn,
		model:           modelStore,
		fieldsReturning: fields,
	}
}

func (sal *StoreAutomatLocation) Get(accountId interface{}) ([]*model.ReturningStruct, error) {
	ctx := context.Background()
	var options map[string]interface{}
	fieldSql, _, _ := sal.Connection.Conn.PrepareGet(options, sal.model)
	sql := fmt.Sprintf("SELECT %v FROM %v.%v", fieldSql, sal.model.GetNameSchema(parserTypes.ParseTypeInString(accountId)), sal.model.GetNameTable())
	return sal.GetDataOfMapStuct(ctx, sql)
}

func (sal *StoreAutomatLocation) AddByParams(parametrs map[string]interface{}) (*model.ReturningStruct, error) {
	if _, exist := parametrs["account_id"]; exist == false {
		return &model.ReturningStruct{}, errors.New("not found paramentr account_id")
	}
	accountId := parserTypes.ParseTypeInString(parametrs["account_id"])
	delete(parametrs, "account_id")
	ctx := context.Background()
	sqlField, sqlValues, values := sal.Connection.Conn.PrepareInsert(parametrs, sal.model)
	sql := fmt.Sprintf("INSERT INTO %v.%v (%v) VALUES (%s) RETURNING %s", sal.model.GetNameSchema(accountId), sal.model.GetNameTable(), sqlField, sqlValues, sal.fieldsReturning)
	result, err := sal.GetDataOfStruct(ctx, sql, values...)
	if err != nil {
		return result, err
	}
	return result, nil
}

func (sal *StoreAutomatLocation) SetByParams(parametrs map[string]interface{}) (*model.ReturningStruct, error) {
	if _, exist := parametrs["id"]; exist != true {
		return &model.ReturningStruct{}, errors.New("not found id in parametrs")
	}
	if _, exist := parametrs["account_id"]; exist == false {
		return &model.ReturningStruct{}, errors.New("not found paramentr account_id")
	}
	delete(parametrs, "account_id")
	accountId := parserTypes.ParseTypeInString(parametrs["account_id"])
	id := parametrs["id"]
	delete(parametrs, "id")
	options := make(map[string]interface{})
	options["id"] = id
	ctx := context.Background()
	sqlValues, sqlWhere, Values := sal.Connection.Conn.PrepareUpdate(parametrs, options, sal.model)
	sql := fmt.Sprintf("UPDATE %v.%v SET %s %v", sal.model.GetNameSchema(accountId), sal.model.GetNameTable(), sqlValues, sqlWhere)
	result, err := sal.GetDataOfStruct(ctx, sql, Values...)
	if err != nil {
		return result, err
	}
	return result, nil
}

func (sal *StoreAutomatLocation) GetOneById(id, accountId interface{}) (*model.ReturningStruct, error) {
	var result *model.ReturningStruct
	var err error
	options := make(map[string]interface{})
	options["id"] = id
	sqlField, sqlWhere, valuesWhere := sal.Connection.Conn.PrepareGet(options, sal.model)
	sql := fmt.Sprintf("SELECT %v FROM %v.%v %v", sqlField, sal.model.GetNameSchema(parserTypes.ParseTypeInString(accountId)), sal.model.GetNameTable(), sqlWhere)
	ctx := context.Background()
	result, err = sal.GetDataOfStruct(ctx, sql, valuesWhere...)
	if err != nil {
		return result, err
	}
	return result, nil
}

func (sal *StoreAutomatLocation) GetDataOfStruct(ctx context.Context, sql string, args ...interface{}) (*model.ReturningStruct, error) {
	var result *model.ReturningStruct
	row, err := sal.Connection.Conn.QueryRow(ctx, sql, args...)
	if err != nil {
		return result, err
	}
	return model.ScanModelRow(row)
}

func (sal *StoreAutomatLocation) GetDataOfMapStuct(ctx context.Context, sql string, args ...interface{}) ([]*model.ReturningStruct, error) {
	var result []*model.ReturningStruct
	rows, err := sal.Connection.Conn.Query(ctx, sql, args...)
	if err != nil {
		return result, err
	}
	return model.ScanModelRows(rows)
}

func (sal *StoreAutomatLocation) GetWithOptions(options map[string]interface{}) ([]*model.ReturningStruct, error) {
	var result []*model.ReturningStruct
	if _, exist := options["account_id"]; exist == false {
		return result, errors.New("not found paramentr account_id")
	}
	accountId := parserTypes.ParseTypeInString(options["account_id"])
	delete(options, "account_id")
	ctx := context.Background()
	sqlField, sqlWhere, valuesWhere := sal.Connection.Conn.PrepareGet(options, sal.model)
	sql := fmt.Sprintf("SELECT %v FROM %v.%v %v", sqlField, sal.model.GetNameSchema(accountId), sal.model.GetNameTable(), sqlWhere)
	log.Println(sql)
	return sal.GetDataOfMapStuct(ctx, sql, valuesWhere...)
}
