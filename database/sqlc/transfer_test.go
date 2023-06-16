package sqlc

import (
	"context"
	"testing"
	"time"

	"github.com/October-9th/simple-bank/util"
	"github.com/stretchr/testify/require"
)

func CreateRandomTransfer(t *testing.T, sender, receiver Account) Transfer {
	amount := util.RandomMoney()

	transfer, err := testQueries.CreateTransfer(context.Background(), CreateTransferParams{
		FromAccountID: sender.ID,
		ToAccountID:   receiver.ID,
		Amount:        amount,
	})

	require.NoError(t, err)
	require.NotEmpty(t, transfer)

	require.Equal(t, transfer.FromAccountID, sender.ID)
	require.Equal(t, transfer.ToAccountID, receiver.ID)
	require.Equal(t, transfer.Amount, amount)

	return transfer
}
func TestCreateTransfer(t *testing.T) {
	CreateRandomTransfer(t, CreateRandomAccount(t), CreateRandomAccount(t))
}
func TestGetTransfer(t *testing.T) {
	sender := CreateRandomAccount(t)
	receiver := CreateRandomAccount(t)
	test_transfer := CreateRandomTransfer(t, sender, receiver)
	transfer, err := testQueries.GetTransfer(context.Background(), test_transfer.ID)

	require.NoError(t, err)
	require.NotEmpty(t, transfer)

	require.Equal(t, transfer.ID, test_transfer.ID)
	require.Equal(t, transfer.FromAccountID, test_transfer.FromAccountID)
	require.Equal(t, transfer.ToAccountID, test_transfer.ToAccountID)
	require.Equal(t, transfer.Amount, test_transfer.Amount)

	require.WithinDuration(t, transfer.CreatedAt, test_transfer.CreatedAt, time.Second)
}
func TestListTransfer(t *testing.T) {
	sender := CreateRandomAccount(t)
	receiver := CreateRandomAccount(t)
	for i := 0; i < 10; i++ {
		CreateRandomTransfer(t, sender, receiver)
	}

	listTransfer, err := testQueries.ListTransfers(context.Background(), ListTransfersParams{
		FromAccountID: sender.ID,
		ToAccountID:   receiver.ID,
		Limit:         5,
		Offset:        5,
	})

	require.NoError(t, err)
	require.NotEmpty(t, listTransfer)
	require.Len(t, listTransfer, 5)
	for _, transfer := range listTransfer {
		require.NotEmpty(t, transfer)
		require.Equal(t, transfer.FromAccountID, sender.ID)
		require.Equal(t, transfer.ToAccountID, receiver.ID)
	}
}
