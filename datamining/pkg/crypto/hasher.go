package crypto

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"

	"github.com/uniris/uniris-core/datamining/pkg/lock"

	"github.com/uniris/uniris-core/datamining/pkg/account"
)

//Hasher defines methods for hashing
type Hasher interface {
	account.KeychainHasher
	account.BiometricHasher
	lock.Hasher
}

type hasher struct{}

//NewHasher creates new hasher
func NewHasher() Hasher {
	return hasher{}
}

func (h hasher) HashLock(txLock lock.TransactionLock) (string, error) {
	b, err := json.Marshal(lockRaw{
		Address:        txLock.Address,
		MasterRobotKey: txLock.MasterRobotKey,
		TxHash:         txLock.TxHash,
	})
	if err != nil {
		return "", err
	}
	return hashBytes(b), nil
}

func (h hasher) HashBiometric(data account.Biometric) (string, error) {
	valids := make([]validation, 0)
	for _, v := range data.Endorsement().Validations() {
		valids = append(valids, validation{
			Pubk:      v.PublicKey(),
			Sig:       v.Signature(),
			Status:    v.Status(),
			Timestamp: v.Timestamp(),
		})
	}
	b, err := json.Marshal(biometric{
		PersonHash:          data.PersonHash(),
		EncryptedAddrPerson: data.CipherAddrPerson(),
		EncryptedAddrRobot:  data.CipherAddrRobot(),
		EncryptedAESKey:     data.CipherAESKey(),
		BIODPublicKey:       data.BiodPublicKey(),
		PersonPublicKey:     data.PersonPublicKey(),
		BIODSignature:       data.Signatures().Biod(),
		PersonSignature:     data.Signatures().Person(),
		Endorsement: endorsement{
			LastTxHash: data.Endorsement().LastTransactionHash(),
			TxHash:     data.Endorsement().TransactionHash(),
			MasterValidation: masterValidation{
				LastTxRvk:   data.Endorsement().MasterValidation().LastTransactionMiners(),
				PowRobotKey: data.Endorsement().MasterValidation().ProofOfWorkRobotKey(),
				PowValid: validation{
					Pubk:      data.Endorsement().MasterValidation().ProofOfWorkValidation().PublicKey(),
					Sig:       data.Endorsement().MasterValidation().ProofOfWorkValidation().Signature(),
					Status:    data.Endorsement().MasterValidation().ProofOfWorkValidation().Status(),
					Timestamp: data.Endorsement().MasterValidation().ProofOfWorkValidation().Timestamp(),
				},
			},
			Validations: valids,
		},
	})
	if err != nil {
		return "", err
	}
	return hashBytes(b), nil
}

func (h hasher) HashKeychain(data account.Keychain) (string, error) {
	valids := make([]validation, 0)
	for _, v := range data.Endorsement().Validations() {
		valids = append(valids, validation{
			Pubk:      v.PublicKey(),
			Sig:       v.Signature(),
			Status:    v.Status(),
			Timestamp: v.Timestamp(),
		})
	}
	b, err := json.Marshal(keychain{
		Address:            data.Address(),
		EncryptedWallet:    data.CipherWallet(),
		EncryptedAddrRobot: data.CipherAddrRobot(),
		BIODPublicKey:      data.BiodPublicKey(),
		PersonPublicKey:    data.PersonPublicKey(),
		BIODSignature:      data.Signatures().Biod(),
		PersonSignature:    data.Signatures().Person(),
		Endorsement: endorsement{
			LastTxHash: data.Endorsement().LastTransactionHash(),
			TxHash:     data.Endorsement().TransactionHash(),
			MasterValidation: masterValidation{
				LastTxRvk:   data.Endorsement().MasterValidation().LastTransactionMiners(),
				PowRobotKey: data.Endorsement().MasterValidation().ProofOfWorkRobotKey(),
				PowValid: validation{
					Pubk:      data.Endorsement().MasterValidation().ProofOfWorkValidation().PublicKey(),
					Sig:       data.Endorsement().MasterValidation().ProofOfWorkValidation().Signature(),
					Status:    data.Endorsement().MasterValidation().ProofOfWorkValidation().Status(),
					Timestamp: data.Endorsement().MasterValidation().ProofOfWorkValidation().Timestamp(),
				},
			},
			Validations: valids,
		},
	})
	if err != nil {
		return "", err
	}
	return hashBytes(b), nil
}

func (h hasher) HashBiometricData(data account.BiometricData) (string, error) {
	b, err := json.Marshal(biometricData{
		PersonHash:          data.PersonHash(),
		EncryptedAddrPerson: data.CipherAddrPerson(),
		EncryptedAddrRobot:  data.CipherAddrRobot(),
		EncryptedAESKey:     data.CipherAESKey(),
		BIODPublicKey:       data.BiodPublicKey(),
		PersonPublicKey:     data.PersonPublicKey(),
		BIODSignature:       data.Signatures().Biod(),
		PersonSignature:     data.Signatures().Person(),
	})
	if err != nil {
		return "", err
	}
	return hashBytes(b), nil
}

func (h hasher) HashKeychainData(data account.KeychainData) (string, error) {
	b, err := json.Marshal(keychainData{
		EncryptedAddrRobot: data.CipherAddrRobot(),
		EncryptedWallet:    data.CipherWallet(),
		BIODPublicKey:      data.BiodPublicKey(),
		PersonPublicKey:    data.PersonPublicKey(),
		BIODSignature:      data.Signatures().Biod(),
		PersonSignature:    data.Signatures().Person(),
	})
	if err != nil {
		return "", err
	}
	return hashBytes(b), nil
}

func hashString(data string) string {
	return hashBytes([]byte(data))
}

func hashBytes(data []byte) string {
	hash := sha256.New()
	hash.Write([]byte(data))
	return hex.EncodeToString(hash.Sum(nil))
}
