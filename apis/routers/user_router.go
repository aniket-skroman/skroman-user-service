package routers

import (
	"github.com/aniket-skroman/skroman-user-service/apis"
	"github.com/aniket-skroman/skroman-user-service/apis/controller"
	"github.com/aniket-skroman/skroman-user-service/apis/repositories"
	"github.com/aniket-skroman/skroman-user-service/apis/services"
	"github.com/gin-gonic/gin"
)

func UserRouters(router *gin.Engine, store *apis.Store) {
	var (
		user_repo   = repositories.NewUserRepository(store)
		jwt_service = services.NewJWTService()
		user_serv   = services.NewUserService(user_repo, jwt_service)
		user_cont   = controller.NewUserController(user_serv)
	)

	user_auth := router.Group("/api")
	{
		user_auth.POST("/create-user", user_cont.CreateNewUser)
		user_auth.GET("/login-user", user_cont.LoginUser)
	}
}
