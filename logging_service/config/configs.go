package config

/*
 *
 * file: 		jwt_auth.go
 * project:		logging_service - NAD-A3
 * programmer: 	Conor Macpherson
 * description: Defines the functions used for reading config values.
 *
 */

import (
	"io/ioutil"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v2"
)

type Values struct {
	Server   server   `yaml:"Server"`
	IO       io       `yaml:"IO"`
	Auth     auth     `yaml:"Auth"`
	Database database `yaml:"Database"`
	Results  results  `yaml:"Results"`
}

type server struct {
	Port string `yaml:"PORT"`
}

type io struct {
	LogDirectory string `yaml:"LOG_DIRECTORY"`
}

type auth struct {
	Auth0Audience string `yaml:"AUTH_0_AUDIENCE"`
	Auth0URI      string `yaml:"AUTH_0_URI"`
}

type database struct {
	DatabaseUsername string `yaml:"DATABASE_USERNAME"`
	DatabasePassword string `yaml:"DATABASE_PASSWORD"`
	DatabaseName     string `yaml:"DATABASE_NAME"`
	DatabaseURL      string `yaml:"DATABASE_URL"`
}

type results struct {
	Limit int64 `yaml:"LIMIT"`
}

// GetConfig reads and unmarshals a yaml file to a config.Values struct.
//
// Returns
//	Values - Config values
//
func GetConfig() Values {
	// configPath := os.Getenv("LOGGING_SERVICE_CONFIG_PATH")
	// if configPath == "" {
	// 	panic(errors.New("LOGGING_SERVICE_CONFIG_PATH not set; config required"))
	// }
	fileName, err := filepath.Abs("config/config.yaml")
	if err != nil {
		panic(err)
	}

	yamlFile, err := ioutil.ReadFile(fileName)
	if err != nil {
		panic(err)
	}
	var config Values
	err = yaml.Unmarshal(yamlFile, &config)
	if err != nil {
		panic(err)
	}

	config.IO.LogDirectory, err = filepath.Abs(config.IO.LogDirectory)
	if err != nil {
		panic(err)
	}

	config.IO.LogDirectory += string(os.PathSeparator)

	return config
}
