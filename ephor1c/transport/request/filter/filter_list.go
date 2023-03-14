package filter

import (
	filter "ephorservices/ephor1c/transport/request/interface/filter"
	"fmt"
)

type FilterList struct {
	Operator string
	Filters  []filter.Filter
}

func New(operator string) *FilterList {
	sliceFilter := make([]filter.Filter, 0, 10)
	return &FilterList{
		Filters:  sliceFilter,
		Operator: operator,
	}
}

func (fl *FilterList) Add(filter filter.Filter) {
	fl.Filters = append(fl.Filters, filter)
}

func (fl *FilterList) Get(name string) filter.Filter {
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

func (fl *FilterList) GetStringUrl() string {
	query := ""
	for _, filter := range fl.Filters {
		if len(query) == 0 {
			query = filter.GetUrl()
		} else {
			query = fmt.Sprintf("%s %s %s", query, fl.Operator, filter.GetUrl())
		}
	}
	return query
}

func (fl *FilterList) SetOperator(operator string) {
	fl.Operator = operator
}
