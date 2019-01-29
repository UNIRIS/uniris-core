package transaction

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/x509"
	"encoding/hex"
	"encoding/json"
	"sort"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/uniris/uniris-core/pkg/crypto"
	"github.com/uniris/uniris-core/pkg/shared"
)

/*
Scenario: Get a keychain transaction by its type and hash
	Given a keychain type and an hash
	When I want to retrieve the transaction
	Then I can get this transaction
*/
func TestGetTransactionByHashAndTypeKeychain(t *testing.T) {
	repo := &mockTxRepository{}

	key, _ := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	pub, _ := x509.MarshalPKIXPublicKey(key.Public())
	pv, _ := x509.MarshalECPrivateKey(key)
	sig, _ := crypto.Sign("hello", hex.EncodeToString(pv))

	sk, _ := shared.NewKeyPair(hex.EncodeToString([]byte("pvKey")), hex.EncodeToString(pub))
	prop, _ := NewProposal(sk)

	data, _ := json.Marshal(keychainData{
		EncryptedAddress: hex.EncodeToString([]byte("addr")),
		EncryptedWallet:  hex.EncodeToString([]byte("wallet")),
	})

	tx, err := New(crypto.HashString("addr"), KeychainType, hex.EncodeToString(data), time.Now(), hex.EncodeToString(pub), sig, sig, prop, crypto.HashString("hash"))

	vBytes, _ := json.Marshal(MinerValidation{
		minerPubk: hex.EncodeToString(pub),
		status:    ValidationOK,
		timestamp: time.Now(),
	})
	vSig, _ := crypto.Sign(string(vBytes), hex.EncodeToString(pv))
	v, _ := NewMinerValidation(ValidationOK, time.Now(), hex.EncodeToString(pub), vSig)
	mv, _ := NewMasterValidation(Pool{}, hex.EncodeToString(pub), v)

	tx.AddMining(mv, []MinerValidation{v})
	keychain, err := NewKeychain(tx)

	assert.Nil(t, repo.StoreKeychain(keychain))

	s := StorageService{
		repo: repo,
	}
	txRes, err := s.getTransactionByHashAndType(tx.TransactionHash(), KeychainType)
	assert.Nil(t, err)
	assert.NotNil(t, txRes)
	assert.Equal(t, hex.EncodeToString(pub), txRes.PublicKey())
	assert.Equal(t, KeychainType, txRes.Type())
}

/*
Scenario: Get a ID transaction by type and its hash
	Given a ID type and an hash
	When I want to retrieve the transaction
	Then I can get this transaction
*/
func TestGetTransactionByHashAndTypeID(t *testing.T) {
	repo := &mockTxRepository{}

	key, _ := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	pub, _ := x509.MarshalPKIXPublicKey(key.Public())
	pv, _ := x509.MarshalECPrivateKey(key)
	sig, _ := crypto.Sign("hello", hex.EncodeToString(pv))

	sk, _ := shared.NewKeyPair(hex.EncodeToString([]byte("pvKey")), hex.EncodeToString(pub))
	prop, _ := NewProposal(sk)

	data, _ := json.Marshal(idData{
		EncryptedAddressByID:    hex.EncodeToString([]byte("addr")),
		EncryptedAddressByRobot: hex.EncodeToString([]byte("addr")),
		EncryptedAESKey:         hex.EncodeToString([]byte("aesKey")),
	})

	tx, err := New(crypto.HashString("addr"), IDType, hex.EncodeToString(data), time.Now(), hex.EncodeToString(pub), sig, sig, prop, crypto.HashString("hash"))

	vBytes, _ := json.Marshal(MinerValidation{
		minerPubk: hex.EncodeToString(pub),
		status:    ValidationOK,
		timestamp: time.Now(),
	})
	vSig, _ := crypto.Sign(string(vBytes), hex.EncodeToString(pv))
	v, _ := NewMinerValidation(ValidationOK, time.Now(), hex.EncodeToString(pub), vSig)
	mv, _ := NewMasterValidation(Pool{}, hex.EncodeToString(pub), v)

	tx.AddMining(mv, []MinerValidation{v})
	id, err := NewID(tx)

	assert.Nil(t, repo.StoreID(id))

	s := StorageService{
		repo: repo,
	}
	txRes, err := s.getTransactionByHashAndType(tx.TransactionHash(), IDType)
	assert.Nil(t, err)
	assert.NotNil(t, txRes)
	assert.Equal(t, hex.EncodeToString(pub), txRes.PublicKey())
	assert.Equal(t, IDType, txRes.Type())
}

