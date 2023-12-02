package controller

import (
	"database/sql"
	"fmt"
	"net/http"
	"strings"

	"github.com/aniket-skroman/skroman-user-service/apis/dtos"
	"github.com/aniket-skroman/skroman-user-service/apis/helper"
	"github.com/aniket-skroman/skroman-user-service/apis/services"
	"github.com/aniket-skroman/skroman-user-service/utils"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type UserController interface {
	CreateNewUser(*gin.Context)
	LoginUser(*gin.Context)
	UpdateUser(*gin.Context)
	FetchAllUsers(*gin.Context)
	DeleteUser(*gin.Context)
	FetchUserById(*gin.Context)
	CountEmployee(ctx *gin.Context)
	CreateSkromanClient(ctx *gin.Context)
	FetchAllClients(ctx *gin.Context)
	DeleteClient(ctx *gin.Context)
	CountOFClients(ctx *gin.Context)
	FetchClientById(ctx *gin.Context)
}

type user_controller struct {
	user_ser services.UserService
	response map[string]interface{}
}

func NewUserController(user_serv services.UserService) UserController {
	return &user_controller{
		user_ser: user_serv,
		response: map[string]interface{}{},
	}
}

func (cont *user_controller) CreateNewUser(ctx *gin.Context) {
	var req dtos.CreateUserRequestDTO

	if err := ctx.ShouldBindJSON(&req); err != nil {
		cont.response = utils.BuildFailedResponse(helper.Handle_required_param_error(err))
		ctx.JSON(http.StatusBadRequest, cont.response)
		return
	}

	user, err := cont.user_ser.CreateNewUser(req)

	if err != nil {
		cont.response = utils.BuildFailedResponse(err.Error())
		if err == helper.Err_Account_Already_Exists {
			ctx.JSON(http.StatusConflict, cont.response)
			return
		}

		ctx.JSON(http.StatusInternalServerError, cont.response)
		return
	}

	cont.response = utils.BuildSuccessResponse(utils.USER_REGISTRATION_SUCCESS, utils.USER_DATA, user)
	ctx.JSON(http.StatusCreated, cont.response)
}

func (cont *user_controller) LoginUser(ctx *gin.Context) {
	var req dtos.LoginUserRequestDTO

	if err := ctx.ShouldBindQuery(&req); err != nil {
		cont.response = utils.BuildFailedResponse(helper.Handle_required_param_error(err))
		ctx.JSON(http.StatusBadRequest, cont.response)
		return
	}

	user, err := cont.user_ser.FetchUserByEmail(req)

	if err != nil {

		if err == sql.ErrNoRows {
			cont.response = utils.BuildFailedResponse("account does not exists")
			ctx.JSON(http.StatusNotFound, cont.response)
			return
		} else if strings.Contains(err.Error(), "does not matched") {
			cont.response = utils.BuildFailedResponse(err.Error())
			ctx.JSON(http.StatusUnauthorized, cont.response)
			return
		}
		cont.response = utils.BuildFailedResponse(err.Error())
		ctx.JSON(http.StatusInternalServerError, cont.response)
		return
	}

	cont.response = utils.BuildSuccessResponse(utils.LOGIN_SUCCESS, utils.USER_DATA, user)
	ctx.JSON(http.StatusOK, cont.response)
}

func (cont *user_controller) UpdateUser(ctx *gin.Context) {
	var req dtos.UpdateUserRequestDTO

	if err := ctx.ShouldBindJSON(&req); err != nil {
		cont.response = utils.BuildFailedResponse(helper.Handle_required_param_error(err))
		ctx.JSON(http.StatusBadRequest, cont.response)
		return
	}

	user, err := cont.user_ser.UpdateUser(req)

	if err != nil {
		fmt.Println("Error : ", err)
		cont.response = utils.BuildFailedResponse(err.Error())
		if strings.Contains(err.Error(), "already used by someone") {
			ctx.JSON(http.StatusConflict, cont.response)
			return
		} else if err == sql.ErrNoRows {
			cont.response = utils.BuildFailedResponse(helper.Err_Data_Not_Found.Error())
			ctx.JSON(http.StatusNotFound, cont.response)
			return
		} else if err == helper.ERR_INVALID_ID {
			ctx.JSON(http.StatusBadRequest, cont.response)
			return
		} else if err == helper.Err_EMP_Code_Exists {
			ctx.JSON(http.StatusConflict, cont.response)
			return
		} else {
			ctx.JSON(http.StatusInternalServerError, cont.response)
			return
		}
	}

	cont.response = utils.BuildSuccessResponse(utils.UPDATE_SUCCESS, utils.USER_DATA, user)
	ctx.JSON(http.StatusOK, cont.response)
}

func (cont *user_controller) FetchAllUsers(ctx *gin.Context) {
	var req dtos.GetUsersRequestParams

	if err := ctx.ShouldBindUri(&req); err != nil {
		cont.response = utils.BuildFailedResponse(helper.Handle_required_param_error(err))
		ctx.JSON(http.StatusBadRequest, cont.response)
		return
	}

	users, err := cont.user_ser.FetchAllUsers(req)

	if err != nil {
		cont.response = utils.BuildFailedResponse(err.Error())
		if err == helper.Err_Data_Not_Found || err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, cont.response)
			return
		}

		ctx.JSON(http.StatusInternalServerError, cont.response)
		return
	}

	cont.response = utils.BuildResponseWithPagination(utils.FETCHED_SUCCESS, "", utils.USER_DATA, users)
	ctx.JSON(http.StatusOK, cont.response)
}

