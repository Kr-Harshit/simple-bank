package db

import (
	"context"
	"testing"
	"time"

	util "github.com/KHarshit1203/simple-bank/util/datagenerator"
	"github.com/jackc/pgx/v4"
	"github.com/stretchr/testify/require"
)

func generateAccount(ctx context.Context, t *testing.T) Account {
	arg := CreateAccountParams{
		OwnerID:  util.RandomUUID(),
		Balance:  util.RandomFloat(1, 1000),
		Currency: util.RandomCurrency(),
	}
	account, err := testQueries.CreateAccount(ctx, arg)

	require.NoError(t, err)
	require.NotEmpty(t, account)

	require.Equal(t, arg.OwnerID, account.OwnerID)
	require.Equal(t, arg.Balance, account.Balance)
	require.Equal(t, arg.Currency, account.Currency)

	require.NotZero(t, account.ID)
	require.NotZero(t, account.CreatedAt)

	return account
}

func deleteGeneratedAccount(ctx context.Context, id int64, t *testing.T) {
	err := testQueries.DeleteAccount(ctx, id)
	require.NoError(t, err)
}

func TestCreateAccount(t *testing.T) {
	ctx := context.Background()

	account := generateAccount(ctx, t)
	deleteGeneratedAccount(ctx, account.ID, t)

}

func TestGetAccount(t *testing.T) {
	ctx := context.Background()
	account1 := generateAccount(ctx, t)

	account2, err := testQueries.GetAccount(ctx, account1.ID)

	require.NoError(t, err)
	require.NotEmpty(t, account2)

	require.Equal(t, account1.ID, account2.ID)
	require.Equal(t, account1.OwnerID, account2.OwnerID)
	require.Equal(t, account1.Balance, account2.Balance)
	require.Equal(t, account1.Currency, account2.Currency)
	require.Equal(t, account1.CreatedAt, account2.CreatedAt)

	deleteGeneratedAccount(ctx, account1.ID, t)
}

func TestListAccounts(t *testing.T) {
	ctx := context.Background()
	var lastAccount Account
	generateAccounts := make([]Account, 10)
	for i := 0; i < 10; i++ {
		lastAccount = generateAccount(ctx, t)
		generateAccounts = append(generateAccounts, lastAccount)
	}

	arg := ListAccountsParams{
		OwnerID: lastAccount.OwnerID,
		Limit:   5,
		Offset:  0,
	}

	accounts, err := testQueries.ListAccounts(ctx, arg)
	require.NoError(t, err)
	require.NotEmpty(t, accounts)

	for _, account := range accounts {
		require.NotEmpty(t, account)
		require.Equal(t, lastAccount.OwnerID, account.OwnerID)
	}

	for _, account := range generateAccounts {
		deleteGeneratedAccount(ctx, account.ID, t)
	}
}

func TestDeleteAccount(t *testing.T) {
	ctx := context.Background()
	account1 := generateAccount(ctx, t)

	err := testQueries.DeleteAccount(ctx, account1.ID)
	require.NoError(t, err)

	account2, err := testQueries.GetAccount(ctx, account1.ID)
	require.Error(t, err)
	require.EqualError(t, err, pgx.ErrNoRows.Error())
	require.Empty(t, account2)
}

func TestPurgeUserAccounts(t *testing.T) {
	ctx := context.Background()
	var lastAccount Account
	generatedAccounts := make([]Account, 10)
	for i := 0; i < 10; i++ {
		lastAccount = generateAccount(ctx, t)
		generatedAccounts = append(generatedAccounts, lastAccount)
	}

	err := testQueries.PurgeUserAccounts(ctx, lastAccount.OwnerID)
	require.NoError(t, err)

	listAccountsAgr := ListAccountsParams{
		OwnerID: lastAccount.OwnerID,
		Limit:   100,
		Offset:  0,
	}

	accounts, _ := testQueries.ListAccounts(ctx, listAccountsAgr)
	require.Empty(t, accounts)

	for _, account := range generatedAccounts {
		if account.ID != lastAccount.ID {
			deleteGeneratedAccount(ctx, account.ID, t)
		}
	}
}

func TestUpdateBalance(t *testing.T) {
	ctx := context.Background()
	account1 := generateAccount(ctx, t)

	arg := UpdateBalanceParams{
		ID:     account1.ID,
		Amount: util.RandomFloat(100, 10000),
	}

	account2, err := testQueries.UpdateBalance(ctx, arg)
	require.NoError(t, err)
	require.NotEmpty(t, account2)

	require.Equal(t, account1.ID, account2.ID)
	require.Equal(t, account1.OwnerID, account2.OwnerID)
	require.Equal(t, account1.Balance+arg.Amount, account2.Balance)
	require.Equal(t, account1.Currency, account2.Currency)

	require.WithinDuration(t, account1.CreatedAt, account2.CreatedAt, time.Second)
}
