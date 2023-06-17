package db

import (
	"context"
	"testing"

	util "github.com/KHarshit1203/simple-bank/util/datagenerator"
	"github.com/stretchr/testify/require"
)

func TestTransferTx(t *testing.T) {
	ctx := context.Background()

	fromAccount := generateAccount(ctx, t)
	toAccount := generateAccount(ctx, t)
	amount := util.RandomFloat(1, 10)
	n := 10
	errs := make(chan error)
	results := make(chan TransferTxResult)

	for i := 0; i < n; i++ {
		// t.Logf("transfering %v from %v account to %v account", amount, fromAccount.ID, toAccount.ID)
		t.Logf("Amount before transfering: fromAccount: %v  toAccount: %v", fromAccount.Balance, toAccount.Balance)

		go func() {
			result, err := testStore.TransferTx(context.Background(), TransferTxParams{
				FromAccountID: fromAccount.ID,
				ToAccountID:   toAccount.ID,
				Amount:        amount,
			})
			// t.Errorf("transfer Tx error: %v", err)
			errs <- err
			results <- result
		}()
	}

	for i := 0; i < n; i++ {
		err := <-errs
		require.NoError(t, err)

		result := <-results
		require.NotEmpty(t, result)

		checkTransferResult(t, fromAccount, toAccount, amount, n, result)
	}
}

func TestTransferTxDeadLock(t *testing.T) {
	ctx := context.Background()

	account1 := generateAccount(ctx, t)
	account2 := generateAccount(ctx, t)
	amount := util.RandomFloat(1, 10)
	n := 10
	errs := make(chan error)
	results := make(chan TransferTxResult)

	for i := 0; i < n; i++ {
		t.Logf("Amount before transfering: account1: [id:%d, amount: %v]  account2: [id: %d, ammount: %v]", account1.ID, account1.Balance, account2.ID, account2.Balance)

		fromAccountID := account1.ID
		toAccountID := account2.ID
		if i%2 == 0 {
			fromAccountID = account2.ID
			toAccountID = account1.ID
		}

		t.Logf("transfering %v from %v account to %v account", amount, fromAccountID, toAccountID)

		go func() {
			result, err := testStore.TransferTx(context.Background(), TransferTxParams{
				FromAccountID: fromAccountID,
				ToAccountID:   toAccountID,
				Amount:        amount,
			})
			// t.Errorf("transfer Tx error: %v", err)
			errs <- err
			results <- result
		}()
	}

	for i := 0; i < n; i++ {
		err := <-errs
		require.NoError(t, err)
	}

	// check the final updated balance
	updatedAccount1, err := testStore.GetAccount(context.Background(), account1.ID)
	require.NoError(t, err)

	updatedAccount2, err := testStore.GetAccount(context.Background(), account2.ID)
	require.NoError(t, err)

	t.Logf(">> after tx baklance, account1: %v, account2: %v", updatedAccount1.Balance, updatedAccount2.Balance)
	require.Equal(t, account1.Balance, updatedAccount1.Balance)
	require.Equal(t, account2.Balance, updatedAccount2.Balance)
}

func checkTransferResult(t *testing.T, fromAccount Account, toAccount Account, amount float32, iteration int, result TransferTxResult) {
	ctx := context.Background()
	var err error

	// check transfer
	require.NotEmpty(t, result.Transfer)
	require.NotZero(t, result.Transfer.ID)
	require.NotZero(t, result.Transfer.CreatedAt)
	require.Equal(t, fromAccount.ID, result.Transfer.FromAccountID)
	require.Equal(t, toAccount.ID, result.Transfer.ToAccountID)
	require.Equal(t, amount, result.Transfer.Amount)

	_, err = testStore.GetTransfer(ctx, result.Transfer.ID)
	require.NoError(t, err)

	// // check from entry
	require.NotEmpty(t, result.FromEntry)
	require.Equal(t, fromAccount.ID, result.FromEntry.AccountID)
	require.Equal(t, result.Transfer.ID, result.FromEntry.TransferID)
	require.Equal(t, amount, result.FromEntry.Amount)
	require.Equal(t, false, result.FromEntry.Credit)
	require.NotZero(t, result.FromEntry.ID)
	require.NotZero(t, result.FromEntry.CreatedAt)

	_, err = testStore.GetEntry(ctx, result.FromEntry.ID)
	require.NoError(t, err)

	// check to entry
	require.NotEmpty(t, result.ToEntry)
	require.Equal(t, toAccount.ID, result.ToEntry.AccountID)
	require.Equal(t, result.Transfer.ID, result.ToEntry.TransferID)
	require.Equal(t, amount, result.ToEntry.Amount)
	require.Equal(t, true, result.ToEntry.Credit)
	require.NotZero(t, result.ToEntry.ID)
	require.NotZero(t, result.ToEntry.CreatedAt)

	_, err = testStore.GetEntry(ctx, result.ToEntry.ID)
	require.NoError(t, err)

	// check accounts
	// t.Logf("result fromAccount: %v", result.FromAccount)
	fromAccount2 := result.FromAccount
	require.NotEmpty(t, fromAccount2)
	require.Equal(t, fromAccount.ID, fromAccount2.ID)

	// t.Logf("result toAccount: %v", result.ToAccount)
	toAccount2 := result.ToAccount
	require.NotEmpty(t, toAccount2)
	require.Equal(t, toAccount.ID, toAccount2.ID)
}
