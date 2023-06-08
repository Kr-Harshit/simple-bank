package db

import (
	"context"
	"testing"

	util "github.com/KHarshit1203/simple-bank/util/datagenerator"
	"github.com/stretchr/testify/require"
)

func generateTransfer(ctx context.Context, t *testing.T, fromAccount, toAccount Account, amount float32) Transfer {
	arg := CreateTransferParams{
		FromAccountID: fromAccount.ID,
		ToAccountID:   toAccount.ID,
		Amount:        amount,
	}

	transfer, err := testQueries.CreateTransfer(ctx, arg)
	require.NoError(t, err)
	require.NotEmpty(t, transfer)

	require.NotEmpty(t, transfer.ID)
	require.Equal(t, fromAccount.ID, transfer.FromAccountID)
	require.Equal(t, toAccount.ID, transfer.ToAccountID)
	require.Equal(t, amount, transfer.Amount)
	require.NotEmpty(t, transfer.CreatedAt)

	return transfer
}

func TestCreateTransfer(t *testing.T) {
	ctx := context.Background()

	account1 := generateAccount(ctx, t)
	account2 := generateAccount(ctx, t)
	amount := util.RandomFloat(0, 1000)

	generateTransfer(ctx, t, account1, account2, amount)
}

func TestGetTransfer(t *testing.T) {
	ctx := context.Background()

	account1 := generateAccount(ctx, t)
	account2 := generateAccount(ctx, t)
	amount := util.RandomFloat(0, 1000)

	transfer1 := generateTransfer(ctx, t, account1, account2, float32(amount))

	transfer2, err := testQueries.GetTransfer(ctx, transfer1.ID)
	require.NoError(t, err)
	require.NotEmpty(t, transfer2)

	require.Equal(t, transfer1.FromAccountID, transfer2.FromAccountID)
	require.Equal(t, transfer1.ToAccountID, transfer2.ToAccountID)
	require.Equal(t, transfer1.Amount, transfer2.Amount)
	require.Equal(t, transfer1.CreatedAt, transfer2.CreatedAt)
}

func TestListTransfers(t *testing.T) {
	ctx := context.Background()

	account1 := generateAccount(ctx, t)
	account2 := generateAccount(ctx, t)

	for i := 0; i < 5; i++ {
		generateTransfer(ctx, t, account1, account2, util.RandomFloat(0, 100))
		generateTransfer(ctx, t, account2, account1, util.RandomFloat(0, 100))
	}

	arg := ListTransfersParams{
		FromAccountID: account1.ID,
		ToAccountID:   account1.ID,
		Limit:         10,
		Offset:        5,
	}

	transfers, err := testQueries.ListTransfers(ctx, arg)
	require.NoError(t, err)
	require.Len(t, transfers, 5)

	for _, transfer := range transfers {
		require.NotEmpty(t, transfer)
		require.True(t, transfer.FromAccountID == account1.ID || transfer.ToAccountID == account1.ID)
	}
}
