package consensus

// import (
// 	"crypto/rand"
// 	"github.com/uniris/uniris-core/pkg/logging"
// 	"log"
// 	"net"
// 	"os"
// 	"testing"

// 	"encoding/hex"
// 	"encoding/json"
// 	"time"

// 	"github.com/stretchr/testify/assert"
// 	"github.com/uniris/uniris-core/pkg/chain"
// 	"github.com/uniris/uniris-core/pkg/crypto"
// 	"github.com/uniris/uniris-core/pkg/shared"
// )

// /*
// Scenario: request transaction validations
// 	Given a transaction to validate
// 	When I aprop validations to a pool
// 	Then I get validations from them
// */
// func TestRequestValidations(t *testing.T) {
// 	poolR := &mockPoolRequester{}
// 	pv, pub, _ := crypto.GenerateECKeyPair(crypto.Ed25519Curve, rand.Reader)
// 	pubB, _ := pub.Marshal()

// 	v, _ := buildValidation(chain.ValidationOK, pub, pv)
// 	wHeaders := chain.NewWelcomeNodeHeader(pub, []chain.NodeHeader{chain.NewNodeHeader(pub, false, false, 0, true)}, []byte("sig"))
// 	vHeaders := []chain.NodeHeader{chain.NewNodeHeader(pub, false, false, 0, true)}
// 	sHeaders := []chain.NodeHeader{chain.NewNodeHeader(pub, false, false, 0, true)}
// 	mv, _ := chain.NewMasterValidation([]crypto.PublicKey{}, pub, v, wHeaders, vHeaders, sHeaders)

// 	prop, _ := shared.NewEmitterCrossKeyPair([]byte("pvkey"), pub)

// 	data := map[string][]byte{
// 		"encrypted_aes_key":         []byte("aesKey"),
// 		"encrypted_address_by_node": []byte("addr"),
// 		"encrypted_address_by_id":   []byte("addr"),
// 	}

// 	txRaw, _ := json.Marshal(map[string]interface{}{
// 		"addr": hex.EncodeToString(crypto.Hash([]byte("addr"))),
// 		"data": map[string]string{
// 			"encrypted_aes_key":         hex.EncodeToString([]byte("aesKey")),
// 			"encrypted_address_by_node": hex.EncodeToString([]byte("addr")),
// 			"encrypted_address_by_id":   hex.EncodeToString([]byte("addr")),
// 		},
// 		"timestamp":  time.Now().Unix(),
// 		"type":       chain.KeychainTransactionType,
// 		"public_key": hex.EncodeToString(pubB),
// 		"em_shared_keys_proposal": map[string]string{
// 			"encrypted_private_key": hex.EncodeToString([]byte("pvkey")),
// 			"public_key":            hex.EncodeToString(pubB),
// 		},
// 	})
// 	sig, _ := pv.Sign(txRaw)
// 	txSigned, _ := json.Marshal(map[string]interface{}{
// 		"addr": hex.EncodeToString(crypto.Hash([]byte("addr"))),
// 		"data": map[string]string{
// 			"encrypted_aes_key":         hex.EncodeToString([]byte("aesKey")),
// 			"encrypted_address_by_node": hex.EncodeToString([]byte("addr")),
// 			"encrypted_address_by_id":   hex.EncodeToString([]byte("addr")),
// 		},
// 		"timestamp":  time.Now().Unix(),
// 		"type":       chain.KeychainTransactionType,
// 		"public_key": hex.EncodeToString(pubB),
// 		"em_shared_keys_proposal": map[string]string{
// 			"encrypted_private_key": hex.EncodeToString([]byte("pvkey")),
// 			"public_key":            hex.EncodeToString(pubB),
// 		},
// 		"signature": hex.EncodeToString(sig),
// 	})
// 	emSig, _ := pv.Sign(txSigned)
// 	txEmSigned, _ := json.Marshal(map[string]interface{}{
// 		"addr": hex.EncodeToString(crypto.Hash([]byte("addr"))),
// 		"data": map[string]string{
// 			"encrypted_aes_key":         hex.EncodeToString([]byte("aesKey")),
// 			"encrypted_address_by_node": hex.EncodeToString([]byte("addr")),
// 			"encrypted_address_by_id":   hex.EncodeToString([]byte("addr")),
// 		},
// 		"timestamp":  time.Now().Unix(),
// 		"type":       chain.KeychainTransactionType,
// 		"public_key": hex.EncodeToString(pubB),
// 		"em_shared_keys_proposal": map[string]string{
// 			"encrypted_private_key": hex.EncodeToString([]byte("pvkey")),
// 			"public_key":            hex.EncodeToString(pubB),
// 		},
// 		"signature":    hex.EncodeToString(sig),
// 		"em_signature": hex.EncodeToString(emSig),
// 	})
// 	tx, err := chain.NewTransaction(crypto.Hash([]byte("addr")), chain.KeychainTransactionType, data, time.Now(), pub, prop, sig, emSig, crypto.Hash(txEmSigned))

