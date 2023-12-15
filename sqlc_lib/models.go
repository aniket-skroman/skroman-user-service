// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.24.0

package db

import (
	"database/sql"
	"time"

	"github.com/google/uuid"
)

type SkromanClient struct {
	ID        uuid.UUID      `json:"id"`
	UserName  string         `json:"user_name"`
	Email     string         `json:"email"`
	Password  sql.NullString `json:"password"`
	Contact   string         `json:"contact"`
	Address   string         `json:"address"`
	City      sql.NullString `json:"city"`
	State     sql.NullString `json:"state"`
	Pincode   sql.NullString `json:"pincode"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
}

type UserFcmData struct {
	ID        uuid.UUID `json:"id"`
	UserID    uuid.UUID `json:"user_id"`
	FcmToken  string    `json:"fcm_token"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type Users struct {
	ID         uuid.UUID `json:"id"`
	FullName   string    `json:"full_name"`
	Email      string    `json:"email"`
	Password   string    `json:"password"`
	Contact    string    `json:"contact"`
	UserType   string    `json:"user_type"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
	Department string    `json:"department"`
	EmpCode    string    `json:"emp_code"`
}
