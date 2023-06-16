package sqlc

import (
	"context"
	"testing"
	"time"

	"github.com/October-9th/simple-bank/util"
	"github.com/stretchr/testify/require"
)

func CreateRandomEntry(t *testing.T, account Account) Entry {
	arg := CreateEntryParams{
		AccountID: account.ID,
		Amount:    util.RandomMoney(),
	}
	entry, err := testQueries.CreateEntry(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, entry)
	require.Equal(t, entry.AccountID, account.ID)
	require.Equal(t, entry.Amount, arg.Amount)
	require.NotZero(t, entry.ID)
	require.NotZero(t, entry.CreatedAt)
	return entry
}
func TestCreateEntry(t *testing.T) {
	account := CreateRandomAccount(t)
	CreateRandomEntry(t, account)
}
func TestGetEntry(t *testing.T) {
	account := CreateRandomAccount(t)
	test_entry := CreateRandomEntry(t, account)
	entry, err := testQueries.GetEntry(context.Background(), test_entry.ID)

	require.NoError(t, err)
	require.NotEmpty(t, entry)
	require.Equal(t, entry.ID, test_entry.ID)
	require.Equal(t, entry.AccountID, test_entry.AccountID)
	require.Equal(t, entry.Amount, test_entry.Amount)
	require.WithinDuration(t, entry.CreatedAt, test_entry.CreatedAt, time.Second)
}
func TestGetListEntry(t *testing.T) {
	account := CreateRandomAccount(t)
	for i := 0; i < 10; i++ {
		CreateRandomEntry(t, account)
	}

	entries, err := testQueries.ListEntries(context.Background(), ListEntriesParams{
		AccountID: account.ID,
		Limit:     5,
		Offset:    5,
	})
	require.NoError(t, err)
	require.NotEmpty(t, entries)

	for _, entry := range entries {
		require.NotEmpty(t, entry)
		require.Equal(t, entry.AccountID, account.ID)
	}
}
