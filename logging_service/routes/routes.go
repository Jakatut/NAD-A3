package routes

import (
	"os"
	"github.com/gin-gonic/gin"
	"logging_service/handlers"
)

type Engine = gin.Engine

func Setup() *Engine {
	port := os.Getenv("PORT")
	router := gin.New()
	router.Use(gin.Logger())

	router.LoadHTMLGlob("public/templates/*.tmpl.html")
	router.Static("public/static", "static")
	enableRoutes(router)

	router.Run(":" + port)
	return router
}

// Add your route here with:
// router.GET("routeName", func(c *gin.Context){})
// Where GET is any of the HTTP methods.
// Define parameters in the route name like: /user/:userid/status
// and get the value with c.Param("userid") in the callback.
// Callbacks should be defined under logging_service/app/handlers
func enableRoutes(router *Engine) {
	
	// root
	router.GET("/", handlers.HandleGetRoot, nil)

	

	// Log Types //
	// router.GET("/log/debug", handlers.HandleGetDebugLog, nil)
	// router.GET("/log/info", handlers.HandleGetInfoLog, nil)
	// router.GET("/log/warn", handlers.HandleGetWarnLog, nil)
	// router.GET("/log/error", handlers.HandleGetErrorLog, nil)
	// router.GET("/log/fatal", handlers.HandleGetFatalLog, nil)
	router.GET("/log", handlers.HandleGetLog, nil)				// ALL
}