// 	valids, err := requestValidations(tx, mv, Pool{}, 1, poolR)
// 	assert.Nil(t, err)
// 	assert.NotEmpty(t, valids)
// 	assert.Equal(t, chain.ValidationOK, valids[0].Status())
// }

// /*
// Scenario: Create a node validation
// 	Given a validation status
// 	When I want to create node validation
// 	Then I get a validation signed
// */
// func TestBuildValidation(t *testing.T) {
// 	pv, pub, _ := crypto.GenerateECKeyPair(crypto.Ed25519Curve, rand.Reader)

// 	v, err := buildValidation(chain.ValidationOK, pub, pv)
// 	assert.Nil(t, err)
// 	assert.Equal(t, pub, v.PublicKey())
// 	assert.Nil(t, err)
// 	assert.Equal(t, time.Now().Unix(), v.Timestamp().Unix())
// 	assert.Equal(t, chain.ValidationOK, v.Status())
// 	ok, err := v.IsValid()
// 	assert.True(t, ok)
// }

// /*
// Scenario: Validate an incoming transaction
// 	Given a valid transaction
// 	When I want to valid the transaction
// 	Then I get a validation with status OK
// */
// func TestValidateTransaction(t *testing.T) {
// 	pv, pub, _ := crypto.GenerateECKeyPair(crypto.Ed25519Curve, rand.Reader)
// 	pubB, _ := pub.Marshal()

// 	prop, _ := shared.NewEmitterCrossKeyPair([]byte("pvkey"), pub)

// 	data := map[string][]byte{
// 		"encrypted_aes_key":         []byte("aesKey"),
// 		"encrypted_address_by_node": []byte("addr"),
// 		"encrypted_address_by_id":   []byte("addr"),
// 	}

// 	txRaw, _ := json.Marshal(map[string]interface{}{
// 		"addr": hex.EncodeToString(crypto.Hash([]byte("addr"))),
// 		"data": map[string]string{
// 			"encrypted_aes_key":         hex.EncodeToString([]byte("aesKey")),
// 			"encrypted_address_by_node": hex.EncodeToString([]byte("addr")),
// 			"encrypted_address_by_id":   hex.EncodeToString([]byte("addr")),
// 		},
// 		"timestamp":  time.Now().Unix(),
// 		"type":       chain.KeychainTransactionType,
// 		"public_key": hex.EncodeToString(pubB),
// 		"em_shared_keys_proposal": map[string]string{
// 			"encrypted_private_key": hex.EncodeToString([]byte("pvkey")),
// 			"public_key":            hex.EncodeToString(pubB),
// 		},
// 	})

