package main

import (
	"logging_service/routes"
	"sync"

	"github.com/gin-gonic/gin"
)

var debugWaitGroup sync.WaitGroup
var warningWaitGroup sync.WaitGroup
var infoWaitGroup sync.WaitGroup
var errorWaitGroup sync.WaitGroup
var fatalWaitGroup sync.WaitGroup

func init() {
	debugWaitGroup = sync.WaitGroup{}
	warningWaitGroup = sync.WaitGroup{}
	infoWaitGroup = sync.WaitGroup{}
	errorWaitGroup = sync.WaitGroup{}
	fatalWaitGroup = sync.WaitGroup{}
}

func main() {
	router := gin.New()
	routes.Setup(router, &debugWaitGroup, &warningWaitGroup, &infoWaitGroup, &errorWaitGroup, &fatalWaitGroup)
}
