package handlers

import (
	"app/messages"
	"net/http"
	"github.com/gin-gonic/gin"
)

// type gin.Context = Context

func handleGetLog(c *gin.Context) {
	c.HTML(http.StatusOk, "index.tmpl.html", nil)
}