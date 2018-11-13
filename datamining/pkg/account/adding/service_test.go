package adding

import (
	"errors"
	"sort"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/uniris/uniris-core/datamining/pkg/account"
	"github.com/uniris/uniris-core/datamining/pkg/account/listing"
	"github.com/uniris/uniris-core/datamining/pkg/mining"
)

/*
Scenario: Store a keychain
	Given a data data
	When I want to store a keychain
	Then the wallet is stored on the database
*/
func TestStoreKeychain(t *testing.T) {
	repo := &databasemock{}

	lister := listing.NewService(repo)
	s := NewService(mockAiClient{}, repo, mockSigVerfier{}, lister, mockHasher{})

	sigs := account.NewSignatures("sig1", "sig2")

	end := mining.NewEndorsement(
		"", "hash", mining.NewMasterValidation([]string{}, "robotkey", mining.NewValidation(mining.ValidationOK, time.Now(), "robotkey", "sig")),
		[]mining.Validation{
			mining.NewValidation(mining.ValidationOK, time.Now(), "pub", "sig"),
		},
	)
	data := account.NewKeychainData("enc addr", "enc wallet", "person pub", "biod pub", sigs)
	kc := account.NewKeychain("addr", data, end)

	err := s.StoreKeychain(kc)
	assert.Nil(t, err)
	assert.Len(t, repo.keychains, 1)
	assert.Equal(t, "addr", repo.keychains[0].Address())
	assert.Equal(t, "hash", repo.keychains[0].Endorsement().TransactionHash())
}

/*
Scenario: Store a keychain with master validation KO
	Given a keychain with a master validation as KO
	When I want to store the keychain
	Then I get the keychain is store on the KO database
*/
func TestStoreKeychainWithMasterValidKO(t *testing.T) {
	repo := &databasemock{}

	lister := listing.NewService(repo)
	s := NewService(mockAiClient{}, repo, mockSigVerfier{}, lister, mockHasher{})

	sigs := account.NewSignatures("sig1", "sig2")

	end := mining.NewEndorsement(
		"", "hash", mining.NewMasterValidation([]string{}, "robotkey", mining.NewValidation(mining.ValidationKO, time.Now(), "robotkey", "sig")),
		[]mining.Validation{
			mining.NewValidation(mining.ValidationOK, time.Now(), "pub", "sig"),
		},
	)
	data := account.NewKeychainData("enc addr", "enc wallet", "person pub", "biod pub", sigs)
	kc := account.NewKeychain("addr", data, end)

	err := s.StoreKeychain(kc)
	assert.Nil(t, err)

	assert.Empty(t, repo.keychains)
	assert.Len(t, repo.keychainsKO, 1)
	assert.Equal(t, "addr", repo.keychainsKO[0].Address())
	assert.Equal(t, "hash", repo.keychainsKO[0].Endorsement().TransactionHash())
}

/*
Scenario: Store a keychain with one slave validation as KO
	Given a keychain with one slave validation as KO
	When I want to store the keychain
	Then I get the keychain is store on the KO database
*/
func TestStoreKeychainWithOneSlaveValidKO(t *testing.T) {
	repo := &databasemock{}

	lister := listing.NewService(repo)
	s := NewService(mockAiClient{}, repo, mockSigVerfier{}, lister, mockHasher{})
	sigs := account.NewSignatures("sig1", "sig2")

	end := mining.NewEndorsement(
		"", "hash", mining.NewMasterValidation([]string{}, "robotkey", mining.NewValidation(mining.ValidationOK, time.Now(), "robotkey", "sig")),
		[]mining.Validation{
			mining.NewValidation(mining.ValidationOK, time.Now(), "pub", "sig"),
			mining.NewValidation(mining.ValidationKO, time.Now(), "pub", "sig"),
		},
	)
	data := account.NewKeychainData("enc addr", "enc wallet", "person pub", "biod pub", sigs)
	kc := account.NewKeychain("addr", data, end)

	err := s.StoreKeychain(kc)
	assert.Nil(t, err)

	assert.Empty(t, repo.keychains)
	assert.Len(t, repo.keychainsKO, 1)
	assert.Equal(t, "addr", repo.keychainsKO[0].Address())
	assert.Equal(t, "hash", repo.keychainsKO[0].Endorsement().TransactionHash())
}

