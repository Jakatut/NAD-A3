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
	"context"
	"logging_service/config"
	"logging_service/core"
	"strings"
	"time"

	"go.mongodb.org/mongo-driver/mongo"

	"github.com/kamva/mgm/v3/operator"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/kamva/mgm/v3"
)

// Log defines the contents of a log
type Log struct {
	ID        primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty" binding:"-"`
	CreatedAt time.Time          `bson:"created_at" json:"created_at"`
	UpdatedAt time.Time          `bson:"updated_at,omitempty" json:"-" form:"-"`
	LogLevel  string             `bson:"log_level" json:"log_level,omitempty" form:"log_level,omitempty" validate:"DEBUG|WARNING|INFO|ERROR|FATAL"`
	Message   string             `bson:"message" json:"message" form:",omitempty"`
	Extra     []string           `bson:"extra,omitempty" json:"extra,omitempty"`
	Location  string             `bson:"location" json:"location" form:"location,omitempty"`
}

// PrepareID method prepares by creating an object id from a string id.
//
// Receiver:
//	*Log				l
//
// Parameters
//	interface{}	-	id	- The id to be prepared.
//
// Returns
//	interface{}	-	The id. as an object id.
//	error		-	Any error that occurs.
//
func (l *Log) PrepareID(id interface{}) (interface{}, error) {
	if idStr, ok := id.(string); ok {
		return primitive.ObjectIDFromHex(idStr)
	}

	// Otherwise id must be ObjectId
	return id, nil
}

// GetID method return model's id
//
// Receiver:
//	*Log				l
//
// Returns
//	interface{}	-	The id.
//
func (l *Log) GetID() interface{} {
	return l.ID
}

// SetID set id value of model's id field.
//
// Receiver:
//	*Log				l
//
// Parameters
//	interface{}	-	id	- The id to be set.
//
func (l *Log) SetID(id interface{}) {
	l.ID = id.(primitive.ObjectID)
}

// Create creates a log in the mongodb log collection.
//
// Receiver:
//	*Log				l
//
// Returns
//	error - Any error that occurs.
//
func (l *Log) Create() error {
	err := mgm.Coll(l, &options.CollectionOptions{}).Create(l)
	return err
}

// Find searches the log collection to find any logs that match the search criteria.
//
// Receiver:
//	*Log				l
//
// Parameters:
//	LogSearchFields		fields - Search fields.
//
// Returns
//	[]map[string]interface{} - List of maps containing mongodb filters.
//
func (l *Log) Find(ctx context.Context, fields LogSearchFields) (core.FindResults, error) {
	configs := config.GetConfig()
	limit := configs.Results.Limit
	suppliedLimit := fields.Limit

	if suppliedLimit != 0 && suppliedLimit < limit {
		limit = suppliedLimit
	}

	findOptions := fields.getFindOptions()
	findOptions.SetLimit(limit)
	findOptions.SetSkip(limit * fields.Page)
	logsColl := mgm.Coll(l)
	logs := []Log{}

	filter := GetFilter(fields)
	err := logsColl.SimpleFind(&logs, filter, findOptions)

	countOptions := options.Count()
	totalDocuments, err := logsColl.CountDocuments(ctx, filter, countOptions)
	countOptions.SetSkip(limit * (fields.Page + 1))
	remainingDocumentCount, err := logsColl.CountDocuments(ctx, filter, countOptions)
	results := core.FindResults{Data: logs, Remaining: remainingDocumentCount, Total: totalDocuments, Limit: configs.Results.Limit}
	return results, err
}

// Count returns the count of logs based on the provided log search fields.
//
// Receiver:
//	*Log				l
//
// Parameters:
//	LogSearchFields		fields - Search fields.
//
// Returns
//	core.CountResults
//  error
//
func (l *Log) Count(ctx context.Context, fields LogSearchFields) (core.CountResults, error) {
	log := Log{}
	logsColl := mgm.Coll(&log)
	countOptions := options.Count()
	_, all := IsValidLogLevel(fields.LogLevel)
	if all {
		fields.LogLevel = ""
	}

	filter := bson.M{operator.And: fields.getFilters()}
	totalDocuments, err := logsColl.CountDocuments(ctx, filter, countOptions)
	results := core.CountResults{}
	results.Count = totalDocuments
	return results, err
}

// Count returns the count of logs based on the provided log search fields.
//
// Receiver:
//	*Log				l
//
// Parameters:
//	LogSearchFields		fields - Search fields.
//
// Returns
//	core.CountResults
//  error
//
func (l *Log) CountByDates(ctx context.Context, fields LogSearchFields) ([]core.CountResultsWithDate, error) {
	log := Log{}
	logsColl := mgm.Coll(&log)
	_, all := IsValidLogLevel(fields.LogLevel)
	if all {
		fields.LogLevel = ""
	}

	filter := GetFilter(fields)
	matchStage := bson.D{{operator.Match, filter}}
	groupStage := bson.D{
		{
			operator.Group, bson.M{
				"_id": bson.M{
					"date": bson.M{
						operator.DateToString: bson.M{
							"format": "%Y-%m-%d", "date": "$created_at",
						},
					},
					"log_level": "$log_level",
				},
				"count": bson.M{operator.Sum: 1},
			},
		},
	}

	counts := []core.CountResultsWithDate{}
	countByDatesCursor, err := logsColl.Aggregate(ctx, mongo.Pipeline{matchStage, groupStage})
	if err != nil {
		return counts, err
	}
	countByDatesCursor.All(ctx, &counts)

	return counts, nil
}

// IsEmptyCreate checks that the struct is not nil, and that the message and location are not empty.
// If any of these are true, ture is returned.
//
// Receiver:
//	*Log				logModel
//
// Returns:
//	[]string	- Slice of validation messages.
//	bool		- True if empty.
//
func (l *Log) IsEmptyCreate() ([]string, bool) {
	if l == nil {
		return []string{"missing field: message", "missing field: location"}, true
	}

	missingFields := []string{}

	if l.Message == "" {
		missingFields = append(missingFields, "missing field: message")
	}
	if l.Location == "" {
		missingFields = append(missingFields, "missing field: location")
	}

	return missingFields, len(missingFields) > 0
}

// IsValidLogLevel check the provided logLevel is one of "DEBUG", "WARNING", "ERROR", "FATAL", "INFO", "ALL"|""
//
// Parameters:
//	string	logLevel	- Log level to get the last file for.
//
// Returns
//	bool - True if the given log level is a valid log level.
//  bool - True if the given log level is ALL
func IsValidLogLevel(logLevel string) (bool, bool) {
	if strings.ToUpper(logLevel) == "ALL" {
		return true, true
	}

	for _, val := range core.LogLevels {
		if strings.ToUpper(logLevel) == val || logLevel == "" {
			return true, false
		}
	}

	return false, false
}

func GetFilter(fields LogSearchFields) bson.M {
	filters := fields.getFilters()
	filter := bson.M{operator.And: filters}
	if len(filters) == 0 {
		filter = bson.M{}
	} else if len(filters) == 1 {
		filter = filters[0]
	}

	return filter
}
