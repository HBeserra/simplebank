package db

import (
	"context"
	"github.com/HBeserra/simplebank/util"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func createRandomAccount(t *testing.T) Account {
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
	createRandomAccount(t)
}

func TestListAccounts(t *testing.T) {
	for i := 0; i < 10; i++ {
		createRandomAccount(t)
	}

	listAccountsArgs := ListAccountsParams{
		Limit:  5,
		Offset: 5,
	}

	accounts, err := testQueries.ListAccounts(context.Background(), listAccountsArgs)
	require.NoError(t, err)
	require.Len(t, accounts, 5)

	for _, account := range accounts {
		require.NotEmpty(t, account)
	}
}

func TestGetAccount(t *testing.T) {
	randomAccount := createRandomAccount(t)

	account, err := testQueries.GetAccount(context.Background(), randomAccount.ID)
	require.NoError(t, err)
	require.NotEmpty(t, account)

	require.Equal(t, randomAccount.ID, account.ID)
	require.Equal(t, randomAccount.Owner, account.Owner)
	require.Equal(t, randomAccount.Balance, account.Balance)
	require.Equal(t, randomAccount.Currency, account.Currency)
	require.WithinDuration(t, randomAccount.CreatedAt, account.CreatedAt, time.Second)
}

func TestUpdateAccount(t *testing.T) {
	randomAccount := createRandomAccount(t)

	updateArgs := UpdateAccountParams{
		ID:      randomAccount.ID,
		Balance: util.RandomMoney(),
	}

	updatedAccount, err := testQueries.UpdateAccount(context.Background(), updateArgs)
	require.NoError(t, err)
	require.NotEmpty(t, updatedAccount)

	require.Equal(t, randomAccount.ID, updatedAccount.ID)
	require.Equal(t, updateArgs.Balance, updatedAccount.Balance)

}

func TestDeleteAccount(t *testing.T) {
	randomAccount := createRandomAccount(t)

	err := testQueries.DeleteAccount(context.Background(), randomAccount.ID)
	require.NoError(t, err)

	account, err := testQueries.GetAccount(context.Background(), randomAccount.ID)
	require.ErrorContains(t, err, "sql: no rows in result set")
	require.Zero(t, account.ID)
	require.Empty(t, account)
}
