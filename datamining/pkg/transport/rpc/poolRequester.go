package rpc

import (
	"errors"

	datamining "github.com/uniris/uniris-core/datamining/pkg"
	"github.com/uniris/uniris-core/datamining/pkg/account"
	"github.com/uniris/uniris-core/datamining/pkg/lock"
	"github.com/uniris/uniris-core/datamining/pkg/mining"
	"github.com/uniris/uniris-core/datamining/pkg/system"
)

//PoolRequester define methods for pool requesting
type PoolRequester interface {
	mining.PoolRequester
	account.PoolRequester
}

type poolR struct {
	conf   system.UnirisConfig
	crypto Crypto
	api    apiBuilder
	data   dataBuilder
}

//NewPoolRequester creates a new pool requester using GRPC
func NewPoolRequester(conf system.UnirisConfig, crypto Crypto) PoolRequester {
	return poolR{
		conf:   conf,
		crypto: crypto,
		api:    apiBuilder{},
		data:   dataBuilder{},
	}
}

func (pR poolR) RequestBiometric(sPool datamining.Pool, personHash string) (account.Biometric, error) {

	biometrics := make([]account.Biometric, 0)

	for _, p := range sPool.Peers() {

		cli := newExternalClient(p.IP.String(), pR.conf.Datamining.ExternalPort, pR.crypto, pR.conf)
		b, err := cli.RequestBiometric(personHash)
		if err != nil {
			return nil, err
		}

		biometrics = append(biometrics, b)
	}

	if len(biometrics) == 0 {
		return nil, errors.New(pR.conf.Datamining.Errors.AccountNotExist)
	}

	return biometrics[0], nil
}

func (pR poolR) RequestKeychain(sPool datamining.Pool, address string) (account.Keychain, error) {

	keychains := make([]account.Keychain, 0)

	for _, p := range sPool.Peers() {

		cli := newExternalClient(p.IP.String(), pR.conf.Datamining.ExternalPort, pR.crypto, pR.conf)
		k, err := cli.RequestKeychain(address)
		if err != nil {
			return nil, err
		}
		keychains = append(keychains, k)
	}

	if len(keychains) == 0 {
		return nil, errors.New(pR.conf.Datamining.Errors.AccountNotExist)
	}

	return keychains[0], nil
}

func (pR poolR) RequestLock(lastValidPool datamining.Pool, txLock lock.TransactionLock) error {

	//TODO: using goroutines
	for _, p := range lastValidPool.Peers() {
		cli := newExternalClient(p.IP.String(), pR.conf.Datamining.ExternalPort, pR.crypto, pR.conf)
		if err := cli.RequestLock(txLock); err != nil {
			return err
		}
	}

	return nil
}
func (pR poolR) RequestUnlock(lastValidPool datamining.Pool, txLock lock.TransactionLock) error {

	//TODO: using goroutines
	for _, p := range lastValidPool.Peers() {
		cli := newExternalClient(p.IP.String(), pR.conf.Datamining.ExternalPort, pR.crypto, pR.conf)
		if err := cli.RequestUnlock(txLock); err != nil {
			return err
		}
	}

	return nil
}

func (pR poolR) RequestValidations(validPool datamining.Pool, txHash string, data interface{}, txType mining.TransactionType) ([]mining.Validation, error) {

	valids := make([]mining.Validation, 0)

	//TODO: using goroutines
	for _, p := range validPool.Peers() {
		cli := newExternalClient(p.IP.String(), pR.conf.Datamining.ExternalPort, pR.crypto, pR.conf)
		v, err := cli.RequestValidation(txType, txHash, data)
		if err != nil {
			return nil, err
		}

		valids = append(valids, v)
	}

	return valids, nil
}

func (pR poolR) RequestStorage(sPool datamining.Pool, data interface{}, end mining.Endorsement, txType mining.TransactionType) error {

	//TODO: using goroutines
	for _, p := range sPool.Peers() {
		cli := newExternalClient(p.IP.String(), pR.conf.Datamining.ExternalPort, pR.crypto, pR.conf)
		if err := cli.RequestStorage(txType, data, end); err != nil {
			return err
		}
	}

	return nil
}
