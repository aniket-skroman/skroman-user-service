package repositories

import (
	db "github.com/aniket-skroman/skroman-user-service/sqlc_lib"
)

func (repo *user_repository) CreateNewUser(args db.CreateNewUserParams) (db.Users, error) {
	ctx, cancel := repo.Init()
	defer cancel()
	return repo.db.Queries.CreateNewUser(ctx, args)
}

func (repo *user_repository) CheckDuplicateUser(args db.CheckFullNameAndMailIDParams) (int64, error) {
	ctx, cancel := repo.Init()
	defer cancel()

	return repo.db.Queries.CheckFullNameAndMailID(ctx, args)
}

func (repo *user_repository) FetchUserByMultipleTag(tag string) (db.Users, error) {
	ctx, cancel := repo.Init()
	defer cancel()
	return repo.db.Queries.GetUserByEmailOrContact(ctx, tag)
}
