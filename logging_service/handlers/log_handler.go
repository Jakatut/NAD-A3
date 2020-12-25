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
	"errors"
	"log"
	"logging_service/core"
	"logging_service/models"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
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
	fields, err := getSearchFields(c)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"Error": err.Error()})
	}

	log := models.Log{}
	results, err := log.Find(fields)
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
	fields, err := getSearchFields(c)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"Error": err.Error()})
	}

	_log := models.Log{}
	count, err := _log.Count(fields)
	if err != nil {
		log.Fatal(err.Error())
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"Error": "internal server error"})
	} else {
		results := core.CountResults{Total: count}
		c.JSON(200, results)
	}
}

// HandleGetLog handles all get requests for any log type.
//
// Parameters:
//	*gin.Context	c	- Handler context from gin.
//
//	models.LogSearchFields - Struct containing log search fields.
//	error				   - Any error that occurs.
//
func getSearchFields(c *gin.Context) (models.LogSearchFields, error) {
	createdAt := c.Query("created_at")
	from := c.Query("from")
	to := c.Query("to")
	page := c.Query("page")
	id := c.Query("id")
	location := c.Query("location")
	logLevel := c.Param("log_level")
	orderBy := c.Query("orderby")
	limit := c.Query("limit")


	searchFields := models.LogSearchFields{}

	createdAtDate, err := time.Parse(core.LogDateFormat, createdAt)
	if createdAt != "" && err != nil {
		return searchFields, errors.New("created_at: invalid date time format")
	}

	fromDate, err := time.Parse(core.LogDateFormat, from)
	if from != "" && err != nil {
		return searchFields, errors.New("from: invalid date time format")
	}

	toDate, err := time.Parse(core.LogDateFormat, to)
	if to != "" && err != nil {
		return searchFields, errors.New("to: invalid date time format")
	}

	pageNumber, err := strconv.Atoi(page)
	if page != "" && err != nil {
		return searchFields, errors.New("page: must be a number")
	}
	limitNumber, err := strconv.Atoi(limit)
	if limit != "" && err != nil {
		return searchFields, errors.New("limit: must be a number")
	}

	if !isOrderByFieldValid(orderBy) {
		return searchFields, errors.New("orderby: must be 'created_at', 'log_level', 'id', or 'location'")
	}

	// Create secondary required date value for to or from if not provided.
	if from != "" && to == "" {
		toDate = time.Now()
	} else if from == "" && to != "" {
		fromDate = time.Unix(0, 0)
	}

	objectID, err := primitive.ObjectIDFromHex(id)
	if id != "" && err != nil {
		return searchFields, errors.New("id: invalid id")
	}

	return models.LogSearchFields{
		CreatedAt: &createdAtDate,
		Location:  location,
		FromDate:  &fromDate,
		ToDate:    &toDate,
		Page:      int64(pageNumber),
		LogLevel:  strings.ToUpper(logLevel),
		ID:        objectID,
		OrderBy:   orderBy,
		Limit:     int64(limitNumber),
	}, nil
}

func isOrderByFieldValid(orderByField string) bool {
	var validOrderByField = false
	searchFields := []string{"created_at", "id", "location", "log_level", ""}
	for _, val := range searchFields {
		if val == orderByField {
			validOrderByField = true
		}
	}
	return validOrderByField
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
	if logLevel != "" && !IsValidLogLevel(logLevel) {
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

// IsValidLogLevel check the provided logLevel is one of "DEBUG", "WARNING", "ERROR", "FATAL" or "INFO"
//
// Parameters:
//	string	logLevel	- Log level to get the last file for.
//
// Returns
//	bool - True if the given log level is a valid log level.
func IsValidLogLevel(logLevel string) bool {
	for _, value := range core.LogLevels {
		if logLevel == "" || strings.ToUpper(value) == strings.ToUpper(logLevel) {
			return true
		}
	}

	return false
}
