package crypto

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/x509"
	"encoding/asn1"
	"encoding/hex"
	"encoding/json"
	"time"

	"testing"

	"github.com/stretchr/testify/assert"
	api "github.com/uniris/uniris-core/datamining/api/protobuf-spec"
	"github.com/uniris/uniris-core/datamining/pkg/account"
	"github.com/uniris/uniris-core/datamining/pkg/mining"
)

/*
Scenario: Sign encrypted data
	Given an encrypted data and a private key
	When I want sign this data
	Then I get the signature and can be verify by the public key associated
*/
func TestSign(t *testing.T) {
	key, _ := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	pvKey, _ := x509.MarshalECPrivateKey(key)
	encData := "uxazexc"

	sig, err := sign(hex.EncodeToString(pvKey), encData)
	assert.Nil(t, err)
	assert.NotEmpty(t, sig)
	var signature ecdsaSignature
	decodesig, _ := hex.DecodeString(string(sig))
	asn1.Unmarshal(decodesig, &signature)

	hash := []byte(hashString(encData))

	assert.True(t, ecdsa.Verify(&key.PublicKey, hash, signature.R, signature.S))
}

/*
Scenario: Verify encrypted data
	Given a data , a signature and a public key
	When I want verify this signature
	Then I get the supposed result
*/
func TestVerify(t *testing.T) {
	key, _ := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	pvKey, _ := x509.MarshalECPrivateKey(key)
	pubKey, _ := x509.MarshalPKIXPublicKey(key.Public())
	encData := struct {
		Message string
	}{Message: "hello"}

	b, _ := json.Marshal(encData)

	sig, _ := sign(hex.EncodeToString(pvKey), string(b))

	err := checkSignature(
		hex.EncodeToString(pubKey),
		string(b),
		sig,
	)
	assert.Nil(t, err)
}

/*
Scenario: Verifies keychain transaction signature
	Given a signed keychain data, a signature and a public key
	When I want to check it's the data matched the signature
	Then I get not error
*/
func TestVerifyTransactionKeychainSignature(t *testing.T) {
	k := account.NewKeychainData("cipher addr", "cipher wallet", "person pub", "biod pub", account.NewSignatures("sig", "sig"))
	key, _ := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	pvKey, _ := x509.MarshalECPrivateKey(key)
	pubKey, _ := x509.MarshalPKIXPublicKey(key.Public())

	data := keychainRaw{
		BIODPublicKey:      k.BiodPublicKey(),
		EncryptedAddrRobot: k.CipherAddrRobot(),
		EncryptedWallet:    k.CipherWallet(),
		PersonPublicKey:    k.PersonPublicKey(),
	}
	b, _ := json.Marshal(data)

	sig, _ := sign(hex.EncodeToString(pvKey), string(b))

	assert.Nil(t, NewSigner().VerifyTransactionDataSignature(mining.KeychainTransaction, hex.EncodeToString(pubKey), k, sig))
}

/*
Scenario: Verify biometric transaction signature
	Given a signed biometric data, a signature and a public key
	When I want to check it's the data matches the signature
	Then I get not error
*/
func TestVerifyTransactionBiometricSignature(t *testing.T) {
	bio := account.NewBiometricData("hash", "cipher addr", "cipher addr", "cipher aes key", "person pub", "biod pub", account.NewSignatures("sig", "sig"))
	key, _ := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	pvKey, _ := x509.MarshalECPrivateKey(key)
	pubKey, _ := x509.MarshalPKIXPublicKey(key.Public())

	data := biometricRaw{
		PersonHash:          bio.PersonHash(),
		EncryptedAESKey:     bio.CipherAESKey(),
		BIODPublicKey:       bio.BiodPublicKey(),
		EncryptedAddrRobot:  bio.CipherAddrRobot(),
		EncryptedAddrPerson: bio.CipherAddrPerson(),
		PersonPublicKey:     bio.PersonPublicKey(),
	}
	b, _ := json.Marshal(data)

	sig, _ := sign(hex.EncodeToString(pvKey), string(b))

	assert.Nil(t, NewSigner().VerifyTransactionDataSignature(mining.BiometricTransaction, hex.EncodeToString(pubKey), bio, sig))
}

