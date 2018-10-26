package datamining

import (
	"encoding/json"
)

//Signatures describe differnet needed signatures
type Signatures struct {
	EmSig   string
	BiodSig string
}

//WalletData describe a decoded wallet data from the wallet request
type WalletData struct {
	WalletAddr      string
	CipherAddrRobot string
	CipherWallet    string
	EmPubk          string
	BiodPubk        string
	Sigs            Signatures
}

//BioData describe a decoded biometric data from the wallet request
type BioData struct {
	BHash           string
	CipherAddrRobot string
	CipherAddrBio   string
	CipherAESKey    string
	EmPubk          string
	BiodPubk        string
	Sigs            Signatures
}

//Wallet describes the stored wallet with its endorsement
type Wallet struct {
	data        *WalletData
	oldTxnHash  string
	endorsement *Endorsement
}

//NewWallet creates a new wallet
func NewWallet(data *WalletData, endor *Endorsement, oldTxnHash string) *Wallet {
	return &Wallet{data, oldTxnHash, endor}
}

//Endorsement returns the wallet endorsement
func (w Wallet) Endorsement() *Endorsement {
	return w.endorsement
}

//CipherWallet returns the encrypted wallet
func (w Wallet) CipherWallet() string {
	return w.data.CipherWallet
}

//WalletAddr returns address of the encrypted wallet
func (w Wallet) WalletAddr() string {
	return w.data.WalletAddr
}

//CipherAddrRobot get the wallet address encrypted for roboto
func (w Wallet) CipherAddrRobot() string {
	return w.data.CipherAddrRobot
}

//PersonPublicKey returns the wallet personal public key
func (w Wallet) PersonPublicKey() string {
	return w.data.EmPubk
}

//BiodPublicKey returns the wallet biod public key
func (w Wallet) BiodPublicKey() string {
	return w.data.BiodPubk
}

//Signatures return the wallet signatures
func (w Wallet) Signatures() Signatures {
	return w.data.Sigs
}

//OldTransactionHash returns the hash of the previous transaction
func (w Wallet) OldTransactionHash() string {
	return w.oldTxnHash
}

func (w Wallet) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		Data                *WalletData  `json:"data"`
		Endorsement         *Endorsement `json:"endorsment"`
		LastTransactionHash string       `json:"last_transaction_hash"`
	}{
		Data:                w.data,
		Endorsement:         w.endorsement,
		LastTransactionHash: w.oldTxnHash,
	})
}

func (w *Wallet) UnmarshalJSON(b []byte) error {

	bData := struct {
		Data                *WalletData  `json:"data"`
		Endorsement         *Endorsement `json:"endorsment"`
		LastTransactionHash string       `json:"last_transaction_hash"`
	}{}

	if err := json.Unmarshal(b, &bData); err != nil {
		return err
	}

	w.data = bData.Data
	w.endorsement = bData.Endorsement
	w.oldTxnHash = bData.LastTransactionHash

	return nil
}

//BioWallet describes the stored biometric wallet with its endorsement
type BioWallet struct {
	data        *BioData
	endorsement *Endorsement
}

//NewBioWallet creates a new bio wallet
func NewBioWallet(data *BioData, endor *Endorsement) *BioWallet {
	return &BioWallet{data, endor}
}

//BiodPublicKey return the biometric public key for the bio wallet
func (b BioWallet) BiodPublicKey() string {
	return b.data.BiodPubk
}

//PersonPublicKey returns person public key for the bio wallet
func (b BioWallet) PersonPublicKey() string {
	return b.data.EmPubk
}

//Signatures returns the bio wallet signatures
func (b BioWallet) Signatures() Signatures {
	return b.data.Sigs
}

//Bhash returns the biometric hash
func (b BioWallet) Bhash() string {
	return b.data.BHash
}

//CipherAddrRobot returns the address of the wallet encrypted with shared robot publickey
func (b BioWallet) CipherAddrRobot() string {
	return b.data.CipherAddrRobot
}

//CipherAddrBio returns the address of the wallet encrypted with person keys
func (b BioWallet) CipherAddrBio() string {
	return b.data.CipherAddrBio
}

//CipherAESKey returns the AES key encrypted with person keys
func (b BioWallet) CipherAESKey() string {
	return b.data.CipherAESKey
}

//Endorsement returns the bio wallet endorsement
func (b BioWallet) Endorsement() *Endorsement {
	return b.endorsement
}

func (b BioWallet) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		Data        *BioData     `json:"data"`
		Endorsement *Endorsement `json:"endorsment"`
	}{
		Data:        b.data,
		Endorsement: b.endorsement,
	})
}

func (b *BioWallet) UnmarshalJSON(bytes []byte) error {
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
