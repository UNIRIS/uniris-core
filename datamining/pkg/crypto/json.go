package crypto

import (
	"github.com/uniris/uniris-core/datamining/pkg/mining"
)

type endorsement struct {
	LastTxHash       string           `json:"last_transaction_hash"`
	TxHash           string           `json:"transaction_hash"`
	MasterValidation masterValidation `json:"master_validation"`
	Validations      []validation     `json:"validations"`
}

type masterValidation struct {
	LastTxRvk []string   `json:"last_transaction_miners"`
	PowKey    string     `json:"pow_key"`
	PowValid  validation `json:"pow_validation"`
}

type validation struct {
	Status    mining.ValidationStatus `json:"status"`
	Timestamp int64                   `json:"timestamp"`
	Pubk      string                  `json:"public_key"`
	Sig       string                  `json:"signature"`
}

type endorsedID struct {
	PublicKey            string      `json:"pubk"`
	Hash                 string      `json:"hash"`
	EncryptedAESKey      string      `json:"encrypted_aes_key"`
	EncryptedAddrByID    string      `json:"encrypted_addr_id"`
	EncryptedAddrByRobot string      `json:"encrypted_addr_robot"`
	Proposal             proposal    `json:"proposal"`
	IDSignature          string      `json:"id_sig"`
	EmitterSignature     string      `json:"em_sig"`
	Endorsement          endorsement `json:"endorsement"`
}

type idWithoutSig struct {
	PublicKey            string   `json:"pubk"`
	Hash                 string   `json:"hash"`
	EncryptedAESKey      string   `json:"encrypted_aes_key"`
	EncryptedAddrByID    string   `json:"encrypted_addr_id"`
	EncryptedAddrByRobot string   `json:"encrypted_addr_robot"`
	Proposal             proposal `json:"proposal"`
}

type id struct {
	PublicKey            string   `json:"pubk"`
	Hash                 string   `json:"hash"`
	EncryptedAESKey      string   `json:"encrypted_aes_key"`
	EncryptedAddrByID    string   `json:"encrypted_addr_id"`
	EncryptedAddrByRobot string   `json:"encrypted_addr_robot"`
	Proposal             proposal `json:"proposal"`
	IDSignature          string   `json:"id_sig"`
	EmitterSignature     string   `json:"em_sig"`
}

type endorsedKeychain struct {
	Address              string      `json:"address"`
	IDPublicKey          string      `json:"id_pubk"`
	EncryptedWallet      string      `json:"encrypted_wal"`
	EncryptedAddrByRobot string      `json:"encrypted_addr_robot"`
	Proposal             proposal    `json:"proposal"`
	IDSignature          string      `json:"id_sig"`
	EmitterSignature     string      `json:"em_sig"`
	Endorsement          endorsement `json:"endorsement"`
}

type keychain struct {
	IDPublicKey          string   `json:"id_pubk"`
	EncryptedWallet      string   `json:"encrypted_wal"`
	EncryptedAddrByRobot string   `json:"encrypted_addr_robot"`
	Proposal             proposal `json:"proposal"`
	EmitterSignature     string   `json:"em_sig"`
	IDSignature          string   `json:"id_sig"`
}

type keychainWithoutSig struct {
	IDPublicKey          string   `json:"id_pubk"`
	EncryptedWallet      string   `json:"encrypted_wal"`
	EncryptedAddrByRobot string   `json:"encrypted_addr_robot"`
	Proposal             proposal `json:"proposal"`
}

type lockRaw struct {
	TxHash         string `json:"transaction_hash"`
	MasterRobotKey string `json:"master_robot_key"`
	Address        string `json:"address"`
}

type validationWithoutSig struct {
	Status    mining.ValidationStatus `json:"status"`
	Timestamp int64                   `json:"timestamp"`
	PublicKey string                  `json:"public_key"`
}

type proposal struct {
	SharedEmitterKeyPair proposalKeypair `json:"shared_emitter_kp"`
}

type proposalKeypair struct {
	EncryptedPrivateKey string `json:"encrypted_private_key"`
	PublicKey           string `json:"public_key"`
}

type transactionResult struct {
	TransactionHash string `json:"transaction_hash"`
	MasterPeerIP    string `json:"master_peer_ip"`
}

type accountSearchResult struct {
	EncryptedAESKey  string `json:"encrypted_aes_key"`
	EncryptedWallet  string `json:"encrypted_wallet"`
	EncryptedAddress string `json:"encrypted_address"`
}

type contractWithoutSig struct {
	Address   string `json:"address"`
	Code      string `json:"code"`
	Event     string `json:"event"`
	PublicKey string `json":public_key"`
}

type contractJSON struct {
	Address          string `json:"address"`
	Code             string `json:"code"`
	Event            string `json:"event"`
	PublicKey        string `json:"public_key"`
	Signature        string `json:"signature"`
	EmitterSignature string `json:"emitter_signature"`
}

type endorsedContractJSON struct {
	Address          string      `json:"address"`
	Code             string      `json:"code"`
	Event            string      `json:"event"`
	PublicKey        string      `json:"public_key"`
	Signature        string      `json:"signature"`
	EmitterSignature string      `json:"emitter_signature"`
	Endorsement      endorsement `json:"endorsement"`
}

type contractMessage struct {
	ContractAddress  string   `json:"contract_address"`
	Method           string   `json:"method"`
	Parameters       []string `json:"parameters"`
	PublicKey        string   `json:"public_key"`
	Signature        string   `json:"signature"`
	EmitterSignature string   `json:"emitter_signature"`
}

type contractMessageWithoutAddress struct {
	Method           string   `json:"method"`
	Parameters       []string `json:"parameters"`
	PublicKey        string   `json:"public_key"`
	Signature        string   `json:"signature"`
	EmitterSignature string   `json:"emitter_signature"`
}

type contractMessageWithSig struct {
	Method     string   `json:"method"`
	Parameters []string `json:"parameters"`
	PublicKey  string   `json:"public_key"`
}

type endorsedContractMessage struct {
	ContractAddress  string      `json:"contract_address"`
	Method           string      `json:"method"`
	Parameters       []string    `json:"parameters"`
	PublicKey        string      `json:"public_key"`
	Signature        string      `json:"signature"`
	EmitterSignature string      `json:"emitter_signature"`
	Endorsement      endorsement `json:"endorsement"`
}
