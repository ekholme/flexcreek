package sqlite

import (
	"context"
	"database/sql"
)

// querier defines the set of methods that both *sql.DB and *sql.Tx satisfy.
// This allows service methods to be agnostic about whether they are operating
// within a transaction or directly on the database connection pool.
type querier interface {
	ExecContext(ctx context.Context, query string, args ...any) (sql.Result, error)
	QueryContext(ctx context.Context, query string, args ...any) (*sql.Rows, error)
	QueryRowContext(ctx context.Context, query string, args ...any) *sql.Row
}
