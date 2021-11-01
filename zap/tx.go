package zap

import (
	"context"
	"database/sql"

	"github.com/morozovcookie/agat-banking/percona"
	"go.uber.org/zap"
)

var _ percona.Tx = (*tx)(nil)

type tx struct {
	wrapped percona.Tx
	logger  *zap.Logger
}

func newTx(wrapped percona.Tx, logger *zap.Logger) *tx {
	return &tx{
		wrapped: wrapped,
		logger:  logger.With(zap.String("component", "Tx")),
	}
}

// PrepareContext returns prepared statement.
func (tx *tx) PrepareContext(ctx context.Context, query string) (percona.Stmt, error) {
	res, err := tx.wrapped.PrepareContext(ctx, query)

	tx.logger.Debug("prepare",
		zap.String("query", query),
		zap.Error(err))

	if err != nil {
		tx.logger.Error("prepare",
			zap.String("query", query),
			zap.Error(err))

		return nil, err // nolint:wrapcheck
	}

	return newStmt(res, query, tx.logger), nil
}

// Commit commits the transaction.
func (tx *tx) Commit(ctx context.Context) error {
	err := tx.wrapped.Commit(ctx)

	tx.logger.Debug("commit",
		zap.Error(err))

	if err != nil {
		tx.logger.Error("commit",
			zap.Error(err))

		return err // nolint:wrapcheck
	}

	return nil
}

// Rollback aborts the transaction.
func (tx *tx) Rollback(ctx context.Context) error {
	err := tx.wrapped.Rollback(ctx)

	tx.logger.Debug("rollback",
		zap.Error(err))

	if err != nil {
		tx.logger.Error("rollback",
			zap.Error(err))

		return err // nolint:wrapcheck
	}

	return nil
}

var _ percona.TxBeginner = (*TxBeginner)(nil)

// TxBeginner represents a service for creating transaction.
type TxBeginner struct {
	wrapped percona.TxBeginner
	logger  *zap.Logger
}

// NewTxBeginner returns a new TxBeginner instance.
func NewTxBeginner(beginner percona.TxBeginner, logger *zap.Logger) *TxBeginner {
	return &TxBeginner{
		wrapped: beginner,
		logger:  logger.With(zap.String("component", "TxBeginner")),
	}
}

// BeginTx starts a transaction.
func (beginner *TxBeginner) BeginTx(ctx context.Context, opts *sql.TxOptions) (percona.Tx, error) {
	res, err := beginner.wrapped.BeginTx(ctx, opts)

	beginner.logger.Debug("begin",
		zap.Error(err))

	if err != nil {
		beginner.logger.Error("begin",
			zap.Error(err))

		return nil, err // nolint:wrapcheck
	}

	return newTx(res, beginner.logger), nil
}
