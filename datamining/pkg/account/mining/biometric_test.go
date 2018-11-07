package mining

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/uniris/uniris-core/datamining/pkg/account"
	"github.com/uniris/uniris-core/datamining/pkg/mining"
)

/*
Scenario: Checks the biometric data integrity
	Given a transaction hash and biometric data
	When I want to check if the data match the transaction
	Then I get no errors
*/
func TestBiometricIntegrity(t *testing.T) {
	miner := biometricMiner{hasher: mockBiometricHasher{}}
	err := miner.checkDataIntegrity("hash", &account.BioData{})
	assert.Nil(t, err)
}

/*
Scenario: Checks the biometric data integrity
	Given a invalid transaction hash for a biometric data
	When I want to check if the data match the transaction
	Then I get an errors
*/
func TestInvalidBiometricIntegrity(t *testing.T) {
	miner := biometricMiner{hasher: mockBadBiometricHasher{}}
	err := miner.checkDataIntegrity("hash", &account.BioData{})
	assert.Equal(t, mining.ErrInvalidTransaction, err)
}

/*
Scenario: Checks the biometric data signature
	Given biometric data
	When I want to check if the signature match the transaction
	Then I get no errors
*/
func TestBiometricSignature(t *testing.T) {
	miner := biometricMiner{signer: mockBiometricSigner{}}
	err := miner.checkDataSignature(&account.BioData{})
	assert.Nil(t, err)
}

/*
Scenario: Check biometric data as master node
	Given a transaction hash and biometric data
	When I want to check it as master
	Then I get not error
*/
func TestBiometricMasterCheck(t *testing.T) {
	miner := NewBiometricMiner(mockBiometricSigner{}, mockBiometricHasher{})
	err := miner.CheckAsMaster("hash", &account.BioData{})
	assert.Nil(t, err)
}

/*
Scenario: Check biometric data as slave node
	Given a transaction hash and biometric data
	When I want to check it as slave
	Then I get not error
*/
func TestBiometricSlaveCheck(t *testing.T) {
	miner := NewBiometricMiner(mockBiometricSigner{}, mockBiometricHasher{})
	err := miner.CheckAsSlave("hash", &account.BioData{})
	assert.Nil(t, err)
}

type mockBiometricHasher struct{}

func (h mockBiometricHasher) HashUnsignedBiometricData(data UnsignedBiometricData) (string, error) {
	return "hash", nil
}

type mockBadBiometricHasher struct{}

func (h mockBadBiometricHasher) HashUnsignedBiometricData(data UnsignedBiometricData) (string, error) {
	return "other hash", nil
}

type mockBiometricSigner struct{}

func (s mockBiometricSigner) CheckBiometricSignature(pubK string, data UnsignedBiometricData, sig string) error {
	return nil
}
