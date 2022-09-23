package store

import (
	"context"
	connectionPostgresql "ephorservices/pkg/db"
	model "ephorservices/pkg/model/schema/account/companypoint/model"
	modelInterface "ephorservices/pkg/model/schema/interface/model"
	parserTypes "ephorservices/pkg/parser/typeParse"
	"errors"
	"fmt"
	"log"
)

type StoreCompanyPoint struct {
	Connection      *connectionPostgresql.Manager
	model           modelInterface.Model
	fieldsReturning string
}

func NewStore(conn *connectionPostgresql.Manager) *StoreCompanyPoint {
	modelStore := model.Init()
	fields := conn.Conn.PrepareReturningFields(modelStore)
	return &StoreCompanyPoint{
		Connection:      conn,
		model:           modelStore,
		fieldsReturning: fields,
	}
}

func (scp *StoreCompanyPoint) Get(accountId interface{}) ([]*model.ReturningStruct, error) {
	ctx := context.Background()
	var options map[string]interface{}
	fieldSql, _, _ := scp.Connection.Conn.PrepareGet(options, scp.model)
	sql := fmt.Sprintf("SELECT %v FROM %v.%v", fieldSql, scp.model.GetNameSchema(parserTypes.ParseTypeInString(accountId)), scp.model.GetNameTable())
	return scp.GetDataOfMapStuct(ctx, sql)
}

func (scp *StoreCompanyPoint) AddByParams(parametrs map[string]interface{}) (*model.ReturningStruct, error) {
	if _, exist := parametrs["account_id"]; !exist {
		return &model.ReturningStruct{}, errors.New("not found paramentr account_id")
	}
	accountId := parserTypes.ParseTypeInString(parametrs["account_id"])
	delete(parametrs, "account_id")
	ctx := context.Background()
	sqlField, sqlValues, values := scp.Connection.Conn.PrepareInsert(parametrs, scp.model)
	sql := fmt.Sprintf("INSERT INTO %v.%v (%v) VALUES (%s) RETURNING %s", scp.model.GetNameSchema(accountId), scp.model.GetNameTable(), sqlField, sqlValues, scp.fieldsReturning)
	result, err := scp.GetDataOfStruct(ctx, sql, values...)
	if err != nil {
		return result, err
	}
	return result, nil
}

func (scp *StoreCompanyPoint) SetByParams(parametrs map[string]interface{}) (*model.ReturningStruct, error) {
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
	sqlValues, sqlWhere, Values := scp.Connection.Conn.PrepareUpdate(parametrs, options, scp.model)
	sql := fmt.Sprintf("UPDATE %v.%v SET %s %v", scp.model.GetNameSchema(accountId), scp.model.GetNameTable(), sqlValues, sqlWhere)
	result, err := scp.GetDataOfStruct(ctx, sql, Values...)
	if err != nil {
		return result, err
	}
	return result, nil
}

func (scp *StoreCompanyPoint) GetOneById(id, accountId interface{}) (*model.ReturningStruct, error) {
	var result *model.ReturningStruct
	var err error
	options := make(map[string]interface{})
	options["id"] = id
	sqlField, sqlWhere, valuesWhere := scp.Connection.Conn.PrepareGet(options, scp.model)
	sql := fmt.Sprintf("SELECT %v FROM %v.%v %v", sqlField, scp.model.GetNameSchema(parserTypes.ParseTypeInString(accountId)), scp.model.GetNameTable(), sqlWhere)
	ctx := context.Background()
	result, err = scp.GetDataOfStruct(ctx, sql, valuesWhere...)
	if err != nil {
		return result, err
	}
	return result, nil
}

func (scp *StoreCompanyPoint) GetDataOfStruct(ctx context.Context, sql string, args ...interface{}) (*model.ReturningStruct, error) {
	var result *model.ReturningStruct
	row, err := scp.Connection.Conn.QueryRow(ctx, sql, args...)
	if err != nil {
		return result, err
	}
	return model.ScanModelRow(row)
}

func (scp *StoreCompanyPoint) GetDataOfMapStuct(ctx context.Context, sql string, args ...interface{}) ([]*model.ReturningStruct, error) {
	var result []*model.ReturningStruct
	rows, err := scp.Connection.Conn.Query(ctx, sql, args...)
	if err != nil {
		return result, err
	}
	return model.ScanModelRows(rows)
}

func (scp *StoreCompanyPoint) GetWithOptions(options map[string]interface{}) ([]*model.ReturningStruct, error) {
	var result []*model.ReturningStruct
	if _, exist := options["account_id"]; !exist {
		return result, errors.New("not found paramentr account_id")
	}
	accountId := parserTypes.ParseTypeInString(options["account_id"])
	delete(options, "account_id")
	ctx := context.Background()
	sqlField, sqlWhere, valuesWhere := scp.Connection.Conn.PrepareGet(options, scp.model)
	sql := fmt.Sprintf("SELECT %v FROM %v.%v %v", sqlField, scp.model.GetNameSchema(accountId), scp.model.GetNameTable(), sqlWhere)
	log.Println(sql)
	return scp.GetDataOfMapStuct(ctx, sql, valuesWhere...)
}
