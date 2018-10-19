package datamining

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
	data       *WalletData
	oldTxnHash string
	endorsment *Endorsement
}

//NewWallet creates a new wallet
func NewWallet(data *WalletData, endor *Endorsement, oldTxnHash string) *Wallet {
	return &Wallet{data, oldTxnHash, endor}
}

//Endorsement returns the wallet endorsement
func (w Wallet) Endorsement() *Endorsement {
	return w.endorsment
}

//CipherWallet returns the encrypted wallet
func (w Wallet) CipherWallet() string {
	return w.data.CipherWallet
}

//WalletAddr returns address of the encrypted wallet
func (w Wallet) WalletAddr() string {
	return w.data.WalletAddr
}

//BioWallet describes the stored biometric wallet with its endorsement
type BioWallet struct {
	data       *BioData
	endorsment *Endorsement
}

//NewBioWallet creates a new bio wallet
func NewBioWallet(data *BioData, endor *Endorsement) *BioWallet {
	return &BioWallet{data, endor}
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
