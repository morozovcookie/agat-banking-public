package prometheus

import (
	"context"
	"database/sql"
	"time"

	"github.com/morozovcookie/agat-banking/opentelemetry"
	"github.com/morozovcookie/agat-banking/percona"
	"github.com/pkg/errors"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/metric"
	"go.opentelemetry.io/otel/metric/unit"
)

var _ percona.Stmt = (*stmt)(nil)

// Stmt is a prepared statement.
type stmt struct {
	wrapped percona.Stmt
	query   string

	attrs []attribute.KeyValue

	queryExecutionDuration metric.Int64Histogram
	queryErrors            metric.Int64Counter
}

func newStmt(
	source percona.Stmt,
	query string,
	meter metric.Meter,
	attrs ...attribute.KeyValue,
) (
	wrapper *stmt,
	err error,
) {
	wrapper = &stmt{
		wrapped: source,
		query:   query,

		attrs: append(attrs, opentelemetry.SQLAttributesFromQuery(query)...),

		queryExecutionDuration: metric.Int64Histogram{},
		queryErrors:            metric.Int64Counter{},
	}

	wrapper.queryExecutionDuration, err = meter.NewInt64Histogram("sql.query.duration",
		metric.WithDescription("measures the duration of the SQL statement execution"),
		metric.WithUnit(unit.Milliseconds))
	if err != nil {
		return nil, errors.Wrap(err, "init stmt")
	}

	wrapper.queryErrors, err = meter.NewInt64Counter("sql.query.error",
		metric.WithDescription("measures the number of SQL statement execution errors"),
		metric.WithUnit(unit.Dimensionless))
	if err != nil {
		return nil, errors.Wrap(err, "init stmt")
	}

	return wrapper, nil
}

// ExecContext executes a prepared statement with the given arguments and returns a Result summarizing the effect of
// the statement.
func (stmt *stmt) ExecContext(ctx context.Context, args ...interface{}) (sql.Result, error) {
	var (
		now      = time.Now()
		res, err = stmt.wrapped.ExecContext(ctx, args...)
	)

	stmt.queryExecutionDuration.Record(ctx, time.Since(now).Milliseconds(), stmt.attrs...)

	if err != nil {
		stmt.queryErrors.Add(ctx, 1, stmt.attrs...)

		return nil, err // nolint:wrapcheck
	}

	return res, nil
}

// QueryContext executes a prepared query statement with the given arguments and returns the query results as a
// *Rows.
func (stmt *stmt) QueryContext(ctx context.Context, args ...interface{}) (*sql.Rows, error) {
	var (
		now     = time.Now()
		rr, err = stmt.wrapped.QueryContext(ctx, args...)
	)

	stmt.queryExecutionDuration.Record(ctx, time.Since(now).Milliseconds(), stmt.attrs...)

	if err != nil {
		stmt.queryErrors.Add(ctx, 1, stmt.attrs...)

		return nil, err // nolint:wrapcheck
	}

	return rr, nil
}

// QueryRowContext executes a prepared query statement with the given arguments.
func (stmt *stmt) QueryRowContext(ctx context.Context, args ...interface{}) *sql.Row {
	var (
		now = time.Now()
		row = stmt.wrapped.QueryRowContext(ctx, args...)
	)

	stmt.queryExecutionDuration.Record(ctx, time.Since(now).Milliseconds(), stmt.attrs...)

	return row
}

// Close closes the statement.
func (stmt *stmt) Close(ctx context.Context) error {
	return stmt.wrapped.Close(ctx)
}

var _ percona.Preparer = (*Preparer)(nil)

// Preparer represents a service for creating prepared statement.
type Preparer struct {
	wrapped percona.Preparer

	meter metric.Meter
	attrs []attribute.KeyValue
}

// NewPreparer returns a new Preparer instance.
func NewPreparer(preparer percona.Preparer, meter metric.Meter, attrs ...attribute.KeyValue) *Preparer {
	return &Preparer{
		wrapped: preparer,

		meter: meter,
		attrs: attrs,
	}
}

// PrepareContext returns prepared statement.
func (p *Preparer) PrepareContext(ctx context.Context, query string) (percona.Stmt, error) {
	res, err := p.wrapped.PrepareContext(ctx, query)
	if err != nil {
		return nil, err // nolint:wrapcheck
	}

	return newStmt(res, query, p.meter, p.attrs...)
}
