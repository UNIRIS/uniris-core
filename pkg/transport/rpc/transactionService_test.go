package rpc

import (
	"bytes"
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"errors"
	"net"
	"time"

	"golang.org/x/crypto/ed25519"
)

// /*
// Scenario: Receive  get last transction about an unknown transaction
// 	Given no transaction store for an address
// 	When I want to request to retrieve the last transaction keychain of this unknown address
// 	Then I get an error
// */
// func TestHandleGetLastTransactionWhenNotExist(t *testing.T) {
// 	pub, pv, _ := ed25519.GenerateKey(rand.Reader)

// 	sharedKeyReader := &mockSharedKeyReader{}
// 	nodeKey, _ := shared.NewNodeCrossKeyPair(pub, pv)
// 	sharedKeyReader.crossNodeKeys = append(sharedKeyReader.crossNodeKeys, nodeKey)

// 	nodeReader := &mockNodeReader{
// 		nodes: []node{
// 			consensus.NewNode(net.ParseIP("127.0.0.1"), 5000, pub, nodeOK, "", 300, "1.0", 0, 1, 30.0, -10.0, consensus.GeoPatch{}, true),
// 		},
// 	}

// 	chainDB := &mockChainDB{}

// 	poolR := &mockPoolRequester{
// 		repo: chainDB,
// 	}

// 	l := logging.NewLogger("stdout", log.New(os.Stdout, "", 0), "test", net.ParseIP("127.0.0.1"), logging.ErrorLogLevel)
// 	txSrv := NewTransactionService(chainDB, sharedKeyReader, nodeReader, poolR, pub, pv, l)

// 	req := &api.GetLastTransactionRequest{
// 		Timestamp:          time.Now().Unix(),
// 		TransactionAddress: []byte("address"),
// 		Type:               api.TransactionType_KEYCHAIN,
// 	}
// 	reqBytes, _ := json.Marshal(req)
// 	sig, _ := pv.Sign(reqBytes)
// 	req.SignatureRequest = sig

// 	_, err := txSrv.GetLastTransaction(context.TODO(), req)
// 	assert.NotNil(t, err)
// 	statusCode, _ := status.FromError(err)
// 	assert.Equal(t, codes.NotFound, statusCode.Code())
// 	assert.Equal(t, statusCode.Message(), "transaction does not exist")
// }

// /*
// Scenario: Receive  get last transaction request
// 	Given a keychain transaction stored
// 	When I want to request to retrieve the last transaction keychain of this address
// 	Then I get an error
// */
// func TestHandleGetLastTransaction(t *testing.T) {

// 	pub, pv, _ := ed25519.GenerateKey(rand.Reader)
// 	pubB, _ := pub.Marshal()

// 	sharedKeyReader := &mockSharedKeyReader{}
// 	nodeKey, _ := shared.NewNodeCrossKeyPair(pub, pv)
// 	sharedKeyReader.crossNodeKeys = append(sharedKeyReader.crossNodeKeys, nodeKey)

// 	nodeReader := &mockNodeReader{
// 		nodes: []node{
// 			consensus.NewNode(net.ParseIP("127.0.0.1"), 5000, pub, nodeOK, "", 300, "1.0", 0, 1, 30.0, -10.0, consensus.GeoPatch{}, true),
// 		},
// 	}

// 	chainDB := &mockChainDB{}

// 	poolR := &mockPoolRequester{
// 		repo: chainDB,
// 	}

// 	l := logging.NewLogger("stdout", log.New(os.Stdout, "", 0), "test", net.ParseIP("127.0.0.1"), logging.ErrorLogLevel)
// 	txSrv := NewTransactionService(chainDB, sharedKeyReader, nodeReader, poolR, pub, pv, l)

// 	data := map[string][]byte{
// 		"encrypted_address_by_node": []byte("addr"),
// 		"encrypted_wallet":          []byte("wallet"),
// 	}

// 	prop, _ := shared.NewEmitterCrossKeyPair([]byte("pvkey"), pub)

// 	txRaw := map[string]interface{}{
// 		"addr": []byte("addr"),
// 		"data": map[string]string{
// 			"encrypted_address_by_node": hex.EncodeToString([]byte("addr")),
// 			"encrypted_wallet":          hex.EncodeToString([]byte("wallet")),
// 			"em_shared_keys_proposal": map[string]string{
// 				"encrypted_private_key": hex.EncodeToString([]byte("pvkey")),
// 				"public_key":            hex.EncodeToString(pubB),
// 			},
// 		},
// 		"timestamp":  time.Now().Unix(),
// 		"type":       chain.Keychainint,
// 		"public_key": hex.EncodeToString(pubB),
// 	}
// 	txBytes, _ := json.Marshal(txRaw)
// 	sig, _ := pv.Sign(txBytes)
// 	txRaw["signature"] = hex.EncodeToString(sig)

