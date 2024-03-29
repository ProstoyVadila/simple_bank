package db

import (
	"context"
	"database/sql"
	"testing"
	"time"

	"github.com/ProstoyVadila/simple_bank/utils"
	"github.com/stretchr/testify/require"
)

func cleanUpEntry(t *testing.T, entry Entry, withAcc bool) {
	t.Cleanup(func() {
		err := testQueries.DeleteEntry(context.Background(), entry.ID)
		require.NoError(t, err)
		if withAcc {
			err = testQueries.DeleteAccount(context.Background(), entry.AccountID)
			require.NoError(t, err)
		}
	})
}

func createRandomEntry(t *testing.T, acc Account) Entry {
	amount := utils.RandomBalance()

	entry, err := testQueries.CreateEntry(context.Background(), CreateEntryParams{
		AccountID: acc.ID,
		Amount:    amount,
	})
	require.NoError(t, err)
	require.NotEmpty(t, entry)
	require.NotZero(t, entry.ID)
	require.NotZero(t, entry.CreatedAt)
	require.Equal(t, amount, entry.Amount)

	return entry
}

func TestCreateEntry(t *testing.T) {
	entry := createRandomEntry(t, createRandomAccount(t))
	cleanUpEntry(t, entry, true)
}

func TestGetEntry(t *testing.T) {
	entry1 := createRandomEntry(t, createRandomAccount(t))
	entry2, err := testQueries.GetEntry(context.Background(), entry1.ID)
	require.NoError(t, err)
	require.NotEmpty(t, entry2)
	require.WithinDuration(t, entry1.CreatedAt, entry2.CreatedAt, time.Second)
	require.Equal(t, entry1, entry2)
	cleanUpEntry(t, entry1, true)
}

func TestGetEntriesByAccount(t *testing.T) {
	acc := createRandomAccount(t)

	var entries1 [5]Entry
	for i := 0; i < len(entries1); i++ {
		entries1[i] = createRandomEntry(t, acc)
	}

	entries2, err := testQueries.GetEntriesByAccount(context.Background(), acc.ID)
	require.NoError(t, err)
	require.Len(t, entries2, len(entries1))
	for i := 0; i < len(entries1); i++ {
		require.NotEmpty(t, entries2[i])
		require.NotZero(t, entries2[i].ID)
		require.NotZero(t, entries2[i].CreatedAt)
		require.Equal(t, entries1[i], entries2[i])
	}
	for _, entry := range entries1 {
		cleanUpEntry(t, entry, false)
	}
	cleanUpAccount(t, acc)
}

func TestListEntries(t *testing.T) {
	acc := createRandomAccount(t)
	n := 10
	entries1 := make([]Entry, n)
	for i := 0; i < n; i++ {
		entries1[i] = createRandomEntry(t, acc)
	}

	entries2, err := testQueries.ListEntries(context.Background(), ListEntriesParams{
		Limit:  5,
		Offset: 5,
	})
	require.NoError(t, err)
	require.Len(t, entries2, 5)
	for _, entry := range entries2 {
		require.NotEmpty(t, entry)
	}
	for _, entry := range entries1 {
		cleanUpEntry(t, entry, false)
	}
	cleanUpAccount(t, acc)
}

func TestUpdateEntry(t *testing.T) {
	entry1 := createRandomEntry(t, createRandomAccount(t))
	args := UpdateEntryParams{
		ID:     entry1.ID,
		Amount: utils.RandomBalance(),
	}

	entry2, err := testQueries.UpdateEntry(context.Background(), args)
	require.NoError(t, err)
	require.NotEmpty(t, entry2)

	require.Equal(t, entry1.ID, entry2.ID)
	require.Equal(t, entry1.AccountID, entry2.AccountID)
	require.Equal(t, entry1.CreatedAt, entry2.CreatedAt)
	require.Equal(t, args.Amount, entry2.Amount)
	require.NotEqual(t, entry1.Amount, entry2.Amount)

	cleanUpEntry(t, entry1, true)
}

func TestDeleteEntry(t *testing.T) {
	acc := createRandomAccount(t)
	entry1 := createRandomEntry(t, acc)
	err := testQueries.DeleteEntry(context.Background(), entry1.ID)
	require.NoError(t, err)

	entry2, err := testQueries.GetEntry(context.Background(), entry1.ID)
	require.Error(t, err)
	require.EqualError(t, err, sql.ErrNoRows.Error())
	require.Empty(t, entry2)

	cleanUpAccount(t, acc)
}

func TestDeleteEntriesByAccount(t *testing.T) {
	acc := createRandomAccount(t)
	for i := 0; i < 10; i++ {
		createRandomEntry(t, acc)
	}
	err := testQueries.DeleteEntriesByAccount(context.Background(), acc.ID)
	require.NoError(t, err)

	entries, err := testQueries.GetEntriesByAccount(context.Background(), acc.ID)
	require.NoError(t, err)
	require.Empty(t, entries)

	cleanUpAccount(t, acc)
}
