package handlers

import (
	"fmt"
	"log"
	"logging_service/core"
	models "logging_service/models"
	"net/http"
	"net/url"
	"strings"

	"github.com/gin-gonic/gin"
	"gopkg.in/bluesuncorp/validator.v5"
)

// HandlePostLog handles all post requests for any log type.
func HandlePostLog(c *gin.Context, mutexPool *core.FileMutexPool, counters *core.LogTypeCounter) {
	logData, err := serializeLogFromRequest(c)
	if err != nil {
		return
	}

	result, err := postRequestWorker(logData, mutexPool, counters)

	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
	}

	c.JSON(200, result)
}

// HandleGetLog handles all get requests for any log type.
func HandleGetLog(c *gin.Context, mutexPool *core.FileMutexPool) {
	logData, err := serializeLogFromRequest(c)
	if err != nil {
		return
	}

	result, err := getRequestWorker(logData, mutexPool)

	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		log.Fatal(err)
	}

	if logData != nil && len(result.Data.([]models.LogModel)) < 0 {
		c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"error": "not found"})

	}

	c.JSON(200, result)
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
	}
	logData.LogLevel = logLevel

	switch method {
	case "POST":
		if logLevel == "ALL" {
			c.AbortWithStatusJSON(http.StatusNotFound, "Not found.")
		}
		if err := c.ShouldBindJSON(logData); err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, err)
		}
	case "GET":
		if err := c.ShouldBindQuery(logData); err != nil {

		}
		logData.Location, _ = url.QueryUnescape(logData.Location)
		logData.Message, _ = url.QueryUnescape(logData.Message)
	}

	checkErrors(c)

	return logData, nil
}

func ValidationErrorToText(e *validator.FieldError) string {
	switch e.Tag {
	case "required":
		return fmt.Sprintf("%s is required", e.Field)
	case "max":
		return fmt.Sprintf("%s cannot be longer than %s", e.Field, e.Param)
	case "min":
		return fmt.Sprintf("%s must be longer than %s", e.Field, e.Param)
	case "email":
		return fmt.Sprintf("Invalid email format")
	case "len":
		return fmt.Sprintf("%s must be %s characters long", e.Field, e.Param)
	}
	return fmt.Sprintf("%s is not valid", e.Field)
}

func checkErrors(c *gin.Context) {
	if len(c.Errors) > 0 {
		status := http.StatusBadRequest
		errors := make([]map[string]string, 10)
		for _, e := range c.Errors {
			switch e.Type {
			case gin.ErrorTypeBind:
				errs := e.Err.(*validator.StructErrors)
				list := make(map[string]string)
				for field, err := range errs.Errors {
					list[field] = ValidationErrorToText(err)
				}
				errors = append(errors, list)
			}
			// Make sure we maintain the preset response status

			if c.Writer.Status() != http.StatusOK {
				status = c.Writer.Status()
			}
		}

		if len(errors) > 0 {
			c.AbortWithStatusJSON(status, gin.H{"Errors": errors})
		}
	}
}
