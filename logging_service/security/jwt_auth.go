package security

/*
 *
 * file: 		jwt_auth.go
 * project:		logging_service - NAD-A3
 * programmer: 	Conor Macpherson
 * description: Defines the middleware for authentication jwt tokens from/with Auth0.
 *
 */

import (
	"log"
	"logging_service/config"
	"net/http"

	"github.com/auth0-community/go-auth0"
	"github.com/gin-gonic/gin"
	"gopkg.in/square/go-jose.v2"
)

// AuthenticateJWT is a gin middleware that authenticates a jwt in the Authorization header before proceeding with processing a request.
//
// Returns
//	gin.HandlerFunc	- next gin handler/middleware.
func AuthenticateJWT() gin.HandlerFunc {

	return func(c *gin.Context) {

		conf := config.GetConfig()
		client := auth0.NewJWKClient(auth0.JWKClientOptions{URI: conf.Auth.Auth0URI + ".well-known/jwks.json"}, nil)
		configuration := auth0.NewConfiguration(client, []string{conf.Auth.Auth0Audience}, conf.Auth.Auth0URI, jose.RS256)
		validator := auth0.NewValidator(configuration, nil)

		_, err := validator.ValidateRequest(c.Request)

		if err != nil {
			log.Println(err)
			terminateWithError(http.StatusUnauthorized, "token is not valid", c)
			return
		}
		c.Next()
	}
}

func terminateWithError(statusCode int, message string, c *gin.Context) {
	c.JSON(statusCode, gin.H{"error": message})
	c.Abort()
}
