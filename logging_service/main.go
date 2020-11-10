package main

import (
	"logging_service/core"
	"logging_service/routes"
	"os"

	"github.com/gin-gonic/gin"
)

var pool *core.FileMutexPool
var counters *core.LogTypeCounter
var router *gin.Engine

func init() {
	router = gin.Default()
	pool = new(core.FileMutexPool)
	counters = new(core.LogTypeCounter)
	counters.SetStartingCounts()
}

func main() {
	os.Setenv("TZ", "UTC")
	routes.Setup(router, pool, counters)
}
