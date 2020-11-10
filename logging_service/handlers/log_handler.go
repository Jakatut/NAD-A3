package handlers

import (
	"errors"
	"logging_service/core"
	models "logging_service/models"
	"net/http"
	"net/url"
	"strings"

	ut "github.com/go-playground/universal-translator"

	"github.com/go-playground/locales/en"
	"github.com/go-playground/validator/v10"

	"github.com/gin-gonic/gin"
)

// HandlePostLog handles all post requests for any log type.
func HandlePostLog(c *gin.Context, mutexPool *core.FileMutexPool, counters *core.LogTypeCounter) {
	logData, err := serializeLogFromRequest(c)
	if err != nil {
		return
	}

	result, err := postRequestWorker(logData, mutexPool, counters)

	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
	} else if result.Data != nil {
		c.JSON(200, result)
	} else {
		c.JSON(200, gin.H{"message": "success"})
	}
}

// HandleGetLog handles all get requests for any log type.
func HandleGetLog(c *gin.Context, mutexPool *core.FileMutexPool) {
	logData, err := serializeLogFromRequest(c)
	if err != nil {
		return
	}

	result, err := getRequestWorker(logData, mutexPool)

	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
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

	logLevel := strings.ToUpper(c.Param("log_level"))
	if !core.IsValidLogLevel(logLevel) {
		return nil, errors.New("invalid log level")
	}
	logData.LogLevel = logLevel

	switch method {
	case "POST":
		if err := c.ShouldBindJSON(logData); err != nil {
			if logLevel == "ALL" {
				c.AbortWithStatusJSON(http.StatusBadRequest, "not implemented.")
			}
			handleError(c, err)
			return nil, err
		}
	case "GET":
		if err := c.ShouldBindQuery(logData); err != nil {
			handleError(c, err)
			return nil, err
		}
		logData.Location, _ = url.QueryUnescape(logData.Location)
		logData.Message, _ = url.QueryUnescape(logData.Message)
	}

	return logData, nil
}

func handleError(c *gin.Context, err error) bool {
	if err != nil {
		en := en.New()
		uni := ut.New(en, en)
		trans, _ := uni.GetTranslator("en")
		// response := core.Response{}
		ve, ok := err.(validator.ValidationErrors)
		if ok {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"Error": ve.Translate(trans)})
			return true
		}
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"Error": err})
	}
	return false
}
