package repositories

import (
	"context"
	"time"

	"github.com/aniket-skroman/skroman-user-service/apis"
	db "github.com/aniket-skroman/skroman-user-service/sqlc_lib"
)

type UserRepository interface {
	Init() (context.Context, context.CancelFunc)
	CreateNewUser(db.CreateNewUserParams) (db.Users, error)
	CheckDuplicateUser(db.CheckFullNameAndMailIDParams) (int64, error)
	FetchUserByMultipleTag(string) (db.Users, error)
}

type user_repository struct {
	db *apis.Store
}

func NewUserRepository(db *apis.Store) UserRepository {
	return &user_repository{
		db: db,
	}
}

func (repo *user_repository) Init() (context.Context, context.CancelFunc) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	return ctx, cancel
}
