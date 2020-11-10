package main

import (
	"logging_service/core"
	"logging_service/routes"
	"os"

	"github.com/gin-gonic/gin"
)

func main() {
	os.Setenv("TZ", "UTC")
	router := gin.New()
	var pool core.FileMutexPool
	var counters core.LogTypeCounter
	routes.Setup(router, &pool, &counters)
}
