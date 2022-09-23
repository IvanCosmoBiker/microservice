package store

import (
	"context"
	connectionPostgresql "ephorservices/pkg/db"
	modelInterface "ephorservices/pkg/model/schema/interface/model"
	model "ephorservices/pkg/model/schema/main/transactionproduct/model"
	"errors"
	"fmt"
	"log"
)

type StoreTransactionProduct struct {
	Connection      *connectionPostgresql.Manager
	model           modelInterface.Model
	fieldsReturning string
}

func NewStore(conn *connectionPostgresql.Manager) *StoreTransactionProduct {
	modelStore := model.Init()
	fields := conn.Conn.PrepareReturningFields(modelStore)
	return &StoreTransactionProduct{
		Connection:      conn,
		model:           modelStore,
		fieldsReturning: fields,
	}
}

func (tp *StoreTransactionProduct) AddByParams(parametrs map[string]interface{}) (*model.ReturningStruct, error) {
	ctx := context.Background()
	sqlField, sqlValues, values := tp.Connection.Conn.PrepareInsert(parametrs, tp.model)
	sql := fmt.Sprintf("INSERT INTO %v.%v (%v) VALUES (%s) RETURNING %s", tp.model.GetNameSchema(""), tp.model.GetNameTable(), sqlField, sqlValues, tp.fieldsReturning)
	result, err := tp.GetDataOfStruct(ctx, sql, values...)
	if err != nil {
		return result, err
	}
	return result, nil
}

func (tp *StoreTransactionProduct) SetByParams(parametrs map[string]interface{}) (*model.ReturningStruct, error) {
	if _, exist := parametrs["id"]; exist != true {
		return &model.ReturningStruct{}, errors.New("not found id in parametrs")
	}
	id := parametrs["id"]
	delete(parametrs, "id")
	options := make(map[string]interface{})
	options["id"] = id
	ctx := context.Background()
	sqlValues, sqlWhere, Values := tp.Connection.Conn.PrepareUpdate(parametrs, options, tp.model)
	sql := fmt.Sprintf("UPDATE %v.%v SET %s %v RETURNING %s", tp.model.GetNameSchema(""), tp.model.GetNameTable(), sqlValues, sqlWhere, tp.fieldsReturning)
	result, err := tp.GetDataOfStruct(ctx, sql, Values...)
	if err != nil {
		return result, err
	}
	return result, nil
}

func (tp *StoreTransactionProduct) GetOneById(id interface{}) (*model.ReturningStruct, error) {
	var result *model.ReturningStruct
	var err error
	options := make(map[string]interface{})
	options["id"] = id
	sqlField, sqlWhere, valuesWhere := tp.Connection.Conn.PrepareGet(options, tp.model)
	sql := fmt.Sprintf("SELECT %v FROM %v.%v %v", sqlField, tp.model.GetNameSchema(""), tp.model.GetNameTable(), sqlWhere)
	ctx := context.Background()
	result, err = tp.GetDataOfStruct(ctx, sql, valuesWhere...)
	if err != nil {
		return result, err
	}
	return result, nil
}

func (tp *StoreTransactionProduct) GetDataOfStruct(ctx context.Context, sql string, args ...interface{}) (*model.ReturningStruct, error) {
	var result *model.ReturningStruct
	row, err := tp.Connection.Conn.QueryRow(ctx, sql, args...)
	if err != nil {
		return result, err
	}
	return model.ScanModelRow(row)
}

func (tp *StoreTransactionProduct) GetDataOfMapStuct(ctx context.Context, sql string, args ...interface{}) ([]*model.ReturningStruct, error) {
	var result []*model.ReturningStruct
	rows, err := tp.Connection.Conn.Query(ctx, sql, args...)
	if err != nil {
		return result, err
	}
	return model.ScanModelRows(rows)
}

func (tp *StoreTransactionProduct) GetWithOptions(options map[string]interface{}) ([]*model.ReturningStruct, error) {
	ctx := context.Background()
	sqlField, sqlWhere, valuesWhere := tp.Connection.Conn.PrepareGet(options, tp.model)
	sql := fmt.Sprintf("SELECT %v FROM %v.%v %v", sqlField, tp.model.GetNameSchema(""), tp.model.GetNameTable(), sqlWhere)
	log.Println(sql)
	return tp.GetDataOfMapStuct(ctx, sql, valuesWhere...)
}
