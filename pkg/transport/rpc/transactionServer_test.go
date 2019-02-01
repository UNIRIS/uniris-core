package rpc

import (
	"context"
	"encoding/hex"
	"encoding/json"
	"log"
	"sort"
	"testing"
	"time"

	"google.golang.org/grpc/codes"

	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc/status"

	api "github.com/uniris/uniris-core/api/protobuf-spec"
	"github.com/uniris/uniris-core/pkg/shared"

	"github.com/uniris/uniris-core/pkg/crypto"
	"github.com/uniris/uniris-core/pkg/transaction"
)

/*
Scenario: Receive  get last transction about an unknown transaction
	Given no transaction store for an address
	When I want to request to retrieve the last transaction keychain of this unknown address
	Then I get an error
*/
func TestHandleGetLastTransactionWhenNotExist(t *testing.T) {
	pub, pv := crypto.GenerateKeys()

	lockRepo := &mockLockRepository{}
	txRepo := &mockTxRepository{}
	sharedRepo := &mockSharedRepo{}
	txSrv := newTransactionServer(txRepo, lockRepo, sharedRepo, pub, pv)

	req := &api.LastTransactionRequest{
		Timestamp:          time.Now().Unix(),
		TransactionAddress: crypto.HashString("address"),
		Type:               api.TransactionType_KEYCHAIN,
	}
	reqBytes, _ := json.Marshal(req)
	sig, _ := crypto.Sign(string(reqBytes), pv)
	req.SignatureRequest = sig

	_, err := txSrv.GetLastTransaction(context.TODO(), req)
	assert.NotNil(t, err)
	statusCode, _ := status.FromError(err)
	assert.Equal(t, codes.NotFound, statusCode.Code())
	assert.Equal(t, statusCode.Message(), "transaction does not exist")
}

/*
Scenario: Receive  get last transaction request
	Given a keychain transaction stored
	When I want to request to retrieve the last transaction keychain of this address
	Then I get an error
*/
func TestHandleGetLastTransaction(t *testing.T) {

	pub, pv := crypto.GenerateKeys()

	lockRepo := &mockLockRepository{}
	txRepo := &mockTxRepository{}
	sharedRepo := &mockSharedRepo{}
	txSrv := newTransactionServer(txRepo, lockRepo, sharedRepo, pub, pv)

	data := map[string]string{
		"encrypted_address": hex.EncodeToString([]byte("addr")),
		"encrypted_wallet":  hex.EncodeToString([]byte("wallet")),
	}

	propKP, _ := shared.NewEmitterKeyPair(hex.EncodeToString([]byte("pvkey")), pub)
	prop, _ := transaction.NewProposal(propKP)

	txRaw := map[string]interface{}{
		"address":    crypto.HashString("addr"),
		"data":       data,
		"timestamp":  time.Now().Unix(),
		"type":       transaction.KeychainType,
		"public_key": pub,
		"proposal":   prop,
	}
	txBytes, _ := json.Marshal(txRaw)
	sig, _ := crypto.Sign(string(txBytes), pv)
	txRaw["signature"] = sig
	txRaw["em_signature"] = sig
	txBytes, _ = json.Marshal(txRaw)

	tx, _ := transaction.New(crypto.HashString("addr"), transaction.KeychainType, data, time.Now(), pub, sig, sig, prop, crypto.HashBytes(txBytes))
	keychain, _ := transaction.NewKeychain(tx)
	txRepo.StoreKeychain(keychain)

	req := &api.LastTransactionRequest{
		Timestamp:          time.Now().Unix(),
		TransactionAddress: crypto.HashString("addr"),
		Type:               api.TransactionType_KEYCHAIN,
	}
	reqBytes, _ := json.Marshal(req)
	sigReq, _ := crypto.Sign(string(reqBytes), pv)
	req.SignatureRequest = sigReq

	res, err := txSrv.GetLastTransaction(context.TODO(), req)
	assert.Nil(t, err)
	assert.NotEmpty(t, res.SignatureResponse)
	assert.NotNil(t, res.Transaction)
	assert.Equal(t, crypto.HashBytes(txBytes), res.Transaction.TransactionHash)

	resBytes, _ := json.Marshal(&api.LastTransactionResponse{
		Timestamp:   res.Timestamp,
		Transaction: res.Transaction,
	})
	assert.Nil(t, crypto.VerifySignature(string(resBytes), pub, res.SignatureResponse))
}

