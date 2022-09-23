package model

type Model interface {
	GetNameSchema(account string) string
	GetNameTable() string
}
