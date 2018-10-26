package leading

import (
	"errors"
	"time"

	"github.com/uniris/uniris-core/datamining/pkg/validating"

	datamining "github.com/uniris/uniris-core/datamining/pkg"
)

//Signer defines methods to handle signatures
type Signer interface {
	SignLock(lock validating.TransactionLock, pvKey string) (string, error)
	PowSigner
}

//Service defines leading methods
type Service interface {
	LeadWalletTransaction(w *datamining.WalletData, txHash string) error
	LeadBioTransaction(b *datamining.BioData, txHash string) error
}

type service struct {
	poolF      PoolFinder
	poolD      PoolDispatcher
	notif      Notifier
	techRepo   TechRepository
	sig        Signer
	h          PreviousDataHasher
	robotKey   string
	robotPvKey string
}

//NewService creates a new leading service
func NewService(pF PoolFinder, pD PoolDispatcher, notif Notifier, sig Signer, h PreviousDataHasher, tRepo TechRepository, rPb, rPv string) Service {
	return &service{
		poolF:      pF,
		poolD:      pD,
		notif:      notif,
		techRepo:   tRepo,
		sig:        sig,
		h:          h,
		robotKey:   rPb,
		robotPvKey: rPv,
	}
}

func (s service) LeadWalletTransaction(data *datamining.WalletData, txHash string) error {
	if err := s.notif.NotifyTransactionStatus(txHash, Pending); err != nil {
		return err
	}

	lastVPool, vPool, sPool, err := s.getPools(data.CipherAddrRobot)
	if err != nil {
		return err
	}

	//Retrieve the last wallet chain and check them
	lastWalletEntries, err := s.poolD.RequestLastWallet(sPool, data.CipherAddrRobot)
	if err != nil {
		return err
	}
	if err := s.compareWallets(lastWalletEntries); err != nil {
		return err
	}

	if err := s.requestLock(txHash, lastVPool); err != nil {
		return err
	}

	//Mine the transaction
	pow := NewPOW(s.techRepo, s.sig, s.robotKey, s.robotPvKey)
	masterValid, err := pow.Execute(txHash, data.Sigs.BiodSig, lastVPool)
	if err != nil {
		return err
	}
	valids, err := s.poolD.RequestWalletValidation(vPool, data)
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

	//Create and store validations
	endorsement := datamining.NewEndorsement(time.Now(), txHash, masterValid, valids)
	w := datamining.NewWallet(data, endorsement, lastWalletEntries[0].OldTransactionHash())
	if err := s.poolD.RequestWalletStorage(sPool, w); err != nil {
		return err
	}

	if err := s.requestUnlock(txHash, lastVPool); err != nil {
		return err
	}
	return s.notif.NotifyTransactionStatus(txHash, Replicated)
}

func (s *service) LeadBioTransaction(data *datamining.BioData, txHash string) error {
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

	//Mine the transaction
	pow := NewPOW(s.techRepo, s.sig, s.robotKey, s.robotPvKey)
	masterValid, err := pow.Execute(txHash, data.Sigs.BiodSig, lastVPool)
	if err != nil {
		return err
	}
	valids, err := s.poolD.RequestBioValidation(vPool, data)
	if err != nil {
		return err
	}

	//Check if the validations passed
	ok, err := s.checkValidations(valids, txHash)
	if err != nil {
		return err
	}
	if !ok {
		if err := s.notif.NotifyTransactionStatus(txHash, Invalid); err != nil {
			return err
		}
	}

	if err := s.notif.NotifyTransactionStatus(txHash, Approved); err != nil {
		return err
	}

	//Create and store validations
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

func (s service) checkValidations(valids []datamining.Validation, txHash string) (bool, error) {
	for _, v := range valids {
		if v.Status() == datamining.ValidationKO {
			return false, nil
		}
	}
	return true, nil
}

func (s service) compareWallets(wallets []*datamining.Wallet) error {
	var invalidLen = 0
	for i := 1; i < len(wallets); i++ {
		if wallets[i] != wallets[0] {
			invalidLen++
			continue
		}
		check := NewPreviousDataChecker(s.h)
		if err := check.CheckPreviousWallet(wallets[i], wallets[i].OldTransactionHash()); err != nil {
			invalidLen++
		}
	}

	if invalidLen > 0 {
		return errors.New("Invalid wallets")
	}

	return nil
}
