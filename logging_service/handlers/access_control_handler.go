package handlers

import (
	"logging_service/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

func HandleGetAccessControl(c *gin.Context) {

	var changes = c.Request.URL.Query()
	changes.Get()

	if len(changes) == 0 {
		c.AbortWithStatus(http.StatusOK)
		return
	}

	var permissions models.AccessControlModel
	if err := permissions.GetPermissions(); err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}

	if err := permissions.GetChanges(changes); err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}

	permissions.Lock.RLock()
	defer permissions.Lock.RUnlock()
	c.HTML(http.StatusOK, "access_control", permissions)
}
