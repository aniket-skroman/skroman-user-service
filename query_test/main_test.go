package querytest

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"testing"

	sqlc_lib "github.com/aniket-skroman/skroman-user-service/sqlc_lib"
	_ "github.com/lib/pq"
)

var (
	db_driver   = "postgres"
	db_source   = "postgresql://postgres:support12@skroman-support-db.ckwveljlsuux.ap-south-1.rds.amazonaws.com:5432/skroman-support-db?sslmode=disable"
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

func TestConnection(t *testing.T) {
	db, err := sql.Open(db_driver, db_source)

	if err != nil {
		fmt.Println(err)
	}

	fmt.Println("connection has been established", db)
}
