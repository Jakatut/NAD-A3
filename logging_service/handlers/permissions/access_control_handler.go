package permissions

import (
	"logging_service/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

func HandleGetAccessControl(c *gin.Context) {

	var changes map[string]string
	if err := c.ShouldBind(&changes); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "internal server error"})
	}

	if len(changes) == 0 {

	}

	var permissions models.AccessControlModel
	if err := permissions.GetPermissions(); err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
	}

	// permMap := structs.Map(permissions)

	c.HTML(http.StatusOK, "access_control", permissions)
}
