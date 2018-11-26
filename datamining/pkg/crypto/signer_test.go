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

	"github.com/uniris/uniris-core/datamining/pkg"

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
	prop := datamining.NewProposal(datamining.NewProposedKeyPair("enc key", "pub"))
	k := account.NewKeychain("enc addr", "enc wallet", "id pub", prop, "id sig", "em sig")
	key, _ := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	pvKey, _ := x509.MarshalECPrivateKey(key)
	pubKey, _ := x509.MarshalPKIXPublicKey(key.Public())

	data := keychainWithoutSig{
		EncryptedAddrByRobot: k.EncryptedAddrByRobot(),
		EncryptedWallet:      k.EncryptedWallet(),
		IDPublicKey:          k.IDPublicKey(),
		Proposal: proposal{
			SharedEmitterKeyPair: proposalKeypair{
				EncryptedPrivateKey: prop.SharedEmitterKeyPair().EncryptedPrivateKey(),
				PublicKey:           prop.SharedEmitterKeyPair().PublicKey(),
			},
		},
	}
	b, _ := json.Marshal(data)

	sig, _ := sign(hex.EncodeToString(pvKey), string(b))

	assert.Nil(t, NewSigner().VerifyTransactionDataSignature(mining.KeychainTransaction, hex.EncodeToString(pubKey), k, sig))
}

