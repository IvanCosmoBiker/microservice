package filter

type Filter interface {
	GetName() string
	GetSql() string
}