// 	txByteWithSign, _ := json.Marshal(txRaw)
// 	emSig, _ := pv.Sign(txByteWithSign)
// 	txRaw["em_signature"] = hex.EncodeToString(emSig)
// 	txBytes, _ = json.Marshal(txRaw)

// 	tx, _ := chain.NewTransaction([]byte("addr"), chain.Keychainint, data, time.Now(), pub, sig, nil, nil)
// 	keychain, _ := chain.NewKeychain(tx)
// 	chainDB.keychains = append(chainDB.keychains, keychain)

// 	req := &api.GetLastTransactionRequest{
// 		Timestamp:          time.Now().Unix(),
// 		TransactionAddress: []byte("addr"),
// 		Type:               api.TransactionType_KEYCHAIN,
// 	}
// 	reqBytes, _ := json.Marshal(req)
// 	sigReq, _ := pv.Sign(reqBytes)
// 	req.SignatureRequest = sigReq

// 	res, err := txSrv.GetLastTransaction(context.TODO(), req)
// 	assert.Nil(t, err)
// 	assert.NotEmpty(t, res.SignatureResponse)
// 	assert.NotNil(t, res.Transaction)
// 	assert.EqualValues(t, []byte(txBytes), res.Transaction.TransactionHash)

// 	resBytes, _ := json.Marshal(&api.GetLastTransactionResponse{
// 		Timestamp:   res.Timestamp,
// 		Transaction: res.Transaction,
// 	})
// 	assert.True(t, pub.Verify(resBytes, res.SignatureResponse))
// }

// /*
// Scenario: Receive get transaction status request
// 	Given no transaction stored
// 	When I want to request the transactions status for this transaction hash
// 	Then I get a status unknown
// */
// func TestHandleGetTransactionStatus(t *testing.T) {

// 	pub, pv, _ := ed25519.GenerateKey(rand.Reader)

// 	sharedKeyReader := &mockSharedKeyReader{}
// 	nodeKey, _ := shared.NewNodeCrossKeyPair(pub, pv)
// 	sharedKeyReader.crossNodeKeys = append(sharedKeyReader.crossNodeKeys, nodeKey)

// 	nodeReader := &mockNodeReader{
// 		nodes: []node{
// 			consensus.NewNode(net.ParseIP("127.0.0.1"), 5000, pub, nodeOK, "", 300, "1.0", 0, 1, 30.0, -10.0, consensus.GeoPatch{}, true),
// 		},
// 	}

// 	chainDB := &mockChainDB{}

// 	poolR := &mockPoolRequester{
// 		repo: chainDB,
// 	}

// 	l := logging.NewLogger("stdout", log.New(os.Stdout, "", 0), "test", net.ParseIP("127.0.0.1"), logging.ErrorLogLevel)
// 	txSrv := NewTransactionService(chainDB, sharedKeyReader, nodeReader, poolR, pub, pv, l)

// 	req := &api.GetTransactionStatusRequest{
// 		Timestamp:       time.Now().Unix(),
// 		TransactionHash: []byte("tx1"),
// 	}
// 	reqBytes, _ := json.Marshal(req)
// 	sig, _ := pv.Sign(reqBytes)
// 	req.SignatureRequest = sig

// 	res, err := txSrv.GetTransactionStatus(context.TODO(), req)
// 	assert.Nil(t, err)
// 	assert.Equal(t, api.TransactionStatus_UNKNOWN, res.Status)
// 	resBytes, _ := json.Marshal(&api.GetTransactionStatusResponse{
// 		Timestamp: res.Timestamp,
// 		Status:    res.Status,
// 	})
// 	assert.True(t, pub.Verify(resBytes, res.SignatureResponse))
// }

// /*
// Scenario: Receive storage  transaction request
// 	Given a transaction
// 	When I want to request to store of the transaction
// 	Then the transaction is stored
// */
// func TestHandleStoreTransaction(t *testing.T) {

// 	pub, pv, _ := ed25519.GenerateKey(rand.Reader)
// 	pubB, _ := pub.Marshal()

// 	sharedKeyReader := &mockSharedKeyReader{}
// 	nodeKey, _ := shared.NewNodeCrossKeyPair(pub, pv)
// 	sharedKeyReader.crossNodeKeys = append(sharedKeyReader.crossNodeKeys, nodeKey)

