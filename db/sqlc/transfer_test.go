package db

import (
	"context"
	"testing"

	"github.com/ProstoyVadila/simple_bank/utils"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
)

func cleanUpTransfer(t *testing.T, transferID uuid.UUID) {
	t.Cleanup(func() {
		err := testQueries.DeleteTransfer(context.Background(), transferID)
		require.NoError(t, err)
	})
}

func createRandomTransfer(t *testing.T, args CreateTransferParams) Transfer {

	trasfer, err := testQueries.CreateTransfer(context.Background(), args)
	require.NoError(t, err)
	require.NotEmpty(t, trasfer)
	require.NotZero(t, trasfer.ID)
	require.NotZero(t, trasfer.CreatedAt)

	require.Equal(t, args.Amount, trasfer.Amount)
	require.Equal(t, args.FromAccountID, trasfer.FromAccountID)
	require.Equal(t, args.ToAccountID, trasfer.ToAccountID)

	return trasfer
}

func TestCreateTransfer(t *testing.T) {
	args := CreateTransferParams{
		FromAccountID: createRandomAccount(t).ID,
		ToAccountID:   createRandomAccount(t).ID,
		Amount:        utils.RandomBalance(),
	}
	transfer := createRandomTransfer(t, args)

	cleanUpTransfer(t, transfer.ID)
}

func TestGetTransfer(t *testing.T) {
	args := CreateTransferParams{
		FromAccountID: createRandomAccount(t).ID,
		ToAccountID:   createRandomAccount(t).ID,
		Amount:        utils.RandomBalance(),
	}
	transfer1 := createRandomTransfer(t, args)

	transfer2, err := testQueries.GetTrasfer(context.Background(), transfer1.ID)
	require.NoError(t, err)
	require.Equal(t, transfer1, transfer2)

	cleanUpTransfer(t, transfer1.ID)
}

func TestGetTransfersByFromAccount(t *testing.T) {
	fromAcc := createRandomAccount(t)

	var transfers1 [5]Transfer
	for i := 0; i < len(transfers1); i++ {
		args := CreateTransferParams{
			FromAccountID: fromAcc.ID,
			ToAccountID:   createRandomAccount(t).ID,
			Amount:        utils.RandomBalance(),
		}
		transfers1[i] = createRandomTransfer(t, args)
	}

	transfers2, err := testQueries.GetTransfersByFromAccount(
		context.Background(),
		fromAcc.ID,
	)
	require.NoError(t, err)
	require.NotEmpty(t, transfers2)
	for i, trasfer := range transfers2 {
		require.Equal(t, fromAcc.ID, trasfer.FromAccountID)
		require.Equal(t, transfers1[i], trasfer)
	}
	for _, transfer := range transfers1 {
		cleanUpTransfer(t, transfer.ID)
	}
}

func TestGetTransfersByToAccount(t *testing.T) {
	toAcc := createRandomAccount(t)

	var transfers1 [5]Transfer
	for i := 0; i < len(transfers1); i++ {
		args := CreateTransferParams{
			ToAccountID:   toAcc.ID,
			FromAccountID: createRandomAccount(t).ID,
			Amount:        utils.RandomBalance(),
		}
		transfers1[i] = createRandomTransfer(t, args)
	}

	transfers2, err := testQueries.GetTransfersByToAccount(
		context.Background(),
		toAcc.ID,
	)
	require.NoError(t, err)
	require.NotEmpty(t, transfers2)
	for i, trasfer := range transfers2 {
		require.Equal(t, toAcc.ID, trasfer.ToAccountID)
		require.Equal(t, transfers1[i], trasfer)
	}

	for _, transfer := range transfers1 {
		cleanUpTransfer(t, transfer.ID)
	}
}

func TestListTrasnfers(t *testing.T) {
	n := 10
	transfers1 := make([]Transfer, n)
	for i := 0; i < n; i++ {
		args := CreateTransferParams{
			FromAccountID: createRandomAccount(t).ID,
			ToAccountID:   createRandomAccount(t).ID,
			Amount:        utils.RandomBalance(),
		}
		transfers1[i] = createRandomTransfer(t, args)
	}

	transfers, err := testQueries.ListTransfers(context.Background(), ListTransfersParams{
		Limit:  5,
		Offset: 5,
	})
	require.NoError(t, err)
	require.NotEmpty(t, transfers)
	require.Len(t, transfers, 5)

	for _, transfer := range transfers1 {
		cleanUpTransfer(t, transfer.ID)
	}
}
