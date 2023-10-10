package controller

import (
	"net/http"

	"github.com/aniket-skroman/skroman-user-service/apis/dtos"
	"github.com/aniket-skroman/skroman-user-service/apis/services"
	"github.com/aniket-skroman/skroman-user-service/utils"
	"github.com/gin-gonic/gin"
)

type UserController interface {
	CreateNewUser(ctx *gin.Context)
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
		response := utils.BuildFailedResponse(utils.FAILED_PROCESS, err.Error(), utils.USER_DATA, utils.EmptyObj{})
		ctx.JSON(http.StatusBadRequest, response)
		return
	}

	user, err := cont.user_ser.CreateNewUser(req)

	if err != nil {
		response := utils.BuildFailedResponse(utils.FAILED_PROCESS, err.Error(), utils.USER_DATA, utils.EmptyObj{})
		ctx.JSON(http.StatusInternalServerError, response)
		return
	}

	response := utils.BuildSuccessResponse(utils.USER_REGISTRATION_SUCCESS, utils.USER_DATA, user)
	ctx.JSON(http.StatusOK, response)
}
