package models

import (
	"errors"
	"logging_service/core"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/globalsign/mgo/bson"
	"github.com/kamva/mgm/v3/operator"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// LogSearchFields defines the fields which users can use to filters logs which contain the same fields when searching.
type LogSearchFields struct {
	ID        primitive.ObjectID
	LogLevel  string
	Location  string
	CreatedAt *time.Time
	FromDate  *time.Time
	ToDate    *time.Time
	OrderBy   string
	Page      int64
	Limit     int64
}

// GetSearchFields all get request fields for a search.
//
// Receiver:
//	*LogSearchFields				lsf
//
// Parameters:
//	*gin.Context	c	- Handler context from gin.
//
//	error				   - Any error that occurs.
//
func (lsf *LogSearchFields) GetSearchFields(c *gin.Context) error {
	createdAt := c.Query("created_at")
	from := c.Query("from")
	to := c.Query("to")
	page := c.Query("page")
	id := c.Query("id")
	location := c.Query("location")
	logLevel := c.Param("log_level")
	orderBy := c.Query("orderby")
	limit := c.Query("limit")

	createdAtDate, err := time.Parse(core.LogDateFormat, createdAt)
	if createdAt != "" && err != nil {
		return errors.New("created_at: invalid date time format")
	}

	fromDate, err := time.Parse(core.LogDateFormat, from)
	if from != "" && err != nil {
		return errors.New("from: invalid date time format")
	}

	toDate, err := time.Parse(core.LogDateFormat, to)
	if to != "" && err != nil {
		return errors.New("to: invalid date time format")
	}

	if valid, _ := IsValidLogLevel(logLevel); !valid {
		return errors.New("log_level: unknown log level")
	}

	pageNumber, err := strconv.Atoi(page)
	if page != "" && err != nil {
		return errors.New("page: must be a number")
	}

	limitNumber, err := strconv.Atoi(limit)
	if limit != "" && err != nil {
		return errors.New("limit: must be a number")
	}

	if !isOrderByFieldValid(orderBy) {
		return errors.New("orderby: must be 'created_at', 'log_level', 'id', or 'location'")
	}

	// Create secondary required date value for to or from if not provided.
	if from != "" && to == "" {
		toDate = time.Now()
	} else if from == "" && to != "" {
		fromDate = time.Unix(0, 0)
	}

	objectID, err := primitive.ObjectIDFromHex(id)
	if id != "" && err != nil {
		return errors.New("id: invalid id")
	}

	lsf.CreatedAt = &createdAtDate
	lsf.Location = location
	lsf.FromDate = &fromDate
	lsf.ToDate = &toDate
	lsf.Page = int64(pageNumber)
	lsf.LogLevel = strings.ToUpper(logLevel)
	lsf.ID = objectID
	lsf.OrderBy = orderBy
	lsf.Limit = int64(limitNumber)

	return nil
}

// getFilters will create mongodb filters for the fields created_at, from, to, location, logLevel, id.
//
// Receiver:
//	*LogSearchFields				lsf
//
// Returns
//	[]map[string]interface{} - List of maps containing mongodb filters.
//
func (lsf *LogSearchFields) getFilters() []map[string]interface{} {
	var createdAtPresent = !lsf.CreatedAt.IsZero()
	var fromDatePresent = lsf.FromDate != nil && !lsf.FromDate.IsZero()
	var toDatePresent = lsf.ToDate != nil && !lsf.ToDate.IsZero()
	var locationPresent = lsf.Location != ""
	var logLevelPresent = lsf.LogLevel != ""
	var searchIDPresent = !lsf.ID.IsZero()

	filters := []map[string]interface{}{}
	if createdAtPresent {
		filters = append(filters, map[string]interface{}{"created_at": lsf.CreatedAt})
	} else if fromDatePresent && toDatePresent {
		filters = append(filters, map[string]interface{}{"created_at": bson.M{operator.Gte: lsf.FromDate, operator.Lte: lsf.ToDate}})
	}
	if locationPresent {
		filters = append(filters, map[string]interface{}{"location": lsf.Location})
	}
	if logLevelPresent {
		filters = append(filters, map[string]interface{}{"log_level": lsf.LogLevel})
	} else {
		filters = append(filters, map[string]interface{}{"log_level": bson.M{operator.In: core.LogLevels}})
	}

	if searchIDPresent {
		filters = append(filters, map[string]interface{}{"_id": lsf.ID})
	}

	return filters
}

func (lsf *LogSearchFields) getFindOptions() *options.FindOptions {
	var orderByPresent = lsf.OrderBy != ""
	options := options.Find()
	if orderByPresent {
		options.SetSort(bson.D{{lsf.OrderBy, -1}})
	}

	return options
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
