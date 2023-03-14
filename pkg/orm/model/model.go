package interface_model

import (
	pgx "github.com/jackc/pgx/v5"
)

type Model interface {
	New() Model
	GetName() string
	GetNameSchema(accountNumber int) string
	GetNameTable() string
	GetFields() []string
	GetIdKey() int64
	ScanModelRow(row pgx.Row) (Model, error)
	ScanModelRows(rows pgx.Rows) ([]Model, error)
	Get() map[string]interface{}
}