/*
Scenario: Receive get transaction status request
	Given no transaction stored
	When I want to request the transactions status for this transaction hash
	Then I get a status unknown
*/
func TestHandleGetTransactionStatus(t *testing.T) {

	pub, pv := crypto.GenerateKeys()

	lockRepo := &mockLockRepository{}
	txRepo := &mockTxRepository{}
	sharedRepo := &mockSharedRepo{}
	txSrv := newTransactionServer(txRepo, lockRepo, sharedRepo, pub, pv)

	req := &api.TransactionStatusRequest{
		Timestamp:       time.Now().Unix(),
		TransactionHash: crypto.HashString("tx"),
	}
	reqBytes, _ := json.Marshal(req)
	sig, _ := crypto.Sign(string(reqBytes), pv)
	req.SignatureRequest = sig

	res, err := txSrv.GetTransactionStatus(context.TODO(), req)
	assert.Nil(t, err)
	assert.Equal(t, api.TransactionStatusResponse_UNKNOWN, res.Status)
	resBytes, _ := json.Marshal(&api.TransactionStatusResponse{
		Timestamp: res.Timestamp,
		Status:    res.Status,
	})
	assert.Nil(t, crypto.VerifySignature(string(resBytes), pub, res.SignatureResponse))
}

/*
Scenario: Receive lock transaction request
	Given a transaction to lock
	When I want to request to lock it
	Then I get not error and the lock is stored
*/
func TestHandleLockTransaction(t *testing.T) {

	pub, pv := crypto.GenerateKeys()

	lockRepo := &mockLockRepository{}
	txRepo := &mockTxRepository{}
	sharedRepo := &mockSharedRepo{}
	txSrv := newTransactionServer(txRepo, lockRepo, sharedRepo, pub, pv)

	req := &api.LockRequest{
		Timestamp:           time.Now().Unix(),
		TransactionHash:     crypto.HashString("tx"),
		MasterPeerPublicKey: pub,
		Address:             crypto.HashString("addr"),
	}
	reqBytes, _ := json.Marshal(req)
	sig, _ := crypto.Sign(string(reqBytes), pv)
	req.SignatureRequest = sig

	res, err := txSrv.LockTransaction(context.TODO(), req)
	assert.Nil(t, err)
	resBytes, _ := json.Marshal(&api.LockResponse{
		Timestamp: res.Timestamp,
	})
	assert.Nil(t, crypto.VerifySignature(string(resBytes), pub, res.SignatureResponse))

	assert.Len(t, lockRepo.locks, 1)
	assert.Equal(t, crypto.HashString("addr"), lockRepo.locks[0].Address())
}

