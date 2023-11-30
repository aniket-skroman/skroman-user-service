package querytest

import (
	"context"
	"database/sql"
	"log"
	"testing"

	db "github.com/aniket-skroman/skroman-user-service/sqlc_lib"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
)

func createSkromanClient(t *testing.T) db.SkromanClient {
	args := db.CreateSkromanUserParams{
		UserName: "Test",
		Email:    "test1@gmail.com",
		Password: sql.NullString{String: "", Valid: true},
		Contact:  "12345465",
		Address:  "test-address",
		City:     sql.NullString{String: "test-city", Valid: true},
		State:    sql.NullString{String: "test-state", Valid: true},
		Pincode:  sql.NullString{String: "414111", Valid: true},
	}

	user, err := testQueries.CreateSkromanUser(context.Background(), args)

	require.NoError(t, err)
	require.NotEmpty(t, user)

	return user
}

func TestCreateSkromanClient(t *testing.T) {
	createSkromanClient(t)
}

func TestFetchAllSkromanClients(t *testing.T) {
	args := db.FetchAllClientsParams{
		Limit:  10,
		Offset: 1,
	}

	clients, err := testQueries.FetchAllClients(context.Background(), args)

	require.NoError(t, err)
	require.NotEmpty(t, clients)

	for i := range clients {
		log.Println("ID \n", clients[i].ID)
	}
}

func TestFetchClientById(t *testing.T) {
	client_id, _ := uuid.Parse("7aa1899a-11de-481d-ba35-97b10c0d1a71")

	client, err := testQueries.FetchClientById(context.Background(), client_id)

	require.NoError(t, err)
	require.NotEmpty(t, client)
}