// 	sig, _ := pv.Sign(txRaw)
// 	txSigned, _ := json.Marshal(map[string]interface{}{
// 		"addr": hex.EncodeToString(crypto.Hash([]byte("addr"))),
// 		"data": map[string]string{
// 			"encrypted_aes_key":         hex.EncodeToString([]byte("aesKey")),
// 			"encrypted_address_by_node": hex.EncodeToString([]byte("addr")),
// 			"encrypted_address_by_id":   hex.EncodeToString([]byte("addr")),
// 		},
// 		"timestamp":  time.Now().Unix(),
// 		"type":       chain.KeychainTransactionType,
// 		"public_key": hex.EncodeToString(pubB),
// 		"em_shared_keys_proposal": map[string]string{
// 			"encrypted_private_key": hex.EncodeToString([]byte("pvkey")),
// 			"public_key":            hex.EncodeToString(pubB),
// 		},
// 		"signature": hex.EncodeToString(sig),
// 	})
// 	emSig, _ := pv.Sign(txSigned)
// 	txEmSigned, _ := json.Marshal(map[string]interface{}{
// 		"addr": hex.EncodeToString(crypto.Hash([]byte("addr"))),
// 		"data": map[string]string{
// 			"encrypted_aes_key":         hex.EncodeToString([]byte("aesKey")),
// 			"encrypted_address_by_node": hex.EncodeToString([]byte("addr")),
// 			"encrypted_address_by_id":   hex.EncodeToString([]byte("addr")),
// 		},
// 		"timestamp":  time.Now().Unix(),
// 		"type":       chain.KeychainTransactionType,
// 		"public_key": hex.EncodeToString(pubB),
// 		"em_shared_keys_proposal": map[string]string{
// 			"encrypted_private_key": hex.EncodeToString([]byte("pvkey")),
// 			"public_key":            hex.EncodeToString(pubB),
// 		},
// 		"signature":    hex.EncodeToString(sig),
// 		"em_signature": hex.EncodeToString(emSig),
// 	})
// 	tx, err := chain.NewTransaction(crypto.Hash([]byte("addr")), chain.KeychainTransactionType, data, time.Now(), pub, prop, sig, emSig, crypto.Hash(txEmSigned))
// 	assert.Nil(t, err)

// 	v, _ := buildValidation(chain.ValidationOK, pub, pv)
// 	wHeaders := chain.NewWelcomeNodeHeader(pub, []chain.NodeHeader{chain.NewNodeHeader(pub, false, false, 0, true)}, []byte("sig"))
// 	vHeaders := []chain.NodeHeader{chain.NewNodeHeader(pub, false, false, 0, true)}
// 	sHeaders := []chain.NodeHeader{chain.NewNodeHeader(pub, false, false, 0, true)}
// 	mv, _ := chain.NewMasterValidation([]crypto.PublicKey{}, pub, v, wHeaders, vHeaders, sHeaders)

// 	l := logging.NewLogger("stdout", log.New(os.Stdout, "", 0), "test", net.ParseIP("127.0.0.1"), logging.ErrorLogLevel)
// 	valid, err := ConfirmTransactionValidation(tx, mv, pub, pv, l)
// 	assert.Nil(t, err)
// 	assert.Equal(t, chain.ValidationOK, valid.Status())
// }

// /*
// Scenario: Validate an incoming transaction with invalid integrity
// 	Given a transaction with invalid transaction hash or signature
// 	When I want to valid the transaction
// 	Then I get a validation with status KO
// */
// func TestValidateTransactionWithBadIntegrity(t *testing.T) {
// 	pv, pub, _ := crypto.GenerateECKeyPair(crypto.Ed25519Curve, rand.Reader)

// 	prop, _ := shared.NewEmitterCrossKeyPair([]byte("pvkey"), pub)

// 	data := map[string][]byte{
// 		"encrypted_aes_key":         []byte("aesKey"),
// 		"encrypted_address_by_node": []byte("addr"),
// 		"encrypted_address_by_id":   []byte("addr"),
// 	}

// 	sig, _ := pv.Sign([]byte("tx"))
// 	tx, _ := chain.NewTransaction(crypto.Hash([]byte("addr")), chain.IDTransactionType, data, time.Now(), pub, prop, sig, sig, crypto.Hash([]byte("hash")))

// 	v, _ := buildValidation(chain.ValidationOK, pub, pv)
// 	wHeaders := chain.NewWelcomeNodeHeader(pub, []chain.NodeHeader{chain.NewNodeHeader(pub, false, false, 0, true)}, []byte("sig"))
// 	vHeaders := []chain.NodeHeader{chain.NewNodeHeader(pub, false, false, 0, true)}
// 	sHeaders := []chain.NodeHeader{chain.NewNodeHeader(pub, false, false, 0, true)}
// 	mv, _ := chain.NewMasterValidation([]crypto.PublicKey{}, pub, v, wHeaders, vHeaders, sHeaders)
// 	l := logging.NewLogger("stdout", log.New(os.Stdout, "", 0), "test", net.ParseIP("127.0.0.1"), logging.ErrorLogLevel)
// 	valid, err := ConfirmTransactionValidation(tx, mv, pub, pv, l)
// 	assert.Nil(t, err)
// 	assert.Equal(t, chain.ValidationKO, valid.Status())
// }

