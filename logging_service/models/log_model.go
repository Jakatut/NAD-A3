package models

import (
	"bufio"
	"errors"
	"fmt"
	"log"
	"logging_service/core"
	"os"
	"path/filepath"
	"strings"
	"sync"
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

var logLevels = []string{"DEBUG", "INFO", "WARN", "ERROR", "FATAL"}

// LogModel defines the contents of a log
type LogModel struct {
	CreatedDate time.Time `json:"created_date,omitempty"`
	Severity    int8      `json:"severity"` // Severity levels are 1-7 (lowest to highest)
	Type        int8      `json:"type"`     // DEBUG, INFO, WARN, ERROR, FATAL, ALL
	Message     string    `json:"message"`
	Location    string    `json:"location"` // Ideally filename or file location from the software using the logging service.
}

const dateFormat = "January-01-1-15:4:5"

// Writing

// WriteLog writes a log to a logfile.
func (logModel *LogModel) WriteLog(*core.FileMutexPool) error {
	logLocation, err := getLogWriteLocation(logModel)

	createLogLevelDirectory(logModel.Type)

	file, err := os.OpenFile(logLocation, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatal(err)
		return err
	}

	defer file.Close()

	bufferedWriter := bufio.NewWriter(file)
	bytesWritten, err := bufferedWriter.Write(
		buildLogMessage(logModel),
	)
	if err != nil {
		log.Fatal(err)
		return err
	}
	log.Printf("Bytes written: %d\n", bytesWritten)
	bufferedWriter.Flush()

	return nil
}

// Reading

// ReadLog reads a log from the log file.
func (logModel *LogModel) ReadLog(mutexPool *core.FileMutexPool) (*LogModel, error) {

	var logs []LogModel

	if !logModel.CreatedDate.IsZero() {
		logLocations := getFileSearchLocations(logModel)
		for _, location := range logLocations {
			mutexPool.Lock.RLock()
			if _, ok := mutexPool.Pool[location]; ok {

			} else {
				mutexPool.Lock.Lock()
				mutexPool.Pool[location] = sync.RWMutex{}
				mutexPool.Lock.Lock()
			}
			if mutex, ok := mutexPool.Pool[location]; ok {
				mutex.RLock()
			}
			logs = append(logs, searchLog(location, logModel)...)
			if mutex, ok := mutexPool.Pool[location]; ok {
				mutex.RUnlock()
			}
			mutexPool.Lock.RUnlock()
		}
	}

	return nil, nil
}

func getFileSearchLocations(logModel *LogModel) []string {
	var logLevel, err = getLogLevelAsString(logModel.Type)
	if err != nil {
		return nil
	}

	var paths []string
	if logLevel != "ALL" {
		paths = walkLogLevelPaths([]string{logLevel})
	} else {
		paths = walkLogLevelPaths(logLevels)
	}
}

func walkLogLevelPaths(logLevels []string) []string {
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

func searchLog(location string, logModel *LogModel) []LogModel {

	return nil
}

/*
 *"DEBUG"
"INFO"
"WARN"
"ERROR"
"FATAL"
 *
*/

func createLogLevelDirectory(logLevel int8) {
	path, err := getLogLevelAsString(logLevel)
	if err != nil {

	}
	if _, err := os.Stat(path); os.IsNotExist(err) {
		os.MkdirAll(path, 0700)
	}
}

func getLogLevelAsString(logLevel int8) (string, error) {
	var logLevels = make(map[int8]string)
	logLevels[0] = "ALL"
	logLevels[1] = "DEBUG"
	logLevels[2] = "INFO"
	logLevels[3] = "WARN"
	logLevels[4] = "ERROR"
	logLevels[5] = "FATAL"

	if logLevel <= 0 || logLevel >= 6 {
		return "", errors.New("logLevel must be 0 to 5")
	}

	return logLevels[logLevel], nil
}

func getLogWriteLocation(message *LogModel) (string, error) {
	logLevel, err := getLogLevelAsString(message.Type)
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("%s/%s-%d.txt", logLevel, message.CreatedDate.Format(dateFormat), message.Severity), nil
}

func buildLogMessage(message *LogModel) []byte {
	location := strings.ReplaceAll(message.Location, "\n", "%0A")
	messageText := strings.ReplaceAll(message.Message, "\n", "%0A")
	return []byte(fmt.Sprintf("[%s]-[%s]-[%d]: %s\n", message.CreatedDate.Format(dateFormat), location, message.Severity, messageText))
}
