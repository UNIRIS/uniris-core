package account

import (
	"encoding/json"

	datamining "github.com/uniris/uniris-core/datamining/pkg"
)

//BioData describe a decoded biometric data
type BioData struct {
	PersonHash      string
	CipherAddrRobot string
	CipherAddrBio   string
	CipherAESKey    string
	PersonPubk      string
	BiodPubk        string
	Sigs            datamining.Signatures
}

type biometric struct {
	data        *BioData
	endorsement datamining.Endorsement
}

//Biometric represents a biometric
type Biometric interface {
	BiodPublicKey() string
	PersonPublicKey() string
	PersonHash() string
	CipherAddrRobot() string
	CipherAddrBio() string
	CipherAESKey() string
	Signatures() datamining.Signatures
	Endorsement() datamining.Endorsement
}

//NewBiometric creates a new biometric
func NewBiometric(data *BioData, endor datamining.Endorsement) Biometric {
	return biometric{data, endor}
}

//BiodPublicKey return the biometric public key for the bio wallet
func (b biometric) BiodPublicKey() string {
	return b.data.BiodPubk
}

//PersonPublicKey returns person public key for the bio wallet
func (b biometric) PersonPublicKey() string {
	return b.data.PersonPubk
}

//Signatures returns the bio wallet signatures
func (b biometric) Signatures() datamining.Signatures {
	return b.data.Sigs
}

//PersonHash returns the person hash
func (b biometric) PersonHash() string {
	return b.data.PersonHash
}

//CipherAddrRobot returns the address of the wallet encrypted with shared robot publickey
func (b biometric) CipherAddrRobot() string {
	return b.data.CipherAddrRobot
}

//CipherAddrBio returns the address of the wallet encrypted with person keys
func (b biometric) CipherAddrBio() string {
	return b.data.CipherAddrBio
}

//CipherAESKey returns the AES key encrypted with person keys
func (b biometric) CipherAESKey() string {
	return b.data.CipherAESKey
}

//Endorsement returns the bio wallet endorsement
func (b biometric) Endorsement() datamining.Endorsement {
	return b.endorsement
}

func (b biometric) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		Data        *BioData               `json:"data"`
		Endorsement datamining.Endorsement `json:"endorsment"`
	}{
		Data:        b.data,
		Endorsement: b.endorsement,
	})
}

func (b *biometric) UnmarshalJSON(bytes []byte) error {
	bData := struct {
		Data        *BioData               `json:"data"`
		Endorsement datamining.Endorsement `json:"endorsment"`
	}{}
	if err := json.Unmarshal(bytes, &bData); err != nil {
		return err
	}

	b.data = bData.Data
	b.endorsement = bData.Endorsement
	return nil
}
