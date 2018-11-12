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
	sig := account.NewSignatures("sig1", "sig2")
	data := account.NewBiometricData("personHash", "enc addr", "enc addr", "enc aes key", "pub", "pub", sig)
	err := miner.checkDataIntegrity("hash", data)
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
	sig := account.NewSignatures("sig1", "sig2")
	data := account.NewBiometricData("personHash", "enc addr", "enc addr", "enc aes key", "pub", "pub", sig)
	err := miner.checkDataIntegrity("hash", data)
	assert.Equal(t, mining.ErrInvalidTransaction, err)
}

/*
Scenario: Verifies the biometric data signature
	Given biometric data
	When I want to check if the signature match the transaction
	Then I get no errors
*/
func TestBiometricSignature(t *testing.T) {
	miner := biometricMiner{signer: mockBiometricSigner{}}
	sig := account.NewSignatures("sig1", "sig2")
	data := account.NewBiometricData("personHash", "enc addr", "enc addr", "enc aes key", "pub", "pub", sig)
	err := miner.verifyDataSignature(data)
	assert.Nil(t, err)
}

/*
Scenario: Check biometric data as master peer
	Given a transaction hash and biometric data
	When I want to check it as master
	Then I get not error
*/
func TestBiometricMasterCheck(t *testing.T) {
	miner := NewBiometricMiner(mockBiometricSigner{}, mockBiometricHasher{})
	sig := account.NewSignatures("sig1", "sig2")
	data := account.NewBiometricData("personHash", "enc addr", "enc addr", "enc aes key", "pub", "pub", sig)
	err := miner.CheckAsMaster("hash", data)
	assert.Nil(t, err)
}

/*
Scenario: Check biometric data as slave peer
	Given a transaction hash and biometric data
	When I want to check it as slave
	Then I get not error
*/
func TestBiometricSlaveCheck(t *testing.T) {
	miner := NewBiometricMiner(mockBiometricSigner{}, mockBiometricHasher{})
	sig := account.NewSignatures("sig1", "sig2")
	data := account.NewBiometricData("personHash", "enc addr", "enc addr", "enc aes key", "pub", "pub", sig)
	err := miner.CheckAsSlave("hash", data)
	assert.Nil(t, err)
}

type mockBiometricHasher struct{}

func (h mockBiometricHasher) NewBiometricDataHash(account.BiometricData) (string, error) {
	return "hash", nil
}

type mockBadBiometricHasher struct{}

func (h mockBadBiometricHasher) NewBiometricDataHash(account.BiometricData) (string, error) {
	return "other hash", nil
}

type mockBiometricSigner struct{}

func (s mockBiometricSigner) VerifyBiometricDataSignature(pubK string, data account.BiometricData, sig string) error {
	return nil
}
