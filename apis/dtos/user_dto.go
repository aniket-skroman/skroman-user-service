package dtos

import (
	"time"

	db "github.com/aniket-skroman/skroman-user-service/sqlc_lib"
	"github.com/google/uuid"
)

type CreateUserRequestDTO struct {
	FullName   string `json:"full_name" binding:"required"`
	Email      string `json:"email" binding:"required,email"`
	Password   string `json:"password" binding:"required"`
	Contact    string `json:"contact" binding:"required,min=10"`
	Department string `json:"department" binding:"required,oneof=SALES INSTALLATION ACCOUNT INVENTORY PRODUCTION"`
	UserType   string `json:"user_type" binding:"required,oneof=EMP ADMIN SUPER-ADMIN"`
}

type UpdateUserRequestDTO struct {
	ID       string `json:"id" binding:"required"`
	FullName string `json:"full_name" binding:"required"`
	Contact  string `json:"contact" binding:"required,min=10"`
	UserType string `json:"user_type" binding:"required"`
}

type LoginUserRequestDTO struct {
	Email    string `form:"email" binding:"required,email"`
	Password string `form:"password" binding:"required"`
}

type GetUsersRequestParams struct {
	PageID     int32  `uri:"page_id" binding:"required"`
	PageSize   int32  `uri:"page_size" binding:"required"`
	Department string `uri:"department" binding:"omitempty,oneof=SALES INSTALLATION ACCOUNT INVENTORY PRODUCTION"`
}

type DeleteUserRequestDTO struct {
	UserId string `uri:"user_id"`
}

type UserDTO struct {
	ID          uuid.UUID  `json:"id"`
	FullName    string     `json:"full_name"`
	Email       string     `json:"email,omitempty"`
	Contact     string     `json:"contact"`
	UserType    string     `json:"user_type"`
	Department  string     `json:"department"`
	AccessToken string     `json:"access_token,omitempty"`
	CreatedAt   *time.Time `json:"created_at,omitempty"`
	UpdatedAt   *time.Time `json:"updated_at,omitempty"`
}

func (user *UserDTO) MakeUserDTO(access_token string, module_data ...db.Users) interface{} {
	if len(module_data) == 1 {
		new_user := UserDTO{
			ID:          module_data[0].ID,
			FullName:    module_data[0].FullName,
			Email:       module_data[0].Email,
			Contact:     module_data[0].Contact,
			UserType:    module_data[0].UserType,
			AccessToken: access_token,
			Department:  module_data[0].Department,
			CreatedAt:   &module_data[0].CreatedAt,
			UpdatedAt:   &module_data[0].UpdatedAt,
		}

		return new_user
	}

	users := make([]UserDTO, len(module_data))

	for i := range module_data {
		users[i] = UserDTO{
			ID:         module_data[i].ID,
			FullName:   module_data[i].FullName,
			Email:      module_data[i].Email,
			Contact:    module_data[i].Contact,
			UserType:   module_data[i].UserType,
			Department: module_data[i].Department,
			CreatedAt:  &module_data[i].CreatedAt,
			UpdatedAt:  &module_data[i].UpdatedAt,
		}
	}

	return users
}

//----------------------------------- 	HANDLE SKROMAN CLIENT OPERATIONS ------------------------------------------- //

type CreateSkromanClientRequestDTO struct {
	UserName string `json:"user_name" binding:"required"`
	Email    string `json:"email"`
	Password string `json:"password"`
	Contact  string `json:"contact" binding:"required"`
	Address  string `json:"address" binding:"required"`
	City     string `json:"city" `
	State    string `json:"state" `
	Pincode  string `json:"pincode" `
}

type SkromanClientDTO struct {
	ID        uuid.UUID `json:"id"`
	UserName  string    `json:"user_name" `
	Email     string    `json:"email" binding:"required,email"`
	Password  string    `json:"password,omitempty"`
	Contact   string    `json:"contact" `
	Address   string    `json:"address" `
	City      string    `json:"city" `
	State     string    `json:"state" `
	Pincode   string    `json:"pincode" `
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
