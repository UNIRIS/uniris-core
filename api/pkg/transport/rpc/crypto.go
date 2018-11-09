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
	CheckAccountSearchResultSignature(pubKey string, res *api.AccountSearchResult) error
	CheckCreationResultSignature(pubKey string, res *api.CreationResult) error
}

type signatureBuilder interface {
	SignAccountResult(res *listing.AccountResult, pvKey string) error
	SignAccountCreationResult(res *adding.AccountCreationResult, pvKey string) error
}
