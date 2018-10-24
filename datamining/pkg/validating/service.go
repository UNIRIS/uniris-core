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
	ValidateWalletData(w *datamining.WalletData, txHash string) (datamining.Validation, error)
	ValidateBioData(b *datamining.BioData, txHash string) (datamining.Validation, error)
	LockTransaction(txLock TransactionLock, sig string) error
	UnlockTransaction(txLock TransactionLock, sig string) error
}

type service struct {
	bioChecks  []checkers.BioDataChecker
	dataChecks []checkers.WalletDataChecker
	robotKey   string
	robotPvKey string
	sig        Signer
	lock       TransactionLocker
}

//NewService creates a approving service
func NewService(sig Signer, lock TransactionLocker, robotKey, robotPvKey string) Service {
	bioChecks := make([]checkers.BioDataChecker, 0)
	dataChecks := make([]checkers.WalletDataChecker, 0)

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

func (s service) ValidateWalletData(w *datamining.WalletData, txHash string) (valid datamining.Validation, err error) {
	for _, c := range s.dataChecks {
		err = c.CheckWalletData(w)
		if err != nil {
			return
		}
	}
	v := Validation{
		PublicKey: s.robotKey,
		Status:    datamining.ValidationOK,
		Timestamp: time.Now(),
	}
	signature, err := s.sig.SignValidation(v, s.robotPvKey)
	if err != nil {
		return
	}
	valid = datamining.NewValidation(
		v.Status,
		v.Timestamp,
		v.PublicKey,
		signature)
	return
}

func (s service) ValidateBioData(bw *datamining.BioData, txHash string) (valid datamining.Validation, err error) {
	for _, c := range s.bioChecks {
		err = c.CheckBioData(bw)
		if err != nil {
			return
		}
	}
	v := Validation{
		PublicKey: s.robotKey,
		Status:    datamining.ValidationOK,
		Timestamp: time.Now(),
	}
	signature, err := s.sig.SignValidation(v, s.robotPvKey)
	if err != nil {
		return
	}
	valid = datamining.NewValidation(
		v.Status,
		v.Timestamp,
		v.PublicKey,
		signature)
	return
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
