package main

import (
	"logging_service/database"
	"logging_service/routes"
	"os"

	"github.com/gin-gonic/gin"
)

var router *gin.Engine

func init() {
	router = gin.Default()
	database.CreateConnectionConfig()
}

func main() {
	// Set the timezone to UTC so incoming datetimes can be compared to the log file's datetimes which do not have a timezone.
	os.Setenv("TZ", "UTC")
	routes.Setup(router)
}
