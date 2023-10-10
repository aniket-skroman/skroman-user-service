package querytest

import (
	"context"
	"fmt"
	"testing"

	db "github.com/aniket-skroman/skroman-user-service/sqlc_lib"
	"github.com/aniket-skroman/skroman-user-service/utils"
	"github.com/stretchr/testify/require"
)

func create_random_user(t *testing.T) db.Users {
	args := db.CreateNewUserParams{
		FullName: "sdsddsdsd",
		Email:    fmt.Sprintf("%s@gmail.com", "dhbhdb"),
		Password: "user123",
		Contact:  utils.RandomInt(),
		UserType: "",
	}

	user, err := testQueries.CreateNewUser(context.Background(), args)
	fmt.Println(err, user)
	require.NoError(t, err)
	require.NotEmpty(t, user)
	require.NotZero(t, user.ID)

	return user
}

func TestCreateNewUser(t *testing.T) {
	create_random_user(t)
}

func TestGetUserByMultipleTags(t *testing.T) {
	user, err := testQueries.GetUserByEmailOrContact(context.Background(), "3265166")

	require.NoError(t, err)
	require.NotEmpty(t, user)

	fmt.Println("user found : ", user)
}
