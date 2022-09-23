package store

import (
	"context"
	connectionPostgresql "ephorservices/pkg/db"
	model "ephorservices/pkg/model/schema/account/automat/model"
	modelInterface "ephorservices/pkg/model/schema/interface/model"
	parserTypes "ephorservices/pkg/parser/typeParse"
	"errors"
	"fmt"
	"log"
)

type StoreAutomat struct {
	Connection      *connectionPostgresql.Manager
	model           modelInterface.Model
	fieldsReturning string
}

func NewStore(conn *connectionPostgresql.Manager) *StoreAutomat {
	modelStore := model.Init()
	fields := conn.Conn.PrepareReturningFields(modelStore)
	return &StoreAutomat{
		Connection:      conn,
		model:           modelStore,
		fieldsReturning: fields,
	}
}

func (sa *StoreAutomat) Get(accountId interface{}) ([]*model.ReturningStruct, error) {
	ctx := context.Background()
	var options map[string]interface{}
	fieldSql, _, _ := sa.Connection.Conn.PrepareGet(options, sa.model)
	sql := fmt.Sprintf("SELECT %v FROM %v.%v", fieldSql, sa.model.GetNameSchema(parserTypes.ParseTypeInString(accountId)), sa.model.GetNameTable())
	return sa.GetDataOfMapStuct(ctx, sql)
}

func (sa *StoreAutomat) AddByParams(parametrs map[string]interface{}) (*model.ReturningStruct, error) {
	if _, exist := parametrs["account_id"]; exist == false {
		return &model.ReturningStruct{}, errors.New("not found paramentr account_id")
	}
	accountId := parserTypes.ParseTypeInString(parametrs["account_id"])
	delete(parametrs, "account_id")
	ctx := context.Background()
	sqlField, sqlValues, values := sa.Connection.Conn.PrepareInsert(parametrs, sa.model)
	sql := fmt.Sprintf("INSERT INTO %v.%v (%v) VALUES (%s) RETURNING %s", sa.model.GetNameSchema(accountId), sa.model.GetNameTable(), sqlField, sqlValues, sa.fieldsReturning)
	result, err := sa.GetDataOfStruct(ctx, sql, values...)
	if err != nil {
		return result, err
	}
	return result, nil
}

func (sa *StoreAutomat) SetByParams(parametrs map[string]interface{}) (*model.ReturningStruct, error) {
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
	sqlValues, sqlWhere, Values := sa.Connection.Conn.PrepareUpdate(parametrs, options, sa.model)
	sql := fmt.Sprintf("UPDATE %v.%v SET %s %v", sa.model.GetNameSchema(accountId), sa.model.GetNameTable(), sqlValues, sqlWhere)
	result, err := sa.GetDataOfStruct(ctx, sql, Values...)
	if err != nil {
		return result, err
	}
	return result, nil
}

func (sa *StoreAutomat) GetOneById(id, accountId interface{}) (*model.ReturningStruct, error) {
	var result *model.ReturningStruct
	var err error
	options := make(map[string]interface{})
	options["id"] = id
	sqlField, sqlWhere, valuesWhere := sa.Connection.Conn.PrepareGet(options, sa.model)
	sql := fmt.Sprintf("SELECT %v FROM %v.%v %v", sqlField, sa.model.GetNameSchema(parserTypes.ParseTypeInString(accountId)), sa.model.GetNameTable(), sqlWhere)
	ctx := context.Background()
	result, err = sa.GetDataOfStruct(ctx, sql, valuesWhere...)
	if err != nil {
		return result, err
	}
	return result, nil
}

func (sa *StoreAutomat) GetDataOfStruct(ctx context.Context, sql string, args ...interface{}) (*model.ReturningStruct, error) {
	var result *model.ReturningStruct
	row, err := sa.Connection.Conn.QueryRow(ctx, sql, args...)
	if err != nil {
		return result, err
	}
	return model.ScanModelRow(row)
}

func (sa *StoreAutomat) GetDataOfMapStuct(ctx context.Context, sql string, args ...interface{}) ([]*model.ReturningStruct, error) {
	var result []*model.ReturningStruct
	rows, err := sa.Connection.Conn.Query(ctx, sql, args...)
	if err != nil {
		return result, err
	}
	return model.ScanModelRows(rows)
}

func (sa *StoreAutomat) GetWithOptions(options map[string]interface{}) ([]*model.ReturningStruct, error) {
	var result []*model.ReturningStruct
	if _, exist := options["account_id"]; exist == false {
		return result, errors.New("not found paramentr account_id")
	}
	accountId := parserTypes.ParseTypeInString(options["account_id"])
	delete(options, "account_id")
	ctx := context.Background()
	sqlField, sqlWhere, valuesWhere := sa.Connection.Conn.PrepareGet(options, sa.model)
	sql := fmt.Sprintf("SELECT %v FROM %v.%v %v", sqlField, sa.model.GetNameSchema(accountId), sa.model.GetNameTable(), sqlWhere)
	log.Println(sql)
	return sa.GetDataOfMapStuct(ctx, sql, valuesWhere...)
}
