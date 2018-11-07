package listing

import (
	"sort"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	datamining "github.com/uniris/uniris-core/datamining/pkg"
	"github.com/uniris/uniris-core/datamining/pkg/account"
)

/*
Scenario: TestGetLastKeychain keychain
	Given a empty database
	When I add two keychain
	Then I get the last keychain
*/
func TestGetLastKeychain(t *testing.T) {
	db := new(databasemock)
	s := NewService(db)

	sigs := account.Signatures{
		BiodSig:   "sig1",
		PersonSig: "sig2",
	}

	keychainData := &account.KeyChainData{
		WalletAddr:      "address",
		CipherAddrRobot: "xxxx",
		CipherWallet:    "xxxx",
		PersonPubk:      "xxxx",
		BiodPubk:        "xxxx",
		Sigs:            sigs,
	}

	masterValid1 := datamining.NewMasterValidation([]string{}, "robotKey", datamining.NewValidation(datamining.ValidationOK, time.Now(), "pubkey", "signature"))
	endors1 := datamining.NewEndorsement("", "hash1", masterValid1, []datamining.Validation{})

	kc1 := account.NewKeychain(keychainData, endors1)

	time.Sleep(1 * time.Second)

	masterValid2 := datamining.NewMasterValidation([]string{}, "robotKey", datamining.NewValidation(datamining.ValidationOK, time.Now(), "pubkey", "signature"))
	endors2 := datamining.NewEndorsement("hash1", "hash2", masterValid2, []datamining.Validation{})

	kc2 := account.NewKeychain(keychainData, endors2)

	db.Keychains = append(db.Keychains, kc1)
	db.Keychains = append(db.Keychains, kc2)

	keychain, err := s.GetLastKeychain("address")
	assert.Nil(t, err)
	assert.NotNil(t, keychain)
	assert.Equal(t, "hash2", keychain.Endorsement().TransactionHash())
	assert.Equal(t, "hash1", keychain.Endorsement().LastTransactionHash())
}

/*
Scenario: Get biometric
	Given a empty database
	When I add a Biometric
	Then return values of a GetBiometric are the exepeted ones
*/
func TestGetBiometric(t *testing.T) {

	db := new(databasemock)
	s := NewService(db)

	sigs := account.Signatures{
		BiodSig:   "sig1",
		PersonSig: "sig2",
	}

	bdata := &account.BioData{
		PersonHash:      "hash1",
		BiodPubk:        "xxxx",
		CipherAddrBio:   "xxxx",
		CipherAddrRobot: "xxxx",
		CipherAESKey:    "xxxx",
		PersonPubk:      "xxxx",
		Sigs:            sigs,
	}

	masterValid := datamining.NewMasterValidation([]string{}, "robotKey", datamining.NewValidation(datamining.ValidationOK, time.Now(), "pubkey", "signature"))
	endors := datamining.NewEndorsement("", "xxx", masterValid, []datamining.Validation{})

	b := account.NewBiometric(bdata, endors)

	db.Biometrics = append(db.Biometrics, b)
	wa, err := s.GetBiometric("hash1")
	assert.Nil(t, err)
	assert.NotNil(t, wa)
}

type databasemock struct {
	Biometrics []account.Biometric
	Keychains  []account.Keychain
}

func (d *databasemock) FindBiometric(bh string) (account.Biometric, error) {
	for _, b := range d.Biometrics {
		if b.PersonHash() == bh {
			return b, nil
		}
	}
	return nil, nil
}

func (d *databasemock) FindLastKeychain(addr string) (account.Keychain, error) {
	sort.Slice(d.Keychains, func(i, j int) bool {
		iTimestamp := d.Keychains[i].Endorsement().MasterValidation().ProofOfWorkValidation().Timestamp().Unix()
		jTimestamp := d.Keychains[j].Endorsement().MasterValidation().ProofOfWorkValidation().Timestamp().Unix()
		return iTimestamp > jTimestamp
	})

	for _, b := range d.Keychains {
		if b.WalletAddr() == addr {
			return b, nil
		}
	}
	return nil, nil
}
