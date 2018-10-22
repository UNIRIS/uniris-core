package system

import (
	"io/ioutil"
	"os"
	"strconv"

	yaml "gopkg.in/yaml.v2"
)

//UnirisConfig describes the uniris robot main configuration
type UnirisConfig struct {
	API        APIConfiguration        `yaml:"api"`
	SharedKeys SharedKeys              `yaml:"sharedKeys"`
	Datamining DataMiningConfiguration `yaml:"datamining"`
}

//SharedKeys describes the uniris shared keys
type SharedKeys struct {
	BiodPublicKey   string `yaml:"biodPublicKey"`
	RobotPrivateKey string `yaml:"robotPrivateKey"`
	RobotPublicKey  string `yaml:"robotPublicKey"`
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

//BuildFromEnv creates configurtion from env variables
func BuildFromEnv() (*UnirisConfig, error) {

	apiport := os.Getenv("UNIRIS_API_PORT")
	_apiport, err := strconv.Atoi(apiport)
	if err != nil {
		return nil, err
	}

	intDataminingPort := os.Getenv("UNIRIS_DATAMINING_INTERNAL_PORT")
	_intDataminingport, err := strconv.Atoi(intDataminingPort)
	if err != nil {
		return nil, err
	}

	dataminingErrAccountNotExist := os.Getenv("UNIRIS_DATAMINING_ERROR_ACCOUNT_NOT_EXIST")

	sharedBiodPublicKey := os.Getenv("UNIRIS_SHARED_KEYS_BIOD_PUBLIC_KEY")
	sharedRobotPublicKey := os.Getenv("UNIRIS_SHARED_KEYS_ROBOT_PUBLIC_KEY")
	sharedRobotPrivateKey := os.Getenv("UNIRIS_SHARED_KEYS_ROBOT_PRIVATE_KEY")

	return &UnirisConfig{
		API: APIConfiguration{
			Port: _apiport,
		},
		Datamining: DataMiningConfiguration{
			InternalPort: _intDataminingport,
			Errors: DataMininingErrors{
				AccountNotExist: dataminingErrAccountNotExist,
			},
		},
		SharedKeys: SharedKeys{
			BiodPublicKey:   sharedBiodPublicKey,
			RobotPublicKey:  sharedRobotPublicKey,
			RobotPrivateKey: sharedRobotPrivateKey,
		},
	}, nil
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
