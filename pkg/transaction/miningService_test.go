package transaction

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/x509"
	"encoding/hex"
	"encoding/json"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/uniris/uniris-core/pkg/crypto"
	"github.com/uniris/uniris-core/pkg/shared"
)

/*
Scenario: Get the minimum number of a transaction replicas
	Given a transaction hash
	When I want to get the minimum replicas
	Then I get a number  valid
	//TODO: to improve when the implementation will be defined
*/
func TestGetMinimumReplicas(t *testing.T) {
	s := MiningService{}
	assert.Equal(t, 1, s.getMinimumReplicas(""))
}

/*
Scenario: Get the minimum validation number
	Given a transaction hash
	When I want to get the validation required number
	Then I get a number  valid
	//TODO: to improve when the implementation will be defined
*/
func TestGetMinimumTransactionValidation(t *testing.T) {
	s := MiningService{}
	assert.Equal(t, 1, s.GetMinimumTransactionValidation(""))
}

/*
Scenario: Create a miner validation
	Given a validation status
	When I want to create miner validation
	Then I get a validation signed
*/
func TestBuildMinerValidation(t *testing.T) {
	pub, pv := crypto.GenerateKeys()

	s := MiningService{
		minerPubK: pub,
		minerPvk:  pv,
	}

	v, err := s.buildMinerValidation(ValidationOK)
	assert.Nil(t, err)
	assert.Equal(t, pub, v.MinerPublicKey())
	assert.Nil(t, err)
	assert.Equal(t, time.Now().Unix(), v.Timestamp().Unix())
	assert.Equal(t, ValidationOK, v.Status())
	ok, err := v.IsValid()
	assert.True(t, ok)
}

/*
Scenario: Validate an incoming transaction
	Given a valid transaction
	When I want to valid the transaction
	Then I get a validation with status OK
*/
func TestValidateTransaction(t *testing.T) {
	pub, pv := crypto.GenerateKeys()

	vBytes, _ := json.Marshal(MinerValidation{
		minerPubk: pub,
		status:    ValidationOK,
		timestamp: time.Now(),
	})
	vSig, _ := crypto.Sign(string(vBytes), pv)
	v, _ := NewMinerValidation(ValidationOK, time.Now(), pub, vSig)

	mv, _ := NewMasterValidation(Pool{}, pub, v)

	sk, _ := shared.NewEmitterKeyPair(hex.EncodeToString([]byte("pvKey")), pub)
	prop, _ := NewProposal(sk)

	data := map[string]string{
		"encrypted_aes_key":          hex.EncodeToString([]byte("aesKey")),
		"encrypted_address_by_robot": hex.EncodeToString([]byte("addr")),
		"encrypted_address_by_id":    hex.EncodeToString([]byte("addr")),
	}

	tx := Transaction{
		address:   crypto.HashString("addr"),
		data:      data,
		txType:    KeychainType,
		timestamp: time.Now(),
		pubKey:    pub,
		prop:      prop,
	}
	txBytesBeforeSig, _ := tx.MarshalBeforeSignature()
	sig, _ := crypto.Sign(string(txBytesBeforeSig), pv)
	tx.emSig = sig
	tx.sig = sig
	txBytes, _ := tx.MarshalHash()
	txHash := crypto.HashBytes(txBytes)
	tx.txHash = txHash

	s := MiningService{
		minerPubK: pub,
		minerPvk:  pv,
	}
	valid, err := s.ValidateTransaction(tx, mv)
	assert.Nil(t, err)
	assert.Equal(t, ValidationOK, valid.Status())
}

/*
Scenario: Validate an incoming transaction with invalid integrity
	Given a transaction with invalid transaction hash or signature
	When I want to valid the transaction
	Then I get a validation with status KO
*/
func TestValidateTransactionWithBadIntegrity(t *testing.T) {
	pub, pv := crypto.GenerateKeys()

	vBytes, _ := json.Marshal(MinerValidation{
		minerPubk: pub,
		status:    ValidationOK,
		timestamp: time.Now(),
	})
	vSig, _ := crypto.Sign(string(vBytes), pv)
	v, _ := NewMinerValidation(ValidationOK, time.Now(), pub, vSig)

	mv, _ := NewMasterValidation(Pool{}, pub, v)

	sk, _ := shared.NewEmitterKeyPair(hex.EncodeToString([]byte("pvKey")), pub)
	prop, _ := NewProposal(sk)

	data := map[string]string{
		"encrypted_aes_key":          hex.EncodeToString([]byte("aesKey")),
		"encrypted_address_by_robot": hex.EncodeToString([]byte("addr")),
		"encrypted_address_by_id":    hex.EncodeToString([]byte("addr")),
	}

	sig, _ := crypto.Sign("hello", pv)
	tx, _ := New(crypto.HashString("addr"), IDType, data, time.Now(), pub, sig, sig, prop, crypto.HashString("hash"))

	s := MiningService{
		minerPvk: pv,
	}
	valid, err := s.ValidateTransaction(tx, mv)
	assert.Nil(t, err)
	assert.Equal(t, ValidationKO, valid.Status())
}

