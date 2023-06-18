package db

import (
	"context"
	"testing"

	util "github.com/KHarshit1203/simple-bank/util/datagenerator"
	"github.com/stretchr/testify/require"
)

func TestCreateUser(t *testing.T) {
	generateUser(t)
}

func TestGetUser(t *testing.T) {
	user1 := generateUser(t)

	User2, err := testStore.GetUser(context.Background(), user1.Username)
	require.NoError(t, err)
	require.NotEmpty(t, User2)

	require.Equal(t, user1.Username, User2.Username)
	require.Equal(t, user1.FullName, User2.FullName)
	require.Equal(t, user1.HashedPassword, User2.HashedPassword)
	require.Equal(t, user1.Email, User2.Email)
	require.Equal(t, user1.CreatedAt, User2.CreatedAt)

}

// generateUser generates random user in DB.
func generateUser(t *testing.T) User {
	arg := CreateUserParams{
		Username:       util.RandomString(10),
		HashedPassword: "secret",
		FullName:       util.RandomString(20),
		Email:          util.RandomEmail(),
	}

	user, err := testStore.CreateUser(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, user)

	require.Equal(t, arg.Username, user.Username)
	require.Equal(t, arg.HashedPassword, user.HashedPassword)
	require.Equal(t, arg.FullName, user.FullName)
	require.Equal(t, arg.Email, user.Email)

	require.True(t, user.PasswordChangetdAt.IsZero())
	require.NotZero(t, user.CreatedAt)

	return user
}
