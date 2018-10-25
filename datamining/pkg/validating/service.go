package validating

import (
	"errors"
	"time"

	"github.com/uniris/uniris-core/datamining/pkg"
	"github.com/uniris/uniris-core/datamining/pkg/validating/checkers"
)

//ErrLockExisting is returned when a lock already exist
var ErrLockExisting = errors.New("A lock already exist for this transaction")

//Validation represents a validation before its signature
type Validation struct {
	Status    datamining.ValidationStatus `json:"status"`
	Timestamp time.Time                   `json:"timestamp"`
	PublicKey string                      `json:"pubk"`
}

//TransactionLock defines a transaction lock
type TransactionLock struct {
	TxHash         string `json:"tx_hash"`
	MasterRobotKey string `json:"master_robot_key"`
}

//TransactionLocker define methods to lock transaction from master peer
type TransactionLocker interface {
	Lock(TransactionLock) error
	Unlock(TransactionLock) error
	ContainsLock(TransactionLock) bool
}

//Signer defines methods to handle signatures
type Signer interface {
	SignValidation(v Validation, pvKey string) (string, error)
	checkers.Signer
}

//Service is the interface that provide methods for wallets validation
type Service interface {
	ValidateWalletData(w *datamining.WalletData) (datamining.Validation, error)
	ValidateBioData(b *datamining.BioData) (datamining.Validation, error)
	LockTransaction(txLock TransactionLock, sig string) error
	UnlockTransaction(txLock TransactionLock, sig string) error
}

type service struct {
	bioChecks  []checkers.Checker
	dataChecks []checkers.Checker
	robotKey   string
	robotPvKey string
	sig        Signer
	lock       TransactionLocker
}

//NewService creates a approving service
func NewService(sig Signer, lock TransactionLocker, robotKey, robotPvKey string) Service {
	bioChecks := make([]checkers.Checker, 0)
	dataChecks := make([]checkers.Checker, 0)

	bioChecks = append(bioChecks, checkers.NewSignatureChecker(sig))
	dataChecks = append(dataChecks, checkers.NewSignatureChecker(sig))

	return &service{
		bioChecks:  bioChecks,
		dataChecks: dataChecks,
		robotKey:   robotKey,
		robotPvKey: robotPvKey,
		sig:        sig,
		lock:       lock,
	}
}

func (s service) ValidateWalletData(w *datamining.WalletData) (valid datamining.Validation, err error) {
	for _, c := range s.dataChecks {
		err = c.CheckData(w)
		if err != nil {
			if c.IsCatchedError(err) {
				return s.buildValidation(datamining.ValidationKO)
			}
			return
		}
	}
	return s.buildValidation(datamining.ValidationOK)
}

func (s service) ValidateBioData(bw *datamining.BioData) (valid datamining.Validation, err error) {
	for _, c := range s.bioChecks {
		err = c.CheckData(bw)
		if err != nil {
			if c.IsCatchedError(err) {
				return s.buildValidation(datamining.ValidationKO)
			}
			return
		}
	}
	return s.buildValidation(datamining.ValidationOK)
}

func (s service) LockTransaction(txLock TransactionLock, sig string) error {
	if err := s.sig.CheckSignature(s.robotKey, txLock, sig); err != nil {
		return err
	}

	if s.lock.ContainsLock(txLock) {
		return ErrLockExisting
	}

	return s.lock.Lock(txLock)
}

func (s service) UnlockTransaction(txLock TransactionLock, sig string) error {
	if err := s.sig.CheckSignature(s.robotKey, txLock, sig); err != nil {
		return err
	}

	return s.lock.Unlock(txLock)
}

func (s service) buildValidation(status datamining.ValidationStatus) (valid datamining.Validation, err error) {
	v := Validation{
		PublicKey: s.robotKey,
		Status:    status,
		Timestamp: time.Now(),
	}
	signature, err := s.sig.SignValidation(v, s.robotPvKey)
	if err != nil {
		return
	}
	return datamining.NewValidation(
		v.Status,
		v.Timestamp,
		v.PublicKey,
		signature), nil
}
