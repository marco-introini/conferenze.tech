package db

import (
	"context"
	"database/sql"
)

type DB struct {
	*Queries
}

func WrapDB(queries *Queries) *DB {
	return &DB{Queries: queries}
}

func (db *DB) WithTransaction(ctx context.Context, fn func(Querier) error) error {
	tx, err := db.db.(interface {
		BeginTx(ctx context.Context, opts *sql.TxOptions) (*sql.Tx, error)
	}).BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	if err := fn(db.WithTx(tx)); err != nil {
		return err
	}

	return tx.Commit()
}
