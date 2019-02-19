package rest

import (
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"

	api "github.com/uniris/uniris-core/api/protobuf-spec"
	"github.com/uniris/uniris-core/pkg/crypto"
)

func encodeTxReceipt(tx *api.Transaction) string {
	return fmt.Sprintf("%s%s", tx.Address, tx.TransactionHash)
}

func decodeTxReceipt(receipt string) (addr, hash string, err error) {
	if _, err = hex.DecodeString(receipt); err != nil {
		err = errors.New("must be hexadecimal")
		return
	}

	/*
		Length from sha256 hash is 64 bytes.
		a transaction receipt is a set of the hash of the address and the hash of the transaction
		So a transaction receipt is 128 bytes
	*/
	if len(receipt) != 128 {
		err = errors.New("invalid length")
		return
	}

	addr = receipt[:64]
	hash = receipt[64:]

	if _, err = crypto.IsHash(addr); err != nil {
		return
	}

	if _, err = crypto.IsHash(hash); err != nil {
		return
	}

	return
}

func decodeTransactionRaw(txEncoded string, pvKey string) (*api.Transaction, error) {
	txJSON, err := crypto.Decrypt(txEncoded, pvKey)
	if err != nil {
		return nil, err
	}
	txHash := crypto.HashString(txJSON)

	var tx struct {
		Address                   string            `json:"addr"`
		Data                      map[string]string `json:"data"`
		Timestamp                 int64             `json:"timestamp"`
		Type                      int               `json:"type"`
		PublicKey                 string            `json:"public_key"`
		SharedKeysEmitterProposal struct {
			EncryptedPrivateKey string `json:"encrypted_private_key"`
			PublicKey           string `json:"public_key"`
		} `json:"em_shared_keys_proposal"`
		Signature        string `json:"signature"`
		EmitterSignature string `json:"em_signature"`
	}
	if err := json.Unmarshal([]byte(txJSON), &tx); err != nil {
		return nil, err
	}

	return &api.Transaction{
		Address:          tx.Address,
		Data:             tx.Data,
		Type:             api.TransactionType(tx.Type),
		Timestamp:        tx.Timestamp,
		PublicKey:        tx.PublicKey,
		Signature:        tx.Signature,
		EmitterSignature: tx.EmitterSignature,
		SharedKeysEmitterProposal: &api.SharedKeyPair{
			EncryptedPrivateKey: tx.SharedKeysEmitterProposal.EncryptedPrivateKey,
			PublicKey:           tx.SharedKeysEmitterProposal.PublicKey,
		},
		TransactionHash: txHash,
	}, nil

}
