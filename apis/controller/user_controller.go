package controller

import (
	"net/http"
	"strings"
	"sync"

	"github.com/aniket-skroman/skroman-user-service/apis/dtos"
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
}

type user_controller struct {
	user_ser services.UserService
}

func NewUserController(user_serv services.UserService) UserController {
	return &user_controller{
		user_ser: user_serv,
	}
}

func (cont *user_controller) CreateNewUser(ctx *gin.Context) {
	var req dtos.CreateUserRequestDTO

	if err := ctx.ShouldBindJSON(&req); err != nil {
		response := utils.BuildFailedResponse(err.Error())
		ctx.JSON(http.StatusBadRequest, response)
		return
	}

	user, err := cont.user_ser.CreateNewUser(req)

	if err != nil {
		response := utils.BuildFailedResponse(err.Error())
		if strings.Contains(err.Error(), "already exits") {
			ctx.JSON(http.StatusConflict, response)
			return
		}

		ctx.JSON(http.StatusInternalServerError, response)
		return
	}

	response := utils.BuildSuccessResponse(utils.USER_REGISTRATION_SUCCESS, utils.USER_DATA, user)
	ctx.JSON(http.StatusCreated, response)
}

func (cont *user_controller) LoginUser(ctx *gin.Context) {
	var req dtos.LoginUserRequestDTO

	if err := ctx.ShouldBindQuery(&req); err != nil {
		response := utils.BuildFailedResponse(err.Error())
		ctx.JSON(http.StatusBadRequest, response)
		return
	}

	user, err := cont.user_ser.FetchUserByEmail(req)

	if err != nil {

		if strings.Contains(err.Error(), "no rows in result set") {
			response := utils.BuildFailedResponse("account does not exists")
			ctx.JSON(http.StatusNotFound, response)
			return
		} else if strings.Contains(err.Error(), "does not matched") {
			response := utils.BuildFailedResponse(err.Error())
			ctx.JSON(http.StatusUnauthorized, response)
			return
		}
		response := utils.BuildFailedResponse(err.Error())
		ctx.JSON(http.StatusInternalServerError, response)
		return
	}

	response := utils.BuildSuccessResponse(utils.LOGIN_SUCCESS, utils.USER_DATA, user)
	ctx.JSON(http.StatusOK, response)
}

func (cont *user_controller) UpdateUser(ctx *gin.Context) {
	var req dtos.UpdateUserRequestDTO

	if err := ctx.ShouldBindJSON(&req); err != nil {
		response := utils.BuildFailedResponse(err.Error())
		ctx.JSON(http.StatusBadRequest, response)
		return
	}

	user, err := cont.user_ser.UpdateUser(req)

	if err != nil {
		response := utils.BuildFailedResponse(err.Error())
		if strings.Contains(err.Error(), "already used by someone") {
			ctx.JSON(http.StatusConflict, response)
			return
		}
		ctx.JSON(http.StatusInternalServerError, response)
		return
	}

	response := utils.BuildSuccessResponse(utils.UPDATE_SUCCESS, utils.USER_DATA, user)
	ctx.JSON(http.StatusOK, response)
}

func (cont *user_controller) FetchAllUsers(ctx *gin.Context) {
	var req dtos.GetUsersRequestParams

	if err := ctx.ShouldBindUri(&req); err != nil {
		response := utils.BuildFailedResponse(err.Error())
		ctx.JSON(http.StatusBadRequest, response)
		return
	}

	wg := sync.WaitGroup{}
	wg.Add(1)

	go func() {
		defer wg.Done()
		count := cont.user_ser.GetUSersCount()
		utils.SetPaginationData(int(req.PageID), int64(count))
	}()

	users, err := cont.user_ser.FetchAllUsers(req)

	if err != nil {
		response := utils.BuildFailedResponse(err.Error())
		if strings.Contains(err.Error(), "data not found") {
			ctx.JSON(http.StatusNotFound, response)
			return
		}

		ctx.JSON(http.StatusInternalServerError, response)
		return
	}

	wg.Wait()

	response := utils.BuildResponseWithPagination(utils.FETCHED_SUCCESS, "", utils.USER_DATA, users)
	ctx.JSON(http.StatusOK, response)
}

func (cont *user_controller) DeleteUser(ctx *gin.Context) {
	var req dtos.DeleteUserRequestDTO

	if err := ctx.ShouldBindUri(&req); err != nil {
		response := utils.BuildFailedResponse(err.Error())
		ctx.JSON(http.StatusBadRequest, response)
		return
	}

	err := cont.user_ser.DeleteUser(req)

	if err != nil {
		response := utils.BuildFailedResponse(err.Error())
		ctx.JSON(http.StatusInternalServerError, response)
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
		response := utils.BuildFailedResponse(err.Error())
		ctx.JSON(http.StatusBadRequest, response)
		return
	}

	user, err := cont.user_ser.FetchUserById(user_id)

	if err != nil {
		response := utils.BuildFailedResponse(err.Error())

		if strings.Contains(err.Error(), "not found") {
			ctx.JSON(http.StatusNotFound, response)
			return
		}

		ctx.JSON(http.StatusBadRequest, response)
		return
	}

	response := utils.BuildSuccessResponse(utils.FETCHED_SUCCESS, utils.USER_DATA, user)
	ctx.JSON(http.StatusOK, response)
}
