package prometheus

import (
	"context"
	"database/sql"
	"time"

	"github.com/morozovcookie/agat-banking/percona"
	"github.com/pkg/errors"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/metric"
	"go.opentelemetry.io/otel/metric/unit"
)

var _ percona.Tx = (*tx)(nil)

type tx struct {
	wrapped percona.Tx

	meter metric.Meter
	attrs []attribute.KeyValue

	commitDuration metric.Int64Histogram
	rollbackCount  metric.Int64Counter
}

func newTx(source percona.Tx, meter metric.Meter, attrs ...attribute.KeyValue) (wrapper *tx, err error) {
	wrapper = &tx{
		wrapped: source,

		meter: meter,
		attrs: attrs,

		commitDuration: metric.Int64Histogram{},
		rollbackCount:  metric.Int64Counter{},
	}

	wrapper.commitDuration, err = meter.NewInt64Histogram("sql.tx.duration",
		metric.WithDescription("measures the duration of the commit SQL transaction"),
		metric.WithUnit(unit.Milliseconds))
	if err != nil {
		return nil, errors.Wrap(err, "init tx")
	}

	wrapper.rollbackCount, err = meter.NewInt64Counter("sql.tx.rollback",
		metric.WithDescription("measures the number of SQL transaction rollback"),
		metric.WithUnit(unit.Dimensionless))
	if err != nil {
		return nil, errors.Wrap(err, "init tx")
	}

	return wrapper, nil
}

// PrepareContext returns prepared statement.
func (tx *tx) PrepareContext(ctx context.Context, query string) (percona.Stmt, error) {
	res, err := tx.wrapped.PrepareContext(ctx, query)
	if err != nil {
		return nil, err // nolint:wrapcheck
	}

	return newStmt(res, query, tx.meter, tx.attrs...)
}

// Commit commits the transaction.
func (tx *tx) Commit(ctx context.Context) error {
	var (
		now = time.Now()
		err = tx.wrapped.Commit(ctx)
	)

	tx.commitDuration.Record(ctx, time.Since(now).Milliseconds(), tx.attrs...)

	return err // nolint:wrapcheck
}

// Rollback aborts the transaction.
func (tx *tx) Rollback(ctx context.Context) error {
	err := tx.wrapped.Rollback(ctx)

	tx.rollbackCount.Add(ctx, 1, tx.attrs...)

	return err // nolint:wrapcheck
}

var _ percona.TxBeginner = (*TxBeginner)(nil)

// TxBeginner represents a service for creating transaction.
type TxBeginner struct {
	wrapped percona.TxBeginner

	meter metric.Meter
	attrs []attribute.KeyValue
}

// NewTxBeginner returns a new TxBeginner instance.
func NewTxBeginner(beginner percona.TxBeginner, meter metric.Meter, attrs ...attribute.KeyValue) *TxBeginner {
	return &TxBeginner{
		wrapped: beginner,

		meter: meter,
		attrs: attrs,
	}
}

// BeginTx starts a transaction.
func (beginner *TxBeginner) BeginTx(ctx context.Context, opts *sql.TxOptions) (percona.Tx, error) {
	res, err := beginner.wrapped.BeginTx(ctx, opts)
	if err != nil {
		return nil, err // nolint:wrapcheck
	}

	return newTx(res, beginner.meter, beginner.attrs...)
}