/*
Scenario: Sign and check hash signature
	Given a hash and a key pair
	When I want to sign the hash and checks the signature generated
	Then I get not error
*/
func TestSignAndVerifyHashSignature(t *testing.T) {
	key, _ := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	pvKey, _ := x509.MarshalECPrivateKey(key)
	pubKey, _ := x509.MarshalPKIXPublicKey(key.Public())

	sig, err := NewSigner().SignHash("hash", hex.EncodeToString(pvKey))
	assert.Nil(t, err)
	assert.NotEmpty(t, sig)

	assert.Nil(t, NewSigner().VerifyHashSignature(hex.EncodeToString(pubKey), "hash", sig))
}

/*
Scenario: Sign and checks keychain validation request signature
	Given a validation request and a key pair
	When I want to sign the request and checks the signature generated
	Then I get not error
*/
func TestSignAndVerifyKeychainValidationRequestSignature(t *testing.T) {
	key, _ := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	pvKey, _ := x509.MarshalECPrivateKey(key)
	pubKey, _ := x509.MarshalPKIXPublicKey(key.Public())

	req := &api.KeychainValidationRequest{
		Data: &api.KeychainData{
			BiodPubk:        "pub",
			CipherAddrRobot: "enc addr",
			CipherWallet:    "enc wallet",
			PersonPubk:      "pub",
			Signature: &api.Signature{
				Biod:   "sig",
				Person: "sig",
			},
		},
		TransactionHash: "txHash",
	}

	err := NewSigner().SignKeychainValidationRequestSignature(req, hex.EncodeToString(pvKey))
	assert.Nil(t, err)
	assert.NotEmpty(t, req.Signature)

	assert.Nil(t, NewSigner().VerifyKeychainValidationRequestSignature(hex.EncodeToString(pubKey), req))
}

/*
Scenario: Sign and checks biometric validation request signature
	Given a validation request and a key pair
	When I want to sign the request and checks the signature generated
	Then I get not error
*/
func TestSignAndVerifyBiometricValidationRequestSignature(t *testing.T) {
	key, _ := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	pvKey, _ := x509.MarshalECPrivateKey(key)
	pubKey, _ := x509.MarshalPKIXPublicKey(key.Public())

	req := &api.BiometricValidationRequest{
		Data: &api.BiometricData{
			BiodPubk:        "pub",
			CipherAddrRobot: "enc addr",
			CipherAddrBio:   "enc addr",
			CipherAESKey:    "enc aes key",
			PersonHash:      "hash",
			PersonPubk:      "pub",
			Signature: &api.Signature{
				Biod:   "sig",
				Person: "sig",
			},
		},
		TransactionHash: "txHash",
	}

	err := NewSigner().SignBiometricValidationRequestSignature(req, hex.EncodeToString(pvKey))
	assert.Nil(t, err)
	assert.NotEmpty(t, req.Signature)

	assert.Nil(t, NewSigner().VerifyBiometricValidationRequestSignature(hex.EncodeToString(pubKey), req))
}

/*
Scenario: Sign and checks keychain storage request signature
	Given a storage request and a key pair
	When I want to sign the request and checks the signature generated
	Then I get not error
*/
func TestSignAndVerifyKeychainStorageRequestSignature(t *testing.T) {
	key, _ := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	pvKey, _ := x509.MarshalECPrivateKey(key)
	pubKey, _ := x509.MarshalPKIXPublicKey(key.Public())

	req := &api.KeychainStorageRequest{
		Data: &api.KeychainData{
			BiodPubk:        "pub",
			CipherAddrRobot: "enc addr",
			CipherWallet:    "enc wallet",
			PersonPubk:      "pub",
			Signature: &api.Signature{
				Biod:   "sig",
				Person: "sig",
			},
		},
		Endorsement: &api.Endorsement{
			LastTransactionHash: "",
			TransactionHash:     "hash",
			MasterValidation: &api.MasterValidation{
				LastTransactionMiners: []string{"hash"},
				ProofOfWorkRobotKey:   "key",
				ProofOfWorkValidation: &api.Validation{
					PublicKey: "key",
					Signature: "sig",
					Timestamp: time.Now().Unix(),
					Status:    api.Validation_OK,
				},
			},
		},
	}

	err := NewSigner().SignKeychainStorageRequestSignature(req, hex.EncodeToString(pvKey))
	assert.Nil(t, err)
	assert.NotEmpty(t, req.Signature)

	assert.Nil(t, NewSigner().VerifyKeychainStorageRequestSignature(hex.EncodeToString(pubKey), req))
}

