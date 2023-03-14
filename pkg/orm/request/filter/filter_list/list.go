package filter

import (
	filter_interface "ephorservices/pkg/orm/request/interface/filter"
	"fmt"
)

type FilterList struct {
	Operator string
	Filters  []filter_interface.Filter
}

func New(operator string) *FilterList {
	return &FilterList{
		Operator: operator,
		Filters:  make([]filter_interface.Filter, 0, 1),
	}
}

func (fl *FilterList) Add(filter ...filter_interface.Filter) {
	fl.Filters = append(fl.Filters, filter...)
}

func (fl *FilterList) Get(name string) filter_interface.Filter {
	for _, filter := range fl.Filters {
		if filter.GetName() == name {
			return filter
		}
	}
	return nil
}

func (fl *FilterList) Remove(name string) {
	for index, filter := range fl.Filters {
		if filter.GetName() == name {
			fl.Filters = append(fl.Filters[:index], fl.Filters[index+1:]...)
			break
		}
	}
}

func (fl *FilterList) getName() string {
	return "list"
}

func (fl *FilterList) GetOperator() string {
	return fl.Operator
}

func (fl *FilterList) GetCount() int {
	return len(fl.Filters)
}

func (fl *FilterList) GetSql() string {
	query := ""
	for _, filter := range fl.Filters {
		if len(query) == 0 {
			query = filter.GetSql()
		} else {
			query += fmt.Sprintf("%s %s", fl.Operator, filter.GetSql())
		}
	}
	if len(query) != 0 {
		query = fmt.Sprintf(" (%s) ", query)
	}
	return query
}
