package transaction

import (
	"encoding/json"
	"errors"

	"github.com/uniris/uniris-core/pkg/shared"
)

//Proposal describe a proposal for a Transaction
type Proposal struct {
	sharedEmitterKP shared.KeyPair
}

//NewProposal create a new proposal for a Transaction
func NewProposal(shdEmitterKP shared.KeyPair) (Proposal, error) {
	if (shdEmitterKP == shared.KeyPair{}) {
		return Proposal{}, errors.New("Transaction proposal: missing shared keys")
	}
	return Proposal{
		sharedEmitterKP: shdEmitterKP,
	}, nil
}

//SharedEmitterKeyPair returns the keypair proposed for the shared emitter keys
func (p Proposal) SharedEmitterKeyPair() shared.KeyPair {
	return p.sharedEmitterKP
}

func (p Proposal) MarshalJSON() ([]byte, error) {
	return json.Marshal(map[string]interface{}{
		"shared_emitter_keys": p.SharedEmitterKeyPair(),
	})
}