// 	nodeReader := &mockNodeReader{
// 		nodes: []node{
// 			consensus.NewNode(net.ParseIP("127.0.0.1"), 5000, pub, nodeOK, "", 300, "1.0", 0, 1, 30.0, -10.0, consensus.GeoPatch{}, true),
// 		},
// 	}

// 	chainDB := &mockChainDB{}

// 	poolR := &mockPoolRequester{
// 		repo: chainDB,
// 	}

// 	l := logging.NewLogger("stdout", log.New(os.Stdout, "", 0), "test", net.ParseIP("127.0.0.1"), logging.ErrorLogLevel)
// 	txSrv := NewTransactionService(chainDB, sharedKeyReader, nodeReader, poolR, pub, pv, l)

// 	prop, _ := shared.NewEmitterCrossKeyPair([]byte("pvkey"), pub)

// 	data := map[string][]byte{
// 		"encrypted_address_by_node": []byte("addr"),
// 		"encrypted_wallet":          []byte("wallet"),
// 	}
// 	txRaw := map[string]interface{}{
// 		"addr": []byte("addr"),
// 		"data": map[string]string{
// 			"encrypted_address_by_node": hex.EncodeToString([]byte("addr")),
// 			"encrypted_wallet":          hex.EncodeToString([]byte("wallet")),
// 		},
// 		"timestamp":  time.Now().Unix(),
// 		"type":       chain.Keychainint,
// 		"public_key": hex.EncodeToString(pubB),
// 		"em_shared_keys_proposal": map[string]string{
// 			"encrypted_private_key": hex.EncodeToString([]byte("pvkey")),
// 			"public_key":            hex.EncodeToString(pubB),
// 		},
// 	}
// 	txBytes, _ := json.Marshal(txRaw)
// 	sig, _ := pv.Sign(txBytes)
// 	txRaw["signature"] = hex.EncodeToString(sig)

// 	txByteWithSign, _ := json.Marshal(txRaw)
// 	emSig, _ := pv.Sign(txByteWithSign)
// 	txRaw["em_signature"] = hex.EncodeToString(emSig)
// 	txBytes, _ = json.Marshal(txRaw)

// 	txBytes, _ = json.Marshal(txRaw)

// 	tx, _ := chain.NewTransaction([]byte("addr"), chain.Keychainint, data, time.Now(), pub, prop, sig, emSig, []byte(txBytes))

// 	vRaw := map[string]interface{}{
// 		"status":     validationStampOK,
// 		"public_key": pubB,
// 		"timestamp":  time.Now().Unix(),
// 	}
// 	vBytes, _ := json.Marshal(vRaw)
// 	vSig, _ := pv.Sign(vBytes)
// 	v, _ := chain.NewValidation(validationStampOK, time.Now(), pub, vSig)
// 	wHeaders := chain.NewWelcomeNodeHeader(pub, []chain.NodeHeader{chain.NewNodeHeader(pub, false, false, 0, true)}, []byte("sig"))
// 	vHeaders := []chain.NodeHeader{chain.NewNodeHeader(pub, false, false, 0, true)}
// 	sHeaders := []chain.NodeHeader{chain.NewNodeHeader(pub, false, false, 0, true)}
// 	mv, _ := chain.NewcoordStampation([]publicKey{}, pub, v, wHeaders, vHeaders, sHeaders)

// 	txf, _ := formatAPITransaction(tx)
// 	mvf, _ := formatAPIcoordStampation(mv)
// 	vf, _ := formatAPIValidation(v)

// 	req := &api.StoreTransactionRequest{
// 		Timestamp: time.Now().Unix(),
// 		MinedTransaction: &api.MinedTransaction{
// 			Transaction:        txf,
// 			coordStampation:    mvf,
// 			ConfirmValidations: []*api.Validation{vf},
// 		},
// 	}

// 	reqBytes, _ := json.Marshal(req)
// 	sigReq, _ := pv.Sign(reqBytes)
// 	req.SignatureRequest = sigReq

// 	res, err := txSrv.StoreTransaction(context.TODO(), req)
// 	assert.Nil(t, err)

// 	resBytes, _ := json.Marshal(&api.StoreTransactionResponse{
// 		Timestamp: res.Timestamp,
// 	})
// 	assert.True(t, pub.Verify(resBytes, res.SignatureResponse))

// 	assert.Len(t, chainDB.keychains, 1)
// 	assert.EqualValues(t, []byte(txBytes), chainDB.keychains[0].TransactionHash())

// }

