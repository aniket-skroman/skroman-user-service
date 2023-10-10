package repositories

import (
	db "github.com/aniket-skroman/skroman-user-service/sqlc_lib"
)

func (repo *user_repository) CreateNewUser(args db.CreateNewUserParams) (db.Users, error) {
	ctx, cancel := repo.Init()
	defer cancel()
	return repo.db.Queries.CreateNewUser(ctx, args)
}
