package routes

import (
	log "logging_service/core"
	"logging_service/handlers"
	"os"

	"github.com/gin-gonic/gin"
)

type Engine = gin.Engine

// Setups routes
func Setup(router *Engine, logCounts *log.LogCounts) {
	port := os.Getenv("PORT")
	router.Use(gin.Logger())

	router.LoadHTMLGlob("public/templates/*.tmpl.html")
	router.Static("public/static", "static")
	enableRoutes(router, logCounts)

	// router.Run(":8080")
	router.Run(":" + port)
}

// Add your route here with:
// router.GET("routeName", func(c *gin.Context){})
// Where GET is any of the HTTP methods.
// Define parameters in the route name like: /user/:userid/status
// and get the value with c.Param("userid") in the callback.
// Callbacks should be defined under logging_service/app/handlers
func enableRoutes(router *Engine, logCounts *log.LogCounts) {

	if router != nil {
		// Root //
		router.GET("/", handlers.HandleGetRoot)

		// Log Types //
		// router.GET("/log/debug", handlers.HandleGetDebugLog, nil)
		// router.GET("/log/info", handlers.HandleGetInfoLog, nil)
		// router.GET("/log/warn", handlers.HandleGetWarnLog, nil)
		// router.GET("/log/error", handlers.HandleGetErrorLog, nil)
		// router.GET("/log/fatal", handlers.HandleGetFatalLog, nil)

		router.GET("/log", handlers.HandleGetLog) // ALL
		router.POST("/log", handlers.HandlePostLog)
	}
}
