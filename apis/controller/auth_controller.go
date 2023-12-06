package controller

import (
	"fmt"
	"net/http"

	"github.com/aniket-skroman/skroman-user-service/apis/services"
	"github.com/aniket-skroman/skroman-user-service/utils"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
)

type AuthController interface {
	RefreshToken(*gin.Context)
}

type auth_cont struct {
	jwt_service services.JWTService
	response    map[string]interface{}
}

func NewAuthController(jwt_service services.JWTService) AuthController {
	return &auth_cont{
		jwt_service: jwt_service,
		response:    map[string]interface{}{},
	}
}

func (cont *auth_cont) RefreshToken(ctx *gin.Context) {
	authHeader := ctx.GetHeader("Authorization")

	if authHeader == "" {
		cont.response = utils.BuildFailedResponse(utils.REQUIRED_PARAMS)
		ctx.JSON(http.StatusBadRequest, cont.response)
		return
	}

	token, _ := cont.jwt_service.ValidateToken(authHeader)
	claims := token.Claims.(jwt.MapClaims)
	user_id := fmt.Sprintf("%v", claims["user_id"])
	user_type := fmt.Sprintf("%v", claims["user_type"])
	dept := fmt.Sprintf("%v", claims["dept"])

	n_token := cont.jwt_service.GenerateToken(user_id, user_type, dept)

	cont.response = utils.BuildSuccessResponse("Token has been refreshed", "access_token", n_token)
	ctx.JSON(http.StatusOK, cont.response)
}
