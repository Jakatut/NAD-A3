package core

import (
	"fmt"
	"os"
	"strings"
	"time"
)

func CreateLogLevelDirectory(logLevel string) {
	path := strings.ToUpper(logLevel)
	if _, err := os.Stat(path); os.IsNotExist(err) {
		os.MkdirAll(path, 0700)
	}
}

func GetLogWriteLocation(logLevel string) (string, error) {
	return fmt.Sprintf("%s/%s.txt", logLevel, time.Now().Format(ResourceFileNameDateFormat)), nil
}
