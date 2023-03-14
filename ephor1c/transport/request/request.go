package request

import (
	"encoding/base64"
	filterList "ephorservices/ephor1c/transport/request/filter"
	date_range "ephorservices/ephor1c/transport/request/filter/date_range"
	filter_param "ephorservices/ephor1c/transport/request/filter/param"
	filter "ephorservices/ephor1c/transport/request/interface/filter"
	"fmt"
)

var (
	Resource_Catalog                    = "Catalog_"
	Resource_Document                   = "Document_"
	Resource_DocumentJournal            = "DocumentJournal_"
	Resource_Constant                   = "Constant_"
	Resource_ExchangePlan               = "ExchangePlan_"
	Resource_ChartOfAccounts            = "ChartOfAccounts_"
	Resource_ChartOfCalculationTypes    = "ChartOfCalculationTypes_"
	Resource_ChartOfCharacteristicTypes = "ChartOfCharacteristicTypes_"
	Resource_InformationRegister        = "InformationRegister_"
	Resource_AccumulationRegister       = "AccumulationRegister_"
	Resource_CalculationRegister        = "CalculationRegister_"
	Resource_AccountingRegister         = "AccountingRegister_"
	Resource_BusinessProcess            = "BusinessProcess_"
	Resource_Task                       = "Task_"
)

var (
	OperatorMore         = "gt"
	OperatorLess         = "lt"
	OperatorEqual        = "eq"
	OperatorNotEqual     = "ne"
	OperatorMoreAndEqual = "ge"
	OperatorLessAndEqual = "le"
	OperatorOr           = "or"
	OperatorAnd          = "and"
	OperatorNot          = "not"
)

type Request1c struct {
	Protocol string
	Host     string
	ODATA    string
	Resource string
	Format   string
	Filter   *filterList.FilterList
}

func New(protocol, host, resource, resourceName string) *Request1c {
	return &Request1c{
		Filter:   filterList.New(OperatorAnd),
		ODATA:    "odata/standard.odata",
		Protocol: protocol,
		Host:     host,
		Format:   "$format=json",
		Resource: fmt.Sprintf("%s%s", resource, resourceName),
	}
}

func (r *Request1c) AddFilterParam(field string, value interface{}, operator string) {
	r.Filter.Add(filter_param.New(field, value, operator))
}

func (r *Request1c) AddFilterRange(field string, from, to string, operator string) {
	r.Filter.Add(date_range.New(field, from, to, operator))
}

func (r *Request1c) ClearFilter() {
	r.Filter = nil
}

func (r *Request1c) RemoveFilter(name string) {
	r.Filter.Remove(name)
}

func (r *Request1c) GetFilter(name string) filter.Filter {
	return r.Filter.Get(name)
}

func (r *Request1c) GetStringUrl() string {
	filterString := fmt.Sprintf("%s://%s/%s/%s?%s", r.Protocol, r.Host, r.ODATA, r.Resource, r.Format)
	filter := r.Filter.GetStringUrl()
	if len(filter) != 0 {
		filterString += "&$filter="
		filterString += filter
	}
	return filterString
}

func (r *Request1c) basicAuth(username, password string) string {
	auth := username + ":" + password
	return base64.StdEncoding.EncodeToString([]byte(auth))
}

func (r *Request1c) SetBasicAuth(username, password string) map[string]string {
	auth := make(map[string]string)
	auth["Authorization"] = "Basic " + r.basicAuth(username, password)
	return auth
}

func (r *Request1c) Clear() {
	r.Filter = nil
	r.Protocol = ""
	r.Host = ""
	r.ODATA = ""
	r.Resource = ""
	r.Format = ""
}
