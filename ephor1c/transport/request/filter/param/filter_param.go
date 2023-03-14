package param

import (
	filter "ephorservices/ephor1c/transport/request/interface/filter"
	"fmt"
)

type FilterParam struct {
	Field    string
	Value    interface{}
	Operator string
}

func New(field string, value interface{}, operator string) filter.Filter {
	return &FilterParam{
		Field:    field,
		Value:    value,
		Operator: operator,
	}
}

func (fp *FilterParam) GetUrl() string {
	url := ""
	url = fmt.Sprintf("%s %s %v", fp.Field, fp.Operator, fp.Value)
	return url
}

func (fp *FilterParam) GetName() string {
	return fp.Field
}
