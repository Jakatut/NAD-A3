package core

import (
	"sync"
)

type Result struct {
	Response *Response
	Err      error
}

type Response struct {
	Data    interface{}
	Message string
}

type HandlerResources struct {
	WaitGroup  *sync.WaitGroup
	LogChannel chan *Result
}

type FileMutexPool struct {
	Pool map[string]sync.RWMutex
	Lock sync.RWMutex
}
