package store

import (
	"context"
	connectionPostgresql "ephorservices/pkg/db"
	model "ephorservices/pkg/model/schema/account/automat/device/model"
	modelInterface "ephorservices/pkg/model/schema/interface/model"
	parserTypes "ephorservices/pkg/parser/typeParse"
	"errors"
	"fmt"
	"log"
)

type StoreAutomatDevice struct {
	Connection      *connectionPostgresql.Manager
	model           modelInterface.Model
	fieldsReturning string
}

func NewStore(conn *connectionPostgresql.Manager) *StoreAutomatDevice {
	modelStore := model.Init()
	fields, _ := conn.Conn.PrepareReturningFields(modelStore)
	return &StoreAutomatDevice{
		Connection:      conn,
		model:           modelStore,
		fieldsReturning: fields,
	}
}

func (sad *StoreAutomatDevice) Get(accountId interface{}) ([]*model.ReturningStruct, error) {
	ctx := context.Background()
	var options map[string]interface{}
	fieldSql, _, _ := sad.Connection.Conn.PrepareGet(options, sad.model)
	sql := fmt.Sprintf("SELECT %v FROM %v.%v", fieldSql, sad.model.GetNameSchema(parserTypes.ParseTypeInString(accountId)), sad.model.GetNameTable())
	return sad.GetDataOfMapStuct(ctx, sql)
}

func (sad *StoreAutomatDevice) AddByParams(parametrs map[string]interface{}) (*model.ReturningStruct, error) {
	if _, exist := parametrs["account_id"]; exist == false {
		return &model.ReturningStruct{}, errors.New("not found paramentr account_id")
	}
	accountId := parserTypes.ParseTypeInString(parametrs["account_id"])
	delete(parametrs, "account_id")
	ctx := context.Background()
	sqlField, sqlValues, values := sad.Connection.Conn.PrepareInsert(parametrs, sad.model)
	sql := fmt.Sprintf("INSERT INTO %v.%v (%v) VALUES (%s) RETURNING %s", sad.model.GetNameSchema(accountId), sad.model.GetNameTable(), sqlField, sqlValues, sad.fieldsReturning)
	result, err := sad.GetDataOfStruct(ctx, sql, values...)
	if err != nil {
		return result, err
	}
	return result, nil
}

func (sad *StoreAutomatDevice) SetByParams(parametrs map[string]interface{}) (*model.ReturningStruct, error) {
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
	sqlValues, sqlWhere, Values := sad.Connection.Conn.PrepareUpdate(parametrs, options, sad.model)
	sql := fmt.Sprintf("UPDATE %v.%v SET %s %v RETURNING %s", sad.model.GetNameSchema(accountId), sad.model.GetNameTable(), sqlValues, sqlWhere, sad.fieldsReturning)
	result, err := sad.GetDataOfStruct(ctx, sql, Values...)
	if err != nil {
		return result, err
	}
	return result, nil
}

func (sad *StoreAutomatDevice) GetOneById(id, accountId interface{}) (*model.ReturningStruct, error) {
	var result *model.ReturningStruct
	var err error
	options := make(map[string]interface{})
	options["id"] = id
	sqlField, sqlWhere, valuesWhere := sad.Connection.Conn.PrepareGet(options, sad.model)
	sql := fmt.Sprintf("SELECT %v FROM %v.%v %v", sqlField, sad.model.GetNameSchema(parserTypes.ParseTypeInString(accountId)), sad.model.GetNameTable(), sqlWhere)
	ctx := context.Background()
	result, err = sad.GetDataOfStruct(ctx, sql, valuesWhere...)
	if err != nil {
		return result, err
	}
	return result, nil
}

func (sad *StoreAutomatDevice) GetDataOfStruct(ctx context.Context, sql string, args ...interface{}) (*model.ReturningStruct, error) {
	var result *model.ReturningStruct
	row, err := sad.Connection.Conn.QueryRow(ctx, sql, args...)
	if err != nil {
		return result, err
	}
	return model.ScanModelRow(row)
}

func (sad *StoreAutomatDevice) GetDataOfMapStuct(ctx context.Context, sql string, args ...interface{}) ([]*model.ReturningStruct, error) {
	var result []*model.ReturningStruct
	rows, err := sad.Connection.Conn.Query(ctx, sql, args...)
	if err != nil {
		return result, err
	}
	return model.ScanModelRows(rows)
}

func (sad *StoreAutomatDevice) GetWithOptions(options map[string]interface{}) ([]*model.ReturningStruct, error) {
	var result []*model.ReturningStruct
	if _, exist := options["account_id"]; exist == false {
		return result, errors.New("not found paramentr account_id")
	}
	accountId := parserTypes.ParseTypeInString(options["account_id"])
	delete(options, "account_id")
	ctx := context.Background()
	sqlField, sqlWhere, valuesWhere := sad.Connection.Conn.PrepareGet(options, sad.model)
	sql := fmt.Sprintf("SELECT %v FROM %v.%v %v", sqlField, sad.model.GetNameSchema(accountId), sad.model.GetNameTable(), sqlWhere)
	log.Println(sql)
	return sad.GetDataOfMapStuct(ctx, sql, valuesWhere...)
}

func (sad *StoreAutomatDevice) GetOneWithOptions(options map[string]interface{}) (*model.ReturningStruct, error) {
	var result *model.ReturningStruct
	if _, exist := options["account_id"]; exist == false {
		return result, errors.New("not found paramentr account_id")
	}
	accountId := parserTypes.ParseTypeInString(options["account_id"])
	delete(options, "account_id")
	ctx := context.Background()
	sqlField, sqlWhere, valuesWhere := sad.Connection.Conn.PrepareGet(options, sad.model)
	sql := fmt.Sprintf("SELECT %v FROM %v.%v %v", sqlField, sad.model.GetNameSchema(accountId), sad.model.GetNameTable(), sqlWhere)
	log.Println(sql)
	return sad.GetDataOfStruct(ctx, sql, valuesWhere...)
}
