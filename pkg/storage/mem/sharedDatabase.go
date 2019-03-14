package memstorage

import (
	"encoding/hex"

	"github.com/uniris/uniris-core/pkg/crypto"
	"github.com/uniris/uniris-core/pkg/shared"
)

//SharedDatabase is a shared database in memory
type SharedDatabase struct {
	nodeCrossKeys      []shared.NodeCrossKeyPair
	emitterCrossKeys   []shared.EmitterCrossKeyPair
	authNodePublicKeys []crypto.PublicKey

	shared.KeyReadWriter
}

//NewSharedDatabase creates a shared database in memory
func NewSharedDatabase() *SharedDatabase {

	// "ed25519"
	// "em pv key": "000c3bb61141f052e1936823a4a56224f2aae04084265655ff4c83d885295b570344657dab453d34f9adc2100a2cb8f38f644ef48e34b1d99d7c4d9371068e9438"
	// "em pubkey": "0044657dab453d34f9adc2100a2cb8f38f644ef48e34b1d99d7c4d9371068e9438"
	// "enc pv key": "fcec12fb1f715926e4489a10411b4a48451b616c35e7ad99dedb8864f71943ce9235d2e5ecb554604167faf90cb9f532313c0c5aea34573eb821eae985c49dc151e45f1e437c8976549c49302458c51f85ac6f83db510d76674fbefa2e19767620d82baa7b58acd433f6cc6755323346d2925e6ecb5e2d2172ec892e8254568a390c28f5701f2d57eace3340629177f3142b00ebc6c66efe87b6880538"

	// "ed25519"
	// "em root pv key": "0077ff2d9233bdad72fb5c7a178003fb8c7fb3a69625075a8abc3064e62e714b66a8e0f20d4da185d0bf8bd0a45995dfc7926d545e5bbff0194fe34c42bf5e221b"
	// "em root pub": "00a8e0f20d4da185d0bf8bd0a45995dfc7926d545e5bbff0194fe34c42bf5e221b"

	//"node"
	//"pub": "00ee7a047a226e08ea14fe60ec4f6d328e56ebdb2ee2b9f5b1120e231e05c956a3"
	//"pv": "0066fc58626e8e245f58c9c705609f697071bde3968d33ab1022ea4488832f3f5aee7a047a226e08ea14fe60ec4f6d328e56ebdb2ee2b9f5b1120e231e05c956a3"

	crossEmPubBytes, _ := hex.DecodeString("0044657dab453d34f9adc2100a2cb8f38f644ef48e34b1d99d7c4d9371068e9438")
	crossEmPub, _ := crypto.ParsePublicKey(crossEmPubBytes)

	crossEmPvBytes, _ := hex.DecodeString("fcec12fb1f715926e4489a10411b4a48451b616c35e7ad99dedb8864f71943ce9235d2e5ecb554604167faf90cb9f532313c0c5aea34573eb821eae985c49dc151e45f1e437c8976549c49302458c51f85ac6f83db510d76674fbefa2e19767620d82baa7b58acd433f6cc6755323346d2925e6ecb5e2d2172ec892e8254568a390c28f5701f2d57eace3340629177f3142b00ebc6c66efe87b6880538")
	crossEmKeys, _ := shared.NewEmitterCrossKeyPair(crossEmPvBytes, crossEmPub)

	crossNodePvBytes, _ := hex.DecodeString("0066fc58626e8e245f58c9c705609f697071bde3968d33ab1022ea4488832f3f5aee7a047a226e08ea14fe60ec4f6d328e56ebdb2ee2b9f5b1120e231e05c956a3")
	crossNodePv, _ := crypto.ParsePrivateKey(crossNodePvBytes)

	crossNodePubBytes, _ := hex.DecodeString("00ee7a047a226e08ea14fe60ec4f6d328e56ebdb2ee2b9f5b1120e231e05c956a3")
	crossNodePub, _ := crypto.ParsePublicKey(crossNodePubBytes)
	crossNodeKeys, _ := shared.NewNodeCrossKeyPair(crossNodePub, crossNodePv)

	//TODO: once the feature is implemented remove it
	return &SharedDatabase{
		nodeCrossKeys:    []shared.NodeCrossKeyPair{crossNodeKeys},
		emitterCrossKeys: []shared.EmitterCrossKeyPair{crossEmKeys},
		authNodePublicKeys: []crypto.PublicKey{
			crossNodeKeys.PublicKey(),
		},
	}
}

//EmitterCrossKeypairs retrieve the list of the cross emitter keys
func (db SharedDatabase) EmitterCrossKeypairs() ([]shared.EmitterCrossKeyPair, error) {
	return db.emitterCrossKeys, nil
}

//FirstEmitterCrossKeypair retrieves the first public key
func (db SharedDatabase) FirstEmitterCrossKeypair() (shared.EmitterCrossKeyPair, error) {
	return db.emitterCrossKeys[0], nil
}

//FirstNodeCrossKeypair retrieve the first shared crosskeys for the nodes
func (db SharedDatabase) FirstNodeCrossKeypair() (shared.NodeCrossKeyPair, error) {
	return db.nodeCrossKeys[0], nil
}

//LastNodeCrossKeypair retrieve the last shared crosskeys for the nodes
func (db SharedDatabase) LastNodeCrossKeypair() (shared.NodeCrossKeyPair, error) {
	return db.nodeCrossKeys[len(db.nodeCrossKeys)-1], nil
}

//AuthorizedNodesPublicKeys retrieves the list of public keys of the authorized nodes
func (db SharedDatabase) AuthorizedNodesPublicKeys() ([]crypto.PublicKey, error) {
	return db.authNodePublicKeys, nil
}

//WriteAuthorizedNode inserts a new node public key as an authorized node
func (db *SharedDatabase) WriteAuthorizedNode(pub crypto.PublicKey) error {
	db.authNodePublicKeys = append(db.authNodePublicKeys, pub)
	return nil
}
