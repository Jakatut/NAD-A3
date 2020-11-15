package security

import (
	"log"
	"net/http"
	"os"

	"github.com/auth0-community/go-auth0"
	"github.com/gin-gonic/gin"
	"gopkg.in/square/go-jose.v2"
)

// AuthenticateJWT is a gin middleware that authenticates a jwt in the Authorization header before proceeding with processing a request.
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

func terminateWithError(statusCode int, message string, c *gin.Context) {
	c.JSON(statusCode, gin.H{"error": message})
	c.Abort()
}
