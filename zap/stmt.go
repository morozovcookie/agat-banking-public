package zap

import (
	"context"
	"database/sql"

	"github.com/morozovcookie/agat-banking/percona"
	"go.uber.org/zap"
)

var _ percona.Stmt = (*stmt)(nil)

type stmt struct {
	loggerCreator LoggerCreator
	wrapped       percona.Stmt
	query         string
}

func newStmt(creator LoggerCreator, wrapped percona.Stmt, query string) *stmt {
	return &stmt{
		loggerCreator: creator,
		wrapped:       wrapped,
		query:         query,
	}
}

// ExecContext executes a prepared statement with the given arguments and returns a Result summarizing the effect of
// the statement.
func (stmt *stmt) ExecContext(ctx context.Context, args ...interface{}) (sql.Result, error) {
	logger := stmt.loggerCreator.CreateLogger(ctx, "Stmt", "ExecContext")

	res, err := stmt.wrapped.ExecContext(ctx, args...)

	logger.Debug("exec", zap.String("query", stmt.query), zap.Any("args", args), zap.Error(err))

	if err != nil {
		logger.Error("exec", zap.String("query", stmt.query), zap.Any("args", args), zap.Error(err))

		return nil, err // nolint:wrapcheck
	}

	return res, nil
}

// QueryContext executes a prepared query statement with the given arguments and returns the query results as a
// *Rows.
func (stmt *stmt) QueryContext(ctx context.Context, args ...interface{}) (*sql.Rows, error) {
	logger := stmt.loggerCreator.CreateLogger(ctx, "Stmt", "QueryContext")

	rr, err := stmt.wrapped.QueryContext(ctx, args...)

	logger.Debug("query", zap.String("query", stmt.query), zap.Any("args", args), zap.Error(err))

	if err != nil {
		logger.Error("query", zap.String("query", stmt.query), zap.Any("args", args), zap.Error(err))

		return nil, err // nolint:wrapcheck
	}

	return rr, nil
}

// QueryRowContext executes a prepared query statement with the given arguments.
func (stmt *stmt) QueryRowContext(ctx context.Context, args ...interface{}) *sql.Row {
	logger := stmt.loggerCreator.CreateLogger(ctx, "Stmt", "QueryRowContext")

	row := stmt.wrapped.QueryRowContext(ctx, args...)

	logger.Debug("query row", zap.String("query", stmt.query), zap.Any("args", args))

	return row
}

// Close closes the statement.
func (stmt *stmt) Close(ctx context.Context) error {
	logger := stmt.loggerCreator.CreateLogger(ctx, "Stmt", "Close")

	err := stmt.wrapped.Close(ctx)

	logger.Debug("close", zap.Error(err))

	if err != nil {
		logger.Error("close", zap.Error(err))

		return err // nolint:wrapchek
	}

	return nil
}

var _ percona.Preparer = (*Preparer)(nil)

// Preparer represents a service for creating prepared statement.
type Preparer struct {
	loggerCreator LoggerCreator
	wrapped       percona.Preparer
}

// NewPreparer returns a new Preparer instance.
func NewPreparer(creator LoggerCreator, preparer percona.Preparer) *Preparer {
	return &Preparer{
		loggerCreator: creator,
		wrapped:       preparer,
	}
}

// PrepareContext returns prepared statement.
func (p *Preparer) PrepareContext(ctx context.Context, query string) (percona.Stmt, error) {
	logger := p.loggerCreator.CreateLogger(ctx, "Preparer", "PrepareContext")

	res, err := p.wrapped.PrepareContext(ctx, query)

	logger.Debug("prepare", zap.String("query", query), zap.Error(err))

	if err != nil {
		logger.Error("prepare", zap.String("query", query), zap.Error(err))

		return nil, err // nolint:wrapcheck
	}

	return newStmt(p.loggerCreator, res, query), nil
}