/*
Scenario: request transaction validations
	Given a transaction to validate
	When I ask validations to a pool
	Then I get validations from them
*/
func TestRequestValidations(t *testing.T) {
	s := MiningService{
		poolR: &mockPoolRequester{},
	}

	pub, pv := crypto.GenerateKeys()

	vBytes, _ := json.Marshal(MinerValidation{
		minerPubk: pub,
		status:    ValidationOK,
		timestamp: time.Now(),
	})
	vSig, _ := crypto.Sign(string(vBytes), pv)
	v, _ := NewMinerValidation(ValidationOK, time.Now(), pub, vSig)

	mv, _ := NewMasterValidation(Pool{}, pub, v)

	sk, _ := shared.NewEmitterKeyPair(hex.EncodeToString([]byte("pvKey")), pub)
	prop, _ := NewProposal(sk)

	data := map[string]string{
		"encrypted_aes_key":          hex.EncodeToString([]byte("aesKey")),
		"encrypted_address_by_robot": hex.EncodeToString([]byte("addr")),
		"encrypted_address_by_id":    hex.EncodeToString([]byte("addr")),
	}

	txBytes, _ := json.Marshal(Transaction{
		address:   crypto.HashString("addr"),
		txType:    IDType,
		data:      data,
		timestamp: time.Now(),
		pubKey:    pub,
		prop:      prop,
	})
	sig, _ := crypto.Sign(string(txBytes), pv)

	tx, _ := New(crypto.HashString("addr"), IDType, data, time.Now(), pub, sig, sig, prop, crypto.HashBytes(txBytes))

	valids, err := s.requestValidations(tx, mv, Pool{}, 1)
	assert.Nil(t, err)
	assert.NotEmpty(t, valids)
	assert.Equal(t, ValidationOK, valids[0].Status())
}

/*
Scenario: request transaction storage
	Given a transaction to store
	When I ask storage to a pool
	Then I get acks from them
*/
func TestRequestStorage(t *testing.T) {
	poolR := &mockPoolRequester{}
	s := MiningService{
		poolR: poolR,
	}

	pub, pv := crypto.GenerateKeys()

	vBytes, _ := json.Marshal(MinerValidation{
		minerPubk: pub,
		status:    ValidationOK,
		timestamp: time.Now(),
	})
	vSig, _ := crypto.Sign(string(vBytes), pv)
	v, _ := NewMinerValidation(ValidationOK, time.Now(), pub, vSig)

	mv, _ := NewMasterValidation(Pool{}, pub, v)

	sk, _ := shared.NewEmitterKeyPair(hex.EncodeToString([]byte("pvKey")), pub)
	prop, _ := NewProposal(sk)

	data := map[string]string{
		"encrypted_aes_key":          hex.EncodeToString([]byte("aesKey")),
		"encrypted_address_by_robot": hex.EncodeToString([]byte("addr")),
		"encrypted_address_by_id":    hex.EncodeToString([]byte("addr")),
	}

	txBytes, _ := json.Marshal(Transaction{
		address:   crypto.HashString("addr"),
		txType:    IDType,
		data:      data,
		timestamp: time.Now(),
		pubKey:    pub,
		prop:      prop,
	})
	sig, _ := crypto.Sign(string(txBytes), pv)

	tx, _ := New(crypto.HashString("addr"), IDType, data, time.Now(), pub, sig, sig, prop, crypto.HashBytes(txBytes))
	tx.AddMining(mv, []MinerValidation{v})

	s.requestTransactionStorage(tx, Pool{})
	assert.Len(t, poolR.stores, 1)
	assert.Equal(t, tx.TransactionHash(), poolR.stores[0].TransactionHash())
}

