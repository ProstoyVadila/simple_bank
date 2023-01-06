package db

import (
	"context"
	"testing"
	"time"

	"github.com/ProstoyVadila/simple_bank/utils"
	"github.com/stretchr/testify/require"
)

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

func TestCreateEnntry(t *testing.T) {
	createRandomEntry(t, createRandomAccount(t))
}

func TestGetEntry(t *testing.T) {
	entry1 := createRandomEntry(t, createRandomAccount(t))
	entry2, err := testQueries.GetEntry(context.Background(), entry1.ID)
	require.NoError(t, err)
	require.NotEmpty(t, entry2)
	require.WithinDuration(t, entry1.CreatedAt, entry2.CreatedAt, time.Second)
	require.Equal(t, entry1, entry2)
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
		require.NotZero(t, entries2[i].ID)
		require.NotZero(t, entries2[i].CreatedAt)
		require.Equal(t, entries1[i], entries2[i])
	}
}
