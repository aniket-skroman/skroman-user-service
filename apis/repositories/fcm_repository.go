package repositories

import (
	"context"
	"time"

	"github.com/aniket-skroman/skroman-user-service/apis"
	db "github.com/aniket-skroman/skroman-user-service/sqlc_lib"
	"github.com/google/uuid"
)

type FCMRepository interface {
	Init() (context.Context, context.CancelFunc)
	CreateUserFCMData(args db.CreateUserFCMDataParams) (db.UserFcmData, error)
	FetchTokenByUsers(user_id uuid.UUID) ([]db.UserFcmData, error)
}

type fcm_repo struct {
	db *apis.Store
}

func NewFCMRepository(db *apis.Store) FCMRepository {
	return &fcm_repo{
		db: db,
	}
}

func (repo *fcm_repo) Init() (context.Context, context.CancelFunc) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	return ctx, cancel
}