/*
Scenario: Store a keychain with a invalid last transaction hash
	Given a previous keychain stored and a new keychain with a invalid last transaction hash
	When I want to store it
	Then I get the error
*/
func TestInvalidLastTransactionKeychain(t *testing.T) {
	repo := &databasemock{}
	lister := listing.NewService(repo)
	s := NewService(mockAiClient{}, repo, mockSigVerfier{}, lister, mockHasher{})

	sigs := account.NewSignatures("sig1", "sig2")
	end1 := mining.NewEndorsement(
		"", "hash", mining.NewMasterValidation([]string{}, "robotkey", mining.NewValidation(mining.ValidationOK, time.Now(), "robotkey", "sig")),
		[]mining.Validation{
			mining.NewValidation(mining.ValidationOK, time.Now(), "pub", "sig"),
		},
	)
	data := account.NewKeychainData("enc addr", "enc wallet", "person pub", "biod pub", sigs)
	kc1 := account.NewKeychain("addr", data, end1)

	assert.Nil(t, s.StoreKeychain(kc1))

	end2 := mining.NewEndorsement(
		"bad last hash", "hash", mining.NewMasterValidation([]string{}, "robotkey", mining.NewValidation(mining.ValidationOK, time.Now(), "robotkey", "sig")),
		[]mining.Validation{
			mining.NewValidation(mining.ValidationOK, time.Now(), "pub", "sig"),
		},
	)
	kc2 := account.NewKeychain("addr", data, end2)

	assert.Equal(t, ErrInvalidDataIntegrity, s.StoreKeychain(kc2))

	end3 := mining.NewEndorsement(
		"", "hash", mining.NewMasterValidation([]string{}, "robotkey", mining.NewValidation(mining.ValidationOK, time.Now(), "robotkey", "sig")),
		[]mining.Validation{
			mining.NewValidation(mining.ValidationOK, time.Now(), "pub", "sig"),
		},
	)

	kc3 := account.NewKeychain("addr", data, end3)

	assert.Equal(t, ErrInvalidDataIntegrity, s.StoreKeychain(kc3))
}

/*
Scenario: Store a keychain with a zero
	Given a keychain without validations
	When I want to store it
	Then I get the error
*/
func TestStoreKeychainWithZeroValidations(t *testing.T) {
	repo := &databasemock{}
	lister := listing.NewService(repo)
	s := NewService(mockAiClient{}, repo, mockSigVerfier{}, lister, mockHasher{})

	sigs := account.NewSignatures("sig1", "sig2")
	end := mining.NewEndorsement(
		"", "hash", mining.NewMasterValidation([]string{}, "robotkey", mining.NewValidation(mining.ValidationOK, time.Now(), "robotkey", "sig")),
		[]mining.Validation{},
	)
	data := account.NewKeychainData("enc addr", "enc wallet", "person pub", "biod pub", sigs)
	kc := account.NewKeychain("addr", data, end)

	assert.Equal(t, ErrInvalidValidationNumber, s.StoreKeychain(kc))

}

/*
Scenario: Store a keychain with a invalid transaction hash
	Given a keychain with invalid tx hash
	When I want to store it
	Then I get the error
*/
func TestStoreKeychainWithInvalidTxHash(t *testing.T) {
	repo := &databasemock{}
	lister := listing.NewService(repo)
	s := NewService(mockAiClient{}, repo, mockSigVerfier{}, lister, mockHasher{})

	sigs := account.NewSignatures("sig1", "sig2")
	end := mining.NewEndorsement(
		"", "bad hash", mining.NewMasterValidation([]string{}, "robotkey", mining.NewValidation(mining.ValidationOK, time.Now(), "robotkey", "sig")),
		[]mining.Validation{
			mining.NewValidation(mining.ValidationOK, time.Now(), "pub", "sig"),
		},
	)
	data := account.NewKeychainData("enc addr", "enc wallet", "person pub", "biod pub", sigs)
	kc := account.NewKeychain("addr", data, end)

	assert.Equal(t, ErrInvalidDataIntegrity, s.StoreKeychain(kc))
}

/*
Scenario: Store a biometric
	Given a bio data
	When I want to store a biometric data
	Then the bio data are stored on the database
*/
func TestStoreBiometric(t *testing.T) {
	repo := &databasemock{}
	lister := listing.NewService(repo)
	s := NewService(mockAiClient{}, repo, mockSigVerfier{}, lister, mockHasher{})

	sigs := account.NewSignatures("sig1", "sig2")

	end := mining.NewEndorsement(
		"", "hash", mining.NewMasterValidation([]string{}, "robotkey", mining.NewValidation(mining.ValidationOK, time.Now(), "robotkey", "sig")),
		[]mining.Validation{
			mining.NewValidation(mining.ValidationOK, time.Now(), "pub", "sig"),
		},
	)

	data := account.NewBiometricData("pHash", "enc addr robot", "enc addr person", "enc aes key", "person pub", "biod pub", sigs)
	bio := account.NewBiometric(data, end)
	err := s.StoreBiometric(bio)
	assert.Nil(t, err)
	l := len(repo.biometrics)
	assert.Equal(t, 1, l)
	assert.Equal(t, 1, l)
	assert.Equal(t, "pHash", repo.biometrics[0].PersonHash())
}

