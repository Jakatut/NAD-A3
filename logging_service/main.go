package main

import (
	"logging_service/core"
	"logging_service/routes"

	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.New()
	var pool core.FileMutexPool
	routes.Setup(router, &pool)
}