// /*
// Scenario: Perform Proof of work
// 	Given a transaction and em chain keypair stored
// 	When I want to perform the proof of work of this transaction
// 	Then I get the valid public key
// */
// func TestPerformPOW(t *testing.T) {

// 	pv, pub, _ := crypto.GenerateECKeyPair(crypto.Ed25519Curve, rand.Reader)
// 	pubB, _ := pub.Marshal()

// 	keyReader := &mockSharedKeyReader{}
// 	emKP, _ := shared.NewEmitterCrossKeyPair([]byte("pvKey"), pub)
// 	keyReader.crossEmitterKeys = append(keyReader.crossEmitterKeys, emKP)

// 	prop, _ := shared.NewEmitterCrossKeyPair([]byte("pvkey"), pub)

// 	data := map[string][]byte{
// 		"encrypted_aes_key":         []byte("aesKey"),
// 		"encrypted_address_by_node": []byte("addr"),
// 		"encrypted_address_by_id":   []byte("addr"),
// 	}

// 	txRaw, _ := json.Marshal(map[string]interface{}{
// 		"addr": hex.EncodeToString(crypto.Hash([]byte("addr"))),
// 		"data": map[string]string{
// 			"encrypted_aes_key":         hex.EncodeToString([]byte("aesKey")),
// 			"encrypted_address_by_node": hex.EncodeToString([]byte("addr")),
// 			"encrypted_address_by_id":   hex.EncodeToString([]byte("addr")),
// 		},
// 		"timestamp":  time.Now().Unix(),
// 		"type":       chain.KeychainTransactionType,
// 		"public_key": hex.EncodeToString(pubB),
// 		"em_shared_keys_proposal": map[string]string{
// 			"encrypted_private_key": hex.EncodeToString([]byte("pvkey")),
// 			"public_key":            hex.EncodeToString(pubB),
// 		},
// 	})
// 	sig, _ := pv.Sign(txRaw)
// 	txSigned, _ := json.Marshal(map[string]interface{}{
// 		"addr": hex.EncodeToString(crypto.Hash([]byte("addr"))),
// 		"data": map[string]string{
// 			"encrypted_aes_key":         hex.EncodeToString([]byte("aesKey")),
// 			"encrypted_address_by_node": hex.EncodeToString([]byte("addr")),
// 			"encrypted_address_by_id":   hex.EncodeToString([]byte("addr")),
// 		},
// 		"timestamp":  time.Now().Unix(),
// 		"type":       chain.KeychainTransactionType,
// 		"public_key": hex.EncodeToString(pubB),
// 		"em_shared_keys_proposal": map[string]string{
// 			"encrypted_private_key": hex.EncodeToString([]byte("pvkey")),
// 			"public_key":            hex.EncodeToString(pubB),
// 		},
// 		"signature": hex.EncodeToString(sig),
// 	})
// 	emSig, _ := pv.Sign(txSigned)
// 	txEmSigned, _ := json.Marshal(map[string]interface{}{
// 		"addr": hex.EncodeToString(crypto.Hash([]byte("addr"))),
// 		"data": map[string]string{
// 			"encrypted_aes_key":         hex.EncodeToString([]byte("aesKey")),
// 			"encrypted_address_by_node": hex.EncodeToString([]byte("addr")),
// 			"encrypted_address_by_id":   hex.EncodeToString([]byte("addr")),
// 		},
// 		"timestamp":  time.Now().Unix(),
// 		"type":       chain.KeychainTransactionType,
// 		"public_key": hex.EncodeToString(pubB),
// 		"em_shared_keys_proposal": map[string]string{
// 			"encrypted_private_key": hex.EncodeToString([]byte("pvkey")),
// 			"public_key":            hex.EncodeToString(pubB),
// 		},
// 		"signature":    hex.EncodeToString(sig),
// 		"em_signature": hex.EncodeToString(emSig),
// 	})
// 	tx, err := chain.NewTransaction(crypto.Hash([]byte("addr")), chain.KeychainTransactionType, data, time.Now(), pub, prop, sig, emSig, crypto.Hash(txEmSigned))
// 	assert.Nil(t, err)

// 	pow, err := proofOfWork(tx, keyReader)
// 	assert.Nil(t, err)
// 	assert.Equal(t, pub, pow)
// }

