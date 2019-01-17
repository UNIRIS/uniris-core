package rpc

import (
	"errors"

	"github.com/uniris/uniris-core/datamining/pkg/contract"

	api "github.com/uniris/uniris-core/datamining/api/protobuf-spec"
	"github.com/uniris/uniris-core/datamining/pkg/account"
	"github.com/uniris/uniris-core/datamining/pkg/lock"
)

//Crypto defines the required crypto handlers
type Crypto struct {
	decrypter Decrypter
	signer    Signer
	hasher    Hasher
}

//NewCrypto creates a new crypto handler for GRPC
func NewCrypto(d Decrypter, s Signer, h Hasher) Crypto {
	return Crypto{d, s, h}
}

//ErrInvalidEncryption is returned when the encrypted data cannot be decrypted
var ErrInvalidEncryption = errors.New("Invalid encryption")

//ErrInvalidSignature is returned when the signature is invalid
var ErrInvalidSignature = errors.New("Invalid signature")

//Decrypter define methods to decrypt data for RPC methods
type Decrypter interface {

	//DecryptHash decrypt a hash using the shared robot private key
	DecryptHash(hash string, pvKey string) (string, error)

	//DecryptKeychain decrypt the account's keychain using the shared robot private key
	DecryptKeychain(data string, pvKey string) (account.Keychain, error)

	//DecryptID decrypt the account's ID using the shared robot private key
	DecryptID(data string, pvKey string) (account.ID, error)

	DecryptContract(data string, pvKey string) (contract.Contract, error)
	DecryptContractMessage(data string, pvKey string) (contract.Message, error)
}

//Hasher define methods to hash incoming data
type Hasher interface {
	account.KeychainHasher
	account.IDHasher
	lock.Hasher
	contract.Hasher
}

//Signer define methods to handle signatures
type Signer interface {
	signatureBuilder
	signatureVerifier
}

type signatureVerifier interface {

	//VerifyKeychainValidationRequestSignature checks the keychain validation request's signature
	//Using the share robot public key
	VerifyKeychainValidationRequestSignature(pubKey string, req *api.KeychainValidationRequest) error

	//VerifyIDValidationRequestSignature checks the ID validation request's signature
	//Using the share robot public key
	VerifyIDValidationRequestSignature(pubKey string, req *api.IDValidationRequest) error

	//VerifyKeychainStorageRequestSignature checks the keychain storage request's signature
	//Using the share robot public key
	VerifyKeychainStorageRequestSignature(pubKey string, req *api.KeychainStorageRequest) error

	//VerifyIDStorageRequestSignature checks the ID storage request's signature
	//Using the share robot public key
	VerifyIDStorageRequestSignature(pubKey string, req *api.IDStorageRequest) error

	//VerifyLockRequestSignature checks the lock request's signature using the share robot public key
	VerifyLockRequestSignature(pubKey string, req *api.LockRequest) error

	//VerifyKeychainLeadRequestSignature checks the keychain lead mining request's signature
	//Using the share robot public key
	VerifyKeychainLeadRequestSignature(pubKey string, req *api.KeychainLeadRequest) error

	//VerifyIDLeadRequestSignature checks the ID lead mining request's signature
	//Using the share robot public key
	VerifyIDLeadRequestSignature(pubKey string, req *api.IDLeadRequest) error

	//VerifyValidationResponseSignature checks the signature of a validation response using the share robot public key
	VerifyValidationResponseSignature(pubKey string, res *api.ValidationResponse) error

	//VerifyHashSignature checks the signature of a hash using the shared robot public key
	VerifyHashSignature(pubKey string, hash string, sig string) error

	//VerifyStorageAckSignature checks the signature of a storage ack using the shared robot public key
	VerifyStorageAckSignature(pubKey string, ack *api.StorageAck) error

	//VerifyLockAckSignature checks the signature of a lock ack using the shared robot public key
	VerifyLockAckSignature(pubKey string, ack *api.LockAck) error

	//VerifyKeychainResponseSignature checks the signature of a keychain response using the shared robot public key
	VerifyKeychainResponseSignature(pubKey string, res *api.KeychainResponse) error

	//VerifyIDResponseSignature checks the signature of a ID response using the shared robot public key
	VerifyIDResponseSignature(pubKey string, res *api.IDResponse) error

	VerifyContractLeadRequestSignature(pubKey string, req *api.ContractLeadRequest) error
	VerifyContractStorageRequestSignature(pubKey string, req *api.ContractStorageRequest) error
	VerifyContractValidationRequestSignature(pubKey string, req *api.ContractValidationRequest) error

	VerifyContractMessageLeadRequestSignature(pubKey string, req *api.ContractMessageLeadRequest) error
	VerifyContractMessageStorageRequestSignature(pubKey string, req *api.ContractMessageStorageRequest) error
	VerifyContractMessageValidationRequestSignature(pubKey string, req *api.ContractMessageValidationRequest) error
}

type signatureBuilder interface {

	//SignHash create a signature of hash using the shared robot private key
	SignHash(text string, pvKey string) (string, error)

	//SignIDResponse create a signature of the ID response using the shared robot private key
	SignIDResponse(res *api.IDResponse, pvKey string) error

	//SignKeychainResponse create a signature of the keychain response using the shared robot private key
	SignKeychainResponse(res *api.KeychainResponse, pvKey string) error

	//SignKeychainLeadRequest create a signature of the keychain lead mining's request
	//Using the shared robot private key
	SignKeychainLeadRequest(req *api.KeychainLeadRequest, pvKey string) error

	//SignIDLeadRequest create a signature of the ID lead mining's request
	//Using the shared robot private key
	SignIDLeadRequest(req *api.IDLeadRequest, pvKey string) error

	//SignKeychainValidationRequestSignature create a signature of the keychain validation's request
	//Using the shared robot private key
	SignKeychainValidationRequestSignature(req *api.KeychainValidationRequest, pvKey string) error

	//SignIDValidationRequestSignature create a signature of the ID validation's request
	//Using the shared robot private key
	SignIDValidationRequestSignature(req *api.IDValidationRequest, pvKey string) error

	//SignKeychainStorageRequestSignature create a signature of the keychain storage's request
	//Using the shared robot private key
	SignKeychainStorageRequestSignature(req *api.KeychainStorageRequest, pvKey string) error

	//SignIDStorageRequestSignature create a signature of ID storage'request
	//Using the shared robot private key
	SignIDStorageRequestSignature(req *api.IDStorageRequest, pvKey string) error

	//SignLockRequest create a signature of lock request using the shared robot private key
	SignLockRequest(req *api.LockRequest, pvKey string) error

	//SignValidationResponse create a signature of validation response using the shared robot private key
	SignValidationResponse(res *api.ValidationResponse, pvKey string) error

	//SignLockAck create a signature of a lock ack using the shared robot private key
	SignLockAck(ack *api.LockAck, pvKey string) error

	//SignStorageAck create a signature of a storage ack using the shared robot private key
	SignStorageAck(ack *api.StorageAck, pvKey string) error

	//SignCreationResult create a signature of a transaction creation result using the shared robot private key
	SignCreationResult(res *api.CreationResult, pvKey string) error

	//SignAccountResult create a signature of a account search using the shared robot private key
	SignAccountSearchResult(res *api.AccountSearchResult, pvKey string) error

	SignContractLeadRequest(req *api.ContractLeadRequest, pvKey string) error
	SignContractMessageLeadRequest(req *api.ContractMessageLeadRequest, pvKey string) error
}
