package master

import (
	"time"

	datamining "github.com/uniris/uniris-core/datamining/pkg"
	"github.com/uniris/uniris-core/datamining/pkg/listing"
	"github.com/uniris/uniris-core/datamining/pkg/mining/master/checks"
	"github.com/uniris/uniris-core/datamining/pkg/mining/master/pool"
	"github.com/uniris/uniris-core/datamining/pkg/mining/master/transactions"
)

//Service defines methods for the master mining
type Service interface {
	LeadMining(txHash string, addr string, bioSig string, data interface{}, txType datamining.TransactionType) error
}

type service struct {
	poolF      pool.Finder
	poolD      pool.Requester
	notif      Notifier
	sig        Signer
	listSrv    listing.Service
	robotKey   string
	robotPvKey string

	txHandlers map[datamining.TransactionType]transactions.Handler
	txCheckers map[datamining.TransactionType][]checks.Handler
}

//NewService creates master mining service
func NewService(poolF pool.Finder, poolD pool.Requester, n Notifier, sig Signer, txHash Hasher, listSrv listing.Service, robotKey, robotPvKey string) Service {

	txHandlers := map[datamining.TransactionType]transactions.Handler{
		datamining.CreateKeychainTransaction: transactions.NewCreateKeychainHandler(),
		datamining.CreateBioTransaction:      transactions.NewCreateBiometricHandler(),
	}

	txCheckers := map[datamining.TransactionType][]checks.Handler{
		datamining.CreateBioTransaction: []checks.Handler{
			checks.NewIntegrityChecker(txHash),
		},
		datamining.CreateKeychainTransaction: []checks.Handler{
			checks.NewIntegrityChecker(txHash),
		},
	}

	return service{poolF, poolD, n, sig, listSrv, robotKey, robotPvKey, txHandlers, txCheckers}
}

func (s service) LeadMining(txHash string, addr string, bioSig string, data interface{}, txType datamining.TransactionType) error {
	if err := s.notif.NotifyTransactionStatus(txHash, Pending); err != nil {
		return err
	}

	lastVPool, vPool, sPool, err := pool.Lookup(addr, s.poolF)
	if err != nil {
		return err
	}

	if err := s.requestLock(txHash, lastVPool); err != nil {
		return err
	}

	masterValid, valids, err := s.mine(txHash, data, bioSig, lastVPool, vPool, txType)
	if err != nil {
		if err == checks.ErrInvalidTransaction {
			if err := s.notif.NotifyTransactionStatus(txHash, Invalid); err != nil {
				return err
			}
			return nil
		}
		return err
	}

	if err := s.notif.NotifyTransactionStatus(txHash, Approved); err != nil {
		return err
	}

	return s.requestStorage(txHash, data, s.txHandlers[txType], masterValid, valids, lastVPool, sPool)
}

func (s service) mine(txHash string, data interface{}, bioSig string, lastVPool, vPool pool.PeerGroup, txType datamining.TransactionType) (*datamining.MasterValidation, []datamining.Validation, error) {
	//Check data before to perform POW
	for _, c := range s.txCheckers[txType] {
		if err := c.CheckData(data, txHash); err != nil {
			return nil, nil, err
		}
	}

	masterValid, err := NewPOW(s.listSrv, s.sig, s.robotKey, s.robotPvKey).Execute(txHash, bioSig, lastVPool)
	if err != nil {
		return nil, nil, err
	}

	//Ask a pool of peers to validate the transaction
	valids, err := s.requestValidations(vPool, txHash, data, s.txHandlers[txType])
	if err != nil {
		return nil, nil, err
	}

	return masterValid, valids, nil
}

func (s service) requestValidations(vPool pool.PeerGroup, txHash string, data interface{}, txHandler transactions.Handler) ([]datamining.Validation, error) {
	valids, err := txHandler.RequestValidations(s.poolD, vPool, data)
	if err != nil {
		return nil, err
	}

	//Check if the validations passed
	for _, v := range valids {
		if v.Status() == datamining.ValidationKO {
			return nil, checks.ErrInvalidTransaction
		}
	}
	return valids, nil
}

func (s service) requestStorage(txHash string, data interface{}, txHandler transactions.Handler, masterValid *datamining.MasterValidation, valids []datamining.Validation, lastVPool, sPool pool.PeerGroup) error {

	endorsement := datamining.NewEndorsement(time.Now(), txHash, masterValid, valids)

	//Execute a storage request to write the data and the validations
	if err := txHandler.RequestStorage(s.poolD, sPool, data, endorsement); err != nil {
		return err
	}

	if err := s.requestUnlock(txHash, lastVPool); err != nil {
		return err
	}
	return s.notif.NotifyTransactionStatus(txHash, Replicated)
}

func (s service) requestLock(txHash string, lastVPool pool.PeerGroup) error {
	lock := pool.TransactionLock{TxHash: txHash, MasterRobotKey: s.robotKey}
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

func (s service) requestUnlock(txHash string, lastVPool pool.PeerGroup) error {
	lock := pool.TransactionLock{TxHash: txHash, MasterRobotKey: s.robotKey}
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