// /*
// Scenario: Pre-validate a transaction
// 	Given a transaction
// 	When I want to prevalidate this transaction
// 	Then I get the node validation and the proof of work
// */
// func TestPreValidateTransaction(t *testing.T) {

// 	pv, pub, _ := crypto.GenerateECKeyPair(crypto.Ed25519Curve, rand.Reader)
// 	pubB, _ := pub.Marshal()
// 	sharedKeyReader := &mockSharedKeyReader{}
// 	emKP, _ := shared.NewEmitterCrossKeyPair([]byte("pvkey"), pub)
// 	sharedKeyReader.crossEmitterKeys = append(sharedKeyReader.crossEmitterKeys, emKP)

// 	prop, _ := shared.NewEmitterCrossKeyPair([]byte("pvkey"), pub)

// 	data := map[string][]byte{
// 		"encrypted_aes_key":         []byte("aesKey"),
// 		"encrypted_address_by_node": []byte("addr"),
// 		"encrypted_address_by_id":   []byte("addr"),
// 	}

// 	txRaw, _ := json.Marshal(map[string]interface{}{
// 		"addr": hex.EncodeToString(crypto.Hash([]byte("addr"))),
// 		"data": map[string]string{
// 			"encrypted_aes_key":         hex.EncodeToString([]byte("aesKey")),
// 			"encrypted_address_by_node": hex.EncodeToString([]byte("addr")),
// 			"encrypted_address_by_id":   hex.EncodeToString([]byte("addr")),
// 		},
// 		"timestamp":  time.Now().Unix(),
// 		"type":       chain.KeychainTransactionType,
// 		"public_key": hex.EncodeToString(pubB),
// 		"em_shared_keys_proposal": map[string]string{
// 			"encrypted_private_key": hex.EncodeToString([]byte("pvkey")),
// 			"public_key":            hex.EncodeToString(pubB),
// 		},
// 	})
// 	sig, _ := pv.Sign(txRaw)
// 	txSigned, _ := json.Marshal(map[string]interface{}{
// 		"addr": hex.EncodeToString(crypto.Hash([]byte("addr"))),
// 		"data": map[string]string{
// 			"encrypted_aes_key":         hex.EncodeToString([]byte("aesKey")),
// 			"encrypted_address_by_node": hex.EncodeToString([]byte("addr")),
// 			"encrypted_address_by_id":   hex.EncodeToString([]byte("addr")),
// 		},
// 		"timestamp":  time.Now().Unix(),
// 		"type":       chain.KeychainTransactionType,
// 		"public_key": hex.EncodeToString(pubB),
// 		"em_shared_keys_proposal": map[string]string{
// 			"encrypted_private_key": hex.EncodeToString([]byte("pvkey")),
// 			"public_key":            hex.EncodeToString(pubB),
// 		},
// 		"signature": hex.EncodeToString(sig),
// 	})
// 	emSig, _ := pv.Sign(txSigned)
// 	txEmSigned, _ := json.Marshal(map[string]interface{}{
// 		"addr": hex.EncodeToString(crypto.Hash([]byte("addr"))),
// 		"data": map[string]string{
// 			"encrypted_aes_key":         hex.EncodeToString([]byte("aesKey")),
// 			"encrypted_address_by_node": hex.EncodeToString([]byte("addr")),
// 			"encrypted_address_by_id":   hex.EncodeToString([]byte("addr")),
// 		},
// 		"timestamp":  time.Now().Unix(),
// 		"type":       chain.KeychainTransactionType,
// 		"public_key": hex.EncodeToString(pubB),
// 		"em_shared_keys_proposal": map[string]string{
// 			"encrypted_private_key": hex.EncodeToString([]byte("pvkey")),
// 			"public_key":            hex.EncodeToString(pubB),
// 		},
// 		"signature":    hex.EncodeToString(sig),
// 		"em_signature": hex.EncodeToString(emSig),
// 	})
// 	tx, err := chain.NewTransaction(crypto.Hash([]byte("addr")), chain.KeychainTransactionType, data, time.Now(), pub, prop, sig, emSig, crypto.Hash(txEmSigned))
// 	assert.Nil(t, err)

