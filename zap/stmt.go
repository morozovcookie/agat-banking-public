package zap

import (
	"context"
	"database/sql"

	"github.com/morozovcookie/agat-banking/percona"
	"go.uber.org/zap"
)

var _ percona.Stmt = (*stmt)(nil)

type stmt struct {
	wrapped percona.Stmt
	query   string

	logger *zap.Logger
}

func newStmt(wrapped percona.Stmt, query string, logger *zap.Logger) *stmt {
	return &stmt{
		wrapped: wrapped,
		query:   query,

		logger: logger.With(zap.String("component", "Stmt")),
	}
}

// ExecContext executes a prepared statement with the given arguments and returns a Result summarizing the effect of
// the statement.
func (stmt *stmt) ExecContext(ctx context.Context, args ...interface{}) (sql.Result, error) {
	res, err := stmt.wrapped.ExecContext(ctx, args...)

	stmt.logger.Debug("exec",
		zap.String("query", stmt.query),
		zap.Any("args", args),
		zap.Error(err))

	if err != nil {
		stmt.logger.Error("exec",
			zap.String("query", stmt.query),
			zap.Any("args", args),
			zap.Error(err))

		return nil, err // nolint:wrapcheck
	}

	return res, nil
}

// QueryContext executes a prepared query statement with the given arguments and returns the query results as a
// *Rows.
func (stmt *stmt) QueryContext(ctx context.Context, args ...interface{}) (*sql.Rows, error) {
	rr, err := stmt.wrapped.QueryContext(ctx, args...)

	stmt.logger.Debug("query",
		zap.String("query", stmt.query),
		zap.Any("args", args),
		zap.Error(err))

	if err != nil {
		stmt.logger.Error("query",
			zap.String("query", stmt.query),
			zap.Any("args", args),
			zap.Error(err))

		return nil, err // nolint:wrapcheck
	}

	return rr, nil
}

// QueryRowContext executes a prepared query statement with the given arguments.
func (stmt *stmt) QueryRowContext(ctx context.Context, args ...interface{}) *sql.Row {
	row := stmt.wrapped.QueryRowContext(ctx, args...)

	stmt.logger.Debug("query row",
		zap.String("query", stmt.query),
		zap.Any("args", args))

	return row
}

// Close closes the statement.
func (stmt *stmt) Close(ctx context.Context) error {
	err := stmt.wrapped.Close(ctx)

	stmt.logger.Debug("close",
		zap.Error(err))

	if err != nil {
		stmt.logger.Error("close",
			zap.Error(err))

		return err // nolint:wrapchek
	}

	return nil
}

var _ percona.Preparer = (*Preparer)(nil)

// Preparer represents a service for creating prepared statement.
type Preparer struct {
	wrapped percona.Preparer
	logger  *zap.Logger
}

// NewPreparer returns a new Preparer instance.
func NewPreparer(preparer percona.Preparer, logger *zap.Logger) *Preparer {
	return &Preparer{
		wrapped: preparer,
		logger:  logger.With(zap.String("component", "Preparer")),
	}
}

// PrepareContext returns prepared statement.
func (p *Preparer) PrepareContext(ctx context.Context, query string) (percona.Stmt, error) {
	res, err := p.wrapped.PrepareContext(ctx, query)

	p.logger.Debug("prepare",
		zap.String("query", query),
		zap.Error(err))

	if err != nil {
		p.logger.Error("prepare",
			zap.String("query", query),
			zap.Error(err))

		return nil, err // nolint:wrapcheck
	}

	return newStmt(res, query, p.logger), nil
}
