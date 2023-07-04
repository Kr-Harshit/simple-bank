package db

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
)

// Store provides all function to execute db queries indiviually and transactions
type Store interface {
	Querier
	TransferTx(ctx context.Context, arg TransferTxParams) (TransferTxResult, error)
}

// Store provides all function to execute SQL queries transactions
type SQLStore struct {
	*Queries
	connPool *pgxpool.Pool
}

// NewStore creates a new store
func NewSQLStore(connPool *pgxpool.Pool) Store {
	return SQLStore{
		connPool: connPool,
		Queries:  New(connPool),
	}
}
