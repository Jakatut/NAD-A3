package handlers

import (
	"errors"
	"fmt"
	"logging_service/core"
	models "logging_service/models"

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
	result := new(core.Result)
	switch method {
	case "POST":
		result = postRequestWorker(logData, mutexPool)
	case "GET":
		result = getRequestWorker(logData, mutexPool)
	}

	if result.Err != nil {
		c.JSON(500, result.Err)
	} else {
		if result.Response.Data == nil {
			fmt.Print("")
		}
		c.JSON(200, result.Response.Data)
	}
}

/*
 *
 * Workers
 *
 */

func getRequestWorker(log *models.LogModel, mutexPool *core.FileMutexPool) *core.Result {

	var result *core.Result

	log.ReadLog()

	logModel := new(models.LogModel)
	logModel.Severity = 1
	response := new(core.Response)
	response.Data = logModel
	response.Message = "success"
	result = new(core.Result)
	result.Err = nil
	result.Response = response
	return result
}

func postRequestWorker(logMessage *models.LogModel, mutexPool *core.FileMutexPool) *core.Result {
	var result *core.Result

	// if err := log.WriteLog(logMessage); err != nil {

	// } else {
	// 	response := new(core.Response)
	// 	response.Data = nil
	// 	response.Message = ""
	// 	result = new(core.Result)
	// 	result.Err = nil
	// 	result.Response = response
	// 	logChannel <- result
	// }

	logModel := new(models.LogModel)
	logModel.Severity = 1
	logMessage.Message = "hello"
	response := new(core.Response)
	response.Data = logModel
	response.Message = ""
	result = new(core.Result)
	result.Err = nil
	result.Response = response
	return result
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
		// if err := c.BindUri(logData); err != nil {
		// 	fmt.Println("Erorr binding to log from GET request.")
		// 	return nil, errors.New("could not bind json to log from GET request")
		// }
	default:
		return nil, errors.New("Invalid request type")
	}

	return logData, nil
}
