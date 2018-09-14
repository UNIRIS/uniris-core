package infrastructure

import (
	"io/ioutil"
)

//GetVersion returns the daemon version
func GetVersion() (string, error) {
	version, err := ioutil.ReadFile("version")
	if err != nil {
		return "", err
	}
	return string(version), nil
}
