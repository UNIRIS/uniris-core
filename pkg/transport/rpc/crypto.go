package rpc

import (
	api "github.com/uniris/uniris-core/api/protobuf-spec"
)

type Decrypter interface {
	DecryptString(str string, pvKey string) (string, error)
}

type Hasher interface {
	HashString(str string) string
}

type SignatureHandler interface {
	SignKeychainRequest(res *api.KeychainRequest, pvKey string) error
	SignKeychainResponse(res *api.KeychainResponse, pvKey string) error
	SignIDRequest(res *api.IDRequest, pvKey string) error
	SignIDResponse(res *api.IDResponse, pvKey string) error
	SignTransactionResult(res *api.TransactionResult, pvKey string) error
	SignAccountResponse(res *api.GetAccountResponse, pvKey string) error
	SignTransactionStatusResponse(res *api.TransactionStatusResponse, pvKey string) error
	SignLockResponse(res *api.LockResponse, pvKey string) error
	SignPreValidationResponse(res *api.PreValidationResponse, pvKey string) error
	SignConfirmValidationResponse(res *api.ConfirmValidationResponse, pvKey string) error
	SignStoreResponse(res *api.StoreResponse, pvKey string) error

	VerifyKeychainRequestSignature(req *api.KeychainRequest, pubKey string) error
	VerifyIDRequestSignature(req *api.IDRequest, pubKey string) error
	VerifyTransactionStatusResponseSignature(res *api.TransactionStatusResponse, pubKey string) error
	VerifyAccountRequestSignature(req *api.GetAccountRequest, pubKey string) error
	VerifyTransactionStatusRequestSignature(req *api.TransactionStatusRequest, pubKey string) error
	VerifyLockRequestSignature(req *api.LockRequest, pubKey string) error
	VerifyPreValidateRequestSignature(req *api.PreValidationRequest, pubKey string) error
	VerifyConfirmValidationRequestSignature(req *api.ConfirmValidationRequest, pubKey string) error
	VerifyStoreRequest(req *api.StoreRequest, pubKey string) error
}
