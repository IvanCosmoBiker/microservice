package filter_param

import (
	filter_interface "ephorservices/pkg/orm/request/interface/filter"
	parserTypes "ephorservices/pkg/parser/typeParse"
	"fmt"
	"strings"
)

var (
	OperatorMore         = ">"
	OperatorLess         = "<"
	OperatorEqual        = "="
	OperatorNotEqual     = "!="
	OperatorMoreAndEqual = ">="
	OperatorLessAndEqual = "<="
)

type FilterParam struct {
	Field         string
	Value         []interface{}
	Operator      string
	CaseSensitive bool
}

func New(field string, operator string, caseSensitive bool, value ...interface{}) filter_interface.Filter {
	values := make([]interface{}, 0, 1)
	values = append(values, value...)
	return &FilterParam{
		Field:         field,
		Value:         values,
		Operator:      operator,
		CaseSensitive: caseSensitive,
	}
}

func (fp *FilterParam) GetName() string {
	return fp.Field
}

func (fp *FilterParam) GetValue() []interface{} {
	return fp.Value
}

func (fp *FilterParam) IsEqual(value interface{}) bool {
	for _, v := range fp.Value {
		if v == value {
			return true
		}
	}
	return false
}

func (fp *FilterParam) GetOperator() string {
	return fp.Operator
}

func (fp *FilterParam) GetSql() string {
	if len(fp.Value) == 0 {
		switch fp.Operator {
		case OperatorEqual:
			return fmt.Sprintf("%s IS NULL", fp.Field)
		case OperatorNotEqual:
			return fmt.Sprintf("%s IS NOT NULL", fp.Field)
		}
	}
	sql := ""
	sqlSlice := make([]string, 0, 1)
	if !fp.CaseSensitive {
		sqlSlice = append(sqlSlice, fmt.Sprintf("LOWER(%s)", fp.Field))
	} else {
		sqlSlice = append(sqlSlice, fp.Field)
	}
	sqlSlice = append(sqlSlice, " IN (")
	isFirstItem := true
	for _, v := range fp.Value {
		if !isFirstItem {
			sqlSlice = append(sqlSlice, ", ")
		} else {
			isFirstItem = false
		}
		value := fp.EscapeString(parserTypes.ParseTypeInString(v))
		if !fp.CaseSensitive {
			sqlSlice = append(sqlSlice, fmt.Sprintf("LOWER(%s)", value))
		} else {
			sqlSlice = append(sqlSlice, value)
		}
	}
	sqlSlice = append(sqlSlice, " )")
	sql = strings.Join(sqlSlice, " ")
	return sql
}

func (fp *FilterParam) EscapeString(value string) string {
	return fmt.Sprintf("'%s'", value)
}
