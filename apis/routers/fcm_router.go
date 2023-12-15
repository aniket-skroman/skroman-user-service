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
	fcm_repo repositories.FCMRepository
	fcm_serv services.FCMService
	fcm_cont controller.FCMController
)

func FcmRouter(route *gin.Engine, db *apis.Store) {
	fcm_repo = repositories.NewFCMRepository(db)
	fcm_serv = services.NewFCMService(fcm_repo)
	fcm_cont = controller.NewFCMController(fcm_serv)

	fcm_data := route.Group("/api", middleware.AuthorizeJWT(jwt_service))
	{
		fcm_data.POST("/fcm-data", fcm_cont.CreateUserFCMData)
		fcm_data.GET("/fcm-data/:user_id", fcm_cont.FetchFCMTokensByUser)
	}
}
