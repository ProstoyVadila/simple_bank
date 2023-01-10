package db

import (
	"context"
	"database/sql"
	"testing"
	"time"

	"github.com/ProstoyVadila/simple_bank/utils"
	"github.com/stretchr/testify/require"
)

func cleanUpAccount(t *testing.T, account Account) {
	time.Sleep(time.Second)
	t.Cleanup(func() {
		err := testQueries.DeleteAccount(context.Background(), account.ID)
		require.NoError(t, err)
	})
	cleanUpUser(t, account.OwnerName)
}

func createRandomAccount(t *testing.T) Account {
	user := createRandomUser(t)
	arg := CreateAccountParams{
		OwnerName: user.Username,
		Balance:   utils.RandomBalance(),
		Currency:  "KZT",
	}
	account, err := testQueries.CreateAccount(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, account)

	require.Equal(t, arg.OwnerName, account.OwnerName)
	require.Equal(t, arg.Currency, account.Currency)
	require.Equal(t, arg.Balance, account.Balance)

	require.NotZero(t, account.ID)
	require.NotZero(t, account.CreatedAt)
	require.Equal(t, "VERSION_4", account.ID.Version().String())

	return account
}

func TestCreateAccount(t *testing.T) {
	account := createRandomAccount(t)
	cleanUpAccount(t, account)
}

func TestGetAccount(t *testing.T) {
	account1 := createRandomAccount(t)
	account2, err := testQueries.GetAccount(context.Background(), account1.ID)
	require.NoError(t, err)
	require.NotEmpty(t, account2)
	require.WithinDuration(t, account1.CreatedAt, account2.CreatedAt, time.Second)
	require.Equal(t, account1, account2)
	cleanUpAccount(t, account1)
}

func TestUpdateAccount(t *testing.T) {
	account1 := createRandomAccount(t)
	args := UpdateAccountParams{
		ID:      account1.ID,
		Balance: utils.RandomBalance(),
	}
	account2, err := testQueries.UpdateAccount(context.Background(), args)
	require.NoError(t, err)
	require.NotEmpty(t, account2)

	require.Equal(t, account1.ID, account2.ID)
	require.Equal(t, account1.OwnerName, account2.OwnerName)
	require.Equal(t, account1.Currency, account2.Currency)

	require.Equal(t, args.Balance, account2.Balance)
	require.NotEqual(t, account1.Balance, account2.Balance)
	cleanUpAccount(t, account1)
}

func TestDeleteAccount(t *testing.T) {
	account1 := createRandomAccount(t)
	err := testQueries.DeleteAccount(context.Background(), account1.ID)
	require.NoError(t, err)

	account2, err := testQueries.GetAccount(context.Background(), account1.ID)
	require.Error(t, err)
	require.EqualError(t, err, sql.ErrNoRows.Error())
	require.Empty(t, account2)
}

func TestListAccounts(t *testing.T) {
	n := 10
	accounts1 := make([]Account, n)
	for i := 0; i < n; i++ {
		accounts1[i] = createRandomAccount(t)
	}

	args := ListAccountsParams{
		Limit:  5,
		Offset: 5,
	}

	accounts2, err := testQueries.ListAccounts(context.Background(), args)
	require.NoError(t, err)
	require.Len(t, accounts2, 5)

	for _, acc := range accounts2 {
		require.NotEmpty(t, acc)
	}
	for _, account := range accounts1 {
		cleanUpAccount(t, account)
	}

}
