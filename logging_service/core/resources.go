package core

import (
	"sync"
)

const LogDateFormat = "2006-01-02T15-04-05"
const ResourceFileNameDateFormat = "2006-01-02"

// Response defines an api request's response. This would be used for successful responses. Any responses that
// indicate a failure or error should use errors.New("") for the response.
type Response struct {
	Data    interface{}
	Message string
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
