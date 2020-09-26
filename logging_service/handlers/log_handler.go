package handlers

import (
	"net/http"
	"github.com/gin-gonic/gin"
)

// handlers must start with an upper case character.
// Symbols with Capital letters at the start are public. Lowercased is private.
func HandleGetLog(c *gin.Context) {
	c.HTML(http.StatusOK, "index.tmpl.html", nil)
}