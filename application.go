package main

import (
	"database/sql"
	"log"
	"net/http"

	"github.com/aniket-skroman/skroman-user-service/apis"
	"github.com/aniket-skroman/skroman-user-service/apis/routers"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
)

var (
	dbDriver = "postgres"
	dbSource = "postgresql://postgres:support12@skroman-user.ckwveljlsuux.ap-south-1.rds.amazonaws.com:5432/skroman_users"
	address  = ":8080"
)

func CORSConfig() cors.Config {
	corsConfig := cors.DefaultConfig()
	corsConfig.AllowOrigins = []string{"http://localhost:3000", "http://3.109.133.20:3000"}
	corsConfig.AllowCredentials = true
	corsConfig.AddAllowHeaders("Access-Control-Allow-Headers", "access-control-allow-origin, access-control-allow-headers", "Content-Type", "X-XSRF-TOKEN", "Accept", "Origin", "X-Requested-With", "Authorization")
	corsConfig.AddAllowMethods("GET", "POST", "PUT", "DELETE")
	return corsConfig
}

func init() {
	// if err := godotenv.Load(".env"); err != nil {
	// 	log.Fatal(err)
	// }

	// dbDriver = os.Getenv("DB_DRIVER")
	// dbSource = os.Getenv("LOCAL_DB_SOURCE")
	// address = os.Getenv("LOCAL_ADDRESS")
}

const (
	ContentTypeBinary = "application/octet-stream"
	ContentTypeForm   = "application/x-www-form-urlencoded"
	ContentTypeJSON   = "application/json"
	ContentTypeHTML   = "text/html; charset=utf-8"
	ContentTypeText   = "text/plain; charset=utf-8"
)

func main() {
	db, err := sql.Open(dbDriver, dbSource)

	if err != nil {
		log.Fatal(err)
	}

	err = db.Ping()
	if err != nil {
		log.Fatal(err)
	}

	router := gin.Default()
	router.Use(cors.New(CORSConfig()))
	router.Static("static", "static")

	store := apis.NewStore(db)

	router.GET("/", func(ctx *gin.Context) {
		ctx.Data(http.StatusOK, ContentTypeHTML, []byte("<html>Program file run...</html>"))
	})

	routers.UserRouters(router, store)

	if err := router.Run(address); err != nil {
		log.Fatal(err)
	}
}
