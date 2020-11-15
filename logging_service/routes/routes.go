package routes

import (
	"logging_service/core"
	"logging_service/handlers"
	"logging_service/security"
	"os"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

// Setup create routes with their handlers.
func Setup(router *gin.Engine, mutexPool *core.FileMutexPool, counters *core.LogTypeCounter) {
	router.Use(
		gin.Logger(),
		cors.New(cors.Config{
			AllowMethods:     []string{"POST", "GET"},
			AllowHeaders:     []string{"Content-Type", "Origin", "Accept", "Authorization", "*"},
			ExposeHeaders:    []string{"Content-Length"},
			AllowCredentials: true,
			AllowOriginFunc: func(origin string) bool {
				return true
			},
		}),
	)

	router.LoadHTMLGlob("public/templates/*.tmpl.html")
	router.Static("/static", "public/static")

	authorized := router.Group("/log")
	authorized.Use(security.AuthenticateJWT())
	authorized.GET("/log/:log_level", func(c *gin.Context) {
		handlers.HandleGetLog(c, mutexPool)
	})
	authorized.POST("/log/:log_level", func(c *gin.Context) {
		handlers.HandlePostLog(c, mutexPool, counters)
	})

	router.GET("/", handlers.HandleGetRoot)

	port := os.Getenv("PORT")
	router.Run(":" + port)
}
