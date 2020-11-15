package handlers

import (
	"logging_service/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

func HandleGetAccessControl(c *gin.Context) {

	var changes = c.Request.URL.Query()
	var permissions models.AccessControlModel
	if err := permissions.GetPermissions(); err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}

	if len(changes) > 0 {
		if err := permissions.GetChanges(changes); err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
			return
		}

		if err := permissions.WritePermissions(); err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
			return
		}
	}

	permissions.Lock.RLock()
	defer permissions.Lock.RUnlock()
	c.HTML(http.StatusOK, "access_control", permissions)
}
