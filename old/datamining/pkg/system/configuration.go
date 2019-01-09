package system

import (
	"io/ioutil"

	yaml "gopkg.in/yaml.v2"
)

//UnirisConfig describes the uniris robot main configuration
type UnirisConfig struct {
	PublicKey  string                `yaml:"publicKey"`
	PrivateKey string                `yaml:"privateKey"`
	SharedKeys SharedKeys            `yaml:"sharedKeys"`
	Services   ServicesConfiguration `yaml:"services"`
}

//ServicesConfiguration describes the services configuration
type ServicesConfiguration struct {
	Datamining DataMiningConfiguration `yaml:"datamining"`
}

//KeyPair represent a keypair
type KeyPair struct {
	PrivateKey string `yaml:"priv"`
	PublicKey  string `yaml:"pub"`
}

//SharedKeys describes the uniris shared keys
type SharedKeys struct {
	EmKeys []KeyPair `yaml:"em"`
	Robot  KeyPair   `yaml:"robot"`
}

//DataMiningConfiguration describes the datamining configuration
type DataMiningConfiguration struct {
	InternalPort int                `yaml:"internalPort"`
	ExternalPort int                `yaml:"externalPort"`
	Errors       DataMininingErrors `yaml:"errors"`
}

//DataMininingErrors defines the datamining errors
type DataMininingErrors struct {
	AccountNotExist string `yaml:"accountNotExist"`
}

//BuildFromFile creates configuration from configuration file
func BuildFromFile(confFilePath string) (*UnirisConfig, error) {
	bytes, err := ioutil.ReadFile(confFilePath)
	if err != nil {
		return nil, err
	}

	var conf UnirisConfig
	err = yaml.Unmarshal(bytes, &conf)
	if err != nil {
		return nil, err
	}

	return &conf, nil
}
