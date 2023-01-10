package db

import (
	"context"
	"testing"
	"time"

	"github.com/ProstoyVadila/simple_bank/utils"
	"github.com/stretchr/testify/require"
)

func cleanUpUser(t *testing.T, username string) {
	t.Cleanup(func() {
		err := testQueries.DeleteUser(context.Background(), username)
		require.NoError(t, err)
	})
}

func createRandomUser(t *testing.T) User {
	arg := CreateUserParams{
		Username:       utils.RandomOwner(),
		HashedPassword: "password",
		FullName:       utils.RandomOwner(),
		Email:          utils.RandomEmail(),
	}
	user, err := testQueries.CreateUser(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, user)

	require.Equal(t, arg.Username, user.Username)
	require.Equal(t, arg.HashedPassword, user.HashedPassword)
	require.Equal(t, arg.FullName, user.FullName)
	require.Equal(t, arg.Email, user.Email)

	require.Zero(t, user.PasswordChangedAt)
	require.NotZero(t, user.CreatedAt)

	return user
}

func TestCreateUser(t *testing.T) {
	user := createRandomUser(t)
	cleanUpUser(t, user.Username)
}

func TestGetUser(t *testing.T) {
	user1 := createRandomUser(t)
	user2, err := testQueries.GetUser(context.Background(), user1.Username)
	require.NoError(t, err)
	require.NotEmpty(t, user2)
	require.WithinDuration(t, user1.CreatedAt, user2.CreatedAt, time.Second)
	require.Equal(t, user1, user2)

	cleanUpUser(t, user1.Username)
}
