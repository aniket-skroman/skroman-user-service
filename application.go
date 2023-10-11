package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	"github.com/aniket-skroman/skroman-user-service/apis"
	"github.com/aniket-skroman/skroman-user-service/apis/routers"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

var (
	dbDriver = ""
	dbSource = ""
	address  = ""
)

func init() {
	if err := godotenv.Load(".env"); err != nil {
		log.Fatal(err)
	}

	dbDriver = os.Getenv("DB_DRIVER")
	dbSource = os.Getenv("SKROMAN_DB")
	address = os.Getenv("LOCAL_ADDRESS")
}

func CORSConfig() cors.Config {
	corsConfig := cors.DefaultConfig()
	corsConfig.AllowOrigins = []string{"http://localhost:3000"}
	corsConfig.AllowCredentials = true
	corsConfig.AddAllowHeaders("Access-Control-Allow-Headers", "access-control-allow-origin, access-control-allow-headers", "Content-Type", "X-XSRF-TOKEN", "Accept", "Origin", "X-Requested-With", "Authorization")
	corsConfig.AddAllowMethods("GET", "POST", "PUT", "DELETE")
	return corsConfig
}
func main() {
	fmt.Println("connecting to db")
	db, err := sql.Open(dbDriver, dbSource)

	if err != nil {
		log.Fatal(err)
	}

	//fmt.Println("connection has been established..", db)

	router := gin.Default()
	store := apis.NewStore(db)

	routers.UserRouters(router, store)

	if err := router.Run(address); err != nil {
		log.Fatal(err)
	}
}
