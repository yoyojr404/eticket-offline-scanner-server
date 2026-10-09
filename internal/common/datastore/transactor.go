package datastore

import (
	"context"
	"fmt"

	"github.com/jmoiron/sqlx"
)

type Transactor interface {
	WithinTx(ctx context.Context, fn func(tx Executor) error) error
}

type transactor struct {
	db *sqlx.DB
}

func NewTransactor(db *sqlx.DB) Transactor {
	return &transactor{db: db}
}

func (t *transactor) WithinTx(ctx context.Context, fn func(tx Executor) error) error {
	tx, err := t.db.BeginTxx(ctx, nil)
	if err != nil {
		return fmt.Errorf("begin transaction: %w", err)
	}

	committed := false
	defer func() {
		if !committed {
			_ = tx.Rollback()
		}
	}()

	// Pass the sqlx.Tx directly to the callback as an Executor interface
	if err := fn(tx); err != nil {
		return err
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("commit transaction: %w", err)
	}

	committed = true
	return nil
}
