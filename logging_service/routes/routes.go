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

	router.GET("/log/:log_level", func(c *gin.Context) {
		handlers.HandleLog(c, mutexPool)
	})
	router.POST("/log/:log_level", func(c *gin.Context) {
		handlers.HandleLog(c, mutexPool)
	})

	router.Run(":8080")
	// router.Run(":" + port)
}
