package main

import (
	_ "github.com/heroku/x/hmetrics/onload"
	"logging_service/routes"
)

func main() {
	routes.Setup()
}
