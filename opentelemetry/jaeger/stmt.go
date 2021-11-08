package jaeger

import (
	"context"
	"database/sql"

	"github.com/morozovcookie/agat-banking/opentelemetry"
	"github.com/morozovcookie/agat-banking/percona"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
)

var _ percona.Stmt = (*stmt)(nil)

// Stmt is a prepared statement.
type stmt struct {
	tracer  trace.Tracer
	wrapped percona.Stmt
	attrs   []attribute.KeyValue
}

func newStmt(tracer trace.Tracer, source percona.Stmt, attrs ...attribute.KeyValue) *stmt {
	return &stmt{
		tracer:  tracer,
		wrapped: source,
		attrs:   attrs,
	}
}

// ExecContext executes a prepared statement with the given arguments and returns a Result summarizing the effect of
// the statement.
func (stmt *stmt) ExecContext(ctx context.Context, args ...interface{}) (sql.Result, error) {
	ctx, span := stmt.tracer.Start(ctx, "Stmt.ExecContext", trace.WithAttributes(stmt.attrs...))
	defer span.End()

	res, err := stmt.wrapped.ExecContext(ctx, args...)
	if err != nil {
		span.SetStatus(codes.Error, err.Error())

		return nil, err // nolint:wrapcheck
	}

	span.SetStatus(codes.Ok, "")

	return res, nil
}

// QueryContext executes a prepared query statement with the given arguments and returns the query results as a
// *Rows.
func (stmt *stmt) QueryContext(ctx context.Context, args ...interface{}) (*sql.Rows, error) {
	ctx, span := stmt.tracer.Start(ctx, "Stmt.QueryContext", trace.WithAttributes(stmt.attrs...))
	defer span.End()

	rr, err := stmt.wrapped.QueryContext(ctx, args...)
	if err != nil {
		span.SetStatus(codes.Error, err.Error())

		return nil, err // nolint:wrapcheck
	}

	span.SetStatus(codes.Ok, "")

	return rr, nil
}

// QueryRowContext executes a prepared query statement with the given arguments.
func (stmt *stmt) QueryRowContext(ctx context.Context, args ...interface{}) *sql.Row {
	ctx, span := stmt.tracer.Start(ctx, "Stmt.QueryRowContext", trace.WithAttributes(stmt.attrs...))
	defer span.End()

	row := stmt.wrapped.QueryRowContext(ctx, args...)

	span.SetStatus(codes.Ok, "")

	return row
}

// Close closes the statement.
func (stmt *stmt) Close(ctx context.Context) error {
	ctx, span := stmt.tracer.Start(ctx, "Stmt.Close", trace.WithAttributes(stmt.attrs...))
	defer span.End()

	if err := stmt.wrapped.Close(ctx); err != nil {
		span.SetStatus(codes.Error, err.Error())

		return err // nolint:wrapcheck
	}

	span.SetStatus(codes.Ok, "")

	return nil
}

var _ percona.Preparer = (*Preparer)(nil)

// Preparer represents a service for creating prepared statement.
type Preparer struct {
	tracer  trace.Tracer
	wrapped percona.Preparer
	attrs   []attribute.KeyValue
}

// NewPreparer returns a new Preparer instance.
func NewPreparer(tracer trace.Tracer, preparer percona.Preparer, attrs ...attribute.KeyValue) *Preparer {
	return &Preparer{
		tracer:  tracer,
		wrapped: preparer,
		attrs:   attrs,
	}
}

// PrepareContext returns prepared statement.
func (p *Preparer) PrepareContext(ctx context.Context, query string) (percona.Stmt, error) {
	attrs := append(p.attrs, opentelemetry.SQLAttributesFromQuery(query)...)

	ctx, span := p.tracer.Start(ctx, "Preparer.PrepareContext", trace.WithAttributes(attrs...))
	defer span.End()

	res, err := p.wrapped.PrepareContext(ctx, query)
	if err != nil {
		span.SetStatus(codes.Error, err.Error())

		return nil, err // nolint:wrapcheck
	}

	span.SetStatus(codes.Ok, "")

	return newStmt(p.tracer, res, attrs...), nil
}
