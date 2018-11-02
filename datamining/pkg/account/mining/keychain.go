package master

import (
	"github.com/uniris/uniris-core/datamining/pkg/mining"
)

type keychainChecker struct {
}

//NewKeychainChecker creates a keychain master checker
func NewKeychainChecker() mining.Checker {
	return keychainChecker{}
}

func (c keychainChecker) CheckAsMaster(txHash string, data interface{}) error {
	return nil
}

func (c keychainChecker) CheckAsSlave(txHash string, data interface{}) error {
	return nil
}
