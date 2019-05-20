package rpc

import (
	"net"
	"time"
)

type transaction interface {
	Address() []byte
	Type() int
	Data() map[string]interface{}
	Timestamp() time.Time
	PreviousPublicKey() interface{}
	Signature() []byte
	OriginSignature() []byte
	CoordinatorStamp() interface{}
	CrossValidations() []interface{}
	MarshalBeforeOriginSignature() ([]byte, error)
	MarshalRoot() ([]byte, error)
}

type publicKey interface {
	Verify(data []byte, sig []byte) (bool, error)
	Marshal() []byte
}

type privateKey interface {
	Sign(data []byte) ([]byte, error)
}

type coordinatorStamp interface {
	PreviousCrossValidators() [][]byte
	ProofOfWork() publicKey
	ValidationStamp() validationStamp
	TransactionHash() []byte
	ElectedCoordinatorNodes() electedNodeList
	ElectedCrossValidationNodes() electedNodeList
	ElectedStorageNodes() electedNodeList
}

type electedNodeList interface {
	Nodes() []electedNode
	CreatorPublicKey() publicKey
	CreatorSignature() []byte
}

type validationStamp interface {
	Status() int
	Timestamp() time.Time
	NodePublicKey() publicKey
	NodeSignature() []byte
}

type electedNode interface {
	MarshalJSON() ([]byte, error)
	IsUnreachable() bool
	IsCoordinator() bool
	IsOK() bool
	PatchNumber() int
	PublicKey() publicKey
}

type sharedKeyReader interface {

	//EmitterCrossKeypairs retrieve the list of the cross emitter keys
	EmitterCrossKeypairs() ([]emitterCrossKeyPair, error)

	//FirstEmitterCrossKeypair retrieves the first public key
	FirstEmitterCrossKeypair() (emitterCrossKeyPair, error)

	//CrossEmitterPublicKeys retrieves the public keys of the cross emitter keys
	CrossEmitterPublicKeys() ([]publicKey, error)

	//FirstNodeCrossKeypair retrieve the first shared crosskeys for the nodes
	FirstNodeCrossKeypair() (publicKey, privateKey, error)

	//LastNodeCrossKeypair retrieve the last shared crosskeys for the nodes
	LastNodeCrossKeypair() (publicKey, privateKey, error)

	//AuthorizedNodesPublicKeys retrieves the list of public keys of the authorized nodes
	AuthorizedNodesPublicKeys() ([]publicKey, error)
}

type emitterCrossKeyPair interface {
	PublicKey() publicKey
	EncryptedPrivateKey() []byte
}

type chainDB interface {
	WriteKOTransaction(t transaction) error
	WriteTransaction(t transaction) error
	GetTransactionByHash(txHash []byte) (transaction, error)
	LastTransaction(genesis []byte, txType int) (transaction, error)
	GetTransactionStatus(txHash []byte) (int, error)
}

type nodeReader interface {
	CountReachables() (int, error)
	Reachables() ([]node, error)
	Unreachables() ([]node, error)
	FindByPublicKey(publicKey publicKey) (node, error)
}

type node interface {
	IP() net.IP
	Port() int
	PatchNumber() int
	PublicKey() publicKey
}
