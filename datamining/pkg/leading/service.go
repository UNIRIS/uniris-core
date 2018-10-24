package leading

import (
	"time"

	datamining "github.com/uniris/uniris-core/datamining/pkg"
)

//Hasher defines methods for generate hash
type Hasher interface {
	HashMasterValidation(v *datamining.MasterValidation) (string, error)
}

//Service defines leading methods
type Service interface {
	ValidateWallet(w *datamining.WalletData, txHash string, lastTxMinerList []string) (*datamining.Endorsement, error)
	ValidateBio(b *datamining.BioData, txHash string, lastTxMinerList []string) (*datamining.Endorsement, error)
	ComputeWallet(w *datamining.WalletData, txHash string) error
	ComputeBio(b *datamining.BioData, txHash string) error
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
func NewService(poolF PoolFinder, poolD PoolDispatcher, notif Notifier, sig Signer, techRepo TechRepository, robotKey, robotPvKey string) Service {
	return service{poolF, poolD, notif, techRepo, sig, robotKey, robotPvKey}
}

func (s service) ValidateWallet(w *datamining.WalletData, txHash string, lastTxMinerList []string) (*datamining.Endorsement, error) {
	pool, err := s.poolF.FindValidationPool()
	if err != nil {
		return nil, err
	}

	if err := s.poolD.RequestLock(pool, txHash); err != nil {
		return nil, err
	}
	if err := s.notif.NotifyTransactionStatus(txHash, Locked); err != nil {
		return nil, err
	}

	pow := NewPOW(s.techRepo, s.sig, s.robotKey, s.robotPvKey)
	masterValid, err := pow.Execute(txHash, w.Sigs.BiodSig, lastTxMinerList)
	if err != nil {
		return nil, err
	}

	valids, err := s.poolD.RequestWalletValidation(pool, w)
	if err != nil {
		return nil, err
	}

	if err := s.notif.NotifyTransactionStatus(txHash, Approved); err != nil {
		return nil, err
	}

	if err := s.poolD.RequestUnlock(pool, txHash); err != nil {
		return nil, err
	}

	return datamining.NewEndorsement(time.Now(),
		txHash,
		masterValid,
		valids), nil
}

func (s service) ValidateBio(bd *datamining.BioData, txHash string, lastTxMinerList []string) (*datamining.Endorsement, error) {

	pool, err := s.poolF.FindValidationPool()
	if err != nil {
		return nil, err
	}

	if err := s.poolD.RequestLock(pool, txHash); err != nil {
		return nil, err
	}
	if err := s.notif.NotifyTransactionStatus(txHash, Locked); err != nil {
		return nil, err
	}

	pow := NewPOW(s.techRepo, s.sig, s.robotKey, s.robotPvKey)
	masterValid, err := pow.Execute(txHash, bd.Sigs.BiodSig, lastTxMinerList)
	if err != nil {
		return nil, err
	}

	valids, err := s.poolD.RequestBioValidation(pool, bd)
	if err != nil {
		return nil, err
	}

	if err := s.notif.NotifyTransactionStatus(txHash, Approved); err != nil {
		return nil, err
	}

	if err := s.poolD.RequestUnlock(pool, txHash); err != nil {
		return nil, err
	}

	return datamining.NewEndorsement(time.Now(),
		txHash,
		masterValid,
		valids), nil
}

func (s service) ComputeWallet(data *datamining.WalletData, txHash string) error {
	if err := s.notif.NotifyTransactionStatus(txHash, Pending); err != nil {
		return err
	}
	if err := s.notif.NotifyTransactionStatus(txHash, Pending); err != nil {
		return err
	}
	pool, err := s.poolF.FindStoragePool()
	if err != nil {
		return err
	}

	oldTxHash, lastMasterValidation, err := s.poolD.RequestLastTx(pool, txHash)
	if err != nil {
		return err
	}

	endorsement, err := s.ValidateWallet(data, txHash, lastMasterValidation.ValidatorKeysOfLastTransaction())
	if err != nil {
		return err
	}

	w := datamining.NewWallet(data, endorsement, oldTxHash)

	if err := s.poolD.RequestWalletStorage(pool, w); err != nil {
		return err
	}

	return nil
}

func (s service) ComputeBio(data *datamining.BioData, txHash string) error {

	pool, err := s.poolF.FindStoragePool()
	if err != nil {
		return err
	}

	_, lastMasterValidation, err := s.poolD.RequestLastTx(pool, txHash)
	if err != nil {
		return err
	}

	endorsement, err := s.ValidateBio(data, txHash, lastMasterValidation.ValidatorKeysOfLastTransaction())
	if err != nil {
		return err
	}

	bw := datamining.NewBioWallet(data, endorsement)

	if err := s.poolD.RequestBioStorage(pool, bw); err != nil {
		return err
	}

	return nil
}
