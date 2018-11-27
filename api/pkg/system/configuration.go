package system

import (
	"io/ioutil"

	yaml "gopkg.in/yaml.v2"
)

//UnirisConfig describes the uniris robot main configuration
type UnirisConfig struct {
	Services ServicesConfiguration `yaml:"services"`
}

//ServicesConfiguration describes the services configuration
type ServicesConfiguration struct {
	API        APIConfiguration        `yaml:"api"`
	Datamining DataMiningConfiguration `yaml:"datamining"`
}

//APIConfiguration describes the api service configuration
type APIConfiguration struct {
	Port int `yaml:"port"`
}

//DataMiningConfiguration describes the datamining configuration
type DataMiningConfiguration struct {
	InternalPort int                `yaml:"internalPort"`
	Errors       DataMininingErrors `yaml:"errors"`
}

//DataMininingErrors defines the datamining errors
type DataMininingErrors struct {
	AccountNotExist string `yaml:"accountNotExist"`
}

//BuildFromFile creates configuration from configuration file
func BuildFromFile(confFilePath string) (conf UnirisConfig, err error) {
	bytes, err := ioutil.ReadFile(confFilePath)
	if err != nil {
		return
	}

	err = yaml.Unmarshal(bytes, &conf)
	if err != nil {
		return
	}

	return conf, nil
}