func (cont *user_controller) DeleteUser(ctx *gin.Context) {
	var req dtos.DeleteUserRequestDTO

	if err := ctx.ShouldBindUri(&req); err != nil {
		response := utils.BuildFailedResponse(helper.Handle_required_param_error(err))
		ctx.JSON(http.StatusBadRequest, response)
		return
	}

	err := cont.user_ser.DeleteUser(req)

	if err != nil {
		cont.response = utils.BuildFailedResponse(err.Error())

		if err == helper.ERR_INVALID_ID {
			ctx.JSON(http.StatusBadRequest, cont.response)
			return
		} else if err == helper.Err_Delete_Failed {
			ctx.JSON(http.StatusConflict, cont.response)
			return
		}

		ctx.JSON(http.StatusInternalServerError, cont.response)
		return
	}

	response := utils.BuildSuccessResponse(utils.DELETE_SUCCESS, utils.USER_DATA, utils.EmptyObj{})
	ctx.JSON(http.StatusOK, response)
}

func (cont *user_controller) FetchUserById(ctx *gin.Context) {
	if utils.TOKEN_ID == "" {
		response := utils.BuildFailedResponse("faild to extract token info")
		ctx.JSON(http.StatusUnauthorized, response)
		return
	}

	user_id, err := uuid.Parse(utils.TOKEN_ID)

	if err != nil {
		cont.response = utils.BuildFailedResponse(helper.ERR_INVALID_ID.Error())
		ctx.JSON(http.StatusBadRequest, cont.response)
		return
	}

	user, err := cont.user_ser.FetchUserById(user_id)

	if err != nil {
		cont.response = utils.BuildFailedResponse(err.Error())

		if err == helper.Err_Data_Not_Found {
			ctx.JSON(http.StatusNotFound, cont.response)
			return
		}

		ctx.JSON(http.StatusBadRequest, cont.response)
		return
	}

	response := utils.BuildSuccessResponse(utils.FETCHED_SUCCESS, utils.USER_DATA, user)
	ctx.JSON(http.StatusOK, response)
}

func (cont *user_controller) CountEmployee(ctx *gin.Context) {
	count := cont.user_ser.GetUsersCount()
	response := utils.BuildSuccessResponse(utils.FETCHED_SUCCESS, "data_count", count)
	ctx.JSON(http.StatusOK, response)
}