/*
Scenario: Verify ID transaction signature
	Given a signed ID data, a signature and a public key
	When I want to check it's the data matches the signature
	Then I get not error
*/
func TestVerifyTransactionIDSignature(t *testing.T) {
	prop := datamining.NewProposal(datamining.NewProposedKeyPair("enc key", "pub"))

	id := account.NewID("hash", "enc addr", "enc addr", "enc aes key", "id pub", prop, "id sig", "em sig")
	key, _ := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	pvKey, _ := x509.MarshalECPrivateKey(key)
	pubKey, _ := x509.MarshalPKIXPublicKey(key.Public())

	data := idWithoutSig{
		Hash:                 id.Hash(),
		EncryptedAESKey:      id.EncryptedAESKey(),
		EncryptedAddrByRobot: id.EncryptedAddrByRobot(),
		EncryptedAddrByID:    id.EncryptedAddrByID(),
		PublicKey:            id.PublicKey(),
		Proposal: proposal{
			SharedEmitterKeyPair: proposalKeypair{
				EncryptedPrivateKey: prop.SharedEmitterKeyPair().EncryptedPrivateKey(),
				PublicKey:           prop.SharedEmitterKeyPair().PublicKey(),
			},
		},
	}
	b, _ := json.Marshal(data)

	sig, _ := sign(hex.EncodeToString(pvKey), string(b))

	assert.Nil(t, NewSigner().VerifyTransactionDataSignature(mining.IDTransaction, hex.EncodeToString(pubKey), id, sig))
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
		Data: &api.Keychain{
			EncryptedAddrByRobot: "enc addr",
			EncryptedWallet:      "enc wallet",
			IDPublicKey:          "pub",
			EmitterSignature:     "sig",
			IDSignature:          "sig",
			Proposal: &api.Proposal{
				SharedEmitterKeyPair: &api.KeyPairProposal{
					EncryptedPrivateKey: "enc pv key",
					PublicKey:           "pub key",
				},
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
Scenario: Sign and checks ID validation request signature
	Given a validation request and a key pair
	When I want to sign the request and checks the signature generated
	Then I get not error
*/
func TestSignAndVerifyIDValidationRequestSignature(t *testing.T) {
	key, _ := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	pvKey, _ := x509.MarshalECPrivateKey(key)
	pubKey, _ := x509.MarshalPKIXPublicKey(key.Public())

	req := &api.IDValidationRequest{
		Data: &api.ID{
			EncryptedAddrByRobot: "enc addr",
			EncryptedAddrByID:    "enc addr",
			EncryptedAESKey:      "enc aes key",
			Hash:                 "hash",
			PublicKey:            "pub",
			EmitterSignature:     "sig",
			IDSignature:          "sig",
			Proposal: &api.Proposal{
				SharedEmitterKeyPair: &api.KeyPairProposal{
					EncryptedPrivateKey: "enc prv key",
					PublicKey:           "pub key",
				},
			},
		},
		TransactionHash: "txHash",
	}

	err := NewSigner().SignIDValidationRequestSignature(req, hex.EncodeToString(pvKey))
	assert.Nil(t, err)
	assert.NotEmpty(t, req.Signature)

	assert.Nil(t, NewSigner().VerifyIDValidationRequestSignature(hex.EncodeToString(pubKey), req))
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
		Data: &api.Keychain{
			EncryptedAddrByRobot: "enc addr",
			EncryptedWallet:      "enc wallet",
			IDPublicKey:          "pub",
			EmitterSignature:     "sig",
			IDSignature:          "sig",
			Proposal: &api.Proposal{
				SharedEmitterKeyPair: &api.KeyPairProposal{
					EncryptedPrivateKey: "enc prv key",
					PublicKey:           "pub key",
				},
			},
		},
		Endorsement: &api.Endorsement{
			LastTransactionHash: "",
			TransactionHash:     "hash",
			MasterValidation: &api.MasterValidation{
				LastTransactionMiners: []string{"hash"},
				ProofOfWorkKey:        "key",
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
Scenario: Sign and checks ID storage request signature
	Given a storage request and a key pair
	When I want to sign the request and checks the signature generated
	Then I get not error
*/
func TestSignAndVerifyIDStorageRequestSignature(t *testing.T) {
	key, _ := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	pvKey, _ := x509.MarshalECPrivateKey(key)
	pubKey, _ := x509.MarshalPKIXPublicKey(key.Public())

	req := &api.IDStorageRequest{
		Data: &api.ID{
			EncryptedAddrByRobot: "enc addr",
			EncryptedAddrByID:    "enc addr",
			EncryptedAESKey:      "enc aes key",
			Hash:                 "hash",
			PublicKey:            "pub",
			EmitterSignature:     "sig",
			IDSignature:          "sig",
			Proposal: &api.Proposal{
				SharedEmitterKeyPair: &api.KeyPairProposal{
					EncryptedPrivateKey: "enc prv key",
					PublicKey:           "pub key",
				},
			},
		},
		Endorsement: &api.Endorsement{
			LastTransactionHash: "",
			TransactionHash:     "hash",
			MasterValidation: &api.MasterValidation{
				LastTransactionMiners: []string{"hash"},
				ProofOfWorkKey:        "key",
				ProofOfWorkValidation: &api.Validation{
					PublicKey: "key",
					Signature: "sig",
					Timestamp: time.Now().Unix(),
					Status:    api.Validation_OK,
				},
			},
		},
	}

	err := NewSigner().SignIDStorageRequestSignature(req, hex.EncodeToString(pvKey))
	assert.Nil(t, err)
	assert.NotEmpty(t, req.Signature)

	assert.Nil(t, NewSigner().VerifyIDStorageRequestSignature(hex.EncodeToString(pubKey), req))
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
		EncryptedKeychain: "enc data",
		TransactionHash:   "txHash",
		ValidatorPeerIPs:  []string{"127.0.0.1"},
	}

	err := NewSigner().SignKeychainLeadRequest(req, hex.EncodeToString(pvKey))
	assert.Nil(t, err)
	assert.NotEmpty(t, req.SignatureRequest)

	assert.Nil(t, NewSigner().VerifyKeychainLeadRequestSignature(hex.EncodeToString(pubKey), req))
}

/*
Scenario: Sign and checks ID lead mining request signature
	Given a ID lead mining request and a key pair
	When I want to sign the request and checks the signature generated
	Then I get not error
*/
func TestSignAndVerifyIDLeadRequestSignature(t *testing.T) {
	key, _ := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	pvKey, _ := x509.MarshalECPrivateKey(key)
	pubKey, _ := x509.MarshalPKIXPublicKey(key.Public())

	req := &api.IDLeadRequest{
		EncryptedID:      "enc data",
		TransactionHash:  "txHash",
		ValidatorPeerIPs: []string{"127.0.0.1"},
	}

	err := NewSigner().SignIDLeadRequest(req, hex.EncodeToString(pvKey))
	assert.Nil(t, err)
	assert.NotEmpty(t, req.SignatureRequest)

	assert.Nil(t, NewSigner().VerifyIDLeadRequestSignature(hex.EncodeToString(pubKey), req))
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
Scenario: Sign and checks ID response signature
	Given a ID response and a key pair
	When I want to sign the response and checks the signature generated
	Then I get not error
*/
func TestSignAndVerifyIDResponseSignature(t *testing.T) {
	key, _ := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	pvKey, _ := x509.MarshalECPrivateKey(key)
	pubKey, _ := x509.MarshalPKIXPublicKey(key.Public())

	res := &api.IDResponse{
		Data: &api.ID{
			EncryptedAddrByRobot: "enc addr",
			EncryptedAddrByID:    "enc addr",
			EncryptedAESKey:      "enc aes key",
			Hash:                 "hash",
			PublicKey:            "pub",
			EmitterSignature:     "sig",
			IDSignature:          "sig",
			Proposal: &api.Proposal{
				SharedEmitterKeyPair: &api.KeyPairProposal{
					EncryptedPrivateKey: "enc prv key",
					PublicKey:           "pub key",
				},
			},
		},
		Endorsement: &api.Endorsement{
			LastTransactionHash: "",
			TransactionHash:     "hash",
			MasterValidation: &api.MasterValidation{
				LastTransactionMiners: []string{"hash"},
				ProofOfWorkKey:        "key",
				ProofOfWorkValidation: &api.Validation{
					PublicKey: "key",
					Signature: "sig",
					Timestamp: time.Now().Unix(),
					Status:    api.Validation_OK,
				},
			},
		},
	}
	assert.Nil(t, NewSigner().SignIDResponse(res, hex.EncodeToString(pvKey)))
	assert.NotEmpty(t, res.Signature)

	assert.Nil(t, NewSigner().VerifyIDResponseSignature(hex.EncodeToString(pubKey), res))
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
		Data: &api.Keychain{
			EncryptedAddrByRobot: "enc addr",
			EncryptedWallet:      "enc wallet",
			IDPublicKey:          "pub",
			EmitterSignature:     "sig",
			IDSignature:          "sig",
			Proposal: &api.Proposal{
				SharedEmitterKeyPair: &api.KeyPairProposal{
					EncryptedPrivateKey: "enc prv key",
					PublicKey:           "pub key",
				},
			},
		},
		Endorsement: &api.Endorsement{
			LastTransactionHash: "",
			TransactionHash:     "hash",
			MasterValidation: &api.MasterValidation{
				LastTransactionMiners: []string{"hash"},
				ProofOfWorkKey:        "key",
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

/*
Scenario: Sign and checks validation signature
	Given a validation and a key pair
	When I want to sign the validation and checks the signature generated
	Then I get not error
*/
func TestSignAndVerifyValidationSignature(t *testing.T) {
	key, _ := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	pvKey, _ := x509.MarshalECPrivateKey(key)
	pubKey, _ := x509.MarshalPKIXPublicKey(key.Public())

	v := mining.NewValidation(mining.ValidationOK, time.Now(), hex.EncodeToString(pubKey), "")
	sValid, err := NewSigner().SignValidation(v, hex.EncodeToString(pvKey))
	assert.Nil(t, err)
	assert.NotEmpty(t, sValid.Signature())

	assert.Nil(t, NewSigner().VerifyValidationSignature(sValid))
}
