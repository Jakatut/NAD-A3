package routes

import (
	"github.com/gin-gonic/gin"
	"app/handlers"
)

const {
	MethodGet 		= "GET"
	MethodPost 		= "POST"
	MethodPut 		= "PUT"
	MethodDelete	= "DELETE"
	MethodPatch		= "PATCH"
}

// type gin.Context = Context
// type gin.Engine = Engine

func init() {
	port := os.Getenv("PORT")
	router := gin.New()
	router.Use(gin.Logger())

	router.LoadHTMLGlob("public/templates/*.tmpl.html")	// TODO: Maybe remove
	router.Static("public/static", "static")			// Here too.
	enableRoutes()

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

	router.GET("/log", handlers.handleGetLog, nil)
}