func (cont *user_controller) CreateSkromanClient(ctx *gin.Context) {
	var req dtos.CreateSkromanClientRequestDTO

	if err := ctx.ShouldBindJSON(&req); err != nil {
		cont.response = utils.BuildFailedResponse(helper.Handle_required_param_error(err))
		ctx.JSON(http.StatusBadRequest, cont.response)
		return
	}

	client, err := cont.user_ser.CreateSkromanClient(req)

	if err != nil {
		cont.response = utils.BuildFailedResponse(err.Error())
		if err == helper.Err_Invalid_Input {
			ctx.JSON(http.StatusBadRequest, cont.response)
			return
		} else if err == helper.Err_Account_Already_Exists {
			ctx.JSON(http.StatusForbidden, cont.response)
			return
		}
		ctx.JSON(http.StatusInternalServerError, cont.response)
		return
	}

	cont.response = utils.BuildSuccessResponse(utils.DATA_INSERTED, utils.USER_DATA, client)
	ctx.JSON(http.StatusCreated, cont.response)
}

func (cont *user_controller) FetchAllClients(ctx *gin.Context) {
	var req dtos.GetUsersRequestParams

	if err := ctx.ShouldBindUri(&req); err != nil {
		cont.response = utils.BuildFailedResponse(helper.Handle_required_param_error(err))
		ctx.JSON(http.StatusBadRequest, cont.response)
		return
	}

	clients, err := cont.user_ser.FetchAllClients(req)

	if err != nil {
		if err == sql.ErrNoRows {
			cont.response = utils.BuildFailedResponse(helper.Err_Data_Not_Found.Error())
			ctx.JSON(http.StatusNotFound, cont.response)
			return
		}

		cont.response = utils.BuildFailedResponse(err.Error())
		ctx.JSON(http.StatusInternalServerError, cont.response)
		return
	}

	cont.response = utils.BuildResponseWithPagination(utils.FETCHED_SUCCESS, "", utils.USER_DATA, clients)
	ctx.JSON(http.StatusOK, cont.response)
}

func (cont *user_controller) DeleteClient(ctx *gin.Context) {
	client_id := ctx.Param("client_id")

	if client_id == "" {
		cont.response = utils.RequestParamsMissingResponse(helper.ERR_REQUIRED_PARAMS)
		ctx.JSON(http.StatusBadRequest, cont.response)
		return
	}

	err := cont.user_ser.DeleteClient(client_id)

	if err != nil {
		cont.response = utils.BuildFailedResponse(err.Error())
		if err == helper.ERR_INVALID_ID {
			ctx.JSON(http.StatusBadRequest, cont.response)
			return
		} else if err == helper.Err_Delete_Failed {
			ctx.JSON(http.StatusNotFound, cont.response)
			return
		}

		ctx.JSON(http.StatusInternalServerError, cont.response)
		return
	}
	cont.response = utils.BuildSuccessResponse(utils.DELETE_SUCCESS, utils.USER_DATA, nil)
	ctx.JSON(http.StatusOK, cont.response)
}

func (cont *user_controller) CountOFClients(ctx *gin.Context) {
	result, err := cont.user_ser.CountOFClients()

	if err != nil {
		cont.response = utils.BuildFailedResponse(err.Error())
		ctx.JSON(http.StatusInternalServerError, cont.response)
		return
	}

	cont.response = utils.BuildSuccessResponse(utils.FETCHED_SUCCESS, "client_count", result)
	ctx.JSON(http.StatusOK, cont.response)
}

func (cont *user_controller) FetchClientById(ctx *gin.Context) {
	client_id := ctx.Param("client_id")

	if client_id == "" {
		cont.response = utils.RequestParamsMissingResponse(helper.ERR_REQUIRED_PARAMS)
		ctx.JSON(http.StatusBadRequest, cont.response)
		return
	}

	result, err := cont.user_ser.FetchClientById(client_id)

	if err != nil {
		cont.response = utils.BuildFailedResponse(err.Error())
		if err == helper.ERR_INVALID_ID {
			ctx.JSON(http.StatusBadRequest, cont.response)
			return
		} else if err == sql.ErrNoRows {
			cont.response = utils.BuildFailedResponse(helper.Err_Data_Not_Found.Error())
			ctx.JSON(http.StatusNotFound, cont.response)
			return
		}
		ctx.JSON(http.StatusInternalServerError, cont.response)
		return
	}

	cont.response = utils.BuildSuccessResponse(utils.FETCHED_SUCCESS, utils.USER_DATA, result)
	ctx.JSON(http.StatusOK, cont.response)
}
