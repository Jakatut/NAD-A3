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
	ALL         = 0
	DEBUG       = 1
	INFO        = 2
	WARN        = 3
	ERROR       = 4
	FATAL       = 5
	MinSeverity = 1
	MaxSeverity = 7
)

var logLevels = []string{"ALL", "DEBUG", "INFO", "WARN", "ERROR", "FATAL"}

//json:"created_date,omitempty" form:"created_date,omitempty"

// LogModel defines the contents of a log
type LogModel struct {
	CreatedDate time.Time `json:",omitempty" form:"created_date,omitempty" time_format:"2006-01-02T15:04:05"`
	Severity    int       `json:"severity,omitempty" form:"severity,omitempty" bidning:"gte=1,lte=7"` // Severity levels are 1-7 (lowest to highest)
	LogLevel    string    `json:"type,omitempty" form:"type,omitempty"`                               // DEBUG, INFO, WARN, ERROR, FATAL, ALL
	Message     string    `json:"message,omitempty" form:"message,omitempty"`
	Location    string    `json:"location,omitempty" form:"location,omitempty"` // Ideally filename or file location from the software using the logging service.
	FromTime    time.Time `json:",omitempty" form:"from,omitempty" time_format:"2006-01-02T15:04:05"`
	ToTime      time.Time `json:",omitempty" form:"to,omitempty" time_format:"2006-01-02T15:04:05"`
}

// FilterWithoutCreatedDate compares the values between two log models: The receiver and the comparison.
// If the two models are the same, true is returned. Otherwise, false.
func (logModel *LogModel) FilterWithoutCreatedDate(comparison *LogModel) bool {

	validLogLevel := false
	if logModel.LogLevel != "ALL" {
		for _, value := range logLevels {
			if value == comparison.LogLevel {
				validLogLevel = true
			}
		}
	}

	if (logModel.Severity >= MinSeverity && logModel.Severity <= MaxSeverity) && logModel.Severity != comparison.Severity {
		return false
	}
	if validLogLevel && logModel.LogLevel != comparison.LogLevel {
		return false
	}
	if logModel.Message != "" && logModel.Message != comparison.Message {
		return false
	}
	if logModel.Location != "" && logModel.Location != comparison.Location {
		return false
	}
	if (logModel.CreatedDate.IsZero() && (!logModel.FromTime.IsZero() && !logModel.ToTime.IsZero())) && !comparison.CreatedDateFallsWithinDateRange(logModel.FromTime, logModel.ToTime) {
		return false
	}
	if !logModel.CreatedDate.IsZero() && comparison.CreatedDate != logModel.CreatedDate {
		return false
	}

	return true
}

// Filter compares the values between two log models: The receiver and the comparison.
// If the two models are the same, true is returned. Otherwise, false.
func (logModel *LogModel) Filter(comparison *LogModel) bool {

	if !logModel.CreatedDate.IsZero() && logModel.CreatedDate != comparison.CreatedDate {
		return false
	}
	return logModel.FilterWithoutCreatedDate(comparison)
}

// CreatedDateFallsWithinDateRange checks if the created date of the reciever, falls within the provided date range.
func (logModel *LogModel) CreatedDateFallsWithinDateRange(fromTime time.Time, toTime time.Time) bool {

	if !logModel.CreatedDate.IsZero() && !fromTime.IsZero() && !toTime.IsZero() {
		return logModel.CreatedDate.After(fromTime) && logModel.CreatedDate.Before(toTime)
	}

	return false
}

// WriteLog writes a log to a logfile.
func (logModel *LogModel) WriteLog(*core.FileMutexPool) error {
	logLocation, err := getLogWriteLocation(logModel)
	createLogLevelDirectory(logModel.LogLevel)

	file, err := os.OpenFile(logLocation, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatal(err)
		return err
	}
	defer file.Close()

	bufferedWriter := bufio.NewWriter(file)
	_, err = bufferedWriter.Write(
		buildLogMessage(logModel),
	)
	if err != nil {
		log.Fatal(err)
		return err
	}
	bufferedWriter.Flush()

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
		paths = getLogLevelPaths(logLevels)
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
		if logModel.FilterWithoutCreatedDate(logLine) {
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

	var logModel = new(LogModel)
	logModel.CreatedDate, _ = time.Parse(core.LogDateFormat, logProperties[0])
	logModel.Location = logProperties[1]
	logModel.Severity, _ = strconv.Atoi(logProperties[2])
	logModel.LogLevel = logType
	logModel.Message = rawLog[logTextIndicator+2 : len(rawLog)-1]

	return logModel
}

func createLogLevelDirectory(logLevel string) {
	path := strings.ToUpper(logLevel)
	if _, err := os.Stat(path); os.IsNotExist(err) {
		os.MkdirAll(path, 0700)
	}
}

func getLogWriteLocation(logModel *LogModel) (string, error) {
	return fmt.Sprintf("%s/%s.txt", logModel.LogLevel, time.Now().Format(core.ResourceFileNameDateFormat)), nil
}

func buildLogMessage(logModel *LogModel) []byte {
	location := logModel.Location
	messageText := strings.Replace(logModel.Message, "\n", "\\n", -1)
	return []byte(fmt.Sprintf("[%s]-[%s]-[%d]:\"%s\"\n", time.Now().Format(core.LogDateFormat), location, logModel.Severity, messageText))
}
