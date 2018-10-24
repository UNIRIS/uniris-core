package leading

import (
	"time"

	datamining "github.com/uniris/uniris-core/datamining/pkg"
)

//Signer defines methods to handle signatures
type Signer interface {
	SignMasterValidation(v Validation, pvKey string) (string, error)
	CheckTransactionSignature(pubKey string, tx string, sig string) error
}

//TechRepository defines methods to query the bank repository
type TechRepository interface {
	ListBiodPubKeys() ([]string, error)
}

//POW defines methods for the POW
type POW interface {
	Execute(txHash, sig string, lastTxMinerList []string) (*datamining.MasterValidation, error)
}

//Validation represents a validation before its signature
type Validation struct {
	Status    datamining.ValidationStatus `json:"status"`
	Timestamp time.Time                   `json:"timestamp"`
	PublicKey string                      `json:"pubk"`
}

type pow struct {
	repo        TechRepository
	sig         Signer
	robotPubKey string
	robotPvKey  string
}

//NewPOW creates a new Proof Of Work handler
func NewPOW(repo TechRepository, sig Signer, robotPubKey, robotPvKey string) POW {
	return pow{repo, sig, robotPubKey, robotPvKey}
}

//Execute the Proof Of Work
func (p pow) Execute(txHash string, sig string, lastTxMinerList []string) (*datamining.MasterValidation, error) {
	keys, err := p.repo.ListBiodPubKeys()
	if err != nil {
		return nil, err
	}

	//Find the public key which matches the transaction signature
	for _, k := range keys {
		err := p.sig.CheckTransactionSignature(k, txHash, sig)
		if err == nil {
			break
		}
	}

	v := Validation{
		PublicKey: p.robotPubKey,
		Status:    datamining.ValidationOK,
		Timestamp: time.Now(),
	}
	signature, err := p.sig.SignMasterValidation(v, p.robotPvKey)
	if err != nil {
		return nil, err
	}

	valid := datamining.NewValidation(
		v.Status,
		v.Timestamp,
		v.PublicKey,
		signature)

	return datamining.NewMasterValidation(lastTxMinerList, p.robotPubKey, valid), nil
}