/*
Scenario: Sign and checks keychain storage request signature
	Given a storage request and a key pair
	When I want to sign the request and checks the signature generated
	Then I get not error
*/
func TestSignAndVerifyBiometricStorageRequestSignature(t *testing.T) {
	key, _ := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	pvKey, _ := x509.MarshalECPrivateKey(key)
	pubKey, _ := x509.MarshalPKIXPublicKey(key.Public())

	req := &api.BiometricStorageRequest{
		Data: &api.BiometricData{
			BiodPubk:        "pub",
			CipherAddrRobot: "enc addr",
			CipherAddrBio:   "enc addr",
			CipherAESKey:    "enc aes key",
			PersonHash:      "hash",
			PersonPubk:      "pub",
			Signature: &api.Signature{
				Biod:   "sig",
				Person: "sig",
			},
		},
		Endorsement: &api.Endorsement{
			LastTransactionHash: "",
			TransactionHash:     "hash",
			MasterValidation: &api.MasterValidation{
				LastTransactionMiners: []string{"hash"},
				ProofOfWorkRobotKey:   "key",
				ProofOfWorkValidation: &api.Validation{
					PublicKey: "key",
					Signature: "sig",
					Timestamp: time.Now().Unix(),
					Status:    api.Validation_OK,
				},
			},
		},
	}

	err := NewSigner().SignBiometricStorageRequestSignature(req, hex.EncodeToString(pvKey))
	assert.Nil(t, err)
	assert.NotEmpty(t, req.Signature)

	assert.Nil(t, NewSigner().VerifyBiometricStorageRequestSignature(hex.EncodeToString(pubKey), req))
}

/*
Scenario: Sign and checks lock request signature
	Given a lock request and a key pair
	When I want to sign the request and checks the signature generated
	Then I get not error
*/
func TestSignVerifyLockRequestSignature(t *testing.T) {
	key, _ := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	pvKey, _ := x509.MarshalECPrivateKey(key)
	pubKey, _ := x509.MarshalPKIXPublicKey(key.Public())

	req := &api.LockRequest{
		Address:         "address",
		MasterRobotKey:  "robotkey",
		TransactionHash: "hash",
	}

	err := NewSigner().SignLockRequest(req, hex.EncodeToString(pvKey))
	assert.Nil(t, err)
	assert.NotEmpty(t, req.Signature)

	assert.Nil(t, NewSigner().VerifyLockRequestSignature(hex.EncodeToString(pubKey), req))
}

/*
Scenario: Sign and checks keychain lead mining request signature
	Given a keychain lead mining request and a key pair
	When I want to sign the request and checks the signature generated
	Then I get not error
*/
func TestSignAndVerifyKeychainLeadRequestSignature(t *testing.T) {
	key, _ := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	pvKey, _ := x509.MarshalECPrivateKey(key)
	pubKey, _ := x509.MarshalPKIXPublicKey(key.Public())

	req := &api.KeychainLeadRequest{
		EncryptedKeychainData: "enc data",
		SignatureKeychainData: &api.Signature{
			Biod:   "sig",
			Person: "sig",
		},
		TransactionHash:  "txHash",
		ValidatorPeerIPs: []string{"127.0.0.1"},
	}

	err := NewSigner().SignKeychainLeadRequest(req, hex.EncodeToString(pvKey))
	assert.Nil(t, err)
	assert.NotEmpty(t, req.SignatureRequest)

	assert.Nil(t, NewSigner().VerifyKeychainLeadRequestSignature(hex.EncodeToString(pubKey), req))
}

