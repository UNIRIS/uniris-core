package transaction

import (
	"encoding/json"
	"errors"

	"github.com/uniris/uniris-core/pkg/shared"
)

//Proposal describe a proposal for a Transaction
type Proposal struct {
	sharedEmitterKP shared.EmitterKeyPair
}

//NewProposal create a new proposal for a Transaction
func NewProposal(shdEmitterKP shared.EmitterKeyPair) (Proposal, error) {
	if (shdEmitterKP == shared.EmitterKeyPair{}) {
		return Proposal{}, errors.New("Transaction proposal: missing shared emitter keys")
	}
	return Proposal{
		sharedEmitterKP: shdEmitterKP,
	}, nil
}

//SharedEmitterKeyPair returns the keypair proposed for the shared emitter keys
func (p Proposal) SharedEmitterKeyPair() shared.EmitterKeyPair {
	return p.sharedEmitterKP
}

//MarshalJSON serialize the transaction proposal into a JSON
func (p Proposal) MarshalJSON() ([]byte, error) {
	return json.Marshal(map[string]interface{}{
		"shared_emitter_keys": p.SharedEmitterKeyPair(),
	})
}
