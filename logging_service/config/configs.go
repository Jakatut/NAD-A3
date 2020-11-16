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

func GetConfig() Values {
	fileName, _ := filepath.Abs("config/config.yaml")
	yamlFile, err := ioutil.ReadFile(fileName)
	var config Values
	err = yaml.Unmarshal(yamlFile, &config)
	if err != nil {
		panic(err)
	}

	config.LogDirectory, err = filepath.Abs(config.LogDirectory)
	if err != nil {
		panic(err)
	}

	return config
}
