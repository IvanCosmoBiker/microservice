package store

import (
	"context"
	"database/sql"
	connectionPostgresql "ephorservices/pkg/orm/db"
	model_interface "ephorservices/pkg/orm/model"
	request "ephorservices/pkg/orm/request"
	parserTypes "ephorservices/pkg/parser/typeParse"
	"errors"
	"fmt"
	"reflect"
	"strings"
)

var (
	Error_Not_Init_Connection = errors.New("Not initialization connection")
	Error_Empty_Data          = errors.New("Empty records")
	Error_Not_Set_Id          = errors.New("Not set value of field id")
)

type Store struct {
	AccountName   string
	AccountNumber int
	ModelName     string
	Model         model_interface.Model
	TableName     string
	conn          *connectionPostgresql.Manager
}

func New(model model_interface.Model) *Store {
	return &Store{
		Model:     model,
		ModelName: model.GetName(),
		TableName: model.GetNameTable(),
		conn:      connectionPostgresql.ConnectionManager,
	}
}

func (s *Store) checkConnection() error {
	if s.conn == nil {
		return Error_Not_Init_Connection
	}
	return nil
}

func (s *Store) SetAccountNumber(accountNumber int) {
	s.AccountNumber = accountNumber
	s.AccountName = s.Model.GetNameSchema(accountNumber)
}

func (s *Store) GetFields() []string {
	return s.Model.GetFields()
}

func (s *Store) GetStringField() string {
	return strings.Join(s.GetFields(), ", ")
}

func (s *Store) CreateTable() string {
	return ""
}

func (s *Store) Get(req *request.Request) ([]model_interface.Model, error) {
	ctx := context.Background()
	if err := s.checkConnection(); err != nil {
		return nil, err
	}
	where := ""
	order := ""
	limit := 0
	offset := 0
	limitString := ""
	if req != nil {
		filter := req.GetFilterlist()
		if filter != nil {
			where = fmt.Sprintf("WHERE %s", filter.GetSql())
		}
		sorter := req.GetSorter()
		if sorter != nil {
			order = fmt.Sprintf("ORDER BY %s", sorter.GetSql())
		}
		if req.GetLimit() != 0 {
			limit = req.GetLimit()
			offset = req.GetOffset()
			limitString = fmt.Sprintf("LIMIT %v OFFSET %v", limit, offset)
		}
	}
	sql := fmt.Sprintf("SELECT %s FROM %s.%s %s %s %s", s.GetStringField(), s.AccountName, s.TableName, where, order, limitString)
	fmt.Println(sql)
	rows, err := s.conn.Conn.Query(ctx, sql)
	if err != nil {
		return nil, err
	}
	return s.Model.ScanModelRows(rows)
}

func (s *Store) GetOneById(id int) (model_interface.Model, error) {

	req := request.New()
	req.AddFilterParam("id", request.OperatorEqual, true, id)
	models, err := s.Get(req)
	if err != nil {
		return nil, err
	}
	if len(models) < 1 {
		return nil, Error_Empty_Data
	}
	return models[0], nil
}

func (s *Store) GetOneBy(req *request.Request) (model_interface.Model, error) {
	models, err := s.Get(req)
	if err != nil {
		return nil, err
	}
	if len(models) == 0 {
		return nil, Error_Empty_Data
	}
	return models[0], nil
}

func (s *Store) Add(model model_interface.Model) (model_interface.Model, error) {
	if err := s.checkConnection(); err != nil {
		return nil, err
	}
	ctx := context.Background()
	fields, fieldSlice := s.MakeAddKeys(model)
	valuesString, valuesSlice := s.MakeAddValues(model, fieldSlice)
	sql := fmt.Sprintf("INSERT INTO %s.%s (%s) VALUES (%s) RETURNING %s", s.Model.GetNameSchema(s.AccountNumber), s.TableName, fields, valuesString, s.GetStringField())
	row, err := s.conn.Conn.QueryRow(ctx, sql, valuesSlice...)
	if err != nil {
		return nil, err
	}
	model, err = s.Model.ScanModelRow(row)
	return model, err
}