/*
Scenario: Get a keychain unknown transaction
	Given a unknwown Keychain transaction hash
	When I want to retrieve the transaction
	Then I get an error
*/
func TestGetTransactionByHashAndTypeKeychainUnknown(t *testing.T) {
	s := StorageService{
		repo: &mockTxRepository{},
	}
	tx, err := s.getTransactionByHashAndType(crypto.HashString("txHash"), KeychainType)
	assert.Nil(t, tx)
	assert.Nil(t, err)
}

/*
Scenario: Get a keychain unknown transaction
	Given a unknwown Keychain transaction hash
	When I want to retrieve the transaction
	Then I get an error
*/
func TestGetTransactionByHashAndTypeIDUnknown(t *testing.T) {
	s := StorageService{
		repo: &mockTxRepository{},
	}
	tx, err := s.getTransactionByHashAndType(crypto.HashString("txHash"), IDType)
	assert.Nil(t, tx)
	assert.Nil(t, err)
}

/*
Scenario: Get transaction keychain by its hash
	Given a keychain tx stored
	When I want to retrieve the transaction by only its hash
	Then I can get it
*/
func TestGetKeychainTransactionByHash(t *testing.T) {
	repo := &mockTxRepository{}

	key, _ := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	pub, _ := x509.MarshalPKIXPublicKey(key.Public())
	pv, _ := x509.MarshalECPrivateKey(key)
	sig, _ := crypto.Sign("hello", hex.EncodeToString(pv))

	sk, _ := shared.NewKeyPair(hex.EncodeToString([]byte("pvKey")), hex.EncodeToString(pub))
	prop, _ := NewProposal(sk)

	data, _ := json.Marshal(keychainData{
		EncryptedAddress: hex.EncodeToString([]byte("addr")),
		EncryptedWallet:  hex.EncodeToString([]byte("wallet")),
	})

	tx, _ := New(crypto.HashString("addr"), KeychainType, hex.EncodeToString(data), time.Now(), hex.EncodeToString(pub), sig, sig, prop, crypto.HashString("hash"))

	vBytes, _ := json.Marshal(MinerValidation{
		minerPubk: hex.EncodeToString(pub),
		status:    ValidationOK,
		timestamp: time.Now(),
	})
	vSig, _ := crypto.Sign(string(vBytes), hex.EncodeToString(pv))
	v, _ := NewMinerValidation(ValidationOK, time.Now(), hex.EncodeToString(pub), vSig)
	mv, _ := NewMasterValidation(Pool{}, hex.EncodeToString(pub), v)

	tx.AddMining(mv, []MinerValidation{v})
	keychain, _ := NewKeychain(tx)

	assert.Nil(t, repo.StoreKeychain(keychain))

	s := StorageService{
		repo: repo,
	}
	txKeychain, err := s.getTransactionByHash(tx.TransactionHash())
	assert.Nil(t, err)
	assert.Equal(t, KeychainType, txKeychain.Type())
	assert.Equal(t, hex.EncodeToString(pub), txKeychain.PublicKey())
}

/*
Scenario: Get transaction ID by its hash
	Given a ID tx stored
	When I want to retrieve the transaction by only its hash
	Then I can get it
*/
func TestGetIDTransactionByHash(t *testing.T) {
	repo := &mockTxRepository{}

	key, _ := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	pub, _ := x509.MarshalPKIXPublicKey(key.Public())
	pv, _ := x509.MarshalECPrivateKey(key)
	sig, _ := crypto.Sign("hello", hex.EncodeToString(pv))

	sk, _ := shared.NewKeyPair(hex.EncodeToString([]byte("pvKey")), hex.EncodeToString(pub))
	prop, _ := NewProposal(sk)

	data, _ := json.Marshal(idData{
		EncryptedAESKey:         hex.EncodeToString([]byte("aeskey")),
		EncryptedAddressByID:    hex.EncodeToString([]byte("addr")),
		EncryptedAddressByRobot: hex.EncodeToString([]byte("addr")),
	})

	tx, _ := New(crypto.HashString("addr"), IDType, hex.EncodeToString(data), time.Now(), hex.EncodeToString(pub), sig, sig, prop, crypto.HashString("hash"))

	vBytes, _ := json.Marshal(MinerValidation{
		minerPubk: hex.EncodeToString(pub),
		status:    ValidationOK,
		timestamp: time.Now(),
	})
	vSig, _ := crypto.Sign(string(vBytes), hex.EncodeToString(pv))
	v, _ := NewMinerValidation(ValidationOK, time.Now(), hex.EncodeToString(pub), vSig)
	mv, _ := NewMasterValidation(Pool{}, hex.EncodeToString(pub), v)

	tx.AddMining(mv, []MinerValidation{v})
	id, _ := NewID(tx)

	assert.Nil(t, repo.StoreID(id))

	s := StorageService{
		repo: repo,
	}
	txID, err := s.getTransactionByHash(tx.TransactionHash())
	assert.Nil(t, err)
	assert.Equal(t, IDType, txID.Type())
	assert.Equal(t, hex.EncodeToString(pub), txID.PublicKey())
}

