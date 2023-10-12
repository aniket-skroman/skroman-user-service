package middleware

import (
	"net/http"

	"github.com/aniket-skroman/skroman-user-service/apis/services"
	"github.com/aniket-skroman/skroman-user-service/utils"
	"github.com/gin-gonic/gin"
)

func AuthorizeJWT(jwtService services.JWTService) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		authHeader := ctx.GetHeader("Authorization")

		if authHeader == "" {
			response := utils.BuildFailedResponse("Token not found")
			ctx.AbortWithStatusJSON(http.StatusBadRequest, response)
			return
		}

		_, err := jwtService.ValidateToken(authHeader)

		if err != nil {
			response := utils.BuildFailedResponse("Invalid token provided !")
			ctx.AbortWithStatusJSON(http.StatusBadRequest, response)
			return
		}

		// if token.Valid {
		// 	//claims := token.Claims.(jwt.MapClaims)
		// 	//userId := fmt.Sprintf("%v", claims["user_id"])
		// 	//userType := fmt.Sprintf("%v", claims["user_type"])

		// }

	}
}
