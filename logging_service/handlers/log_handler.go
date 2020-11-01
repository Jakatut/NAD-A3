package handlers

import (
	"errors"
	"fmt"
	log "logging_service/core"
	"logging_service/messages"
	"net/http"
	"sync"

	"github.com/gin-gonic/gin"
)

// Gets all logs
func HandleGetLog(c *Context, waitGroup *sync.WaitGroup) {
	logData := getLogFromRequest(c)
	waitGroup.Add(1)
	go getRequestWorker(logData)
}

// Post a log
func HandlePostLog(c *Context, waitGroup *sync.WaitGroup) {
	logData := getLogFromRequest(c)
	waitGroup.Add(1)
	go postRequestWorker(logData)
	c.JSON(200, logData)
}

// Workers //

func getRequestWorker(log *messages.Log) {

}

func postRequestWorker(log *messages.Log) {

}

// Helpers //
func getLogFromRequest(c *gin.Context) *messages.Log, error {
	method :=  c.Request.Method
	logData := new(messages.Log)

	switch method {
		case "POST":
			if err := c.BindJSON(logData); err != nil {
				fmt.Println("Error binding to log from POST request.")
				return nil, errors.New("Could not bind json to log from POST request.")
			}
		case "GET":
			if err := c.BindUri(logData); err != nil {
				fmt.Println("Erorr binding to log from GET request.")
				return nil, errors.New("Could not bind json to log from GET request.")
			}
		case default:
			return nil, errors.New("Invalid request type")
	}

	return logData, nil
}

func extractLogFromGetRequest(c *gin.Context) messages.Log, error {

}
