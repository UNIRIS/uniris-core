package mining

import (
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	datamining "github.com/uniris/uniris-core/datamining/pkg"
	biodlisting "github.com/uniris/uniris-core/datamining/pkg/biod/listing"
	"github.com/uniris/uniris-core/datamining/pkg/lock"
	memstorage "github.com/uniris/uniris-core/datamining/pkg/storage/mem"
)

/*
Scenario: Mine a transaction
	Given a transaction hash, data, biod sig, pools
	When I want to mine the transaction
	Then I get a master valid and a list of validations
*/
func TestMine(t *testing.T) {
	biodLister := biodlisting.NewService(memstorage.NewDatabase())

	s := service{
		notif:  &mockNotifier{},
		signer: mockSrvSigner{},
		poolR:  mockPoolRequester{},
		checks: map[TransactionType]Checker{
			CreateKeychainTransaction: mockCheck{},
		},
		biodLister: biodLister,
	}

	masterValid, valids, err := s.mine("txHash", "fake data", "biod sig",
		NewPool(Peer{}),
		NewPool(Peer{}),
		CreateKeychainTransaction)

	assert.Nil(t, err)
	assert.NotNil(t, masterValid)
	assert.NotEmpty(t, valids)

	assert.Equal(t, datamining.ValidationOK, masterValid.ProofOfWorkValidation().Status())
	assert.Equal(t, datamining.ValidationOK, valids[0].Status())
}

/*
Scenario: Lock a transaction
	Given a transaction hash, an address a validation pool
	When I want to lock the transaction
	Then the transaction is locked
¨*/
func TestLock(t *testing.T) {
	notif := &mockNotifier{}
	s := service{
		notif:  notif,
		signer: mockSrvSigner{},
		poolR:  mockPoolRequester{},
	}

	err := s.requestLock("txHash", "addr", NewPool(Peer{}))
	assert.Nil(t, err)
	assert.Equal(t, "Transaction txHash with status Locked", notif.lastNotif)
}

/*
Scenario: TestUnlock a transaction
	Given a transaction hash, an address a validation pool
	When I want to unlock the transaction
	Then the transaction is unlocked
¨*/
func TestUnlock(t *testing.T) {
	notif := &mockNotifier{}
	s := service{
		notif:  notif,
		signer: mockSrvSigner{},
		poolR:  mockPoolRequester{},
	}

	err := s.requestUnlock("txHash", "addr", NewPool(Peer{}))
	assert.Nil(t, err)
	assert.Equal(t, "Transaction txHash with status Unlocked", notif.lastNotif)
}

/*
Scenario: Validate data from a kind of transaction
	Given a transaction hash, data and a transaction type
	When I want to validate the data
	Then I get a validation with a status OK
*/
func TestValidateTx(t *testing.T) {
	s := service{
		checks: map[TransactionType]Checker{
			CreateKeychainTransaction: mockCheck{},
		},
		robotKey:   "pub key",
		robotPvKey: "pv key",
		signer:     mockSrvSigner{},
	}

	valid, err := s.Validate("hash", "fake data", CreateKeychainTransaction)
	assert.Nil(t, err)
	assert.NotNil(t, valid)
	assert.Equal(t, datamining.ValidationOK, valid.Status())
	assert.Equal(t, "sig", valid.Signature())
	assert.Equal(t, "pub key", valid.PublicKey())
	assert.NotEqual(t, time.Now(), valid.Timestamp())
}

/*
Scenario: Validate invalid data from a kind of transaction
	Given a transaction hash, invalid data and a transaction type
	When I want to validate the data
	Then I get a validation with a status NO ok
*/
func TestValidateInvalidTx(t *testing.T) {
	s := service{
		checks: map[TransactionType]Checker{
			CreateKeychainTransaction: mockBadCheck{},
		},
		robotKey:   "pub key",
		robotPvKey: "pv key",
		signer:     mockSrvSigner{},
	}

	valid, err := s.Validate("hash", "fake data", CreateKeychainTransaction)
	assert.Nil(t, err)
	assert.NotNil(t, valid)
	assert.Equal(t, datamining.ValidationKO, valid.Status())
	assert.Equal(t, "sig", valid.Signature())
	assert.Equal(t, "pub key", valid.PublicKey())
	assert.NotEqual(t, time.Now(), valid.Timestamp())
}

type mockSrvSigner struct{}

func (s mockSrvSigner) CheckTransactionSignature(pubk string, tx string, der string) error {
	return nil
}

func (s mockSrvSigner) SignValidation(v UnsignedValidation, pvKey string) (string, error) {
	return "sig", nil
}

func (s mockSrvSigner) SignLock(lock lock.TransactionLock, pvKey string) (string, error) {
	return "sig", nil
}

type mockCheck struct{}

func (c mockCheck) CheckAsMaster(txHash string, data interface{}) error {
	return nil
}

func (c mockCheck) CheckAsSlave(txHash string, data interface{}) error {
	return nil
}

type mockBadCheck struct{}

func (c mockBadCheck) CheckAsMaster(txHash string, data interface{}) error {
	return ErrInvalidTransaction
}

func (c mockBadCheck) CheckAsSlave(txHash string, data interface{}) error {
	return ErrInvalidTransaction
}

type mockNotifier struct {
	lastNotif string
}

func (n *mockNotifier) NotifyTransactionStatus(tx string, status TransactionStatus) error {
	n.lastNotif = fmt.Sprintf("Transaction %s with status %s", tx, status.String())
	return nil
}

type mockPoolRequester struct {
}

func (r mockPoolRequester) RequestLock(Pool, lock.TransactionLock, string) error {
	return nil
}

func (r mockPoolRequester) RequestUnlock(Pool, lock.TransactionLock, string) error {
	return nil
}

func (r mockPoolRequester) RequestValidations(sPool Pool, data interface{}, txType TransactionType) ([]datamining.Validation, error) {
	return []datamining.Validation{
		datamining.NewValidation(
			datamining.ValidationOK,
			time.Now(),
			"pubkey",
			"fake sig",
		)}, nil
}

func (r mockPoolRequester) RequestStorage(sPool Pool, data interface{}, end datamining.Endorsement, txType TransactionType) error {
	return nil
}
