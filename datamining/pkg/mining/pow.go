package mining

import (
	"time"

	datamining "github.com/uniris/uniris-core/datamining/pkg"
	"github.com/uniris/uniris-core/datamining/pkg/listing"
	"github.com/uniris/uniris-core/datamining/pkg/mining/pool"
)

//PowSigner defines methods to handle signatures
type PowSigner interface {
	SignMasterValidation(v Validation, pvKey string) (string, error)
	CheckTransactionSignature(pubKey string, tx string, sig string) error
}

//POW defines methods for the POW
type POW interface {
	Execute(txHash, sig string, lastValidationPool pool.PeerCluster) (*datamining.MasterValidation, error)
}

//Validation represents a validation before its signature
type Validation struct {
	Status    datamining.ValidationStatus `json:"status"`
	Timestamp time.Time                   `json:"timestamp"`
	PublicKey string                      `json:"pubk"`
}

type pow struct {
	list        listing.Service
	sig         PowSigner
	robotPubKey string
	robotPvKey  string
}

//NewPOW creates a new Proof Of Work handler
func NewPOW(list listing.Service, sig PowSigner, robotPubKey, robotPvKey string) POW {
	return pow{list, sig, robotPubKey, robotPvKey}
}

//Execute the Proof Of Work
func (p pow) Execute(txHash string, sig string, lastValidationPool pool.PeerCluster) (*datamining.MasterValidation, error) {
	keys, err := p.list.ListBiodPubKeys()
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

	lastTxMiners := make([]string, 0)
	for _, p := range lastValidationPool.Peers {
		lastTxMiners = append(lastTxMiners, p.PublicKey)
	}

	return datamining.NewMasterValidation(lastTxMiners, p.robotPubKey, valid), nil
}
