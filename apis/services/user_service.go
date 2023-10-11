package services

import (
	"errors"

	"github.com/aniket-skroman/skroman-user-service/apis/dtos"
	"github.com/aniket-skroman/skroman-user-service/apis/helper"
	"github.com/aniket-skroman/skroman-user-service/apis/repositories"
	db "github.com/aniket-skroman/skroman-user-service/sqlc_lib"
	"github.com/aniket-skroman/skroman-user-service/utils"
)

type UserService interface {
	CreateNewUser(dtos.CreateUserRequestDTO) (dtos.UserDTO, error)
	FetchUserByEmail(dtos.LoginUserRequestDTO) (dtos.UserDTO, error)
}

type user_service struct {
	user_repo   repositories.UserRepository
	jwt_service JWTService
}

func NewUserService(repo repositories.UserRepository, jwt_service JWTService) UserService {
	return &user_service{
		user_repo:   repo,
		jwt_service: jwt_service,
	}
}

func (ser *user_service) CreateNewUser(req dtos.CreateUserRequestDTO) (dtos.UserDTO, error) {
	if req.UserType == "" || len(req.UserType) == 0 {
		return dtos.UserDTO{}, errors.New("user type can not be empty")
	}

	//		check for duplicate account create
	dup_args := db.CheckFullNameAndMailIDParams{
		Email:    req.Email,
		FullName: req.FullName,
	}

	result, err := ser.user_repo.CheckDuplicateUser(dup_args)

	if err != nil {
		return dtos.UserDTO{}, err
	}
	if result != 0 {
		return dtos.UserDTO{}, errors.New("this account already exits")
	}

	hash_password := utils.Hash_password(req.Password)

	args := db.CreateNewUserParams{
		FullName: req.FullName,
		Email:    req.Email,
		Password: hash_password,
		Contact:  req.Contact,
		UserType: req.UserType,
	}

	user, err := ser.user_repo.CreateNewUser(args)
	err = helper.Handle_DBError(err)

	if err != nil {
		return dtos.UserDTO{}, err
	}

	return dtos.UserDTO{
		ID:        user.ID,
		FullName:  user.FullName,
		Email:     user.Email,
		Contact:   user.Contact,
		UserType:  user.UserType,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}, nil
}

func (ser *user_service) FetchUserByEmail(req dtos.LoginUserRequestDTO) (dtos.UserDTO, error) {
	user, err := ser.user_repo.FetchUserByMultipleTag(req.Email)

	if err != nil {
		return dtos.UserDTO{}, err
	}

	if !utils.Compare_password(user.Password, []byte(req.Password)) {
		return dtos.UserDTO{}, errors.New("password does not matched")
	}

	token := ser.jwt_service.GenerateToken(user.ID.String(), user.UserType)

	return dtos.UserDTO{
		ID:          user.ID,
		FullName:    user.FullName,
		Email:       user.Email,
		Contact:     user.Contact,
		UserType:    user.UserType,
		AccessToken: token,
		CreatedAt:   user.CreatedAt,
		UpdatedAt:   user.UpdatedAt,
	}, nil
}
