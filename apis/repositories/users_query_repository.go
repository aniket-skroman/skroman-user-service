package repositories

import (
	"database/sql"

	db "github.com/aniket-skroman/skroman-user-service/sqlc_lib"
	"github.com/google/uuid"
)

func (repo *user_repository) CreateNewUser(args db.CreateNewUserParams) (db.Users, error) {
	ctx, cancel := repo.Init()
	defer cancel()
	return repo.db.Queries.CreateNewUser(ctx, args)
}

func (repo *user_repository) CheckDuplicateUser(args db.CheckEmailOrContactExistsParams) (int64, error) {
	ctx, cancel := repo.Init()
	defer cancel()

	return repo.db.Queries.CheckEmailOrContactExists(ctx, args)
}

func (repo *user_repository) FetchUserByMultipleTag(tag string) (db.Users, error) {
	ctx, cancel := repo.Init()
	defer cancel()
	return repo.db.Queries.GetUserByEmailOrContact(ctx, tag)
}

func (repo *user_repository) UpdateUser(args db.UpdateUserParams) (sql.Result, error) {
	ctx, cancel := repo.Init()
	defer cancel()

	return repo.db.Queries.UpdateUser(ctx, args)
}

func (repo *user_repository) CheckForContact(args db.CheckForContactParams) (db.Users, error) {
	ctx, cancel := repo.Init()
	defer cancel()

	return repo.db.Queries.CheckForContact(ctx, args)
}

func (repo *user_repository) FetchAllUsers(args db.FetchAllUsersParams) ([]db.Users, error) {
	ctx, cancel := repo.Init()
	defer cancel()

	return repo.db.Queries.FetchAllUsers(ctx, args)
}

func (repo *user_repository) DeleteUser(userId uuid.UUID) (int64, error) {
	ctx, cancel := repo.Init()
	defer cancel()

	return repo.db.Queries.DeleteUser(ctx, userId)
}

func (repo *user_repository) CountUsers() (int64, error) {
	ctx, cancel := repo.Init()
	defer cancel()

	return repo.db.Queries.CountUsers(ctx)
}

func (repo *user_repository) FetchUserById(user_id uuid.UUID) (db.Users, error) {
	ctx, cancel := repo.Init()
	defer cancel()

	return repo.db.Queries.GetUserById(ctx, user_id)
}

// users by department
func (repo *user_repository) FetchUsersByDepartment(args db.UsersByDepartmentParams) ([]db.Users, error) {
	ctx, cancel := repo.Init()
	defer cancel()

	return repo.db.Queries.UsersByDepartment(ctx, args)
}

// count users by department
func (repo *user_repository) CountUserByDepartment(dept_name string) (int64, error) {
	ctx, cancel := repo.Init()
	defer cancel()

	return repo.db.Queries.CountUsersByDepartment(ctx, dept_name)
}

//----------------------------------- 	HANDLE SKROMAN CLIENT OPERATIONS ------------------------------------------- //

func (repo *user_repository) CreateSkromanClient(args db.CreateSkromanUserParams) (db.SkromanClient, error) {
	ctx, cancel := repo.Init()
	defer cancel()

	return repo.db.Queries.CreateSkromanUser(ctx, args)
}

func (repo *user_repository) FetchAllClients(args db.FetchAllClientsParams) ([]db.SkromanClient, error) {
	ctx, cancel := repo.Init()
	defer cancel()

	return repo.db.Queries.FetchAllClients(ctx, args)
}

func (repo *user_repository) CountOfClient() (int64, error) {
	ctx, cancel := repo.Init()
	defer cancel()

	return repo.db.Queries.CountOFClients(ctx)
}

func (repo *user_repository) DeleteClient(client_id uuid.UUID) (sql.Result, error) {
	ctx, cancle := repo.Init()
	defer cancle()

	return repo.db.Queries.DeleteClient(ctx, client_id)
}
