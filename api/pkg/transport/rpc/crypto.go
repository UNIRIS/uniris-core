package rpc

import (
	adding "github.com/uniris/uniris-core/api/pkg/adding"
	"github.com/uniris/uniris-core/api/pkg/listing"
	api "github.com/uniris/uniris-core/datamining/api/protobuf-spec"
)

//SignatureHandler define methods to handle signature
type SignatureHandler interface {
	signatureChecker
	signatureBuilder
}

type signatureChecker interface {

	//VerifyAccountSearchResultSignature checks the signature of the account search result
	VerifyAccountSearchResultSignature(pubKey string, res *api.AccountSearchResult) error

	//VerifyCreationResultSignature checks the signature of the creation result
	VerifyCreationResultSignature(pubKey string, res *api.CreationResult) error
}

type signatureBuilder interface {

	//SignAccountResult signs account research result
	SignAccountResult(res *listing.AccountResult, pvKey string) error

	//SignAccountCreationResult signs the account creation result
	SignAccountCreationResult(res *adding.AccountCreationResult, pvKey string) error
}
