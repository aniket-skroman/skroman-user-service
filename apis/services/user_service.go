package services

import (
	"database/sql"
	"errors"
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

	CreateSkromanClient(req dtos.CreateSkromanClientRequestDTO) (dtos.SkromanClientDTO, error)
	FetchAllClients(req dtos.GetUsersRequestParams) ([]dtos.SkromanClientDTO, error)
	DeleteClient(client_id string) error
	CountOFClients() (int64, error)
	FetchClientById(client_id string) (dtos.SkromanClientDTO, error)
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
		return dtos.UserDTO{}, helper.Err_Account_Already_Exists
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
		EmpCode:    req.EmpCode,
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
		return dtos.UserDTO{}, helper.ERR_INVALID_ID
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
		EmpCode:  req.EmpCode,
	}

	result, err := ser.user_repo.UpdateUser(args)
	err = helper.Handle_DBError(err)
	if err != nil {
		return dtos.UserDTO{}, err
	}

	rows_affected, _ := result.RowsAffected()

	if rows_affected == 0 {
		return dtos.UserDTO{}, helper.Err_Update_Failed
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

	if err != nil && err != sql.ErrNoRows {
		return err
	}

	if user.ID != uuid.Nil && user.ID != user_id {
		return errors.New("contact is already used by someone")
	}

	return nil
}

func (ser *user_service) FetchAllUsers(req dtos.GetUsersRequestParams) ([]dtos.UserDTO, error) {
	// fetch count
	wg := sync.WaitGroup{}
	wg.Add(2)
	var users []dtos.UserDTO
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

		users = make([]dtos.UserDTO, len(result))
		t_wg := sync.WaitGroup{}

		for i := range result {
			t_wg.Add(1)
			go ser.setUserData(&t_wg, &users[i], &result[i])
		}

		t_wg.Wait()

	}()
	wg.Wait()
	err = helper.Handle_DBError(err)

	if err != nil {
		return nil, err
	}

	if len(result) == 0 {
		return nil, helper.Err_Data_Not_Found
	}

	return users, nil
}

func (ser *user_service) setUserData(wg *sync.WaitGroup, result *dtos.UserDTO, data *db.Users) {
	defer wg.Done()

	*result = dtos.UserDTO{
		ID:          data.ID,
		FullName:    data.FullName,
		Email:       data.Email,
		Contact:     data.Contact,
		UserType:    data.UserType,
		AccessToken: "",
		EmpCode:     data.EmpCode,
		CreatedAt:   &data.CreatedAt,
		UpdatedAt:   &data.UpdatedAt,
	}
}

