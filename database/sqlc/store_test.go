package sqlc

import (
	"context"
	"log"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestTransferTx(t *testing.T) {
	store := NewStore(testDB)

	// Create account for transaction
	account_1, account_2 := CreateRandomAccount(t), CreateRandomAccount(t)
	log.Println(">> Before: ", account_1.Balance, account_2.Balance)
	// Run n concurrent transfer transactions
	n := 10
	amount := int64(10)

	// Declare a channel to receive the TransferTx error
	errors := make(chan error)

	// Declare another channel to receive the TransferTx result
	results := make(chan TransferTxResult)
	existed := make(map[int]bool)
	for i := 0; i < n; i++ {
		// Since this function is running inside a different goroutine from the TestTransferTx is running
		// So there will be no gurantee that it will stop the whole test if a condition is not satisfied
		// So we will have to send them back to the main go routine

		//  use log to print which transaction is calling which query in which order

		// asign a name for the transaction and pass it into the TransferTx function
		go func() {
			ctx := context.Background()
			result, err := store.TransferTx(ctx, TransferTxParams{
				FromAccountID: account_1.ID,
				ToAccountID:   account_2.ID,
				Amount:        amount,
			})
			errors <- err
			results <- result
		}()
	}
	// Check the result
	for i := 0; i < n; i++ {
		err := <-errors
		require.NoError(t, err)

		result := <-results
		require.NotEmpty(t, result)

		// Check transfer
		transfer := result.Transfer
		require.Equal(t, account_1.ID, transfer.FromAccountID)
		require.Equal(t, account_2.ID, transfer.ToAccountID)
		require.Equal(t, amount, transfer.Amount)
		require.NotZero(t, transfer.ID)
		require.NotZero(t, transfer.CreatedAt)
		_, err = store.GetTransfer(context.Background(), transfer.ID)
		require.NoError(t, err)
		// check entries
		fromEntry := result.FromEntry

		require.NotEmpty(t, fromEntry)
		require.Equal(t, account_1.ID, fromEntry.AccountID)
		require.Equal(t, -amount, fromEntry.Amount)
		require.NotZero(t, fromEntry.ID)
		require.NotZero(t, fromEntry.CreatedAt)
		_, err = store.GetEntry(context.Background(), fromEntry.ID)
		require.NoError(t, err)

		toEntry := result.ToEntry

		require.NotEmpty(t, toEntry)
		require.Equal(t, account_2.ID, toEntry.AccountID)
		require.Equal(t, amount, toEntry.Amount)
		require.NotZero(t, toEntry.ID)
		require.NotZero(t, toEntry.CreatedAt)
		_, err = store.GetEntry(context.Background(), toEntry.ID)
		require.NoError(t, err)

		// Check account
		fromAccount := result.FromAccount
		require.NotEmpty(t, fromAccount)
		require.Equal(t, account_1.ID, fromAccount.ID)

		toAccount := result.ToAccount
		require.NotEmpty(t, toAccount)
		require.Equal(t, account_2.ID, toAccount.ID)

		// check account'balance
		log.Println(">>tx: ", fromAccount.Balance, toAccount.Balance)
		diff_1 := account_1.Balance - fromAccount.Balance
		diff_2 := toAccount.Balance - account_2.Balance
		require.Equal(t, diff_1, diff_2)
		require.True(t, diff_1 > 0)
		require.True(t, diff_1%amount == 0) // amount, 2 * amount, 3 * amount, ..., n * amount

		k := int(diff_1 / amount)
		require.True(t, k >= 1 && k <= n)
		require.NotContains(t, existed, k)
		existed[k] = true
	}

	// Check the final updated balance

	updatedAccount_1, err := testQueries.GetAccount(context.Background(), account_1.ID)
	require.NoError(t, err)
	updatedAccount_2, err := testQueries.GetAccount(context.Background(), account_2.ID)
	require.NoError(t, err)
	log.Println(">>tx: ", updatedAccount_1.Balance, updatedAccount_2.Balance)
	require.Equal(t, account_1.Balance-int64(n)*amount, updatedAccount_1.Balance)
	require.Equal(t, account_2.Balance+int64(n)*amount, updatedAccount_2.Balance)
}

func TestTransferTxDeadlock(t *testing.T) {
	store := NewStore(testDB)

	// Create account for transaction
	account_1, account_2 := CreateRandomAccount(t), CreateRandomAccount(t)
	log.Println(">> Before: ", account_1.Balance, account_2.Balance)
	// Run n concurrent transfer transactions
	n := 10
	amount := int64(10)

	// Declare a channel to receive the TransferTx error
	errors := make(chan error)

	// Run n concurrent transfer transactions
	for i := 0; i < n; i++ {
		fromAccountId := account_1.ID
		toAccountId := account_2.ID

		// Reverse the order of the transfer transactions
		if i%2 == 1 {
			fromAccountId = account_2.ID
			toAccountId = account_1.ID
		}

		go func() {
			_, err := store.TransferTx(context.Background(), TransferTxParams{
				FromAccountID: fromAccountId,
				ToAccountID:   toAccountId,
				Amount:        amount,
			})

			errors <- err
		}()
	}
	// Check the result
	for i := 0; i < n; i++ {
		err := <-errors
		require.NoError(t, err)
	}

	// Check the final updated balance
	updatedAccount_1, err := store.GetAccount(context.Background(), account_1.ID)
	require.NoError(t, err)
	updatedAccount_2, err := store.GetAccount(context.Background(), account_2.ID)
	require.NoError(t, err)

	log.Println(">>tx: ", updatedAccount_1.Balance, updatedAccount_2.Balance)

	require.Equal(t, account_1.Balance, updatedAccount_1.Balance)
	require.Equal(t, account_2.Balance, updatedAccount_2.Balance)
}
