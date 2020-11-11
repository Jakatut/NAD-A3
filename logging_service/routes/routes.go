package routes

import (
	"logging_service/core"
	"logging_service/handlers"

	"logging_service/security"

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
	router.Static("public/static", "static")
	router.GET("/log/:log_level", security.CheckJWT(), func(c *gin.Context) {
		handlers.HandleGetLog(c, mutexPool)
	})
	router.POST("/log/:log_level", security.CheckJWT(), func(c *gin.Context) {

		handlers.HandlePostLog(c, mutexPool, counters)
	})

	// port := os.Getenv("PORT")
	router.Run(":8080")
}