/*
Scenario: Receive unlock transaction request
	Given a transaction already
	When I want to request to unlock
	Then I get not error and the lock is removed
*/
func TestHandleUnlockTransaction(t *testing.T) {

	pub, pv := crypto.GenerateKeys()

	lockRepo := &mockLockRepository{}
	txRepo := &mockTxRepository{}
	sharedRepo := &mockSharedRepo{}
	txSrv := newTransactionServer(txRepo, lockRepo, sharedRepo, pub, pv)

	req := &api.LockRequest{
		Timestamp:           time.Now().Unix(),
		TransactionHash:     crypto.HashString("tx"),
		MasterPeerPublicKey: pub,
		Address:             crypto.HashString("addr"),
	}
	reqBytes, _ := json.Marshal(req)
	sig, _ := crypto.Sign(string(reqBytes), pv)
	req.SignatureRequest = sig

	res, err := txSrv.LockTransaction(context.TODO(), req)
	assert.Nil(t, err)
	resBytes, _ := json.Marshal(&api.LockResponse{
		Timestamp: res.Timestamp,
	})
	assert.Nil(t, crypto.VerifySignature(string(resBytes), pub, res.SignatureResponse))

	assert.Len(t, lockRepo.locks, 1)
	assert.Equal(t, crypto.HashString("addr"), lockRepo.locks[0].Address())

	req2 := &api.LockRequest{
		Timestamp:           time.Now().Unix(),
		TransactionHash:     crypto.HashString("tx"),
		MasterPeerPublicKey: pub,
		Address:             crypto.HashString("addr"),
	}
	reqBytes2, _ := json.Marshal(req2)
	sig2, _ := crypto.Sign(string(reqBytes2), pv)
	req2.SignatureRequest = sig2

	res2, err := txSrv.UnlockTransaction(context.TODO(), req)
	assert.Nil(t, err)
	resBytes2, _ := json.Marshal(&api.LockResponse{
		Timestamp: res.Timestamp,
	})
	assert.Nil(t, crypto.VerifySignature(string(resBytes2), pub, res2.SignatureResponse))

	assert.Len(t, lockRepo.locks, 0)
}

/*
Scenario: Receive lead mining transaction request
	Given a transaction to validate
	When I want to request to lead mining of the transaction
	Then I get not error
*/
func TestHandleLeadTransactionMining(t *testing.T) {

	pub, pv := crypto.GenerateKeys()

	lockRepo := &mockLockRepository{}
	txRepo := &mockTxRepository{}
	sharedRepo := &mockSharedRepo{}
	txSrv := newTransactionServer(txRepo, lockRepo, sharedRepo, pub, pv)

	data := map[string]string{
		"encrypted_address": hex.EncodeToString([]byte("addr")),
		"encrypted_wallet":  hex.EncodeToString([]byte("wallet")),
	}

	propKP, _ := shared.NewEmitterKeyPair(hex.EncodeToString([]byte("pvkey")), pub)
	prop, _ := transaction.NewProposal(propKP)
	txRaw := map[string]interface{}{
		"address":    crypto.HashString("addr"),
		"data":       data,
		"timestamp":  time.Now().Unix(),
		"type":       transaction.KeychainType,
		"public_key": pub,
		"proposal":   prop,
	}
	txBytes, _ := json.Marshal(txRaw)
	sig, _ := crypto.Sign(string(txBytes), pv)
	txRaw["signature"] = sig
	txRaw["em_signature"] = sig
	txBytes, _ = json.Marshal(txRaw)

	tx, _ := transaction.New(crypto.HashString("addr"), transaction.KeychainType, data, time.Now(), pub, sig, sig, prop, crypto.HashBytes(txBytes))
	req := &api.LeadMiningRequest{
		Timestamp:          time.Now().Unix(),
		MinimumValidations: 1,
		Transaction:        formatAPITransaction(tx),
	}

	reqBytes, _ := json.Marshal(req)
	sigReq, _ := crypto.Sign(string(reqBytes), pv)
	req.SignatureRequest = sigReq

	res, err := txSrv.LeadTransactionMining(context.TODO(), req)
	assert.Nil(t, err)

	time.Sleep(1 * time.Second)

	resBytes, _ := json.Marshal(&api.LeadMiningResponse{
		Timestamp: res.Timestamp,
	})
	assert.Nil(t, crypto.VerifySignature(string(resBytes), pub, res.SignatureResponse))

	assert.Len(t, txRepo.keychains, 1)
	assert.Equal(t, crypto.HashString("addr"), txRepo.keychains[0].Address())
}

