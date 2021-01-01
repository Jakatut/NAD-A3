package handlers

/*
 *
 * file: 		log_handler.go
 * project:		logging_service - NAD-A3
 * programmer: 	Conor Macpherson
 * description: Defines the log handler for handling requests from gin's router.
 *
 */

import (
	"log"
	"logging_service/models"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/kamva/mgm/v3"
)

// HandlePostLog handles all post requests for any log type.
//
// Parameters:
//	*gin.Context			c			- Handler context from gin.
//	*core.FileMutexPool		mutexPool	- Contains Read/Write mutexes for each log type.
//	*core.LogTypeCounter	counters	- Contains id counters for each log type.
//
func HandlePostLog(c *gin.Context) {
	logData, err := getNewLog(c)
	if err != nil || logData == nil {
		return
	}

	if err := logData.Create(); err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
	} else {
		c.JSON(200, logData)
	}
}

// HandleGetLog handles all get requests for any log type.
//
// Parameters:
//	*gin.Context	c	- Handler context from gin.
//
func HandleGetLog(c *gin.Context) {
	fields := models.LogSearchFields{}
	err := fields.GetSearchFields(c)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"Error": err.Error()})
	}

	ctx := mgm.Ctx()
	log := models.Log{}
	results, err := log.Find(ctx, fields)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"Error": "internal server error"})
	} else {
		c.JSON(200, results)
	}
}

// HandleGetLogCount handles getting the count of logs based on the provided parameters.
//
// Parameters:
//	*gin.Context	c	- Handler context from gin.
//
func HandleGetLogCount(c *gin.Context) {
	fields := models.LogSearchFields{}
	err := fields.GetSearchFields(c)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"Error": err.Error()})
	}

	_log := models.Log{}
	ctx := mgm.Ctx()
	countType := strings.Trim(c.Param("type"), "/")
	var count interface{}
	switch countType {
	case "date":
		count, err = _log.CountByDates(ctx, fields)
		break
	default:
		count, err = _log.Count(ctx, fields)
		break
	}

	if err != nil {
		log.Fatal(err.Error())
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"Error": "internal server error"})
	} else {
		c.JSON(200, count)
	}
}

/*
 *
 * Helpers
 *
 */

// getNewLog converts a json payload to a log model.
//
// Parameters:
//	*gin.Context	c	- Handler context from gin.
//
// Returns
//	*models.Log	- Serialized log model.
//	error		- Error that occurs or nil.
//
func getNewLog(c *gin.Context) (*models.Log, error) {
	logData := new(models.Log)

	// Check the log level.
	logLevel := strings.ToUpper(c.Param("log_level"))
	valid, _ := models.IsValidLogLevel(logLevel)
	if logLevel != "" && !valid {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"Error": "invalid log level"})
		return logData, nil
	}
	logData.LogLevel = logLevel

	if logData.LogLevel == "" {
		c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"Message": "Not found."})
		return nil, nil
	}
	if err := c.ShouldBindJSON(logData); err != nil {
		if err.Error() == "EOF" {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"Errors": "Missing payload"})
			return nil, nil
		} else if errors, empty := logData.IsEmptyCreate(); empty {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"Errors": errors})
			return nil, nil
		}
		c.AbortWithStatusJSON(http.StatusBadRequest, err)
		return nil, nil
	}

	logData.CreatedAt = time.Now()

	return logData, nil
}
