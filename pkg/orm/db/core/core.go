package core

import (
	"context"
	conInit "ephorservices/pkg/orm/db/coninit"
	parserTypes "ephorservices/pkg/parser/typeParse"
	"fmt"
	"log"
	"reflect"
	"strings"
	"time"

	pgx "github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	pgxpool "github.com/jackc/pgx/v5/pgxpool"
)

const (
	Status_Idle int = iota
	Status_Work
	Status_Err
	Status_Warning
)

type CoreConnection struct {
	Status                   int
	Connection               *conInit.Connection
	ctx                      context.Context
	BackOffPolicyAcquireConn []time.Duration
}

func NewCoreConnection(ctx context.Context, conn *conInit.Connection) *CoreConnection {
	backOffPolicyAcquireConn := setBackOffPolicyAcquireConn()
	core := &CoreConnection{
		Status:                   Status_Work,
		Connection:               conn,
		ctx:                      ctx,
		BackOffPolicyAcquireConn: backOffPolicyAcquireConn,
	}
	return core
}

func setBackOffPolicyAcquireConn() []time.Duration {
	return []time.Duration{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
}

func (cc *CoreConnection) GetConnectionPool() *pgxpool.Pool {
	return cc.Connection.GetConn()
}

func (cc *CoreConnection) acquire(ctx context.Context) (connection *pgxpool.Pool, err error) {
	connection = cc.Connection.GetConn()
	return connection, err
}

func (cc *CoreConnection) Query(ctx context.Context, sql string, args ...interface{}) (pgx.Rows, error) {
	var rows pgx.Rows
	var errRows error
	conn, err := cc.acquire(ctx)
	if err != nil {
		log.Println(fmt.Errorf("Unable to acquire a database connection: %s", err.Error()))
		return rows, err
	}
	rows, errRows = conn.Query(ctx, sql, args...)
	if errRows != nil {
		log.Println(errRows)
		return rows, errRows
	}
	return rows, nil
}

func (cc *CoreConnection) QueryRow(ctx context.Context, sql string, args ...interface{}) (pgx.Row, error) {
	var row pgx.Row
	conn, err := cc.acquire(ctx)
	if err != nil {
		log.Println(fmt.Errorf("Unable to acquire a database connection: %s", err.Error()))
		return row, err
	}
	row = conn.QueryRow(ctx, sql, args...)
	return row, nil
}

func (cc *CoreConnection) Exec(ctx context.Context, sql string, args ...interface{}) (int64, error) {
	var tag pgconn.CommandTag
	var errTag error
	conn, err := cc.acquire(ctx)
	if err != nil {
		log.Println(fmt.Errorf("Unable to acquire a database connection: %s", err.Error()))
		return 0, err
	}
	tag, errTag = conn.Exec(ctx, sql, args...)
	if errTag != nil {
		log.Println(errTag)
		return 0, errTag
	}
	count := tag.RowsAffected()
	return count, nil
}

func (cc *CoreConnection) Close() {
	cc.Connection.Close()
}

func (cc *CoreConnection) PrepareGet(options map[string]interface{}, modelData interface{}) (sqlField string, sqlWhere string, valuesWhere []interface{}) {
	parametrs := make(map[string]interface{})
	var valuesWherePrepare []string
	fields, err := cc.PrepareFields(parametrs, modelData)
	if err != nil {
		return sqlField, sqlWhere, valuesWhere
	}
	sqlField = strings.Join(fields, ", ")
	if len(options) < 1 {
		emptyWhereValues := make([]interface{}, 0)
		return sqlField, "", emptyWhereValues
	}
	countFields := 0
	valuesWherePrepare, valuesWhere = cc.PrepareWhere(options, countFields)
	sqlWhere = " WHERE "
	sqlWhere += strings.Join(valuesWherePrepare, " ")
	return sqlField, sqlWhere, valuesWhere
}

func (cc *CoreConnection) PrepareUpdate(parametrs, options map[string]interface{}, modelData interface{}) (sqlValues string, sqlWhere string, Values []interface{}) {
	fields, err := cc.PrepareFields(parametrs, modelData)
	if err != nil {
		return sqlValues, sqlWhere, Values
	}
	countFields := len(parametrs)
	valuesPrepare, values := cc.PrepareValuesUpdate(parametrs, fields)
	if len(options) < 1 {
		emptyWhereValues := make([]interface{}, 0)
		sqlValues := strings.Join(valuesPrepare, ", ")
		Values = append(Values, values...)
		return sqlValues, "", emptyWhereValues
	}
	valuesWherePrepare, valuesWhere := cc.PrepareWhere(options, countFields)
	sqlWhere = " WHERE "
	sqlWhere += strings.Join(valuesWherePrepare, " ")
	sqlValues = strings.Join(valuesPrepare, ", ")
	Values = append(Values, values...)
	Values = append(Values, valuesWhere...)
	return sqlValues, sqlWhere, Values
}

func (cc *CoreConnection) PrepareInsert(parametrs map[string]interface{}, modelData interface{}) (sqlField string, sqlValues string, values []interface{}) {
	var valuesPrepare []string
	fields, err := cc.PrepareFields(parametrs, modelData)
	if err != nil {
		return sqlField, sqlValues, values
	}
	valuesPrepare, values = cc.PrepareValues(parametrs, fields)
	sqlField = strings.Join(fields, ", ")
	sqlValues = strings.Join(valuesPrepare, ", ")
	return sqlField, sqlValues, values
}

func (cc *CoreConnection) PrepareInsertListParamentrs(parametrs []map[string]interface{}, modelData interface{}) (sqlField string, sqlValues string, Values []interface{}) {
	var count int = 0
	fields, err := cc.PrepareFields(parametrs[0], modelData)
	if err != nil {
		return sqlField, sqlValues, Values
	}
	sqlField = strings.Join(fields, ", ")
	for i := 0; i < len(parametrs); i++ {
		valuesPrepare, values := cc.PrepareListValues(parametrs[i], fields, &count)
		Values = append(Values, values)
		sqlValues += strings.Join(valuesPrepare, ", ")
	}
	return sqlField, sqlValues, Values
}

func (cc *CoreConnection) PrepareReturningFields(modelData interface{}) (fields string, err error) {
	defer func() (err error) {
		if r := recover(); r != nil {
			err = fmt.Errorf("%v", r)
			return
		}
		return err
	}()
	v := reflect.ValueOf(modelData).Elem()
	typeOfS := v.Type()
	var stringField []string
	for i := 0; i < v.NumField(); i++ {
		field := typeOfS.Field(i).Name
		stringField = append(stringField, strings.ToLower(field))
	}
	fields = strings.Join(stringField, ", ")
	return fields, err
}

func (cc *CoreConnection) PrepareFields(parametrs map[string]interface{}, modelData interface{}) (stringField []string, err error) {
	defer func() (err error) {
		if r := recover(); r != nil {
			err = fmt.Errorf("%v", r)
			return
		}
		return err
	}()
	v := reflect.ValueOf(modelData).Elem()
	typeOfS := v.Type()
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
	return stringField, nil
}

func (cc *CoreConnection) PrepareWhere(options map[string]interface{}, countField int) (valuesPrepare []string, values []interface{}) {
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

func (cc *CoreConnection) PrepareValuesUpdate(parametrs map[string]interface{}, fields []string) (valuesPrepare []string, values []interface{}) {
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

func (cc *CoreConnection) PrepareValues(parametrs map[string]interface{}, fields []string) (valuesPrepare []string, values []interface{}) {
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

func (cc *CoreConnection) PrepareListValues(parametrs map[string]interface{}, fields []string, count *int) (valuesPrepare []string, values []interface{}) {
	for i := 0; i < len(fields); i++ {
		value, exist := parametrs[fields[i]]
		if !exist {
			continue
		} else {
			valuesPrepare = append(valuesPrepare, fmt.Sprintf("$%v", *count+1))
			values = append(values, value)
		}
	}
	return valuesPrepare, values
}

func (cc *CoreConnection) CheckType(parametrs map[string]interface{}, field reflect.StructField, namefield string) {
	switch field.Type.Name() {
	case "NullInt64", "NullInt32", "int", "int8", "int16", "int32", "int64", "uint", "uint8", "uint32", "uint16", "uint64":
		parametrs[namefield] = parserTypes.ParseTypeInterfaceToInt(parametrs[namefield])
	case "NullString", "string":
		if parametrs[namefield] == nil {
			parametrs[namefield] = "NULL"
		} else {
			parametrs[namefield] = parserTypes.ParseTypeInString(parametrs[namefield])
		}
	}
}
