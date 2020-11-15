package routes

/*
 *
 * file: 		routes.go
 * project:		logging_service - NAD-A3
 * programmer: 	Conor Macpherson
 * description: Defines routes used in the logging service and initializes the logger, cors, and jwt token authentication.
 *
 */

import (
	"logging_service/config"
	"logging_service/core"
	"logging_service/handlers"
	"logging_service/security"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

// Setup configures the router and assigned routes their middleware & handlers.
//
// Parameters:
//	*gin.Engine				router		- gin router
//	*core.FileMutexPool		mutexPool	- contains Read/Write mutexes for each log type.
//	*core.LogTypeCounter	counters	- contains id counters for each log type
//
func Setup(router *gin.Engine, mutexPool *core.FileMutexPool, counters *core.LogTypeCounter) {

	// Add logger, cross origin restrictions.
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

	// Create a log group and give JWT auth middleware.
	// authorized := router.Group("/log")
	router.Use(security.AuthenticateJWT())
	router.GET("/log/:log_level", func(c *gin.Context) {
		handlers.HandleGetLog(c, mutexPool)
	})
	router.POST("/log/:log_level", func(c *gin.Context) {
		handlers.HandlePostLog(c, mutexPool, counters)
	})

	router.Run(":" + config.GetConfig().Port)
}
