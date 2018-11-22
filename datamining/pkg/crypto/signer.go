package crypto

import (
	"crypto/ecdsa"
	"crypto/rand"
	"crypto/x509"
	"encoding/asn1"
	"encoding/hex"
	"encoding/json"
	"errors"
	"math/big"

	api "github.com/uniris/uniris-core/datamining/api/protobuf-spec"
	"github.com/uniris/uniris-core/datamining/pkg/account"
	"github.com/uniris/uniris-core/datamining/pkg/mining"
	"github.com/uniris/uniris-core/datamining/pkg/transport/rpc"
)

type ecdsaSignature struct {
	R, S *big.Int
}

//Signer defines methods to handle signatures
type Signer interface {
	mining.PowSigVerifier
	mining.ValidationVerifier
	mining.ValidationSigner
	account.KeychainSignatureVerifier
	account.IDSignatureVerifier
	rpc.Signer
}

type signer struct{}

//NewSigner creates a new signer
func NewSigner() Signer {
	return signer{}
}

func (s signer) VerifyTransactionDataSignature(txType mining.TransactionType, pubKey string, data interface{}, sig string) error {
	switch txType {
	case mining.KeychainTransaction:
		kc := data.(account.Keychain)
		b, err := json.Marshal(keychainWithoutSig{
			EncryptedWallet:      kc.EncryptedWallet(),
			EncryptedAddrByRobot: kc.EncryptedAddrByRobot(),
			IDPublicKey:          kc.IDPublicKey(),
		})
		if err != nil {
			return err
		}
		return checkSignature(pubKey, string(b), sig)
	case mining.IDTransaction:
		id := data.(account.ID)
		b, err := json.Marshal(idWithoutSig{
			EncryptedAddrByID:    id.EncryptedAddrByID(),
			EncryptedAddrByRobot: id.EncryptedAddrByRobot(),
			EncryptedAESKey:      id.EncryptedAESKey(),
			Hash:                 id.Hash(),
			PublicKey:            id.PublicKey(),
		})
		if err != nil {
			return err
		}
		return checkSignature(pubKey, string(b), sig)
	}

	return mining.ErrUnsupportedTransaction
}

func (s signer) VerifyIDSignatures(id account.ID) error {
	b, err := json.Marshal(idWithoutSig{
		EncryptedAddrByID:    id.EncryptedAddrByID(),
		EncryptedAddrByRobot: id.EncryptedAddrByRobot(),
		EncryptedAESKey:      id.EncryptedAESKey(),
		Hash:                 id.Hash(),
		PublicKey:            id.PublicKey(),
	})
	if err != nil {
		return err
	}
	if err := checkSignature(id.PublicKey(), string(b), id.IDSignature()); err != nil {
		return err
	}
	return nil
}

func (s signer) VerifyKeychainSignatures(kc account.Keychain) error {
	b, err := json.Marshal(keychainWithoutSig{
		EncryptedWallet:      kc.EncryptedWallet(),
		EncryptedAddrByRobot: kc.EncryptedAddrByRobot(),
		IDPublicKey:          kc.IDPublicKey(),
	})
	if err != nil {
		return err
	}
	if err := checkSignature(kc.IDPublicKey(), string(b), kc.IDSignature()); err != nil {
		return err
	}
	return nil
}

func (s signer) VerifyHashSignature(pubKey string, hash string, sig string) error {
	return checkSignature(pubKey, hash, sig)
}

func (s signer) VerifyKeychainValidationRequestSignature(pubKey string, req *api.KeychainValidationRequest) error {
	b, err := json.Marshal(&api.KeychainValidationRequest{
		Data:            req.Data,
		TransactionHash: req.TransactionHash,
	})
	if err != nil {
		return err
	}
	return checkSignature(pubKey, string(b), req.Signature)
}

func (s signer) VerifyIDValidationRequestSignature(pubKey string, req *api.IDValidationRequest) error {
	b, err := json.Marshal(&api.IDValidationRequest{
		Data:            req.Data,
		TransactionHash: req.TransactionHash,
	})
	if err != nil {
		return err
	}
	return checkSignature(pubKey, string(b), req.Signature)
}

