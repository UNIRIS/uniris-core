package rpc

import (
	api "github.com/uniris/uniris-core/datamining/api/protobuf-spec"
	"github.com/uniris/uniris-core/datamining/pkg/account"
	"github.com/uniris/uniris-core/datamining/pkg/lock"
	"github.com/uniris/uniris-core/datamining/pkg/mining"
)

type mockHasher struct{}

func (h mockHasher) HashKeychain(account.Keychain) (string, error) {
	return "hash", nil
}
func (h mockHasher) HashBiometric(account.Biometric) (string, error) {
	return "hash", nil
}

func (h mockHasher) NewKeychainDataHash(account.KeychainData) (string, error) {
	return "hash", nil
}
func (h mockHasher) NewBiometricDataHash(account.BiometricData) (string, error) {
	return "hash", nil
}

func (h mockHasher) HashKeychainData(account.KeychainData) (string, error) {
	return "hash", nil
}
func (h mockHasher) HashBiometricData(account.BiometricData) (string, error) {
	return "hash", nil
}

func (h mockHasher) HashLock(lock.TransactionLock) (string, error) {
	return "hash", nil
}

type mockDecrypter struct{}

func (d mockDecrypter) DecryptHash(hash string, pvKey string) (string, error) {
	return "hash", nil
}

func (d mockDecrypter) DecryptKeychainData(data string, pvKey string) (account.KeychainData, error) {
	return account.NewKeychainData("", "", "", "", account.NewSignatures("", "")), nil
}

func (d mockDecrypter) DecryptBiometricData(data string, pvKey string) (account.BiometricData, error) {
	return account.NewBiometricData("personHash", "", "", "", "", "", account.NewSignatures("", "")), nil
}

type mockSigner struct{}

func (s mockSigner) CheckTransactionDataSignature(txType mining.TransactionType, pubKey string, data interface{}, sig string) error {
	return nil
}

func (s mockSigner) CheckHashSignature(pubKey string, text string, sig string) error {
	return nil
}

func (s mockSigner) CheckKeychainLeadRequestSignature(pubKey string, data *api.KeychainLeadRequest) error {
	return nil
}
func (s mockSigner) CheckBiometricLeadRequestSignature(pubKey string, data *api.BiometricLeadRequest) error {
	return nil
}

func (s mockSigner) CheckKeychainValidationRequestSignature(pubKey string, data *api.KeychainValidationRequest) error {
	return nil
}
func (s mockSigner) CheckBiometricValidationRequestSignature(pubKey string, data *api.BiometricValidationRequest) error {
	return nil
}
func (s mockSigner) CheckKeychainStorageRequestSignature(pubKey string, data *api.KeychainStorageRequest) error {
	return nil
}
func (s mockSigner) CheckBiometricStorageRequestSignature(pubKey string, data *api.BiometricStorageRequest) error {
	return nil
}

func (s mockSigner) CheckLockRequestSignature(pubkey string, req *api.LockRequest) error {
	return nil
}

func (s mockSigner) CheckBiometricDataSignature(pubKey string, data account.BiometricData, sig string) error {
	return nil
}

func (s mockSigner) CheckKeychainDataSignature(pubKey string, data account.KeychainData, sig string) error {
	return nil
}

func (s mockSigner) CheckLockAckSignature(pubKey string, ack *api.LockAck) error {
	return nil
}

func (s mockSigner) CheckStorageAckSignature(pubKey string, ack *api.StorageAck) error {
	return nil
}

func (s mockSigner) CheckValidationResponseSignature(pubKey string, res *api.ValidationResponse) error {
	return nil
}

func (s mockSigner) CheckKeychainResponseSignature(pubKey string, res *api.KeychainResponse) error {
	return nil
}

func (s mockSigner) CheckBiometricResponseSignature(pubKey string, res *api.BiometricResponse) error {
	return nil
}

func (s mockSigner) SignBiometricResponse(res *api.BiometricResponse, pvKey string) error {
	res.Signature = "sig"
	return nil
}
func (s mockSigner) SignKeychainResponse(res *api.KeychainResponse, pvKey string) error {
	res.Signature = "sig"
	return nil
}

func (s mockSigner) SignValidation(v mining.Validation, pvKey string) (string, error) {
	return "sig", nil
}

func (s mockSigner) SignValidationResponse(res *api.ValidationResponse, pvKey string) error {
	res.Signature = "sig"
	return nil
}

func (s mockSigner) SignKeychainLeadRequest(req *api.KeychainLeadRequest, pvKey string) error {
	req.SignatureRequest = "sig"
	return nil
}
func (s mockSigner) SignBiometricLeadRequest(req *api.BiometricLeadRequest, pvKey string) error {
	req.SignatureRequest = "sig"
	return nil
}
func (s mockSigner) SignKeychainValidationRequestSignature(req *api.KeychainValidationRequest, pvKey string) error {
	req.Signature = "sig"
	return nil
}
func (s mockSigner) SignBiometricValidationRequestSignature(req *api.BiometricValidationRequest, pvKey string) error {
	req.Signature = "sig"
	return nil
}
func (s mockSigner) SignKeychainStorageRequestSignature(req *api.KeychainStorageRequest, pvKey string) error {
	req.Signature = "sig"
	return nil
}
func (s mockSigner) SignBiometricStorageRequestSignature(req *api.BiometricStorageRequest, pvKey string) error {
	req.Signature = "sig"
	return nil
}
func (s mockSigner) SignLockRequest(req *api.LockRequest, pvKey string) error {
	req.Signature = "sig"
	return nil
}

func (s mockSigner) SignHash(hash string, pvKey string) (string, error) {
	return "sig", nil
}

func (s mockSigner) SignLockAck(ack *api.LockAck, pvKey string) error {
	ack.Signature = "sig"
	return nil
}

func (s mockSigner) SignStorageAck(ack *api.StorageAck, pvKey string) error {
	ack.Signature = "sig"
	return nil
}

func (s mockSigner) SignAccountSearchResult(res *api.AccountSearchResult, pvKey string) error {
	res.Signature = "sig"
	return nil
}

func (s mockSigner) SignCreationResult(res *api.CreationResult, pvKey string) error {
	res.Signature = "sig"
	return nil
}