/*
Scenario: Receive confirmation of validations transaction request
	Given a transaction to validate
	When I want to request to validation of the transaction
	Then I get the miner validation
*/
func TestHandleConfirmValiation(t *testing.T) {

	pub, pv := crypto.GenerateKeys()

	lockRepo := &mockLockRepository{}
	txRepo := &mockTxRepository{}
	sharedRepo := &mockSharedRepo{}
	txSrv := newTransactionServer(txRepo, lockRepo, sharedRepo, pub, pv)

	propKP, _ := shared.NewEmitterKeyPair(hex.EncodeToString([]byte("pvkey")), pub)
	data := map[string]string{
		"encrypted_address": hex.EncodeToString([]byte("addr")),
		"encrypted_wallet":  hex.EncodeToString([]byte("wallet")),
	}

	prop, _ := transaction.NewProposal(propKP)
	txRaw := map[string]interface{}{
		"address":    crypto.HashString("addr"),
		"data":       data,
		"timestamp":  time.Now().Unix(),
		"type":       transaction.KeychainType,
		"public_key": pub,
		"proposal":   prop,
	}
	txBytes, _ := json.Marshal(txRaw)
	sig, _ := crypto.Sign(string(txBytes), pv)
	txRaw["signature"] = sig
	txRaw["em_signature"] = sig
	txBytes, _ = json.Marshal(txRaw)
	tx, _ := transaction.New(crypto.HashString("addr"), transaction.KeychainType, data, time.Now(), pub, sig, sig, prop, crypto.HashBytes(txBytes))

	vRaw := map[string]interface{}{
		"status":     transaction.ValidationOK,
		"public_key": pub,
		"timestamp":  time.Now().Unix(),
	}

	vBytes, _ := json.Marshal(vRaw)
	vSig, _ := crypto.Sign(string(vBytes), pv)
	v, _ := transaction.NewMinerValidation(transaction.ValidationOK, time.Now(), pub, vSig)
	mv, _ := transaction.NewMasterValidation(transaction.Pool{}, pub, v)
	req := &api.ConfirmValidationRequest{
		Transaction:      formatAPITransaction(tx),
		Timestamp:        time.Now().Unix(),
		MasterValidation: formatAPIMasterValidation(mv),
	}

	reqBytes, _ := json.Marshal(req)
	sigReq, _ := crypto.Sign(string(reqBytes), pv)
	req.SignatureRequest = sigReq

	res, err := txSrv.ConfirmTransactionValidation(context.TODO(), req)
	assert.Nil(t, err)

	resBytes, _ := json.Marshal(&api.ConfirmValidationResponse{
		Timestamp:  res.Timestamp,
		Validation: res.Validation,
	})
	assert.Nil(t, crypto.VerifySignature(string(resBytes), pub, res.SignatureResponse))

	assert.NotNil(t, res.Validation)
	assert.Equal(t, api.MinerValidation_OK, res.Validation.Status)
	assert.Equal(t, pub, res.Validation.PublicKey)
}

