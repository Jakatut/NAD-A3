package routes

import (
	"logging_service/core"
	"logging_service/handlers"

	"github.com/gin-gonic/gin"
)

// Setup create routes with their handlers.
func Setup(router *gin.Engine, mutexPool *core.FileMutexPool) {
	// port := os.Getenv("PORT")
	router.Use(gin.Logger())

	router.LoadHTMLGlob("public/templates/*.tmpl.html")
	router.Static("public/static", "static")
	resources := []string{"/log/debug", "/log/warninig", "/log/info", "/log/error", "/log/fatal"}

	for _, route := range resources {
		// router.GET(route, func(c *gin.Context) {
		// 	handlers.HandleLog(c, resource)
		// })
		router.POST(route, func(c *gin.Context) {
			handlers.HandleLog(c, mutexPool)
		})
	}

	router.Run(":8080")
	// router.Run(":" + port)
}