/*
Scenario: Get unknown transaction by its hash
	Given no tx stored
	When I want to retrieve the transaction by only its hash
	Then I can get an error
*/
func TestGetUnknownTransactionByHash(t *testing.T) {
	s := StorageService{
		repo: &mockTxRepository{},
	}

	_, err := s.getTransactionByHash(crypto.HashString("hash"))
	assert.EqualError(t, err, "unknown transaction")
}

/*
Scenario: Get transaction status pending
	Given a transaction stored in pending
	When I want to get its status
	Then I get pending
*/
func TestGetTransactionStatusPending(t *testing.T) {
	repo := &mockTxRepository{
		pendings: []Transaction{
			Transaction{
				txHash: crypto.HashString("hash"),
			},
		},
	}

	s := StorageService{
		repo: repo,
	}

	status, err := s.GetTransactionStatus(crypto.HashString("hash"))
	assert.Nil(t, err)
	assert.Equal(t, StatusPending, status)
}

/*
Scenario: Get transaction status KO
	Given a transaction stored in KO
	When I want to get its status
	Then I get failure
*/
func TestGetTransactionStatusFailure(t *testing.T) {
	repo := &mockTxRepository{
		kos: []Transaction{
			Transaction{
				txHash: crypto.HashString("hash"),
			},
		},
	}

	s := StorageService{
		repo: repo,
	}

	status, err := s.GetTransactionStatus(crypto.HashString("hash"))
	assert.Nil(t, err)
	assert.Equal(t, StatusFailure, status)
}

/*
Scenario: Get transaction status unknown
	Given a transaction stored in KO
	When I want to get its status
	Then I get failure
*/
func TestGetTransactionStatusUnknown(t *testing.T) {
	repo := &mockTxRepository{}

	s := StorageService{
		repo: repo,
	}

	status, err := s.GetTransactionStatus(crypto.HashString("hash"))
	assert.Nil(t, err)
	assert.Equal(t, StatusUnknown, status)
}

/*
Scenario: Get transaction status success
	Given a transaction stored
	When I want to get its status
	Then I get success
*/
func TestGetTransactionStatusSuccess(t *testing.T) {
	repo := &mockTxRepository{}

	key, _ := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	pub, _ := x509.MarshalPKIXPublicKey(key.Public())
	pv, _ := x509.MarshalECPrivateKey(key)
	sig, _ := crypto.Sign("hello", hex.EncodeToString(pv))

	sk, _ := shared.NewKeyPair(hex.EncodeToString([]byte("pvKey")), hex.EncodeToString(pub))
	prop, _ := NewProposal(sk)

	data, _ := json.Marshal(idData{
		EncryptedAESKey:         hex.EncodeToString([]byte("aeskey")),
		EncryptedAddressByID:    hex.EncodeToString([]byte("addr")),
		EncryptedAddressByRobot: hex.EncodeToString([]byte("addr")),
	})

	tx, _ := New(crypto.HashString("addr"), IDType, hex.EncodeToString(data), time.Now(), hex.EncodeToString(pub), sig, sig, prop, crypto.HashString("hash"))

	vBytes, _ := json.Marshal(MinerValidation{
		minerPubk: hex.EncodeToString(pub),
		status:    ValidationOK,
		timestamp: time.Now(),
	})
	vSig, _ := crypto.Sign(string(vBytes), hex.EncodeToString(pv))
	v, _ := NewMinerValidation(ValidationOK, time.Now(), hex.EncodeToString(pub), vSig)
	mv, _ := NewMasterValidation(Pool{}, hex.EncodeToString(pub), v)

	tx.AddMining(mv, []MinerValidation{v})
	id, _ := NewID(tx)

	assert.Nil(t, repo.StoreID(id))

	s := StorageService{
		repo: repo,
	}

	status, err := s.GetTransactionStatus(tx.TransactionHash())
	assert.Nil(t, err)
	assert.Equal(t, StatusSuccess, status)
}

/*
Scenario: Get transaction status with invalid hash
	Given an invalid hash
	When I want to get its status
	Then I get success
*/
func TestGetTransactionStatusInvalidHash(t *testing.T) {
	s := StorageService{}
	_, err := s.GetTransactionStatus(hex.EncodeToString([]byte("hash")))
	assert.EqualError(t, err, "get transaction status: hash is not valid")
}

