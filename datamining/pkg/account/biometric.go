package account

import (
	"github.com/uniris/uniris-core/datamining/pkg/mining"
)

//BiometricHasher defines methods to handle hash of biometric data
type BiometricHasher interface {
	HashBiometric(Biometric) (string, error)
	HashBiometricData(BiometricData) (string, error)
}

//BiometricSigner define methods to handle signatures of biometric data
type BiometricSigner interface {

	//VerifyBiometricDataSignature checks the signatures
	VerifyBiometricDataSignatures(BiometricData) error
}

//BiometricData describe a biometric data
type BiometricData interface {

	//PersonHash returns the person hash
	PersonHash() string

	//CipherAddrRobot returns the account's address encrypted with shared robot publickey
	CipherAddrRobot() string

	//CipherAddrPerson returns the account's address encrypted with person public key
	CipherAddrPerson() string

	//CipherAESKey returns the AES key encrypted with person public key
	CipherAESKey() string

	//BiodPublicKey return the biometric device public key
	BiodPublicKey() string

	//PersonPublicKey returns person public key
	PersonPublicKey() string

	//Signatures returns the signatures of the biometric data
	Signatures() Signatures
}

type biodata struct {
	personHash       string
	cipherAddrRobot  string
	cipherAddrPerson string
	cipherAESKey     string
	personPubk       string
	biodPubk         string
	sigs             Signatures
}

//NewBiometricData create new biometric
func NewBiometricData(personHash, cipherAddrRobot, cipherAddrPerson, cipherAesKey, personPubk, biodPubk string, sigs Signatures) BiometricData {
	return biodata{personHash, cipherAddrRobot, cipherAddrPerson, cipherAesKey, personPubk, biodPubk, sigs}
}

func (b biodata) PersonHash() string {
	return b.personHash
}

func (b biodata) CipherAddrRobot() string {
	return b.cipherAddrRobot
}

func (b biodata) CipherAddrPerson() string {
	return b.cipherAddrPerson
}

func (b biodata) CipherAESKey() string {
	return b.cipherAESKey
}

func (b biodata) PersonPublicKey() string {
	return b.personPubk
}

func (b biodata) BiodPublicKey() string {
	return b.biodPubk
}

func (b biodata) Signatures() Signatures {
	return b.sigs
}

//Biometric aggregates biometric data and its endorsement
type Biometric interface {
	BiometricData

	//Endorsement return the biometric data endorsement
	Endorsement() mining.Endorsement
}

type biometric struct {
	data        BiometricData
	endorsement mining.Endorsement
}

//NewBiometric creates a new biometric
func NewBiometric(data BiometricData, endor mining.Endorsement) Biometric {
	return biometric{data, endor}
}

func (b biometric) PersonHash() string {
	return b.data.PersonHash()
}

func (b biometric) CipherAddrRobot() string {
	return b.data.CipherAddrRobot()
}

func (b biometric) CipherAddrPerson() string {
	return b.data.CipherAddrPerson()
}

func (b biometric) CipherAESKey() string {
	return b.data.CipherAESKey()
}

func (b biometric) PersonPublicKey() string {
	return b.data.PersonPublicKey()
}

func (b biometric) BiodPublicKey() string {
	return b.data.BiodPublicKey()
}

func (b biometric) Signatures() Signatures {
	return b.data.Signatures()
}

func (b biometric) Endorsement() mining.Endorsement {
	return b.endorsement
}
