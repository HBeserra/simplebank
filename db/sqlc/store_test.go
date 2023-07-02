package db

import (
	"context"
	"fmt"
	"github.com/HBeserra/simplebank/util"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestTransferTx(t *testing.T) {
	// run n concurrent transfer transactions
	const n = 50

	store := NewStore(testDB)

	account1 := createRandomAccount(t)
	account2 := createRandomAccount(t)

	fmt.Printf("Accounts balance: \t%d - %d\n", account1.Balance, account2.Balance)

	amount := util.RandomMoney()

	errs := make(chan error)               // Chanel for go routines return erros
	results := make(chan TransferTxResult) // Chanel for go routines return transaction result

	println("starting test")
	for i := 0; i < n; i++ {
		txName := fmt.Sprintf("tx %d", i+1)
		println(txName, "started")
		go func() {
			result, err := store.TransferTx(context.Background(), TransferTxParams{
				FromAccountID: account1.ID,
				ToAccountID:   account2.ID,
				Amount:        amount,
			})

			errs <- err
			results <- result
		}()
	}

	fmt.Println("Transfer results:")
	// Check results
	existed := make(map[int]bool)

	for i := 0; i < n; i++ {
		err := <-errs
		require.NoError(t, err)

		result := <-results
		require.NotEmpty(t, result)

		// check transfer
		transfer := result.Transfer
		require.NotEmpty(t, transfer)
		require.Equal(t, transfer.FromAccountID, account1.ID)
		require.Equal(t, transfer.ToAccountID, account2.ID)
		require.Equal(t, transfer.Amount, amount)
		require.NotZero(t, transfer.ID)
		require.NotZero(t, transfer.CreatedAt)

		_, err = store.GetTransfer(context.Background(), transfer.ID)
		require.NoError(t, err)

		//Check Entry
		fromEntry := result.FromEntry
		require.NotEmpty(t, fromEntry)
		require.NotZero(t, fromEntry.ID)
		require.Equal(t, fromEntry.AccountID, account1.ID)
		require.Equal(t, -amount, fromEntry.Amount)
		require.NotZero(t, fromEntry.CreatedAt)

		_, err = store.GetEntry(context.Background(), fromEntry.ID)
		require.NoError(t, err)

		//Check to entry
		toEntry := result.ToEntry
		require.NotEmpty(t, toEntry)
		require.NotZero(t, toEntry.ID)
		require.Equal(t, toEntry.AccountID, account2.ID)
		require.Equal(t, toEntry.Amount, amount)
		require.NotZero(t, toEntry.CreatedAt)

		_, err = store.GetEntry(context.Background(), toEntry.ID)
		require.NoError(t, err)

		// Check fromAccount balance
		fromAccount := result.FromAccount
		require.NotEmpty(t, fromAccount)
		require.NotZero(t, fromAccount.ID)
		require.Equal(t, fromAccount.ID, account1.ID)

		// Check fromAccount balance
		toAccount := result.ToAccount
		require.NotEmpty(t, toAccount)
		require.NotZero(t, toAccount.ID)
		require.Equal(t, toAccount.ID, account2.ID)

		diff1 := account1.Balance - fromAccount.Balance
		diff2 := toAccount.Balance - account2.Balance

		require.Equal(t, diff1, diff2)
		require.True(t, diff1 > 0)
		require.True(t, diff1%amount == 0)

		k := int(diff1 / amount)
		require.True(t, k >= 1 && k <= n)
		require.NotContains(t, existed, k)
		existed[k] = true

		//fmt.Printf("\ttransaction [%d]: %+v\n", i+1, result)
		fmt.Printf("\t [%d] %d: %d | %d\n", i+1, amount, result.FromAccount.Balance, result.ToAccount.Balance)
	}
}
