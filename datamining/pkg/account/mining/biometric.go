package master

import "github.com/uniris/uniris-core/datamining/pkg/mining"

type bioChecker struct {
}

//NewBiometricChecker creates a biometric master checker
func NewBiometricChecker() mining.Checker {
	return bioChecker{}
}

func (c bioChecker) CheckAsMaster(txHash string, data interface{}) error {
	return nil
}

func (c bioChecker) CheckAsSlave(txHash string, data interface{}) error {
	return nil
}
