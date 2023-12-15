package controller

import (
	"net/http"

	"github.com/aniket-skroman/skroman-user-service/apis/dtos"
	"github.com/aniket-skroman/skroman-user-service/apis/helper"
	"github.com/aniket-skroman/skroman-user-service/apis/services"
	"github.com/aniket-skroman/skroman-user-service/utils"
	"github.com/gin-gonic/gin"
)

type FCMController interface {
	CreateUserFCMData(ctx *gin.Context)
	FetchFCMTokensByUser(ctx *gin.Context)
}

type fcm_cont struct {
	fcm_serv services.FCMService
	response map[string]interface{}
}

func NewFCMController(serv services.FCMService) FCMController {
	return &fcm_cont{
		fcm_serv: serv,
		response: make(map[string]interface{}),
	}
}

func (cont *fcm_cont) CreateUserFCMData(ctx *gin.Context) {
	var req dtos.CreateUserFCMDataRequestDTO

	if err := ctx.ShouldBindJSON(&req); err != nil {
		cont.response = utils.BuildFailedResponse(helper.Handle_required_param_error(err))
		ctx.JSON(http.StatusBadRequest, cont.response)
		return
	}

	result, err := cont.fcm_serv.CreateUserFCMData(req)

	if err != nil {
		cont.response = utils.BuildFailedResponse(err.Error())

		if err == helper.ERR_INVALID_ID {
			ctx.JSON(http.StatusBadRequest, cont.response)
			return
		}

		ctx.JSON(http.StatusInternalServerError, cont.response)
		return
	}

	cont.response = utils.BuildSuccessResponse(utils.DATA_INSERTED, utils.USER_DATA, result)
	ctx.JSON(http.StatusCreated, cont.response)
}

func (cont *fcm_cont) FetchFCMTokensByUser(ctx *gin.Context) {
	user_id := ctx.Param("user_id")

	result, err := cont.fcm_serv.FetchTokenByUsers(user_id)

	if err != nil {
		cont.response = utils.BuildFailedResponse(err.Error())

		if err == helper.ERR_INVALID_ID {
			ctx.JSON(http.StatusBadRequest, cont.response)
			return
		} else if err == helper.Err_Data_Not_Found {
			ctx.JSON(http.StatusNotFound, cont.response)
			return
		}

		ctx.JSON(http.StatusInternalServerError, cont.response)
		return
	}

	cont.response = utils.BuildSuccessResponse(utils.FETCHED_SUCCESS, utils.USER_DATA, result)
	ctx.JSON(http.StatusOK, cont.response)
}
