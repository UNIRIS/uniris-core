package rpc

import (
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"log"
	"net"
	"os"
	"path/filepath"
	"testing"
	"time"

	"golang.org/x/crypto/ed25519"

	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc"

	api "github.com/uniris/uniris-core/api/protobuf-spec"
	"github.com/uniris/uniris-core/pkg/logging"
)

func TestMain(m *testing.M) {
	dir, _ := os.Getwd()
	os.Setenv("PLUGINS_DIR", filepath.Join(dir, "../../plugins"))
	m.Run()
}

/*
Scenario: Request transction lock on a pool
	Given a transaction to lock
	When I request to lock it
	Then the lock is stored in the database
*/
func TestRequestTransactionLock(t *testing.T) {

	pub, pv, _ := ed25519.GenerateKey(rand.Reader)

	sharedKeyReader := &mockSharedKeyReader{
		crossNodePubKeys: []publicKey{
			mockPublicKey{bytes: pub},
		},
		crossNodePvKeys: []privateKey{
			mockPrivateKey{bytes: pv},
		},
	}

	nodeReader := &mockNodeReader{
		nodes: []node{
			mockNode{
				ip:        net.ParseIP("127.0.0.1"),
				port:      5000,
				publicKey: mockPublicKey{bytes: pub},
				patchNb:   1,
			},
		},
	}

	chainDB := &mockChainDB{}

	l := logging.NewLogger("stdout", log.New(os.Stdout, "", 0), "test", net.ParseIP("127.0.0.1"), logging.ErrorLogLevel)
	pr := PoolRequester{
		Logger:          l,
		SharedKeyReader: sharedKeyReader,
		nodeReader:      nodeReader,
	}

	txSrv := NewTransactionService(chainDB, nil, sharedKeyReader, nodeReader, pr, mockPublicKey{bytes: pub}, mockPrivateKey{bytes: pv}, l)

	lis, _ := net.Listen("tcp", ":5000")
	defer lis.Close()
	grpcServer := grpc.NewServer()
	api.RegisterTransactionServiceServer(grpcServer, txSrv)
	go grpcServer.Serve(lis)

	pool := mockElectedNodeList{
		nodes: []interface{}{
			mockElectedNode{
				publicKey: mockPublicKey{bytes: pub},
			},
		},
	}

	assert.Nil(t, pr.RequestTransactionTimeLock(pool, []byte("tx hash"), []byte("addr"), mockPublicKey{bytes: pub}))
}

