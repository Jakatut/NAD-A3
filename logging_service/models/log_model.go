package models

import (
	"bufio"
	"fmt"
	"log"
	"logging_service/core"
	"os"
	"path/filepath"
	"strconv"
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
	ID          uint       `json:",omitempty" form:"id,omitempty"`
	LogLevel    string     `json:"type,omitempty" form:"type,omitempty" validate:"DEBUG|WARNING|INFO|ERROR|FATAL"` // DEBUG, INFO, WARN, ERROR, FATAL, ALL
	Message     string     `json:"message" form:"message,omitempty"`
	Location    string     `json:"location,omitempty" form:"location,omitempty"`
	FromDate    *time.Time `json:",omitempty" form:"from,omitempty" binding:"omitempty" time_format:"2006-01-02T15:04:05Z" binding:"required"`
	ToDate      *time.Time `json:",omitempty" form:"to,omitempty" binding:"omitempty" time_format:"2006-01-02T15:04:05Z" binding:"required"`
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

// WriteLog writes a log to a logfile.
func (logModel *LogModel) WriteLog(mutexPool *core.FileMutexPool) error {
	logLocation, err := core.GetLogWriteLocation(logModel.LogLevel)
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

	logLocations := getSearchFilePaths(logModel.LogLevel)
	for _, location := range logLocations {
		mutexPool.LockReadFileMutex(location)
		logs = append(logs, searchLog(location, logModel)...)
		mutexPool.UnlockReadFileMutex(location)
	}

	return logs, nil
}

/*
 *
 * Helpers
 *
 */

// getSearchFilePaths get a list of file paths for the files that need to be searched.
func getSearchFilePaths(logLevel string) []string {
	var paths []string
	if logLevel != "ALL" {
		paths = getLogLevelPaths([]string{logLevel})
	} else {
		paths = getLogLevelPaths(core.LogLevels)
	}

	return paths
}

// getLogLevelPaths goes through the list of log levels, and finds the path of every log in that level's directory.
func getLogLevelPaths(logLevels []string) []string {
	var paths []string
	for _, level := range logLevels {
		err := filepath.Walk(""+level+"/", func(path string, info os.FileInfo, err error) error {
			paths = append(paths, path)
			return nil
		})
		if err != nil {
			log.Fatal(err)
			return nil
		}
	}

	return paths
}

// search log
func searchLog(location string, logModel *LogModel) []LogModel {
	var foundLogs []LogModel
	file, _ := os.Open(location)
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		fmt.Println(scanner.Text())
		logLine := rawLogToModel(scanner.Text(), logModel.LogLevel)
		if logModel.filter(logLine) {
			logLine.Message = strings.Replace(logLine.Message, "\\n", "\n", -1)
			foundLogs = append(foundLogs, *logLine)
		}
	}

	return foundLogs
}

func rawLogToModel(rawLog string, logType string) *LogModel {
	logTextIndicator := strings.Index(rawLog, ":")
	// Remove leading and trailing braces, removes the content of the log, and splits the details.
	logProperties := strings.Split(rawLog[1:logTextIndicator-1], "]-[")

	createdDate := new(time.Time)
	*createdDate, _ = time.Parse(core.LogDateFormat, logProperties[0])

	var logModel = new(LogModel)
	logModel.CreatedDate = createdDate
	logModel.Location = strings.Replace(logProperties[1], "\\-", ":", -1)
	id, _ := strconv.Atoi(logProperties[2])

	logModel.ID = uint(id)
	logModel.LogLevel = logType
	logModel.Message = rawLog[logTextIndicator+2 : len(rawLog)-1]

	return logModel
}

func (logModel *LogModel) buildLogMessage() []byte {
	location := strings.Replace(logModel.Location, ":", "\\-", -1)
	messageText := strings.Replace(logModel.Message, "\n", "\\n", -1)
	return []byte(fmt.Sprintf("[%s]-[%s]-[%d]:\"%s\"\n", time.Now().Format(core.LogDateFormat), location, logModel.ID, messageText))
}
