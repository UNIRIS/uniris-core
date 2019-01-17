package crypto

import (
	"crypto/ecdsa"
	"crypto/rand"
	"crypto/sha256"
	"crypto/x509"
	"encoding/asn1"
	"encoding/hex"
	"encoding/json"
	"errors"
	"math/big"

	api "github.com/uniris/uniris-core/api/pkg"
	"github.com/uniris/uniris-core/api/pkg/adding"
	"github.com/uniris/uniris-core/api/pkg/listing"
)

//ErrInvalidSignature is returned when the request contains invalid signatures
var ErrInvalidSignature = errors.New("Invalid signature")

type ecdsaSignature struct {
	R, S *big.Int
}

type signer struct {
}

//NewSigner create a new signer
func NewSigner() signer {
	return signer{}
}

func (s signer) VerifyAccountCreationRequestSignature(req adding.AccountCreationRequest, pubKey string) error {
	b, err := json.Marshal(accountCreationRequest{
		EncryptedID:       req.EncryptedID(),
		EncryptedKeychain: req.EncryptedKeychain(),
	})
	if err != nil {
		return err
	}

	return verifySignature(pubKey, string(b), req.Signature())
}

func (s signer) VerifyHashSignature(hashedData string, pubKey string, sig string) error {
	return verifySignature(pubKey, hashedData, sig)
}

func (s signer) VerifyAccountResultSignature(res listing.AccountResult, pubKey string) error {
	b, err := json.Marshal(accountResult{
		EncryptedAddress: res.EncryptedAddress(),
		EncryptedAESKey:  res.EncryptedAESKey(),
		EncryptedWallet:  res.EncryptedWallet(),
	})
	if err != nil {
		return err
	}
	return verifySignature(pubKey, string(b), res.Signature())
}

func (s signer) VerifyCreationTransactionResultSignature(res api.TransactionResult, pubKey string) error {
	b, err := json.Marshal(transactionResult{
		MasterPeerIP:    res.MasterPeerIP(),
		TransactionHash: res.TransactionHash(),
	})
	if err != nil {
		return err
	}
	return verifySignature(pubKey, string(b), res.Signature())
}

func (s signer) VerifyContractCreationRequestSignature(req adding.ContractCreationRequest, pubKey string) error {
	b, err := json.Marshal(contractCreationRequest{
		EncryptedContract: req.EncryptedContract(),
	})
	if err != nil {
		return err
	}
	return verifySignature(pubKey, string(b), req.Signature())
}

func (s signer) SignAccountCreationResult(res adding.AccountCreationResult, pvKey string) (adding.AccountCreationResult, error) {
	b, err := json.Marshal(accountCreationResult{
		Transactions: accountCreationTransactionsResult{
			ID: transactionResult{
				MasterPeerIP:    res.ResultTransactions().ID().MasterPeerIP(),
				Signature:       res.ResultTransactions().ID().Signature(),
				TransactionHash: res.ResultTransactions().ID().TransactionHash(),
			},
			Keychain: transactionResult{
				MasterPeerIP:    res.ResultTransactions().Keychain().MasterPeerIP(),
				TransactionHash: res.ResultTransactions().Keychain().TransactionHash(),
				Signature:       res.ResultTransactions().Keychain().Signature(),
			},
		},
	})
	sig, err := sign(pvKey, string(b))
	if err != nil {
		return nil, err
	}

	return adding.NewAccountCreationResult(res.ResultTransactions(), sig), nil
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

func verifySignature(pubk string, data string, sig string) error {
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

	return ErrInvalidSignature
}

func hashString(data string) string {
	hash := sha256.New()
	hash.Write([]byte(data))
	return hex.EncodeToString(hash.Sum(nil))
}