/*
Scenario: Receive storage  transaction request
	Given a transaction
	When I want to request to store of the transaction
	Then the transaction is stored
*/
func TestHandleStoreTransaction(t *testing.T) {

	pub, pv := crypto.GenerateKeys()

	lockRepo := &mockLockRepository{}
	txRepo := &mockTxRepository{}
	sharedRepo := &mockSharedRepo{}
	txSrv := newTransactionServer(txRepo, lockRepo, sharedRepo, pub, pv)

	propKP, _ := shared.NewEmitterKeyPair(hex.EncodeToString([]byte("pvkey")), pub)

	data := map[string]string{
		"encrypted_address": hex.EncodeToString([]byte("addr")),
		"encrypted_wallet":  hex.EncodeToString([]byte("wallet")),
	}
	prop, _ := transaction.NewProposal(propKP)
	txRaw := map[string]interface{}{
		"address":    crypto.HashString("addr"),
		"data":       data,
		"timestamp":  time.Now().Unix(),
		"type":       transaction.KeychainType,
		"public_key": pub,
		"proposal":   prop,
	}
	txBytes, _ := json.Marshal(txRaw)
	sig, _ := crypto.Sign(string(txBytes), pv)
	txRaw["signature"] = sig
	txRaw["em_signature"] = sig
	txBytes, _ = json.Marshal(txRaw)

	tx, _ := transaction.New(crypto.HashString("addr"), transaction.KeychainType, data, time.Now(), pub, sig, sig, prop, crypto.HashBytes(txBytes))

	vRaw := map[string]interface{}{
		"status":     transaction.ValidationOK,
		"public_key": pub,
		"timestamp":  time.Now().Unix(),
	}
	vBytes, _ := json.Marshal(vRaw)
	vSig, _ := crypto.Sign(string(vBytes), pv)
	v, _ := transaction.NewMinerValidation(transaction.ValidationOK, time.Now(), pub, vSig)
	mv, _ := transaction.NewMasterValidation(transaction.Pool{}, pub, v)

	req := &api.StoreRequest{
		Timestamp: time.Now().Unix(),
		MinedTransaction: &api.MinedTransaction{
			Transaction:        formatAPITransaction(tx),
			MasterValidation:   formatAPIMasterValidation(mv),
			ConfirmValidations: []*api.MinerValidation{formatAPIValidation(v)},
		},
	}

	reqBytes, _ := json.Marshal(req)
	sigReq, _ := crypto.Sign(string(reqBytes), pv)
	req.SignatureRequest = sigReq

	res, err := txSrv.StoreTransaction(context.TODO(), req)
	assert.Nil(t, err)

	resBytes, _ := json.Marshal(&api.StoreResponse{
		Timestamp: res.Timestamp,
	})
	assert.Nil(t, crypto.VerifySignature(string(resBytes), pub, res.SignatureResponse))

	assert.Len(t, txRepo.keychains, 1)
	assert.Equal(t, crypto.HashBytes(txBytes), txRepo.keychains[0].TransactionHash())

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

type mockPoolRequester struct {
	stores []transaction.Transaction
	repo   *mockTxRepository
}

func (pr mockPoolRequester) RequestTransactionLock(pool transaction.Pool, txLock transaction.Lock) error {
	return nil
}

func (pr mockPoolRequester) RequestTransactionUnlock(pool transaction.Pool, txLock transaction.Lock) error {
	return nil
}

func (pr mockPoolRequester) RequestTransactionValidations(pool transaction.Pool, tx transaction.Transaction, masterValid transaction.MasterValidation, validChan chan<- transaction.MinerValidation) {
	pub, pv := crypto.GenerateKeys()

	vRaw := map[string]interface{}{
		"status":     transaction.ValidationOK,
		"public_key": pub,
		"timestamp":  time.Now().Unix(),
	}
	vBytes, _ := json.Marshal(vRaw)
	sig, _ := crypto.Sign(string(vBytes), pv)
	v, _ := transaction.NewMinerValidation(transaction.ValidationOK, time.Now(), pub, sig)

	validChan <- v
}

func (pr *mockPoolRequester) RequestTransactionStorage(pool transaction.Pool, tx transaction.Transaction, ackChan chan<- bool) {
	pr.stores = append(pr.stores, tx)
	if tx.Type() == transaction.KeychainType {
		k, _ := transaction.NewKeychain(tx)
		pr.repo.keychains = append(pr.repo.keychains, k)
	}
	if tx.Type() == transaction.IDType {
		id, err := transaction.NewID(tx)
		log.Print(err)
		pr.repo.ids = append(pr.repo.ids, id)
	}
	ackChan <- true
}

type mockTxRepository struct {
	pendings  []transaction.Transaction
	kos       []transaction.Transaction
	keychains []transaction.Keychain
	ids       []transaction.ID
}

func (r mockTxRepository) FindPendingTransaction(txHash string) (*transaction.Transaction, error) {
	for _, tx := range r.pendings {
		if tx.TransactionHash() == txHash {
			return &tx, nil
		}
	}
	return nil, nil
}

func (r mockTxRepository) GetKeychain(txAddr string) (*transaction.Keychain, error) {
	sort.Slice(r.keychains, func(i, j int) bool {
		return r.keychains[i].Timestamp().Unix() > r.keychains[j].Timestamp().Unix()
	})

	if len(r.keychains) > 0 {
		return &r.keychains[0], nil
	}
	return nil, nil
}

func (r mockTxRepository) FindKeychainByHash(txHash string) (*transaction.Keychain, error) {
	for _, tx := range r.keychains {
		if tx.TransactionHash() == txHash {
			return &tx, nil
		}
	}
	return nil, nil
}

func (r mockTxRepository) FindLastKeychain(addr string) (*transaction.Keychain, error) {

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

func (r mockTxRepository) FindIDByHash(txHash string) (*transaction.ID, error) {
	for _, tx := range r.ids {
		if tx.TransactionHash() == txHash {
			return &tx, nil
		}
	}
	return nil, nil
}

func (r mockTxRepository) FindIDByAddress(addr string) (*transaction.ID, error) {
	for _, tx := range r.ids {
		if tx.Address() == addr {
			return &tx, nil
		}
	}
	return nil, nil
}

func (r mockTxRepository) FindKOTransaction(txHash string) (*transaction.Transaction, error) {
	for _, tx := range r.kos {
		if tx.TransactionHash() == txHash {
			return &tx, nil
		}
	}
	return nil, nil
}

func (r *mockTxRepository) StoreKeychain(kc transaction.Keychain) error {
	r.keychains = append(r.keychains, kc)
	return nil
}

func (r *mockTxRepository) StoreID(id transaction.ID) error {
	r.ids = append(r.ids, id)
	return nil
}

func (r *mockTxRepository) StoreKO(tx transaction.Transaction) error {
	r.kos = append(r.kos, tx)
	return nil
}

type mockLockRepository struct {
	locks []transaction.Lock
}

func (r *mockLockRepository) StoreLock(l transaction.Lock) error {
	r.locks = append(r.locks, l)
	return nil
}

func (r *mockLockRepository) RemoveLock(l transaction.Lock) error {
	pos := r.findLockPosition(l)
	if pos > -1 {
		r.locks = append(r.locks[:pos], r.locks[pos+1:]...)
	}
	return nil
}

func (r mockLockRepository) ContainsLock(l transaction.Lock) (bool, error) {
	return r.findLockPosition(l) > -1, nil
}

func (r mockLockRepository) findLockPosition(l transaction.Lock) int {
	for i, lock := range r.locks {
		if lock.TransactionHash() == l.TransactionHash() && l.MasterRobotKey() == lock.MasterRobotKey() {
			return i
		}
	}
	return -1
}

func newTransactionServer(txRepo *mockTxRepository, lockRepo *mockLockRepository, sharedRepo *mockSharedRepo, pub string, pv string) api.TransactionServiceServer {

	encPv, _ := crypto.Encrypt(pv, pub)
	minerKP, _ := shared.NewMinerKeyPair(pub, pv)
	propKP, _ := shared.NewEmitterKeyPair(encPv, pub)

	sharedRepo.emKeys = shared.EmitterKeys{propKP}
	sharedRepo.minerKeys = minerKP

	sharedService := shared.NewService(sharedRepo)
	poolR := &mockPoolRequester{
		repo: txRepo,
	}

	lockSrv := transaction.NewLockService(lockRepo)
	poolFindSrv := transaction.NewPoolFindingService(NewPoolRetriever(sharedService))
	miningSrv := transaction.NewMiningService(poolR, poolFindSrv, sharedService, "127.0.0.1", pub, pv)
	storeSrv := transaction.NewStorageService(txRepo, miningSrv)

	return NewTransactionServer(storeSrv, lockSrv, miningSrv, sharedService)
}
