package request

import (
	filter_list "ephorservices/pkg/orm/request/filter/filter_list"
	filter_param "ephorservices/pkg/orm/request/filter/filter_param"
	filter_range "ephorservices/pkg/orm/request/filter/filter_range"
	filter "ephorservices/pkg/orm/request/interface/filter"
	sorter "ephorservices/pkg/orm/request/sorter"
)

var (
	OperatorAnd          = "AND"
	OperatorOr           = "OR"
	OperatorMore         = ">"
	OperatorLess         = "<"
	OperatorEqual        = "="
	OperatorNotEqual     = "!="
	OperatorMoreAndEqual = ">="
	OperatorLessAndEqual = "<="

	/* Sorter variable */
	Asc  = "ASC"  // по возростанию
	Desc = "DESC" // по убыванию
)

type Operators struct {
	OperatorAnd          string
	OperatorOr           string
	OperatorMore         string
	OperatorLess         string
	OperatorEqual        string
	OperatorNotEqual     string
	OperatorMoreAndEqual string
	OperatorLessAndEqual string
	Asc                  string
	Desc                 string
}

type Request struct {
	Operator   *Operators
	Distinct   string
	Sorter     *sorter.Sorter
	Limit      int
	Offset     int
	FilterList *filter_list.FilterList
}

func New() *Request {
	return &Request{
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
}

func (r *Request) AddFilterParam(field string, operator string, caseSensitive bool, value ...interface{}) {
	r.FilterList.Add(filter_param.New(field, operator, caseSensitive, value...))
}

func (r *Request) AddFilterRange(field, from, to string) {
	r.FilterList.Add(filter_range.New(field, from, to))
}

func (r *Request) GetFilter(name string) filter.Filter {
	return r.FilterList.Get(name)
}

func (r *Request) GetFilterlist() *filter_list.FilterList {
	return r.FilterList
}

func (r *Request) RemoveFilter(name string) {
	r.FilterList.Remove(name)
}

func (r *Request) ClearFilter() {
	r.FilterList = nil
}

func (r *Request) SetSorter(field string, direction string, alias ...string) {
	if r.Sorter == nil {
		r.Sorter = sorter.New(field, direction, alias...)
	}
}

func (r *Request) GetSorter() *sorter.Sorter {
	return r.Sorter
}

func (r *Request) ClearSorter() {
	r.Sorter = nil
}

func (r *Request) SetLimit(limit int, offset ...int) {
	Offset := 0
	if len(offset) != 0 {
		Offset = offset[0]
	}
	r.Limit = limit
	r.Offset = Offset
}

func (r *Request) GetLimit() int {
	return r.Limit
}

func (r *Request) GetOffset() int {
	return r.Offset
}

func (r *Request) ClearLimit() {
	r.Limit = 0
	r.Offset = 0
}

func (r *Request) Clear() {
	r.ClearFilter()
	r.ClearLimit()
	r.ClearSorter()
}
