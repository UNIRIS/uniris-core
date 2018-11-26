package system

import (
	"io/ioutil"
	"sort"

	yaml "gopkg.in/yaml.v2"
)

//UnirisConfig describes the uniris robot main configuration
type UnirisConfig struct {
	Services   ServicesConfiguration `yaml:"services"`
	SharedKeys SharedKeys            `yaml:"sharedKeys"`
}

//ServicesConfiguration describes the services configuration
type ServicesConfiguration struct {
	API        APIConfiguration        `yaml:"api"`
	Datamining DataMiningConfiguration `yaml:"datamining"`
}

//KeyPair represent a keypair
type KeyPair struct {
	PrivateKey string `yaml:"priv"`
	PublicKey  string `yaml:"pub"`
}

//SharedKeys describes the uniris shared keys
type SharedKeys struct {
	Emitter []KeyPair `yaml:"em"`
	Robot   KeyPair   `yaml:"robot"`
}

//SortEmitterKeys sorts the emitter shared keys by their public key
func (sh *SharedKeys) SortEmitterKeys() {
	sort.Slice(sh.Emitter, func(i, j int) bool {
		return sh.Emitter[i].PublicKey < sh.Emitter[j].PublicKey
	})
}

//EmitterRequestKey returns the shared emitter key for the request
func (sh SharedKeys) EmitterRequestKey() KeyPair {
	return sh.Emitter[0]
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

	conf.SharedKeys.SortEmitterKeys()

	return conf, nil
}
