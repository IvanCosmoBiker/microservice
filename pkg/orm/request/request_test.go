package request

import (
	filter_list "ephorservices/pkg/orm/request/filter/filter_list"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestRequestNew(t *testing.T) {
	request := New()
	expRequest := &Request{
		Operator: &Operators{
			OperatorAnd:          "AND",
			OperatorOr:           "OR",
			OperatorMore:         ">",
			OperatorLess:         "<",
			OperatorEqual:        "=",
			OperatorNotEqual:     "!=",
			OperatorMoreAndEqual: ">=",
			OperatorLessAndEqual: "<=",
			Asc:                  "ASC",
			Desc:                 "DESC",
		},
		FilterList: filter_list.New(OperatorAnd),
		Limit:      1000,
		Offset:     0,
	}
	assert.Equal(t, expRequest, request)
}

func TestAddFilterParam(t *testing.T) {
	request := New()
	request.AddFilterParam("test", OperatorMore, true, 1, 2, 3, 4, 5)
	filterTest := request.GetFilter("test")
	assert.NotEmpty(t, filterTest)
}

func TestAddFilterRange(t *testing.T) {
	request := New()
	request.AddFilterRange("date", "2023-01-01 00:00:00", "2023-01-01 23:59:59")
	filterTest := request.GetFilter("date")
	assert.NotEmpty(t, filterTest)
}

func TestGetSqlFilterParam(t *testing.T) {
	request := New()
	request.AddFilterParam("test", OperatorMore, true, 1, 2, 3, 4, 5)
	filterTest := request.GetFilter("test")
	sqlTest := filterTest.GetSql()
	expected := "test  IN ( '1' ,  '2' ,  '3' ,  '4' ,  '5'  )"
	assert.Equal(t, expected, sqlTest)
}

func TestGetSqlFilterRange(t *testing.T) {
	request := New()
	request.AddFilterRange("date", "2023-01-01 00:00:00", "2023-01-01 23:59:59")
	filterTest := request.GetFilter("date")
	sqlTest := filterTest.GetSql()
	expected := "'2023-01-01 00:00:00' <= date AND date <= '2023-01-01 23:59:59'"
	assert.Equal(t, expected, sqlTest)
}

func TestRemoveFilter(t *testing.T) {
	request := New()
	request.AddFilterRange("date", "2023-01-01 00:00:00", "2023-01-01 23:59:59")
	request.AddFilterParam("test", OperatorMore, true, 1, 2, 3, 4, 5)
	request.RemoveFilter("date")
	filterTest := request.GetFilter("date")
	assert.Empty(t, filterTest)
	filterTest2 := request.GetFilter("test")
	assert.NotEmpty(t, filterTest2)
	expected := "test  IN ( '1' ,  '2' ,  '3' ,  '4' ,  '5'  )"
	assert.Equal(t, expected, filterTest2.GetSql())
}

func TestNullFilter(t *testing.T) {
	request := New()
	request.AddFilterParam("date", request.Operator.OperatorEqual, true)
	filterTest := request.GetFilter("date")
	expected := "date IS NULL"
	filterTest.GetSql()
	assert.Equal(t, expected, filterTest.GetSql())
}
