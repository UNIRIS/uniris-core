package formater

import (
	robot "github.com/uniris/uniris-core/datamining/pkg"
)

//FormatedWallet describe a decode wallet request
type FormatedWallet struct {
	WalletAddr      []byte
	CipherAddrRobot []byte             `json:"@wallet_encrypted_robot"`
	CWallet         robot.CipherWallet `json:"encrypted_wallet"`
	EmPubk          robot.PublicKey    `json:"person_pubk"`
	BiodPubk        robot.PublicKey    `json:"biod_pubk"`
	Sigs            Signatures         `json:"signature_wallet"`
}

//Signatures describe differnet needed signatures
type Signatures struct {
	EmSig   robot.Signature `json:"person_sig"`
	BiodSig robot.Signature `json:"biod_sig"`
}

//FormatedBioWallet describe a decode wallet request
type FormatedBioWallet struct {
	BHash           robot.BioHash
	CipherAddrRobot []byte          `json:"@wallet_encrypted_person"`
	CipherAddrBio   []byte          `json:"@wallet_encrypted_robot"`
	CipherAesKey    []byte          `json:"encrypted_aes_key"`
	EmPubk          robot.PublicKey `json:"person_pubk"`
	BiodPubk        robot.PublicKey `json:"biod_pubk"`
	Sigs            Signatures      `json:"signature_bio"`
}
