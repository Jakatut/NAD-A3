package handlers

import (
	"fmt"
	"logging_service/messages"
	"net/http"
)

// Gets all logs
func HandleGetLog(c *Context) {
	// fmt.Println("Handle get log")
	// reqBody := &messages.Log{}
	// err := c.Bind(reqBody)
	// if err != nil {
	// 	fmt.Println(err)
	// }
	// fmt.Println(reqBody)
	c.HTML(http.StatusOK, "index.tmpl.html", nil)
}

// Post a log
func HandlePostLog(c *Context) {
	// fmt.Println("Handle get log")
	// reqBody := &messages.Log{}
	// err := c.Bind(reqBody)
	// if err != nil {
	// 	fmt.Println(err)
	// }
	// c.JSON(200, reqBody)
	// fmt.Println(reqBody)
	// c.HTML(http.StatusOK, "index.tmpl.html", nil)\
	c.HTML(http.StatusOK, "index.tmpl.html", nil)
	log := new(messages.Log)
	err := c.BindJSON(log)
	if err != nil {
		fmt.Println("Error binding to log.")
	} else {
		fmt.Println(log)
	}
}
