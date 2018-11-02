package mining

import (
	"errors"
	"time"

	datamining "github.com/uniris/uniris-core/datamining/pkg"
	biodlisting "github.com/uniris/uniris-core/datamining/pkg/biod/listing"
)

//ErrInvalidTransaction is returned a transaction is invalid
var ErrInvalidTransaction = errors.New("Invalid transaction")

//Checker define methods for a checker to implement
type Checker interface {
	CheckAsMaster(txHash string, data interface{}) error
	CheckAsSlave(txHash string, data interface{}) error
}

//Signer defines methods to handle lead mining signing
type Signer interface {
	LockSigner
	PowSigner
}

//LockSigner define method to sign lock transaction
type LockSigner interface {
	SignLock(lock datamining.TransactionLock, pvKey string) (string, error)
}

//Service defines methods for global mining
type Service interface {
	LeadMining(txHash string, addr string, data interface{}, txType datamining.TransactionType, biodSig string) error
	Validate(txHash string, data interface{}, txType datamining.TransactionType) (datamining.Validation, error)
}

type service struct {
	notif      datamining.Notifier
	poolF      datamining.PoolFinder
	poolR      datamining.PoolRequester
	signer     Signer
	biodLister biodlisting.Service
	robotKey   string
	robotPvKey string
	checks     map[datamining.TransactionType]Checker
}

//NewService creates a new global mining service
func NewService() Service {
	return service{}
}

func (s service) LeadMining(txHash string, addr string, data interface{}, txType datamining.TransactionType, biodSig string) error {
	if err := s.notif.NotifyTransactionStatus(txHash, datamining.TxPending); err != nil {
		return err
	}

	lastVPool, vPool, sPool, err := s.findPools(addr)
	if err != nil {
		return err
	}

	if err := s.requestLock(txHash, addr, lastVPool); err != nil {
		return err
	}

	masterValid, valids, err := s.mine(txHash, data, biodSig, lastVPool, vPool, txType)
	if err != nil {
		if err == ErrInvalidTransaction {
			if err := s.notif.NotifyTransactionStatus(txHash, datamining.TxInvalid); err != nil {
				return err
			}
			return nil
		}
		return err
	}

	if err := s.notif.NotifyTransactionStatus(txHash, datamining.TxApproved); err != nil {
		return err
	}

	endorsement := datamining.NewEndorsement(time.Now(), txHash, masterValid, valids)
	if err := s.poolR.RequestStorage(sPool, data, endorsement, txType); err != nil {
		return err
	}
	if err := s.notif.NotifyTransactionStatus(txHash, datamining.TxReplicated); err != nil {
		return err
	}

	return s.requestUnlock(txHash, addr, lastVPool)
}

func (s service) findPools(addr string) (lastVPool datamining.Pool, vPool datamining.Pool, sPool datamining.Pool, err error) {
	lastVPool, err = s.poolF.FindLastValidationPool(addr)
	if err != nil {
		return
	}

	vPool, err = s.poolF.FindValidationPool()
	if err != nil {
		return
	}

	sPool, err = s.poolF.FindStoragePool()
	if err != nil {
		return
	}

	return
}

func (s service) requestLock(txHash string, addr string, lastVPool datamining.Pool) error {
	//Build lock transaction
	lock := datamining.TransactionLock{TxHash: txHash, MasterRobotKey: s.robotKey, Address: addr}
	sigLock, err := s.signer.SignLock(lock, s.robotPvKey)
	if err != nil {
		return err
	}

	if err := s.poolR.RequestLock(lastVPool, lock, sigLock); err != nil {
		return err
	}
	if err := s.notif.NotifyTransactionStatus(txHash, datamining.TxLocked); err != nil {
		return err
	}

	return err
}

func (s service) requestUnlock(txHash string, addr string, lastVPool datamining.Pool) error {
	//Build unlock transaction
	lock := datamining.TransactionLock{TxHash: txHash, MasterRobotKey: s.robotKey, Address: addr}
	sigLock, err := s.signer.SignLock(lock, s.robotPvKey)
	if err != nil {
		return err
	}

	if err := s.poolR.RequestUnlock(lastVPool, lock, sigLock); err != nil {
		return err
	}
	if err := s.notif.NotifyTransactionStatus(txHash, datamining.TxUnlocked); err != nil {
		return err
	}

	return err
}

func (s service) mine(txHash string, data interface{}, biodSig string, lastVPool, vPool datamining.Pool, txType datamining.TransactionType) (datamining.MasterValidation, []datamining.Validation, error) {
	//Execute transaction specific master checks
	if err := s.checks[txType].CheckAsMaster(txHash, data); err != nil {
		return nil, nil, ErrInvalidTransaction
	}

	//Execute the Proof of Work
	masterValid, err := NewPOW(s.biodLister, s.signer, s.robotKey, s.robotPvKey).Execute(txHash, biodSig, lastVPool)
	if err != nil {
		return nil, nil, err
	}

	//Ask a pool of peers to validate the transaction
	valids, err := s.poolR.RequestValidations(vPool, data, txType)
	if err != nil {
		return nil, nil, err
	}

	//Check if the validations passed
	for _, v := range valids {
		if v.Status() == datamining.ValidationKO {
			return nil, nil, ErrInvalidTransaction
		}
	}

	return masterValid, valids, nil
}

func (s service) Validate(txHash string, data interface{}, txType datamining.TransactionType) (datamining.Validation, error) {
	if err := s.checks[txType].CheckAsSlave(txHash, data); err != nil {
		if err == ErrInvalidTransaction {
			return s.buildSlaveValidation(datamining.ValidationKO)
		}
		return nil, err
	}
	return s.buildSlaveValidation(datamining.ValidationOK)
}

func (s service) buildSlaveValidation(status datamining.ValidationStatus) (valid datamining.Validation, err error) {
	v := UnsignedValidation{
		PublicKey: s.robotKey,
		Status:    status,
		Timestamp: time.Now(),
	}
	signature, err := s.signer.SignValidation(v, s.robotPvKey)
	if err != nil {
		return
	}
	return datamining.NewValidation(v.Status, v.Timestamp, v.PublicKey, signature), nil
}
