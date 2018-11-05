package mining

import (
	"github.com/uniris/uniris-core/datamining/pkg/account/mining/checks"
	"github.com/uniris/uniris-core/datamining/pkg/mining"
)

type keychainValidator struct {
	sig checks.Signer
}

type keychainChecker struct {
	masterChecks []checks.Handler
	slaveChecks  []checks.Handler
}

//NewKeychainChecker creates a keychain master checker
func NewKeychainChecker(sig checks.Signer, h checks.TransactionDataHasher) mining.Checker {
	masterChecks := []checks.Handler{
		checks.NewIntegrityChecker(h),
	}

	slaveChecks := []checks.Handler{
		checks.NewSignatureChecker(sig),
	}

	return keychainChecker{
		masterChecks: masterChecks,
		slaveChecks:  slaveChecks,
	}
}

func (c keychainChecker) CheckAsMaster(txHash string, data interface{}) error {
	for _, checks := range c.masterChecks {
		if err := checks.CheckData(data, txHash); err != nil {
			return err
		}
	}

	return nil
}

func (c keychainChecker) CheckAsSlave(txHash string, data interface{}) error {
	for _, checks := range c.slaveChecks {
		if err := checks.CheckData(data, txHash); err != nil {
			return err
		}
	}

	return nil
}
