package filter

type Filter interface {
	GetName() string
	GetUrl() string
}
