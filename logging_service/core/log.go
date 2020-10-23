package log

import (
	"bufio"
	"fmt"
	"log"
	"logging_service/messages"
	"os"
	"sync"
	"time"
)

type SafeCounter struct {
	startIndex   int
	endIndex     int
	currentIndex int
	mux          sync.Mutex
}

type LogCounts struct {
	DEBUG   SafeCounter
	INFO    SafeCounter
	WARNING SafeCounter
	ERROR   SafeCounter
	FATAL   SafeCounter
}

// Log message.
type Log = messages.Log
type Writer int

// Writing

// WriteLog writes a log to a logfile.
func WriteLog(message *Log) {
	logLocation := getLogWriteLocation(message)

	file, err := os.OpenFile(logLocation, os.O_APPEND, 0666)
	if err != nil {
		log.Fatal(err)
	}

	defer file.Close()

	bufferedWriter := bufio.NewWriter(file)
	bytesWritten, err := bufferedWriter.Write(
		buildLogMessage(message),
	)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("Bytes written: %d\n", bytesWritten)
	bufferedWriter.Flush()

}

// WriteLogs writes a bulk of logs to log files.
func WriteLogs(message *Log) {

}

func getLogTypeAsString(logType int) string {
	var logLevels = make(map[int]string)
	logLevels[1] = "DEBUG"
	logLevels[2] = "INFO"
	logLevels[3] = "WARN"
	logLevels[4] = "ERROR"
	logLevels[5] = "FATAL"
	logLevels[6] = "ALL"

	return logLevels[logType]
}

func getLogWriteLocation(message *Log) string {
	return fmt.Sprintf("%s/%s-%d.txt", getLogTypeAsString(int(message.Type)), message.CreatedDate.String(), message.Severity)
}

func buildLogMessage(message *Log) []byte {
	return []byte(fmt.Sprintf("[%s]-[%s]-[%d]: %s", message.CreatedDate.Format(time.RFC3339), message.OriginLocation, message.Severity, message.Message))
}

// Reading
