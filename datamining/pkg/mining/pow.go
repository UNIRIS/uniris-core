package mining

import (
	"time"

	datamining "github.com/uniris/uniris-core/datamining/pkg"
	"github.com/uniris/uniris-core/datamining/pkg/biod/listing"
)

//PowSigner defines methods to handle signatures
type PowSigner interface {

	//SignValidation create signature for a validation data
	SignValidation(v Validation, pvKey string) (string, error)

	//CheckTransactionDataSignature checks the transaction data signature
	CheckTransactionDataSignature(txType TransactionType, pubKey string, data interface{}, sig string) error
}

type pow struct {
	txType      TransactionType
	txData      interface{}
	lastVPool   datamining.Pool
	txBiodSig   string
	lister      listing.Service
	signer      PowSigner
	robotPubKey string
	robotPvKey  string
}

func (p pow) execute() (MasterValidation, error) {
	keys, err := p.lister.ListBiodPubKeys()
	if err != nil {
		return nil, err
	}

	//Find the public key which matches the transaction signature
	status := ValidationKO
	for _, k := range keys {
		err := p.signer.CheckTransactionDataSignature(p.txType, k, p.txData, p.txBiodSig)
		if err == nil {
			status = ValidationOK
			break
		}
	}

	v := validation{
		pubk:      p.robotPubKey,
		status:    status,
		timestamp: time.Now(),
	}
	signature, err := p.signer.SignValidation(v, p.robotPvKey)
	if err != nil {
		return nil, err
	}

	valid := NewValidation(
		v.status,
		v.timestamp,
		v.pubk,
		signature)

	lastTxMiners := make([]string, 0)
	for _, peer := range p.lastVPool.Peers() {
		lastTxMiners = append(lastTxMiners, peer.PublicKey)
	}

	return NewMasterValidation(lastTxMiners, p.robotPubKey, valid), nil
}