/*
Scenario: Store a biometric with a zero
	Given a biometric without validations
	When I want to store it
	Then I get the error
*/
func TestStoreBiometricWithZeroValidations(t *testing.T) {
	repo := &databasemock{}
	lister := listing.NewService(repo)
	s := NewService(mockAiClient{}, repo, mockSigVerfier{}, lister, mockHasher{})

	sigs := account.NewSignatures("sig1", "sig2")
	end := mining.NewEndorsement(
		"", "hash", mining.NewMasterValidation([]string{}, "robotkey", mining.NewValidation(mining.ValidationOK, time.Now(), "robotkey", "sig")),
		[]mining.Validation{},
	)
	data := account.NewBiometricData("pHash", "enc addr robot", "enc addr person", "enc aes key", "person pub", "biod pub", sigs)
	bio := account.NewBiometric(data, end)

	assert.Equal(t, ErrInvalidValidationNumber, s.StoreBiometric(bio))

}

/*
Scenario: Store a biometric with master validation KO
	Given a biometric with a master validation as KO
	When I want to store the biometric
	Then I get the biometric is store on the KO database
*/
func TestStoreBiometricWithMasterValidKO(t *testing.T) {
	repo := &databasemock{}

	lister := listing.NewService(repo)
	s := NewService(mockAiClient{}, repo, mockSigVerfier{}, lister, mockHasher{})

	sigs := account.NewSignatures("sig1", "sig2")

	end := mining.NewEndorsement(
		"", "hash", mining.NewMasterValidation([]string{}, "robotkey", mining.NewValidation(mining.ValidationKO, time.Now(), "robotkey", "sig")),
		[]mining.Validation{
			mining.NewValidation(mining.ValidationOK, time.Now(), "pub", "sig"),
		},
	)
	data := account.NewBiometricData("pHash", "enc addr robot", "enc addr person", "enc aes key", "person pub", "biod pub", sigs)
	bio := account.NewBiometric(data, end)

	err := s.StoreBiometric(bio)
	assert.Nil(t, err)

	assert.Empty(t, repo.biometrics)
	assert.Len(t, repo.biometricsKO, 1)
	assert.Equal(t, "enc aes key", repo.biometricsKO[0].CipherAESKey())
}

/*
Scenario: Store a biometric with one slave validation as KO
	Given a biometric with one slave validation as KO
	When I want to store the keychain
	Then I get the biometric is store on the KO database
*/
func TestStoreBiometricWithOneSlaveValidKO(t *testing.T) {
	repo := &databasemock{}

	lister := listing.NewService(repo)
	s := NewService(mockAiClient{}, repo, mockSigVerfier{}, lister, mockHasher{})

	sigs := account.NewSignatures("sig1", "sig2")

	end := mining.NewEndorsement(
		"", "hash", mining.NewMasterValidation([]string{}, "robotkey", mining.NewValidation(mining.ValidationOK, time.Now(), "robotkey", "sig")),
		[]mining.Validation{
			mining.NewValidation(mining.ValidationOK, time.Now(), "pub", "sig"),
			mining.NewValidation(mining.ValidationKO, time.Now(), "pub", "sig"),
		},
	)
	data := account.NewBiometricData("pHash", "enc addr robot", "enc addr person", "enc aes key", "person pub", "biod pub", sigs)
	bio := account.NewBiometric(data, end)

	err := s.StoreBiometric(bio)
	assert.Nil(t, err)

	assert.Empty(t, repo.biometrics)
	assert.Len(t, repo.biometricsKO, 1)
	assert.Equal(t, "enc aes key", repo.biometricsKO[0].CipherAESKey())
}

/*
Scenario: Store a biometric with a invalid transaction hash
	Given a biometric with invalid tx hash
	When I want to store it
	Then I get the error
*/
func TestStoreBiometricWithInvalidTxHash(t *testing.T) {
	repo := &databasemock{}
	lister := listing.NewService(repo)
	s := NewService(mockAiClient{}, repo, mockSigVerfier{}, lister, mockHasher{})

	sigs := account.NewSignatures("sig1", "sig2")
	end := mining.NewEndorsement(
		"", "bad hash", mining.NewMasterValidation([]string{}, "robotkey", mining.NewValidation(mining.ValidationOK, time.Now(), "robotkey", "sig")),
		[]mining.Validation{
			mining.NewValidation(mining.ValidationOK, time.Now(), "pub", "sig"),
		},
	)
	data := account.NewBiometricData("pHash", "enc addr robot", "enc addr person", "enc aes key", "person pub", "biod pub", sigs)
	bio := account.NewBiometric(data, end)

	assert.Equal(t, ErrInvalidDataIntegrity, s.StoreBiometric(bio))
}

