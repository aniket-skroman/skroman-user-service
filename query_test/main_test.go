package querytest

import (
	"database/sql"
	"log"
	"os"
	"testing"

	sqlc_lib "github.com/aniket-skroman/skroman-user-service/sqlc_lib"
	_ "github.com/lib/pq"
)

var (
	db_driver   = "postgres"
	db_source   = "postgresql://postgres:root@localhost:5432/skroman_users?sslmode=disable"
	testQueries *sqlc_lib.Queries
	testDB      *sql.DB
)

func TestMain(m *testing.M) {
	var err error
	testDB, err = sql.Open(db_driver, db_source)

	if err != nil {
		log.Fatal(err)
	}

	testQueries = sqlc_lib.New(testDB)
	os.Exit(m.Run())
}
