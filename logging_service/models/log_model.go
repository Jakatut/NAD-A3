package messages

import "time"

// Log Types
const (
	DEBUG = 1
	INFO  = 2
	WARN  = 3
	ERROR = 4
	FATAL = 5
	ALL   = 6
)

// Log message
type LogModel struct {
	CreatedDate   time.Time `json:"created_date"`
	Severity      int8      `json:"severity"` // Severity levels are 1-7 (lowest to highest)
	Type          int8      `json:"type"`     // DEBUG, INFO, WARN, ERROR, FATAL, ALL
	Message       string    `json:"message"`
	Location      string    `json:"location"` // Ideally filename or file location from the software using the logging service.
	MessageNumber int       `omitempty`       // maybe don't ommit?
}
