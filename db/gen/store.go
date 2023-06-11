package db

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v4"
)

// Store provides all function to execute db queries indiviually and transactions
type Store interface {
	Querier
	TransferTx(ctx context.Context, arg TransferTxParams) (TransferTxResult, error)
}

// Store provides all function to execute SQL queries transactions
type SQLStore struct {
	*Queries
	db *pgx.Conn
}

// NewStore creates a new store
func NewStore(db *pgx.Conn) Store {
	return SQLStore{
		db:      db,
		Queries: New(db),
	}
}

// execTx executes a function within a database transaction
func (store *SQLStore) execTx(ctx context.Context, fn func(*Queries) error) error {
	// var txOptions pgx.TxOptions
	tx, err := store.db.Begin(ctx)
	if err != nil {
		return err
	}
	q := New(tx)
	err = fn(q)
	if err != nil {
		if rollBackErr := tx.Rollback(ctx); rollBackErr != nil {
			return fmt.Errorf("tx err: %v, rollback err: %v", err, rollBackErr)
		}
		return err
	}
	return tx.Commit(ctx)
}
