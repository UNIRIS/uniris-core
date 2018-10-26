package mining

import (
	"errors"
	"time"

	datamining "github.com/uniris/uniris-core/datamining/pkg"
	"github.com/uniris/uniris-core/datamining/pkg/listing"
	"github.com/uniris/uniris-core/datamining/pkg/mining/lock"
	"github.com/uniris/uniris-core/datamining/pkg/mining/pool"
	"github.com/uniris/uniris-core/datamining/pkg/mining/transactions"
	"github.com/uniris/uniris-core/datamining/pkg/mining/validations"
)

//ErrLockExisting is returned when a lock already exist
var ErrLockExisting = errors.New("A lock already exist for this transaction")

//PoolDispatcher wraps peer cluster requester
type PoolDispatcher interface {
	lock.Locker
	transactions.Requester
}

//Signer defines methods to handle signatures
type Signer interface {
	SignLock(lock lock.TransactionLock, pvKey string) (string, error)
	SignValidation(v Validation, pvKey string) (string, error)
	PowSigner
	validations.Signer
}

//Service defines leading methods
type Service interface {
	Lead(txHash string, addr string, bioSig string, data interface{}, txType transactions.Type) error
	Validate(data interface{}, txType transactions.Type) (datamining.Validation, error)
	LockTransaction(txLock lock.TransactionLock) error
	UnlockTransaction(txLock lock.TransactionLock) error
}

type service struct {
	list       listing.Service
	poolF      pool.Finder
	poolD      PoolDispatcher
	notif      Notifier
	locker     lock.TransactionLocker
	sig        Signer
	h          PreviousDataHasher
	robotKey   string
	robotPvKey string

	checks     map[transactions.Type][]validations.Handler
	txHandlers map[transactions.Type]transactions.Handler
}

//NewService creates a new leading service
func NewService(list listing.Service, pF pool.Finder, pD PoolDispatcher, locker lock.TransactionLocker, notif Notifier, sig Signer, h PreviousDataHasher, rPb, rPv string) Service {

	checks := map[transactions.Type][]validations.Handler{
		transactions.CreateWallet: []validations.Handler{
			validations.NewSignatureValidation(sig),
		},
		transactions.CreateBio: []validations.Handler{
			validations.NewSignatureValidation(sig),
		},
	}

	txHandlers := map[transactions.Type]transactions.Handler{
		transactions.CreateWallet: transactions.NewCreateWalletHandler(),
		transactions.CreateBio:    transactions.NewCreateBioHandler(),
	}

	return &service{
		list:       list,
		poolF:      pF,
		poolD:      pD,
		locker:     locker,
		notif:      notif,
		sig:        sig,
		h:          h,
		robotKey:   rPb,
		robotPvKey: rPv,
		txHandlers: txHandlers,
		checks:     checks,
	}
}

func (s service) Lead(txHash string, addr string, bioSig string, data interface{}, txType transactions.Type) error {
	if err := s.notif.NotifyTransactionStatus(txHash, Pending); err != nil {
		return err
	}

	lastVPool, vPool, sPool, err := pool.GetPools(addr, s.poolF)
	if err != nil {
		return err
	}

	if err := s.requestLock(txHash, lastVPool); err != nil {
		return err
	}

	//Mine the transaction
	pow := NewPOW(s.list, s.sig, s.robotKey, s.robotPvKey)
	masterValid, err := pow.Execute(txHash, bioSig, lastVPool)
	if err != nil {
		return err
	}

	var valids []datamining.Validation
	valids, err = s.txHandlers[txType].RequestValidations(s.poolD, vPool, data, txType)
	if err != nil {
		return err
	}

	//Check if the validations passed
	for _, v := range valids {
		if v.Status() == datamining.ValidationKO {
			if err := s.notif.NotifyTransactionStatus(txHash, Invalid); err != nil {
				return err
			}
			return nil
		}
	}
	if err := s.notif.NotifyTransactionStatus(txHash, Approved); err != nil {
		return err
	}

	if err := s.notif.NotifyTransactionStatus(txHash, Approved); err != nil {
		return err
	}

	//Wraps validations
	endorsement := datamining.NewEndorsement(time.Now(), txHash, masterValid, valids)

	//Execute a storage request to write the data and the validations
	if err := s.txHandlers[txType].RequestStorage(s.poolD, sPool, data, endorsement, txType); err != nil {
		return err
	}

	if err := s.requestUnlock(txHash, lastVPool); err != nil {
		return err
	}
	return s.notif.NotifyTransactionStatus(txHash, Replicated)
}

func (s service) Validate(data interface{}, txType transactions.Type) (valid datamining.Validation, err error) {
	for _, c := range s.checks[txType] {
		err = c.CheckData(data)
		if err != nil {
			if c.IsCatchedError(err) {
				return s.buildValidation(datamining.ValidationKO)
			}
			return
		}
	}
	return s.buildValidation(datamining.ValidationOK)
}

func (s service) LockTransaction(txLock lock.TransactionLock) error {
	if s.locker.ContainsLock(txLock) {
		return ErrLockExisting
	}

	return s.locker.Lock(txLock)
}

func (s service) UnlockTransaction(txLock lock.TransactionLock) error {
	return s.locker.Unlock(txLock)
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

func (s service) requestLock(txHash string, lastVPool pool.PeerCluster) error {
	lock := lock.TransactionLock{TxHash: txHash, MasterRobotKey: s.robotKey}
	sigLock, err := s.sig.SignLock(lock, s.robotPvKey)
	if err != nil {
		return err
	}

	if err := s.poolD.RequestLock(lastVPool, lock, sigLock); err != nil {
		return err
	}
	if err := s.notif.NotifyTransactionStatus(txHash, Locked); err != nil {
		return err
	}

	return err
}

func (s service) requestUnlock(txHash string, lastVPool pool.PeerCluster) error {
	lock := lock.TransactionLock{TxHash: txHash, MasterRobotKey: s.robotKey}
	sigLock, err := s.sig.SignLock(lock, s.robotPvKey)
	if err != nil {
		return err
	}
	if err := s.poolD.RequestUnlock(lastVPool, lock, sigLock); err != nil {
		return err
	}
	if err := s.notif.NotifyTransactionStatus(txHash, Unlocked); err != nil {
		return err
	}

	return err
}