/*
Scenario: Check if the miner is authorized to store the transaction
	Given a transaction hash
	When I want to check if I can store this transaction
	Then I get a true
	//TODO: to improve when the implementation will be defined
*/
func TestIsAuthorizedToStore(t *testing.T) {
	s := StorageService{}
	assert.True(t, s.isAuthorizedToStoreTx(""))
}

/*
Scenario: Get the last keychain transaction
	Given two keychain transaction chained
	When I want to get the last
	Then I get the 2nd
*/
func TestGetLastKeychainTransaction(t *testing.T) {
	repo := &mockTxRepository{}
	s := StorageService{
		repo: repo,
	}

	key, _ := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	pub, _ := x509.MarshalPKIXPublicKey(key.Public())
	pv, _ := x509.MarshalECPrivateKey(key)
	sig, _ := crypto.Sign("hello", hex.EncodeToString(pv))

	sk, _ := shared.NewKeyPair(hex.EncodeToString([]byte("pvKey")), hex.EncodeToString(pub))
	prop, _ := NewProposal(sk)

	data, _ := json.Marshal(keychainData{
		EncryptedAddress: hex.EncodeToString([]byte("addr")),
		EncryptedWallet:  hex.EncodeToString([]byte("wallet")),
	})

	tx1, _ := New(crypto.HashString("addr"), KeychainType, hex.EncodeToString(data), time.Now(), hex.EncodeToString(pub), sig, sig, prop, crypto.HashString("hash"))

	vBytes, _ := json.Marshal(MinerValidation{
		minerPubk: hex.EncodeToString(pub),
		status:    ValidationOK,
		timestamp: time.Now(),
	})
	vSig, _ := crypto.Sign(string(vBytes), hex.EncodeToString(pv))
	v, _ := NewMinerValidation(ValidationOK, time.Now(), hex.EncodeToString(pub), vSig)
	mv, _ := NewMasterValidation(Pool{}, hex.EncodeToString(pub), v)

	tx1.AddMining(mv, []MinerValidation{v})
	keychain1, _ := NewKeychain(tx1)

	assert.Nil(t, repo.StoreKeychain(keychain1))

	time.Sleep(1 * time.Second)

	tx2, err := New(crypto.HashString("addr"), KeychainType, hex.EncodeToString(data), time.Now(), hex.EncodeToString(pub), sig, sig, prop, crypto.HashString("hash2"))
	tx2.AddMining(mv, []MinerValidation{v})
	keychain2, err := NewKeychain(tx2)

	assert.Nil(t, repo.StoreKeychain(keychain2))
	assert.Len(t, repo.keychains, 2)

	lastTx, err := s.GetLastTransaction(crypto.HashString("addr"), KeychainType)
	assert.Nil(t, err)
	assert.NotNil(t, lastTx)
	assert.Equal(t, KeychainType, lastTx.Type())
	assert.Equal(t, crypto.HashString("hash2"), lastTx.TransactionHash())

}

/*
Scenario: Get the last ID transaction
	Given two ID transaction
	When I want to get the last
	Then I get the one I reached (because ID are not chained)
*/
func TestGetLastIDTransaction(t *testing.T) {
	repo := &mockTxRepository{}
	s := StorageService{
		repo: repo,
	}

	key, _ := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	pub, _ := x509.MarshalPKIXPublicKey(key.Public())
	pv, _ := x509.MarshalECPrivateKey(key)
	sig, _ := crypto.Sign("hello", hex.EncodeToString(pv))

	sk, _ := shared.NewKeyPair(hex.EncodeToString([]byte("pvKey")), hex.EncodeToString(pub))
	prop, _ := NewProposal(sk)

	data, _ := json.Marshal(idData{
		EncryptedAESKey:         hex.EncodeToString([]byte("aes")),
		EncryptedAddressByID:    hex.EncodeToString([]byte("addr")),
		EncryptedAddressByRobot: hex.EncodeToString([]byte("addr")),
	})

	tx1, _ := New(crypto.HashString("addr"), IDType, hex.EncodeToString(data), time.Now(), hex.EncodeToString(pub), sig, sig, prop, crypto.HashString("hash"))

	vBytes, _ := json.Marshal(MinerValidation{
		minerPubk: hex.EncodeToString(pub),
		status:    ValidationOK,
		timestamp: time.Now(),
	})
	vSig, _ := crypto.Sign(string(vBytes), hex.EncodeToString(pv))
	v, _ := NewMinerValidation(ValidationOK, time.Now(), hex.EncodeToString(pub), vSig)
	mv, _ := NewMasterValidation(Pool{}, hex.EncodeToString(pub), v)

	tx1.AddMining(mv, []MinerValidation{v})
	id1, _ := NewID(tx1)

	assert.Nil(t, repo.StoreID(id1))

	time.Sleep(1 * time.Second)

	tx2, err := New(crypto.HashString("addr2"), IDType, hex.EncodeToString(data), time.Now(), hex.EncodeToString(pub), sig, sig, prop, crypto.HashString("hash2"))
	tx2.AddMining(mv, []MinerValidation{v})
	id2, err := NewID(tx2)

	assert.Nil(t, repo.StoreID(id2))
	assert.Len(t, repo.ids, 2)

	lastTx, err := s.GetLastTransaction(crypto.HashString("addr"), IDType)
	assert.Nil(t, err)
	assert.NotNil(t, lastTx)
	assert.Equal(t, IDType, lastTx.Type())
	assert.Equal(t, crypto.HashString("hash"), lastTx.TransactionHash())

	lastTx, err = s.GetLastTransaction(crypto.HashString("addr2"), IDType)
	assert.Nil(t, err)
	assert.NotNil(t, lastTx)
	assert.Equal(t, IDType, lastTx.Type())
	assert.Equal(t, crypto.HashString("hash2"), lastTx.TransactionHash())

}

