package querytest

import (
	"context"
	"database/sql"
	"fmt"
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

func TestUpdateSkromanClient(t *testing.T) {
	client_id, _ := uuid.Parse("1706dcea-4d9a-4d06-bf1b-830e10611d6f")

	args := db.UpdateSkromanClientInfoParams{
		ID:       client_id,
		UserName: "Test",
		Email:    "test1@gmail.com",
		Password: sql.NullString{String: "", Valid: true},
		Contact:  "8668342234",
		Address:  "test-address",
		City:     sql.NullString{String: "test-city", Valid: true},
		State:    sql.NullString{String: "test-state", Valid: true},
		Pincode:  sql.NullString{String: "414111", Valid: true},
	}

	client, err := testQueries.UpdateSkromanClientInfo(context.Background(), args)
	require.NoError(t, err)
	require.NotEmpty(t, client)
}

func TestSearchSkromanClient(t *testing.T) {
	data := "na@na.com"
	args := db.SearchClientParams{
		Limit:   100,
		Offset:  1,
		Column3: sql.NullString{String: data, Valid: true},
	}
	result, err := testQueries.SearchClient(context.Background(), args)
	require.NoError(t, err)

	fmt.Println("Count of data", len(result))

	// count
	count, err := testQueries.CountOfSearchClient(context.Background(), sql.NullString{String: data, Valid: true})
	require.NoError(t, err)
	require.NotZero(t, count)

	fmt.Println("Total count : ", count)
}
