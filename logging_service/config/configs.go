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
	"errors"
	"io/ioutil"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v2"
)

type Values struct {
	Port          string `yaml:"Port"`
	Auth0Audience string `yaml:"Auth0Audience"`
	Auth0URI      string `yaml:"Auth0URI"`
	LogDirectory  string `yaml:"LogDirectory"`
}

// GetConfig reads and unmarshals a yaml file to a config.Values struct.
//
// Returns
//	Values - Config values
//
func GetConfig() Values {
	configPath := os.Getenv("LOGGING_SERVICE_CONFIG_PATH")
	if configPath == "" {
		panic(errors.New("LOGGING_SERVICE_CONFIG_PATH not set; config required"))
	}
	fileName, err := filepath.Abs(configPath)
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

	config.LogDirectory, err = filepath.Abs(config.LogDirectory)
	if err != nil {
		panic(err)
	}

	config.LogDirectory += string(os.PathSeparator)

	return config
}
