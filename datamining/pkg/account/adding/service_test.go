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
	s := NewService(mockAiClient{}, repo, lister, mockSigVerfier{}, mockHasher{})

	end := mining.NewEndorsement(
		"", "hash",
		mining.NewMasterValidation([]string{}, "key1", mining.NewValidation(mining.ValidationOK, time.Now(), "robotkey", "sig")),
		[]mining.Validation{
			mining.NewValidation(mining.ValidationOK, time.Now(), "pub", "sig"),
		},
	)
	kc := account.NewKeychain("enc addr", "enc wallet", "id pub", "id sig", "em sig")
	eKc := account.NewEndorsedKeychain("addr", kc, end)

	err := s.StoreKeychain(eKc)
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
	s := NewService(mockAiClient{}, repo, lister, mockSigVerfier{}, mockHasher{})

	end := mining.NewEndorsement(
		"", "hash",
		mining.NewMasterValidation([]string{}, "key1", mining.NewValidation(mining.ValidationKO, time.Now(), "robotkey", "sig")),
		[]mining.Validation{
			mining.NewValidation(mining.ValidationOK, time.Now(), "pub", "sig"),
		},
	)
	kc := account.NewKeychain("enc addr", "enc wallet", "id pub", "id sig", "em sig")
	eKC := account.NewEndorsedKeychain("addr", kc, end)

	err := s.StoreKeychain(eKC)
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
	s := NewService(mockAiClient{}, repo, lister, mockSigVerfier{}, mockHasher{})

	end := mining.NewEndorsement(
		"", "hash",
		mining.NewMasterValidation([]string{}, "key1", mining.NewValidation(mining.ValidationOK, time.Now(), "robotkey", "sig")),
		[]mining.Validation{
			mining.NewValidation(mining.ValidationOK, time.Now(), "pub", "sig"),
			mining.NewValidation(mining.ValidationKO, time.Now(), "pub", "sig"),
		},
	)
	kc := account.NewKeychain("enc addr", "enc wallet", "id pub", "id sig", "em sig")
	eKC := account.NewEndorsedKeychain("addr", kc, end)

	err := s.StoreKeychain(eKC)
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
	s := NewService(mockAiClient{}, repo, lister, mockSigVerfier{}, mockHasher{})

	end1 := mining.NewEndorsement(
		"", "hash",
		mining.NewMasterValidation([]string{}, "key1", mining.NewValidation(mining.ValidationOK, time.Now(), "robotkey", "sig")),
		[]mining.Validation{
			mining.NewValidation(mining.ValidationOK, time.Now(), "pub", "sig"),
		},
	)

	kc := account.NewKeychain("enc addr", "enc wallet", "id pub", "id sig", "em sig")
	eKc1 := account.NewEndorsedKeychain("addr", kc, end1)

	assert.Nil(t, s.StoreKeychain(eKc1))

	end2 := mining.NewEndorsement(
		"bad last hash", "hash",
		mining.NewMasterValidation([]string{}, "key1", mining.NewValidation(mining.ValidationOK, time.Now(), "robotkey", "sig")),
		[]mining.Validation{
			mining.NewValidation(mining.ValidationOK, time.Now(), "pub", "sig"),
		},
	)

	eKc2 := account.NewEndorsedKeychain("addr", kc, end2)

	assert.Equal(t, ErrInvalidDataIntegrity, s.StoreKeychain(eKc2))

	end3 := mining.NewEndorsement(
		"", "hash",
		mining.NewMasterValidation([]string{}, "key1", mining.NewValidation(mining.ValidationOK, time.Now(), "robotkey", "sig")),
		[]mining.Validation{
			mining.NewValidation(mining.ValidationOK, time.Now(), "pub", "sig"),
		},
	)

	eKc3 := account.NewEndorsedKeychain("addr", kc, end3)

	assert.Equal(t, ErrInvalidDataIntegrity, s.StoreKeychain(eKc3))
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
	s := NewService(mockAiClient{}, repo, lister, mockSigVerfier{}, mockHasher{})

	end := mining.NewEndorsement(
		"", "hash",
		mining.NewMasterValidation([]string{}, "key1", mining.NewValidation(mining.ValidationOK, time.Now(), "robotkey", "sig")),
		[]mining.Validation{},
	)
	kc := account.NewKeychain("enc addr", "enc wallet", "id pub", "id sig", "em sig")
	eKC := account.NewEndorsedKeychain("addr", kc, end)

	assert.Equal(t, ErrInvalidValidationNumber, s.StoreKeychain(eKC))

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
	s := NewService(mockAiClient{}, repo, lister, mockSigVerfier{}, mockHasher{})

	end := mining.NewEndorsement(
		"", "bad hash",
		mining.NewMasterValidation([]string{}, "key1", mining.NewValidation(mining.ValidationOK, time.Now(), "robotkey", "sig")),
		[]mining.Validation{
			mining.NewValidation(mining.ValidationOK, time.Now(), "pub", "sig"),
		},
	)

	kc := account.NewKeychain("enc addr", "enc wallet", "id pub", "id sig", "em sig")
	eKC := account.NewEndorsedKeychain("addr", kc, end)

	assert.Equal(t, ErrInvalidDataIntegrity, s.StoreKeychain(eKC))
}

