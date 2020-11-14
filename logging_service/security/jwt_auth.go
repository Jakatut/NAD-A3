package security

import (
	"log"
	"net/http"

	"github.com/auth0-community/go-auth0"
	"github.com/gin-gonic/gin"
	"gopkg.in/square/go-jose.v2"
)

func AuthenticateJWT() gin.HandlerFunc {
	return func(c *gin.Context) {

		client := auth0.NewJWKClient(auth0.JWKClientOptions{URI: "https://ca-logging.us.auth0.com/" + ".well-known/jwks.json"}, nil)
		configuration := auth0.NewConfiguration(client, []string{"http://localhost:8000"}, "https://ca-logging.us.auth0.com/", jose.RS256)
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
