package models

/*
 *
 * file: 		log_model.go
 * project:		logging_service - NAD-A3
 * programmer: 	Conor Macpherson
 * description: Defines the log data structure, and attaches receiver methods to the struct.
 *
 */

import (
	"logging_service/config"
	"logging_service/core"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/kamva/mgm/v3/operator"

	"github.com/globalsign/mgo/bson"

	"github.com/kamva/mgm/v3"
)

// LogSearchFields defines the fields which users can use to filters logs which contain the same fields when searching.
type LogSearchFields struct {
	ID        primitive.ObjectID
	LogLevel  string
	Location  string
	CreatedAt *time.Time
	FromDate  *time.Time
	ToDate    *time.Time
	Page      int64
}

// Log defines the contents of a log
type Log struct {
	ID        primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty" binding:"-"`
	CreatedAt time.Time          `bson:"created_at" json:"created_at"`
	UpdatedAt time.Time          `bson:"updated_at,omitempty" json:"-" form:"-"`
	LogLevel  string             `bson:"log_level" json:"log_level,omitempty" form:"log_level,omitempty" validate:"DEBUG|WARNING|INFO|ERROR|FATAL|ALL"`
	Message   string             `bson:"message" json:"message" form:",omitempty"`
	Extra     []string           `bson:"extra,omitempty" json:"extra,omitempty"`
	Location  string             `bson:"location" json:"location" form:"location,omitempty"`
}

// PrepareID method prepares by creating an object id from a string id.
//
// Receiver:
//	*Log				_log
//
// Parameters
//	interface{}	-	id	- The id to be prepared.
//
// Returns
//	interface{}	-	The id. as an object id.
//	error		-	Any error that occurs.
//
func (_log *Log) PrepareID(id interface{}) (interface{}, error) {
	if idStr, ok := id.(string); ok {
		return primitive.ObjectIDFromHex(idStr)
	}

	// Otherwise id must be ObjectId
	return id, nil
}

// GetID method return model's id
//
// Receiver:
//	*Log				_log
//
// Returns
//	interface{}	-	The id.
//
func (_log *Log) GetID() interface{} {
	return _log.ID
}

// SetID set id value of model's id field.
//
// Receiver:
//	*Log				_log
//
// Parameters
//	interface{}	-	id	- The id to be set.
//
func (_log *Log) SetID(id interface{}) {
	_log.ID = id.(primitive.ObjectID)
}

// Create creates a log in the mongodb log collection.
//
// Receiver:
//	*Log				_log
//
// Returns
//	error - Any error that occurs.
//
func (_log *Log) Create() error {
	err := mgm.Coll(_log, &options.CollectionOptions{}).Create(_log)
	return err
}

// Find searches the log collection to find any logs that match the search criteria.
//
// Receiver:
//	*Log				_log
//
// Parameters:
//	LogSearchFields		fields - Search fields.
//
// Returns
//	[]map[string]interface{} - List of maps containing mongodb filters.
//
func (_log *Log) Find(fields LogSearchFields) (core.FindResults, error) {

	limit := config.GetConfig().Results.Limit
	findOptions := options.Find()
	findOptions.SetLimit(limit)
	findOptions.SetSkip(limit * fields.Page)

	logsColl := mgm.Coll(_log)
	logs := []Log{}

	filters := fields.getFilters()
	filter := bson.M{operator.And: filters}
	if len(filters) == 0 {
		filter = bson.M{}
	} else if len(filters) == 1 {
		filter = filters[0]
	}
	err := logsColl.SimpleFind(&logs, filter, findOptions)

	countOptions := options.Count()
	countOptions.SetSkip(limit * (fields.Page + 1))
	remainingDocumentCount, err := logsColl.CountDocuments(mgm.Ctx(), filter, countOptions)
	results := core.FindResults{Data: logs, RemainingDocuments: remainingDocumentCount}
	return results, err
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

// IsEmptyCreate checks that the struct is not nil, and that the message and location are not empty.
// If any of these are true, ture is returned.
//
// Receiver:
//	*Log				logModel
//
// Returns
//	[]string	- Slice of validation messages.
//	bool		- True if empty.
//
func (_log *Log) IsEmptyCreate() ([]string, bool) {
	if _log == nil {
		return []string{"missing field: message", "missing field: location"}, true
	}

	missingFields := []string{}

	if _log.Message == "" {
		missingFields = append(missingFields, "missing field: message")
	}
	if _log.Location == "" {
		missingFields = append(missingFields, "missing field: location")
	}

	return missingFields, len(missingFields) > 0
}
