package mock

import (
	api "github.com/uniris/uniris-core/datamining/api/protobuf-spec"
	"github.com/uniris/uniris-core/datamining/pkg/account"
	"github.com/uniris/uniris-core/datamining/pkg/mining"
)

type mockSigner struct{}

//NewSigner create new mocked signer
func NewSigner() mockSigner {
	return mockSigner{}
}

func (s mockSigner) VerifyTransactionDataSignature(txType mining.TransactionType, pubKey string, data interface{}, sig string) error {
	return nil
}

func (s mockSigner) VerifyHashSignature(pubKey string, text string, sig string) error {
	return nil
}

func (s mockSigner) VerifyKeychainLeadRequestSignature(pubKey string, data *api.KeychainLeadRequest) error {
	return nil
}
func (s mockSigner) VerifyIDLeadRequestSignature(pubKey string, data *api.IDLeadRequest) error {
	return nil
}

func (s mockSigner) VerifyKeychainValidationRequestSignature(pubKey string, data *api.KeychainValidationRequest) error {
	return nil
}
func (s mockSigner) VerifyIDValidationRequestSignature(pubKey string, data *api.IDValidationRequest) error {
	return nil
}
func (s mockSigner) VerifyKeychainStorageRequestSignature(pubKey string, data *api.KeychainStorageRequest) error {
	return nil
}
func (s mockSigner) VerifyIDStorageRequestSignature(pubKey string, data *api.IDStorageRequest) error {
	return nil
}

func (s mockSigner) VerifyLockRequestSignature(pubkey string, req *api.LockRequest) error {
	return nil
}

func (s mockSigner) VerifyIDSignatures(account.ID) error {
	return nil
}

func (s mockSigner) VerifyKeychainSignatures(account.Keychain) error {
	return nil
}

func (s mockSigner) VerifyLockAckSignature(pubKey string, ack *api.LockAck) error {
	return nil
}

func (s mockSigner) VerifyStorageAckSignature(pubKey string, ack *api.StorageAck) error {
	return nil
}

func (s mockSigner) VerifyValidationResponseSignature(pubKey string, res *api.ValidationResponse) error {
	return nil
}

func (s mockSigner) VerifyKeychainResponseSignature(pubKey string, res *api.KeychainResponse) error {
	return nil
}

func (s mockSigner) VerifyIDResponseSignature(pubKey string, res *api.IDResponse) error {
	return nil
}

func (s mockSigner) VerifyValidationSignature(v mining.Validation) error {
	return nil
}

func (s mockSigner) SignIDResponse(res *api.IDResponse, pvKey string) error {
	res.Signature = "sig"
	return nil
}
func (s mockSigner) SignKeychainResponse(res *api.KeychainResponse, pvKey string) error {
	res.Signature = "sig"
	return nil
}

func (s mockSigner) SignValidation(v mining.Validation, pvKey string) (mining.Validation, error) {
	return mining.NewValidation(v.Status(), v.Timestamp(), v.PublicKey(), "sig"), nil
}

func (s mockSigner) SignValidationResponse(res *api.ValidationResponse, pvKey string) error {
	res.Signature = "sig"
	return nil
}

func (s mockSigner) SignKeychainLeadRequest(req *api.KeychainLeadRequest, pvKey string) error {
	req.SignatureRequest = "sig"
	return nil
}
func (s mockSigner) SignIDLeadRequest(req *api.IDLeadRequest, pvKey string) error {
	req.SignatureRequest = "sig"
	return nil
}
func (s mockSigner) SignKeychainValidationRequestSignature(req *api.KeychainValidationRequest, pvKey string) error {
	req.Signature = "sig"
	return nil
}
func (s mockSigner) SignIDValidationRequestSignature(req *api.IDValidationRequest, pvKey string) error {
	req.Signature = "sig"
	return nil
}
func (s mockSigner) SignKeychainStorageRequestSignature(req *api.KeychainStorageRequest, pvKey string) error {
	req.Signature = "sig"
	return nil
}
func (s mockSigner) SignIDStorageRequestSignature(req *api.IDStorageRequest, pvKey string) error {
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
