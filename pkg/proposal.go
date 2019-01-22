package uniris

import "encoding/json"

//TransactionProposal describe a proposal for a Transaction
type TransactionProposal interface {

	//SharedEmitterKeyPair returns the keypair proposed for the shared emitter keys
	SharedEmitterKeyPair() SharedKeys
}

type txProp struct {
	sharedEmitterKP SharedKeys
}

//NewTransactionProposal create a new proposal for a Transaction
func NewTransactionProposal(shdEmitterKP SharedKeys) TransactionProposal {
	return txProp{
		sharedEmitterKP: shdEmitterKP,
	}
}

func (p txProp) SharedEmitterKeyPair() SharedKeys {
	return p.sharedEmitterKP
}

func (p txProp) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		SharedEmitterKP SharedKeys `json:"shared_emitter_keys"`
	}{
		SharedEmitterKP: p.SharedEmitterKeyPair(),
	})
}
