package datastore

import (
	"context"
	"database/sql"

	"github.com/jmoiron/sqlx"
)

type Executor interface {
	ExecContext(context.Context, string, ...any) (sql.Result, error)
	QueryContext(context.Context, string, ...any) (*sql.Rows, error)
	QueryxContext(context.Context, string, ...any) (*sqlx.Rows, error)
	QueryRowContext(context.Context, string, ...any) *sql.Row
	QueryRowxContext(context.Context, string, ...any) *sqlx.Row
	Rebind(query string) string
	SelectContext(ctx context.Context, dest interface{}, query string, args ...any) error
	GetContext(ctx context.Context, dest interface{}, query string, args ...any) error
	NamedExecContext(ctx context.Context, query string, args any) (sql.Result, error)
}
