package core

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
func CreateLogLevelDirectory(logLevel string) {
	path := strings.ToUpper("logs/" + logLevel)
	if _, err := os.Stat(path); os.IsNotExist(err) {
		os.MkdirAll(path, 0700)
	}
}

// GetLogWriteLocation finds the location of current log file to send logs to for the provided logLevel.
func GetLogWriteLocation(logLevel string) (string, error) {
	location := fmt.Sprintf("logs/%s/%s.txt", logLevel, time.Now().Format(ResourceFileNameDateFormat))
	dir := "logs/" + logLevel
	_, err := os.Stat("logs/" + logLevel + "/")
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
func GetLastLogFileLocation(logLevel string) (string, error) {
	var paths = GetLogLevelPaths([]string{logLevel})
	if len(paths) == 0 {
		return "", errors.New("Does not exist")
	}

	return paths[len(paths)-1], nil
}

// IsValidLogLevel check the provided logLevel is one of "DEBUG", "WARNING", "ERROR", "FATAL", "INFO", or "ALL"
func IsValidLogLevel(logLevel string) bool {
	for _, value := range LogLevels {
		if strings.Compare(strings.ToUpper(value), strings.ToUpper(logLevel)) == 0 || strings.Compare(strings.ToUpper(value), "ALL") == 0 {
			return true
		}
	}

	return false
}

// GetLastLogID find the last log under the given log level, and returns that plus 1.
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

	details, err := GetLogDetailsFromRawLog(rawLog, logLevel)
	if err != nil {
		return 1
	}

	return uint(details["id"].(uint64))
}

// GetLogLevelPaths goes through the list of log levels, and finds the path of every log in that level's directory.
func GetLogLevelPaths(logLevels []string) []string {
	var paths []string
	for _, level := range logLevels {
		_ = filepath.Walk("logs/"+level+"/", func(path string, info os.FileInfo, err error) error {
			if path != "logs/"+level+"/" {
				paths = append(paths, path)
			}
			return nil
		})
	}

	return paths
}

// GetSearchFilePaths get a list of file paths for the files that need to be searched.
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

// GetLogDetailsFromRawLog extracts the dat, id, location, and message of a raw log.
func GetLogDetailsFromRawLog(rawLog string, logType string) (map[string]interface{}, error) {

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
