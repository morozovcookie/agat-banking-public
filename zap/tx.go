package zap

import (
	"context"
	"database/sql"

	"github.com/morozovcookie/agat-banking/percona"
	"go.uber.org/zap"
)

var _ percona.Tx = (*tx)(nil)

type tx struct {
	loggerCreator LoggerCreator
	wrapped       percona.Tx
}

func newTx(creator LoggerCreator, wrapped percona.Tx) *tx {
	return &tx{
		loggerCreator: creator,
		wrapped:       wrapped,
	}
}

// PrepareContext returns prepared statement.
func (tx *tx) PrepareContext(ctx context.Context, query string) (percona.Stmt, error) {
	logger := tx.loggerCreator.CreateLogger(ctx, "Tx", "PrepareContext")

	res, err := tx.wrapped.PrepareContext(ctx, query)

	logger.Debug("prepare", zap.String("query", query), zap.Error(err))

	if err != nil {
		logger.Error("prepare", zap.String("query", query), zap.Error(err))

		return nil, err // nolint:wrapcheck
	}

	return newStmt(tx.loggerCreator, res, query), nil
}

// Commit commits the transaction.
func (tx *tx) Commit(ctx context.Context) error {
	logger := tx.loggerCreator.CreateLogger(ctx, "Tx", "Commit")

	err := tx.wrapped.Commit(ctx)

	logger.Debug("commit", zap.Error(err))

	if err != nil {
		logger.Error("commit", zap.Error(err))

		return err // nolint:wrapcheck
	}

	return nil
}

// Rollback aborts the transaction.
func (tx *tx) Rollback(ctx context.Context) error {
	logger := tx.loggerCreator.CreateLogger(ctx, "Tx", "Rollback")

	err := tx.wrapped.Rollback(ctx)

	logger.Debug("rollback", zap.Error(err))

	if err != nil {
		logger.Error("rollback", zap.Error(err))

		return err // nolint:wrapcheck
	}

	return nil
}

var _ percona.TxBeginner = (*TxBeginner)(nil)

// TxBeginner represents a service for creating transaction.
type TxBeginner struct {
	loggerCreator LoggerCreator
	wrapped       percona.TxBeginner
}

// NewTxBeginner returns a new TxBeginner instance.
func NewTxBeginner(creator LoggerCreator, beginner percona.TxBeginner) *TxBeginner {
	return &TxBeginner{
		loggerCreator: creator,
		wrapped:       beginner,
	}
}

// BeginTx starts a transaction.
func (beginner *TxBeginner) BeginTx(ctx context.Context, opts *sql.TxOptions) (percona.Tx, error) {
	logger := beginner.loggerCreator.CreateLogger(ctx, "TxBeginner", "BeginTx")

	res, err := beginner.wrapped.BeginTx(ctx, opts)

	logger.Debug("begin", zap.Error(err))

	if err != nil {
		logger.Error("begin", zap.Error(err))

		return nil, err // nolint:wrapcheck
	}

	return newTx(beginner.loggerCreator, res), nil
}
