package apis

import (
	"database/sql"

	db "github.com/aniket-skroman/skroman-user-service/sqlc_lib"
)

type Store struct {
	*db.Queries
	db *sql.DB
}

func NewStore(dbs *sql.DB) *Store {
	return &Store{
		db:      dbs,
		Queries: db.New(dbs),
	}
}

func (s *Store) DBTransaction() (*sql.Tx, error) {
	return s.db.Begin()
}
