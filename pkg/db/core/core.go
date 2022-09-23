package core

import (
	"context"
	conInit "ephorservices/pkg/db/coninit"
	parserTypes "ephorservices/pkg/parser/typeParse"
	"fmt"
	"reflect"
	"strings"

	"github.com/jackc/pgconn"
	pgx "github.com/jackc/pgx/v4"
	pgxpool "github.com/jackc/pgx/v4/pgxpool"
)

const (
	Status_Idle int = iota
	Status_Work
	Status_Err
)

type CoreConnection struct {
	Status     int
	Connection *conInit.Connection
	ctx        context.Context
}

func NewCoreConnection(ctx context.Context, conn *conInit.Connection) *CoreConnection {
	return &CoreConnection{
		Status:     Status_Work,
		Connection: conn,
		ctx:        ctx,
	}
}

func (cc *CoreConnection) GetConnectionPool() *pgxpool.Pool {
	return cc.Connection.GetConn()
}

func (cc *CoreConnection) acquire() (*pgxpool.Conn, error) {
	connection, err := cc.Connection.Acquire()
	if err != nil {
		cc.Status = Status_Err
	}
	return connection, err
}

func (cc *CoreConnection) Query(ctx context.Context, sql string, args ...interface{}) (pgx.Rows, error) {
	var rows pgx.Rows
	var errRows error
	rows, errRows = cc.Connection.GetConn().Query(ctx, sql, args)
	if errRows != nil {
		return rows, errRows
	}
	return rows, nil
}

func (cc *CoreConnection) QueryRow(ctx context.Context, sql string, args ...interface{}) (pgx.Row, error) {
	row := cc.Connection.GetConn().QueryRow(ctx, sql, args...)
	return row, nil
}

func (cc *CoreConnection) Exec(ctx context.Context, sql string, args ...interface{}) (pgconn.CommandTag, error) {
	var tag pgconn.CommandTag
	var errTag error
	tag, errTag = cc.Connection.GetConn().Exec(ctx, sql, args)
	return tag, errTag
}

func (cc *CoreConnection) Close() {
	cc.Connection.Close()
}

func (cc *CoreConnection) PrepareGet(options map[string]interface{}, modelData interface{}) (string, string, []interface{}) {
	parametrs := make(map[string]interface{})
	fields := cc.PrepareFields(parametrs, modelData)
	sqlField := strings.Join(fields, ", ")
	if len(options) < 1 {
		emptyWhereValues := make([]interface{}, 0)
		return sqlField, "", emptyWhereValues
	}
	countFields := 0
	valuesWherePrepare, valuesWhere := cc.PrepareWhere(options, countFields)
	sqlWhere := " WHERE "
	sqlWhere += strings.Join(valuesWherePrepare, " ")
	return sqlField, sqlWhere, valuesWhere
}

func (cc *CoreConnection) PrepareUpdate(parametrs, options map[string]interface{}, modelData interface{}) (string, string, []interface{}) {
	fields := cc.PrepareFields(parametrs, modelData)
	countFields := len(parametrs)
	valuesPrepare, values := cc.PrepareValuesUpdate(parametrs, fields)
	var Values []interface{}
	if len(options) < 1 {
		emptyWhereValues := make([]interface{}, 0, 0)
		sqlValues := strings.Join(valuesPrepare, ", ")
		Values = append(Values, values...)
		return sqlValues, "", emptyWhereValues
	}
	valuesWherePrepare, valuesWhere := cc.PrepareWhere(options, countFields)
	sqlWhere := " WHERE "
	sqlWhere += strings.Join(valuesWherePrepare, " ")
	sqlValues := strings.Join(valuesPrepare, ", ")
	Values = append(Values, values...)
	Values = append(Values, valuesWhere...)
	return sqlValues, sqlWhere, Values
}

