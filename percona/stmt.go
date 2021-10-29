package percona

import (
	"context"
	"database/sql"
)

// Stmt is a prepared statement.
type Stmt interface {
	// ExecContext executes a prepared statement with the given arguments and returns a Result summarizing the effect of
	// the statement.
	ExecContext(ctx context.Context, args ...interface{}) (sql.Result, error)

	// QueryContext executes a prepared query statement with the given arguments and returns the query results as a
	// *Rows.
	QueryContext(ctx context.Context, args ...interface{}) (*sql.Rows, error)

	// QueryRowContext executes a prepared query statement with the given arguments.
	QueryRowContext(ctx context.Context, args ...interface{}) *sql.Row

	// Close closes the statement.
	Close() error
}

// Preparer represents a service for creating prepared statement.
type Preparer interface {
	// PrepareContext returns prepared statement.
	PrepareContext(ctx context.Context, query string) (Stmt, error)
}
