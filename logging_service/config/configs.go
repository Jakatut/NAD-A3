package config

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
