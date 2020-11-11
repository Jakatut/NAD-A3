package handlers

import (
	"log"
	"logging_service/core"
	models "logging_service/models"
	"net/http"
	"net/url"
	"strings"

	"github.com/gin-gonic/gin"
)

// HandlePostLog handles all post requests for any log type.
func HandlePostLog(c *gin.Context, mutexPool *core.FileMutexPool, counters *core.LogTypeCounter) {
	logData, err := serializeLogFromRequest(c)
	if err != nil || logData == nil {
		return
	}

	result, err := postRequestWorker(logData, mutexPool, counters)

	if err != nil || result == nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		log.Fatal(err)
	} else {
		c.JSON(200, result)
	}
}

// HandleGetLog handles all get requests for any log type.
func HandleGetLog(c *gin.Context, mutexPool *core.FileMutexPool) {
	logData, err := serializeLogFromRequest(c)
	if err != nil || logData == nil {
		return
	}

	result, err := getRequestWorker(logData, mutexPool)

	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		log.Fatal(err)
	} else if result != nil && len(result.Data.([]models.LogModel)) < 1 {
		c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"error": "not found"})
	} else {
		c.JSON(200, result)
	}
}

/*
 *
 * Workers
 *
 */

func getRequestWorker(log *models.LogModel, mutexPool *core.FileMutexPool) (*core.Response, error) {
	readResult, err := log.ReadLog(mutexPool)
	if err != nil {
		return nil, err
	}

	var response = new(core.Response)
	response.Data = readResult
	return response, nil
}

func postRequestWorker(logModel *models.LogModel, mutexPool *core.FileMutexPool, counters *core.LogTypeCounter) (*core.Response, error) {

	logModel.ID = counters.AddCount(logModel.LogLevel)
	if err := logModel.WriteLog(mutexPool); err != nil {
		return nil, err
	}

	response := new(core.Response)
	response.Data = logModel
	response.Message = "success"
	return response, nil
}

/*
 *
 * Helpers
 *
 */

func serializeLogFromRequest(c *gin.Context) (*models.LogModel, error) {
	method := c.Request.Method
	logData := new(models.LogModel)

	logLevel := strings.ToUpper(c.Param("log_level"))
	if !core.IsValidLogLevel(logLevel) {
		c.AbortWithStatusJSON(http.StatusBadRequest, "invalid log level.")
		return nil, nil
	}
	logData.LogLevel = logLevel

	switch method {
	case "POST":

		if logLevel == "ALL" {
			c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"Message": "Not found."})
			return nil, nil
		}
		if err := c.ShouldBindJSON(logData); err != nil {
			if err.Error() == "EOF" {
				c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"Errors": "Missing payload"})
			}
			c.AbortWithStatusJSON(http.StatusBadRequest, err)
			return nil, nil
		}
	case "GET":
		if err := c.ShouldBindQuery(logData); err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, err)
			return nil, nil
		}
		logData.Location, _ = url.QueryUnescape(logData.Location)
		logData.Message, _ = url.QueryUnescape(logData.Message)
	}

	// checkErrors(c)

	missing, empty := logData.IsEmptyCreate()
	if empty {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"Errors": missing})
		return nil, nil
	}

	return logData, nil
}
