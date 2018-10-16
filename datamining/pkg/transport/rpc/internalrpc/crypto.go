package internalrpc

import (
	"encoding/json"

	api "github.com/uniris/uniris-core/datamining/api/protobuf-spec"
	crypto "github.com/uniris/uniris-core/datamining/pkg/crypto"
)

//DecryptWallet decrypt bio and wallet data from the robot shared private key
func DecryptWallet(w *api.Wallet, pvKey []byte) ([]byte, []byte, error) {

	bioData, err := crypto.Decrypt(pvKey, w.EncryptedBioData)
	if err != nil {
		return nil, nil, err
	}
	walletData, err := crypto.Decrypt(pvKey, w.EncryptedWalletData)
	if err != nil {
		return nil, nil, err
	}

	return bioData, walletData, nil
}

//VerifyBioSignatures checks if the bio data is valid from the given signatures
func VerifyBioSignatures(bio BioDataFromJSON, sigBio *api.Signature) error {
	bioRaw, err := json.Marshal(bio)
	if err != nil {
		return err
	}

	if err := crypto.Verify([]byte(bio.PersonPublicKey), sigBio.Person, crypto.Hash(bioRaw)); err != nil {
		return err
	}

	if err := crypto.Verify([]byte(bio.BiodPublicKey), sigBio.Biod, crypto.Hash(bioRaw)); err != nil {
		return err
	}

	return nil
}

//VerifyWalSignatures checks if the wallet data is valid from the given signatures
func VerifyWalSignatures(wal WalletDataFromJSON, sigWal *api.Signature) error {
	walRaw, err := json.Marshal(wal)
	if err != nil {
		return err
	}

	if err := crypto.Verify([]byte(wal.PersonPublicKey), sigWal.Person, crypto.Hash(walRaw)); err != nil {
		return err
	}

	if err := crypto.Verify([]byte(wal.BiodPublicKey), sigWal.Biod, crypto.Hash(walRaw)); err != nil {
		return err
	}

	return nil
}
