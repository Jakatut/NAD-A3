package handlers

import (
	"net/http"
	"github.com/gin-gonic/gin"
)

func handleGetLog(c *gin.Context) {
	c.HTML(http.StatusOk, "index.tmpl.html", nil)
}