/*
Scenario: Get the chain of a transaction
	Given transactions chained
	When I want to get the chain
	Then I get recursively the transactions linked to each other
*/
func TestGetTransactionChain(t *testing.T) {
	repo := &mockTxRepository{}
	s := StorageService{
		repo: repo,
	}

	key, _ := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	pub, _ := x509.MarshalPKIXPublicKey(key.Public())
	pv, _ := x509.MarshalECPrivateKey(key)

	sk, _ := shared.NewKeyPair(hex.EncodeToString([]byte("pvKey")), hex.EncodeToString(pub))
	prop, _ := NewProposal(sk)

	data, _ := json.Marshal(keychainData{
		EncryptedAddress: hex.EncodeToString([]byte("addr")),
		EncryptedWallet:  hex.EncodeToString([]byte("wallet")),
	})

	tx1Bytes, _ := json.Marshal(Transaction{
		address:   crypto.HashString("addr"),
		txType:    KeychainType,
		data:      hex.EncodeToString(data),
		timestamp: time.Now(),
		pubKey:    hex.EncodeToString(pub),
		prop:      prop,
	})
	sig1, _ := crypto.Sign(string(tx1Bytes), hex.EncodeToString(pv))

	tx1, _ := New(crypto.HashString("addr"), KeychainType, hex.EncodeToString(data), time.Now(), hex.EncodeToString(pub), sig1, sig1, prop, crypto.HashBytes(tx1Bytes))

	vBytes, _ := json.Marshal(MinerValidation{
		minerPubk: hex.EncodeToString(pub),
		status:    ValidationOK,
		timestamp: time.Now(),
	})
	vSig, _ := crypto.Sign(string(vBytes), hex.EncodeToString(pv))
	v, _ := NewMinerValidation(ValidationOK, time.Now(), hex.EncodeToString(pub), vSig)
	mv, _ := NewMasterValidation(Pool{}, hex.EncodeToString(pub), v)

	tx1.AddMining(mv, []MinerValidation{v})
	keychain1, _ := NewKeychain(tx1)

	assert.Nil(t, repo.StoreKeychain(keychain1))

	time.Sleep(1 * time.Second)

	tx2Bytes, _ := json.Marshal(Transaction{
		address:   crypto.HashString("addr"),
		txType:    KeychainType,
		data:      hex.EncodeToString(data),
		timestamp: time.Now(),
		pubKey:    hex.EncodeToString(pub),
		prop:      prop,
	})
	sig2, _ := crypto.Sign(string(tx2Bytes), hex.EncodeToString(pv))

	tx2, _ := New(crypto.HashString("addr"), KeychainType, hex.EncodeToString(data), time.Now(), hex.EncodeToString(pub), sig2, sig2, prop, crypto.HashBytes(tx2Bytes))
	tx2.AddMining(mv, []MinerValidation{v})
	keychain2, _ := NewKeychain(tx2)
	assert.Nil(t, keychain2.Chain(&tx1))

	assert.Nil(t, repo.StoreKeychain(keychain2))
	assert.Len(t, repo.keychains, 2)

	chain, err := s.GetTransactionChain(crypto.HashString("addr"), KeychainType)
	assert.Nil(t, err)
	assert.NotNil(t, chain)
	assert.Equal(t, crypto.HashBytes(tx2Bytes), chain.TransactionHash())
	assert.NotNil(t, chain.PreviousTransaction())
	assert.Equal(t, crypto.HashBytes(tx1Bytes), chain.PreviousTransaction().TransactionHash())
}

