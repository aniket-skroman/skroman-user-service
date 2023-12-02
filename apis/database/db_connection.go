package database

import (
	"database/sql"
)

var (
	dbDriver = "postgres"
	dbSource = "postgresql://postgres:support12@skroman-user.ckwveljlsuux.ap-south-1.rds.amazonaws.com:5432/skroman_users"
)

func make_db_connection() (*sql.DB, error) {
	db, err := sql.Open(dbDriver, dbSource)

	if err != nil {
		return nil, err
	}

	err = db.Ping()
	return db, err
}

var DB_INSTANCE = make_db_connection

func CloseDBConnection(db *sql.DB) error {
	if err := db.Close(); err != nil {
		return err
	}

	return nil
}
