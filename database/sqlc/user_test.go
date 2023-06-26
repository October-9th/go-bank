package sqlc

import (
	"context"
	"testing"
	"time"

	"github.com/October-9th/simple-bank/util"
	"github.com/stretchr/testify/require"
)

func CreateRandomUser(t *testing.T) User {
	hashedPassword, err := util.HashPassword(util.RandomString(10))
	require.NoError(t, err)

	arg := CreateUserParams{
		Username:       util.RandomOwner(),
		HashedPassword: hashedPassword,
		FullName:       util.RandomOwner(),
		Email:          util.RandomEmail(),
	}
	user, err := testQueries.CreateUser(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, user)

	require.Equal(t, arg.Username, user.Username)
	require.Equal(t, arg.HashedPassword, user.HashedPassword)
	require.Equal(t, arg.Email, user.Email)

	require.NotZero(t, user.CreatedAt)
	require.True(t, user.PasswordChangedAt.IsZero())
	return user
}
func TestCreateUser(t *testing.T) {
	CreateRandomUser(t)
}

func TestGetUser(t *testing.T) {
	user_test := CreateRandomUser(t)
	user, err := testQueries.GetUser(context.Background(), user_test.Username)

	require.NoError(t, err)
	require.NotEmpty(t, user)

	require.Equal(t, user.Username, user_test.Username)
	require.Equal(t, user.HashedPassword, user_test.HashedPassword)
	require.Equal(t, user.FullName, user_test.FullName)
	require.Equal(t, user.Email, user_test.Email)
	// Use require.WithinDuration to check 2 timestamps are different by at most some delta duration
	require.WithinDuration(t, user.CreatedAt, user_test.CreatedAt, time.Second)
	require.WithinDuration(t, user.PasswordChangedAt, user_test.PasswordChangedAt, time.Second)
}
