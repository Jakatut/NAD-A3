package security

import (
	"os"

	jwtmiddleware "github.com/auth0/go-jwt-middleware"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
)

var authMiddlware = jwtmiddleware.New(jwtmiddleware.Options{
	ValidationKeyGetter: func(token *jwt.Token) (interface{}, error) {
		return []byte(os.Getenv("AUTH0_CLIENT_SECRET")), nil
	},
	SigningMethod: jwt.SigningMethodHS256,
})

// CheckJWT checks the Authorization headers contains a valid Auth0 JWT.
func CheckJWT() gin.HandlerFunc {
	return func(c *gin.Context) {
		jwtMiddleware := *authMiddlware
		if err := jwtMiddleware.CheckJWT(c.Writer, c.Request); err != nil {
			c.AbortWithStatus(401)
		}
	}
}
