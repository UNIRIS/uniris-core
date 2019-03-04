package rest

import (
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"

	validator "gopkg.in/go-playground/validator.v9"

	api "github.com/uniris/uniris-core/api/protobuf-spec"
	"github.com/uniris/uniris-core/pkg/crypto"
)

func encodeTxReceipt(tx *api.Transaction) string {
	return fmt.Sprintf("%x%x", tx.Address, tx.TransactionHash)
}

func decodeTxReceipt(receipt string) (addr, txHash crypto.VersionnedHash, err error) {
	rBytes, err := hex.DecodeString(receipt)
	if err != nil {
		err = fmt.Errorf("transaction receipt is not in hexadecimal")
		return nil, nil, err
	}

	defer func() {
		if r := recover(); r != nil {
			err = errors.New("transaction receipt is invalid")
			return
		}
	}()

	h := crypto.VersionnedHash(rBytes)
	hSize := h.Algorithm().Size()

	//Including the versionning algo byte
	addr = h[:hSize+1]

	//Each hash is versionned so
	//
	// Hash1: [Algo byte][Digest bytes]
	// Hash2: [Algo byte][Digest bytes]
	//Exemple: of SHA256, 32 bytes for the digest. For 2 hashes: 64 (digest) + 2 (algo)
	if len(rBytes) != ((hSize * 2) + 2) {
		err = errors.New("transaction receipt is invalid")
		return
	}

	//Including the versionning algo byte
	txHash = h[hSize+1:]

	if !addr.IsValid() {
		return nil, nil, errors.New("transaction address is an invalid hash")
	}

	if !txHash.IsValid() {
		return nil, nil, errors.New("transaction hash is an invalid hash")
	}

	return
}

func decodeTransactionRaw(txEncoded []byte, pvKey crypto.PrivateKey) (*api.Transaction, *httpError) {
	txJSON, err := pvKey.Decrypt(txEncoded)
	if err != nil {
		return nil, &httpError{
			code:      http.StatusBadRequest,
			Timestamp: time.Now().Unix(),
			Status:    http.StatusText(http.StatusBadRequest),
			Error:     err.Error(),
		}
	}
	if !json.Valid(txJSON) {
		return nil, &httpError{
			code:      http.StatusBadRequest,
			Timestamp: time.Now().Unix(),
			Status:    http.StatusText(http.StatusBadRequest),
			Error:     "invalid JSON",
		}
	}

	var tx txRaw
	if err := json.Unmarshal(txJSON, &tx); err != nil {
		return nil, &httpError{
			code:      http.StatusBadRequest,
			Timestamp: time.Now().Unix(),
			Status:    http.StatusText(http.StatusBadRequest),
			Error:     err.Error(),
		}
	}

	if err := validator.New().Struct(tx); err != nil {
		return nil, &httpError{
			code:      http.StatusBadRequest,
			Timestamp: time.Now().Unix(),
			Status:    http.StatusText(http.StatusBadRequest),
			Error:     err.Error(),
		}
	}

	txJSONForHash, err := json.Marshal(map[string]interface{}{
		"addr":       tx.Address,
		"data":       tx.Data,
		"timestamp":  tx.Timestamp,
		"type":       tx.Type,
		"public_key": tx.PublicKey,
		"em_shared_keys_proposal": map[string]string{
			"encrypted_private_key": tx.SharedKeysEmitterProposal.EncryptedPrivateKey,
			"public_key":            tx.SharedKeysEmitterProposal.PublicKey,
		},
		"signature":    tx.Signature,
		"em_signature": tx.EmitterSignature,
	})
	if err != nil {
		return nil, &httpError{
			code:      http.StatusInternalServerError,
			Timestamp: time.Now().Unix(),
			Status:    http.StatusText(http.StatusBadRequest),
			Error:     err.Error(),
		}
	}

	addr, _ := hex.DecodeString(tx.Address)
	pubK, _ := hex.DecodeString(tx.PublicKey)
	sig, _ := hex.DecodeString(tx.Signature)
	emSig, _ := hex.DecodeString(tx.EmitterSignature)
	sharedPropPv, _ := hex.DecodeString(tx.SharedKeysEmitterProposal.EncryptedPrivateKey)
	sharedPropPub, _ := hex.DecodeString(tx.SharedKeysEmitterProposal.PublicKey)

	data := make(map[string][]byte)
	for k, v := range tx.Data {
		vB, _ := hex.DecodeString(v)
		data[k] = vB
	}

	return &api.Transaction{
		Address:          addr,
		Data:             data,
		Type:             api.TransactionType(tx.Type),
		Timestamp:        tx.Timestamp,
		PublicKey:        pubK,
		Signature:        sig,
		EmitterSignature: emSig,
		SharedKeysEmitterProposal: &api.SharedKeyPair{
			EncryptedPrivateKey: sharedPropPv,
			PublicKey:           sharedPropPub,
		},
		TransactionHash: crypto.Hash(txJSONForHash),
	}, nil

}
