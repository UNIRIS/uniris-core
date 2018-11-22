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
	account.IDHasher
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

func (h hasher) HashEndorsedID(id account.EndorsedID) (string, error) {
	valids := make([]validation, 0)
	for _, v := range id.Endorsement().Validations() {
		valids = append(valids, validation{
			Pubk:      v.PublicKey(),
			Sig:       v.Signature(),
			Status:    v.Status(),
			Timestamp: v.Timestamp().Unix(),
		})
	}
	b, err := json.Marshal(endorsedID{
		Hash:                 id.Hash(),
		EncryptedAddrByID:    id.EncryptedAddrByID(),
		EncryptedAddrByRobot: id.EncryptedAddrByRobot(),
		EncryptedAESKey:      id.EncryptedAESKey(),
		PublicKey:            id.PublicKey(),
		EmitterSignature:     id.EmitterSignature(),
		IDSignature:          id.IDSignature(),
		Endorsement: endorsement{
			LastTxHash: id.Endorsement().LastTransactionHash(),
			TxHash:     id.Endorsement().TransactionHash(),
			MasterValidation: masterValidation{
				LastTxRvk: id.Endorsement().MasterValidation().LastTransactionMiners(),
				PowKey:    id.Endorsement().MasterValidation().ProofOfWorkKey(),
				PowValid: validation{
					Pubk:      id.Endorsement().MasterValidation().ProofOfWorkValidation().PublicKey(),
					Sig:       id.Endorsement().MasterValidation().ProofOfWorkValidation().Signature(),
					Status:    id.Endorsement().MasterValidation().ProofOfWorkValidation().Status(),
					Timestamp: id.Endorsement().MasterValidation().ProofOfWorkValidation().Timestamp().Unix(),
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

func (h hasher) HashEndorsedKeychain(kc account.EndorsedKeychain) (string, error) {
	valids := make([]validation, 0)
	for _, v := range kc.Endorsement().Validations() {
		valids = append(valids, validation{
			Pubk:      v.PublicKey(),
			Sig:       v.Signature(),
			Status:    v.Status(),
			Timestamp: v.Timestamp().Unix(),
		})
	}
	b, err := json.Marshal(endorsedKeychain{
		Address:              kc.Address(),
		EncryptedWallet:      kc.EncryptedWallet(),
		EncryptedAddrByRobot: kc.EncryptedAddrByRobot(),
		IDPublicKey:          kc.IDPublicKey(),
		EmitterSignature:     kc.EmitterSignature(),
		IDSignature:          kc.IDSignature(),
		Endorsement: endorsement{
			LastTxHash: kc.Endorsement().LastTransactionHash(),
			TxHash:     kc.Endorsement().TransactionHash(),
			MasterValidation: masterValidation{
				LastTxRvk: kc.Endorsement().MasterValidation().LastTransactionMiners(),
				PowKey:    kc.Endorsement().MasterValidation().ProofOfWorkKey(),
				PowValid: validation{
					Pubk:      kc.Endorsement().MasterValidation().ProofOfWorkValidation().PublicKey(),
					Sig:       kc.Endorsement().MasterValidation().ProofOfWorkValidation().Signature(),
					Status:    kc.Endorsement().MasterValidation().ProofOfWorkValidation().Status(),
					Timestamp: kc.Endorsement().MasterValidation().ProofOfWorkValidation().Timestamp().Unix(),
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

func (h hasher) HashID(data account.ID) (string, error) {
	b, err := json.Marshal(id{
		Hash:                 data.Hash(),
		EncryptedAddrByID:    data.EncryptedAddrByID(),
		EncryptedAddrByRobot: data.EncryptedAddrByRobot(),
		EncryptedAESKey:      data.EncryptedAESKey(),
		PublicKey:            data.PublicKey(),
		EmitterSignature:     data.EmitterSignature(),
		IDSignature:          data.IDSignature(),
	})
	if err != nil {
		return "", err
	}
	return hashBytes(b), nil
}

func (h hasher) HashKeychain(kc account.Keychain) (string, error) {
	b, err := json.Marshal(keychain{
		EncryptedAddrByRobot: kc.EncryptedAddrByRobot(),
		EncryptedWallet:      kc.EncryptedWallet(),
		IDPublicKey:          kc.IDPublicKey(),
		EmitterSignature:     kc.EmitterSignature(),
		IDSignature:          kc.IDSignature(),
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
