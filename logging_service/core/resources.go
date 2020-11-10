package core

import (
	"sync"
)

// LogDateFormat used when writing content to log files. Includes time.
const LogDateFormat = "2006-01-02T15-04-05"

const CreatedDayFormat = "2006-01-02"
const CreatedTimeFormat = "15-04-05"

// ResourceFileNameDateFormat used in the file name when creating log files
const ResourceFileNameDateFormat = "2006-01-02"

// LogLevels defines all available log level types.
var LogLevels = []string{"DEBUG", "INFO", "WARNING", "ERROR", "FATAL"}

// Response defines an api request's response. This would be used for successful responses. Any responses that
// indicate a failure or error should use errors.New("") for the response.
type Response struct {
	Data    interface{} `json:"data,omitempty"`
	Message string      `json:"message,omitempty"`
	Error   string      `json:"error,omitempty"`
}

// FileMutexPool is a map of strings to Read Write mutexes used to control concurrent access to log files.
type FileMutexPool struct {
	Pool map[string]*sync.RWMutex
	Lock sync.RWMutex
}

// AddMutex adds a new mutex to the pool map with the key fileName
// If the key already exists, nothing happens.
func (fmp *FileMutexPool) addMutex(fileName string) {
	if fmp.Pool == nil {
		fmp.Pool = make(map[string]*sync.RWMutex)
	}
	fmp.Lock.RLock()
	if _, ok := fmp.Pool[fileName]; !ok {
		fmp.Lock.RUnlock()
		fmp.Lock.Lock()
		defer fmp.Lock.Unlock()
		fmp.Pool[fileName] = new(sync.RWMutex)
	} else {
		fmp.Lock.RUnlock()
	}
}

// LockReadFileMutex locks a log file's read mutex.
func (fmp *FileMutexPool) LockReadFileMutex(fileName string) {
	fmp.addMutex(fileName)
	if _, ok := fmp.Pool[fileName]; ok {
		fmp.Pool[fileName].RLock()
	}
}

// UnlockReadFileMutex unlocks a log file's read mutex.
func (fmp *FileMutexPool) UnlockReadFileMutex(fileName string) {
	if _, ok := fmp.Pool[fileName]; ok {
		fmp.Pool[fileName].RUnlock()
	}
}

// LockWriteFileMutex locks a log file's write mutex.
func (fmp *FileMutexPool) LockWriteFileMutex(fileName string) {
	fmp.addMutex(fileName)
	if _, ok := fmp.Pool[fileName]; ok {
		fmp.Pool[fileName].Lock()
	}
}

// UnlockWirterFileMutex unlocks a log file's write mutex.
func (fmp *FileMutexPool) UnlockWirterFileMutex(fileName string) {
	if _, ok := fmp.Pool[fileName]; ok {
		fmp.Pool[fileName].Unlock()
	}
}

// LogTypeCounter keeps track of log counts.
type LogTypeCounter struct {
	Counters map[string]uint
	Lock     sync.RWMutex
}

// AddCount adds 1 to the count for logType.
func (ltc *LogTypeCounter) AddCount(logType string) uint {

	if ltc.Counters == nil {
		ltc.Counters = make(map[string]uint)
	}
	ltc.Lock.RLock()
	if _, ok := ltc.Counters[logType]; ok {
		ltc.Lock.RUnlock()
		ltc.Lock.Lock()
		defer ltc.Lock.Unlock()
		ltc.Counters[logType] = ltc.Counters[logType] + 1
	} else {
		ltc.Lock.RUnlock()
		ltc.Lock.Lock()
		defer ltc.Lock.Unlock()
		ltc.Counters[logType] = 1
	}

	return ltc.Counters[logType]
}

// GetCount returns the current log count for logType.
func (ltc *LogTypeCounter) GetCount(logType string) uint {
	if ltc.Counters == nil {
		ltc.Counters = make(map[string]uint)
		return 0
	}
	ltc.Lock.RLock()
	defer ltc.Lock.RUnlock()
	if value, ok := ltc.Counters[logType]; ok {
		return value
	}

	return 0
}

// SubtractCount removes 1 from the log types value in the counters map.
func (ltc *LogTypeCounter) SubtractCount(logType string) uint {

	if ltc.Counters == nil {
		ltc.Counters = make(map[string]uint)
	}
	ltc.Lock.RLock()
	if _, ok := ltc.Counters[logType]; ok {
		var greaterThanOne = ltc.Counters[logType] > 1
		ltc.Lock.RUnlock()
		ltc.Lock.Lock()
		defer ltc.Lock.Unlock()
		if greaterThanOne {
			ltc.Counters[logType] = ltc.Counters[logType] - 1
		}
	} else {
		defer ltc.Lock.RUnlock()
	}

	return ltc.Counters[logType]
}

func (ltc *LogTypeCounter) resetCounts() {
	for _, logLevel := range LogLevels {
		ltc.Counters[logLevel] = 0
	}
}

// SetStartingCounts sets the log type counters values to either 1 or to the last count.
func (ltc *LogTypeCounter) SetStartingCounts() {

	if ltc.Counters == nil {
		ltc.Counters = make(map[string]uint)
	}
	for _, logLevel := range LogLevels {
		location, err := GetLastLogFileLocation(logLevel)
		if err != nil {
			ltc.Counters[logLevel] = 0
		} else {
			// ltc.Counters[logLevel] =
			GetLastLogId(location)
		}
	}
}
