package routes

import (
	"logging_service/core"
	"logging_service/handlers"
	"sync"

	"github.com/gin-gonic/gin"
)

type Engine = gin.Engine

// Setups routes
func Setup(router *Engine, debugWaitGroup *sync.WaitGroup, warningWaitGroup *sync.WaitGroup, infoWaitGroup *sync.WaitGroup, errorWaitGroup *sync.WaitGroup, fatalWaitGroup *sync.WaitGroup) {
	// port := os.Getenv("PORT")
	router.Use(gin.Logger())

	router.LoadHTMLGlob("public/templates/*.tmpl.html")
	router.Static("public/static", "static")

	resources := make(map[string]core.HandlerResources)
	resources["/log/debug"] = core.HandlerResources{WaitGroup: debugWaitGroup, LogChannel: make(chan *core.Result)}
	resources["/log/warninig"] = core.HandlerResources{WaitGroup: warningWaitGroup, LogChannel: make(chan *core.Result)}
	resources["/log/info"] = core.HandlerResources{WaitGroup: infoWaitGroup, LogChannel: make(chan *core.Result)}
	resources["/log/error"] = core.HandlerResources{WaitGroup: errorWaitGroup, LogChannel: make(chan *core.Result)}
	resources["/log/fatal"] = core.HandlerResources{WaitGroup: fatalWaitGroup, LogChannel: make(chan *core.Result)}

	for route, resource := range resources {
		// router.GET(route, func(c *gin.Context) {
		// 	handlers.HandleLog(c, resource)
		// })
		router.POST(route, func(c *gin.Context) {
			handlers.HandleLog(c, resource)
		})
	}

	router.Run(":8080")
	// router.Run(":" + port)
}