/*
Scenario: Sign and checks biometric lead mining request signature
	Given a biometric lead mining request and a key pair
	When I want to sign the request and checks the signature generated
	Then I get not error
*/
func TestSignAndVerifyBiometricLeadRequestSignature(t *testing.T) {
	key, _ := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	pvKey, _ := x509.MarshalECPrivateKey(key)
	pubKey, _ := x509.MarshalPKIXPublicKey(key.Public())

	req := &api.BiometricLeadRequest{
		EncryptedBioData: "enc data",
		SignatureBioData: &api.Signature{
			Biod:   "sig",
			Person: "sig",
		},
		TransactionHash:  "txHash",
		ValidatorPeerIPs: []string{"127.0.0.1"},
	}

	err := NewSigner().SignBiometricLeadRequest(req, hex.EncodeToString(pvKey))
	assert.Nil(t, err)
	assert.NotEmpty(t, req.SignatureRequest)

	assert.Nil(t, NewSigner().VerifyBiometricLeadRequestSignature(hex.EncodeToString(pubKey), req))
}

/*
Scenario: Sign and checks validation response signature
	Given a validation response and a key pair
	When I want to sign the response and checks the signature generated
	Then I get not error
*/
func TestSignAndVerifyValidationResponseSignature(t *testing.T) {
	key, _ := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	pvKey, _ := x509.MarshalECPrivateKey(key)
	pubKey, _ := x509.MarshalPKIXPublicKey(key.Public())

	res := &api.ValidationResponse{
		Validation: &api.Validation{
			PublicKey: "pub",
			Status:    api.Validation_OK,
			Timestamp: time.Now().Unix(),
			Signature: "sig",
		},
	}

	err := NewSigner().SignValidationResponse(res, hex.EncodeToString(pvKey))
	assert.Nil(t, err)
	assert.NotEmpty(t, res.Signature)

	assert.Nil(t, NewSigner().VerifyValidationResponseSignature(hex.EncodeToString(pubKey), res))
}

/*
Scenario: Sign and checks lock ack signature
	Given a lock ack and a key pair
	When I want to sign the response and checks the signature generated
	Then I get not error
*/
func TestSignAndChecLockAckSignature(t *testing.T) {
	key, _ := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	pvKey, _ := x509.MarshalECPrivateKey(key)
	pubKey, _ := x509.MarshalPKIXPublicKey(key.Public())

	ack := &api.LockAck{
		LockHash: "hash",
	}

	err := NewSigner().SignLockAck(ack, hex.EncodeToString(pvKey))
	assert.Nil(t, err)
	assert.NotEmpty(t, ack.Signature)

	assert.Nil(t, NewSigner().VerifyLockAckSignature(hex.EncodeToString(pubKey), ack))
}

/*
Scenario: Sign and checks storage ack signature
	Given a storage ack and a key pair
	When I want to sign the response and checks the signature generated
	Then I get not error
*/
func TestSignAndVerifyStorageAckSignature(t *testing.T) {
	key, _ := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	pvKey, _ := x509.MarshalECPrivateKey(key)
	pubKey, _ := x509.MarshalPKIXPublicKey(key.Public())

	ack := &api.StorageAck{
		StorageHash: "hash",
	}

	err := NewSigner().SignStorageAck(ack, hex.EncodeToString(pvKey))
	assert.Nil(t, err)
	assert.NotEmpty(t, ack.Signature)

	assert.Nil(t, NewSigner().VerifyStorageAckSignature(hex.EncodeToString(pubKey), ack))
}

/*
Scenario: Sign and checks account search result signature
	Given a account search result and a key pair
	When I want to sign the response and checks the signature generated
	Then I get not error
*/
func TestSignAndVerifyAccountSearchResultSignature(t *testing.T) {
	key, _ := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	pvKey, _ := x509.MarshalECPrivateKey(key)
	pubKey, _ := x509.MarshalPKIXPublicKey(key.Public())

	res := &api.AccountSearchResult{
		EncryptedAddress: "enc address",
		EncryptedAESkey:  "enc aes key",
		EncryptedWallet:  "enc wallet",
	}
	b, _ := json.Marshal(res)

	assert.Nil(t, NewSigner().SignAccountSearchResult(res, hex.EncodeToString(pvKey)))
	assert.NotEmpty(t, res.Signature)

	checkSignature(hex.EncodeToString(pubKey), string(b), res.Signature)
}