/*
Scenario: Check a valid transaction before store
	Given a transaction
	When I want to check before storage
	Then I get not error
*/
func TestCheckTransactionBeforeStore(t *testing.T) {
	key, _ := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	pub, _ := x509.MarshalPKIXPublicKey(key.Public())
	pv, _ := x509.MarshalECPrivateKey(key)

	sk, _ := shared.NewKeyPair(hex.EncodeToString([]byte("pvKey")), hex.EncodeToString(pub))
	prop, _ := NewProposal(sk)

	data, _ := json.Marshal(keychainData{
		EncryptedAddress: hex.EncodeToString([]byte("addr")),
		EncryptedWallet:  hex.EncodeToString([]byte("wallet")),
	})

	txBytes, _ := json.Marshal(Transaction{
		address:   crypto.HashString("addr"),
		txType:    KeychainType,
		data:      hex.EncodeToString(data),
		timestamp: time.Now(),
		pubKey:    hex.EncodeToString(pub),
		prop:      prop,
	})
	sig, _ := crypto.Sign(string(txBytes), hex.EncodeToString(pv))

	tx, _ := New(crypto.HashString("addr"), KeychainType, hex.EncodeToString(data), time.Now(), hex.EncodeToString(pub), sig, sig, prop, crypto.HashBytes(txBytes))

	vBytes, _ := json.Marshal(MinerValidation{
		minerPubk: hex.EncodeToString(pub),
		status:    ValidationOK,
		timestamp: time.Now(),
	})
	vSig, _ := crypto.Sign(string(vBytes), hex.EncodeToString(pv))
	v, _ := NewMinerValidation(ValidationOK, time.Now(), hex.EncodeToString(pub), vSig)
	mv, _ := NewMasterValidation(Pool{}, hex.EncodeToString(pub), v)
	tx.AddMining(mv, []MinerValidation{v})

	s := StorageService{}

	assert.Nil(t, s.checkTransactionBeforeStorage(tx))
}

/*
Scenario: Check a transaction before store with misssing validations
	Given a transaction with missing confirmations validations
	When I want to check before storage
	Then I get an error
*/
func TestCheckTransactionBeforeStoreWithMissingValidations(t *testing.T) {
	key, _ := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	pub, _ := x509.MarshalPKIXPublicKey(key.Public())
	pv, _ := x509.MarshalECPrivateKey(key)

	sk, _ := shared.NewKeyPair(hex.EncodeToString([]byte("pvKey")), hex.EncodeToString(pub))
	prop, _ := NewProposal(sk)

	data, _ := json.Marshal(keychainData{
		EncryptedAddress: hex.EncodeToString([]byte("addr")),
		EncryptedWallet:  hex.EncodeToString([]byte("wallet")),
	})

	txBytes, _ := json.Marshal(Transaction{
		address:   crypto.HashString("addr"),
		txType:    KeychainType,
		data:      hex.EncodeToString(data),
		timestamp: time.Now(),
		pubKey:    hex.EncodeToString(pub),
		prop:      prop,
	})
	sig, _ := crypto.Sign(string(txBytes), hex.EncodeToString(pv))

	tx, _ := New(crypto.HashString("addr"), KeychainType, hex.EncodeToString(data), time.Now(), hex.EncodeToString(pub), sig, sig, prop, crypto.HashBytes(txBytes))

	vBytes, _ := json.Marshal(MinerValidation{
		minerPubk: hex.EncodeToString(pub),
		status:    ValidationOK,
		timestamp: time.Now(),
	})
	vSig, _ := crypto.Sign(string(vBytes), hex.EncodeToString(pv))
	v, _ := NewMinerValidation(ValidationOK, time.Now(), hex.EncodeToString(pub), vSig)
	mv, _ := NewMasterValidation(Pool{}, hex.EncodeToString(pub), v)
	tx.AddMining(mv, []MinerValidation{})

	s := StorageService{}

	assert.EqualError(t, s.checkTransactionBeforeStorage(tx), "transaction: invalid number of validations")
}

