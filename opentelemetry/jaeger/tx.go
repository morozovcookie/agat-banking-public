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

var _ percona.Tx = (*tx)(nil)

type tx struct {
	tracer  trace.Tracer
	wrapped percona.Tx
	attrs   []attribute.KeyValue
}

func newTx(tracer trace.Tracer, source percona.Tx, attrs ...attribute.KeyValue) *tx {
	return &tx{
		tracer:  tracer,
		wrapped: source,
		attrs:   attrs,
	}
}

// PrepareContext returns prepared statement.
func (tx *tx) PrepareContext(ctx context.Context, query string) (percona.Stmt, error) {
	attrs := append(tx.attrs, opentelemetry.SQLAttributesFromQuery(query)...)

	ctx, span := tx.tracer.Start(ctx, "Tx.PrepareContext", trace.WithAttributes(attrs...))
	defer span.End()

	res, err := tx.wrapped.PrepareContext(ctx, query)
	if err != nil {
		span.SetStatus(codes.Error, err.Error())

		return nil, err // nolint:wrapcheck
	}

	span.SetStatus(codes.Ok, "")

	return newStmt(tx.tracer, res, tx.attrs...), nil
}

// Commit commits the transaction.
func (tx *tx) Commit(ctx context.Context) error {
	ctx, span := tx.tracer.Start(ctx, "Tx.Commit", trace.WithAttributes(tx.attrs...))
	defer span.End()

	if err := tx.wrapped.Commit(ctx); err != nil {
		span.SetStatus(codes.Error, err.Error())

		return err // nolint:wrapcheck
	}

	span.SetStatus(codes.Ok, "")

	return nil
}

// Rollback aborts the transaction.
func (tx *tx) Rollback(ctx context.Context) error {
	ctx, span := tx.tracer.Start(ctx, "Tx.Rollback", trace.WithAttributes(tx.attrs...))
	defer span.End()

	if err := tx.wrapped.Rollback(ctx); err != nil {
		span.SetStatus(codes.Error, err.Error())

		return err // nolint:wrapcheck
	}

	span.SetStatus(codes.Ok, "")

	return nil
}

var _ percona.TxBeginner = (*TxBeginner)(nil)

// TxBeginner represents a service for creating transaction.
type TxBeginner struct {
	tracer  trace.Tracer
	wrapped percona.TxBeginner
	attrs   []attribute.KeyValue
}

// NewTxBeginner returns a new TxBeginner instance.
func NewTxBeginner(tracer trace.Tracer, beginner percona.TxBeginner, attrs ...attribute.KeyValue) *TxBeginner {
	return &TxBeginner{
		tracer:  tracer,
		wrapped: beginner,
		attrs:   attrs,
	}
}

// BeginTx starts a transaction.
func (beginner *TxBeginner) BeginTx(ctx context.Context, opts *sql.TxOptions) (percona.Tx, error) {
	ctx, span := beginner.tracer.Start(ctx, "TxBeginner.BeginTx", trace.WithAttributes(beginner.attrs...))
	defer span.End()

	res, err := beginner.wrapped.BeginTx(ctx, opts)
	if err != nil {
		span.SetStatus(codes.Error, err.Error())

		return nil, err // nolint:wrapcheck
	}

	span.SetStatus(codes.Ok, "")

	return newTx(beginner.tracer, res, beginner.attrs...), nil
}
