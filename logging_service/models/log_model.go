package models

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"
	"time"
)

// Log Types
const (
	DEBUG = 1
	INFO  = 2
	WARN  = 3
	ERROR = 4
	FATAL = 5
	ALL   = 6
)

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
func (logModel *LogModel) WriteLog() error {
	logLocation := getLogWriteLocation(logModel)

	createLogTypeDirectory(logModel.Type)

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

func createLogTypeDirectory(logType int8) {
	path := getLogTypeAsString(logType)
	if _, err := os.Stat(path); os.IsNotExist(err) {
		os.MkdirAll(path, 0700)
	}
}

func getLogTypeAsString(logType int8) string {
	var logLevels = make(map[int8]string)
	logLevels[1] = "DEBUG"
	logLevels[2] = "INFO"
	logLevels[3] = "WARN"
	logLevels[4] = "ERROR"
	logLevels[5] = "FATAL"
	logLevels[6] = "ALL"

	return logLevels[logType]
}

func getLogWriteLocation(message *LogModel) string {
	return fmt.Sprintf("%s/%s-%d.txt", getLogTypeAsString(message.Type), message.CreatedDate.Format(dateFormat), message.Severity)
}

func buildLogMessage(message *LogModel) []byte {
	location := strings.ReplaceAll(message.Location, "\n", "%0A")
	messageText := strings.ReplaceAll(message.Message, "\n", "%0A")
	return []byte(fmt.Sprintf("[%s]-[%s]-[%d]: %s\n", message.CreatedDate.Format(dateFormat), location, message.Severity, messageText))
}

// Reading

// ReadLog reads a log from the log file.
func (logModel *LogModel) ReadLog() *LogModel {
	if !logModel.CreatedDate.IsZero() {
		logLocation := getLogWriteLocation(logModel)
		readLogWithKnownLocation(logLocation)
	}

	return nil
}

func readLogWithKnownLocation(location string) {

}