// 	wHeaders := chain.NewWelcomeNodeHeader(pub, []chain.NodeHeader{chain.NewNodeHeader(pub, false, false, 0, true)}, []byte("sig"))

// 	sPool := Pool{Node{publicKey: pub}}
// 	vPool := Pool{Node{publicKey: pub}}
// 	lastVPool := Pool{}

// 	mv, err := preValidateTransaction(tx, wHeaders, sPool, vPool, lastVPool, pub, pv, sharedKeyReader, mockNodeReader{})
// 	assert.Nil(t, err)
// 	assert.Equal(t, pub, mv.ProofOfWork())
// 	assert.EqualValues(t, pub, mv.Validation().PublicKey())
// 	assert.Equal(t, chain.ValidationOK, mv.Validation().Status())
// 	ok, err := mv.Validation().IsValid()
// 	assert.True(t, ok)
// 	assert.Nil(t, err)

// 	assert.Len(t, mv.ValidationHeaders(), 2) //vPool member + master node
// }

// /*
// Scenario: Lead transaction mining
// 	Given a valid transaction
// 	When I want to lead its mining
// 	Then the transaction is mined and stored
// */
// func TestLeadMining(t *testing.T) {

// 	pv, pub, _ := crypto.GenerateECKeyPair(crypto.Ed25519Curve, rand.Reader)
// 	pubB, _ := pub.Marshal()
// 	sharedKeyReader := &mockSharedKeyReader{}
// 	emKP, _ := shared.NewEmitterCrossKeyPair([]byte("pvkey"), pub)
// 	nodeKP, _ := shared.NewNodeCrossKeyPair(pub, pv)

// 	sharedKeyReader.crossEmitterKeys = append(sharedKeyReader.crossEmitterKeys, emKP)
// 	sharedKeyReader.crossNodeKeys = append(sharedKeyReader.crossNodeKeys, nodeKP)
// 	sharedKeyReader.authKeys = append(sharedKeyReader.authKeys, pub)

// 	nodeReader := &mockNodeReader{
// 		nodes: []Node{
// 			Node{
// 				publicKey:   pub,
// 				isReachable: true,
// 			},
// 		},
// 	}

// 	prop, _ := shared.NewEmitterCrossKeyPair([]byte("pvkey"), pub)

// 	data := map[string][]byte{
// 		"encrypted_aes_key":         []byte("aesKey"),
// 		"encrypted_address_by_node": []byte("addr"),
// 		"encrypted_address_by_id":   []byte("addr"),
// 	}

// 	txRaw, _ := json.Marshal(map[string]interface{}{
// 		"addr": hex.EncodeToString(crypto.Hash([]byte("addr"))),
// 		"data": map[string]string{
// 			"encrypted_aes_key":         hex.EncodeToString([]byte("aesKey")),
// 			"encrypted_address_by_node": hex.EncodeToString([]byte("addr")),
// 			"encrypted_address_by_id":   hex.EncodeToString([]byte("addr")),
// 		},
// 		"timestamp":  time.Now().Unix(),
// 		"type":       chain.KeychainTransactionType,
// 		"public_key": hex.EncodeToString(pubB),
// 		"em_shared_keys_proposal": map[string]string{
// 			"encrypted_private_key": hex.EncodeToString([]byte("pvkey")),
// 			"public_key":            hex.EncodeToString(pubB),
// 		},
// 	})
// 	sig, _ := pv.Sign(txRaw)
// 	txSigned, _ := json.Marshal(map[string]interface{}{
// 		"addr": hex.EncodeToString(crypto.Hash([]byte("addr"))),
// 		"data": map[string]string{
// 			"encrypted_aes_key":         hex.EncodeToString([]byte("aesKey")),
// 			"encrypted_address_by_node": hex.EncodeToString([]byte("addr")),
// 			"encrypted_address_by_id":   hex.EncodeToString([]byte("addr")),
// 		},
// 		"timestamp":  time.Now().Unix(),
// 		"type":       chain.KeychainTransactionType,
// 		"public_key": hex.EncodeToString(pubB),
// 		"em_shared_keys_proposal": map[string]string{
// 			"encrypted_private_key": hex.EncodeToString([]byte("pvkey")),
// 			"public_key":            hex.EncodeToString(pubB),
// 		},
// 		"signature": hex.EncodeToString(sig),
// 	})
// 	emSig, _ := pv.Sign(txSigned)
// 	txEmSigned, _ := json.Marshal(map[string]interface{}{
// 		"addr": hex.EncodeToString(crypto.Hash([]byte("addr"))),
// 		"data": map[string]string{
// 			"encrypted_aes_key":         hex.EncodeToString([]byte("aesKey")),
// 			"encrypted_address_by_node": hex.EncodeToString([]byte("addr")),
// 			"encrypted_address_by_id":   hex.EncodeToString([]byte("addr")),
// 		},
// 		"timestamp":  time.Now().Unix(),
// 		"type":       chain.KeychainTransactionType,
// 		"public_key": hex.EncodeToString(pubB),
// 		"em_shared_keys_proposal": map[string]string{
// 			"encrypted_private_key": hex.EncodeToString([]byte("pvkey")),
// 			"public_key":            hex.EncodeToString(pubB),
// 		},
// 		"signature":    hex.EncodeToString(sig),
// 		"em_signature": hex.EncodeToString(emSig),
// 	})
// 	tx, err := chain.NewTransaction(crypto.Hash([]byte("addr")), chain.KeychainTransactionType, data, time.Now(), pub, prop, sig, emSig, crypto.Hash(txEmSigned))
// 	assert.Nil(t, err)
// 	poolR := &mockPoolRequester{}
// 	l := logging.NewLogger("stdout", log.New(os.Stdout, "", 0), "test", net.ParseIP("127.0.0.1"), logging.ErrorLogLevel)
// 	wHeaders := chain.NewWelcomeNodeHeader(pub, []chain.NodeHeader{chain.NewNodeHeader(pub, false, false, 0, true)}, []byte("sig"))
// 	assert.Nil(t, LeadMining(tx, 1, wHeaders, poolR, pub, pv, sharedKeyReader, nodeReader, l))