// /*
// Scenario: Receive lock transaction request
// 	Given a transaction to lock
// 	When I want to request to lock it
// 	Then I get not error and the lock is stored
// */
// func TestHandleLockTransaction(t *testing.T) {

// 	pub, pv, _ := ed25519.GenerateKey(rand.Reader)

// 	chainDB := &mockChainDB{}
// 	sharedKeyReader := &mockSharedKeyReader{}
// 	nodeKey, _ := shared.NewNodeCrossKeyPair(pub, pv)
// 	sharedKeyReader.crossNodeKeys = append(sharedKeyReader.crossNodeKeys, nodeKey)

// 	nodeReader := &mockNodeReader{
// 		nodes: []node{
// 			consensus.NewNode(net.ParseIP("127.0.0.1"), 5000, pub, nodeOK, "", 300, "1.0", 0, 1, 30.0, -10.0, consensus.GeoPatch{}, true),
// 		},
// 	}

// 	poolR := &mockPoolRequester{}
// 	l := logging.NewLogger("stdout", log.New(os.Stdout, "", 0), "test", net.ParseIP("127.0.0.1"), logging.ErrorLogLevel)
// 	txSrv := NewTransactionService(chainDB, sharedKeyReader, nodeReader, poolR, pub, pv, l)

// 	pubB, _ := pub.Marshal()

// 	req := &api.TimeLockTransactionRequest{
// 		Timestamp:           time.Now().Unix(),
// 		TransactionHash:     []byte("tx1"),
// 		MasterNodePublicKey: pubB,
// 		Address:             []byte("addr1"),
// 	}
// 	reqBytes, _ := json.Marshal(req)
// 	sig, _ := pv.Sign(reqBytes)
// 	req.SignatureRequest = sig

// 	res, err := txSrv.TimeLockTransaction(context.TODO(), req)
// 	assert.Nil(t, err)
// 	resBytes, _ := json.Marshal(&api.TimeLockTransactionResponse{
// 		Timestamp: res.Timestamp,
// 	})
// 	assert.True(t, pub.Verify(resBytes, res.SignatureResponse))
// 	assert.True(t, chain.ContainsTimeLock([]byte("tx1"), []byte("addr1")))
// }

// /*
// Scenario: Receive lead mining transaction request
// 	Given a transaction to validate
// 	When I want to request to lead mining of the transaction
// 	Then I get not error
// */
// func TestHandleLeadTransactionMining(t *testing.T) {

// 	pub, pv, _ := ed25519.GenerateKey(rand.Reader)
// 	pubB, _ := pub.Marshal()

// 	sharedKeyReader := &mockSharedKeyReader{}
// 	nodeKey, _ := shared.NewNodeCrossKeyPair(pub, pv)
// 	sharedKeyReader.crossNodeKeys = append(sharedKeyReader.crossNodeKeys, nodeKey)
// 	emKey, _ := shared.NewEmitterCrossKeyPair([]byte("encpv"), pub)
// 	sharedKeyReader.crossEmitterKeys = append(sharedKeyReader.crossEmitterKeys, emKey)
// 	sharedKeyReader.authKeys = append(sharedKeyReader.authKeys, pub)

// 	nodeReader := &mockNodeReader{
// 		nodes: []node{
// 			consensus.NewNode(net.ParseIP("127.0.0.1"), 5000, pub, nodeOK, "", 300, "1.0", 0, 1, 30.0, -10.0, consensus.GeoPatch{}, true),
// 		},
// 	}

// 	chainDB := &mockChainDB{}

// 	poolR := &mockPoolRequester{
// 		repo: chainDB,
// 	}
// 	l := logging.NewLogger("stdout", log.New(os.Stdout, "", 0), "test", net.ParseIP("127.0.0.1"), logging.ErrorLogLevel)
// 	txSrv := NewTransactionService(chainDB, sharedKeyReader, nodeReader, poolR, pub, pv, l)
// 	data := map[string][]byte{
// 		"encrypted_address_by_node": []byte("addr"),
// 		"encrypted_wallet":          []byte("wallet"),
// 	}

// 	prop, _ := shared.NewEmitterCrossKeyPair([]byte("pvkey"), pub)
// 	txRaw := map[string]interface{}{
// 		"addr": []byte("addr"),
// 		"data": map[string]string{
// 			"encrypted_address_by_node": hex.EncodeToString([]byte("addr")),
// 			"encrypted_wallet":          hex.EncodeToString([]byte("wallet")),
// 		},
// 		"timestamp":  time.Now().Unix(),
// 		"type":       chain.Keychainint,
// 		"public_key": hex.EncodeToString(pubB),
// 		"em_shared_keys_proposal": map[string]string{
// 			"encrypted_private_key": hex.EncodeToString([]byte("pvkey")),
// 			"public_key":            hex.EncodeToString(pubB),
// 		},
// 	}
// 	txBytes, _ := json.Marshal(txRaw)
// 	sig, _ := pv.Sign(txBytes)
// 	txRaw["signature"] = hex.EncodeToString(sig)