/*
Scenario: Store a KO transaction
	Given a transaction not valid
	When I want to store it
	Then the transaction is stored on the KO db
*/
func TestStoreKOTransaction(t *testing.T) {
	repo := &mockTxRepository{}
	s := StorageService{repo: repo}

	key, _ := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	pub, _ := x509.MarshalPKIXPublicKey(key.Public())
	pv, _ := x509.MarshalECPrivateKey(key)

	vBytes, _ := json.Marshal(MinerValidation{
		minerPubk: hex.EncodeToString(pub),
		status:    ValidationKO,
		timestamp: time.Now(),
	})
	vSig, _ := crypto.Sign(string(vBytes), hex.EncodeToString(pv))
	v, _ := NewMinerValidation(ValidationKO, time.Now(), hex.EncodeToString(pub), vSig)

	vBytes2, _ := json.Marshal(MinerValidation{
		minerPubk: hex.EncodeToString(pub),
		status:    ValidationOK,
		timestamp: time.Now(),
	})
	vSig2, _ := crypto.Sign(string(vBytes2), hex.EncodeToString(pv))
	v2, _ := NewMinerValidation(ValidationOK, time.Now(), hex.EncodeToString(pub), vSig2)
	mv, _ := NewMasterValidation(Pool{}, hex.EncodeToString(pub), v)

	sk, _ := shared.NewKeyPair(hex.EncodeToString([]byte("pvKey")), hex.EncodeToString(pub))
	prop, _ := NewProposal(sk)

	data, _ := json.Marshal(keychainData{
		EncryptedAddress: hex.EncodeToString([]byte("addr")),
		EncryptedWallet:  hex.EncodeToString([]byte("wallet")),
	})

	txBytes, _ := json.Marshal(Transaction{
		address:   crypto.HashString("addr"),
		txType:    KeychainType,
		data:      hex.EncodeToString(data),
		timestamp: time.Now(),
		pubKey:    hex.EncodeToString(pub),
		prop:      prop,
	})
	sig, _ := crypto.Sign(string(txBytes), hex.EncodeToString(pv))

	tx1, _ := New(crypto.HashString("addr"), KeychainType, hex.EncodeToString(data), time.Now(), hex.EncodeToString(pub), sig, sig, prop, crypto.HashBytes(txBytes))
	tx1.AddMining(mv, []MinerValidation{v2})

	assert.Nil(t, s.StoreTransaction(tx1))
	assert.Len(t, repo.kos, 1)
	assert.Equal(t, crypto.HashBytes(txBytes), repo.kos[0].txHash)
}

/*
Scenario: Store a Keychain transaction
	Given a keychain transaction
	When I want to store it
	Then the transaction is stored on the keychain db
*/
func TestStoreKeychainTransaction(t *testing.T) {
	repo := &mockTxRepository{}
	s := StorageService{repo: repo}

	key, _ := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	pub, _ := x509.MarshalPKIXPublicKey(key.Public())
	pv, _ := x509.MarshalECPrivateKey(key)

	vBytes, _ := json.Marshal(MinerValidation{
		minerPubk: hex.EncodeToString(pub),
		status:    ValidationOK,
		timestamp: time.Now(),
	})
	vSig, _ := crypto.Sign(string(vBytes), hex.EncodeToString(pv))
	v, _ := NewMinerValidation(ValidationOK, time.Now(), hex.EncodeToString(pub), vSig)

	vBytes2, _ := json.Marshal(MinerValidation{
		minerPubk: hex.EncodeToString(pub),
		status:    ValidationOK,
		timestamp: time.Now(),
	})
	vSig2, _ := crypto.Sign(string(vBytes2), hex.EncodeToString(pv))
	v2, _ := NewMinerValidation(ValidationOK, time.Now(), hex.EncodeToString(pub), vSig2)
	mv, _ := NewMasterValidation(Pool{}, hex.EncodeToString(pub), v)

	sk, _ := shared.NewKeyPair(hex.EncodeToString([]byte("pvKey")), hex.EncodeToString(pub))
	prop, _ := NewProposal(sk)

	data, _ := json.Marshal(keychainData{
		EncryptedAddress: hex.EncodeToString([]byte("addr")),
		EncryptedWallet:  hex.EncodeToString([]byte("wallet")),
	})

	txBytes, _ := json.Marshal(Transaction{
		address:   crypto.HashString("addr"),
		txType:    KeychainType,
		data:      hex.EncodeToString(data),
		timestamp: time.Now(),
		pubKey:    hex.EncodeToString(pub),
		prop:      prop,
	})
	sig, _ := crypto.Sign(string(txBytes), hex.EncodeToString(pv))

	tx1, _ := New(crypto.HashString("addr"), KeychainType, hex.EncodeToString(data), time.Now(), hex.EncodeToString(pub), sig, sig, prop, crypto.HashBytes(txBytes))
	tx1.AddMining(mv, []MinerValidation{v2})

	assert.Nil(t, s.StoreTransaction(tx1))
	assert.Len(t, repo.keychains, 1)
	assert.Equal(t, crypto.HashBytes(txBytes), repo.keychains[0].txHash)
}

