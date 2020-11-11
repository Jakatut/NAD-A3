package models

import (
	"bufio"
	"fmt"
	"log"
	"logging_service/core"
	"os"
	"strings"
	"time"
)

// Log Types
const (
	ALL   = 0
	DEBUG = 1
	INFO  = 2
	WARN  = 3
	ERROR = 4
	FATAL = 5
)

//json:"created_date,omitempty" form:"created_date,omitempty"

// LogModel defines the contents of a log
type LogModel struct {
	CreatedDate *time.Time `json:",omitempty" form:"created_date,omitempty" time_format:"2006-01-02T15:04:05Z"`
	CreatedDay  *time.Time `json:",omitempty" form:"created_day,omitempty" time_format:"2006-01-02"`
	CreatedTime *time.Time `json:",omitempty" form:"created_time,omitempty" time_format:"15:04:05Z"`
	ID          uint       `json:"id,omitempty" form:"id,omitempty"`
	LogLevel    string     `json:"type,omitempty" form:"type,omitempty" validate:"DEBUG|WARNING|INFO|ERROR|FATAL"` // DEBUG, INFO, WARN, ERROR, FATAL, ALL
	Message     string     `json:"message" form:"message,omitempty"`
	Location    string     `json:"location,omitempty" form:"location,omitempty"`
	FromDate    *time.Time `json:",omitempty" form:"from,omitempty" binding:"omitempty" time_format:"2006-01-02T15:04:05Z" binding:"required"`
	ToDate      *time.Time `json:",omitempty" form:"to,omitempty" binding:"omitempty" time_format:"2006-01-02T15:04:05Z" binding:"required"`
}

// WriteLog writes a log to a logfile.
func (logModel *LogModel) WriteLog(mutexPool *core.FileMutexPool) error {
	logLocation, err := core.GetLogWriteLocation(logModel.LogLevel)
	if err != nil {
		return err
	}
	core.CreateLogLevelDirectory(logModel.LogLevel)

	mutexPool.LockWriteFileMutex(logLocation)
	file, err := os.OpenFile(logLocation, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatal(err)
		return err
	}
	defer file.Close()

	bufferedWriter := bufio.NewWriter(file)
	_, err = bufferedWriter.Write(
		logModel.buildLogMessage(),
	)
	if err != nil {
		log.Fatal(err)
		return err
	}
	bufferedWriter.Flush()
	mutexPool.UnlockWirterFileMutex(logLocation)

	return nil
}

// ReadLog reads a log from the log file.
func (logModel *LogModel) ReadLog(mutexPool *core.FileMutexPool) ([]LogModel, error) {

	var logs = []LogModel{}

	logLocations := core.GetSearchFilePaths(logModel.LogLevel)
	for _, location := range logLocations {
		mutexPool.LockReadFileMutex(location)
		foundLogs, err := searchLog(location, logModel)
		if err != nil {
			return nil, err
		}
		logs = append(logs, foundLogs...)
		mutexPool.UnlockReadFileMutex(location)
	}

	return logs, nil
}

// filter compares the values between two log models: The receiver and the comparison.
// If the two models are the same, true is returned. Otherwise, false.
func (logModel *LogModel) filter(comparison *LogModel) bool {
	return logModel.compareCreatedDateValues(comparison) && logModel.filterWithoutCreatedDate(comparison)
}

// filterWithoutCreatedDate compares the values between two log models: The receiver and the comparison.
// If the two models are the same, true is returned. Otherwise, false.
func (logModel *LogModel) filterWithoutCreatedDate(comparison *LogModel) bool {

	if (logModel.ID != 0) && logModel.ID != comparison.ID {
		return false
	}
	if (logModel.Message != "") && logModel.Message != comparison.Message {
		return false
	}
	if (logModel.Location != "") && logModel.Location != comparison.Location {
		return false
	}

	return true
}

