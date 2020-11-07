package handlers

import (
	"fmt"
	"logging_service/core"
	models "logging_service/models"
	"net/url"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

// HandleLog handles all post/get requests for any log type.
func HandleLog(c *gin.Context, mutexPool *core.FileMutexPool) {
	logData, err := serializeLogFromRequest(c)
	if handleError(c, err) {
		return
	}

	var method = c.Request.Method
	result := new(core.Response)
	switch method {
	case "POST":
		result, err = postRequestWorker(logData, mutexPool)
	case "GET":
		result, err = getRequestWorker(logData, mutexPool)
	}

	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
	} else if result != nil {
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
	response.Message = "success"
	response.Data = readResult
	return response, nil
}

func postRequestWorker(logModel *models.LogModel, mutexPool *core.FileMutexPool) (*core.Response, error) {

	if err := logModel.WriteLog(mutexPool); err != nil {
		return nil, err
	}

	response := new(core.Response)
	response.Data = nil
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

	switch method {
	case "POST":
		if err := c.BindJSON(logData); err != nil {
			fmt.Println("Error binding to log from POST request.")
		}
	case "GET":
		if err := validateRequestDateStrings(c); err != nil {
			return nil, err
		}
		if err := c.BindQuery(logData); err != nil {
			fmt.Println("Erorr binding to log from GET request.")
			return nil, err
		}
		logData.Location, _ = url.QueryUnescape(logData.Location)
		logData.Message, _ = url.QueryUnescape(logData.Message)
	}

	logData.LogLevel = strings.ToUpper(c.Param("log_level"))

	return logData, nil
}

func validateRequestDateStrings(c *gin.Context) error {
	var keys = []string{"created_date", "from", "to"}
	for _, value := range keys {
		err := validateRequestDateString(c, value)
		if err != nil {
			return err
		}
	}

	return nil
}

func validateRequestDateString(c *gin.Context, key string) error {
	date, exists := c.Get(key)
	if exists {
		if _, err := time.Parse(core.LogDateFormat, date.(string)); err != nil {
			return err
		}
	}

	return nil
}

func checkDateString(date string) bool {
	if _, err := time.Parse(core.LogDateFormat, date); err != nil {
		return false
	}

	return true
}

func handleError(c *gin.Context, err error) bool {
	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return true
	}
	return false
}
