package main

import (
	"logging_service/routes"

	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.New()
	routes.Setup(router)
}
