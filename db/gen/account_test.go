package db

import (
	"context"
	"testing"
	"time"

	util "github.com/KHarshit1203/simple-bank/util"
	"github.com/jackc/pgx"
	"github.com/stretchr/testify/require"
)

func generateAccount(t *testing.T) Account {
	user := generateUser(t)

	arg := CreateAccountParams{
		Owner:    user.Username,
		Balance:  util.RandomFloat(100, 1000),
		Currency: util.RandomCurrency(),
	}

	account, err := testStore.CreateAccount(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, account)

	require.Equal(t, arg.Owner, account.Owner)
	require.Equal(t, arg.Balance, account.Balance)
	require.Equal(t, arg.Currency, account.Currency)

	require.NotZero(t, account.ID)
	require.NotZero(t, account.CreatedAt)

	return account
}

func TestCreateAccount(t *testing.T) {
	generateAccount(t)
}

func TestGetAccount(t *testing.T) {
	account1 := generateAccount(t)

	account2, err := testStore.GetAccount(context.Background(), account1.ID)
	require.NoError(t, err)
	require.NotEmpty(t, account2)

	require.Equal(t, account1.ID, account2.ID)
	require.Equal(t, account1.Owner, account2.Owner)
	require.Equal(t, account1.Balance, account2.Balance)
	require.Equal(t, account1.Currency, account2.Currency)
	require.Equal(t, account1.CreatedAt, account2.CreatedAt)
}

func TestListAccounts(t *testing.T) {
	var lastAccount Account
	generatedAccounts := make([]Account, 10)

	for i := 0; i < 10; i++ {
		lastAccount = generateAccount(t)
		generatedAccounts = append(generatedAccounts, lastAccount)
	}

	arg := ListAccountsParams{
		Owner:  lastAccount.Owner,
		Limit:  5,
		Offset: 0,
	}

	accounts, err := testStore.ListAccounts(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, accounts)

	for _, account := range accounts {
		require.NotEmpty(t, account)
		require.Equal(t, lastAccount.Owner, account.Owner)
	}
}

func TestDeleteAccount(t *testing.T) {
	account1 := generateAccount(t)

	err := testStore.DeleteAccount(context.Background(), account1.ID)
	require.NoError(t, err)

	account2, err := testStore.GetAccount(context.Background(), account1.ID)
	require.Error(t, err)
	require.EqualError(t, err, pgx.ErrNoRows.Error())
	require.Empty(t, account2)
}

func TestPurgeUserAccounts(t *testing.T) {
	var lastAccount Account
	generatedAccounts := make([]Account, 10)
	for i := 0; i < 10; i++ {
		lastAccount = generateAccount(t)
		generatedAccounts = append(generatedAccounts, lastAccount)
	}

	err := testStore.PurgeUserAccounts(context.Background(), lastAccount.Owner)
	require.NoError(t, err)

	listAccountsAgr := ListAccountsParams{
		Owner:  lastAccount.Owner,
		Limit:  100,
		Offset: 0,
	}

	accounts, _ := testStore.ListAccounts(context.Background(), listAccountsAgr)
	require.Empty(t, accounts)
}

func TestUpdateBalance(t *testing.T) {
	account1 := generateAccount(t)

	arg := UpdateBalanceParams{
		ID:     account1.ID,
		Amount: util.RandomFloat(100, 10000),
	}

	account2, err := testStore.UpdateBalance(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, account2)

	require.Equal(t, account1.ID, account2.ID)
	require.Equal(t, account1.Owner, account2.Owner)
	require.Equal(t, account1.Balance+arg.Amount, account2.Balance)
	require.Equal(t, account1.Currency, account2.Currency)

	require.WithinDuration(t, account1.CreatedAt, account2.CreatedAt, time.Second)
}
