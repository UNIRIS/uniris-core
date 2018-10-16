package externalrpc

import (
	datamining "github.com/uniris/uniris-core/datamining/pkg"
	"github.com/uniris/uniris-core/datamining/pkg/validating"
)

type validatorRequester struct{}

//NewValidatorRequest creates a validator request
func NewValidatorRequest() validating.ValidationRequester {
	return validatorRequester{}
}

func (v validatorRequester) RequestWalletValidation(validating.Peer, datamining.WalletData) (datamining.Validation, error) {
	var validation datamining.Validation

	//TODO

	return validation, nil
}
func (v validatorRequester) RequestBioValidation(validating.Peer, datamining.BioData) (datamining.Validation, error) {
	var validation datamining.Validation

	//TODO

	return validation, nil
}
