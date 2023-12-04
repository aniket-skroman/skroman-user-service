package main

import (
	"database/sql"
	"log"
	"net/http"

	"github.com/aniket-skroman/skroman-user-service/apis"
	"github.com/aniket-skroman/skroman-user-service/apis/database"
	"github.com/aniket-skroman/skroman-user-service/apis/routers"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
)

var (
	address = ":8080"
)

func CORSConfig() cors.Config {
	corsConfig := cors.DefaultConfig()
	corsConfig.AllowOrigins = []string{"http://localhost:3000", "http://3.109.133.20:3000", "http://13.234.110.115:3000"}
	corsConfig.AllowCredentials = true
	corsConfig.AddAllowHeaders("Access-Control-Allow-Headers", "access-control-allow-origin, access-control-allow-headers", "Content-Type", "X-XSRF-TOKEN", "Accept", "Origin", "X-Requested-With", "Authorization")
	corsConfig.AddAllowMethods("GET", "POST", "PUT", "DELETE")
	return corsConfig
}

const (
	ContentTypeBinary = "application/octet-stream"
	ContentTypeForm   = "application/x-www-form-urlencoded"
	ContentTypeJSON   = "application/json"
	ContentTypeHTML   = "text/html; charset=utf-8"
	ContentTypeText   = "text/plain; charset=utf-8"
)

type APIServer struct{}

func (api *APIServer) make_db_connection() (*sql.DB, error) {
	db, err := database.DB_INSTANCE()
	if err != nil {
		log.Fatal(err)
	}
	return db, nil
}

func (api *APIServer) init_app_route() *gin.Engine {
	r := gin.New()
	r.Use(cors.New(CORSConfig()))

	return r
}

func (api *APIServer) make_app_route(route *gin.Engine, db *sql.DB) {
	route.GET("/", func(ctx *gin.Context) {
		ctx.Data(http.StatusOK, ContentTypeHTML, []byte("<html>Program file run...</html>"))
	})

	store := apis.NewStore(db)

	routers.UserRouters(route, store)

}

func (api *APIServer) run_app(route *gin.Engine) error {
	return route.Run(address)
}

func main() {

	app_server := APIServer{}

	db, _ := app_server.make_db_connection()

	defer database.CloseDBConnection(db)

	route := app_server.init_app_route()

	app_server.make_app_route(route, db)

	if err := app_server.run_app(route); err != nil {
		log.Fatal(err)
	}
}
