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
	"net/http"
	"os"

	"github.com/auth0-community/go-auth0"
	"github.com/gin-gonic/gin"
	"gopkg.in/square/go-jose.v2"
)

// AuthenticateJWT is a gin middleware that authenticates a jwt in the Authorization header before proceeding with processing a request.
//
// Returns
//	gin.HandlerFunc	- next gin handler/middleware.
func AuthenticateJWT() gin.HandlerFunc {
	auth0URI := os.Getenv("AUTH0_URI")
	auth0Audience := os.Getenv("AUTH0_AUDIENCE")

	return func(c *gin.Context) {

		client := auth0.NewJWKClient(auth0.JWKClientOptions{URI: auth0URI + ".well-known/jwks.json"}, nil)
		configuration := auth0.NewConfiguration(client, []string{auth0Audience}, auth0URI, jose.RS256)
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

// WriteLog writes a log to a logfile.
//
// Receiver:
//	*LogModel				logModel
//
// Parameters:
//	int				statusCode	- http status code.
//	string			message		- response message.
//	*gin.Context	c			- handler context from gin.
//
func terminateWithError(statusCode int, message string, c *gin.Context) {
	c.JSON(statusCode, gin.H{"error": message})
	c.Abort()
}