// 	txByteWithSig, _ := json.Marshal(txRaw)
// 	emSig, _ := pv.Sign(txByteWithSig)
// 	txRaw["em_signature"] = hex.EncodeToString(emSig)

// 	txBytes, _ = json.Marshal(txRaw)

// 	tx, _ := chain.NewTransaction([]byte("addr"), chain.Keychainint, data, time.Now(), pub, prop, sig, emSig, []byte(txBytes))
// 	txf, _ := formatAPITransaction(tx)
// 	ml := []*api.NodeHeader{
// 		&api.NodeHeader{
// 			IsMaster:      true,
// 			IsUnreachable: false,
// 			PatchNumber:   1,
// 			PublicKey:     pubB,
// 		}}
// 	req := &api.LeadTransactionMiningRequest{
// 		Timestamp:          time.Now().Unix(),
// 		MinimumValidations: 1,
// 		WelcomeHeaders: &api.WelcomeNodeHeader{
// 			PublicKey:   pubB,
// 			MastersList: ml,
// 			Signature:   []byte("sig"),
// 		},
// 		Transaction: txf,
// 	}

// 	reqBytes, _ := json.Marshal(req)
// 	sigReq, _ := pv.Sign(reqBytes)
// 	req.SignatureRequest = sigReq

// 	res, err := txSrv.LeadTransactionMining(context.TODO(), req)
// 	assert.Nil(t, err)

// 	time.Sleep(1 * time.Second)

// 	resBytes, _ := json.Marshal(&api.LeadTransactionMiningResponse{
// 		Timestamp: res.Timestamp,
// 	})
// 	assert.True(t, pub.Verify(resBytes, res.SignatureResponse))

// 	assert.Len(t, chainDB.keychains, 1)
// 	assert.EqualValues(t, []byte("addr"), chainDB.keychains[0].Address())
// }

// /*
// Scenario: Receive confirmation of validations transaction request
// 	Given a transaction to validate
// 	When I want to request to validation of the transaction
// 	Then I get the node validation
// */
// func TestHandleConfirmValiation(t *testing.T) {

// 	pub, pv, _ := ed25519.GenerateKey(rand.Reader)
// 	pubB, _ := pub.Marshal()

// 	sharedKeyReader := &mockSharedKeyReader{}
// 	nodeKey, _ := shared.NewNodeCrossKeyPair(pub, pv)
// 	sharedKeyReader.crossNodeKeys = append(sharedKeyReader.crossNodeKeys, nodeKey)
// 	emKey, _ := shared.NewEmitterCrossKeyPair([]byte("encpv"), pub)
// 	sharedKeyReader.crossEmitterKeys = append(sharedKeyReader.crossEmitterKeys, emKey)

// 	nodeReader := &mockNodeReader{
// 		nodes: []node{
// 			consensus.NewNode(net.ParseIP("127.0.0.1"), 5000, pub, nodeOK, "", 300, "1.0", 0, 1, 30.0, -10.0, consensus.GeoPatch{}, true),
// 		},
// 	}

// 	chainDB := &mockChainDB{}

// 	poolR := &mockPoolRequester{
// 		repo: chainDB,
// 	}
// 	l := logging.NewLogger("stdout", log.New(os.Stdout, "", 0), "test", net.ParseIP("127.0.0.1"), logging.ErrorLogLevel)
// 	txSrv := NewTransactionService(chainDB, sharedKeyReader, nodeReader, poolR, pub, pv, l)

// 	prop, _ := shared.NewEmitterCrossKeyPair([]byte("pvkey"), pub)
// 	data := map[string][]byte{
// 		"encrypted_address_by_node": []byte("addr"),
// 		"encrypted_wallet":          []byte("wallet"),
// 	}