func (s *Store) AddByParams(params map[string]interface{}) (model_interface.Model, error) {
	model, err := s.AddParamModel(params)
	if err != nil {
		return model, err
	}
	model, err = s.Add(model)
	return model, err
}

func (s *Store) Set(model model_interface.Model) (model_interface.Model, error) {
	if err := s.checkConnection(); err != nil {
		return nil, err
	}
	ctx := context.Background()
	fields, fieldSlice, WhereId := s.MakeSetKeys(model)
	valuesString, valuesSlice := s.MakeSetValues(model, fieldSlice)
	sql := fmt.Sprintf("UPDATE %s.%s SET (%s) = (%s) WHERE %s RETURNING %s", s.Model.GetNameSchema(s.AccountNumber), s.TableName, fields, valuesString, WhereId, s.GetStringField())
	row, err := s.conn.Conn.QueryRow(ctx, sql, valuesSlice...)
	if err != nil {
		return nil, err
	}
	model, err = s.Model.ScanModelRow(row)
	return model, err
}

func (s *Store) SetByParams(params map[string]interface{}) (model model_interface.Model, err error) {
	_, ok := params["id"]
	if !ok {
		err = Error_Not_Set_Id
		return
	}
	model, err = s.AddParamModel(params)
	if err != nil {
		return
	}
	model, err = s.Set(model)
	return
}

func (s *Store) Delete(req *request.Request) error {
	ctx := context.Background()
	if err := s.checkConnection(); err != nil {
		return err
	}
	where := ""
	if req != nil {
		filter := req.GetFilterlist()
		if filter != nil {
			where = fmt.Sprintf("WHERE %s", filter.GetSql())
		}
	}
	sql := fmt.Sprintf("DELETE FROM %s.%s %s", s.Model.GetNameSchema(s.AccountNumber), s.Model.GetNameTable(), where)
	fmt.Println(sql)
	_, err := s.conn.Conn.Exec(ctx, sql)
	return err
}

func (s *Store) DeleteById(id int) error {
	req := request.New()
	req.AddFilterParam("id", request.OperatorEqual, true, id)
	return s.Delete(req)
}

func (s *Store) MakeAddKeys(model model_interface.Model) (string, []string) {
	result := ""
	fieldSlice := make([]string, 0, 1)
	v := reflect.ValueOf(model).Elem()
	for _, field := range s.GetFields() {
		if strings.ToLower(field) == "id" {
			continue
		}
		FieldValue := v.FieldByName(field)
		resultCheck := s.CheckValueField(&FieldValue)
		if resultCheck {
			if len(result) != 0 {
				result += ","
			}
			result += strings.ToLower(field)
			fieldSlice = append(fieldSlice, field)
		}
	}
	return result, fieldSlice
}

func (s *Store) MakeSetKeys(model model_interface.Model) (string, []string, string) {
	fieldSlice := make([]string, 0, 1)
	whereId := ""
	result := ""
	v := reflect.ValueOf(model).Elem()
	for _, field := range s.GetFields() {
		FieldValue := v.FieldByName(field)
		if strings.ToLower(field) == "id" {
			//result += strings.ToLower(field)
			whereId += fmt.Sprintf("%s=%v", strings.ToLower(field), FieldValue.Interface())
			continue
		}

		if FieldValue.Kind() != reflect.Struct {
			continue
		}
		resultCheck := s.CheckValueField(&FieldValue)
		if resultCheck {
			if len(result) != 0 {
				result += ","
			}
			result += strings.ToLower(field)
			fieldSlice = append(fieldSlice, field)
		}
	}
	return result, fieldSlice, whereId
}

func (s *Store) MakeAddValues(model model_interface.Model, fields []string) (string, []interface{}) {
	valuesPrepare, values := s.PrepareValues(model, fields)
	sqlValues := strings.Join(valuesPrepare, ", ")
	return sqlValues, values
}