/*
Scenario: Request transaction validation confirmation
	Given a transaction to validate
	When I request to confirm the validation
	Then I get a validation
*/
func TestRequestConfirmValidation(t *testing.T) {

	pub, pv, _ := ed25519.GenerateKey(rand.Reader)

	sharedKeyReader := &mockSharedKeyReader{
		crossNodePubKeys: []publicKey{
			mockPublicKey{bytes: pub},
		},
		crossNodePvKeys: []privateKey{
			mockPrivateKey{bytes: pv},
		},
	}

	nodeReader := &mockNodeReader{
		nodes: []node{
			mockNode{
				ip:        net.ParseIP("127.0.0.1"),
				port:      5000,
				publicKey: mockPublicKey{bytes: pub},
				patchNb:   1,
			},
		},
	}

	l := logging.NewLogger("stdout", log.New(os.Stdout, "", 0), "test", net.ParseIP("127.0.0.1"), logging.ErrorLogLevel)

	pr := PoolRequester{
		SharedKeyReader: sharedKeyReader,
		Logger:          l,
		nodeReader:      nodeReader,
	}

	miningSrv := NewTransactionService(nil, nil, sharedKeyReader, nodeReader, pr, mockPublicKey{bytes: pub}, mockPrivateKey{bytes: pv}, l)

	lis, err := net.Listen("tcp", ":5000")
	defer lis.Close()
	grpcServer := grpc.NewServer()
	api.RegisterTransactionServiceServer(grpcServer, miningSrv)
	go grpcServer.Serve(lis)

	txRaw := map[string]interface{}{
		"addr": []byte("addr"),
		"data": map[string]interface{}{
			"encrypted_address_by_node": []byte("addr"),
			"encrypted_wallet":          []byte("wallet"),
			"em_shared_keys_proposal": map[string]interface{}{
				"encrypted_private_key": []byte("pvkey"),
				"public_key":            pub,
			},
		},
		"timestamp":  time.Now().Unix(),
		"type":       0,
		"public_key": pub,
	}
	txBytes, _ := json.Marshal(txRaw)
	sig := ed25519.Sign(pv, txBytes)
	txRaw["signature"] = hex.EncodeToString(sig)
	txByteWithSig, _ := json.Marshal(txRaw)
	emSig := ed25519.Sign(pv, txByteWithSig)
	txRaw["em_signature"] = hex.EncodeToString(emSig)
	txBytes, _ = json.Marshal(txRaw)

	tx := mockTransaction{
		addr:   []byte("addr"),
		txType: 0,
		data: map[string]interface{}{
			"encrypted_address_by_node": []byte("addr"),
			"encrypted_wallet":          []byte("wallet"),
			"em_shared_keys_proposal": map[string]interface{}{
				"encrypted_private_key": []byte("pvkey"),
				"public_key":            pub,
			},
		},
		timestamp: time.Now(),
		pubKey:    mockPublicKey{bytes: pub},
		sig:       sig,
		originSig: emSig,
	}

	v := mockValidationStamp{
		nodePub:   mockPublicKey{bytes: pub},
		status:    1,
		timestamp: time.Now(),
	}
	vBytes, _ := json.Marshal(map[string]interface{}{
		"status":     1,
		"public_key": v.nodePub.Marshal(),
		"timestamp":  v.timestamp.Unix(),
	})
	vSig := ed25519.Sign(pv, vBytes)
	v.sig = vSig

	coordN := mockElectedNodeList{
		nodes: []interface{}{
			mockElectedNode{
				publicKey: mockPublicKey{bytes: pub},
			},
		},
		pubK: mockPublicKey{bytes: pub},
	}
	coorNB, _ := json.Marshal(coordN.nodes)
	coordN.sig = ed25519.Sign(pv, coorNB)

	crossN := mockElectedNodeList{
		nodes: []interface{}{
			mockElectedNode{
				publicKey: mockPublicKey{bytes: pub},
			},
		},
		pubK: mockPublicKey{bytes: pub},
	}
	crossNB, _ := json.Marshal(crossN.nodes)
	crossN.sig = ed25519.Sign(pv, crossNB)

	storN := mockElectedNodeList{
		nodes: []interface{}{
			mockElectedNode{
				publicKey: mockPublicKey{bytes: pub},
			},
		},
		pubK: mockPublicKey{bytes: pub},
	}
	storNB, _ := json.Marshal(storN.nodes)
	storN.sig = ed25519.Sign(pv, storNB)

	coordStmp := mockCoordinatorStamp{
		coordN:     coordN,
		crossVN:    crossN,
		storN:      storN,
		pow:        mockPublicKey{bytes: pub},
		stmp:       v,
		txHash:     []byte("hash"),
		prevCrossV: nil,
	}

	pool := mockElectedNodeList{
		nodes: []interface{}{
			mockElectedNode{
				publicKey: mockPublicKey{bytes: pub},
			},
		},
	}

	valids, err := pr.RequestTransactionValidations(pool, tx, 1, coordStmp)
	assert.Nil(t, err)

	assert.Len(t, valids, 1)
	assert.Equal(t, pub, valids[0].NodePublicKey().(publicKey).Marshal())
	assert.Equal(t, 1, valids[0].Status())
}

