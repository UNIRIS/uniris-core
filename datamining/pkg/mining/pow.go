package mining

import (
	"time"

	datamining "github.com/uniris/uniris-core/datamining/pkg"
	"github.com/uniris/uniris-core/datamining/pkg/biod/listing"
)

type pow struct {
	txType      TransactionType
	txData      interface{}
	lastVPool   datamining.Pool
	txBiodSig   string
	lister      listing.Service
	signer      signer
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
	var matchedKey string
	for _, k := range keys {
		err := p.signer.VerifyTransactionDataSignature(p.txType, k, p.txData, p.txBiodSig)
		if err == nil {
			matchedKey = k
			status = ValidationOK
			break
		}
	}

	v := validation{
		pubk:      p.robotPubKey,
		status:    status,
		timestamp: time.Now(),
	}
	sValid, err := p.signer.SignValidation(v, p.robotPvKey)
	if err != nil {
		return nil, err
	}

	lastTxMiners := make([]string, 0)
	for _, peer := range p.lastVPool.Peers() {
		lastTxMiners = append(lastTxMiners, peer.PublicKey)
	}

	return NewMasterValidation(lastTxMiners, matchedKey, sValid), nil
}
