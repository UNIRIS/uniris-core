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

	accountMining "github.com/uniris/uniris-core/datamining/pkg/account/mining"
)

type ecdsaSignature struct {
	R, S *big.Int
}

//Signer defines methods to handle signatures
type Signer interface {
	mining.Signer
	accountMining.KeychainSigner
	accountMining.BiometricSigner
	rpc.Signer
}

type signer struct{}

//NewSigner creates a new signer
func NewSigner() Signer {
	return signer{}
}

func (s signer) CheckTransactionDataSignature(txType mining.TransactionType, pubKey string, data interface{}, sig string) error {
	switch txType {
	case mining.KeychainTransaction:
		return s.CheckKeychainDataSignature(pubKey, data.(account.KeychainData), sig)
	case mining.BiometricTransaction:
		return s.CheckBiometricDataSignature(pubKey, data.(account.BiometricData), sig)
	}

	return mining.ErrUnsupportedTransaction
}

func (s signer) CheckBiometricDataSignature(pubKey string, data account.BiometricData, sig string) error {
	b, err := json.Marshal(biometricRaw{
		BIODPublicKey:       data.BiodPublicKey(),
		EncryptedAddrPerson: data.CipherAddrPerson(),
		EncryptedAddrRobot:  data.CipherAddrRobot(),
		EncryptedAESKey:     data.CipherAESKey(),
		PersonHash:          data.PersonHash(),
		PersonPublicKey:     data.PersonPublicKey(),
	})
	if err != nil {
		return err
	}
	return checkSignature(pubKey, string(b), sig)
}

func (s signer) CheckKeychainDataSignature(pubKey string, data account.KeychainData, sig string) error {
	b, err := json.Marshal(keychainRaw{
		BIODPublicKey:      data.BiodPublicKey(),
		EncryptedWallet:    data.CipherWallet(),
		EncryptedAddrRobot: data.CipherAddrRobot(),
		PersonPublicKey:    data.PersonPublicKey(),
	})
	if err != nil {
		return err
	}
	return checkSignature(pubKey, string(b), sig)
}

func (s signer) CheckHashSignature(pubKey string, hash string, sig string) error {
	return checkSignature(pubKey, hash, sig)
}

func (s signer) CheckKeychainValidationRequestSignature(pubKey string, req *api.KeychainValidationRequest) error {
	b, err := json.Marshal(&api.KeychainValidationRequest{
		Data:            req.Data,
		TransactionHash: req.TransactionHash,
	})
	if err != nil {
		return err
	}
	return checkSignature(pubKey, string(b), req.Signature)
}

func (s signer) CheckBiometricValidationRequestSignature(pubKey string, req *api.BiometricValidationRequest) error {
	b, err := json.Marshal(&api.BiometricValidationRequest{
		Data:            req.Data,
		TransactionHash: req.TransactionHash,
	})
	if err != nil {
		return err
	}
	return checkSignature(pubKey, string(b), req.Signature)
}

func (s signer) CheckKeychainStorageRequestSignature(pubKey string, req *api.KeychainStorageRequest) error {
	b, err := json.Marshal(&api.KeychainStorageRequest{
		Data:        req.Data,
		Endorsement: req.Endorsement,
	})
	if err != nil {
		return err
	}
	return checkSignature(pubKey, string(b), req.Signature)
}

func (s signer) CheckBiometricStorageRequestSignature(pubKey string, req *api.BiometricStorageRequest) error {
	b, err := json.Marshal(&api.BiometricStorageRequest{
		Data:        req.Data,
		Endorsement: req.Endorsement,
	})
	if err != nil {
		return err
	}
	return checkSignature(pubKey, string(b), req.Signature)
}

