package store

import (
	"context"
	connectionPostgresql "ephorservices/pkg/db"
	modelInterface "ephorservices/pkg/model/schema/interface/model"
	model "ephorservices/pkg/model/schema/main/transaction/model"
	"errors"
	"fmt"
	"log"
)

type StoreTransaction struct {
	Connection      *connectionPostgresql.Manager
	model           modelInterface.Model
	fieldsReturning string
}

func NewStore(conn *connectionPostgresql.Manager) *StoreTransaction {
	modelStore := model.Init()
	fields := conn.Conn.PrepareReturningFields(modelStore)
	return &StoreTransaction{
		Connection:      conn,
		model:           modelStore,
		fieldsReturning: fields,
	}
}

func (t *StoreTransaction) AddByParams(parametrs map[string]interface{}) (*model.ReturningStruct, error) {
	ctx := context.Background()
	sqlField, sqlValues, values := t.Connection.Conn.PrepareInsert(parametrs, t.model)
	fmt.Printf("\n%s\n", sqlField)
	fmt.Printf("\n%s\n", sqlValues)
	fmt.Printf("\n%+v\n", values)
	sql := fmt.Sprintf("INSERT INTO %v.%v (%v) VALUES (%s) RETURNING %s", t.model.GetNameSchema(""), t.model.GetNameTable(), sqlField, sqlValues, t.fieldsReturning)
	result, err := t.GetDataOfStruct(ctx, sql, values...)
	log.Println(sql)
	if err != nil {
		return result, err
	}
	return result, nil
}

func (t *StoreTransaction) SetByParams(parametrs map[string]interface{}) (*model.ReturningStruct, error) {
	if _, exist := parametrs["id"]; exist != true {
		return &model.ReturningStruct{}, errors.New("not found id in parametrs")
	}
	id := parametrs["id"]
	delete(parametrs, "id")
	options := make(map[string]interface{})
	options["id"] = id
	ctx := context.Background()
	sqlValues, sqlWhere, Values := t.Connection.Conn.PrepareUpdate(parametrs, options, t.model)
	sql := fmt.Sprintf("UPDATE %v.%v SET %s %v RETURNING %s", t.model.GetNameSchema(""), t.model.GetNameTable(), sqlValues, sqlWhere, t.fieldsReturning)
	fmt.Printf("\n%v\n", sql)
	result, err := t.GetDataOfStruct(ctx, sql, Values...)
	if err != nil {
		fmt.Printf("\n%v\n", err)
		return result, err
	}
	return result, nil
}

func (t *StoreTransaction) GetOneById(id interface{}) (*model.ReturningStruct, error) {
	var result *model.ReturningStruct
	var err error
	options := make(map[string]interface{})
	options["id"] = id
	sqlField, sqlWhere, valuesWhere := t.Connection.Conn.PrepareGet(options, t.model)
	sql := fmt.Sprintf("SELECT %v FROM %v.%v %v", sqlField, t.model.GetNameSchema(""), t.model.GetNameTable(), sqlWhere)
	ctx := context.Background()
	result, err = t.GetDataOfStruct(ctx, sql, valuesWhere...)
	if err != nil {
		return result, err
	}
	return result, nil
}

func (t *StoreTransaction) GetDataOfStruct(ctx context.Context, sql string, args ...interface{}) (*model.ReturningStruct, error) {
	var result *model.ReturningStruct
	fmt.Printf("\n%s\n", sql)
	row, err := t.Connection.Conn.QueryRow(ctx, sql, args...)
	if err != nil {
		return result, err
	}
	return model.ScanModelRow(row)
}

func (t *StoreTransaction) GetDataOfMapStuct(ctx context.Context, sql string, args ...interface{}) ([]*model.ReturningStruct, error) {
	var result []*model.ReturningStruct
	rows, err := t.Connection.Conn.Query(ctx, sql, args...)
	if err != nil {
		return result, err
	}
	return model.ScanModelRows(rows)
}

func (t *StoreTransaction) GetWithOptions(options map[string]interface{}) ([]*model.ReturningStruct, error) {
	ctx := context.Background()
	sqlField, sqlWhere, valuesWhere := t.Connection.Conn.PrepareGet(options, t.model)
	sql := fmt.Sprintf("SELECT %v FROM %v.%v %v", sqlField, t.model.GetNameSchema(""), t.model.GetNameTable(), sqlWhere)
	log.Println(sql)
	return t.GetDataOfMapStuct(ctx, sql, valuesWhere...)
}
