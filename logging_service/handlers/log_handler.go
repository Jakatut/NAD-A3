package handlers

import (
	"net/http"
)

// Gets all logs
func HandleGetLog(c *Context) {
	c.HTML(http.StatusOK, "index.tmpl.html", nil)
}