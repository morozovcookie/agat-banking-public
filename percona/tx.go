package percona

import (
	"context"
	"database/sql"

	"github.com/pkg/errors"
)

// Tx is an in-progress database transaction.
type Tx interface {
	Preparer

	// Commit commits the transaction.
	Commit(ctx context.Context) error

	// Rollback aborts the transaction.
	Rollback(ctx context.Context) error
}

// TxBeginner represents a service for creating transaction.
type TxBeginner interface {
	// BeginTx starts a transaction.
	BeginTx(ctx context.Context, opts *sql.TxOptions) (Tx, error)
}

type tx struct {
	*sql.Tx
}

// PrepareContext returns prepared statement.
func (tx *tx) PrepareContext(ctx context.Context, query string) (Stmt, error) {
	res, err := tx.Tx.PrepareContext(ctx, query) // nolint:sqlclosecheck
	if err != nil {
		return nil, errors.Wrap(err, "tx prepare")
	}

	return &stmt{
		Stmt: res,
	}, nil
}

// Commit commits the transaction.
func (tx *tx) Commit(_ context.Context) error {
	return tx.Tx.Commit()
}

// Rollback aborts the transaction.
func (tx *tx) Rollback(_ context.Context) error {
	return tx.Tx.Rollback()
}
