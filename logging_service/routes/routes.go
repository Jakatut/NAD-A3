package routes

import (
	"logging_service/handlers"
	"sync"

	"github.com/gin-gonic/gin"
)

type Engine = gin.Engine

var endpoints map[string]*sync.WaitGroup

// Setups routes
func Setup(router *Engine, debugWaitGroup *sync.WaitGroup, warningWaitGroup *sync.WaitGroup, infoWaitGroup *sync.WaitGroup, errorWaitGroup *sync.WaitGroup, fatalWaitGroup *sync.WaitGroup) {
	// port := os.Getenv("PORT")
	router.Use(gin.Logger())

	router.LoadHTMLGlob("public/templates/*.tmpl.html")
	router.Static("public/static", "static")

	endpoints = make(map[string]*sync.WaitGroup)
	endpoints["/log/debug"] = debugWaitGroup
	endpoints["/log/warninig"] = warningWaitGroup
	endpoints["/log/info"] = infoWaitGroup
	endpoints["/log/error"] = errorWaitGroup
	endpoints["/log/fatal"] = fatalWaitGroup
	for route, waitGroup := range endpoints {
		router.GET(route, func(c *gin.Context) {
			handlers.HandleGetLog(c, waitGroup)
		})
		router.POST(route, func(c *gin.Context) {
			handlers.HandlePostLog(c, waitGroup)
		})
	}

	router.Run(":8080")
	// router.Run(":" + port)
}
