package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"

	"github.com/aniket-skroman/skroman-user-service/apis"
	"github.com/aniket-skroman/skroman-user-service/apis/routers"
	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
)

var (
	dbDriver = ""
	dbSource = ""
	address  = ""
)

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
	fmt.Println("connecting to db")
	db, err := sql.Open(dbDriver, dbSource)

	if err != nil {
		log.Fatal(err)
	}

	//fmt.Println("connection has been established..", db)

	router := gin.Default()
	store := apis.NewStore(db)

	router.GET("/", func(ctx *gin.Context) {
		ctx.Data(http.StatusOK, ContentTypeHTML, []byte("<html>Program file run...</html>"))
	})

	routers.UserRouters(router, store)

	if err := router.Run(address); err != nil {
		log.Fatal(err)
	}
}
