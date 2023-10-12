package services

import (
	"errors"
	"fmt"
	"reflect"
	"strings"

	"github.com/aniket-skroman/skroman-user-service/apis/dtos"
	"github.com/aniket-skroman/skroman-user-service/apis/helper"
	"github.com/aniket-skroman/skroman-user-service/apis/repositories"
	db "github.com/aniket-skroman/skroman-user-service/sqlc_lib"
	"github.com/aniket-skroman/skroman-user-service/utils"
	"github.com/google/uuid"
)

type UserService interface {
	CreateNewUser(dtos.CreateUserRequestDTO) (dtos.UserDTO, error)
	FetchUserByEmail(dtos.LoginUserRequestDTO) (dtos.UserDTO, error)
	UpdateUser(dtos.UpdateUserRequestDTO) (dtos.UserDTO, error)
	FetchAllUsers(dtos.GetUsersRequestParams) ([]dtos.UserDTO, error)
	DeleteUser(dtos.DeleteUserRequestDTO) error
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
		CreatedAt: &user.CreatedAt,
		UpdatedAt: &user.UpdatedAt,
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
		CreatedAt:   &user.CreatedAt,
		UpdatedAt:   &user.UpdatedAt,
	}, nil
}

func (ser *user_service) UpdateUser(req dtos.UpdateUserRequestDTO) (dtos.UserDTO, error) {
	user_id, err := uuid.Parse(req.ID)

	if err != nil {
		return dtos.UserDTO{}, err
	}

	// check the contact should be unique
	err = ser.check_duplicate_contact(req.Contact, user_id)

	if err != nil {
		return dtos.UserDTO{}, err
	}

	args := db.UpdateUserParams{
		ID:       user_id,
		FullName: req.FullName,
		Contact:  req.Contact,
		UserType: req.UserType,
	}

	result, err := ser.user_repo.UpdateUser(args)

	if err != nil {
		return dtos.UserDTO{}, err
	}

	rows_affected, _ := result.RowsAffected()

	if rows_affected == 0 {
		return dtos.UserDTO{}, errors.New("faild to update")
	}

	return dtos.UserDTO{
		ID:       user_id,
		FullName: req.FullName,
		Contact:  req.Contact,
		UserType: req.UserType,
	}, nil
}

func (ser *user_service) check_duplicate_contact(contact string, user_id uuid.UUID) error {

	args := db.CheckForContactParams{
		Contact: contact,
		ID:      user_id,
	}

	user, err := ser.user_repo.CheckForContact(args)

	if err != nil && !strings.Contains(err.Error(), " no rows") {
		return err
	}
	fmt.Println("User found : ", user)
	if !reflect.DeepEqual(user, db.Users{}) && user.ID != user_id {
		return errors.New("contact is already used by someone")
	}

	return nil
}

func (ser *user_service) FetchAllUsers(req dtos.GetUsersRequestParams) ([]dtos.UserDTO, error) {
	args := db.FetchAllUsersParams{
		Limit:  req.PageSize,
		Offset: (req.PageID - 1) * req.PageSize,
	}

	result, err := ser.user_repo.FetchAllUsers(args)
	err = helper.Handle_DBError(err)

	if err != nil {
		return nil, err
	}

	if len(result) == 0 {
		return nil, errors.New(utils.DATA_NOT_FOUND)
	}

	users := new(dtos.UserDTO).MakeUserDTO("", result...)

	if _, ok := users.([]dtos.UserDTO); ok {
		return users.([]dtos.UserDTO), nil
	}
	s_user := users.(dtos.UserDTO)
	return []dtos.UserDTO{
		{
			ID:          s_user.ID,
			FullName:    s_user.FullName,
			Email:       s_user.Email,
			Contact:     s_user.Contact,
			UserType:    s_user.UserType,
			AccessToken: "",
			CreatedAt:   s_user.CreatedAt,
			UpdatedAt:   s_user.UpdatedAt,
		},
	}, nil
}

func (ser *user_service) DeleteUser(req dtos.DeleteUserRequestDTO) error {
	user_id, err := uuid.Parse(req.UserId)

	if err != nil {
		return err
	}

	result, err := ser.user_repo.DeleteUser(user_id)

	if err != nil {
		return err
	}

	if result == 0 {
		return errors.New(utils.DELETE_FAILED)
	}

	return nil
}
