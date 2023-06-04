package db

import (
	"context"
	"fmt"
	"testing"

	util "github.com/Kr-Harshit/learn/simple-bank/util/datagenerator"
	"github.com/stretchr/testify/require"
)

func TestTransferTx(t *testing.T) {
	ctx := context.Background()
	testStore := NewStore(testDb)

	account1 := generateAccount(ctx, t)
	account2 := generateAccount(ctx, t)
	t.Logf("Ammount before transfering: account1: %v  account2: %v", account1.Balance, account2.Balance)

	amount := util.RandomFloat(10, 1000)
	n := 10
	errs := make(chan error)
	results := make(chan TransferTxResult)

	for i := 0; i < n; i++ {
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
			t.Errorf("trnasfer Tx error: %+v", err)

			errs <- err
			results <- result
		}()
	}

	for i := 0; i < n; i++ {
		err := <-errs
		require.NoError(t, err)

		result := <-results
		require.NotEmpty(t, result)

		// check transfer
		require.NotEmpty(t, result.Transfer)
		require.Equal(t, account1.ID, result.Transfer.FromAccountID)
		require.Equal(t, account2.ID, result.Transfer.ToAccountID)
		require.NotZero(t, result.Transfer.ID)
		require.NotZero(t, result.Transfer.CreatedAt)
		require.Equal(t, amount, result.Transfer.Amount)

		_, err = testStore.GetTransfer(ctx, result.Transfer.ID)
		require.NoError(t, err)

		// check from entry
		require.NotEmpty(t, result.FromEntry)
		require.Equal(t, account1.ID, result.FromEntry.AccountID)
		require.Equal(t, result.Transfer.ID, result.FromEntry.TransferID)
		require.Equal(t, amount, result.FromEntry.Amount)
		require.Equal(t, false, result.FromEntry.Credit)
		require.NotZero(t, result.FromEntry.ID)
		require.NotZero(t, result.FromEntry.CreatedAt)

		_, err = testStore.GetEntry(ctx, result.FromEntry.ID)
		require.NoError(t, err)

		// check to entry
		require.NotEmpty(t, result.ToEntry)
		require.Equal(t, account2.ID, result.ToEntry.AccountID)
		require.Equal(t, result.Transfer.ID, result.ToEntry.TransferID)
		require.Equal(t, amount, result.ToEntry.Amount)
		require.Equal(t, true, result.ToEntry.Credit)
		require.NotZero(t, result.ToEntry.ID)
		require.NotZero(t, result.ToEntry.CreatedAt)

		_, err = testStore.GetEntry(ctx, result.ToEntry.ID)
		require.NoError(t, err)

		// check accounts
		fromAccount := result.FromAccount
		require.NotEmpty(t, fromAccount)
		require.Equal(t, account1.ID, fromAccount.ID)

		toAccount := result.ToAccount
		require.NotEmpty(t, toAccount)
		require.Equal(t, account2.ID, toAccount.ID)

		// check balance
		fmt.Println(">> tx:", fromAccount.Balance, toAccount.Balance)

		diff1 := account1.Balance - fromAccount.Balance
		diff2 := toAccount.Balance - account2.Balance

		require.Equal(t, diff1, diff2)
		require.True(t, diff1 > 0)

		k := int(diff1 / amount)
		require.True(t, k >= 1 && k <= n)
	}
}