func (cc *CoreConnection) PrepareInsert(parametrs map[string]interface{}, modelData interface{}) (string, string, []interface{}) {
	fields := cc.PrepareFields(parametrs, modelData)
	valuesPrepare, values := cc.PrepareValues(parametrs, fields)
	sqlField := strings.Join(fields, ", ")
	sqlValues := strings.Join(valuesPrepare, ", ")
	return sqlField, sqlValues, values
}

func (cc *CoreConnection) PrepareReturningFields(modelData interface{}) string {
	v := reflect.ValueOf(modelData).Elem()
	typeOfS := v.Type()
	var stringField []string
	for i := 0; i < v.NumField(); i++ {
		field := typeOfS.Field(i).Name
		stringField = append(stringField, strings.ToLower(field))
	}
	return strings.Join(stringField, ", ")
}

func (cc *CoreConnection) PrepareFields(parametrs map[string]interface{}, modelData interface{}) []string {
	v := reflect.ValueOf(modelData).Elem()
	typeOfS := v.Type()
	var stringField []string
	for i := 0; i < v.NumField(); i++ {
		field := typeOfS.Field(i).Name
		if len(parametrs) >= 1 {
			_, exist := parametrs[strings.ToLower(field)]
			if exist != true {
				continue
			}
			cc.CheckType(parametrs, typeOfS.Field(i), strings.ToLower(field))
			stringField = append(stringField, strings.ToLower(field))
		} else {
			stringField = append(stringField, strings.ToLower(field))
		}
	}
	return stringField
}

func (cc *CoreConnection) PrepareWhere(options map[string]interface{}, countField int) ([]string, []interface{}) {
	var values []interface{}
	var valuesPrepare []string
	count := 0
	for key, _ := range options {
		value := options[key]
		if value != nil {
			countField += 1
			if count == 0 {
				valuesPrepare = append(valuesPrepare, fmt.Sprintf(" %s = $%v", key, countField))
			} else {
				valuesPrepare = append(valuesPrepare, fmt.Sprintf(" AND %s = $%v", key, countField))
			}
			values = append(values, value)
			count++
		} else {
			if count == 0 {
				valuesPrepare = append(valuesPrepare, fmt.Sprintf(" %s IS NULL", key))
			} else {
				valuesPrepare = append(valuesPrepare, fmt.Sprintf(" AND %s IS NULL", key))
			}
			count++
		}

	}
	return valuesPrepare, values
}

func (cc *CoreConnection) PrepareValuesUpdate(parametrs map[string]interface{}, fields []string) ([]string, []interface{}) {
	var values []interface{}
	var valuesPrepare []string
	count := 0
	for i := 0; i < len(fields); i++ {
		value, exist := parametrs[fields[i]]
		if !exist {
			continue
		} else {
			valuesPrepare = append(valuesPrepare, fmt.Sprintf("%s = $%v", fields[i], count+1))
			values = append(values, value)
			count++
		}
	}
	return valuesPrepare, values
}

func (cc *CoreConnection) PrepareValues(parametrs map[string]interface{}, fields []string) ([]string, []interface{}) {
	var values []interface{}
	var valuesPrepare []string
	for i := 0; i < len(fields); i++ {
		value, exist := parametrs[fields[i]]
		if !exist {
			continue
		} else {
			valuesPrepare = append(valuesPrepare, fmt.Sprintf("$%v", i+1))
			values = append(values, value)
		}
	}
	return valuesPrepare, values
}

func (cc *CoreConnection) CheckType(parametrs map[string]interface{}, field reflect.StructField, namefield string) {
	switch field.Type.Name() {
	case "NullInt64", "NullInt32":
		parametrs[namefield] = parserTypes.ParseTypeInterfaceToInt(parametrs[namefield])
	case "NullString":
		if parametrs[namefield] == nil {
			parametrs[namefield] = "NULL"
		} else {
			parametrs[namefield] = parserTypes.ParseTypeInString(parametrs[namefield])
		}
	}
}