/*
Scenario: Request transaction store
	Given a transaction to store
	When I request to store the validation
	Then the transaction is stored
*/
func TestRequestStorage(t *testing.T) {

	pub, pv, _ := ed25519.GenerateKey(rand.Reader)

	sharedKeyReader := &mockSharedKeyReader{
		crossNodePubKeys: []publicKey{
			mockPublicKey{bytes: pub},
		},
		crossNodePvKeys: []privateKey{
			mockPrivateKey{bytes: pv},
		},
	}

	nodeReader := &mockNodeReader{
		nodes: []node{
			mockNode{
				ip:        net.ParseIP("127.0.0.1"),
				port:      5000,
				publicKey: mockPublicKey{bytes: pub},
				patchNb:   1,
			},
		},
	}

	l := logging.NewLogger("stdout", log.New(os.Stdout, "", 0), "test", net.ParseIP("127.0.0.1"), logging.ErrorLogLevel)

	pr := PoolRequester{
		SharedKeyReader: sharedKeyReader,
		Logger:          l,
		nodeReader:      nodeReader,
	}

	chainDB := &mockChainDB{}

	txSrv := NewTransactionService(chainDB, nil, sharedKeyReader, nodeReader, pr, mockPublicKey{bytes: pub}, mockPrivateKey{bytes: pv}, l)

	lis, _ := net.Listen("tcp", ":5000")
	defer lis.Close()
	grpcServer := grpc.NewServer()
	api.RegisterTransactionServiceServer(grpcServer, txSrv)
	go grpcServer.Serve(lis)

	v := mockValidationStamp{
		nodePub:   mockPublicKey{bytes: pub},
		status:    1,
		timestamp: time.Now(),
	}
	vBytes, _ := json.Marshal(map[string]interface{}{
		"status":     1,
		"public_key": v.nodePub.Marshal(),
		"timestamp":  v.timestamp.Unix(),
	})
	vSig := ed25519.Sign(pv, vBytes)
	v.sig = vSig

	coordN := mockElectedNodeList{
		nodes: []interface{}{
			mockElectedNode{
				publicKey: mockPublicKey{bytes: pub},
			},
		},
		pubK: mockPublicKey{bytes: pub},
	}
	coorNB, _ := json.Marshal(coordN.nodes)
	coordN.sig = ed25519.Sign(pv, coorNB)

	crossN := mockElectedNodeList{
		nodes: []interface{}{
			mockElectedNode{
				publicKey: mockPublicKey{bytes: pub},
			},
		},
		pubK: mockPublicKey{bytes: pub},
	}
	crossNB, _ := json.Marshal(crossN.nodes)
	crossN.sig = ed25519.Sign(pv, crossNB)

	storN := mockElectedNodeList{
		nodes: []interface{}{
			mockElectedNode{
				publicKey: mockPublicKey{bytes: pub},
			},
		},
		pubK: mockPublicKey{bytes: pub},
	}
	storNB, _ := json.Marshal(storN.nodes)
	storN.sig = ed25519.Sign(pv, storNB)

	coordStmp := mockCoordinatorStamp{
		coordN:     coordN,
		crossVN:    crossN,
		storN:      storN,
		pow:        mockPublicKey{bytes: pub},
		stmp:       v,
		txHash:     []byte("hash"),
		prevCrossV: nil,
	}

	txRaw := map[string]interface{}{
		"addr": []byte("addr"),
		"data": map[string]interface{}{
			"encrypted_address_by_node": []byte("addr"),
			"encrypted_wallet":          []byte("wallet"),
			"em_shared_keys_proposal": map[string]interface{}{
				"encrypted_private_key": []byte("pvkey"),
				"public_key":            pub,
			},
		},
		"timestamp":  time.Now().Unix(),
		"type":       0,
		"public_key": pub,
	}
	txBytes, _ := json.Marshal(txRaw)
	sig := ed25519.Sign(pv, txBytes)
	txRaw["signature"] = sig
	txByteWithSig, _ := json.Marshal(txRaw)
	emSig := ed25519.Sign(pv, txByteWithSig)
	txRaw["em_signature"] = emSig
	txBytes, _ = json.Marshal(txRaw)

	tx := mockTransaction{
		addr: []byte("addr"),
		data: map[string]interface{}{
			"encrypted_address_by_node": []byte("addr"),
			"encrypted_wallet":          []byte("wallet"),
			"em_shared_keys_proposal": map[string]interface{}{
				"encrypted_private_key": []byte("pvkey"),
				"public_key":            pub,
			},
		},
		timestamp: time.Now(),
		txType:    0,
		pubKey:    mockPublicKey{bytes: pub},
		coordStmp: coordStmp,
		crossB:    []interface{}{v},
	}

	pool := mockElectedNodeList{
		nodes: []interface{}{
			mockElectedNode{
				publicKey: mockPublicKey{bytes: pub},
			},
		},
	}

	assert.Nil(t, pr.RequestTransactionStorage(pool, 1, tx))

	assert.Len(t, chainDB.keychains, 1)
	assert.EqualValues(t, txBytes, chainDB.keychains[0].CoordinatorStamp().(coordinatorStamp).TransactionHash())
}

