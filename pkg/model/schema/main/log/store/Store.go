package model

import (
	"context"
	connectionPostgresql "ephorservices/pkg/db"
	modelInterface "ephorservices/pkg/model/schema/interface/model"
	model "ephorservices/pkg/model/schema/main/transaction/model"
	"errors"
	"fmt"
	"log"
)

type StoreLog struct {
	Connection      *connectionPostgresql.Manager
	model           modelInterface.Model
	fieldsReturning string
}

func NewStore(conn *connectionPostgresql.Manager) *StoreLog {
	modelStore := model.Init()
	fields := conn.Conn.PrepareReturningFields(modelStore)
	return &StoreLog{
		Connection:      conn,
		model:           modelStore,
		fieldsReturning: fields,
	}
}

func (sl *StoreLog) AddByParams(parametrs map[string]interface{}) (*model.ReturningStruct, error) {
	ctx := context.Background()
	sqlField, sqlValues, values := sl.Connection.Conn.PrepareInsert(parametrs, sl.model)
	sql := fmt.Sprintf("INSERT INTO %v.%v (%v) VALUES (%s) RETURNING %s", sl.model.GetNameSchema(0), sl.model.GetNameTable(), sqlField, sqlValues, sl.fieldsReturning)
	result, err := sl.GetDataOfStruct(ctx, sql, values...)
	if err != nil {
		return result, err
	}
	return result, nil
}

func (sl *StoreLog) SetByParams(parametrs map[string]interface{}) (*model.ReturningStruct, error) {
	if _, exist := parametrs["id"]; exist != true {
		return &model.ReturningStruct{}, errors.New("not found id in parametrs")
	}
	id := parametrs["id"]
	delete(parametrs, "id")
	options := make(map[string]interface{})
	options["id"] = id
	ctx := context.Background()
	sqlValues, sqlWhere, Values := sl.Connection.Conn.PrepareUpdate(parametrs, options, sl.model)
	sql := fmt.Sprintf("UPDATE %v.%v SET %s %v", sl.model.GetNameSchema(0), sl.model.GetNameTable(), sqlValues, sqlWhere)
	result, err := sl.GetDataOfStruct(ctx, sql, Values...)
	if err != nil {
		return result, err
	}
	return result, nil
}

func (sl *StoreLog) GetOneById(id interface{}) (*model.ReturningStruct, error) {
	var result *model.ReturningStruct
	var err error
	options := make(map[string]interface{})
	options["id"] = id
	sqlField, sqlWhere, valuesWhere := sl.Connection.Conn.PrepareGet(options, sl.model)
	sql := fmt.Sprintf("SELECT %v FROM %v.%v %v", sqlField, sl.model.GetNameSchema(0), sl.model.GetNameTable(), sqlWhere)
	ctx := context.Background()
	result, err = sl.GetDataOfStruct(ctx, sql, valuesWhere...)
	if err != nil {
		return result, err
	}
	return result, nil
}

func (sl *StoreLog) GetDataOfStruct(ctx context.Context, sql string, args ...interface{}) (*model.ReturningStruct, error) {
	var result *model.ReturningStruct
	row, err := sl.Connection.Conn.QueryRow(ctx, sql, args...)
	if err != nil {
		return result, err
	}
	return model.ScanModelRow(row)
}

func (sl *StoreLog) GetDataOfMapStuct(ctx context.Context, sql string, args ...interface{}) ([]*model.ReturningStruct, error) {
	var result []*model.ReturningStruct
	rows, err := sl.Connection.Conn.Query(ctx, sql, args...)
	if err != nil {
		return result, err
	}
	return model.ScanModelRows(rows)
}

func (sl *StoreLog) GetWithOptions(options map[string]interface{}) ([]*model.ReturningStruct, error) {
	ctx := context.Background()
	sqlField, sqlWhere, valuesWhere := sl.Connection.Conn.PrepareGet(options, sl.model)
	sql := fmt.Sprintf("SELECT %v FROM %v.%v %v", sqlField, sl.model.GetNameSchema(0), sl.model.GetNameTable(), sqlWhere)
	log.Println(sql)
	return sl.GetDataOfMapStuct(ctx, sql, valuesWhere...)
}
