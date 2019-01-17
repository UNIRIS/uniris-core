package crypto

import (
	"crypto/x509"
	"encoding/hex"
	"encoding/json"

	"github.com/uniris/uniris-core/datamining/pkg/contract"

	datamining "github.com/uniris/uniris-core/datamining/pkg"

	"github.com/uniris/uniris-core/datamining/pkg/account"
	"github.com/uniris/uniris-core/datamining/pkg/transport/rpc"

	ecies "github.com/uniris/ecies/pkg"
)

//Decrypter defines methods to handle decryption
type Decrypter interface {
	rpc.Decrypter
}

type decrypter struct{}

//NewDecrypter create a new decrypter
func NewDecrypter() Decrypter {
	return decrypter{}
}

func (d decrypter) DecryptHash(hash string, pvKey string) (string, error) {
	return decrypt(pvKey, hash)
}

func (d decrypter) DecryptID(data string, pvKey string) (account.ID, error) {
	clear, err := decrypt(pvKey, data)
	if err != nil {
		return nil, err
	}
	var id id
	err = json.Unmarshal([]byte(clear), &id)
	if err != nil {
		return nil, err
	}

	prop := datamining.NewProposal(
		datamining.NewProposedKeyPair(
			id.Proposal.SharedEmitterKeyPair.EncryptedPrivateKey,
			id.Proposal.SharedEmitterKeyPair.PublicKey),
	)

	return account.NewID(
		id.Hash,
		id.EncryptedAddrByRobot,
		id.EncryptedAddrByID,
		id.EncryptedAESKey,
		id.PublicKey,
		prop,
		id.IDSignature,
		id.EmitterSignature), nil
}

func (d decrypter) DecryptKeychain(data string, pvKey string) (account.Keychain, error) {
	clear, err := decrypt(pvKey, data)
	if err != nil {
		return nil, err
	}
	var kc keychain
	err = json.Unmarshal([]byte(clear), &kc)
	if err != nil {
		return nil, err
	}

	prop := datamining.NewProposal(
		datamining.NewProposedKeyPair(
			kc.Proposal.SharedEmitterKeyPair.EncryptedPrivateKey,
			kc.Proposal.SharedEmitterKeyPair.PublicKey),
	)

	return account.NewKeychain(
		kc.EncryptedAddrByRobot,
		kc.EncryptedWallet,
		kc.IDPublicKey,
		prop,
		kc.IDSignature,
		kc.EmitterSignature), nil
}

func (d decrypter) DecryptContract(data string, pvKey string) (contract.Contract, error) {
	clear, err := decrypt(pvKey, data)
	if err != nil {
		return nil, err
	}
	var c contractJSON
	err = json.Unmarshal([]byte(clear), &c)
	if err != nil {
		return nil, err
	}

	return contract.New(c.Address, c.Code, c.Event, c.PublicKey, c.Signature, c.EmitterSignature), nil
}

func (d decrypter) DecryptContractMessage(data string, pvKey string) (contract.Message, error) {
	clear, err := decrypt(pvKey, data)
	if err != nil {
		return nil, err
	}
	var c contractMessageWithoutAddress
	err = json.Unmarshal([]byte(clear), &c)
	if err != nil {
		return nil, err
	}

	return contract.NewMessage("", c.Method, c.Parameters, c.PublicKey, c.Signature, c.EmitterSignature), nil
}

func decrypt(privk string, data string) (string, error) {
	decodeKey, err := hex.DecodeString(privk)
	if err != nil {
		return "", err
	}

	robotKey, err := x509.ParseECPrivateKey(decodeKey)
	if err != nil {
		return "", err
	}

	decodeCipher, err := hex.DecodeString(data)
	if err != nil {
		return "", err
	}

	robotEciesKey := ecies.ImportECDSA(robotKey)
	message, err := robotEciesKey.Decrypt(decodeCipher, nil, nil)
	if err != nil {
		return "", err
	}
	return string(message), nil
}
