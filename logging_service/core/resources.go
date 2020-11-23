package core

/*
 *
 * file: 		resources.go
 * project:		logging_service - NAD-A3
 * programmer: 	Conor Macpherson
 * description: Defines core resources such as response structs, date time layouts and mutex functions.
 *
 */

// LogDateFormat used when writing content to log files. Includes time.
const LogDateFormat = "2006-01-02T15:04:05Z"

// CreatedDayFormat used when extracting or fomratting to year-month-day
const CreatedDayFormat = "2006-01-02"

// CreatedTimeFormat used when extracting or formatting to hour-minute-second
const CreatedTimeFormat = "15-04-05"

// ResourceFileNameDateFormat used in the file name when creating log files
const ResourceFileNameDateFormat = "2006-01-02"

// LogLevels defines all available log level types.
var LogLevels = []string{"DEBUG", "INFO", "WARNING", "ERROR", "FATAL"}

// Response defines an api request's response. This would be used for successful responses. Any responses that
// indicate a failure or error should use errors.New("") for the response.
type Response struct {
	Data    interface{} `json:"data,omitempty"`
	Message string      `json:"message,omitempty"`
	Error   string      `json:"error,omitempty"`
}

// Error designed an api's error response.
type Error struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

// FindResults defines the results from a mongodb find. It includes the number of remaining documents and the
// found data.
type FindResults struct {
	RemainingDocuments int64
	Data               interface{}
}
