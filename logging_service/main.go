package main

import (
	log "logging_service/core"
	"logging_service/routes"

	"github.com/gin-gonic/gin"
)

func main() {

	logCounts := new(log.LogCounts)

	router := gin.New()
	routes.Setup(router, logCounts)
}
