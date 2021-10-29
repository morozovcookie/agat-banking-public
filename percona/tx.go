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
	Commit() error

	// Rollback aborts the transaction.
	Rollback() error
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
	stmt, err := tx.Tx.PrepareContext(ctx, query)
	if err != nil {
		return nil, errors.Wrap(err, "tx prepare")
	}

	return stmt, nil
}
