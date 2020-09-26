# NAD-A3
Logging service written in Golang with Gin-gonic and data creator client written in Java (3 billion devices btw)
By Conor and Attila MacPherson

## Creating a new route

A new route is easy to add. Go to routes/routes.go.
Add your endpoint to the router in enableRoutes:
```
func enableRoutes(router *Engine) {

	router.GET("/log", handlers.HandleGetLog, nil)
}

```

Add your route handler to the handlers directory.
Make sure the function for your route handler is capitalized.
Go makes things public if their first letter is a capital. Otherwise, it's private. If the function is private, it won't be visible after it's imported.
```
package handlers

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

type Context = gin.Context

// handlers must start with an upper case character.
// Symbols with Capital letters at the start are public. Lowercased is private.
func HandleGetLog(c *Context) {
	c.HTML(http.StatusOK, "index.tmpl.html", nil)
}
```

[Context documentation](https://godoc.org/github.com/gin-gonic/gin#Context)


# Logging

## Format
The log format will be as follows:

`Date \<severity-level>-\<message-number>:\<message-text>`

e.g.:

`Nov 03 2003 21:21:21  1-752:Test Log`

---
## Severity levels
Severity levels allow you to quickly identiy a log that may have critical information.

Severity levels range from 1 (the lowest level) to 7.

---
## Log Storage Format
Logs are stored as plaintext in flat files.

Each severity level has it's own directory at the top level of the log storage directory.

Each log file will hold up to 2000 logs. Once capacity is reached, a new log file will be created under the same directory.

The log file's will have a format of:

`severity<severity-level_<first-message-number>-<last-message-number>.txt`

e.g.:

`severity1_1-2000.txt`, `severity1_2001-4000.txt`, etc.

This will allow quick access to log files via query, but may add some overhead if there are a lot of logs stored because opening files may be slow.
