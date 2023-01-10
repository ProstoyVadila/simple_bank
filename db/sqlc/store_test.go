package db

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestTransferTx(t *testing.T) {
	store := NewStore(testDB)
	account1 := createRandomAccount(t)
	account2 := createRandomAccount(t)

	n := 5
	amount := int64(10)
	transfers := make([]TransferTxResult, n)

	errs := make(chan error)
	results := make(chan TransferTxResult)
	defer close(errs)
	defer close(results)

	for i := 0; i < n; i++ {
		go func() {
			res, err := store.TransferTx(
				context.Background(),
				TransferTxParams{
					FromAccountID: account1.ID,
					ToAccountID:   account2.ID,
					Amount:        amount,
					Currency:      account1.Currency,
				})
			errs <- err
			results <- res
		}()
	}

	existed := make(map[int]bool)
	for i := 0; i < n; i++ {
		err := <-errs
		require.NoError(t, err)

		result := <-results
		require.NotEmpty(t, result)
		transfers[i] = result

		// Checking transfer data
		require.NotEmpty(t, result.Transfer)
		require.Equal(t, account1.ID, result.Transfer.FromAccountID)
		require.Equal(t, account2.ID, result.Transfer.ToAccountID)
		require.Equal(t, amount, result.Transfer.Amount)
		require.NotEmpty(t, result.Transfer.ID)
		require.NotZero(t, result.Transfer.CreatedAt)

		_, err = store.GetTrasfer(context.Background(), result.Transfer.ID)
		require.NoError(t, err)

		// Checking entries data
		// FromEntry
		require.NotEmpty(t, result.FromEntry)
		require.NotZero(t, result.FromEntry.ID)
		require.NotZero(t, result.FromEntry.CreatedAt)
		require.Equal(t, account1.ID, result.FromEntry.AccountID)
		require.Equal(t, -amount, result.FromEntry.Amount)

		_, err = store.GetEntry(context.Background(), result.FromEntry.ID)
		require.NoError(t, err)

		// ToEntry
		require.NotEmpty(t, result.ToEntry)
		require.NotZero(t, result.ToEntry.ID)
		require.NotZero(t, result.ToEntry.CreatedAt)
		require.Equal(t, account2.ID, result.ToEntry.AccountID)
		require.Equal(t, amount, result.ToEntry.Amount)

		_, err = store.GetEntry(context.Background(), result.ToEntry.ID)
		require.NoError(t, err)

		// Accounts
		require.NotEmpty(t, result.FromAccount)
		require.Equal(t, account1.ID, result.FromAccount.ID)
		require.NotEmpty(t, result.ToAccount)
		require.Equal(t, account2.ID, result.ToAccount.ID)

		// Balance
		diff1 := account1.Balance - result.FromAccount.Balance
		diff2 := result.ToAccount.Balance - account2.Balance
		require.True(t, diff1 > 0)
		require.Equal(t, diff1, diff2)
		require.True(t, diff1%amount == 0)
		k := int(diff1 / amount)
		require.True(t, k >= 1 && k <= n)
		require.NotContains(t, existed, k)
		existed[k] = true

	}

	// check final updated balance
	updatedAccount1, err := testQueries.GetAccount(context.Background(), account1.ID)
	require.NoError(t, err)
	updatedAccount2, err := testQueries.GetAccount(context.Background(), account2.ID)
	require.NoError(t, err)

	require.Equal(t, account1.Balance-int64(n)*amount, updatedAccount1.Balance)
	require.Equal(t, account2.Balance+int64(n)*amount, updatedAccount2.Balance)

	for _, transfer := range transfers {
		cleanUpTransfer(t, transfer.Transfer.ID)
		cleanUpEntry(t, transfer.FromEntry, false)
		cleanUpEntry(t, transfer.ToEntry, false)
	}
	// cleanUpAccount(t, account1)
	// cleanUpAccount(t, account2)

}

func TestTransferTxDeadlock(t *testing.T) {
	store := NewStore(testDB)
	account1 := createRandomAccount(t)
	account2 := createRandomAccount(t)

	n := 10
	amount := int64(10)
	transfers := make([]TransferTxResult, n)

	errs := make(chan error)
	results := make(chan TransferTxResult)
	defer close(errs)

	for i := 0; i < n; i++ {
		fromAccount := account1
		toAccount := account2
		if i%2 == 1 {
			fromAccount = account2
			toAccount = account1
		}

		go func() {
			transfer, err := store.TransferTx(
				context.Background(),
				TransferTxParams{
					FromAccountID: fromAccount.ID,
					ToAccountID:   toAccount.ID,
					Amount:        amount,
					Currency:      fromAccount.Currency,
				})
			errs <- err
			results <- transfer
		}()
	}

	for i := 0; i < n; i++ {
		err := <-errs
		transfers[i] = <-results
		require.NoError(t, err)
	}

	// check final updated balance
	updatedAccount1, err := testQueries.GetAccount(context.Background(), account1.ID)
	require.NoError(t, err)
	updatedAccount2, err := testQueries.GetAccount(context.Background(), account2.ID)
	require.NoError(t, err)

	require.Equal(t, account1.Balance, updatedAccount1.Balance)
	require.Equal(t, account2.Balance, updatedAccount2.Balance)

	for _, transfer := range transfers {
		cleanUpTransfer(t, transfer.Transfer.ID)
		cleanUpEntry(t, transfer.FromEntry, false)
		cleanUpEntry(t, transfer.ToEntry, false)
	}
	// cleanUpUser(t, account1.OwnerName)
	// cleanUpUser(t, account2.OwnerName)
}
