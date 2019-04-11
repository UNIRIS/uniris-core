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

	crossNodePvBytes2, _ := hex.DecodeString("0077ff2d9233bdad72fb5c7a178003fb8c7fb3a69625075a8abc3064e62e714b66a8e0f20d4da185d0bf8bd0a45995dfc7926d545e5bbff0194fe34c42bf5e221b")
	crossNodePv2, _ := crypto.ParsePrivateKey(crossNodePvBytes2)

	crossNodePvBytes3, _ := hex.DecodeString("00b274c45f99fe42de79433e8c8b71edcf4e61b6b8e55406883ba2c423785705af29dc9568f67727da6b0d93ca8538fe171a6bcd239dfd08dea748986848f24466")
	crossNodePv3, _ := crypto.ParsePrivateKey(crossNodePvBytes3)

	crossNodePvBytes4, _ := hex.DecodeString("00b84ac63695ef472b9a23dab78e388cb6ba2a7a0578b51c3ba4e4d7e5b23ee28a947b6ca58d87a126fc410858046c3edaea3dd1d570275502e7b45331c47fb655")
	crossNodePv4, _ := crypto.ParsePrivateKey(crossNodePvBytes4)

	crossNodePvBytes5, _ := hex.DecodeString("00292fab83ef20397b07d2f9f0adaab2ff093d269111c70d3f80d1b7cc9ee6c1b7ca8a194ecb4ecc61287124a7f5d4db80a1a3d203ed43ace2e4fabf5e78a6bc83")
	crossNodePv5, _ := crypto.ParsePrivateKey(crossNodePvBytes5)

	crossNodePvBytes6, _ := hex.DecodeString("003dfcbdfb38b042a9c319677a93999c5e3125a64e1441c7409d4c7db5434021295475d994fecb8492dfe05df3025c7fdfacdec323970b1c96c2d75f7ac5ef0fc3")
	crossNodePv6, _ := crypto.ParsePrivateKey(crossNodePvBytes6)

	crossNodePvBytes7, _ := hex.DecodeString("00a5805009c794e7a90203417405c2ace5189af3a9fd68bc4ccfb5f1ccf4a0e00cb43683f3f473af5719f294f93519c579fbcf1f603080df8b0ccfe042aa5896b4")
	crossNodePv7, _ := crypto.ParsePrivateKey(crossNodePvBytes7)

	crossNodePvBytes8, _ := hex.DecodeString("0076b7c35321d6f4e126bd6efaff8e33114f56a309756ef492d206b38e5a70df25254bd5c54d29afc54156042f50f0b3424c4f1a60882a3e2ea71d4e803ae301c5")
	crossNodePv8, _ := crypto.ParsePrivateKey(crossNodePvBytes8)

	crossNodePubBytes, _ := hex.DecodeString("00ee7a047a226e08ea14fe60ec4f6d328e56ebdb2ee2b9f5b1120e231e05c956a3")
	crossNodePub, _ := crypto.ParsePublicKey(crossNodePubBytes)
	crossNodeKeys, _ := shared.NewNodeCrossKeyPair(crossNodePub, crossNodePv)

	crossNodePubBytes2, _ := hex.DecodeString("00a8e0f20d4da185d0bf8bd0a45995dfc7926d545e5bbff0194fe34c42bf5e221b")
	crossNodePub2, _ := crypto.ParsePublicKey(crossNodePubBytes2)
	crossNodeKeys2, _ := shared.NewNodeCrossKeyPair(crossNodePub2, crossNodePv2)

	crossNodePubBytes3, _ := hex.DecodeString("0029dc9568f67727da6b0d93ca8538fe171a6bcd239dfd08dea748986848f24466")
	crossNodePub3, _ := crypto.ParsePublicKey(crossNodePubBytes3)
	crossNodeKeys3, _ := shared.NewNodeCrossKeyPair(crossNodePub3, crossNodePv3)

	crossNodePubBytes4, _ := hex.DecodeString("00947b6ca58d87a126fc410858046c3edaea3dd1d570275502e7b45331c47fb655")
	crossNodePub4, _ := crypto.ParsePublicKey(crossNodePubBytes4)
	crossNodeKeys4, _ := shared.NewNodeCrossKeyPair(crossNodePub4, crossNodePv4)

	crossNodePubBytes5, _ := hex.DecodeString("00ca8a194ecb4ecc61287124a7f5d4db80a1a3d203ed43ace2e4fabf5e78a6bc83")
	crossNodePub5, _ := crypto.ParsePublicKey(crossNodePubBytes5)
	crossNodeKeys5, _ := shared.NewNodeCrossKeyPair(crossNodePub5, crossNodePv5)

	crossNodePubBytes6, _ := hex.DecodeString("005475d994fecb8492dfe05df3025c7fdfacdec323970b1c96c2d75f7ac5ef0fc3")
	crossNodePub6, _ := crypto.ParsePublicKey(crossNodePubBytes6)
	crossNodeKeys6, _ := shared.NewNodeCrossKeyPair(crossNodePub6, crossNodePv6)

	crossNodePubBytes7, _ := hex.DecodeString("00b43683f3f473af5719f294f93519c579fbcf1f603080df8b0ccfe042aa5896b4")
	crossNodePub7, _ := crypto.ParsePublicKey(crossNodePubBytes7)
	crossNodeKeys7, _ := shared.NewNodeCrossKeyPair(crossNodePub7, crossNodePv7)

	crossNodePubBytes8, _ := hex.DecodeString("00254bd5c54d29afc54156042f50f0b3424c4f1a60882a3e2ea71d4e803ae301c5")
	crossNodePub8, _ := crypto.ParsePublicKey(crossNodePubBytes8)
	crossNodeKeys8, _ := shared.NewNodeCrossKeyPair(crossNodePub8, crossNodePv8)

	//TODO: once the feature is implemented remove it
	return &SharedDatabase{
		nodeCrossKeys:    []shared.NodeCrossKeyPair{crossNodeKeys},
		emitterCrossKeys: []shared.EmitterCrossKeyPair{crossEmKeys},
		authNodePublicKeys: []crypto.PublicKey{
			crossNodeKeys.PublicKey(),
			crossNodeKeys2.PublicKey(),
			crossNodeKeys3.PublicKey(),
			crossNodeKeys4.PublicKey(),
			crossNodeKeys5.PublicKey(),
			crossNodeKeys6.PublicKey(),
			crossNodeKeys7.PublicKey(),
			crossNodeKeys8.PublicKey(),
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
	var found bool
	for _, k := range db.authNodePublicKeys {
		if k.Equals(pub) {
			found = true
			break
		}
	}

	if !found {
		db.authNodePublicKeys = append(db.authNodePublicKeys, pub)
	}

	return nil
}

//IsAuthorizedNode check if the public Key is on the authorized list
func (db *SharedDatabase) IsAuthorizedNode(pub crypto.PublicKey) bool {
	found := false
	for _, k := range db.authNodePublicKeys {
		if k.Equals(pub) {
			found = true
			break
		}
	}
	return found
}