/*
Scenario: Sign and checks creation result signature
	Given a creation result and a key pair
	When I want to sign the response and checks the signature generated
	Then I get not error
*/
func TestSignAndVerifyCreationResultSignature(t *testing.T) {
	key, _ := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	pvKey, _ := x509.MarshalECPrivateKey(key)
	pubKey, _ := x509.MarshalPKIXPublicKey(key.Public())

	res := &api.CreationResult{
		MasterPeerIP:    "127.0.0.1",
		TransactionHash: "hash",
	}
	b, _ := json.Marshal(res)

	assert.Nil(t, NewSigner().SignCreationResult(res, hex.EncodeToString(pvKey)))
	assert.NotEmpty(t, res.Signature)

	checkSignature(hex.EncodeToString(pubKey), string(b), res.Signature)
}

/*
Scenario: Sign and checks biometric response signature
	Given a biometric response and a key pair
	When I want to sign the response and checks the signature generated
	Then I get not error
*/
func TestSignAndVerifyBiometricResponseSignature(t *testing.T) {
	key, _ := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	pvKey, _ := x509.MarshalECPrivateKey(key)
	pubKey, _ := x509.MarshalPKIXPublicKey(key.Public())

	res := &api.BiometricResponse{
		Data: &api.BiometricData{
			BiodPubk:        "pub",
			CipherAddrRobot: "enc addr",
			CipherAddrBio:   "enc addr",
			CipherAESKey:    "enc aes key",
			PersonHash:      "hash",
			PersonPubk:      "pub",
			Signature: &api.Signature{
				Biod:   "sig",
				Person: "sig",
			},
		},
		Endorsement: &api.Endorsement{
			LastTransactionHash: "",
			TransactionHash:     "hash",
			MasterValidation: &api.MasterValidation{
				LastTransactionMiners: []string{"hash"},
				ProofOfWorkRobotKey:   "key",
				ProofOfWorkValidation: &api.Validation{
					PublicKey: "key",
					Signature: "sig",
					Timestamp: time.Now().Unix(),
					Status:    api.Validation_OK,
				},
			},
		},
	}
	assert.Nil(t, NewSigner().SignBiometricResponse(res, hex.EncodeToString(pvKey)))
	assert.NotEmpty(t, res.Signature)

	assert.Nil(t, NewSigner().VerifyBiometricResponseSignature(hex.EncodeToString(pubKey), res))
}

/*
Scenario: Sign and checks keychain response signature
	Given a keychain response and a key pair
	When I want to sign the response and checks the signature generated
	Then I get not error
*/
func TestSignAndVerifyKeychainResponseSignature(t *testing.T) {
	key, _ := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	pvKey, _ := x509.MarshalECPrivateKey(key)
	pubKey, _ := x509.MarshalPKIXPublicKey(key.Public())

	res := &api.KeychainResponse{
		Data: &api.KeychainData{
			BiodPubk:        "pub",
			CipherAddrRobot: "enc addr",
			CipherWallet:    "enc wallet",
			PersonPubk:      "pub",
			Signature: &api.Signature{
				Biod:   "sig",
				Person: "sig",
			},
		},
		Endorsement: &api.Endorsement{
			LastTransactionHash: "",
			TransactionHash:     "hash",
			MasterValidation: &api.MasterValidation{
				LastTransactionMiners: []string{"hash"},
				ProofOfWorkRobotKey:   "key",
				ProofOfWorkValidation: &api.Validation{
					PublicKey: "key",
					Signature: "sig",
					Timestamp: time.Now().Unix(),
					Status:    api.Validation_OK,
				},
			},
		},
	}
	assert.Nil(t, NewSigner().SignKeychainResponse(res, hex.EncodeToString(pvKey)))
	assert.NotEmpty(t, res.Signature)

	assert.Nil(t, NewSigner().VerifyKeychainResponseSignature(hex.EncodeToString(pubKey), res))
}
