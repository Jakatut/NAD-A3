package validators

import "time"

type LogModel struct {
	CreatedDate time.Time `json:",omitempty" form:"created_date,omitempty" time_format:"2006-01-02T15:04:05"`
	Severity    int       `json:"severity,omitempty" form:"severity,omitempty" binding:"gte=1,lte=7"` // Severity levels are 1-7 (lowest to highest)
	LogLevel    string    `validate:"DEBUG|WARNING|INFO|ERROR|FATAL"`                                 // DEBUG, INFO, WARN, ERROR, FATAL, ALL
	Message     string    `json:"message,omitempty" form:"message,omitempty"`
	Location    string    `json:"location,omitempty" form:"location,omitempty"` // Ideally filename or file location from the software using the logging service.
	FromTime    time.Time `json:",omitempty" form:"from,omitempty" binding:"omitempty" time_format:"2006-01-02T15:04:05"`
	ToTime      time.Time `json:",omitempty" form:"to,omitempty" binding:"omitempty,gtfield=FromTime" time_format:"2006-01-02T15:04:05"`
}
