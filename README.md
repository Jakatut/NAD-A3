# NAD-A3
Logging service written in Golang with Gin-gonic and a Client written in C# with xaml
By Conor and Attila MacPherson

## Run it

Set an environment variable LOGGING_SERVICE_CONFIG_PATH pointing to your yaml config.
```
export LOGGING_SERVICE_CONFIG_PATH=config/config.yaml
```

Your config.yaml should have 4 items:
```
Port: 8080
Auth0Audience: eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9
Auth0URI: https://ca-logging.us.auth0.com/
LogDirectory: LOGS/
```

Linux/Mac:
```
make build
bin/logging_service
```

(idk how to run on windows lol)
Windows: 
```
go run main.go
```




## Debugging

### Create your launch.json

Under .vscode, if it does not already exists, write your launch.json file:

```
{
	"version": "0.2.0",
	"configurations": [
		{
			"name": "Launch",
			"type": "go",
			"request": "launch",
			"mode": "auto",
			"program": "${workspaceFolder}/logging_service",
			"env": {},
			"args": []
		}
	]
}
```

### Install dlv:

!!!!! OUTSIDE OF THE PROJECT, RUN THIS GO COMMAND !!!!!

```
go get github.com/go-delve/delve/cmd/dlv
```

Make sure you have a GOPATH setup.

The GOPATH environment variable specifies the location of your workspace. It defaults to a directory named go inside your home directory, so $HOME/go on Unix, $home/go on Plan 9, and %USERPROFILE%\go (usually C:\Users\YourName\go) on Windows.

Under the debug tab in vscode, click launch. You can set breakpoints and watch variables in this tab.