// 	txRaw := map[string]interface{}{
// 		"addr": []byte("addr"),
// 		"data": map[string]string{
// 			"encrypted_address_by_node": hex.EncodeToString([]byte("addr")),
// 			"encrypted_wallet":          hex.EncodeToString([]byte("wallet")),
// 		},
// 		"timestamp":  time.Now().Unix(),
// 		"type":       chain.Keychainint,
// 		"public_key": hex.EncodeToString(pubB),
// 		"em_shared_keys_proposal": map[string]string{
// 			"encrypted_private_key": hex.EncodeToString([]byte("pvkey")),
// 			"public_key":            hex.EncodeToString(pubB),
// 		},
// 	}
// 	txBytes, _ := json.Marshal(txRaw)
// 	sig, _ := pv.Sign(txBytes)
// 	txRaw["signature"] = hex.EncodeToString(sig)
// 	txByteWithSig, _ := json.Marshal(txRaw)
// 	emSig, _ := pv.Sign(txByteWithSig)
// 	txRaw["em_signature"] = hex.EncodeToString(emSig)
// 	txBytes, _ = json.Marshal(txRaw)
// 	tx, _ := chain.NewTransaction([]byte("addr"), chain.Keychainint, data, time.Now(), pub, prop, sig, emSig, []byte(txBytes))

// 	vRaw := map[string]interface{}{
// 		"status":     validationStampOK,
// 		"public_key": pubB,
// 		"timestamp":  time.Now().Unix(),
// 	}

// 	vBytes, _ := json.Marshal(vRaw)
// 	vSig, _ := pv.Sign(vBytes)
// 	v, _ := chain.NewValidation(validationStampOK, time.Now(), pub, vSig)
// 	wHeaders := chain.NewWelcomeNodeHeader(pub, []chain.NodeHeader{chain.NewNodeHeader(pub, false, false, 0, true)}, []byte("sig"))
// 	vHeaders := []chain.NodeHeader{chain.NewNodeHeader(pub, false, false, 0, true)}
// 	sHeaders := []chain.NodeHeader{chain.NewNodeHeader(pub, false, false, 0, true)}
// 	mv, _ := chain.NewcoordStampation([]publicKey{}, pub, v, wHeaders, vHeaders, sHeaders)

// 	txv, _ := formatAPITransaction(tx)
// 	mvf, _ := formatAPIcoordStampation(mv)

// 	req := &api.ConfirmTransactionValidationRequest{
// 		Transaction:     txv,
// 		Timestamp:       time.Now().Unix(),
// 		coordStampation: mvf,
// 	}

// 	reqBytes, _ := json.Marshal(req)
// 	sigReq, _ := pv.Sign(reqBytes)
// 	req.SignatureRequest = sigReq

// 	res, err := txSrv.ConfirmTransactionValidation(context.TODO(), req)
// 	assert.Nil(t, err)

// 	resBytes, _ := json.Marshal(&api.ConfirmTransactionValidationResponse{
// 		Timestamp:  res.Timestamp,
// 		Validation: res.Validation,
// 	})
// 	assert.True(t, pub.Verify(resBytes, res.SignatureResponse))

// 	assert.NotNil(t, res.Validation)
// 	assert.Equal(t, api.Validation_OK, res.Validation.Status)
// 	assert.EqualValues(t, pubB, res.Validation.PublicKey)
// }

type mockPoolRequester struct {
	stores []transaction
	repo   *mockChainDB
}

func (pr mockPoolRequester) RequestLastTransaction(pool electedNodeList, txAddr []byte, txType int) (transaction, error) {
	return nil, nil
}

func (pr mockPoolRequester) RequestTransactionTimeLock(pool electedNodeList, txHash []byte, txAddr []byte, masterPublicKey publicKey) error {
	return nil
}

func (pr mockPoolRequester) RequestTransactionUnlock(pool electedNodeList, txHash []byte, txAddr []byte) error {
	return nil
}

func (pr mockPoolRequester) RequestTransactionValidations(pool electedNodeList, tx transaction, minValids int, coordStamp coordinatorStamp) ([]validationStamp, error) {
	pub, pv, _ := ed25519.GenerateKey(rand.Reader)

	vRaw := map[string]interface{}{
		"status":     1,
		"public_key": pub,
		"timestamp":  time.Now().Unix(),
	}
	vBytes, _ := json.Marshal(vRaw)
	sig := ed25519.Sign(pv, vBytes)
	v := mockValidationStamp{
		nodePub:   mockPublicKey{bytes: pub},
		sig:       sig,
		status:    1,
		timestamp: time.Now(),
	}

	return []validationStamp{v}, nil
}

func (pr *mockPoolRequester) RequestTransactionStorage(pool electedNodeList, minReplicas int, tx transaction) error {
	pr.stores = append(pr.stores, tx)
	if tx.Type() == 0 {
		pr.repo.keychains = append(pr.repo.keychains, tx)
	}
	if tx.Type() == 1 {
		pr.repo.ids = append(pr.repo.ids, tx)
	}
	return nil
}

