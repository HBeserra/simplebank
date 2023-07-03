package db

import (
	"context"
	"database/sql"
	"errors"
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
			Amount:    -arg.Amount,
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

		// Update accounts balance
		if arg.FromAccountID == arg.ToAccountID {
			return errors.New("transfer requires different sender and recipient accounts")
		}
		var err1, err2 error
		if arg.FromAccountID < arg.ToAccountID {
			result.FromAccount, err1 = addMoney(ctx, q, arg.FromAccountID, -arg.Amount)
			result.ToAccount, err2 = addMoney(ctx, q, arg.ToAccountID, arg.Amount)
		} else {
			result.ToAccount, err2 = addMoney(ctx, q, arg.ToAccountID, arg.Amount)
			result.FromAccount, err1 = addMoney(ctx, q, arg.FromAccountID, -arg.Amount)
		}

		if err1 != nil {
			return err1
		}
		if err2 != nil {
			return err2
		}

		return nil
	})

	return result, err
}

func addMoney(
	ctx context.Context,
	q *Queries,
	accountID int64,
	amount int64,
) (Account, error) {
	return q.AddAccountBalance(ctx, AddAccountBalanceParams{
		ID:     accountID,
		Amount: amount,
	})
}
