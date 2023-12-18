package routers

import (
	"github.com/aniket-skroman/skroman-user-service/apis"
	"github.com/aniket-skroman/skroman-user-service/apis/controller"
	"github.com/aniket-skroman/skroman-user-service/apis/middleware"
	"github.com/aniket-skroman/skroman-user-service/apis/repositories"
	"github.com/aniket-skroman/skroman-user-service/apis/services"
	"github.com/gin-gonic/gin"
)

var (
	user_repo   repositories.UserRepository
	jwt_service services.JWTService
	user_serv   services.UserService
	auth_cont   controller.AuthController
	user_cont   controller.UserController
)

func UserRouters(router *gin.Engine, store *apis.Store) {
	user_repo = repositories.NewUserRepository(store)
	jwt_service = services.NewJWTService()
	user_serv = services.NewUserService(user_repo, jwt_service)
	auth_cont = controller.NewAuthController(jwt_service)
	user_cont = controller.NewUserController(user_serv)

	user_auth := router.Group("/api")
	{
		user_auth.POST("/create-user", user_cont.CreateNewUser)
		user_auth.GET("/login-user", user_cont.LoginUser)
	}

	user := router.Group("/api", middleware.AuthorizeJWT(jwt_service))
	{
		user.PUT("/update-user", user_cont.UpdateUser)
		user.GET("/fetch-users/:page_id/:page_size/:department", user_cont.FetchAllUsers)
		user.DELETE("/delete-user/:user_id", user_cont.DeleteUser)
		user.GET("/fetch-user", user_cont.FetchUserById)
	}

	user_proxy := router.Group("/api", middleware.AuthorizeJWT(jwt_service))
	{
		user_proxy.GET("/emp-count", user_cont.CountEmployee)
	}

	token_val := router.Group("/api", middleware.AuthorizeJWT(jwt_service))
	{
		token_val.GET("/validate-token", user_cont.FetchUserById)
	}

	router.GET("/api/refresh-token", auth_cont.RefreshToken)

	skroman_client := router.Group("/api", middleware.AuthorizeJWT(jwt_service))
	{
		skroman_client.POST("/clients", user_cont.CreateSkromanClient)
		skroman_client.GET("/clients/:page_id/:page_size", user_cont.FetchAllClients)
		skroman_client.DELETE("/clients/:client_id", user_cont.DeleteClient)
		skroman_client.GET("/count-clients", user_cont.CountOFClients)
		skroman_client.GET("/client/:client_id", user_cont.FetchClientById)
		skroman_client.PUT("client/:client_id", user_cont.UpdateSkromanClientInfo)
	}

	search := router.Group("/api/search", middleware.AuthorizeJWT(jwt_service))
	{
		search.GET("/client/:page_id/:page_size/:search_data", user_cont.SearchClient)
		search.GET("/user/:page_id/:page_size/:search_data", user_cont.SearchUsers)
	}

	router.GET("/api/user/:user_id", user_cont.EMPById)
}