/*
Scenario: Verify an endorsement with different masterpublic key
	Given an endorsement with an pow validation differents than the pow public key
	When I want to verify it
	Then I get an error
*/
func TestEndorsementWithDifferenteMasterPubKey(t *testing.T) {
	end := mining.NewEndorsement("", "hash",
		mining.NewMasterValidation([]string{}, "pubkey", mining.NewValidation(mining.ValidationOK, time.Now(), "other pub", "sig")), nil)

	s := service{}
	assert.Equal(t, ErrInvalidDataMining, s.verifyEndorsementSignatures(end))
}

/*
Scenario: Verify an endorsement with different masterpublic key
	Given an endorsement with an pow validation differents than the pow public key
	When I want to verify it
	Then I get an error
*/
func TestEndorsementWithBadMasterSignature(t *testing.T) {
	end := mining.NewEndorsement("", "hash",
		mining.NewMasterValidation([]string{}, "pubkey", mining.NewValidation(mining.ValidationOK, time.Now(), "pubkey", "invalid sig")), nil)

	s := service{
		sigVerif: mockBadSigVerfier{},
	}
	assert.Equal(t, ErrInvalidDataMining, s.verifyEndorsementSignatures(end))

	end2 := mining.NewEndorsement("", "hash",
		mining.NewMasterValidation([]string{}, "pubkey", mining.NewValidation(mining.ValidationOK, time.Now(), "pubkey", "sig")),
		[]mining.Validation{
			mining.NewValidation(mining.ValidationOK, time.Now(), "pub", "invalid sig"),
		},
	)
	assert.Equal(t, ErrInvalidDataMining, s.verifyEndorsementSignatures(end2))
}

type databasemock struct {
	biometrics   []account.Biometric
	keychains    []account.Keychain
	biometricsKO []account.Biometric
	keychainsKO  []account.Keychain
}

func (d *databasemock) StoreKeychain(kc account.Keychain) error {
	d.keychains = append(d.keychains, kc)
	return nil
}

func (d *databasemock) StoreBiometric(b account.Biometric) error {
	d.biometrics = append(d.biometrics, b)
	return nil
}

func (d *databasemock) StoreKOKeychain(kc account.Keychain) error {
	d.keychainsKO = append(d.keychainsKO, kc)
	return nil
}

func (d *databasemock) StoreKOBiometric(b account.Biometric) error {
	d.biometricsKO = append(d.biometricsKO, b)
	return nil
}

func (d *databasemock) FindBiometric(bh string) (account.Biometric, error) {
	for _, b := range d.biometrics {
		if b.PersonHash() == bh {
			return b, nil
		}
	}
	return nil, nil
}

func (d *databasemock) FindLastKeychain(addr string) (account.Keychain, error) {
	sort.Slice(d.keychains, func(i, j int) bool {
		iTimestamp := d.keychains[i].Endorsement().MasterValidation().ProofOfWorkValidation().Timestamp().Unix()
		jTimestamp := d.keychains[j].Endorsement().MasterValidation().ProofOfWorkValidation().Timestamp().Unix()
		return iTimestamp > jTimestamp
	})

	for _, b := range d.keychains {
		if b.Address() == addr {
			return b, nil
		}
	}
	return nil, nil
}

type mockSigVerfier struct{}

func (v mockSigVerfier) VerifyKeychainSignatures(account.Keychain) error {
	return nil
}
func (v mockSigVerfier) VerifyBiometricSignatures(account.Biometric) error {
	return nil
}
func (v mockSigVerfier) VerifyValidationSignature(mining.Validation) error {
	return nil
}

type mockBadSigVerfier struct{}

func (v mockBadSigVerfier) VerifyKeychainSignatures(account.Keychain) error {
	return nil
}
func (v mockBadSigVerfier) VerifyBiometricSignatures(account.Biometric) error {
	return nil
}
func (v mockBadSigVerfier) VerifyValidationSignature(valid mining.Validation) error {
	if valid.Signature() == "invalid sig" {
		return errors.New("invalid signature")
	}
	return nil
}

type mockAiClient struct{}

func (ai mockAiClient) CheckStorageAuthorization(txHash string) error {
	return nil
}

func (ai mockAiClient) GetMininumValidations(txHash string) (int, error) {
	return 1, nil
}

type mockHasher struct{}

func (h mockHasher) HashKeychainData(account.KeychainData) (string, error) {
	return "hash", nil
}

func (h mockHasher) HashBiometricData(account.BiometricData) (string, error) {
	return "hash", nil
}
