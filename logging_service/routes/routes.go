package routes

import (
	"os"
	"github.com/gin-gonic/gin"
	"logging_service/handlers"
)

// type gin.Context = Context
// type gin.Engine = Engine

func Setup() *gin.Engine {
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
func enableRoutes(router *gin.Engine) {

	router.GET("/log", handlers.HandleGetLog, nil)
}