func (s signer) CheckLockRequestSignature(pubKey string, req *api.LockRequest) error {
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

func (s signer) CheckKeychainLeadRequestSignature(pubKey string, req *api.KeychainLeadRequest) error {
	b, err := json.Marshal(&api.KeychainLeadRequest{
		EncryptedKeychainData: req.EncryptedKeychainData,
		SignatureKeychainData: req.SignatureKeychainData,
		TransactionHash:       req.TransactionHash,
		ValidatorPeerIPs:      req.ValidatorPeerIPs,
	})
	if err != nil {
		return err
	}
	return checkSignature(pubKey, string(b), req.SignatureRequest)
}

func (s signer) CheckBiometricLeadRequestSignature(pubKey string, req *api.BiometricLeadRequest) error {
	b, err := json.Marshal(&api.BiometricLeadRequest{
		EncryptedBioData: req.EncryptedBioData,
		SignatureBioData: req.SignatureBioData,
		TransactionHash:  req.TransactionHash,
		ValidatorPeerIPs: req.ValidatorPeerIPs,
	})
	if err != nil {
		return err
	}
	return checkSignature(pubKey, string(b), req.SignatureRequest)
}

func (s signer) CheckValidationResponseSignature(pubKey string, res *api.ValidationResponse) error {
	b, err := json.Marshal(&api.ValidationResponse{
		Validation: res.Validation,
	})
	if err != nil {
		return err
	}
	return checkSignature(pubKey, string(b), res.Signature)
}

func (s signer) CheckLockAckSignature(pubKey string, ack *api.LockAck) error {
	b, err := json.Marshal(&api.LockAck{
		LockHash: ack.LockHash,
	})
	if err != nil {
		return err
	}
	return checkSignature(pubKey, string(b), ack.Signature)
}

func (s signer) CheckStorageAckSignature(pubKey string, ack *api.StorageAck) error {
	b, err := json.Marshal(&api.StorageAck{
		StorageHash: ack.StorageHash,
	})
	if err != nil {
		return err
	}
	return checkSignature(pubKey, string(b), ack.Signature)
}

func (s signer) CheckKeychainResponseSignature(pubKey string, res *api.KeychainResponse) error {
	b, err := json.Marshal(&api.KeychainResponse{
		Data:        res.Data,
		Endorsement: res.Endorsement,
	})
	if err != nil {
		return err
	}
	return checkSignature(pubKey, string(b), res.Signature)
}

func (s signer) CheckBiometricResponseSignature(pubKey string, res *api.BiometricResponse) error {
	b, err := json.Marshal(&api.BiometricResponse{
		Data:        res.Data,
		Endorsement: res.Endorsement,
	})
	if err != nil {
		return err
	}
	return checkSignature(pubKey, string(b), res.Signature)
}

func (s signer) SignBiometricResponse(res *api.BiometricResponse, pvKey string) error {
	b, err := json.Marshal(&api.BiometricResponse{
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
		EncryptedKeychainData: req.EncryptedKeychainData,
		SignatureKeychainData: req.SignatureKeychainData,
		TransactionHash:       req.TransactionHash,
		ValidatorPeerIPs:      req.ValidatorPeerIPs,
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

func (s signer) SignBiometricLeadRequest(req *api.BiometricLeadRequest, pvKey string) error {
	b, err := json.Marshal(&api.BiometricLeadRequest{
		EncryptedBioData: req.EncryptedBioData,
		SignatureBioData: req.SignatureBioData,
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

func (s signer) SignBiometricValidationRequestSignature(req *api.BiometricValidationRequest, pvKey string) error {
	b, err := json.Marshal(&api.BiometricValidationRequest{
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

func (s signer) SignBiometricStorageRequestSignature(req *api.BiometricStorageRequest, pvKey string) error {
	b, err := json.Marshal(&api.BiometricStorageRequest{
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

func (s signer) SignValidation(data mining.Validation, pvKey string) (string, error) {
	b, err := json.Marshal(validationRaw{
		PublicKey: data.PublicKey(),
		Status:    data.Status(),
		Timestamp: data.Timestamp(),
	})
	if err != nil {
		return "", err
	}
	return sign(pvKey, string(b))
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
