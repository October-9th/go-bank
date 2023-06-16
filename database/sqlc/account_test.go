package sqlc

import (
	"context"
	"database/sql"
	"testing"
	"time"

	"github.com/October-9th/simple-bank/util"
	"github.com/stretchr/testify/require"
)

func CreateRandomAccount(t *testing.T) Account {
	arg := CreateAccountParams{
		Owner:    util.RandomOwner(),
		Balance:  util.RandomMoney(),
		Currency: util.RandomCurrency(),
	}
	account, err := testQueries.CreateAccount(context.Background(), arg)
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
	CreateRandomAccount(t)
}

func TestGetAccount(t *testing.T) {
	account_test := CreateRandomAccount(t)
	account, err := testQueries.GetAccount(context.Background(), account_test.ID)

	require.NoError(t, err)
	require.NotEmpty(t, account)

	require.Equal(t, account.Owner, account_test.Owner)
	require.Equal(t, account.Balance, account_test.Balance)
	require.Equal(t, account.Currency, account_test.Currency)
	// Use require.WithinDuration to check 2 timestamps are different by at most some delta duration
	require.WithinDuration(t, account.CreatedAt, account_test.CreatedAt, time.Second)
}
func TestListAccounts(t *testing.T) {
	for i := 0; i < 10; i++ {
		CreateRandomAccount(t)
	}
	listAccounts, err := testQueries.ListAccounts(context.Background(), ListAccountsParams{
		Limit:  5,
		Offset: 5,
	})
	require.NoError(t, err)
	require.Len(t, listAccounts, 5)

	for _, account := range listAccounts {
		require.NotEmpty(t, account)
	}
}

func TestDeleteAccount(t *testing.T) {
	account_test := CreateRandomAccount(t)
	err := testQueries.DeleteAccount(context.Background(), account_test.ID)
	require.NoError(t, err)

	account_after_deleted, err := testQueries.GetAccount(context.Background(), account_test.ID)
	require.Error(t, err)
	require.Error(t, err, sql.ErrNoRows.Error())
	require.Empty(t, account_after_deleted)
}

func TestUpdateAccount(t *testing.T) {
	account_test := CreateRandomAccount(t)
	arg := UpdateAccountParams{
		ID:      account_test.ID,
		Balance: util.RandomMoney(),
	}
	updatedAccount, err := testQueries.UpdateAccount(context.Background(), arg)
	require.NoError(t, err)
	require.Equal(t, updatedAccount.ID, account_test.ID)
	require.Equal(t, updatedAccount.Owner, account_test.Owner)
	require.Equal(t, arg.Balance, updatedAccount.Balance)
	require.Equal(t, updatedAccount.Currency, account_test.Currency)
	require.NotEmpty(t, updatedAccount)

	require.WithinDuration(t, updatedAccount.CreatedAt, account_test.CreatedAt, time.Second)
}
