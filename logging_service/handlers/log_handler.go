package handlers

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

type Context = gin.Context

// handlers must start with an upper case character.
// Symbols with Capital letters at the start are public. Lowercased is private.
func HandleGetLog(c *Context) {
	c.HTML(http.StatusOK, "index.tmpl.html", nil)
}