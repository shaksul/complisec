package repo

import (
	"context"
	"database/sql"
)

// DBInterface - интерфейс для базы данных
type DBInterface interface {
	Query(query string, args ...interface{}) (*sql.Rows, error)
	QueryContext(ctx context.Context, query string, args ...interface{}) (*sql.Rows, error)
	QueryRow(query string, args ...interface{}) *sql.Row
	Exec(query string, args ...interface{}) (sql.Result, error)
	ExecContext(ctx context.Context, query string, args ...interface{}) (sql.Result, error)
}

// RiskRepoInterface - интерфейс для RiskRepo
type RiskRepoInterface interface {
	GetByIDWithTenant(ctx context.Context, id, tenantID string) (*Risk, error)
}