func (s *Store) MakeSetValues(model model_interface.Model, fields []string) (string, []interface{}) {
	valuesPrepare, values := s.PrepareValues(model, fields)
	sqlValues := strings.Join(valuesPrepare, ", ")
	return sqlValues, values
}

func (s *Store) PrepareValues(model model_interface.Model, fields []string) ([]string, []interface{}) {
	v := reflect.ValueOf(model).Elem()
	values := make([]interface{}, 0, len(fields))
	valuesPrepare := make([]string, 0, len(fields))
	for i := 0; i < len(fields); i++ {
		value := v.FieldByName(fields[i])
		valuesPrepare = append(valuesPrepare, fmt.Sprintf("$%v", i+1))
		typeValue := value.Type().Name()
		values = append(values, s.EscapeAndGetValue(&value, typeValue))
	}
	return valuesPrepare, values
}

func (s *Store) AddParamModel(params map[string]interface{}) (model_interface.Model, error) {
	model := s.Model.New()
	ps := reflect.ValueOf(model)
	sValue := ps.Elem()
	fmt.Println(sValue.Type().Name())
	if sValue.Kind() == reflect.Struct {
		for _, field := range s.GetFields() {
			f := sValue.FieldByName(field)
			if f.IsValid() {
				if f.CanSet() {
					_, exist := params[strings.ToLower(field)]
					if !exist {
						continue
					}
					s.SetValueField(&f, params[strings.ToLower(field)])
				}
			}
		}
	}
	return model, nil
}

func (s *Store) SetValueField(f *reflect.Value, value interface{}) {
	if f.Kind() != reflect.Struct {
		switch f.Kind() {
		case reflect.String:
			f.SetString(parserTypes.ParseTypeInString(value))
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			f.SetInt(parserTypes.ParseTypeInterfaceToInt64(value))
		case reflect.Float32, reflect.Float64:
			f.SetFloat(parserTypes.ParseTypeInFloat64(value))
		case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
			f.SetInt(parserTypes.ParseTypeInterfaceToInt64(value))
		}
	} else {
		typeOfS := f.Type()
		for i := 0; i < f.NumField(); i++ {
			field := typeOfS.Field(i).Name
			FieldValue := f.FieldByName(field)
			switch field {
			case "Int32":
				FieldValue.SetInt(int64(parserTypes.ParseTypeInterfaceToInt32(value)))
			case "Valid":
				FieldValue.SetBool(true)
			case "String":
				FieldValue.SetString(parserTypes.ParseTypeInString(value))
			case "Int64":
				FieldValue.SetInt(parserTypes.ParseTypeInterfaceToInt64(value))
			}
		}
	}

}

func (s *Store) CheckValueField(f *reflect.Value) bool {
	typeOfS := f.Type()
	switch typeOfS.Name() {
	case "NullInt32":
		return f.Interface().(sql.NullInt32).Valid
	case "NullInt64":
		return f.Interface().(sql.NullInt64).Valid
	case "NullString":
		return f.Interface().(sql.NullString).Valid
	}
	return false
}

func (s *Store) EscapeAndGetValue(f *reflect.Value, typeName string) interface{} {
	switch typeName {
	case "NullInt32":
		return s.GetValueOfTypeSql(f.Interface().(sql.NullInt32).Int32, f.Interface().(sql.NullInt32).Valid)
	case "NullInt64":
		return s.GetValueOfTypeSql(f.Interface().(sql.NullInt64).Int64, f.Interface().(sql.NullInt64).Valid)
	case "NullString":
		return s.GetValueOfTypeSql(f.Interface().(sql.NullString).String, f.Interface().(sql.NullString).Valid)
	case "int", "int32", "int16", "int64", "int8", "string", "uint8", "uint", "uint16", "uint32", "uint64":
		return f.Interface()
	}
	return "null"
}

func (s *Store) GetValueOfTypeSql(value interface{}, valid bool) interface{} {
	if !valid {
		return "null"
	}
	return value
}
