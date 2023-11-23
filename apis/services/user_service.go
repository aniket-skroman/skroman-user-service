package services

import (
	"errors"
	"fmt"
	"reflect"
	"strings"
	"sync"

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
	GetUsersCount() int32
	FetchUserById(uuid.UUID) (dtos.UserDTO, error)
	GetUsersByDepartmentCount(dep_name string) (int64, error)
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
	dupl_args := db.CheckEmailOrContactExistsParams{
		Email:   req.Email,
		Contact: req.Contact,
	}
	result, err := ser.user_repo.CheckDuplicateUser(dupl_args)

	if err != nil {
		return dtos.UserDTO{}, err
	}
	if result != 0 {
		return dtos.UserDTO{}, errors.New("this account already exits")
	}

	// check for duplicate

	hash_password := utils.Hash_password(req.Password)

	args := db.CreateNewUserParams{
		FullName:   req.FullName,
		Email:      req.Email,
		Password:   hash_password,
		Contact:    req.Contact,
		UserType:   req.UserType,
		Department: req.Department,
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
	err_chan := make(chan error)
	wg := sync.WaitGroup{}
	var n_user interface{}
	user, err := ser.user_repo.FetchUserByMultipleTag(req.Email)

	wg.Add(2)

	go func() {
		defer wg.Done()
		if !utils.Compare_password(user.Password, []byte(req.Password)) {
			err_chan <- errors.New("password does not matched")
		}
	}()

	go func() {
		defer wg.Done()
		user_id, _ := helper.EncryptData(user.ID.String())
		user_type, _ := helper.EncryptData(user.UserType)
		dept, _ := helper.EncryptData(user.Department)
		token := ser.jwt_service.GenerateToken(user_id, user_type, dept)
		n_user = new(dtos.UserDTO).MakeUserDTO(token, user)
	}()

	go func() {
		err_chan <- err
		wg.Wait()
		close(err_chan)
	}()

	for p_err := range err_chan {
		if p_err != nil {
			return dtos.UserDTO{}, p_err
		}
	}

	return n_user.(dtos.UserDTO), nil
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
	// fetch count
	wg := sync.WaitGroup{}
	wg.Add(2)
	var result []db.Users
	var err error

	go func(dept_name string) {
		defer wg.Done()
		count, _ := ser.GetUsersByDepartmentCount(dept_name)
		utils.SetPaginationData(int(req.PageID), int64(count))
	}(req.Department)

	go func() {
		defer wg.Done()
		args := db.UsersByDepartmentParams{
			Limit:      req.PageSize,
			Offset:     (req.PageID - 1) * req.PageSize,
			Department: req.Department,
		}
		result, err = ser.user_repo.FetchUsersByDepartment(args)
	}()
	wg.Wait()
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

func (ser *user_service) GetUsersCount() int32 {
	count, err := ser.user_repo.CountUsers()

	if err != nil {
		return 0
	}
	return int32(count)
}

func (ser *user_service) GetUsersByDepartmentCount(dep_name string) (int64, error) {
	return ser.user_repo.CountUserByDepartment(dep_name)
}

func (ser *user_service) FetchUserById(user_id uuid.UUID) (dtos.UserDTO, error) {
	result, err := ser.user_repo.FetchUserById(user_id)

	if err != nil {
		return dtos.UserDTO{}, err
	}

	user := new(dtos.UserDTO).MakeUserDTO("", result)

	if reflect.DeepEqual(user.(dtos.UserDTO), dtos.UserDTO{}) {
		return dtos.UserDTO{}, errors.New("user not found")
	}

	return user.(dtos.UserDTO), nil
}
