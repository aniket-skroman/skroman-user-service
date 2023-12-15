package repositories

import (
	db "github.com/aniket-skroman/skroman-user-service/sqlc_lib"
	"github.com/google/uuid"
)

func (repo *fcm_repo) CreateUserFCMData(args db.CreateUserFCMDataParams) (db.UserFcmData, error) {
	ctx, cancel := repo.Init()
	defer cancel()

	return repo.db.Queries.CreateUserFCMData(ctx, args)
}

func (repo *fcm_repo) FetchTokenByUsers(user_id uuid.UUID) ([]db.UserFcmData, error) {
	ctx, cancel := repo.Init()
	defer cancel()

	return repo.db.Queries.FetchFCMTokensByUser(ctx, user_id)
}
