package querytest

import (
	"context"
	"fmt"
	"testing"

	db "github.com/aniket-skroman/skroman-user-service/sqlc_lib"
	"github.com/aniket-skroman/skroman-user-service/utils"
	"github.com/google/uuid"
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

func TestFetchAllUsers(t *testing.T) {
	args := db.FetchAllUsersParams{
		Limit:  10,
		Offset: 0,
	}

	users, err := testQueries.FetchAllUsers(context.Background(), args)

	require.NoError(t, err)
	fmt.Println("Users : ", users)
	require.NotEmpty(t, users)
}

func TestUpdateUser(t *testing.T) {
	user_id, err := uuid.Parse("72577a45-a4dc-4734-a74c-f2102e9e1381")

	require.NoError(t, err)

	tx, err := testDB.Begin()
	require.NoError(t, err)

	qtx := testQueries.WithTx(tx)

	require.NotEmpty(t, qtx)

	// first check for new contact is a unique
	cont_args := db.CheckForContactParams{
		Contact: "7720830160",
		ID:      user_id,
	}

	user, err := qtx.CheckForContact(context.Background(), cont_args)

	require.Error(t, err)
	require.Empty(t, user)

	args := db.UpdateUserParams{
		ID:       user_id,
		FullName: "test",
		Contact:  "7720830160",
		UserType: "EMP",
	}

	result, err := qtx.UpdateUser(context.Background(), args)

	if err != nil {
		if err := tx.Rollback(); err != nil {
			fmt.Println(err)
			return
		}
	}

	affected_rows, err := result.RowsAffected()
	require.NoError(t, err)
	require.NotEqual(t, affected_rows, 0)

	tx.Commit()
}

func TestUserCount(t *testing.T) {
	count, err := testQueries.CountUsers(context.Background())

	require.NoError(t, err)
	require.NotZero(t, count)
}

func TestGetUserById(t *testing.T) {
	user_id, err := uuid.Parse("2c87109d-52c3-4b2f-b22e-98d61fec12b3")
	require.NoError(t, err)

	user, err := testQueries.GetUserById(context.Background(), user_id)

	require.NoError(t, err)
	require.NotEmpty(t, user)

}
