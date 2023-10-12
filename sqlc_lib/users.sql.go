// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.22.0
// source: users.sql

package db

import (
	"context"
	"database/sql"

	"github.com/google/uuid"
)

const checkEmailOrContactExists = `-- name: CheckEmailOrContactExists :execrows
select id, full_name, email, password, contact, user_type, created_at, updated_at from users
where email=$1 or contact=$2
limit 1
`

type CheckEmailOrContactExistsParams struct {
	Email   string `json:"email"`
	Contact string `json:"contact"`
}

func (q *Queries) CheckEmailOrContactExists(ctx context.Context, arg CheckEmailOrContactExistsParams) (int64, error) {
	result, err := q.db.ExecContext(ctx, checkEmailOrContactExists, arg.Email, arg.Contact)
	if err != nil {
		return 0, err
	}
	return result.RowsAffected()
}

const checkForContact = `-- name: CheckForContact :one
select id, full_name, email, password, contact, user_type, created_at, updated_at from users 
where contact = $1
and id <> $2
`

type CheckForContactParams struct {
	Contact string    `json:"contact"`
	ID      uuid.UUID `json:"id"`
}

func (q *Queries) CheckForContact(ctx context.Context, arg CheckForContactParams) (Users, error) {
	row := q.db.QueryRowContext(ctx, checkForContact, arg.Contact, arg.ID)
	var i Users
	err := row.Scan(
		&i.ID,
		&i.FullName,
		&i.Email,
		&i.Password,
		&i.Contact,
		&i.UserType,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}

const countUsers = `-- name: CountUsers :one
select count(*) from users
`

func (q *Queries) CountUsers(ctx context.Context) (int64, error) {
	row := q.db.QueryRowContext(ctx, countUsers)
	var count int64
	err := row.Scan(&count)
	return count, err
}

const createNewUser = `-- name: CreateNewUser :one
insert into users (
    full_name,
    email,
    password,
    contact,
    user_type
) values (
    $1,$2,$3,$4,$5
) returning id, full_name, email, password, contact, user_type, created_at, updated_at
`

type CreateNewUserParams struct {
	FullName string `json:"full_name"`
	Email    string `json:"email"`
	Password string `json:"password"`
	Contact  string `json:"contact"`
	UserType string `json:"user_type"`
}

func (q *Queries) CreateNewUser(ctx context.Context, arg CreateNewUserParams) (Users, error) {
	row := q.db.QueryRowContext(ctx, createNewUser,
		arg.FullName,
		arg.Email,
		arg.Password,
		arg.Contact,
		arg.UserType,
	)
	var i Users
	err := row.Scan(
		&i.ID,
		&i.FullName,
		&i.Email,
		&i.Password,
		&i.Contact,
		&i.UserType,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}

const deleteUser = `-- name: DeleteUser :execrows
delete from users 
where id = $1
`

func (q *Queries) DeleteUser(ctx context.Context, id uuid.UUID) (int64, error) {
	result, err := q.db.ExecContext(ctx, deleteUser, id)
	if err != nil {
		return 0, err
	}
	return result.RowsAffected()
}

const fetchAllUsers = `-- name: FetchAllUsers :many
select id, full_name, email, password, contact, user_type, created_at, updated_at from users 
order by created_at desc 
limit $1
offset $2
`

type FetchAllUsersParams struct {
	Limit  int32 `json:"limit"`
	Offset int32 `json:"offset"`
}

func (q *Queries) FetchAllUsers(ctx context.Context, arg FetchAllUsersParams) ([]Users, error) {
	rows, err := q.db.QueryContext(ctx, fetchAllUsers, arg.Limit, arg.Offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []Users{}
	for rows.Next() {
		var i Users
		if err := rows.Scan(
			&i.ID,
			&i.FullName,
			&i.Email,
			&i.Password,
			&i.Contact,
			&i.UserType,
			&i.CreatedAt,
			&i.UpdatedAt,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const getUserByEmailOrContact = `-- name: GetUserByEmailOrContact :one
select id, full_name, email, password, contact, user_type, created_at, updated_at from users
where email=$1 or contact = $1
limit 1
`

func (q *Queries) GetUserByEmailOrContact(ctx context.Context, email string) (Users, error) {
	row := q.db.QueryRowContext(ctx, getUserByEmailOrContact, email)
	var i Users
	err := row.Scan(
		&i.ID,
		&i.FullName,
		&i.Email,
		&i.Password,
		&i.Contact,
		&i.UserType,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}

const updateUser = `-- name: UpdateUser :execresult
update users 
set full_name=$2,
contact=$3,
user_type=$4,
updated_at = CURRENT_TIMESTAMP
where id=$1 
returning id, full_name, email, password, contact, user_type, created_at, updated_at
`

type UpdateUserParams struct {
	ID       uuid.UUID `json:"id"`
	FullName string    `json:"full_name"`
	Contact  string    `json:"contact"`
	UserType string    `json:"user_type"`
}

func (q *Queries) UpdateUser(ctx context.Context, arg UpdateUserParams) (sql.Result, error) {
	return q.db.ExecContext(ctx, updateUser,
		arg.ID,
		arg.FullName,
		arg.Contact,
		arg.UserType,
	)
}
