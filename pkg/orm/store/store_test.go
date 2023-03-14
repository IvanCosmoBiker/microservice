package store

import (
	"context"
	"database/sql"
	connectionPostgresql "ephorservices/pkg/orm/db"
	model_interface "ephorservices/pkg/orm/model/interface_model"
	"fmt"
	pgx "github.com/jackc/pgx/v5"
	"github.com/stretchr/testify/assert"
	"reflect"
	"strings"
	"testing"
)

var ConnectionDb *connectionPostgresql.Manager

func init() {
	var err error
	ConnectionDb, err = connectionPostgresql.Init("postgres", "123", "127.0.0.1", "local", uint16(5432), uint16(10), uint16(10), uint16(10), false, true, context.Background())
	if err != nil {
		panic(err)
	}
}

type ModelTest struct {
	Id     int64
	Field1 sql.NullInt32
	Field2 sql.NullString
}

func NewModel() model_interface.Model {
	return &ModelTest{}
}

func (mt *ModelTest) New() model_interface.Model {
	return &ModelTest{}
}

func (mt *ModelTest) GetNameSchema(accountNumber int) string {
	return fmt.Sprintf("%s%v", "test", accountNumber)
}

func (mt *ModelTest) GetNameTable() string {
	return "model_test"
}

func (mt *ModelTest) GetFields() []string {
	fields := make([]string, 0, 1)
	v := reflect.ValueOf(mt).Elem()
	typeOfS := v.Type()
	for i := 0; i < v.NumField(); i++ {
		field := typeOfS.Field(i).Name
		fields = append(fields, field)
	}
	return fields
}

func (mt *ModelTest) GetIdKey() int64 {
	return int64(mt.Id)
}

func (mt *ModelTest) ScanModelRow(row pgx.Row) (model_interface.Model, error) {
	model := &ModelTest{}
	err := row.Scan(model.Id,
		model.Field1,
		model.Field2)
	if err != nil {
		return model, err
	}
	return model, nil
}

func (mt *ModelTest) ScanModelRows(rows pgx.Rows) ([]model_interface.Model, error) {
	defer rows.Close()
	Result := make([]model_interface.Model, 0, 1)
	for rows.Next() {
		resultStruct, err := mt.ScanModelRow(rows)
		if err != nil {
			continue
		}
		Result = append(Result, resultStruct)
	}
	return Result, nil
}

func (mt *ModelTest) GetName() string {
	return "ModelTest"
}

func (whm *ModelTest) Get() map[string]interface{} {
	model := make(map[string]interface{})
	v := reflect.ValueOf(whm).Elem()
	typeOfS := v.Type()
	for i := 0; i < v.NumField(); i++ {
		field, ok := whm.GetNameFieldOfTag(typeOfS.Field(i))
		if !ok {
			field = strings.ToLower(typeOfS.Field(i).Name)
		}
		FieldValue := v.FieldByName(field)
		model[field] = whm.GetValueField(&FieldValue)
	}
	return model
}

func (whm *ModelTest) GetNameFieldOfTag(field reflect.StructField) (string, bool) {
	if value, ok := field.Tag.Lookup("json"); ok {
		return value, ok
	} else {
		return "", false
	}

}

func (whm *ModelTest) GetValueField(f *reflect.Value) interface{} {
	typeOfS := f.Type()
	switch typeOfS.Name() {
	case "NullInt32":
		return f.Interface().(sql.NullInt32).Int32
	case "NullString":
		return f.Interface().(sql.NullString).String
	case "int", "int32", "int8", "int16", "int64":
		return f.Int()
	case "uint", "uint8", "uint16", "uint32", "uint64":
		return f.Uint()
	case "string":
		return f.String()
	}
	return false
}

func TestAddParamModel(t *testing.T) {
	model := NewModel()
	store := New(model)
	store.SetAccountNumber(0)
	param := make(map[string]interface{})
	param["field1"] = 79
	param["field2"] = "test"
	testModel, _ := store.AddParamModel(param)
	expected := &ModelTest{
		Id: 0,
	}
	expected.Field1.Scan(79)
	expected.Field2.Scan("test")
	assert.Equal(t, expected, testModel)
}

func TestMakeSetKeys(t *testing.T) {
	model := NewModel()
	store := New(model)
	store.SetAccountNumber(0)
	test := &ModelTest{
		Id: 0,
	}
	test.Field2.Scan("test")
	fields, fieldSlice, _ := store.MakeSetKeys(test)
	expectedFields := "id,field2"
	expected := []string{"Field2"}
	assert.Equal(t, expected, fieldSlice)
	assert.Equal(t, expectedFields, fields)
}

func TestMakeSetValues(t *testing.T) {
	model := NewModel()
	store := New(model)
	store.SetAccountNumber(0)
	test := &ModelTest{
		Id: 0,
	}
	test.Field2.Scan("test")
	_, fields, _ := store.MakeSetKeys(test)
	sqlValues, values := store.MakeSetValues(test, fields)
	expectedSqlValues := "$1"
	expectedValues := make([]interface{}, 0, 1)
	expectedValues = append(expectedValues, int64(0))
	expectedValues = append(expectedValues, "test")
	assert.Equal(t, expectedSqlValues, sqlValues)
	assert.Equal(t, expectedValues, values)
}

func TestMakeAddValues(t *testing.T) {
	model := NewModel()
	store := New(model)
	store.SetAccountNumber(0)
	test := &ModelTest{
		Id: 5,
	}
	test.Field2.Scan("test")
	_, fieldSlice := store.MakeAddKeys(model)
	stringFields, fieldsSlice := store.MakeAddValues(test, fieldSlice)
	expected := make([]interface{}, 0, 1)
	expected = append(expected, int32(0), "test")
	assert.Equal(t, expected, fieldsSlice)
	assert.Equal(t, "$1, $2", stringFields)
}
