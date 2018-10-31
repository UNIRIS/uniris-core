package master

import (
	"time"

	datamining "github.com/uniris/uniris-core/datamining/pkg"
	"github.com/uniris/uniris-core/datamining/pkg/listing"
	"github.com/uniris/uniris-core/datamining/pkg/mining/master/pool"
)

//PowSigner defines methods to handle signatures
type PowSigner interface {
	SignMasterValidation(v Validation, pvKey string) (string, error)
	CheckTransactionSignature(pubKey string, tx string, sig string) error
}

//POW defines methods for the POW
type POW interface {
	Execute(txHash, biodSig string, lastValidationPool pool.PeerGroup) (*datamining.MasterValidation, error)
}

//Validation represents a validation before its signature
type Validation struct {
	Status    datamining.ValidationStatus `json:"status"`
	Timestamp time.Time                   `json:"timestamp"`
	PublicKey string                      `json:"pubk"`
}

type pow struct {
	lister      listing.Service
	signer      PowSigner
	robotPubKey string
	robotPvKey  string
}

//NewPOW creates a new Proof Of Work handler
func NewPOW(lister listing.Service, signer PowSigner, robotPubKey, robotPvKey string) POW {
	return pow{lister, signer, robotPubKey, robotPvKey}
}

//Execute the Proof Of Work
func (p pow) Execute(txHash string, biodSig string, lastValidationPool pool.PeerGroup) (*datamining.MasterValidation, error) {
	keys, err := p.lister.ListBiodPubKeys()
	if err != nil {
		return nil, err
	}

	//Find the public key which matches the transaction signature
	status := datamining.ValidationKO
	for _, k := range keys {
		err := p.signer.CheckTransactionSignature(k, txHash, biodSig)
		if err == nil {
			status = datamining.ValidationOK
			break
		}
	}

	v := Validation{
		PublicKey: p.robotPubKey,
		Status:    status,
		Timestamp: time.Now(),
	}
	signature, err := p.signer.SignMasterValidation(v, p.robotPvKey)
	if err != nil {
		return nil, err
	}

	valid := datamining.NewValidation(
		v.Status,
		v.Timestamp,
		v.PublicKey,
		signature)

	lastTxMiners := make([]string, 0)
	for _, p := range lastValidationPool.Peers {
		lastTxMiners = append(lastTxMiners, p.PublicKey)
	}

	return datamining.NewMasterValidation(lastTxMiners, p.robotPubKey, valid), nil
}
