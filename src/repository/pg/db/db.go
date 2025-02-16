package db

import (
	"context"
	"database/sql"
)

type TxProvider struct {
	db *sql.DB
}

func NewTxProvider(db *sql.DB) *TxProvider {
	return &TxProvider{db}
}

func (t *TxProvider) InWriteTx(ctx context.Context, f func(ctx context.Context, tx *sql.Tx) error) error {
	transaction, err := t.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	defer func() {
		if p := recover(); p != nil {
			_ = transaction.Rollback()

			panic(p)
		} else if err != nil {
			err = transaction.Rollback()
		} else {
			err = transaction.Commit()
		}
	}()

	err = f(ctx, transaction)

	return err
}

func (t *TxProvider) InReadTx(ctx context.Context, f func(ctx context.Context, tx *sql.Tx) error) error {
	transaction, err := t.db.BeginTx(ctx, &sql.TxOptions{ReadOnly: true})
	if err != nil {
		return err
	}

	defer func() {
		if p := recover(); p != nil {
			_ = transaction.Rollback()

			panic(p)
		} else if err != nil {
			err = transaction.Rollback()
		} else {
			err = transaction.Commit()
		}
	}()

	err = f(ctx, transaction)

	return err
}
