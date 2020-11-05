package handlers

import (
	"errors"
	"fmt"
	"logging_service/core"
	models "logging_service/models"
	"sync"

	"github.com/gin-gonic/gin"
)

// // HandleGetLog gets all logs
// func HandleGetLog(c *Context, resources core.HandlerResources) {
// 	logData, err := getLogFromRequest(c)
// 	if err != nil {
// 		c.JSON(500, errors.New("issue reading data in log request"))
// 	}
// 	resources.WaitGroup.Add(1)
// 	go getRequestWorker(logData, resources.WaitGroup, resources.LogChannel)
// 	result := <-resources.LogChannel
// 	c.JSON(200, result)
// }

// // HandlePostLog posts a log
// func HandlePostLog(c *Context, resources core.HandlerResources) {
// 	logData, err := getLogFromRequest(c)
// 	if err != nil {
// 		c.JSON(500, errors.New("issue reading data in log request"))
// 	}
// 	resources.WaitGroup.Add(1)
// 	go postRequestWorker(logData, resources.WaitGroup, resources.LogChannel)
// 	result := <-resources.LogChannel
// 	c.JSON(200, result)
// }

// HandleLog handles all post/get requests for any log type.
func HandleLog(c *Context, resources core.HandlerResources) {
	logData, err := getLogFromRequest(c)
	if err != nil {
		c.JSON(500, errors.New("issue reading data in log request"))
	}

	var method = c.Request.Method
	if method != "GET" && method != "POST" {
		c.JSON(400, errors.New("unsupported methods"))
	}
	resources.WaitGroup.Add(1)
	result := new(core.Result)
	switch method {
	case "POST":
		result = postRequestWorker(logData, resources.WaitGroup)
	case "GET":
		result = getRequestWorker(logData, resources.WaitGroup)
	}

	// result := <-resources.LogChannel

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

func getRequestWorker(logMessage *models.LogModel, waitGroup *sync.WaitGroup) *core.Result {
	defer waitGroup.Done()

	var result *core.Result

	// if logModel := log.ReadLog(logMessage); logModel != nil {

	// } else {
	logModel := new(models.LogModel)
	logModel.Severity = 1
	response := new(core.Response)
	response.Data = logModel
	response.Message = "success"
	result = new(core.Result)
	result.Err = nil
	result.Response = response
	return result

	// logChannel <- result
	// }
}

func postRequestWorker(logMessage *models.LogModel, waitGroup *sync.WaitGroup) *core.Result {
	defer waitGroup.Done()

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
	// logChannel <- result
}

/*
 *
 * Helpers
 *
 */

func getLogFromRequest(c *gin.Context) (*models.LogModel, error) {
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
