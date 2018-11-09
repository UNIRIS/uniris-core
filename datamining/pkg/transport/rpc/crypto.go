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

	//HashKeychain produces a hash of the keychain
	HashKeychain(account.Keychain) (string, error)

	//HashBiometric produces a hash of the biometric
	HashBiometric(account.Biometric) (string, error)

	//HashKeychainData produces a hash of the account's keychain data
	HashKeychainData(account.KeychainData) (string, error)

	//HashBiometricData produces a hash of the account's biometric data
	HashBiometricData(account.BiometricData) (string, error)

	//HashLock produces a hash of the lock transaction
	HashLock(txLock lock.TransactionLock) (string, error)
}

//Signer define methods to handle signatures
type Signer interface {
	signatureBuilder
	signatureChecker
}

type signatureChecker interface {

	//CheckKeychainValidationRequestSignature checks the keychain validation request's signature
	//Using the share robot public key
	CheckKeychainValidationRequestSignature(pubKey string, req *api.KeychainValidationRequest) error

	//CheckKeychainValidationRequestSignature checks the biometric validation request's signature
	//Using the share robot public key
	CheckBiometricValidationRequestSignature(pubKey string, req *api.BiometricValidationRequest) error

	//CheckKeychainStorageRequestSignature checks the keychain storage request's signature
	//Using the share robot public key
	CheckKeychainStorageRequestSignature(pubKey string, req *api.KeychainStorageRequest) error

	//CheckBiometricStorageRequestSignature checks the biometric storage request's signature
	//Using the share robot public key
	CheckBiometricStorageRequestSignature(pubKey string, req *api.BiometricStorageRequest) error

	//CheckLockRequestSignature checks the lock request's signature using the share robot public key
	CheckLockRequestSignature(pubKey string, req *api.LockRequest) error

	//CheckKeychainLeadRequestSignature checks the keychain lead mining request's signature
	//Using the share robot public key
	CheckKeychainLeadRequestSignature(pubKey string, req *api.KeychainLeadRequest) error

	//CheckBiometricLeadRequestSignature checks the biometric lead mining request's signature
	//Using the share robot public key
	CheckBiometricLeadRequestSignature(pubKey string, req *api.BiometricLeadRequest) error

	//CheckValidationResponseSignature checks the signature of a validation response using the share robot public key
	CheckValidationResponseSignature(pubKey string, res *api.ValidationResponse) error

	//CheckHashSignature checks the signature of a hash using the shared robot public key
	CheckHashSignature(pubKey string, hash string, sig string) error

	//CheckStorageAckSignature checks the signature of a storage ack using the shared robot public key
	CheckStorageAckSignature(pubKey string, ack *api.StorageAck) error

	//CheckLockAckSignature checks the signature of a lock ack using the shared robot public key
	CheckLockAckSignature(pubKey string, ack *api.LockAck) error

	//CheckKeychainResponseSignature checks the signature of a keychain response using the shared robot public key
	CheckKeychainResponseSignature(pubKey string, res *api.KeychainResponse) error

	//CheckBiometricResponseSignature checks the signature of a biometric response using the shared robot public key
	CheckBiometricResponseSignature(pubKey string, res *api.BiometricResponse) error
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
