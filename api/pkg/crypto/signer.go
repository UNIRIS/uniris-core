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

	"github.com/uniris/uniris-core/api/pkg/adding"
	"github.com/uniris/uniris-core/api/pkg/listing"
	"github.com/uniris/uniris-core/api/pkg/transport/rpc"
	api "github.com/uniris/uniris-core/datamining/api/protobuf-spec"
)

//ErrInvalidSignature is returned when the request contains invalid signatures
var ErrInvalidSignature = errors.New("Invalid signature")

type ecdsaSignature struct {
	R, S *big.Int
}

//Signer define methods to handle signatures
type Signer interface {
	rpc.SignatureHandler
	adding.SignatureVerifier
	listing.SignatureVerifier
}

type signer struct {
}

//NewSigner create a new signer
func NewSigner() Signer {
	return signer{}
}

func (s signer) VerifyAccountCreationRequestSignature(data adding.AccountCreationRequest, pubKey string) error {
	b, err := json.Marshal(adding.AccountCreationRequest{
		EncryptedID:       data.EncryptedID,
		EncryptedKeychain: data.EncryptedKeychain,
	})
	if err != nil {
		return err
	}

	return verifySignature(pubKey, string(b), data.Signature)
}

func (s signer) VerifyHashSignature(hashedData string, pubKey string, sig string) error {
	return verifySignature(pubKey, hashedData, sig)
}

func (s signer) VerifyAccountSearchResultSignature(pubKey string, res *api.AccountSearchResult) error {
	b, err := json.Marshal(&api.AccountSearchResult{
		EncryptedAddress: res.EncryptedAddress,
		EncryptedAESkey:  res.EncryptedAESkey,
		EncryptedWallet:  res.EncryptedWallet,
	})
	if err != nil {
		return err
	}
	return verifySignature(pubKey, string(b), res.Signature)
}
func (s signer) VerifyCreationResultSignature(pubKey string, res *api.CreationResult) error {
	b, err := json.Marshal(&api.CreationResult{
		MasterPeerIP:    res.MasterPeerIP,
		TransactionHash: res.TransactionHash,
	})
	if err != nil {
		return err
	}
	return verifySignature(pubKey, string(b), res.Signature)
}

func (s signer) SignAccountResult(res *listing.AccountResult, pvKey string) error {
	b, err := json.Marshal(struct {
		EncryptedAESKey  string `json:"encrypted_aes_key"`
		EncryptedWallet  string `json:"encrypted_wallet"`
		EncryptedAddress string `json:"encrypted_address"`
	}{
		EncryptedAddress: res.EncryptedAddress,
		EncryptedAESKey:  res.EncryptedAESKey,
		EncryptedWallet:  res.EncryptedWallet,
	})
	sig, err := sign(pvKey, string(b))
	if err != nil {
		return err
	}
	res.Signature = sig
	return nil
}
func (s signer) SignAccountCreationResult(res *adding.AccountCreationResult, pvKey string) error {
	b, err := json.Marshal(adding.AccountCreationResult{
		Transactions: adding.AccountCreationTransactionsResult{
			ID: adding.TransactionResult{
				MasterPeerIP:    res.Transactions.ID.MasterPeerIP,
				Signature:       res.Transactions.ID.Signature,
				TransactionHash: res.Transactions.ID.TransactionHash,
			},
			Keychain: adding.TransactionResult{
				MasterPeerIP:    res.Transactions.Keychain.MasterPeerIP,
				TransactionHash: res.Transactions.Keychain.TransactionHash,
				Signature:       res.Transactions.Keychain.Signature,
			},
		},
	})
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
