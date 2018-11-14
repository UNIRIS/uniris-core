package rpc

import (
	"errors"

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

	//DecryptKeychainData decrypt the account's keychain using the shared robot private key
	DecryptKeychainData(data string, pvKey string) (account.KeychainData, error)

	//DecryptBiometricData decrypt the account's biometric using the shared robot private key
	DecryptBiometricData(data string, pvKey string) (account.BiometricData, error)
}

//Hasher define methods to hash incoming data
type Hasher interface {
	account.KeychainHasher
	account.BiometricHasher
	lock.Hasher

	//HashBiodPublicKey produces a hash of the decrypted biometric device's public key
	HashBiodPublicKey(string) string
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

	//VerifyKeychainValidationRequestSignature checks the biometric validation request's signature
	//Using the share robot public key
	VerifyBiometricValidationRequestSignature(pubKey string, req *api.BiometricValidationRequest) error

	//VerifyKeychainStorageRequestSignature checks the keychain storage request's signature
	//Using the share robot public key
	VerifyKeychainStorageRequestSignature(pubKey string, req *api.KeychainStorageRequest) error

	//VerifyBiometricStorageRequestSignature checks the biometric storage request's signature
	//Using the share robot public key
	VerifyBiometricStorageRequestSignature(pubKey string, req *api.BiometricStorageRequest) error

	//VerifyLockRequestSignature checks the lock request's signature using the share robot public key
	VerifyLockRequestSignature(pubKey string, req *api.LockRequest) error

	//VerifyKeychainLeadRequestSignature checks the keychain lead mining request's signature
	//Using the share robot public key
	VerifyKeychainLeadRequestSignature(pubKey string, req *api.KeychainLeadRequest) error

	//VerifyBiometricLeadRequestSignature checks the biometric lead mining request's signature
	//Using the share robot public key
	VerifyBiometricLeadRequestSignature(pubKey string, req *api.BiometricLeadRequest) error

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

	//VerifyBiometricResponseSignature checks the signature of a biometric response using the shared robot public key
	VerifyBiometricResponseSignature(pubKey string, res *api.BiometricResponse) error
}

type signatureBuilder interface {

	//SignHash create a signature of hash using the shared robot private key
	SignHash(text string, pvKey string) (string, error)

	//SignBiometricResponse create a signature of the biometric response using the shared robot private key
	SignBiometricResponse(res *api.BiometricResponse, pvKey string) error

	//SignKeychainResponse create a signature of the keychain response using the shared robot private key
	SignKeychainResponse(res *api.KeychainResponse, pvKey string) error

	//SignKeychainLeadRequest create a signature of the keychain lead mining's request
	//Using the shared robot private key
	SignKeychainLeadRequest(req *api.KeychainLeadRequest, pvKey string) error

	//SignBiometricLeadRequest create a signature of the biometric lead mining's request
	//Using the shared robot private key
	SignBiometricLeadRequest(req *api.BiometricLeadRequest, pvKey string) error

	//SignKeychainValidationRequestSignature create a signature of the keychain validation's request
	//Using the shared robot private key
	SignKeychainValidationRequestSignature(req *api.KeychainValidationRequest, pvKey string) error

	//SignBiometricValidationRequestSignature create a signature of the biometric validation's request
	//Using the shared robot private key
	SignBiometricValidationRequestSignature(req *api.BiometricValidationRequest, pvKey string) error

	//SignKeychainStorageRequestSignature create a signature of the keychain storage's request
	//Using the shared robot private key
	SignKeychainStorageRequestSignature(req *api.KeychainStorageRequest, pvKey string) error

	//SignBiometricStorageRequestSignature create a signature of biometric storage'request
	//Using the shared robot private key
	SignBiometricStorageRequestSignature(req *api.BiometricStorageRequest, pvKey string) error

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
}