/*
Scenario: Perform Proof of work
	Given a transaction and em shared keypair stored
	When I want to perform the proof of work of this transaction
	Then I get the valid public key
*/
func TestPerformPOW(t *testing.T) {

	emKey, _ := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	emPub, _ := x509.MarshalPKIXPublicKey(emKey.Public())
	emPv, _ := x509.MarshalECPrivateKey(emKey)

	propKey, _ := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	propPub, _ := x509.MarshalPKIXPublicKey(propKey.Public())

	txKey, _ := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	txKeyPub, _ := x509.MarshalPKIXPublicKey(txKey.Public())
	txKeyPv, _ := x509.MarshalECPrivateKey(txKey)

	sharedRepo := &mockSharedRepo{}

	emKP, _ := shared.NewEmitterKeyPair(hex.EncodeToString([]byte("pvKey")), hex.EncodeToString(emPub))
	sharedRepo.StoreSharedEmitterKeyPair(emKP)

	s := MiningService{
		sharedSrv: shared.NewService(sharedRepo),
	}

	sk, _ := shared.NewEmitterKeyPair(hex.EncodeToString([]byte("pvKey")), hex.EncodeToString(propPub))
	prop, _ := NewProposal(sk)

	data := map[string]string{
		"encrypted_aes_key":          hex.EncodeToString([]byte("aesKey")),
		"encrypted_address_by_robot": hex.EncodeToString([]byte("addr")),
		"encrypted_address_by_id":    hex.EncodeToString([]byte("addr")),
	}

	tx := Transaction{
		address:   crypto.HashString("addr"),
		data:      data,
		txType:    IDType,
		timestamp: time.Now(),
		pubKey:    hex.EncodeToString(txKeyPub),
		prop:      prop,
	}
	txBytesBeforeSig, _ := tx.MarshalBeforeSignature()
	emSig, _ := crypto.Sign(string(txBytesBeforeSig), hex.EncodeToString(emPv))
	txSig, _ := crypto.Sign(string(txBytesBeforeSig), hex.EncodeToString(txKeyPv))
	tx.emSig = emSig
	tx.sig = txSig
	txBytes, _ := tx.MarshalHash()
	txHash := crypto.HashBytes(txBytes)
	tx.txHash = txHash

	pow, err := s.performPow(tx)
	assert.Nil(t, err)
	assert.Equal(t, hex.EncodeToString(emPub), pow)
}

/*
Scenario: Pre-validate a transaction
	Given a transaction
	When I want to prevalidate this transaction
	Then I get the miner validation and the proof of work
*/
func TestPreValidateTransaction(t *testing.T) {
	minerKey, _ := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	minerPub, _ := x509.MarshalPKIXPublicKey(minerKey.Public())
	minerPv, _ := x509.MarshalECPrivateKey(minerKey)

	emKey, _ := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	emPub, _ := x509.MarshalPKIXPublicKey(emKey.Public())
	emPv, _ := x509.MarshalECPrivateKey(emKey)

	propKey, _ := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	propPub, _ := x509.MarshalPKIXPublicKey(propKey.Public())

	txKey, _ := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	txKeyPub, _ := x509.MarshalPKIXPublicKey(txKey.Public())
	txKeyPv, _ := x509.MarshalECPrivateKey(txKey)

	sharedRepo := &mockSharedRepo{}

	emKP, _ := shared.NewEmitterKeyPair(hex.EncodeToString([]byte("pvKey")), hex.EncodeToString(emPub))
	sharedRepo.StoreSharedEmitterKeyPair(emKP)

	s := MiningService{
		sharedSrv: shared.NewService(sharedRepo),
		minerPubK: hex.EncodeToString(minerPub),
		minerPvk:  hex.EncodeToString(minerPv),
	}

	sk, _ := shared.NewEmitterKeyPair(hex.EncodeToString([]byte("pvKey")), hex.EncodeToString(propPub))
	prop, _ := NewProposal(sk)

	data := map[string]string{
		"encrypted_aes_key":          hex.EncodeToString([]byte("aesKey")),
		"encrypted_address_by_robot": hex.EncodeToString([]byte("addr")),
		"encrypted_address_by_id":    hex.EncodeToString([]byte("addr")),
	}

	tx := Transaction{
		address:   crypto.HashString("addr"),
		data:      data,
		txType:    IDType,
		timestamp: time.Now(),
		pubKey:    hex.EncodeToString(txKeyPub),
		prop:      prop,
	}
	txBytesBeforeSig, _ := tx.MarshalBeforeSignature()
	emSig, _ := crypto.Sign(string(txBytesBeforeSig), hex.EncodeToString(emPv))
	txSig, _ := crypto.Sign(string(txBytesBeforeSig), hex.EncodeToString(txKeyPv))
	tx.emSig = emSig
	tx.sig = txSig
	txBytes, _ := tx.MarshalHash()
	txHash := crypto.HashBytes(txBytes)
	tx.txHash = txHash

	v, pow, err := s.preValidateTx(tx)
	assert.Nil(t, err)
	assert.Equal(t, hex.EncodeToString(emPub), pow)
	assert.Equal(t, hex.EncodeToString(minerPub), v.MinerPublicKey())
	assert.Equal(t, ValidationOK, v.Status())
	ok, err := v.IsValid()
	assert.True(t, ok)
	assert.Nil(t, err)
}

