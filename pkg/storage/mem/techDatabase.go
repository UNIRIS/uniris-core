package memstorage

import (
	"encoding/hex"

	"github.com/uniris/uniris-core/pkg/crypto"
	"github.com/uniris/uniris-core/pkg/shared"
)

type techDB struct {
	emitterKeys shared.EmitterKeys
	nodeKeys    []shared.NodeKeyPair
}

//NewTechDatabase creates a tech database in memory
func NewTechDatabase() shared.TechDatabaseReader {
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

	//TODO: once the feature is implemented
	emPubBytes, _ := hex.DecodeString("0044657dab453d34f9adc2100a2cb8f38f644ef48e34b1d99d7c4d9371068e9438")
	emPub, _ := crypto.ParsePublicKey(emPubBytes)

	emKPBytes, _ := hex.DecodeString("fcec12fb1f715926e4489a10411b4a48451b616c35e7ad99dedb8864f71943ce9235d2e5ecb554604167faf90cb9f532313c0c5aea34573eb821eae985c49dc151e45f1e437c8976549c49302458c51f85ac6f83db510d76674fbefa2e19767620d82baa7b58acd433f6cc6755323346d2925e6ecb5e2d2172ec892e8254568a390c28f5701f2d57eace3340629177f3142b00ebc6c66efe87b6880538")
	emKP, _ := shared.NewEmitterKeyPair(emKPBytes, emPub)

	nodePvBytes, _ := hex.DecodeString("0066fc58626e8e245f58c9c705609f697071bde3968d33ab1022ea4488832f3f5aee7a047a226e08ea14fe60ec4f6d328e56ebdb2ee2b9f5b1120e231e05c956a3")
	nodePv, _ := crypto.ParsePrivateKey(nodePvBytes)

	nodePubBytes, _ := hex.DecodeString("00ee7a047a226e08ea14fe60ec4f6d328e56ebdb2ee2b9f5b1120e231e05c956a3")
	nodePub, _ := crypto.ParsePublicKey(nodePubBytes)

	nodeKP, _ := shared.NewNodeKeyPair(nodePub, nodePv)

	return &techDB{
		emitterKeys: shared.EmitterKeys{emKP},
		nodeKeys:    []shared.NodeKeyPair{nodeKP},
	}
}

func (db techDB) EmitterKeys() (shared.EmitterKeys, error) {
	return db.emitterKeys, nil
}

func (db *techDB) NodeLastKeys() (shared.NodeKeyPair, error) {
	return db.nodeKeys[len(db.nodeKeys)-1], nil
}
