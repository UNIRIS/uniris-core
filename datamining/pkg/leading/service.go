package leading

import (
	"time"

	"github.com/uniris/uniris-core/datamining/pkg/validating"

	datamining "github.com/uniris/uniris-core/datamining/pkg"
)

//Signer defines methods to handle signatures
type Signer interface {
	SignLock(lock validating.TransactionLock, pvKey string) (string, error)
	PowSigner
}

//Hasher defines methods for generate hash
type Hasher interface {
	HashMasterValidation(v *datamining.MasterValidation) (string, error)
}

//Service defines leading methods
type Service interface {
	NewWallet(w *datamining.WalletData, txHash string) error
	NewBio(b *datamining.BioData, txHash string) error
}

type service struct {
	poolF      PoolFinder
	poolD      PoolDispatcher
	notif      Notifier
	techRepo   TechRepository
	sig        Signer
	robotKey   string
	robotPvKey string
}

//NewService creates a new leading service
func NewService(pF PoolFinder, pD PoolDispatcher, notif Notifier, sig Signer, tRepo TechRepository, rPb, rPv string) Service {
	return &service{
		poolF:      pF,
		poolD:      pD,
		notif:      notif,
		techRepo:   tRepo,
		sig:        sig,
		robotKey:   rPb,
		robotPvKey: rPv,
	}
}

func (s service) NewWallet(data *datamining.WalletData, txHash string) error {
	if err := s.notif.NotifyTransactionStatus(txHash, Pending); err != nil {
		return err
	}

	lastVPool, vPool, sPool, err := s.getPools(data.CipherAddrRobot)
	if err != nil {
		return err
	}

	oldTxHash, err := s.poolD.RequestLastTx(sPool, txHash)
	if err != nil {
		return err
	}

	if err := s.requestLock(txHash, lastVPool); err != nil {
		return err
	}

	pow := NewPOW(s.techRepo, s.sig, s.robotKey, s.robotPvKey)
	masterValid, err := pow.Execute(txHash, data.Sigs.BiodSig, lastVPool)
	if err != nil {
		return err
	}

	valids, err := s.poolD.RequestWalletValidation(vPool, data, txHash)
	if err != nil {
		return err
	}

	if err := s.notif.NotifyTransactionStatus(txHash, Approved); err != nil {
		return err
	}

	endorsement := datamining.NewEndorsement(time.Now(), txHash, masterValid, valids)
	w := datamining.NewWallet(data, endorsement, oldTxHash)

	if err := s.poolD.RequestWalletStorage(sPool, w); err != nil {
		return err
	}

	if err := s.requestUnlock(txHash, lastVPool); err != nil {
		return err
	}
	return s.notif.NotifyTransactionStatus(txHash, Replicated)
}

func (s *service) NewBio(data *datamining.BioData, txHash string) error {
	if err := s.notif.NotifyTransactionStatus(txHash, Pending); err != nil {
		return err
	}

	lastVPool, vPool, sPool, err := s.getPools(data.CipherAddrRobot)
	if err != nil {
		return err
	}

	if err := s.requestLock(txHash, lastVPool); err != nil {
		return err
	}

	pow := NewPOW(s.techRepo, s.sig, s.robotKey, s.robotPvKey)
	masterValid, err := pow.Execute(txHash, data.Sigs.BiodSig, lastVPool)
	if err != nil {
		return err
	}

	valids, err := s.poolD.RequestBioValidation(vPool, data, txHash)
	if err != nil {
		return err
	}

	if err := s.notif.NotifyTransactionStatus(txHash, Approved); err != nil {
		return err
	}

	endorsement := datamining.NewEndorsement(time.Now(), txHash, masterValid, valids)
	bw := datamining.NewBioWallet(data, endorsement)

	if err := s.poolD.RequestBioStorage(sPool, bw); err != nil {
		return err
	}

	if err := s.requestUnlock(txHash, lastVPool); err != nil {
		return err
	}
	return s.notif.NotifyTransactionStatus(txHash, Replicated)
}

func (s service) getPools(addr string) (lastVPool Pool, vPool Pool, sPool Pool, err error) {
	sPool, err = s.poolF.FindStoragePool()
	if err != nil {
		return
	}

	lastVPool, err = s.poolF.FindLastValidationPool(addr)
	if err != nil {
		return
	}

	vPool, err = s.poolF.FindValidationPool()
	if err != nil {
		return
	}

	return
}

func (s service) requestLock(txHash string, lastVPool Pool) error {
	lock := validating.TransactionLock{TxHash: txHash, MasterRobotKey: s.robotKey}
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

func (s service) requestUnlock(txHash string, lastVPool Pool) error {
	lock := validating.TransactionLock{TxHash: txHash, MasterRobotKey: s.robotKey}
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
