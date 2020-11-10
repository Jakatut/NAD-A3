package core

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

func CreateLogLevelDirectory(logLevel string) {
	path := strings.ToUpper("logs/" + logLevel)
	if _, err := os.Stat(path); os.IsNotExist(err) {
		os.MkdirAll(path, 0700)
	}
}

func GetLogWriteLocation(logLevel string) (string, error) {
	location := fmt.Sprintf("logs/%s/%s.txt", logLevel, time.Now().Format(ResourceFileNameDateFormat))
	_, err := os.Stat(location)
	if os.IsNotExist(err) {
		return "", err
	}

	return location, nil
}

func GetLastLogFileLocation(logLevel string) (string, error) {
	var paths = GetLogLevelPaths([]string{logLevel})
	if len(paths) == 0 {
		return "", errors.New("Does not exist")
	}

	for _, value := range paths {
		fmt.Println(value)
	}

	return "", nil
}

func IsValidLogLevel(logLevel string) bool {
	for _, value := range LogLevels {
		if strings.Compare(strings.ToUpper(value), strings.ToUpper(logLevel)) == 0 || strings.Compare(strings.ToUpper(value), "ALL") == 0 {
			return true
		}
	}

	return false
}

func GetLastLogId(location string) uint {
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

	logTextIndicator := strings.Index(rawLog, ":")
	// Remove leading and trailing braces, removes the content of the log, and splits the details.
	logProperties := strings.Split(rawLog[1:logTextIndicator-1], "]-[")
	id, err := strconv.Atoi(logProperties[2])
	if err != nil {
		return 1
	}

	return uint(id)
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
