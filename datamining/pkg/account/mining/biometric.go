package mining

import (
	"github.com/uniris/uniris-core/datamining/pkg/account/mining/checks"
	"github.com/uniris/uniris-core/datamining/pkg/mining"
)

type bioChecker struct {
	masterChecks []checks.Handler
	slaveChecks  []checks.Handler
}

//NewBiometricChecker creates a biometric master checker
func NewBiometricChecker(sig checks.Signer, h checks.TransactionDataHasher) mining.Checker {
	masterChecks := []checks.Handler{
		checks.NewIntegrityChecker(h),
	}

	slaveChecks := []checks.Handler{
		checks.NewSignatureChecker(sig),
	}

	return bioChecker{
		masterChecks: masterChecks,
		slaveChecks:  slaveChecks,
	}
}

func (c bioChecker) CheckAsMaster(txHash string, data interface{}) error {
	for _, checks := range c.masterChecks {
		if err := checks.CheckData(data, txHash); err != nil {
			return err
		}
	}

	return nil
}

func (c bioChecker) CheckAsSlave(txHash string, data interface{}) error {
	for _, checks := range c.slaveChecks {
		if err := checks.CheckData(data, txHash); err != nil {
			return err
		}
	}

	return nil
}
