package sqlc

import (
	"context"
	"database/sql"
	"fmt"
)

// Store provides all function to execute database queries and transactions
type Store struct {
	*Queries
	db *sql.DB // Required for creating database transactions
}

// NewStore creates a new store
func NewStore(db *sql.DB) *Store {
	return &Store{
		db:      db,
		Queries: New(db),
	}
}

// This function take the context and a call back function as input,
// then it will start a new db transaction -> create a new query object with that transaction
// and call the callback function with the created query -> finally commit or rollback the transaction
// based on the error returnd by the callback function
func (s *Store) execTx(ctx context.Context, fn func(*Queries) error) error {
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	// Call New function with the created transaction and get back a new query object
	q := New(tx)
	// When we have the query within transaction we can now call the input function with that query
	err = fn(q)
	// if the error is not nil, then we have to roll back the transaction
	if err != nil {
		if rbErr := tx.Rollback(); rbErr != nil {
			// if the rollback error is not nil
			// then we should combine them into 1 single error before returning
			return fmt.Errorf("tx error: %v, roll back error: %v", err, rbErr)
		}
		return err
	}
	// finally if all the operation in the transaction is successful
	// commit it
	return tx.Commit()
}

// TransferTxParams contains the input parameters of the transfer transaction
type TransferTxParams struct {
	FromAccountID int64 `json:"from_account_id"`
	ToAccountID   int64 `json:"to_account_id"`
	Amount        int64 `json:"amount"`
}

// TransferTxResult is the result of the transfer transaction
type TransferTxResult struct {
	Transfer    Transfer `json:"transfer"`     // the created  transfer record
	FromAccount Account  `json:"from_account"` // from account after it balance has been updated
	ToAccount   Account  `json:"to_account"`   // the same as from account
	FromEntry   Entry    `json:"from_entry"`   // the entry of the from account which record that money is moving out
	ToEntry     Entry    `json:"to_entry"`     // the entry of the to account which record that money is moving in
}

// TransferTx performs a money transfer from one account to another
// It creates a transfer record, add account entries and update account's balance within a single database transaction
func (s *Store) TransferTx(ctx context.Context, arg TransferTxParams) (TransferTxResult, error) {
	var txResult TransferTxResult

	err := s.execTx(ctx, func(q *Queries) error {
		var err error

		// First step create a transfer record
		createTransferParams := CreateTransferParams(arg)

		txResult.Transfer, err = q.CreateTransfer(ctx, createTransferParams)
		if err != nil {
			return err
		}

		// Second step add account entry
		txResult.FromEntry, err = q.CreateEntry(ctx, CreateEntryParams{
			AccountID: arg.FromAccountID,
			Amount:    -arg.Amount,
		})
		if err != nil {
			return err
		}

		txResult.ToEntry, err = q.CreateEntry(ctx, CreateEntryParams{
			AccountID: arg.ToAccountID,
			Amount:    arg.Amount,
		})
		if err != nil {
			return err
		}

		// Last step updating account balance, invoking locking and preventing deadlock
		// Step: get account -> update its balance

		if arg.FromAccountID < arg.ToAccountID {
			txResult.FromAccount, txResult.ToAccount, err = UpdateAccountBalance(arg.FromAccountID, arg.ToAccountID, -arg.Amount, arg.Amount, ctx, q)
			if err != nil {
				return err
			}
		} else {
			txResult.ToAccount, txResult.FromAccount, err = UpdateAccountBalance(arg.ToAccountID, arg.FromAccountID, arg.Amount, -arg.Amount, ctx, q)
			if err != nil {
				return err
			}
		}
		return nil
	})
	return txResult, err

}
func UpdateAccountBalance(
	accountID_1,
	accountID_2 int64,
	amount1,
	amount2 int64,
	ctx context.Context,
	q *Queries,
) (account1, account2 Account, err error) {
	account1, err = q.AddAccountBalance(ctx, AddAccountBalanceParams{
		ID:     accountID_1,
		Amount: amount1,
	})
	if err != nil {
		return
	}
	account2, err = q.AddAccountBalance(ctx, AddAccountBalanceParams{
		ID:     accountID_2,
		Amount: amount2,
	})
	if err != nil {
		return
	}
	return

}
