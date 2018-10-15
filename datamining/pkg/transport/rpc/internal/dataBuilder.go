package internalrpc

import (
	"encoding/hex"
	"encoding/json"

	api "github.com/uniris/uniris-core/datamining/api/protobuf-spec"
	"github.com/uniris/uniris-core/datamining/pkg"
)

//DataBuilder defines methods to transform API entities for the domain layer
type DataBuilder struct {
	decrypt Decrypter
}

//BuildWalletResult constitue details for the WalletResult rpc command
func (b DataBuilder) BuildWalletResult(w datamining.Wallet, bioW datamining.BioWallet) *api.WalletResult {
	return &api.WalletResult{
		EncryptedAESkey:        bioW.CipherAESKey(),
		EncryptedWallet:        w.CipherWallet(),
		EncryptedWalletAddress: w.WalletAddr(),
	}
}

//BuildWallet constitue details for the StoreWallet rpc command
func (b DataBuilder) BuildWallet(p *api.Wallet) (w datamining.WalletData, bw datamining.BioData, err error) {
	bioData, err := b.decrypt.Decipher(p.EncryptedBioData)
	if err != nil {
		return
	}
	walletData, err := b.decrypt.Decipher(p.EncryptedWalletData)
	if err != nil {
		return
	}

	bioDataB, err := hex.DecodeString(string(bioData))
	if err != nil {
		return
	}

	var bio BioDataFromJSON
	err = json.Unmarshal(bioDataB, &bio)
	if err != nil {
		return
	}

	walletDataB, err := hex.DecodeString(string(walletData))
	if err != nil {
		return
	}

	var wallet WalletDataFromJSON
	err = json.Unmarshal(walletDataB, &wallet)
	if err != nil {
		return
	}

	bioSig := Signatures{
		Person: string(p.SignatureBioData.Person),
		Biod:   string(p.SignatureBioData.Biod),
	}

	walletSig := Signatures{
		Person: string(p.SignatureWalletData.Person),
		Biod:   string(p.SignatureWalletData.Biod),
	}

	w = datamining.WalletData{
		BiodPubk:        datamining.PublicKey(wallet.BiodPublicKey),
		CipherAddrRobot: datamining.WalletAddr(wallet.EncryptedAddrRobot),
		CipherWallet:    datamining.CipherWallet(wallet.EncryptedWallet),
		EmPubk:          datamining.PublicKey(wallet.PersonPublicKey),
		Sigs: datamining.Signatures{
			BiodSig: []byte(walletSig.Biod),
			EmSig:   []byte(walletSig.Person),
		},
	}

	bw = datamining.BioData{
		BHash:           datamining.BioHash(bio.PersonHash),
		BiodPubk:        datamining.PublicKey(bio.BiodPublicKey),
		CipherAddrBio:   datamining.WalletAddr(bio.EncryptedAddrPerson),
		CipherAddrRobot: datamining.WalletAddr(bio.EncryptedAddrRobot),
		CipherAESKey:    []byte(bio.EncryptedAESKey),
		EmPubk:          datamining.PublicKey(bio.PersonPublicKey),
		Sigs: datamining.Signatures{
			BiodSig: []byte(bioSig.Biod),
			EmSig:   []byte(bioSig.Person),
		},
	}

	return
}

//WalletDataFromJSON represents wallet data JSON
type WalletDataFromJSON struct {
	PersonPublicKey    string     `json:"person_pubk"`
	BiodPublicKey      string     `json:"biod_pubk"`
	EncryptedWallet    string     `json:"encrypted_wallet"`
	EncryptedAddrRobot string     `json:"encrypted_addr_robot"`
	Sigs               Signatures `json:"signature_wallet"`
}

//BioDataFromJSON represents bio data JSON
type BioDataFromJSON struct {
	PersonPublicKey     string     `json:"person_pubk"`
	BiodPublicKey       string     `json:"biod_pubk"`
	PersonHash          string     `json:"person_hash"`
	EncryptedAESKey     string     `json:"encrypted_aes_key"`
	EncryptedAddrPerson string     `json:"encrypted_addr_person"`
	EncryptedAddrRobot  string     `json:"encrypted_addr_robot"`
	Sigs                Signatures `json:"signature_bio"`
}

//Signatures represents signatures JSON
type Signatures struct {
	Person string `json:"person_sig"`
	Biod   string `json:"biod_sig"`
}