type mockChainDB struct {
	kos       []transaction
	keychains []transaction
	ids       []transaction
}

func (r *mockChainDB) FindKeychainByAddr(addr []byte) (interface{}, error) {
	for _, tx := range r.keychains {
		if bytes.Equal(tx.Address(), addr) {
			return tx, nil
		}
	}
	return nil, errors.New("transaction not found")
}
func (r *mockChainDB) FindKeychainByHash(txHash []byte) (interface{}, error) {
	for _, tx := range r.keychains {
		if bytes.Equal(tx.CoordinatorStamp().(coordinatorStamp).TransactionHash(), txHash) {
			return tx, nil
		}
	}
	return nil, errors.New("transaction not found")
}
func (r mockChainDB) FindIDByHash(txHash []byte) (interface{}, error) {
	for _, tx := range r.ids {
		if bytes.Equal(tx.CoordinatorStamp().(coordinatorStamp).TransactionHash(), txHash) {
			return tx, nil
		}
	}
	return nil, errors.New("transaction not found")
}
func (r mockChainDB) FindIDByAddr(addr []byte) (interface{}, error) {
	for _, tx := range r.ids {
		if bytes.Equal(tx.Address(), addr) {
			return tx, nil
		}
	}
	return nil, errors.New("transaction not found")
}
func (r mockChainDB) FindKOByHash(txHash []byte) (interface{}, error) {
	for _, tx := range r.kos {
		if bytes.Equal(tx.CoordinatorStamp().(coordinatorStamp).TransactionHash(), txHash) {
			return tx, nil
		}
	}
	return nil, errors.New("transaction not found")
}
func (r mockChainDB) FindKOByAddr(addr []byte) (interface{}, error) {
	for _, tx := range r.kos {
		if bytes.Equal(tx.Address(), addr) {
			return tx, nil
		}
	}
	return nil, errors.New("transaction not found")
}

func (r *mockChainDB) WriteKeychain(tx interface{}) error {
	r.keychains = append(r.keychains, tx.(transaction))
	return nil
}
func (r *mockChainDB) WriteID(tx interface{}) error {
	r.ids = append(r.ids, tx.(transaction))
	return nil
}

func (r *mockChainDB) WriteKO(tx interface{}) error {
	r.kos = append(r.kos, tx.(transaction))
	return nil
}

type mockIndexDB struct {
	rows map[string][]byte
}

func (db mockIndexDB) FindLastTransactionAddr(genesis []byte) ([]byte, error) {
	return db.rows[hex.EncodeToString(genesis)], nil
}

type mockSharedKeyReader struct {
	crossNodePvKeys  []privateKey
	crossNodePubKeys []publicKey
	crossEmitterKeys []emitterCrossKeyPair
	authKeys         []publicKey
}

func (r mockSharedKeyReader) EmitterCrossKeypairs() ([]emitterCrossKeyPair, error) {
	return r.crossEmitterKeys, nil
}

func (r mockSharedKeyReader) FirstNodeCrossKeypair() (publicKey, privateKey, error) {
	return r.crossNodePubKeys[0], r.crossNodePvKeys[0], nil
}

func (r mockSharedKeyReader) LastNodeCrossKeypair() (publicKey, privateKey, error) {
	return r.crossNodePubKeys[len(r.crossNodePubKeys)-1], r.crossNodePvKeys[len(r.crossNodePvKeys)-1], nil
}

func (r mockSharedKeyReader) AuthorizedNodesPublicKeys() ([]publicKey, error) {
	return r.authKeys, nil
}

func (r mockSharedKeyReader) CrossEmitterPublicKeys() (pubKeys []publicKey, err error) {
	for _, kp := range r.crossEmitterKeys {
		pubKeys = append(pubKeys, kp.PublicKey())
	}
	return
}

func (r mockSharedKeyReader) FirstEmitterCrossKeypair() (emitterCrossKeyPair, error) {
	return r.crossEmitterKeys[0], nil
}

type mockNodeReader struct {
	nodes []node
}

func (db mockNodeReader) Reachables() (reachables []node, err error) {
	for _, n := range db.nodes {
		if n.IsReachable() {
			reachables = append(reachables, n)
		}
	}
	return
}

func (db mockNodeReader) Unreachables() (unreachables []node, err error) {
	for _, n := range db.nodes {
		if !n.IsReachable() {
			unreachables = append(unreachables, n)
		}
	}
	return
}

func (db mockNodeReader) CountReachables() (nb int, err error) {
	for _, n := range db.nodes {
		if n.IsReachable() {
			nb++
		}
	}
	return
}