// 	time.Sleep(1 * time.Second)

// 	assert.Len(t, poolR.stores, 1)
// }

// /*
// Scenario: Get the minimum validation number for a system transaction with tiny network
// 	Given a system transaction with 1 nodes and 1 reachable
// 	When I want to get the validation required number
// 	Then I get 1
// */
// func TestTestValidationNumberNetworkBasedWithTinyNetwork(t *testing.T) {
// 	nbValidations, err := requiredValidationNumberForSysTX(1, 1)
// 	assert.Nil(t, err)
// 	assert.Equal(t, 1, nbValidations)
// }

// /*
// Scenario: Get the minimum validation number for a system transaction with small network
// 	Given a system transaction with less than 5 nodes and 2 reachables
// 	When I want to get the validation required number
// 	Then I get 2
// */
// func TestTestValidationNumberNetworkBasedWithSmallNetwork(t *testing.T) {

// 	nbValidations, err := requiredValidationNumberForSysTX(5, 2)
// 	assert.Nil(t, err)
// 	assert.Equal(t, 2, nbValidations)
// }

// /*
// Scenario: Get the minimum validation number for a system transaction with normal network
// 	Given a system transaction with 10 nodes and 10 reachables
// 	When I want to get the validation required number
// 	Then I get 2
// */
// func TestValidationNumberNetworkBasedWithNormalNetwork(t *testing.T) {
// 	nbValidations, err := requiredValidationNumberForSysTX(10, 10)
// 	assert.Nil(t, err)
// 	assert.Equal(t, 5, nbValidations)
// }

// /*
// Scenario: Get the minimum validation number for a system transaction with too less nodes
// 	Given a system transaction with 5 nodes and 1 reachable
// 	When I want to get the validation required number
// 	Then I get an error
// */
// func TestValidationNumberNetworkBasedWithUnsufficientNetwork(t *testing.T) {
// 	_, err := requiredValidationNumberForSysTX(6, 1)
// 	assert.EqualError(t, err, "no enough nodes in the network to validate this transaction")
// }

// /*
// SCenario: Get the minimum validation number for transaction with fees
// 	Given a transaction fees as 1 UCO and 10 reachables nodes
// 	When I want to get the validation required number
// 	Then I get 9 validations neeed
// */
// func TestValidationNumberFeesBasedFor1UCOFeesWith10Nodes(t *testing.T) {
// 	assert.Equal(t, 9, requiredValidationNumberWithFees(1, 10))
// }

