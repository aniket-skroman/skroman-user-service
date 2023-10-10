package dtos

import (
	"time"

	"github.com/google/uuid"
)

type CreateUserRequestDTO struct {
	FullName string `json:"full_name" binding:"required"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
	Contact  string `json:"contact" binding:"required, min=10"`
	UserType string `json:"user_type" binding:"required"`
}

type UserDTO struct {
	ID        uuid.UUID `json:"id"`
	FullName  string    `json:"full_name"`
	Email     string    `json:"email"`
	Password  string    `json:"-,omitempty"`
	Contact   string    `json:"contact"`
	UserType  string    `json:"user_type"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
