package mock

import (
	"time"

	datamining "github.com/uniris/uniris-core/datamining/pkg"
	"github.com/uniris/uniris-core/datamining/pkg/account"
	"github.com/uniris/uniris-core/datamining/pkg/lock"
	"github.com/uniris/uniris-core/datamining/pkg/mining"
)

//PoolRequester define methods to handle pool requesting
type PoolRequester interface {
	mining.PoolRequester
	account.PoolRequester
}

//NewPoolRequester create a mock pool requester
func NewPoolRequester(db Repo) PoolRequester {
	return mockPoolRequester{db}
}

type mockPoolRequester struct {
	Repo Repo
}

func (r mockPoolRequester) RequestBiometric(pool datamining.Pool, personHash string) (account.Biometric, error) {
	return account.NewBiometric(
		&account.BioData{
			BiodPubk:        "pub",
			CipherAddrBio:   "enc addr",
			CipherAddrRobot: "enc addr",
			CipherAESKey:    "enc aes key",
			PersonHash:      personHash,
			PersonPubk:      "pub",
		},
		datamining.NewEndorsement(
			"",
			"hash",
			datamining.NewMasterValidation(
				[]string{"hash"},
				"robotkey",
				datamining.NewValidation(datamining.ValidationOK, time.Now(), "pubkey", "sig"),
			),
			[]datamining.Validation{},
		),
	), nil
}

func (r mockPoolRequester) RequestKeychain(pool datamining.Pool, address string) (account.Keychain, error) {
	return account.NewKeychain(
		&account.KeyChainData{
			BiodPubk:        "pub",
			CipherAddrRobot: "enc addr",
			CipherWallet:    "enc wallet",
			PersonPubk:      "pub",
		},
		datamining.NewEndorsement(
			"",
			"hash",
			datamining.NewMasterValidation(
				[]string{"hash"},
				"robotkey",
				datamining.NewValidation(datamining.ValidationOK, time.Now(), "pubkey", "sig"),
			),
			[]datamining.Validation{},
		),
	), nil
}

func (r mockPoolRequester) RequestLock(datamining.Pool, lock.TransactionLock, string) error {
	return nil
}

func (r mockPoolRequester) RequestUnlock(datamining.Pool, lock.TransactionLock, string) error {
	return nil
}

func (r mockPoolRequester) RequestValidations(sPool datamining.Pool, txHash string, data interface{}, txType mining.TransactionType) ([]datamining.Validation, error) {
	return []datamining.Validation{
		datamining.NewValidation(
			datamining.ValidationOK,
			time.Now(),
			"pubkey",
			"fake sig",
		)}, nil
}

func (r mockPoolRequester) RequestStorage(sPool datamining.Pool, data interface{}, end datamining.Endorsement, txType mining.TransactionType) error {
	switch data.(type) {
	case *account.KeyChainData:
		data := data.(*account.KeyChainData)
		kc := account.NewKeychain(data, end)
		r.Repo.StoreKeychain(kc)
	case *account.BioData:
		data := data.(*account.BioData)
		bio := account.NewBiometric(data, end)
		r.Repo.StoreBiometric(bio)
	}

	return nil
}