/*
Scenario: Store a ID transaction
	Given a ID transaction
	When I want to store it
	Then the transaction is stored on the ID db
*/
func TestStoreIDTransaction(t *testing.T) {
	repo := &mockTxRepository{}
	s := StorageService{repo: repo}

	key, _ := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	pub, _ := x509.MarshalPKIXPublicKey(key.Public())
	pv, _ := x509.MarshalECPrivateKey(key)

	vBytes, _ := json.Marshal(MinerValidation{
		minerPubk: hex.EncodeToString(pub),
		status:    ValidationOK,
		timestamp: time.Now(),
	})
	vSig, _ := crypto.Sign(string(vBytes), hex.EncodeToString(pv))
	v, _ := NewMinerValidation(ValidationOK, time.Now(), hex.EncodeToString(pub), vSig)

	vBytes2, _ := json.Marshal(MinerValidation{
		minerPubk: hex.EncodeToString(pub),
		status:    ValidationOK,
		timestamp: time.Now(),
	})
	vSig2, _ := crypto.Sign(string(vBytes2), hex.EncodeToString(pv))
	v2, _ := NewMinerValidation(ValidationOK, time.Now(), hex.EncodeToString(pub), vSig2)
	mv, _ := NewMasterValidation(Pool{}, hex.EncodeToString(pub), v)

	sk, _ := shared.NewKeyPair(hex.EncodeToString([]byte("pvKey")), hex.EncodeToString(pub))
	prop, _ := NewProposal(sk)

	data, _ := json.Marshal(idData{
		EncryptedAddressByID:    hex.EncodeToString([]byte("addr")),
		EncryptedAddressByRobot: hex.EncodeToString([]byte("addr")),
		EncryptedAESKey:         hex.EncodeToString([]byte("aesKey")),
	})

	txBytes, _ := json.Marshal(Transaction{
		address:   crypto.HashString("addr"),
		txType:    IDType,
		data:      hex.EncodeToString(data),
		timestamp: time.Now(),
		pubKey:    hex.EncodeToString(pub),
		prop:      prop,
	})
	sig, _ := crypto.Sign(string(txBytes), hex.EncodeToString(pv))

	tx1, _ := New(crypto.HashString("addr"), IDType, hex.EncodeToString(data), time.Now(), hex.EncodeToString(pub), sig, sig, prop, crypto.HashBytes(txBytes))
	tx1.AddMining(mv, []MinerValidation{v2})

	assert.Nil(t, s.StoreTransaction(tx1))
	assert.Len(t, repo.ids, 1)
	assert.Equal(t, crypto.HashBytes(txBytes), repo.ids[0].txHash)
}

type mockTxRepository struct {
	pendings  []Transaction
	kos       []Transaction
	keychains []Keychain
	ids       []ID
}

func (r mockTxRepository) FindPendingTransaction(txHash string) (*Transaction, error) {
	for _, tx := range r.pendings {
		if tx.TransactionHash() == txHash {
			return &tx, nil
		}
	}
	return nil, nil
}

func (r mockTxRepository) GetKeychain(txAddr string) (*Keychain, error) {
	sort.Slice(r.keychains, func(i, j int) bool {
		return r.keychains[i].Timestamp().Unix() > r.keychains[j].Timestamp().Unix()
	})

	if len(r.keychains) > 0 {
		return &r.keychains[0], nil
	}
	return nil, nil
}

func (r mockTxRepository) FindKeychainByHash(txHash string) (*Keychain, error) {
	for _, tx := range r.keychains {
		if tx.TransactionHash() == txHash {
			return &tx, nil
		}
	}
	return nil, nil
}

func (r mockTxRepository) FindLastKeychain(addr string) (*Keychain, error) {

	sort.Slice(r.keychains, func(i, j int) bool {
		return r.keychains[i].Timestamp().Unix() > r.keychains[j].Timestamp().Unix()
	})

	for _, tx := range r.keychains {
		if tx.Address() == addr {
			return &tx, nil
		}
	}
	return nil, nil
}

func (r mockTxRepository) FindIDByHash(txHash string) (*ID, error) {
	for _, tx := range r.ids {
		if tx.TransactionHash() == txHash {
			return &tx, nil
		}
	}
	return nil, nil
}

func (r mockTxRepository) FindIDByAddress(addr string) (*ID, error) {
	for _, tx := range r.ids {
		if tx.Address() == addr {
			return &tx, nil
		}
	}
	return nil, nil
}

func (r mockTxRepository) FindKOTransaction(txHash string) (*Transaction, error) {
	for _, tx := range r.kos {
		if tx.TransactionHash() == txHash {
			return &tx, nil
		}
	}
	return nil, nil
}

func (r *mockTxRepository) StoreKeychain(kc Keychain) error {
	r.keychains = append(r.keychains, kc)
	return nil
}

func (r *mockTxRepository) StoreID(id ID) error {
	r.ids = append(r.ids, id)
	return nil
}

func (r *mockTxRepository) StoreKO(tx Transaction) error {
	r.kos = append(r.kos, tx)
	return nil
}
