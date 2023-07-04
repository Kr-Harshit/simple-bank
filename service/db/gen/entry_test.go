package db

import (
	"context"
	"testing"

	"github.com/KHarshit1203/simple-bank/service/util"
	"github.com/stretchr/testify/require"
)

func TestCreateEntry(t *testing.T) {
	account1 := generateAccount(t)
	account2 := generateAccount(t)
	amount := util.RandomFloat(0, 100)
	transfer := generateTransfer(t, account1, account2, amount)

	generateEntries(t, account1, transfer, false)
}

func TestGetEntry(t *testing.T) {
	account1 := generateAccount(t)
	account2 := generateAccount(t)
	amount := util.RandomFloat(0, 100)

	transfer1 := generateTransfer(t, account1, account2, amount)
	entry1 := generateEntries(t, account1, transfer1, false)

	entry2, err := testStore.GetEntry(context.Background(), entry1.ID)
	require.NoError(t, err)

	require.NotEmpty(t, entry2)
	require.Equal(t, account1.ID, entry2.AccountID)
	require.NotEmpty(t, entry2.CreatedAt)
	require.NotEmpty(t, entry2.TransferID)
	require.Equal(t, false, entry2.Credit)
	require.Equal(t, amount, entry2.Amount)
}

func TestListEntries(t *testing.T) {
	account1 := generateAccount(t)
	account2 := generateAccount(t)

	for i := 0; i < 10; i++ {
		transfer1 := generateTransfer(t, account1, account2, util.RandomFloat(0, 100))
		generateEntries(t, account1, transfer1, false)
		transfer2 := generateTransfer(t, account2, account1, util.RandomFloat(0, 100))
		generateEntries(t, account1, transfer2, true)
	}

	arg := ListEntriesParams{
		AccountID: account1.ID,
		Limit:     20,
		Offset:    0,
	}

	entries, err := testStore.ListEntries(context.Background(), arg)
	require.NoError(t, err)
	require.Len(t, entries, 20)

	for _, entry := range entries {
		require.NotEmpty(t, entry)
		require.Equal(t, entry.AccountID, account1.ID)
		require.NotEmpty(t, entry.CreatedAt)
		require.NotEmpty(t, entry.TransferID)
	}

}

func generateEntries(t *testing.T, account Account, transfer Transfer, credit bool) Entry {
	amount := transfer.Amount

	arg1 := CreateEntryParams{
		AccountID:  account.ID,
		Amount:     amount,
		TransferID: transfer.ID,
		Credit:     credit,
	}

	entry, err := testStore.CreateEntry(context.Background(), arg1)
	require.NoError(t, err)
	require.NotEmpty(t, entry)

	require.NotEmpty(t, entry.ID)
	require.Equal(t, account.ID, entry.AccountID)
	require.Equal(t, amount, entry.Amount)
	require.Equal(t, transfer.ID, entry.TransferID)
	require.NotEmpty(t, entry.CreatedAt)
	require.Equal(t, credit, entry.Credit)

	return entry
}