/*
Scenario: Mine transaction
	Given a valid transaction
	When I want to mine this transaction
	Then I get a master validation with the right POW, and validations confirmations
*/
func TestMineTransaction(t *testing.T) {
	minerKey, _ := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	minerPub, _ := x509.MarshalPKIXPublicKey(minerKey.Public())
	minerPv, _ := x509.MarshalECPrivateKey(minerKey)

	emKey, _ := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	emPub, _ := x509.MarshalPKIXPublicKey(emKey.Public())
	emPv, _ := x509.MarshalECPrivateKey(emKey)

	propKey, _ := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	propPub, _ := x509.MarshalPKIXPublicKey(propKey.Public())

	txKey, _ := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	txKeyPub, _ := x509.MarshalPKIXPublicKey(txKey.Public())
	txKeyPv, _ := x509.MarshalECPrivateKey(txKey)

	sharedRepo := &mockSharedRepo{}

	emKP, _ := shared.NewEmitterKeyPair(hex.EncodeToString([]byte("pvKey")), hex.EncodeToString(emPub))
	sharedRepo.StoreSharedEmitterKeyPair(emKP)

	s := MiningService{
		sharedSrv: shared.NewService(sharedRepo),
		minerPubK: hex.EncodeToString(minerPub),
		minerPvk:  hex.EncodeToString(minerPv),
		poolR:     &mockPoolRequester{},
	}

	sk, _ := shared.NewEmitterKeyPair(hex.EncodeToString([]byte("pvKey")), hex.EncodeToString(propPub))
	prop, _ := NewProposal(sk)

	data := map[string]string{
		"encrypted_aes_key":          hex.EncodeToString([]byte("aesKey")),
		"encrypted_address_by_robot": hex.EncodeToString([]byte("addr")),
		"encrypted_address_by_id":    hex.EncodeToString([]byte("addr")),
	}

	tx := Transaction{
		address:   crypto.HashString("addr"),
		data:      data,
		txType:    IDType,
		timestamp: time.Now(),
		pubKey:    hex.EncodeToString(txKeyPub),
		prop:      prop,
	}
	txBytesBeforeSig, _ := tx.MarshalBeforeSignature()
	emSig, _ := crypto.Sign(string(txBytesBeforeSig), hex.EncodeToString(emPv))
	txSig, _ := crypto.Sign(string(txBytesBeforeSig), hex.EncodeToString(txKeyPv))
	tx.emSig = emSig
	tx.sig = txSig
	txBytes, _ := tx.MarshalHash()
	txHash := crypto.HashBytes(txBytes)
	tx.txHash = txHash

	masterValid, confs, err := s.mineTransaction(tx, Pool{}, Pool{}, 1)
	assert.Nil(t, err)
	assert.Equal(t, hex.EncodeToString(emPub), masterValid.ProofOfWork())
	assert.Equal(t, ValidationOK, masterValid.Validation().Status())
	assert.Len(t, confs, 1)
}

