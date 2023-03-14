package date_range

import (
	filter "ephorservices/ephor1c/transport/request/interface/filter"
	"fmt"
)

type FilterRange struct {
	Field    string
	From     string
	To       string
	Operator string
}

func New(field string, from, to string, operator string) filter.Filter {
	return &FilterRange{
		Field:    field,
		From:     from,
		To:       to,
		Operator: operator,
	}
}

func (fp *FilterRange) GetUrl() string {
	url := ""
	url = fmt.Sprintf("%s ge %s and %s le %s", fp.Field, fp.From, fp.Field, fp.To)
	return url
}

func (fp *FilterRange) GetName() string {
	return fp.Field
}