func (s signer) VerifyKeychainStorageRequestSignature(pubKey string, req *api.KeychainStorageRequest) error {
	b, err := json.Marshal(&api.KeychainStorageRequest{
		Data:        req.Data,
		Endorsement: req.Endorsement,
	})
	if err != nil {
		return err
	}
	return checkSignature(pubKey, string(b), req.Signature)
}

func (s signer) VerifyIDStorageRequestSignature(pubKey string, req *api.IDStorageRequest) error {
	b, err := json.Marshal(&api.IDStorageRequest{
		Data:        req.Data,
		Endorsement: req.Endorsement,
	})
	if err != nil {
		return err
	}
	return checkSignature(pubKey, string(b), req.Signature)
}

func (s signer) VerifyLockRequestSignature(pubKey string, req *api.LockRequest) error {
	b, err := json.Marshal(&api.LockRequest{
		Address:         req.Address,
		MasterRobotKey:  req.MasterRobotKey,
		TransactionHash: req.TransactionHash,
	})
	if err != nil {
		return err
	}
	return checkSignature(pubKey, string(b), req.Signature)
}

func (s signer) VerifyKeychainLeadRequestSignature(pubKey string, req *api.KeychainLeadRequest) error {
	b, err := json.Marshal(&api.KeychainLeadRequest{
		EncryptedKeychain: req.EncryptedKeychain,
		TransactionHash:   req.TransactionHash,
		ValidatorPeerIPs:  req.ValidatorPeerIPs,
	})
	if err != nil {
		return err
	}
	return checkSignature(pubKey, string(b), req.SignatureRequest)
}

func (s signer) VerifyIDLeadRequestSignature(pubKey string, req *api.IDLeadRequest) error {
	b, err := json.Marshal(&api.IDLeadRequest{
		EncryptedID:      req.EncryptedID,
		TransactionHash:  req.TransactionHash,
		ValidatorPeerIPs: req.ValidatorPeerIPs,
	})
	if err != nil {
		return err
	}
	return checkSignature(pubKey, string(b), req.SignatureRequest)
}

func (s signer) VerifyValidationResponseSignature(pubKey string, res *api.ValidationResponse) error {
	b, err := json.Marshal(&api.ValidationResponse{
		Validation: res.Validation,
	})
	if err != nil {
		return err
	}
	return checkSignature(pubKey, string(b), res.Signature)
}

func (s signer) VerifyLockAckSignature(pubKey string, ack *api.LockAck) error {
	b, err := json.Marshal(&api.LockAck{
		LockHash: ack.LockHash,
	})
	if err != nil {
		return err
	}
	return checkSignature(pubKey, string(b), ack.Signature)
}

func (s signer) VerifyStorageAckSignature(pubKey string, ack *api.StorageAck) error {
	b, err := json.Marshal(&api.StorageAck{
		StorageHash: ack.StorageHash,
	})
	if err != nil {
		return err
	}
	return checkSignature(pubKey, string(b), ack.Signature)
}

func (s signer) VerifyKeychainResponseSignature(pubKey string, res *api.KeychainResponse) error {
	b, err := json.Marshal(&api.KeychainResponse{
		Data:        res.Data,
		Endorsement: res.Endorsement,
	})
	if err != nil {
		return err
	}
	return checkSignature(pubKey, string(b), res.Signature)
}

func (s signer) VerifyIDResponseSignature(pubKey string, res *api.IDResponse) error {
	b, err := json.Marshal(&api.IDResponse{
		Data:        res.Data,
		Endorsement: res.Endorsement,
	})
	if err != nil {
		return err
	}
	return checkSignature(pubKey, string(b), res.Signature)
}

func (s signer) VerifyValidationSignature(v mining.Validation) error {
	b, err := json.Marshal(validationWithoutSig{
		PublicKey: v.PublicKey(),
		Status:    v.Status(),
		Timestamp: v.Timestamp().Unix(),
	})
	if err != nil {
		return err
	}

	return checkSignature(v.PublicKey(), string(b), v.Signature())
}

