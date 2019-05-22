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
	ProofOfWork() interface{}
	ValidationStamp() interface{}
	TransactionHash() []byte
	ElectedCoordinatorNodes() interface{}
	ElectedCrossValidationNodes() interface{}
	ElectedStorageNodes() interface{}
}

type electedNodeList interface {
	Nodes() []interface{}
	CreatorPublicKey() interface{}
	CreatorSignature() []byte
}

type validationStamp interface {
	Status() int
	Timestamp() time.Time
	NodePublicKey() interface{}
	NodeSignature() []byte
}

type electedNode interface {
	MarshalJSON() ([]byte, error)
	IsUnreachable() bool
	IsCoordinator() bool
	IsOK() bool
	PatchNumber() int
	PublicKey() interface{}
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
	IsReachable() bool
}

type chainDB interface {
	WriteKeychain(tx interface{}) error
	WriteID(tx interface{}) error
	WriteKO(tx interface{}) error

	FindKeychainByAddr(addr []byte) (interface{}, error)
	FindKeychainByHash(txHash []byte) (interface{}, error)
	FindIDByHash(txHash []byte) (interface{}, error)
	FindIDByAddr(addr []byte) (interface{}, error)
	FindKOByHash(txHash []byte) (interface{}, error)
	FindKOByAddr(addr []byte) (interface{}, error)
}

type indexDB interface {
	FindLastTransactionAddr(genesis []byte) ([]byte, error)
}
