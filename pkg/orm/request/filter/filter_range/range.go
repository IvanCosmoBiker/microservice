package filter_range

import (
	filter_interface "ephorservices/pkg/orm/request/interface/filter"
	"fmt"
)

type FilterRange struct {
	Field string
	From  string
	To    string
}

func New(field, from, to string) filter_interface.Filter {
	return &FilterRange{
		Field: field,
		From:  from,
		To:    to,
	}
}

func (fr *FilterRange) GetName() string {
	return fr.Field
}

func (fr *FilterRange) GetFrom() string {
	return fr.From
}

func (fr *FilterRange) GetTo() string {
	return fr.To
}

func (fr *FilterRange) GetSql() string {
	sql := fmt.Sprintf("'%s' <= %s", fr.From, fr.Field)
	if len(fr.To) != 0 {
		sql += fmt.Sprintf(" AND %s <= '%s'", fr.Field, fr.To)
	}
	return sql
}
