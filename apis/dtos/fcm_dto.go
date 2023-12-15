package dtos

import (
	"time"

	"github.com/google/uuid"
)

type CreateUserFCMDataRequestDTO struct {
	UserID   string `json:"user_id" binding:"required"`
	FcmToken string `json:"fcm_token" binding:"required"`
}

type UserFcmData struct {
	ID        uuid.UUID `json:"id"`
	UserID    uuid.UUID `json:"user_id"`
	FcmToken  string    `json:"fcm_token"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
