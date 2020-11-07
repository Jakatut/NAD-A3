package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// HandleGetRoot serves the html index.
func HandleGetRoot(c *gin.Context) {
	c.HTML(http.StatusOK, "index.tmpl.html", nil)
}