/*
Scenario: Mine transaction where miner returns a KO validation
	Given a valid transaction
	When I want to mine this transaction and get a KO validation
	Then I get an error
*/
func TestMineTransactionWithKO(t *testing.T) {
	minerKey, _ := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	minerPub, _ := x509.MarshalPKIXPublicKey(minerKey.Public())
	minerPv, _ := x509.MarshalECPrivateKey(minerKey)

	emKey, _ := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	emPub, _ := x509.MarshalPKIXPublicKey(emKey.Public())
	emPv, _ := x509.MarshalECPrivateKey(emKey)

	propKey, _ := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	propPub, _ := x509.MarshalPKIXPublicKey(propKey.Public())

	txKey, _ := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	txKeyPub, _ := x509.MarshalPKIXPublicKey(txKey.Public())
	txKeyPv, _ := x509.MarshalECPrivateKey(txKey)

	sharedRepo := &mockSharedRepo{}

	emKP, _ := shared.NewEmitterKeyPair(hex.EncodeToString([]byte("pvKey")), hex.EncodeToString(emPub))
	sharedRepo.StoreSharedEmitterKeyPair(emKP)

	s := MiningService{
		sharedSrv: shared.NewService(sharedRepo),
		minerPubK: hex.EncodeToString(minerPub),
		minerPvk:  hex.EncodeToString(minerPv),
		poolR:     &mockPoolRequesterKO{},
	}

	sk, _ := shared.NewEmitterKeyPair(hex.EncodeToString([]byte("pvKey")), hex.EncodeToString(propPub))
	prop, _ := NewProposal(sk)

	data := map[string]string{
		"encrypted_aes_key":          hex.EncodeToString([]byte("aesKey")),
		"encrypted_address_by_robot": hex.EncodeToString([]byte("addr")),
		"encrypted_address_by_id":    hex.EncodeToString([]byte("addr")),
	}

	tx := Transaction{
		address:   crypto.HashString("addr"),
		data:      data,
		txType:    IDType,
		timestamp: time.Now(),
		pubKey:    hex.EncodeToString(txKeyPub),
		prop:      prop,
	}
	txBytesBeforeSig, _ := tx.MarshalBeforeSignature()
	emSig, _ := crypto.Sign(string(txBytesBeforeSig), hex.EncodeToString(emPv))
	txSig, _ := crypto.Sign(string(txBytesBeforeSig), hex.EncodeToString(txKeyPv))
	tx.emSig = emSig
	tx.sig = txSig
	txBytes, _ := tx.MarshalHash()
	txHash := crypto.HashBytes(txBytes)
	tx.txHash = txHash

	_, _, err := s.mineTransaction(tx, Pool{}, Pool{}, 1)
	assert.Equal(t, err, ErrInvalidTransaction)
}

/*
Scenario: Find pool for transaction mining
	Given a transaction
	When I want to find the pools
	Then I get the last validation pool, the validation pool and the storage pool
*/
func TestFindPools(t *testing.T) {
	s := MiningService{
		poolFSrv: PoolFindingService{
			pRetr: mockPoolRetriever{},
		},
	}

	lastVPool, validPool, storagePool, err := s.findPools(Transaction{
		address: "addr",
	})

	assert.Nil(t, err)
	assert.Equal(t, "127.0.0.1", lastVPool[0].IP().String())
	assert.Equal(t, "127.0.0.1", validPool[0].IP().String())
	assert.Equal(t, "127.0.0.1", storagePool[0].IP().String())
}

