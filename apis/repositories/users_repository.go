package repositories

import (
	"context"
	"database/sql"
	"time"

	"github.com/aniket-skroman/skroman-user-service/apis"
	db "github.com/aniket-skroman/skroman-user-service/sqlc_lib"
	"github.com/google/uuid"
)

type UserRepository interface {
	Init() (context.Context, context.CancelFunc)
	CreateNewUser(db.CreateNewUserParams) (db.Users, error)
	CheckDuplicateUser(db.CheckEmailOrContactExistsParams) (int64, error)
	FetchUserByMultipleTag(string) (db.Users, error)
	UpdateUser(db.UpdateUserParams) (sql.Result, error)
	CheckForContact(db.CheckForContactParams) (db.Users, error)
	FetchAllUsers(db.FetchAllUsersParams) ([]db.Users, error)
	DeleteUser(uuid.UUID) (int64, error)
	CountUsers() (int64, error)
	FetchUserById(uuid.UUID) (db.Users, error)
	FetchUsersByDepartment(args db.UsersByDepartmentParams) ([]db.Users, error)
	CountUserByDepartment(dept_name string) (int64, error)

	CreateSkromanClient(args db.CreateSkromanUserParams) (db.SkromanClient, error)
	FetchAllClients(args db.FetchAllClientsParams) ([]db.SkromanClient, error)
	CountOfClient() (int64, error)
	DeleteClient(client_id uuid.UUID) (sql.Result, error)
	FetchClientById(client_id uuid.UUID) (db.SkromanClient, error)
	UpdateSkromanClientInfo(args db.UpdateSkromanClientInfoParams) (db.SkromanClient, error)
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
