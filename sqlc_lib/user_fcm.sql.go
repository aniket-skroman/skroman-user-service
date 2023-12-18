// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.24.0
// source: user_fcm.sql

package db

import (
	"context"

	"github.com/google/uuid"
)

const createUserFCMData = `-- name: CreateUserFCMData :one
insert into user_fcm_data (
    user_id,
    fcm_token
) values (
    $1,$2
) returning id, user_id, fcm_token, created_at, updated_at
`

type CreateUserFCMDataParams struct {
	UserID   uuid.UUID `json:"user_id"`
	FcmToken string    `json:"fcm_token"`
}

func (q *Queries) CreateUserFCMData(ctx context.Context, arg CreateUserFCMDataParams) (UserFcmData, error) {
	row := q.db.QueryRowContext(ctx, createUserFCMData, arg.UserID, arg.FcmToken)
	var i UserFcmData
	err := row.Scan(
		&i.ID,
		&i.UserID,
		&i.FcmToken,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}

const fetchFCMTokensByUser = `-- name: FetchFCMTokensByUser :many
select id, user_id, fcm_token, created_at, updated_at from user_fcm_data
where user_id = $1
`

func (q *Queries) FetchFCMTokensByUser(ctx context.Context, userID uuid.UUID) ([]UserFcmData, error) {
	rows, err := q.db.QueryContext(ctx, fetchFCMTokensByUser, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []UserFcmData{}
	for rows.Next() {
		var i UserFcmData
		if err := rows.Scan(
			&i.ID,
			&i.UserID,
			&i.FcmToken,
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