package handlers

import (
	"net/http"
)

// Get root
func HandleGetRoot(c *Context) {
	c.HTML(http.StatusOK, "index.tmpl.html", nil)
}