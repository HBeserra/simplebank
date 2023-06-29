package db

import (
	"context"
	"database/sql"
	"fmt"
)

// Store provides all functions to execute db queries and transactions
type Store struct {
	*Queries
	db *sql.DB
}

// NewStore return a new Store
func NewStore(db *sql.DB) *Store {
	return &Store{
		db:      db,
		Queries: New(db),
	}
}

//execTransaction executes a function within a database transaction
func (s *Store) execTransaction(ctx context.Context, fn func(*Queries) error) error {
	transaction, err := s.db.BeginTx(ctx, nil)

	if err != nil {
		return err
	}

	q := New(transaction)
	err = fn(q)

	if err != nil {
		if rbErr := transaction.Rollback(); rbErr != nil {
			return fmt.Errorf("tx err: %v, rb err: %v", err, rbErr)
		}
		return err
	}

	return transaction.Commit()
}

// TransferTxParams Contains the input parameters of the transaction
type TransferTxParams struct {
	FromAccountID int64 `json:"from_account_id"`
	ToAccountID   int64 `json:"to_account_id"`
	Amount        int64 `json:"amount"`
}

// TransferTxResult is the result of the transaction
type TransferTxResult struct {
	Transfer    Transfer `json:"transfer"`
	FromAccount Account  `json:"from_account"`
	ToAccount   Account  `json:"to_account"`
	FromEntry   Entry    `json:"from_entry"`
	ToEntry     Entry    `json:"to_entry"`
}

// TransferTx performs a money transfer from one account to the other
// It creates a transfer record, add account entries and update accounts balance within a single db transaction
func (s Store) TransferTx(ctx context.Context, arg TransferTxParams) (TransferTxResult, error) {
	var result TransferTxResult

	// Call the transaction execution func and pass the context and a callback func.
	// The callback func execute the querys in the transaction
	err := s.execTransaction(ctx, func(q *Queries) error {
		var err error

		//Create the transfer record
		result.Transfer, err = q.CreateTransfer(ctx, CreateTransferParams{
			FromAccountID: arg.FromAccountID,
			ToAccountID:   arg.ToAccountID,
			Amount:        arg.Amount,
		})

		if err != nil {
			return err
		}

		//Create from account entry
		result.FromEntry, err = q.CreateEntry(ctx, CreateEntryParams{
			AccountID: arg.FromAccountID,
			Amount:    arg.Amount,
		})

		if err != nil {
			return err
		}

		//Create to account entry
		result.ToEntry, err = q.CreateEntry(ctx, CreateEntryParams{
			AccountID: arg.ToAccountID,
			Amount:    arg.Amount,
		})

		if err != nil {
			return err
		}

		// ToDo: Update accounts balance and use lock

		//Retrieve from account result
		result.FromAccount, err = q.GetAccount(ctx, arg.FromAccountID)

		if err != nil {
			return err
		}

		//Retrieve to account result
		result.ToAccount, err = q.GetAccount(ctx, arg.ToAccountID)

		if err != nil {
			return err
		}

		return nil
	})

	return result, err
}