/*
Scenario: Store an ID
	Given an ID
	When I want to store it
	Then the ID is stored on the database
*/
func TestStoreID(t *testing.T) {
	repo := &databasemock{}
	lister := listing.NewService(repo)
	s := NewService(mockAiClient{}, repo, lister, mockSigVerfier{}, mockHasher{})

	end := mining.NewEndorsement(
		"", "hash",
		mining.NewMasterValidation([]string{}, "key1", mining.NewValidation(mining.ValidationOK, time.Now(), "robotkey", "sig")),
		[]mining.Validation{
			mining.NewValidation(mining.ValidationOK, time.Now(), "pub", "sig"),
		},
	)

	id := account.NewID("hash", "enc addr robot", "enc addr person", "enc aes key", "id pub", "id sig", "em pub")
	eID := account.NewEndorsedID(id, end)
	err := s.StoreID(eID)
	assert.Nil(t, err)
	l := len(repo.ids)
	assert.Equal(t, 1, l)
	assert.Equal(t, 1, l)
	assert.Equal(t, "hash", repo.ids[0].Hash())
}

/*
Scenario: Store a id with a zero
	Given a id without validations
	When I want to store it
	Then I get the error
*/
func TestStoreIDWithZeroValidations(t *testing.T) {
	repo := &databasemock{}
	lister := listing.NewService(repo)
	s := NewService(mockAiClient{}, repo, lister, mockSigVerfier{}, mockHasher{})

	end := mining.NewEndorsement(
		"", "hash",
		mining.NewMasterValidation([]string{}, "key1", mining.NewValidation(mining.ValidationOK, time.Now(), "robotkey", "sig")),
		[]mining.Validation{},
	)

	id := account.NewID("hash", "enc addr robot", "enc addr person", "enc aes key", "id pub", "id sig", "em pub")
	eID := account.NewEndorsedID(id, end)

	assert.Equal(t, ErrInvalidValidationNumber, s.StoreID(eID))

}

/*
Scenario: Store a ID with master validation KO
	Given an id with a master validation as KO
	When I want to store the id
	Then I get the id is store on the KO database
*/
func TestStoreIDWithMasterValidKO(t *testing.T) {
	repo := &databasemock{}

	lister := listing.NewService(repo)
	s := NewService(mockAiClient{}, repo, lister, mockSigVerfier{}, mockHasher{})

	end := mining.NewEndorsement(
		"", "hash",
		mining.NewMasterValidation([]string{}, "key1", mining.NewValidation(mining.ValidationKO, time.Now(), "robotkey", "sig")),
		[]mining.Validation{
			mining.NewValidation(mining.ValidationOK, time.Now(), "pub", "sig"),
		},
	)
	id := account.NewID("hash", "enc addr robot", "enc addr person", "enc aes key", "id pub", "id sig", "em pub")
	eID := account.NewEndorsedID(id, end)

	err := s.StoreID(eID)
	assert.Nil(t, err)

	assert.Empty(t, repo.ids)
	assert.Len(t, repo.idsKO, 1)
	assert.Equal(t, "enc aes key", repo.idsKO[0].EncryptedAESKey())
}

/*
Scenario: Store a id with one slave validation as KO
	Given a id with one slave validation as KO
	When I want to store the keychain
	Then I get the id is store on the KO database
*/
func TestStoreIDWithOneSlaveValidKO(t *testing.T) {
	repo := &databasemock{}

	lister := listing.NewService(repo)
	s := NewService(mockAiClient{}, repo, lister, mockSigVerfier{}, mockHasher{})

	end := mining.NewEndorsement(
		"", "hash",
		mining.NewMasterValidation([]string{}, "key1", mining.NewValidation(mining.ValidationOK, time.Now(), "robotkey", "sig")),
		[]mining.Validation{
			mining.NewValidation(mining.ValidationKO, time.Now(), "pub", "sig"),
		},
	)
	id := account.NewID("hash", "enc addr robot", "enc addr person", "enc aes key", "id pub", "id sig", "em pub")
	eID := account.NewEndorsedID(id, end)

	err := s.StoreID(eID)
	assert.Nil(t, err)

	assert.Empty(t, repo.ids)
	assert.Len(t, repo.idsKO, 1)
	assert.Equal(t, "enc aes key", repo.idsKO[0].EncryptedAESKey())
}

