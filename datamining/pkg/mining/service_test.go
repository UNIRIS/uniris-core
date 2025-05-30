package mining

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	datamining "github.com/uniris/uniris-core/datamining/pkg"
	"github.com/uniris/uniris-core/datamining/pkg/emitter"
	emListing "github.com/uniris/uniris-core/datamining/pkg/emitter/listing"
	"github.com/uniris/uniris-core/datamining/pkg/lock"
	"github.com/uniris/uniris-core/datamining/pkg/system"
)

/*
Scenario: Mine a transaction
	Given a transaction hash, data, biod sig, pools
	When I want to mine the transaction
	Then I get a master valid and a list of validations
*/
func TestMine(t *testing.T) {
	emLister := emListing.NewService(&mockEmDatabase{})

	s := service{
		aiClient: mockAIClient{},
		signer:   mockSrvSigner{},
		poolR:    mockPoolRequester{},
		txMiners: map[TransactionType]TransactionMiner{
			KeychainTransaction: mockMiner{},
		},
		emLister: emLister,
	}

	endorsement, err := s.mine("txHash", "fake data", "addr", "biod sig",
		datamining.NewPool(datamining.Peer{}),
		datamining.NewPool(datamining.Peer{}),
		KeychainTransaction)

	assert.Nil(t, err)
	assert.NotNil(t, endorsement.MasterValidation())
	assert.NotEmpty(t, endorsement.Validations())

	assert.Equal(t, ValidationOK, endorsement.MasterValidation().ProofOfWorkValidation().Status())
	assert.Equal(t, ValidationOK, endorsement.Validations()[0].Status())
}

/*
Scenario: Lock a transaction
	Given a transaction hash, an address a validation pool
	When I want to lock the transaction
	Then the transaction is locked
¨*/
func TestLock(t *testing.T) {
	s := service{
		signer: mockSrvSigner{},
		poolR:  mockPoolRequester{},
	}

	err := s.requestLock("txHash", "addr", datamining.NewPool(datamining.Peer{}))
	assert.Nil(t, err)
}

/*
Scenario: TestUnlock a transaction
	Given a transaction hash, an address a validation pool
	When I want to unlock the transaction
	Then the transaction is unlocked
¨*/
func TestUnlock(t *testing.T) {
	s := service{
		signer: mockSrvSigner{},
		poolR:  mockPoolRequester{},
	}

	err := s.requestUnlock("txHash", "addr", datamining.NewPool(datamining.Peer{}))
	assert.Nil(t, err)
}

/*
Scenario: Validate data from a kind of transaction
	Given a transaction hash, data and a transaction type
	When I want to validate the data
	Then I get a validation with a status OK
*/
func TestValidateTx(t *testing.T) {
	s := service{
		txMiners: map[TransactionType]TransactionMiner{
			KeychainTransaction: mockMiner{},
		},
		config: system.UnirisConfig{
			PublicKey: "pub key",
			SharedKeys: system.SharedKeys{
				Robot: system.KeyPair{
					PublicKey: "pub key",
				},
			},
		},
		signer: mockSrvSigner{},
	}

	valid, err := s.Validate("hash", "fake data", KeychainTransaction)
	assert.Nil(t, err)
	assert.NotNil(t, valid)
	assert.Equal(t, ValidationOK, valid.Status())
	assert.Equal(t, "sig", valid.Signature())
	assert.Equal(t, "pub key", valid.PublicKey())
	assert.Equal(t, time.Now().Unix(), valid.Timestamp().Unix())
}

/*
Scenario: Validate invalid data from a kind of transaction
	Given a transaction hash, invalid data and a transaction type
	When I want to validate the data
	Then I get a validation with a status NO ok
*/
func TestValidateInvalidTx(t *testing.T) {
	s := service{
		txMiners: map[TransactionType]TransactionMiner{
			KeychainTransaction: mockBadMiner{},
		},
		config: system.UnirisConfig{
			PublicKey: "pub key",
			SharedKeys: system.SharedKeys{
				Robot: system.KeyPair{
					PublicKey: "pub key",
				},
			},
		},
		signer: mockSrvSigner{},
	}

	valid, err := s.Validate("hash", "fake data", KeychainTransaction)
	assert.Nil(t, err)
	assert.NotNil(t, valid)
	assert.Equal(t, ValidationKO, valid.Status())
	assert.Equal(t, "sig", valid.Signature())
	assert.Equal(t, "pub key", valid.PublicKey())
	assert.Equal(t, time.Now().Unix(), valid.Timestamp().Unix())
}

type mockSrvSigner struct{}

func (s mockSrvSigner) VerifyTransactionDataSignature(txType TransactionType, pubk string, data interface{}, der string) error {
	return nil
}

func (s mockSrvSigner) SignValidation(v Validation, pvKey string) (Validation, error) {
	return NewValidation(v.Status(), v.Timestamp(), v.PublicKey(), "sig"), nil
}

type mockMiner struct{}

func (c mockMiner) CheckAsMaster(txHash string, data interface{}) error {
	return nil
}

func (c mockMiner) CheckAsSlave(txHash string, data interface{}) error {
	return nil
}

func (c mockMiner) GetLastTransactionHash(addr string) (string, error) {
	return "", nil
}

type mockBadMiner struct{}

func (c mockBadMiner) CheckAsMaster(txHash string, data interface{}) error {
	return ErrInvalidTransaction
}

func (c mockBadMiner) CheckAsSlave(txHash string, data interface{}) error {
	return ErrInvalidTransaction
}

func (c mockBadMiner) GetLastTransactionHash(addr string) (string, error) {
	return "", nil
}

type mockPoolRequester struct {
}

func (r mockPoolRequester) RequestLock(datamining.Pool, lock.TransactionLock) error {
	return nil
}

func (r mockPoolRequester) RequestUnlock(datamining.Pool, lock.TransactionLock) error {
	return nil
}

func (r mockPoolRequester) RequestValidations(minValid int, sPool datamining.Pool, txHash string, data interface{}, txType TransactionType) ([]Validation, error) {
	return []Validation{
		NewValidation(
			ValidationOK,
			time.Now(),
			"pubkey",
			"fake sig",
		)}, nil
}

func (r mockPoolRequester) RequestStorage(minReplicas int, sPool datamining.Pool, data interface{}, end Endorsement, txType TransactionType) error {
	return nil
}

type mockEmDatabase struct{}

func (d *mockEmDatabase) ListSharedEmitterKeyPairs() ([]emitter.SharedKeyPair, error) {
	return []emitter.SharedKeyPair{
		emitter.SharedKeyPair{
			PublicKey: "key1",
		}, emitter.SharedKeyPair{
			PublicKey: "key2",
		}, emitter.SharedKeyPair{
			PublicKey: "key3",
		}}, nil
}

type mockAIClient struct{}

func (ai mockAIClient) GetMininumValidations(txHash string) (int, error) {
	return 1, nil
}

func (ai mockAIClient) GetMininumReplications(txHash string) (int, error) {
	return 1, nil
}
