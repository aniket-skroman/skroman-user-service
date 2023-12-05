package controller

import (
	"fmt"
	"net/http"
	"strings"

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

	token, err := cont.jwt_service.ValidateToken(authHeader)
	var n_token string
	if err != nil {
		if strings.Contains(err.Error(), "expired") {
			// will generate new token
			claims := token.Claims.(jwt.MapClaims)
			user_id := fmt.Sprintf("%v", claims["user_id"])
			user_type := fmt.Sprintf("%v", claims["user_type"])
			dept := fmt.Sprintf("%v", claims["dept"])

			n_token = cont.jwt_service.GenerateToken(user_id, user_type, dept)
		}
		response := utils.BuildFailedResponse("Failed to refresh token !")
		ctx.AbortWithStatusJSON(http.StatusBadRequest, response)
		return
	}

	cont.response = utils.BuildSuccessResponse("Token has been refreshed", "access_token", n_token)
	ctx.JSON(http.StatusOK, cont.response)
}
