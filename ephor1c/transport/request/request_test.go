package request

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestAddFilterParam(t *testing.T) {
	request := New("http", "OData_Test_Infobase", Resource_Catalog, "Склад")
	request.AddFilterParam("Catalog", 50, OperatorEqual)
	stringUrl := request.GetStringUrl()
	stringtest := "http://OData_Test_Infobase/odata/standard.odata/Catalog_Склад?$format=json&$filter=Catalog eq 50"
	assert.Equal(t, stringtest, stringUrl)
}

func TestRemoveFilter(t *testing.T) {
	request := New("http", "OData_Test_Infobase", Resource_Catalog, "Склад")
	request.AddFilterParam("Catalog", 50, OperatorEqual)
	request.RemoveFilter("Catalog")
	filterCatalog := request.GetFilter("Catalog")
	assert.Equal(t, nil, filterCatalog)
}

func TestAddFilterRange(t *testing.T) {
	dateFrom := "2023-02-09 07:13:25"
	dateTo := "2023-02-09 08:36:45"
	request := New("http", "OData_Test_Infobase", Resource_Catalog, "Склад")
	request.AddFilterRange("Catalog", dateFrom, dateTo, OperatorAnd)
	stringUrl := request.GetStringUrl()
	stringtest := "http://OData_Test_Infobase/odata/standard.odata/Catalog_Склад?$format=json&$filter=Catalog ge 2023-02-09 07:13:25 and Catalog le 2023-02-09 08:36:45"
	assert.Equal(t, stringtest, stringUrl)
}

func TestEmptyFilter(t *testing.T) {
	request := New("http", "OData_Test_Infobase", Resource_Catalog, "Склад")
	stringUrl := request.GetStringUrl()
	stringtest := "http://OData_Test_Infobase/odata/standard.odata/Catalog_Склад?$format=json"
	assert.Equal(t, stringtest, stringUrl)
}