func (logModel *LogModel) compareCreatedDateValues(comparison *LogModel) bool {
	var createdDatePresent = logModel.CreatedDate != nil && !logModel.CreatedDate.IsZero()
	var createdTimePresent = logModel.CreatedTime != nil && !logModel.CreatedTime.IsZero()
	var createdDayPresent = logModel.CreatedDay != nil && !logModel.CreatedDay.IsZero()
	var fromDatePresent = logModel.FromDate != nil && !logModel.FromDate.IsZero()
	var toDatePresent = logModel.ToDate != nil && !logModel.ToDate.IsZero()

	if !createdDatePresent &&
		!createdTimePresent &&
		!createdDayPresent &&
		!fromDatePresent &&
		!toDatePresent {

		return true
	}

	var logModelCreatedDate = ""
	if logModel.CreatedDate != nil {
		logModelCreatedDate = logModel.CreatedDate.Format(core.LogDateFormat)
	}
	var comparisonCreatedDate = ""
	if comparison.CreatedDate != nil {
		comparisonCreatedDate = comparison.CreatedDate.Format(core.LogDateFormat)
	}
	var logModelCreatedTime = ""
	if logModel.CreatedTime != nil {
		logModelCreatedTime = logModel.CreatedTime.Format(core.CreatedTimeFormat)
	}
	var comparisonCreatedTime = ""
	if comparison.CreatedDate != nil {
		comparisonCreatedTime = comparison.CreatedDate.Format(core.CreatedTimeFormat)
	}
	var logModelCreatedDay = ""
	if logModel.CreatedDay != nil {
		logModelCreatedDay = logModel.CreatedDay.Format(core.CreatedDayFormat)
	}
	var comparisonCreatedDay = ""
	if comparison.CreatedDate != nil {
		comparisonCreatedDay = comparison.CreatedDate.Format(core.CreatedDayFormat)
	}

	if createdDatePresent && logModelCreatedDate == comparisonCreatedDate {
		return true
	}
	if (createdTimePresent) && (logModelCreatedTime == comparisonCreatedTime) {
		return true
	}
	if (createdDayPresent) && (logModelCreatedDay == comparisonCreatedDay) {
		return true
	}
	if (fromDatePresent && toDatePresent) && comparison.createdDateFallsWithinDateRange(*logModel.FromDate, *logModel.ToDate) {
		return true
	}

	return false
}

// createdDateFallsWithinDateRange checks if the created date of the reciever, falls within the provided date range.
func (logModel *LogModel) createdDateFallsWithinDateRange(fromTime time.Time, toTime time.Time) bool {
	if !logModel.CreatedDate.IsZero() && !fromTime.IsZero() && !toTime.IsZero() {
		return logModel.CreatedDate.After(fromTime.Add(-time.Second*1)) && logModel.CreatedDate.Before(toTime.Add(time.Second*1))
	}

	return false
}

/*
 *
 * Helpers
 *
 */

// search log
func searchLog(location string, logModel *LogModel) ([]LogModel, error) {
	var foundLogs []LogModel
	file, err := os.Open(location)
	if err != nil {
		return foundLogs, err
	}
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		logLine, err := rawLogToModel(scanner.Text(), logModel.LogLevel)
		if err != nil {
			return foundLogs, err
		}
		if logLine != nil && logModel.filter(logLine) {
			logLine.Message = strings.Replace(logLine.Message, "\\n", "\n", -1)
			foundLogs = append(foundLogs, *logLine)
		}
	}

	return foundLogs, nil
}

func rawLogToModel(rawLog string, logType string) (*LogModel, error) {
	logModel := new(LogModel)
	details, err := core.GetLogDetailsFromRawLog(rawLog, logType)
	if err != nil {
		return nil, err
	}

	logModel.CreatedDate = new(time.Time)
	*logModel.CreatedDate = (details["created_date"]).(time.Time)
	logModel.ID = uint((details["id"]).(uint64))
	logModel.Location = (details["location"]).(string)
	logModel.Message = (details["message"]).(string)
	logModel.LogLevel = logType

	return logModel, nil
}

func (logModel *LogModel) buildLogMessage() []byte {
	messageText := strings.Replace(logModel.Message, "\n", "\\n", -1)
	messageText = strings.Replace(messageText, "\"", "\\\"", -1)
	location := strings.Replace(logModel.Location, "\"", "\\\"", -1)
	return []byte(fmt.Sprintf("[date=\"%s\"  id=\"%d\" location=\"%s\"]:\"%s\"\n", time.Now().Format(core.LogDateFormat), logModel.ID, location, messageText))
}
