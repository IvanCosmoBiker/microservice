package store

import (
	"context"
	connectionPostgresql "ephorservices/pkg/db"
	model "ephorservices/pkg/model/schema/account/automatevent/model"
	modelInterface "ephorservices/pkg/model/schema/interface/model"
	parserTypes "ephorservices/pkg/parser/typeParse"
	"errors"
	"fmt"
	"log"
)

type StoreAutomatEvent struct {
	Connection      *connectionPostgresql.Manager
	model           modelInterface.Model
	fieldsReturning string
}

func NewStore(conn *connectionPostgresql.Manager) *StoreAutomatEvent {
	modelStore := model.Init()
	fields := conn.Conn.PrepareReturningFields(modelStore)
	return &StoreAutomatEvent{
		Connection:      conn,
		model:           modelStore,
		fieldsReturning: fields,
	}
}

func (sae *StoreAutomatEvent) Get(accountId interface{}) ([]*model.ReturningStruct, error) {
	ctx := context.Background()
	var options map[string]interface{}
	fieldSql, _, _ := sae.Connection.Conn.PrepareGet(options, sae.model)
	sql := fmt.Sprintf("SELECT %v FROM %v.%v", fieldSql, sae.model.GetNameSchema(parserTypes.ParseTypeInString(accountId)), sae.model.GetNameTable())
	return sae.GetDataOfMapStuct(ctx, sql)
}

func (sae *StoreAutomatEvent) AddByParams(parametrs map[string]interface{}) (*model.ReturningStruct, error) {
	if _, exist := parametrs["account_id"]; exist == false {
		return &model.ReturningStruct{}, errors.New("not found paramentr account_id")
	}
	accountId := parserTypes.ParseTypeInString(parametrs["account_id"])
	delete(parametrs, "account_id")
	ctx := context.Background()
	sqlField, sqlValues, values := sae.Connection.Conn.PrepareInsert(parametrs, sae.model)
	sql := fmt.Sprintf("INSERT INTO %v.%v (%v) VALUES (%s) RETURNING %s", sae.model.GetNameSchema(accountId), sae.model.GetNameTable(), sqlField, sqlValues, sae.fieldsReturning)
	result, err := sae.GetDataOfStruct(ctx, sql, values...)
	if err != nil {
		return result, err
	}
	return result, nil
}

func (sae *StoreAutomatEvent) SetByParams(parametrs map[string]interface{}) (*model.ReturningStruct, error) {
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
	sqlValues, sqlWhere, Values := sae.Connection.Conn.PrepareUpdate(parametrs, options, sae.model)
	sql := fmt.Sprintf("UPDATE %v.%v SET %s %v", sae.model.GetNameSchema(accountId), sae.model.GetNameTable(), sqlValues, sqlWhere)
	result, err := sae.GetDataOfStruct(ctx, sql, Values...)
	if err != nil {
		return result, err
	}
	return result, nil
}

func (sae *StoreAutomatEvent) GetOneById(id, accountId interface{}) (*model.ReturningStruct, error) {
	var result *model.ReturningStruct
	var err error
	options := make(map[string]interface{})
	options["id"] = id
	sqlField, sqlWhere, valuesWhere := sae.Connection.Conn.PrepareGet(options, sae.model)
	sql := fmt.Sprintf("SELECT %v FROM %v.%v %v", sqlField, sae.model.GetNameSchema(parserTypes.ParseTypeInString(accountId)), sae.model.GetNameTable(), sqlWhere)
	ctx := context.Background()
	result, err = sae.GetDataOfStruct(ctx, sql, valuesWhere...)
	if err != nil {
		return result, err
	}
	return result, nil
}

func (sae *StoreAutomatEvent) GetDataOfStruct(ctx context.Context, sql string, args ...interface{}) (*model.ReturningStruct, error) {
	var result *model.ReturningStruct
	row, err := sae.Connection.Conn.QueryRow(ctx, sql, args...)
	if err != nil {
		return result, err
	}
	return model.ScanModelRow(row)
}

func (sae *StoreAutomatEvent) GetDataOfMapStuct(ctx context.Context, sql string, args ...interface{}) ([]*model.ReturningStruct, error) {
	var result []*model.ReturningStruct
	rows, err := sae.Connection.Conn.Query(ctx, sql, args...)
	if err != nil {
		return result, err
	}
	return model.ScanModelRows(rows)
}

func (sae *StoreAutomatEvent) GetWithOptions(options map[string]interface{}) ([]*model.ReturningStruct, error) {
	var result []*model.ReturningStruct
	if _, exist := options["account_id"]; exist == false {
		return result, errors.New("not found paramentr account_id")
	}
	accountId := parserTypes.ParseTypeInString(options["account_id"])
	delete(options, "account_id")
	ctx := context.Background()
	sqlField, sqlWhere, valuesWhere := sae.Connection.Conn.PrepareGet(options, sae.model)
	sql := fmt.Sprintf("SELECT %v FROM %v.%v %v", sqlField, sae.model.GetNameSchema(accountId), sae.model.GetNameTable(), sqlWhere)
	log.Println(sql)
	return sae.GetDataOfMapStuct(ctx, sql, valuesWhere...)
}
