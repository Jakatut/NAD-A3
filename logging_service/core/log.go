package core

import (
	"bufio"
	"fmt"
	"log"
	models "logging_service/models"
	"os"
	"strings"
)

// Log message.
type Log = models.LogModel

type Writer int

const dateFormat = "January-01-1-15:4:5"

// Writing

// WriteLog writes a log to a logfile.
func WriteLog(message *models.LogModel) error {
	logLocation := getLogWriteLocation(message)

	createLogTypeDirectory(message.Type)

	file, err := os.OpenFile(logLocation, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatal(err)
		return err
	}

	defer file.Close()

	bufferedWriter := bufio.NewWriter(file)
	bytesWritten, err := bufferedWriter.Write(
		buildLogMessage(message),
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

func getLogWriteLocation(message *models.LogModel) string {
	return fmt.Sprintf("%s/%s-%d.txt", getLogTypeAsString(message.Type), message.CreatedDate.Format(dateFormat), message.Severity)
}

func buildLogMessage(message *models.LogModel) []byte {
	location := strings.ReplaceAll(message.Location, "\n", "%0A")
	messageText := strings.ReplaceAll(message.Message, "\n", "%0A")
	return []byte(fmt.Sprintf("[%s]-[%s]-[%d]: %s\n", message.CreatedDate.Format(dateFormat), location, message.Severity, messageText))
}

// Reading

func ReadLog(message *models.LogModel) *models.LogModel {
	if !message.CreatedDate.IsZero() {
		logLocation := getLogWriteLocation(message)
		readLogWithKnownLocation(logLocation)
	}

	return nil
}

func readLogWithKnownLocation(location string) {

}
