package routers

import (
	"github.com/aniket-skroman/skroman-user-service/apis"
	"github.com/aniket-skroman/skroman-user-service/apis/controller"
	"github.com/aniket-skroman/skroman-user-service/apis/middleware"
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

	user := router.Group("/api", middleware.AuthorizeJWT(jwt_service))
	{
		user.PUT("/update-user", user_cont.UpdateUser)
		user.GET("/fetch-users/:page_id/:page_size", user_cont.FetchAllUsers)
		user.DELETE("/delete-user/:user_id", user_cont.DeleteUser)
		user.GET("/fetch-user", user_cont.FetchUserById)
	}

	token_val := router.Group("/api", middleware.AuthorizeJWT(jwt_service))
	{
		token_val.GET("/validate-token", user_cont.FetchUserById)
	}
}
