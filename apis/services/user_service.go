package services

import (
	"errors"

	"github.com/aniket-skroman/skroman-user-service/apis/dtos"
	"github.com/aniket-skroman/skroman-user-service/apis/repositories"
	db "github.com/aniket-skroman/skroman-user-service/sqlc_lib"
)

type UserService interface {
	CreateNewUser(dtos.CreateUserRequestDTO) (dtos.UserDTO, error)
}

type user_service struct {
	user_repo repositories.UserRepository
}

func NewUserService(repo repositories.UserRepository) UserService {
	return &user_service{
		user_repo: repo,
	}
}

func (ser *user_service) CreateNewUser(req dtos.CreateUserRequestDTO) (dtos.UserDTO, error) {
	if req.UserType == "" || len(req.UserType) == 0 {
		return dtos.UserDTO{}, errors.New("user type can not be empty")
	}

	args := db.CreateNewUserParams{
		FullName: req.FullName,
		Email:    req.Email,
		Password: req.Password,
		Contact:  req.Contact,
		UserType: req.UserType,
	}

	user, err := ser.user_repo.CreateNewUser(args)

	if err != nil {
		return dtos.UserDTO{}, err
	}

	return dtos.UserDTO(user), nil
}