func (ser *user_service) DeleteUser(req dtos.DeleteUserRequestDTO) error {
	user_id, err := uuid.Parse(req.UserId)

	if err != nil {
		return helper.ERR_INVALID_ID
	}

	result, err := ser.user_repo.DeleteUser(user_id)

	if err != nil {
		return err
	}

	if result == 0 {
		return helper.Err_Delete_Failed
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

	// if reflect.DeepEqual(user.(dtos.UserDTO), dtos.UserDTO{}) {
	// 	return dtos.UserDTO{}, helper.Err_Data_Not_Found
	// }

	if (user.(dtos.UserDTO) == dtos.UserDTO{}) {
		return dtos.UserDTO{}, helper.Err_Data_Not_Found
	}

	return user.(dtos.UserDTO), nil
}

//----------------------------------- 	HANDLE SKROMAN CLIENT OPERATIONS ------------------------------------------- //

func (ser *user_service) CreateSkromanClient(req dtos.CreateSkromanClientRequestDTO) (dtos.SkromanClientDTO, error) {
	isValid, _ := helper.ValidateInput(req.Contact)
	if !isValid.(bool) {
		return dtos.SkromanClientDTO{}, helper.Err_Invalid_Input
	}

	args := db.CreateSkromanUserParams{
		UserName: req.UserName,
		Email:    req.Email,
		Password: sql.NullString{String: req.Password, Valid: true},
		Contact:  req.Contact,
		Address:  req.Address,
		City:     sql.NullString{String: req.City, Valid: true},
		State:    sql.NullString{String: req.State, Valid: true},
		Pincode:  sql.NullString{String: req.Pincode, Valid: true},
	}

	result, err := ser.user_repo.CreateSkromanClient(args)
	err = helper.Handle_DBError(err)

	if err != nil {
		return dtos.SkromanClientDTO{}, err
	}

	return dtos.SkromanClientDTO{
		ID:        result.ID,
		UserName:  result.UserName,
		Email:     result.Email,
		Contact:   result.Contact,
		Address:   result.Address,
		City:      result.City.String,
		State:     result.State.String,
		Pincode:   result.Pincode.String,
		CreatedAt: result.CreatedAt,
		UpdatedAt: result.UpdatedAt,
	}, nil
}

func (ser *user_service) FetchAllClients(req dtos.GetUsersRequestParams) ([]dtos.SkromanClientDTO, error) {
	var skroman_clinets []dtos.SkromanClientDTO
	err_chan := make(chan error)
	wg := sync.WaitGroup{}
	wg.Add(2)

	// fetch clients
	go func() {
		defer wg.Done()
		args := db.FetchAllClientsParams{
			Limit:  req.PageSize,
			Offset: (req.PageID - 1) * req.PageSize,
		}

		result, err := ser.user_repo.FetchAllClients(args)

		if err != nil {
			err_chan <- err
			return
		}

		if len(result) == 0 {
			err_chan <- sql.ErrNoRows
			return
		}

		t_wg := sync.WaitGroup{}
		skroman_clinets = make([]dtos.SkromanClientDTO, len(result))

		for i := range result {
			t_wg.Add(1)
			go ser.setSkromanClientData(&t_wg, &skroman_clinets[i], &result[i])
		}
		t_wg.Wait()
	}()

	// count client
	go func() {
		defer wg.Done()
		count, err := ser.user_repo.CountOfClient()
		if err != nil {
			err_chan <- err
		}
		utils.SetPaginationData(int(req.PageID), count)
	}()

	go func() {
		wg.Wait()
		close(err_chan)
	}()

	for data_err := range err_chan {
		return nil, data_err
	}

	return skroman_clinets, nil
}

func (ser *user_service) setSkromanClientData(wg *sync.WaitGroup, result *dtos.SkromanClientDTO, data *db.SkromanClient) {
	defer wg.Done()

	*result = dtos.SkromanClientDTO{
		ID:        data.ID,
		UserName:  data.UserName,
		Email:     data.Email,
		Contact:   data.Contact,
		Address:   data.Address,
		City:      data.City.String,
		State:     data.State.String,
		Pincode:   data.Pincode.String,
		CreatedAt: data.CreatedAt,
		UpdatedAt: data.UpdatedAt,
	}
}

func (ser *user_service) DeleteClient(client_id string) error {
	client_obj_id, err := uuid.Parse(client_id)

	if err != nil {
		return helper.ERR_INVALID_ID
	}

	result, err := ser.user_repo.DeleteClient(client_obj_id)

	if err != nil {
		return err
	}

	if a_r, _ := result.RowsAffected(); a_r == 0 {
		return helper.Err_Delete_Failed
	}

	return nil
}

func (ser *user_service) FetchClientById(client_id string) (dtos.SkromanClientDTO, error) {
	client_obj_id, err := uuid.Parse(client_id)

	if err != nil {
		return dtos.SkromanClientDTO{}, helper.ERR_INVALID_ID
	}

	result, err := ser.user_repo.FetchClientById(client_obj_id)

	if err != nil {
		return dtos.SkromanClientDTO{}, err
	}

	return dtos.SkromanClientDTO{
		ID:        result.ID,
		UserName:  result.UserName,
		Email:     result.Email,
		Contact:   result.Contact,
		Address:   result.Address,
		City:      result.City.String,
		State:     result.State.String,
		Pincode:   result.Pincode.String,
		CreatedAt: result.CreatedAt,
		UpdatedAt: result.UpdatedAt,
	}, nil

}

func (ser *user_service) CountOFClients() (int64, error) {
	return ser.user_repo.CountOfClient()
}
