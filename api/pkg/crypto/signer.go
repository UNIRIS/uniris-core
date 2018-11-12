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

type signatures struct {
	BiodSig   string `json:"biod_sig"`
	PersonSig string `json:"person_sig"`
}

type accountCreationTransactions struct {
	Biometric creationTransaction `json:"biometric"`
	Keychain  creationTransaction `json:"keychain"`
}

type creationTransaction struct {
	TransactionHash string `json:"transaction_hash"`
	MasterPeerIP    string `json:"master_peer_ip"`
	Signature       string `json:"signature"`
}

//Signer define methods to handle signatures
type Signer interface {
	rpc.SignatureHandler
	adding.SignatureChecker
	listing.SignatureChecker
}

type signer struct {
}

//NewSigner create a new signer
func NewSigner() Signer {
	return signer{}
}

func (s signer) CheckAccountCreationRequestSignature(data adding.AccountCreationRequest, pubKey string) error {
	b, err := json.Marshal(struct {
		EncryptedBioData      string     `json:"encrypted_bio_data"`
		EncryptedKeychainData string     `json:"encrypted_keychain_data"`
		SignaturesBio         signatures `json:"signatures_bio"`
		SignaturesKeychain    signatures `json:"signatures_keychain"`
	}{
		EncryptedBioData:      data.EncryptedBioData,
		EncryptedKeychainData: data.EncryptedKeychainData,
		SignaturesBio: signatures{
			BiodSig:   data.SignaturesBio.BiodSig,
			PersonSig: data.SignaturesBio.PersonSig,
		},
		SignaturesKeychain: signatures{
			BiodSig:   data.SignaturesKeychain.BiodSig,
			PersonSig: data.SignaturesKeychain.PersonSig,
		},
	})
	if err != nil {
		return err
	}

	return checkSignature(pubKey, string(b), data.SignatureRequest)
}

func (s signer) CheckHashSignature(hashedData string, pubKey string, sig string) error {
	return checkSignature(pubKey, hashedData, sig)
}

func (s signer) CheckAccountSearchResultSignature(pubKey string, res *api.AccountSearchResult) error {
	b, err := json.Marshal(&api.AccountSearchResult{
		EncryptedAddress: res.EncryptedAddress,
		EncryptedAESkey:  res.EncryptedAESkey,
		EncryptedWallet:  res.EncryptedWallet,
	})
	if err != nil {
		return err
	}
	return checkSignature(pubKey, string(b), res.Signature)
}
func (s signer) CheckCreationResultSignature(pubKey string, res *api.CreationResult) error {
	b, err := json.Marshal(&api.CreationResult{
		MasterPeerIP:    res.MasterPeerIP,
		TransactionHash: res.TransactionHash,
	})
	if err != nil {
		return err
	}
	return checkSignature(pubKey, string(b), res.Signature)
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
	res.SignatureRequest = sig
	return nil
}
func (s signer) SignAccountCreationResult(res *adding.AccountCreationResult, pvKey string) error {
	b, err := json.Marshal(struct {
		Transactions accountCreationTransactions `json:"transactions"`
	}{
		Transactions: accountCreationTransactions{
			Biometric: creationTransaction{
				MasterPeerIP:    res.Transactions.Biometric.MasterPeerIP,
				Signature:       res.Transactions.Biometric.Signature,
				TransactionHash: res.Transactions.Biometric.TransactionHash,
			},
			Keychain: creationTransaction{
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

	return ErrInvalidSignature
}

func hashString(data string) string {
	hash := sha256.New()
	hash.Write([]byte(data))
	return hex.EncodeToString(hash.Sum(nil))
}
