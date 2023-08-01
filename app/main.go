package main

import (
	"humidity_service/main/server"
	"log"

	"github.com/joho/godotenv"
)

// @title Humidity service API
// @version 1.0
// @description This is a simple server to get and store data from multiple sensor stations.
// @host localhost:8080
// @BasePath /api/v1

// @contact.name XXXX
// @contact.url http://www.XXXX.com
// @contact.email xxxx@xxxx.icom

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html
func main() {
	err := godotenv.Load(".env")

	if err != nil {
		log.Fatalf("Error loading .env file")
	}

	server.Init()
}
