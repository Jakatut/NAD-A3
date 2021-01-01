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
	"logging_service/handlers"
	"logging_service/security"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

// Setup configures the router and assigned routes their middleware & handlers.
//
// Parameters:
//	*gin.Engine				router		- gin router
//
func Setup(router *gin.Engine) {

	var configs = config.GetConfig()

	// Add logger, cross origin restrictions.
	router.Use(
		cors.New(cors.Config{
			AllowMethods:     []string{"POST", "GET"},
			AllowHeaders:     []string{"Content-Type", "Origin", "Accept", "Authorization", "*"},
			ExposeHeaders:    []string{"Content-Length"},
			AllowCredentials: true,
			AllowOriginFunc: func(origin string) bool {
				for _, val := range configs.Server.AllowedOrigins {
					if val == origin {
						return true
					}
				}
				return false
			},
		}),
	)

	router.LoadHTMLGlob("public/templates/*.tmpl.html")
	router.Static("/static", "public/static")

	router.Use(security.AuthenticateJWT())
	router.GET("/log", handlers.HandleGetLog)
	router.GET("/log/:log_level", handlers.HandleGetLog)
	router.POST("/log/:log_level", handlers.HandlePostLog)
	router.GET("/log/:log_level/count/*type", handlers.HandleGetLogCount)

	router.Run(":" + configs.Server.Port)
}