func (s signer) SignIDResponse(res *api.IDResponse, pvKey string) error {
	b, err := json.Marshal(&api.IDResponse{
		Data:        res.Data,
		Endorsement: res.Endorsement,
	})
	if err != nil {
		return err
	}
	sig, err := sign(pvKey, string(b))
	if err != nil {
		return err
	}
	res.Signature = sig
	return nil
}

func (s signer) SignKeychainResponse(res *api.KeychainResponse, pvKey string) error {
	b, err := json.Marshal(&api.KeychainResponse{
		Data:        res.Data,
		Endorsement: res.Endorsement,
	})
	if err != nil {
		return err
	}
	sig, err := sign(pvKey, string(b))
	if err != nil {
		return err
	}
	res.Signature = sig
	return nil
}

func (s signer) SignHash(hash string, pvKey string) (string, error) {
	return sign(pvKey, hash)
}

func (s signer) SignKeychainLeadRequest(req *api.KeychainLeadRequest, pvKey string) error {
	b, err := json.Marshal(&api.KeychainLeadRequest{
		EncryptedKeychain: req.EncryptedKeychain,
		TransactionHash:   req.TransactionHash,
		ValidatorPeerIPs:  req.ValidatorPeerIPs,
	})
	if err != nil {
		return err
	}
	sig, err := sign(pvKey, string(b))
	if err != nil {
		return err
	}
	req.SignatureRequest = sig
	return nil
}

func (s signer) SignIDLeadRequest(req *api.IDLeadRequest, pvKey string) error {
	b, err := json.Marshal(&api.IDLeadRequest{
		EncryptedID:      req.EncryptedID,
		TransactionHash:  req.TransactionHash,
		ValidatorPeerIPs: req.ValidatorPeerIPs,
	})
	if err != nil {
		return err
	}
	sig, err := sign(pvKey, string(b))
	if err != nil {
		return err
	}
	req.SignatureRequest = sig
	return nil
}

func (s signer) SignKeychainValidationRequestSignature(req *api.KeychainValidationRequest, pvKey string) error {
	b, err := json.Marshal(&api.KeychainValidationRequest{
		Data:            req.Data,
		TransactionHash: req.TransactionHash,
	})
	if err != nil {
		return err
	}
	sig, err := sign(pvKey, string(b))
	if err != nil {
		return err
	}
	req.Signature = sig
	return nil
}

func (s signer) SignIDValidationRequestSignature(req *api.IDValidationRequest, pvKey string) error {
	b, err := json.Marshal(&api.IDValidationRequest{
		Data:            req.Data,
		TransactionHash: req.TransactionHash,
	})
	if err != nil {
		return err
	}
	sig, err := sign(pvKey, string(b))
	if err != nil {
		return err
	}
	req.Signature = sig
	return nil
}

func (s signer) SignKeychainStorageRequestSignature(req *api.KeychainStorageRequest, pvKey string) error {
	b, err := json.Marshal(&api.KeychainStorageRequest{
		Data:        req.Data,
		Endorsement: req.Endorsement,
	})
	if err != nil {
		return err
	}
	sig, err := sign(pvKey, string(b))
	if err != nil {
		return err
	}
	req.Signature = sig
	return nil
}

func (s signer) SignIDStorageRequestSignature(req *api.IDStorageRequest, pvKey string) error {
	b, err := json.Marshal(&api.IDStorageRequest{
		Data:        req.Data,
		Endorsement: req.Endorsement,
	})
	if err != nil {
		return err
	}
	sig, err := sign(pvKey, string(b))
	if err != nil {
		return err
	}
	req.Signature = sig
	return nil
}

func (s signer) SignLockRequest(req *api.LockRequest, pvKey string) error {
	b, err := json.Marshal(&api.LockRequest{
		Address:         req.Address,
		MasterRobotKey:  req.MasterRobotKey,
		TransactionHash: req.TransactionHash,
	})
	if err != nil {
		return err
	}
	sig, err := sign(pvKey, string(b))
	if err != nil {
		return err
	}
	req.Signature = sig
	return nil
}

