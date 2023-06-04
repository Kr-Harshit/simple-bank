// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.17.2

package db

import (
	"time"
)

type Account struct {
	ID             int64     `db:"id" json:"id"`
	OwnerID        string    `db:"owner_id" json:"owner_id"`
	Balance        float32   `db:"balance" json:"balance"`
	Currency       string    `db:"currency" json:"currency"`
	CreatedAt      time.Time `db:"created_at" json:"created_at"`
	LastModifiedAt time.Time `db:"last_modified_at" json:"last_modified_at"`
}

type Entry struct {
	ID        int64 `db:"id" json:"id"`
	AccountID int64 `db:"account_id" json:"account_id"`
	// can be negative and positive
	Amount     float32   `db:"amount" json:"amount"`
	TransferID int64     `db:"transfer_id" json:"transfer_id"`
	CreatedAt  time.Time `db:"created_at" json:"created_at"`
	Credit     bool      `db:"credit" json:"credit"`
}

type Transfer struct {
	ID            int64 `db:"id" json:"id"`
	FromAccountID int64 `db:"from_account_id" json:"from_account_id"`
	ToAccountID   int64 `db:"to_account_id" json:"to_account_id"`
	// must be positive
	Amount    float32   `db:"amount" json:"amount"`
	CreatedAt time.Time `db:"created_at" json:"created_at"`
}
