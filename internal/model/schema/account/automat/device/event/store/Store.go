package store

import (
	"context"
	connectionPostgresql "ephorservices/pkg/db"
	model "ephorservices/pkg/model/schema/account/automat/device/event/model"
	modelInterface "ephorservices/pkg/model/schema/interface/model"
	parserTypes "ephorservices/pkg/parser/typeParse"
	"errors"
	"fmt"
	"log"
)

type StoreAutomatDeviceEvent struct {
	Connection      *connectionPostgresql.Manager
	model           modelInterface.Model
	fieldsReturning string
}

func NewStore(conn *connectionPostgresql.Manager) *StoreAutomatDeviceEvent {
	modelStore := model.Init()
	fields, _ := conn.Conn.PrepareReturningFields(modelStore)
	return &StoreAutomatDeviceEvent{
		Connection:      conn,
		model:           modelStore,
		fieldsReturning: fields,
	}
}

func (sade *StoreAutomatDeviceEvent) Get(accountId interface{}) ([]*model.ReturningStruct, error) {
	ctx := context.Background()
	var options map[string]interface{}
	fieldSql, _, _ := sade.Connection.Conn.PrepareGet(options, sade.model)
	sql := fmt.Sprintf("SELECT %v FROM %v.%v", fieldSql, sade.model.GetNameSchema(parserTypes.ParseTypeInString(accountId)), sade.model.GetNameTable())
	return sade.GetDataOfMapStuct(ctx, sql)
}

func (sade *StoreAutomatDeviceEvent) AddByParams(parametrs map[string]interface{}) (*model.ReturningStruct, error) {
	if _, exist := parametrs["account_id"]; exist == false {
		return &model.ReturningStruct{}, errors.New("not found paramentr account_id")
	}
	accountId := parserTypes.ParseTypeInString(parametrs["account_id"])
	delete(parametrs, "account_id")
	ctx := context.Background()
	sqlField, sqlValues, values := sade.Connection.Conn.PrepareInsert(parametrs, sade.model)
	sql := fmt.Sprintf("INSERT INTO %v.%v (%v) VALUES (%s) RETURNING %s", sade.model.GetNameSchema(accountId), sade.model.GetNameTable(), sqlField, sqlValues, sade.fieldsReturning)
	result, err := sade.GetDataOfStruct(ctx, sql, values...)
	if err != nil {
		return result, err
	}
	return result, nil
}

func (sade *StoreAutomatDeviceEvent) SetByParams(parametrs map[string]interface{}) (*model.ReturningStruct, error) {
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
	sqlValues, sqlWhere, Values := sade.Connection.Conn.PrepareUpdate(parametrs, options, sade.model)
	sql := fmt.Sprintf("UPDATE %v.%v SET %s %v RETURNING %s", sade.model.GetNameSchema(accountId), sade.model.GetNameTable(), sqlValues, sqlWhere, sade.fieldsReturning)
	result, err := sade.GetDataOfStruct(ctx, sql, Values...)
	if err != nil {
		return result, err
	}
	return result, nil
}

func (sade *StoreAutomatDeviceEvent) GetOneById(id, accountId interface{}) (*model.ReturningStruct, error) {
	var result *model.ReturningStruct
	var err error
	options := make(map[string]interface{})
	options["id"] = id
	sqlField, sqlWhere, valuesWhere := sade.Connection.Conn.PrepareGet(options, sade.model)
	sql := fmt.Sprintf("SELECT %v FROM %v.%v %v", sqlField, sade.model.GetNameSchema(parserTypes.ParseTypeInString(accountId)), sade.model.GetNameTable(), sqlWhere)
	ctx := context.Background()
	result, err = sade.GetDataOfStruct(ctx, sql, valuesWhere...)
	if err != nil {
		return result, err
	}
	return result, nil
}

func (sade *StoreAutomatDeviceEvent) GetDataOfStruct(ctx context.Context, sql string, args ...interface{}) (*model.ReturningStruct, error) {
	var result *model.ReturningStruct
	row, err := sade.Connection.Conn.QueryRow(ctx, sql, args...)
	if err != nil {
		return result, err
	}
	return model.ScanModelRow(row)
}

func (sade *StoreAutomatDeviceEvent) GetDataOfMapStuct(ctx context.Context, sql string, args ...interface{}) ([]*model.ReturningStruct, error) {
	var result []*model.ReturningStruct
	rows, err := sade.Connection.Conn.Query(ctx, sql, args...)
	if err != nil {
		return result, err
	}
	return model.ScanModelRows(rows)
}

func (sade *StoreAutomatDeviceEvent) GetWithOptions(options map[string]interface{}) ([]*model.ReturningStruct, error) {
	var result []*model.ReturningStruct
	if _, exist := options["account_id"]; exist == false {
		return result, errors.New("not found paramentr account_id")
	}
	accountId := parserTypes.ParseTypeInString(options["account_id"])
	delete(options, "account_id")
	ctx := context.Background()
	sqlField, sqlWhere, valuesWhere := sade.Connection.Conn.PrepareGet(options, sade.model)
	sql := fmt.Sprintf("SELECT %v FROM %v.%v %v", sqlField, sade.model.GetNameSchema(accountId), sade.model.GetNameTable(), sqlWhere)
	log.Println(sql)
	return sade.GetDataOfMapStuct(ctx, sql, valuesWhere...)
}

func (sade *StoreAutomatDeviceEvent) GetOneWithOptions(options map[string]interface{}) (*model.ReturningStruct, error) {
	var result *model.ReturningStruct
	if _, exist := options["account_id"]; exist == false {
		return result, errors.New("not found paramentr account_id")
	}
	accountId := parserTypes.ParseTypeInString(options["account_id"])
	delete(options, "account_id")
	ctx := context.Background()
	sqlField, sqlWhere, valuesWhere := sade.Connection.Conn.PrepareGet(options, sade.model)
	sql := fmt.Sprintf("SELECT %v FROM %v.%v %v", sqlField, sade.model.GetNameSchema(accountId), sade.model.GetNameTable(), sqlWhere)
	log.Println(sql)
	return sade.GetDataOfStruct(ctx, sql, valuesWhere...)
}