/*
Scenario: Send request to get last transaction
	Given a keychain transaction stored
	When I want to request a node to get the last transaction from the address
	Then I get the last transaction
*/
func TestSendGetLastTransaction(t *testing.T) {
	pub, pv, _ := ed25519.GenerateKey(rand.Reader)

	sharedKeyReader := &mockSharedKeyReader{
		crossNodePubKeys: []publicKey{
			mockPublicKey{bytes: pub},
		},
		crossNodePvKeys: []privateKey{
			mockPrivateKey{bytes: pv},
		},
	}

	nodeReader := &mockNodeReader{
		nodes: []node{
			mockNode{
				ip:        net.ParseIP("127.0.0.1"),
				port:      5000,
				publicKey: mockPublicKey{bytes: pub},
				patchNb:   1,
			},
		},
	}

	l := logging.NewLogger("stdout", log.New(os.Stdout, "", 0), "test", net.ParseIP("127.0.0.1"), logging.ErrorLogLevel)

	pr := PoolRequester{
		SharedKeyReader: sharedKeyReader,
		Logger:          l,
		nodeReader:      nodeReader,
	}

	chainDB := &mockChainDB{}
	indexDB := &mockIndexDB{
		rows: make(map[string][]byte, 0),
	}

	txSrv := NewTransactionService(chainDB, indexDB, sharedKeyReader, nodeReader, pr, mockPublicKey{bytes: pub}, mockPrivateKey{bytes: pv}, l)

	lis, _ := net.Listen("tcp", ":5000")
	defer lis.Close()
	grpcServer := grpc.NewServer()
	api.RegisterTransactionServiceServer(grpcServer, txSrv)
	go grpcServer.Serve(lis)

	v := mockValidationStamp{
		nodePub:   mockPublicKey{bytes: pub},
		status:    1,
		timestamp: time.Now(),
	}
	vBytes, _ := json.Marshal(map[string]interface{}{
		"status":     1,
		"public_key": v.nodePub.Marshal(),
		"timestamp":  v.timestamp.Unix(),
	})
	vSig := ed25519.Sign(pv, vBytes)
	v.sig = vSig

	coordN := mockElectedNodeList{
		nodes: []interface{}{
			mockElectedNode{
				publicKey: mockPublicKey{bytes: pub},
			},
		},
		pubK: mockPublicKey{bytes: pub},
	}
	coorNB, _ := json.Marshal(coordN.nodes)
	coordN.sig = ed25519.Sign(pv, coorNB)

	crossN := mockElectedNodeList{
		nodes: []interface{}{
			mockElectedNode{
				publicKey: mockPublicKey{bytes: pub},
			},
		},
		pubK: mockPublicKey{bytes: pub},
	}
	crossNB, _ := json.Marshal(crossN.nodes)
	crossN.sig = ed25519.Sign(pv, crossNB)

	storN := mockElectedNodeList{
		nodes: []interface{}{
			mockElectedNode{
				publicKey: mockPublicKey{bytes: pub},
			},
		},
		pubK: mockPublicKey{bytes: pub},
	}
	storNB, _ := json.Marshal(storN.nodes)
	storN.sig = ed25519.Sign(pv, storNB)

	coordStmp := mockCoordinatorStamp{
		coordN:     coordN,
		crossVN:    crossN,
		storN:      storN,
		pow:        mockPublicKey{bytes: pub},
		stmp:       v,
		txHash:     []byte("hash"),
		prevCrossV: nil,
	}

	txRaw := map[string]interface{}{
		"addr": []byte("addr"),
		"data": map[string]interface{}{
			"encrypted_address_by_node": []byte("addr"),
			"encrypted_wallet":          []byte("wallet"),
			"em_shared_keys_proposal": map[string]interface{}{
				"encrypted_private_key": []byte("pvkey"),
				"public_key":            pub,
			},
		},
		"timestamp":  time.Now().Unix(),
		"type":       0,
		"public_key": pub,
	}
	txBytes, _ := json.Marshal(txRaw)
	sig := ed25519.Sign(pv, txBytes)
	txRaw["signature"] = sig
	txByteWithSig, _ := json.Marshal(txRaw)
	emSig := ed25519.Sign(pv, txByteWithSig)
	txRaw["em_signature"] = emSig
	txBytes, _ = json.Marshal(txRaw)

	tx := mockTransaction{
		addr: []byte("addr"),
		data: map[string]interface{}{
			"encrypted_address_by_node": []byte("addr"),
			"encrypted_wallet":          []byte("wallet"),
			"em_shared_keys_proposal": map[string]interface{}{
				"encrypted_private_key": []byte("pvkey"),
				"public_key":            pub,
			},
		},
		timestamp: time.Now(),
		txType:    0,
		pubKey:    mockPublicKey{bytes: pub},
		coordStmp: coordStmp,
		crossB:    []interface{}{v},
	}
	chainDB.WriteKeychain(tx)
	indexDB.rows[hex.EncodeToString([]byte("addr"))] = []byte("addr")

	pool := mockElectedNodeList{
		nodes: []interface{}{
			mockElectedNode{
				publicKey: mockPublicKey{bytes: pub},
			},
		},
	}

	txRes, err := pr.RequestLastTransaction(pool, []byte("addr"), 1)
	assert.Nil(t, err)
	assert.Equal(t, 0, txRes.Type())
	assert.Equal(t, txBytes, txRes.CoordinatorStamp().(coordinatorStamp).TransactionHash())
}

