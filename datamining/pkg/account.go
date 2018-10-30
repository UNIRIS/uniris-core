package datamining

import (
	"encoding/json"
)

//Signatures describe differnet needed signatures
type Signatures struct {
	BiodSig   string
	PersonSig string
}

//KeyChainData describe a decoded keychain data
type KeyChainData struct {
	WalletAddr      string
	CipherAddrRobot string
	CipherWallet    string
	PersonPubk      string
	BiodPubk        string
	Sigs            Signatures
}

//BioData describe a decoded biometric data
type BioData struct {
	PersonHash      string
	CipherAddrRobot string
	CipherAddrBio   string
	CipherAESKey    string
	PersonPubk      string
	BiodPubk        string
	Sigs            Signatures
}

//Keychain describes the keychain data with its endorsement
type Keychain struct {
	data        *KeyChainData
	oldTxnHash  string
	endorsement *Endorsement
}

//NewKeychain creates a new keychain
func NewKeychain(data *KeyChainData, endor *Endorsement, oldTxnHash string) *Keychain {
	return &Keychain{data, oldTxnHash, endor}
}

//Endorsement returns the wallet endorsement
func (k Keychain) Endorsement() *Endorsement {
	return k.endorsement
}

//CipherWallet returns the encrypted wallet
func (k Keychain) CipherWallet() string {
	return k.data.CipherWallet
}

//WalletAddr returns address of the encrypted wallet
func (k Keychain) WalletAddr() string {
	return k.data.WalletAddr
}

//CipherAddrRobot get the wallet address encrypted for roboto
func (k Keychain) CipherAddrRobot() string {
	return k.data.CipherAddrRobot
}

//PersonPublicKey returns the wallet personal public key
func (k Keychain) PersonPublicKey() string {
	return k.data.PersonPubk
}

//BiodPublicKey returns the wallet biod public key
func (k Keychain) BiodPublicKey() string {
	return k.data.BiodPubk
}

//Signatures return the wallet signatures
func (k Keychain) Signatures() Signatures {
	return k.data.Sigs
}

//OldTransactionHash returns the hash of the previous transaction
func (k Keychain) OldTransactionHash() string {
	return k.oldTxnHash
}

func (k Keychain) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		Data                *KeyChainData `json:"data"`
		Endorsement         *Endorsement  `json:"endorsment"`
		LastTransactionHash string        `json:"last_transaction_hash"`
	}{
		Data:                k.data,
		Endorsement:         k.endorsement,
		LastTransactionHash: k.oldTxnHash,
	})
}

func (k *Keychain) UnmarshalJSON(b []byte) error {

	bData := struct {
		Data                *KeyChainData `json:"data"`
		Endorsement         *Endorsement  `json:"endorsment"`
		LastTransactionHash string        `json:"last_transaction_hash"`
	}{}

	if err := json.Unmarshal(b, &bData); err != nil {
		return err
	}

	k.data = bData.Data
	k.endorsement = bData.Endorsement
	k.oldTxnHash = bData.LastTransactionHash

	return nil
}

//Biometric describes the biometric data with its endorsement
type Biometric struct {
	data        *BioData
	endorsement *Endorsement
}

//NewBiometric creates a new biometric
func NewBiometric(data *BioData, endor *Endorsement) *Biometric {
	return &Biometric{data, endor}
}

//BiodPublicKey return the biometric public key for the bio wallet
func (b Biometric) BiodPublicKey() string {
	return b.data.BiodPubk
}

//PersonPublicKey returns person public key for the bio wallet
func (b Biometric) PersonPublicKey() string {
	return b.data.PersonPubk
}

//Signatures returns the bio wallet signatures
func (b Biometric) Signatures() Signatures {
	return b.data.Sigs
}

//PersonHash returns the person hash
func (b Biometric) PersonHash() string {
	return b.data.PersonHash
}

//CipherAddrRobot returns the address of the wallet encrypted with shared robot publickey
func (b Biometric) CipherAddrRobot() string {
	return b.data.CipherAddrRobot
}

//CipherAddrBio returns the address of the wallet encrypted with person keys
func (b Biometric) CipherAddrBio() string {
	return b.data.CipherAddrBio
}

//CipherAESKey returns the AES key encrypted with person keys
func (b Biometric) CipherAESKey() string {
	return b.data.CipherAESKey
}

//Endorsement returns the bio wallet endorsement
func (b Biometric) Endorsement() *Endorsement {
	return b.endorsement
}

func (b Biometric) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		Data        *BioData     `json:"data"`
		Endorsement *Endorsement `json:"endorsment"`
	}{
		Data:        b.data,
		Endorsement: b.endorsement,
	})
}

func (b *Biometric) UnmarshalJSON(bytes []byte) error {
	bData := struct {
		Data        *BioData     `json:"data"`
		Endorsement *Endorsement `json:"endorsment"`
	}{}
	if err := json.Unmarshal(bytes, &bData); err != nil {
		return err
	}

	b.data = bData.Data
	b.endorsement = bData.Endorsement
	return nil
}