/*
Scenario: Store a id with a invalid transaction hash
	Given a id with invalid tx hash
	When I want to store it
	Then I get the error
*/
func TestStoreIDWithInvalidTxHash(t *testing.T) {
	repo := &databasemock{}
	lister := listing.NewService(repo)
	s := NewService(mockAiClient{}, repo, lister, mockSigVerfier{}, mockHasher{})

	end := mining.NewEndorsement(
		"", "bad hash",
		mining.NewMasterValidation([]string{}, "key1", mining.NewValidation(mining.ValidationOK, time.Now(), "robotkey", "sig")),
		[]mining.Validation{
			mining.NewValidation(mining.ValidationKO, time.Now(), "pub", "sig"),
		},
	)
	id := account.NewID("hash", "enc addr robot", "enc addr person", "enc aes key", "id pub", "id sig", "em pub")
	eID := account.NewEndorsedID(id, end)

	assert.Equal(t, ErrInvalidDataIntegrity, s.StoreID(eID))
}

/*
Scenario: Verify an endorsement with different masterpublic key
	Given an endorsement with an pow validation differents than the pow public key
	When I want to verify it
	Then I get an error
*/
func TestEndorsementWithBadMasterSignature(t *testing.T) {

	end := mining.NewEndorsement(
		"", "hash",
		mining.NewMasterValidation([]string{}, "key1", mining.NewValidation(mining.ValidationOK, time.Now(), "robotkey", "invalid sig")),
		[]mining.Validation{
			mining.NewValidation(mining.ValidationKO, time.Now(), "pub", "sig"),
		},
	)

	s := service{
		sigVerif: mockBadSigVerfier{},
	}
	assert.Equal(t, ErrInvalidDataMining, s.verifyEndorsementSignatures(end))

	end2 := mining.NewEndorsement(
		"", "hash",
		mining.NewMasterValidation([]string{}, "key1", mining.NewValidation(mining.ValidationOK, time.Now(), "robotkey", "invalid sig")),
		[]mining.Validation{
			mining.NewValidation(mining.ValidationKO, time.Now(), "pub", "sig"),
		},
	)
	assert.Equal(t, ErrInvalidDataMining, s.verifyEndorsementSignatures(end2))
}

type databasemock struct {
	ids         []account.EndorsedID
	keychains   []account.EndorsedKeychain
	idsKO       []account.EndorsedID
	keychainsKO []account.EndorsedKeychain
}

func (d *databasemock) StoreKeychain(kc account.EndorsedKeychain) error {
	d.keychains = append(d.keychains, kc)
	return nil
}

func (d *databasemock) StoreID(id account.EndorsedID) error {
	d.ids = append(d.ids, id)
	return nil
}

func (d *databasemock) StoreKOKeychain(kc account.EndorsedKeychain) error {
	d.keychainsKO = append(d.keychainsKO, kc)
	return nil
}

func (d *databasemock) StoreKOID(id account.EndorsedID) error {
	d.idsKO = append(d.idsKO, id)
	return nil
}

func (d *databasemock) FindID(hash string) (account.EndorsedID, error) {
	for _, id := range d.ids {
		if id.Hash() == hash {
			return id, nil
		}
	}
	return nil, nil
}

func (d *databasemock) FindLastKeychain(addr string) (account.EndorsedKeychain, error) {
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
func (v mockSigVerfier) VerifyIDSignatures(account.ID) error {
	return nil
}
func (v mockSigVerfier) VerifyValidationSignature(mining.Validation) error {
	return nil
}

func (v mockSigVerfier) VerifyTransactionDataSignature(txType mining.TransactionType, pubk string, data interface{}, der string) error {
	return nil
}

type mockBadSigVerfier struct{}

func (v mockBadSigVerfier) VerifyKeychainSignatures(account.Keychain) error {
	return nil
}
func (v mockBadSigVerfier) VerifyIDSignatures(account.ID) error {
	return nil
}
func (v mockBadSigVerfier) VerifyValidationSignature(valid mining.Validation) error {
	if valid.Signature() == "invalid sig" {
		return errors.New("invalid signature")
	}
	return nil
}

func (v mockBadSigVerfier) VerifyTransactionDataSignature(txType mining.TransactionType, pubk string, data interface{}, der string) error {
	return errors.New("Invalid signature")
}

type mockAiClient struct{}

func (ai mockAiClient) CheckStorageAuthorization(txHash string) error {
	return nil
}

func (ai mockAiClient) GetMininumValidations(txHash string) (int, error) {
	return 1, nil
}

type mockHasher struct{}

func (h mockHasher) HashKeychain(account.Keychain) (string, error) {
	return "hash", nil
}

func (h mockHasher) HashEndorsedKeychain(account.EndorsedKeychain) (string, error) {
	return "hash", nil
}

func (h mockHasher) HashID(account.ID) (string, error) {
	return "hash", nil
}

func (h mockHasher) HashEndorsedID(account.EndorsedID) (string, error) {
	return "hash", nil
}
