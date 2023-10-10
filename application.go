package main

import (
	"log"
	"os"

	"github.com/gin-contrib/cors"
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
	dbSource = os.Getenv("LOCAL_DB_SOURCE")
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

}
