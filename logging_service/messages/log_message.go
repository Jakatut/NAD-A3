package messages

import "time"

type log struct {
	CreatedDate		time.Time	`json:"created_date"`
	UpdatedDate 	time.Time	`json:"updated_date"`
	Severity 		int8		`json:"severity"` // Severity levels are 1-7 (lowest to highest)
	Message 		string		`json:"message"`
	OriginLocation 	string		`json:"orign_location"` // Ideally filename or file location from the software using the logging service.
	MessageNumber 	int			`omitempty`
}