func (db *mockNodeReader) FindByPublicKey(publicKey publicKey) (found node, err error) {
	for _, n := range db.nodes {
		if bytes.Equal(n.PublicKey().Marshal(), publicKey.Marshal()) {
			return n, nil
		}
	}
	return
}

type mockNode struct {
	ip          net.IP
	port        int
	patchNb     int
	publicKey   publicKey
	isReachable bool
}

func (n mockNode) IP() net.IP {
	return n.ip
}
func (n mockNode) Port() int {
	return n.port
}
func (n mockNode) PatchNumber() int {
	return n.patchNb
}
func (n mockNode) PublicKey() publicKey {
	return n.publicKey
}

func (n mockNode) IsReachable() bool {
	return n.isReachable
}

type mockPublicKey struct {
	bytes []byte
	curve int
}

func (pb mockPublicKey) Marshal() []byte {
	out := make([]byte, 1+len(pb.bytes))
	out[0] = byte(int(pb.curve))
	copy(out[1:], pb.bytes)
	return out
}

func (pb mockPublicKey) Verify(data []byte, sig []byte) (bool, error) {
	return ed25519.Verify(pb.bytes, data, sig), nil
}

type mockPrivateKey struct {
	bytes []byte
}

func (pv mockPrivateKey) Sign(data []byte) ([]byte, error) {
	return ed25519.Sign(pv.bytes, data), nil
}

type mockTransaction struct {
	addr      []byte
	txType    int
	data      map[string]interface{}
	timestamp time.Time
	pubKey    interface{}
	sig       []byte
	originSig []byte
	coordStmp interface{}
	crossB    []interface{}
}

func (t mockTransaction) Address() []byte {
	return t.addr
}
func (t mockTransaction) Type() int {
	return t.txType
}
func (t mockTransaction) Data() map[string]interface{} {
	return t.data
}
func (t mockTransaction) Timestamp() time.Time {
	return t.timestamp
}
func (t mockTransaction) PreviousPublicKey() interface{} {
	return t.pubKey
}
func (t mockTransaction) Signature() []byte {
	return t.sig
}
func (t mockTransaction) OriginSignature() []byte {
	return t.originSig
}
func (t mockTransaction) CoordinatorStamp() interface{} {
	return t.coordStmp
}
func (t mockTransaction) CrossValidations() []interface{} {
	return t.crossB
}
func (t mockTransaction) MarshalBeforeOriginSignature() ([]byte, error) {
	return json.Marshal(map[string]interface{}{
		"addr":       t.addr,
		"data":       t.data,
		"timestamp":  t.timestamp.Unix(),
		"type":       t.txType,
		"public_key": t.pubKey.(publicKey).Marshal(),
		"signature":  t.sig,
	})
}
func (t mockTransaction) MarshalRoot() ([]byte, error) {
	return json.Marshal(map[string]interface{}{
		"addr":             t.addr,
		"data":             t.data,
		"timestamp":        t.timestamp.Unix(),
		"type":             t.txType,
		"public_key":       t.pubKey.(publicKey).Marshal(),
		"signature":        t.sig,
		"origin_signature": t.originSig,
	})
}

type mockCoordinatorStamp struct {
	prevCrossV [][]byte
	pow        interface{}
	stmp       interface{}
	txHash     []byte
	coordN     interface{}
	crossVN    interface{}
	storN      interface{}
}

func (s mockCoordinatorStamp) PreviousCrossValidators() [][]byte {
	return s.prevCrossV
}
func (s mockCoordinatorStamp) ProofOfWork() interface{} {
	return s.pow
}
func (s mockCoordinatorStamp) ValidationStamp() interface{} {
	return s.stmp
}
func (s mockCoordinatorStamp) TransactionHash() []byte {
	return s.txHash
}
func (s mockCoordinatorStamp) ElectedCoordinatorNodes() interface{} {
	return s.coordN
}
func (s mockCoordinatorStamp) ElectedCrossValidationNodes() interface{} {
	return s.crossVN
}
func (s mockCoordinatorStamp) ElectedStorageNodes() interface{} {
	return s.storN
}

type mockValidationStamp struct {
	status    int
	timestamp time.Time
	nodePub   publicKey
	sig       []byte
}

func (v mockValidationStamp) Status() int {
	return v.status
}
func (v mockValidationStamp) Timestamp() time.Time {
	return v.timestamp
}
func (v mockValidationStamp) NodePublicKey() interface{} {
	return v.nodePub
}
func (v mockValidationStamp) NodeSignature() []byte {
	return v.sig
}