type mockElectedNodeList struct {
	nodes []interface{}
	pubK  interface{}
	sig   []byte
}

func (l mockElectedNodeList) Nodes() []interface{} {
	return l.nodes
}
func (l mockElectedNodeList) CreatorPublicKey() interface{} {
	return l.pubK
}
func (l mockElectedNodeList) CreatorSignature() []byte {
	return l.sig
}

type mockElectedNode struct {
	publicKey     interface{}
	isUnreachable bool
	isCoord       bool
	patchNb       int
	isOk          bool
}

func (n mockElectedNode) MarshalJSON() ([]byte, error) {
	return json.Marshal(map[string]interface{}{
		"publicKey":     n.publicKey.(publicKey).Marshal(),
		"isUnreachable": n.isUnreachable,
		"isCoordinator": n.isCoord,
		"patchNumber":   n.patchNb,
		"isOk":          n.isOk,
	})
}
func (n mockElectedNode) IsUnreachable() bool {
	return n.isUnreachable
}
func (n mockElectedNode) IsCoordinator() bool {
	return n.isCoord
}
func (n mockElectedNode) IsOK() bool {
	return n.isOk
}
func (n mockElectedNode) PatchNumber() int {
	return n.patchNb
}
func (n mockElectedNode) PublicKey() interface{} {
	return n.publicKey
}