/*
Scenario: Lead the transaction validation
	Given an incoming transaction
	When I want to lead the mining
	Then after 1 second (because asynchronous), the transaction is stored
*/
func TestLeadTransactionValidation(t *testing.T) {
	minerKey, _ := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	minerPub, _ := x509.MarshalPKIXPublicKey(minerKey.Public())
	minerPv, _ := x509.MarshalECPrivateKey(minerKey)

	emKey, _ := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	emPub, _ := x509.MarshalPKIXPublicKey(emKey.Public())
	emPv, _ := x509.MarshalECPrivateKey(emKey)

	propKey, _ := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	propPub, _ := x509.MarshalPKIXPublicKey(propKey.Public())

	txKey, _ := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	txKeyPub, _ := x509.MarshalPKIXPublicKey(txKey.Public())
	txKeyPv, _ := x509.MarshalECPrivateKey(txKey)

	sharedRepo := &mockSharedRepo{}

	emKP, _ := shared.NewEmitterKeyPair(hex.EncodeToString([]byte("pvKey")), hex.EncodeToString(emPub))
	sharedRepo.StoreSharedEmitterKeyPair(emKP)

	poolR := &mockPoolRequester{}

	s := MiningService{
		sharedSrv: shared.NewService(sharedRepo),
		minerPubK: hex.EncodeToString(minerPub),
		minerPvk:  hex.EncodeToString(minerPv),
		poolR:     poolR,
		poolFSrv: PoolFindingService{
			pRetr: mockPoolRetriever{},
		},
	}

	sk, _ := shared.NewEmitterKeyPair(hex.EncodeToString([]byte("pvKey")), hex.EncodeToString(propPub))
	prop, _ := NewProposal(sk)

	data := map[string]string{
		"encrypted_aes_key":          hex.EncodeToString([]byte("aesKey")),
		"encrypted_address_by_robot": hex.EncodeToString([]byte("addr")),
		"encrypted_address_by_id":    hex.EncodeToString([]byte("addr")),
	}

	tx := Transaction{
		address:   crypto.HashString("addr"),
		data:      data,
		txType:    IDType,
		timestamp: time.Now(),
		pubKey:    hex.EncodeToString(txKeyPub),
		prop:      prop,
	}
	txBytesBeforeSig, _ := tx.MarshalBeforeSignature()
	emSig, _ := crypto.Sign(string(txBytesBeforeSig), hex.EncodeToString(emPv))
	txSig, _ := crypto.Sign(string(txBytesBeforeSig), hex.EncodeToString(txKeyPv))
	tx.emSig = emSig
	tx.sig = txSig
	txBytes, _ := tx.MarshalHash()
	txHash := crypto.HashBytes(txBytes)
	tx.txHash = txHash

	s.LeadTransactionValidation(tx, 1)

	time.Sleep(1 * time.Second)
	assert.Len(t, poolR.stores, 1)
	assert.Equal(t, tx.TransactionHash(), poolR.stores[0].TransactionHash())
}

type mockPoolRequester struct {
	stores []Transaction
}

func (pr mockPoolRequester) RequestTransactionLock(pool Pool, txLock Lock) error {
	return nil
}

func (pr mockPoolRequester) RequestTransactionUnlock(pool Pool, txLock Lock) error {
	return nil
}

func (pr mockPoolRequester) RequestTransactionValidations(pool Pool, tx Transaction, masterValid MasterValidation, validChan chan<- MinerValidation) {
	pub, pv := crypto.GenerateKeys()

	v := MinerValidation{
		minerPubk: pub,
		status:    ValidationOK,
		timestamp: time.Now(),
	}
	vBytes, _ := json.Marshal(v)
	sig, _ := crypto.Sign(string(vBytes), pv)
	v, _ = NewMinerValidation(v.status, v.timestamp, v.minerPubk, sig)

	validChan <- v
}

func (pr *mockPoolRequester) RequestTransactionStorage(pool Pool, tx Transaction, ackChan chan<- bool) {
	pr.stores = append(pr.stores, tx)
	ackChan <- true
}

type mockPoolRequesterKO struct {
	stores []Transaction
}

func (pr mockPoolRequesterKO) RequestTransactionLock(pool Pool, txLock Lock) error {
	return nil
}

func (pr mockPoolRequesterKO) RequestTransactionUnlock(pool Pool, txLock Lock) error {
	return nil
}

func (pr mockPoolRequesterKO) RequestTransactionValidations(pool Pool, tx Transaction, masterValid MasterValidation, validChan chan<- MinerValidation) {
	pub, pv := crypto.GenerateKeys()

	v := MinerValidation{
		minerPubk: pub,
		status:    ValidationKO,
		timestamp: time.Now(),
	}
	vBytes, _ := json.Marshal(v)
	sig, _ := crypto.Sign(string(vBytes), pv)
	v, _ = NewMinerValidation(v.status, v.timestamp, v.minerPubk, sig)

	validChan <- v
}

func (pr *mockPoolRequesterKO) RequestTransactionStorage(pool Pool, tx Transaction, ackChan chan<- bool) {
}

type mockSharedRepo struct {
	emKeys    shared.EmitterKeys
	minerKeys shared.MinerKeyPair
}

func (r mockSharedRepo) ListSharedEmitterKeyPairs() (shared.EmitterKeys, error) {
	return r.emKeys, nil
}
func (r *mockSharedRepo) StoreSharedEmitterKeyPair(kp shared.EmitterKeyPair) error {
	r.emKeys = append(r.emKeys, kp)
	return nil
}

func (r *mockSharedRepo) GetLastSharedMinersKeyPair() (shared.MinerKeyPair, error) {
	return r.minerKeys, nil
}
