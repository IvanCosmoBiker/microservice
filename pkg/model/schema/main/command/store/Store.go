package store

import (
	"context"
	connectionPostgresql "ephorservices/pkg/db"
	modelInterface "ephorservices/pkg/model/schema/interface/model"
	model "ephorservices/pkg/model/schema/main/command/model"
	"errors"
	"fmt"
	"log"
)

type StoreCommand struct {
	Connection      *connectionPostgresql.Manager
	model           modelInterface.Model
	fieldsReturning string
}

func NewStore(conn *connectionPostgresql.Manager) *StoreCommand {
	modelStore := model.Init()
	fields := conn.Conn.PrepareReturningFields(modelStore)
	return &StoreCommand{
		Connection:      conn,
		model:           modelStore,
		fieldsReturning: fields,
	}
}

func (sc *StoreCommand) Get() ([]*model.ReturningStruct, error) {
	ctx := context.Background()
	var options map[string]interface{}
	fieldSql, _, _ := sc.Connection.Conn.PrepareGet(options, sc.model)
	sql := fmt.Sprintf("SELECT %v FROM %v.%v", fieldSql, sc.model.GetNameSchema(""), sc.model.GetNameTable())
	return sc.GetDataOfMapStuct(ctx, sql)
}

func (sc *StoreCommand) AddByParams(parametrs map[string]interface{}) (*model.ReturningStruct, error) {
	ctx := context.Background()
	sqlField, sqlValues, values := sc.Connection.Conn.PrepareInsert(parametrs, sc.model)
	sql := fmt.Sprintf("INSERT INTO %v.%v (%v) VALUES (%s) RETURNING %s", sc.model.GetNameSchema(""), sc.model.GetNameTable(), sqlField, sqlValues, sc.fieldsReturning)
	result, err := sc.GetDataOfStruct(ctx, sql, values...)
	if err != nil {
		return result, err
	}
	return result, nil
}

func (sc *StoreCommand) SetByParams(parametrs map[string]interface{}) (*model.ReturningStruct, error) {
	if _, exist := parametrs["id"]; exist != true {
		return &model.ReturningStruct{}, errors.New("not found id in parametrs")
	}
	id := parametrs["id"]
	delete(parametrs, "id")
	options := make(map[string]interface{})
	options["id"] = id
	ctx := context.Background()
	sqlValues, sqlWhere, Values := sc.Connection.Conn.PrepareUpdate(parametrs, options, sc.model)
	sql := fmt.Sprintf("UPDATE %v.%v SET %s %v", sc.model.GetNameSchema(""), sc.model.GetNameTable(), sqlValues, sqlWhere)
	result, err := sc.GetDataOfStruct(ctx, sql, Values...)
	if err != nil {
		return result, err
	}
	return result, nil
}

func (sc *StoreCommand) GetOneById(id interface{}) (*model.ReturningStruct, error) {
	var result *model.ReturningStruct
	var err error
	options := make(map[string]interface{})
	options["id"] = id
	sqlField, sqlWhere, valuesWhere := sc.Connection.Conn.PrepareGet(options, sc.model)
	sql := fmt.Sprintf("SELECT %v FROM %v.%v %v", sqlField, sc.model.GetNameSchema(""), sc.model.GetNameTable(), sqlWhere)
	ctx := context.Background()
	result, err = sc.GetDataOfStruct(ctx, sql, valuesWhere...)
	if err != nil {
		return result, err
	}
	return result, nil
}

func (sc *StoreCommand) GetDataOfStruct(ctx context.Context, sql string, args ...interface{}) (*model.ReturningStruct, error) {
	var result *model.ReturningStruct
	row, err := sc.Connection.Conn.QueryRow(ctx, sql, args...)
	if err != nil {
		return result, err
	}
	return model.ScanModelRow(row)
}

func (sc *StoreCommand) GetDataOfMapStuct(ctx context.Context, sql string, args ...interface{}) ([]*model.ReturningStruct, error) {
	var result []*model.ReturningStruct
	rows, err := sc.Connection.Conn.Query(ctx, sql, args...)
	if err != nil {
		return result, err
	}
	return model.ScanModelRows(rows)
}

func (sc *StoreCommand) GetWithOptions(options map[string]interface{}) ([]*model.ReturningStruct, error) {
	ctx := context.Background()
	sqlField, sqlWhere, valuesWhere := sc.Connection.Conn.PrepareGet(options, sc.model)
	sql := fmt.Sprintf("SELECT %v FROM %v.%v %v", sqlField, sc.model.GetNameSchema(""), sc.model.GetNameTable(), sqlWhere)
	log.Println(sql)
	return sc.GetDataOfMapStuct(ctx, sql, valuesWhere...)
}
