package core

/*
 *
 * file: 		helpers.go
 * project:		logging_service - NAD-A3
 * programmer: 	Conor Macpherson
 * description: Defines helper functions for things like finiding file locations, validating log level, and extracting log details.
 *
 */

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"time"
)

// CreateLogLevelDirectory creates the directory for storing logs for the level logLevel
//
// Parameters:
//	string	logLevel	- Log level directory to create.
//
func CreateLogLevelDirectory(logLevel string) {
	path := strings.ToUpper("logs/" + logLevel)
	if _, err := os.Stat(path); os.IsNotExist(err) {
		os.MkdirAll(path, 0700)
	}
}

// GetLogWriteLocation finds the location of current log file to send logs to for the provided logLevel.
//
// Parameters:
//	string	logLevel	- Log level directory to create.
//
// Returns
//	string	- The write location for the given log level.
//	error	- Any errors that occur.
//
func GetLogWriteLocation(logLevel string) (string, error) {
	dir := strings.ToUpper("logs/" + logLevel)
	location := strings.ToUpper(dir+"/"+time.Now().Format(ResourceFileNameDateFormat)) + ".txt"
	_, err := os.Stat(dir)
	if os.IsNotExist(err) {
		os.MkdirAll(dir, 0700)
	}
	_, err = os.Stat(location)
	if os.IsNotExist(err) {
		os.Create(location)
	}
	return location, nil
}

// GetLastLogFileLocation gets the last log in the given log level's directory.
//
// Parameters:
//	string	logLevel	- Log level to get the last file for.
//
// Returns
//	string	- The location of the last file added to the logs for the given log level.
//	error	- Any errors that occur.
//
func GetLastLogFileLocation(logLevel string) (string, error) {
	var paths = GetLogLevelPaths([]string{logLevel})
	if len(paths) == 0 {
		return "", errors.New("Does not exist")
	}

	return paths[len(paths)-1], nil
}

// IsValidLogLevel check the provided logLevel is one of "DEBUG", "WARNING", "ERROR", "FATAL", "INFO", or "ALL"
//
// Parameters:
//	string	logLevel	- Log level to get the last file for.
//
// Returns
//	bool - True if the given log level is a valid log level.
func IsValidLogLevel(logLevel string) bool {
	for _, value := range LogLevels {
		if strings.Compare(strings.ToUpper(value), strings.ToUpper(logLevel)) == 0 || strings.Compare(strings.ToUpper(value), "ALL") == 0 {
			return true
		}
	}

	return false
}

// GetLastLogID find the last log under the given log level, and returns that plus 1.
//
// Parameters:
//	string	location	- Location of the logs to check.
//	string	logLevel	- Log level to get the last file for.
//
// Returns
//	uint - the last id of the given log level.
//
func GetLastLogID(location string, logLevel string) uint {
	file, err := os.Open(location)
	if err != nil {
		fmt.Print(err)
	}
	scanner := bufio.NewScanner(file)
	var rawLog string
	for scanner.Scan() {
		rawLog = scanner.Text()
	}

	if rawLog == "" {
		return 1
	}

	details, err := GetLogDetailsFromRawLog(rawLog)
	if err != nil {
		return 1
	}

	return uint(details["id"].(uint64))
}

// GetLogLevelPaths goes through the list of log levels, and finds the path of every log in that level's directory.
//
// Parameters:
//	[]string	logLevels	- Slice of log levels to get the path for.
//
// Returns
//	[]string - Slice of strings containing paths for each provided log level.
//
func GetLogLevelPaths(logLevels []string) []string {
	var paths []string
	for _, level := range logLevels {
		dir := strings.ToUpper("logs/" + level + "/")
		_ = filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
			if path != dir {
				paths = append(paths, path)
			}
			return nil
		})
	}

	return paths
}

// GetSearchFilePaths get a list of file paths for the files that need to be searched.
//
// Parameters:
//	string	logLevel	- The log level to get the search file paths for.
//
// Returns
//	[]string - Slice of strings containing paths for the provided log level.
//
func GetSearchFilePaths(logLevel string) []string {
	var paths []string
	var logLevels = append(LogLevels, "All")
	if logLevel != "ALL" {
		paths = GetLogLevelPaths([]string{logLevel})
	} else {
		paths = GetLogLevelPaths(logLevels)
	}

	return paths
}

// GetLogDetailsFromRawLog extracts the created date, id, location, and message of a raw log.
//
// Parameters:
//	string	rawLog	-	The raw log from the log file.
//
// Returns
//	map[string]interface{}	- Slice of strings containing paths for the provided log level.
//	error					- Any errors that occur.
//
func GetLogDetailsFromRawLog(rawLog string) (map[string]interface{}, error) {

	regex, _ := regexp.Compile("\\w+\\s*=\\s*")
	// Only find the first 3 keys.
	foundKeys := regex.FindAllString(rawLog, 3)

	// Validate the keys.
	for _, value := range foundKeys {
		if value != "date=" && value != "id=" && value != "location=" {
			return nil, errors.New("unknown property")
		}
	}

	dateIndex := strings.Index(rawLog, foundKeys[0])
	idIndex := strings.Index(rawLog, foundKeys[1]) + len(foundKeys[1])
	locationIndex := strings.Index(rawLog, foundKeys[2]) + len(foundKeys[2])

	dateString := rawLog[dateIndex+len(foundKeys[0]) : idIndex-len("date=")]
	dateString = strings.Trim(dateString, "\"")
	createdDate, err := time.Parse(LogDateFormat, dateString)
	if err != nil {
		return nil, err
	}

	idString := rawLog[idIndex+1 : locationIndex-len("location=")-2]
	idString = strings.Trim(idString, "\"")
	id, err := strconv.ParseUint(idString, 0, 64)
	if err != nil {
		return nil, err
	}

	location := rawLog[locationIndex:]
	numberOfEscapedQuotes := strings.Count(location, "\\\"")
	location = strings.Replace(location, "\\\"", "", -1)
	locationEndIndex := strings.Index(location, "\"]:") + numberOfEscapedQuotes + len(rawLog[0:locationIndex-1]) + 1
	location = rawLog[locationIndex:locationEndIndex]
	location = strings.TrimLeft(location, "\"")
	location = strings.Replace(location, "\\\"", "\"", -1)

	message := rawLog[locationEndIndex+3:]
	message = strings.Trim(message, "\"")
	message = strings.Replace(message, "\\\"", "\"", -1)

	details := make(map[string]interface{})
	details["created_date"] = createdDate
	details["id"] = id
	details["location"] = location
	details["message"] = message

	return details, nil
}