// /*
// SCenario: Get the minimum validation number for transaction with fees
// 	Given a transaction fees as 1 UCO and 8 reachables nodes
// 	When I want to get the validation required number
// 	Then I get 9 validations neeed
// */
// func TestValidationNumberFeesBasedFor1UCOFeesWith8Nodes(t *testing.T) {
// 	assert.Equal(t, 8, requiredValidationNumberWithFees(1, 8))
// }

// /*
// Scenario: Get the minimum validation for normal transaction with less 3 nodes
// 	Given a system transaction and less 3 nodes
// 	When I want to get the validation required number
// 	Then I get an error
// */
// func TestRequiredValidationNumberNormalTxWithLess3Nodes(t *testing.T) {

// 	nodeReader := mockNodeReader{
// 		nodes: []Node{
// 			Node{isReachable: true},
// 			Node{isReachable: true},
// 		},
// 	}

// 	keyReader := mockSharedKeyReader{
// 		authKeys: make([]crypto.PublicKey, 2),
// 	}

// 	_, err := RequiredValidationNumber(chain.KeychainTransactionType, 0.001, nodeReader, keyReader)
// 	assert.EqualError(t, err, "no enough nodes in the network to validate this transaction")
// }

// /*
// Scenario: Get the minimum validation for system transaction
// 	Given a system transaction and 5 nodes and 5 reachables
// 	When I want to get the validation required number
// 	Then I get 5 validation
// */
// func TestRequiredValidationNumberSystemTxWith5Nodes(t *testing.T) {
// 	nodeReader := mockNodeReader{
// 		nodes: []Node{
// 			Node{isReachable: true},
// 			Node{isReachable: true},
// 			Node{isReachable: true},
// 			Node{isReachable: true},
// 			Node{isReachable: true},
// 		},
// 	}

// 	keyReader := mockSharedKeyReader{
// 		authKeys: make([]crypto.PublicKey, 5),
// 	}

// 	nbValidations, err := RequiredValidationNumber(chain.SystemTransactionType, 0, nodeReader, keyReader)
// 	assert.Nil(t, err)
// 	assert.Equal(t, 5, nbValidations)
// }

// /*
// Scenario: Get the minimum validation  for normal transaction
// 	Given a normal transaction (minium fees: 0.001 => 3 validations)
// 	When I want to get the validation required number
// 	Then I get 3 validation
// */
// func TestRequiredValidationNumberWith5UCO(t *testing.T) {

// 	nodeReader := mockNodeReader{
// 		nodes: []Node{
// 			Node{isReachable: true},
// 			Node{isReachable: true},
// 			Node{isReachable: true},
// 			Node{isReachable: true},
// 			Node{isReachable: true},
// 		},
// 	}

// 	keyReader := mockSharedKeyReader{
// 		authKeys: make([]crypto.PublicKey, 5),
// 	}

// 	nbValidations, err := RequiredValidationNumber(chain.ContractTransactionType, 0.001, nodeReader, keyReader)
// 	assert.Nil(t, err)
// 	assert.Equal(t, 3, nbValidations)
// }

// type mockPoolRequester struct {
// 	stores []chain.Transaction
// 	ko     []chain.Transaction
// }

// func (pr mockPoolRequester) RequestLastTransaction(pool Pool, txAddr crypto.VersionnedHash, txType chain.TransactionType) (*chain.Transaction, error) {
// 	return nil, nil
// }

// func (pr mockPoolRequester) RequestTransactionTimeLock(pool Pool, txHash crypto.VersionnedHash, txAddr crypto.VersionnedHash, masterPublicKey crypto.PublicKey) error {
// 	return nil
// }

// func (pr mockPoolRequester) RequestTransactionValidations(pool Pool, tx chain.Transaction, minValid int, masterValid chain.MasterValidation) ([]chain.Validation, error) {
// 	pv, pub, _ := crypto.GenerateECKeyPair(crypto.Ed25519Curve, rand.Reader)

// 	v, _ := buildValidation(chain.ValidationOK, pub, pv)
// 	return []chain.Validation{v}, nil
// }

// func (pr *mockPoolRequester) RequestTransactionStorage(pool Pool, minReplicas int, tx chain.Transaction) error {
// 	pr.stores = append(pr.stores, tx)
// 	return nil
// }
