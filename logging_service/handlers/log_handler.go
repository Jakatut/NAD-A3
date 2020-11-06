package handlers

import (
	"errors"
	"fmt"
	"logging_service/core"
	models "logging_service/models"
	"net/url"

	"github.com/gin-gonic/gin"
)

// HandleLog handles all post/get requests for any log type.
func HandleLog(c *Context, mutexPool *core.FileMutexPool) {
	logData, err := serializeLogFromRequest(c)
	if err != nil {
		c.JSON(500, errors.New("issue reading data in log request"))
	}

	var method = c.Request.Method
	if method != "GET" && method != "POST" {
		c.JSON(400, errors.New("unsupported methods"))
	}
	result := new(core.Response)
	switch method {
	case "POST":
		result, err = postRequestWorker(logData, mutexPool)
	case "GET":
		result, err = getRequestWorker(logData, mutexPool)
	}

	if err != nil {
		c.JSON(500, err)
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
			return nil, errors.New("could not bind json to log from POST request")
		}
	case "GET":
		// logData.CreatedDate = time.Parse(c.Get("created_date"))
		if err := c.BindQuery(logData); err != nil {
			fmt.Println("Erorr binding to log from GET request.")
			return nil, err
		}
		logData.Location, _ = url.QueryUnescape(logData.Location)
		logData.Message, _ = url.QueryUnescape(logData.Message)
	default:
		return nil, errors.New("Invalid request type")
	}

	return logData, nil
}