func (s signer) SignValidation(v mining.Validation, pvKey string) (mining.Validation, error) {
	b, err := json.Marshal(validationWithoutSig{
		PublicKey: v.PublicKey(),
		Status:    v.Status(),
		Timestamp: v.Timestamp().Unix(),
	})
	if err != nil {
		return nil, err
	}

	sig, err := sign(pvKey, string(b))
	if err != nil {
		return nil, err
	}

	return mining.NewValidation(v.Status(), v.Timestamp(), v.PublicKey(), sig), nil
}

func (s signer) SignValidationResponse(res *api.ValidationResponse, pvKey string) error {
	b, err := json.Marshal(&api.ValidationResponse{
		Validation: res.Validation,
	})
	if err != nil {
		return err
	}
	sig, err := sign(pvKey, string(b))
	if err != nil {
		return err
	}
	res.Signature = sig
	return nil
}

func (s signer) SignLockAck(ack *api.LockAck, pvKey string) error {
	b, err := json.Marshal(&api.LockAck{
		LockHash: ack.LockHash,
	})
	if err != nil {
		return err
	}
	sig, err := sign(pvKey, string(b))
	if err != nil {
		return err
	}
	ack.Signature = sig
	return nil
}

func (s signer) SignStorageAck(ack *api.StorageAck, pvKey string) error {
	b, err := json.Marshal(&api.StorageAck{
		StorageHash: ack.StorageHash,
	})
	if err != nil {
		return err
	}
	sig, err := sign(pvKey, string(b))
	if err != nil {
		return err
	}
	ack.Signature = sig
	return nil
}

func (s signer) SignAccountSearchResult(res *api.AccountSearchResult, pvKey string) error {
	b, err := json.Marshal(&api.AccountSearchResult{
		EncryptedAddress: res.EncryptedAddress,
		EncryptedAESkey:  res.EncryptedAESkey,
		EncryptedWallet:  res.EncryptedWallet,
	})
	if err != nil {
		return err
	}
	sig, err := sign(pvKey, string(b))
	if err != nil {
		return err
	}
	res.Signature = sig
	return nil
}

func (s signer) SignCreationResult(res *api.CreationResult, pvKey string) error {
	b, err := json.Marshal(&api.CreationResult{
		MasterPeerIP:    res.MasterPeerIP,
		TransactionHash: res.TransactionHash,
	})
	if err != nil {
		return err
	}
	sig, err := sign(pvKey, string(b))
	if err != nil {
		return err
	}
	res.Signature = sig
	return nil
}

func sign(privk string, data string) (string, error) {
	pvDecoded, err := hex.DecodeString(privk)
	if err != nil {
		return "", err
	}

	pv, err := x509.ParseECPrivateKey(pvDecoded)
	if err != nil {
		return "", err
	}

	hash := []byte(hashString(data))

	r, s, err := ecdsa.Sign(rand.Reader, pv, hash)
	if err != nil {
		return "", err
	}

	sig, err := asn1.Marshal(ecdsaSignature{r, s})
	if err != nil {
		return "", err
	}

	return hex.EncodeToString(sig), nil

}

func checkSignature(pubk string, data string, sig string) error {
	var signature ecdsaSignature

	decodedkey, err := hex.DecodeString(pubk)
	if err != nil {
		return err
	}

	decodedsig, err := hex.DecodeString(sig)
	if err != nil {
		return err
	}

	pu, err := x509.ParsePKIXPublicKey(decodedkey)
	if err != nil {
		return err
	}

	ecdsaPublic := pu.(*ecdsa.PublicKey)
	asn1.Unmarshal(decodedsig, &signature)

	hash := []byte(hashString(data))

	if ecdsa.Verify(ecdsaPublic, hash, signature.R, signature.S) {
		return nil
	}

	return errors.New("Invalid signature")
}
