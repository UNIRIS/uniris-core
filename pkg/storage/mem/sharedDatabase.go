package memstorage

import (
	"github.com/uniris/uniris-core/pkg/shared"
)

//SharedDatabase is a shared database in memory
type SharedDatabase struct {
	nodeCrossKeys      []shared.NodeCrossKeyPair
	emitterCrossKeys   []shared.EmitterCrossKeyPair
	authNodePublicKeys []string

	shared.KeyReader
}

//NewSharedDatabase creates a shared database in memory
func NewSharedDatabase() SharedDatabase {

	crossEmKeys, _ := shared.NewEmitterCrossKeyPair("04a9094a400f48aaea9f04719c295c0d1d1e4d26517310196e816939e6c924b62550f972c5e400725931797d22caae4bb6501b087bf3c898f8355639dc04265fb29c668a72bee6344265f9d713ca7636000a9634021671dce67ae696f7083259802b82e13b237a71d53a1083a932216417f7e84428ba937bf2669e0439bfdced6afe7621ecf00b41dc4bca18e7ff19c1cf7466822da2c2c0386706442b546570bf9b990c5d0d480b42802102f9797e3fc9ed3cc85955f51bebd123ad999dc87cefc27a090c5ec6034ecac0db726eca657dd9cd873020038151d7e3e44c71dd0db19caabe3620a2d91ee5127fa2ef16527074d0f6ce412ec42625e82edb756fe5940acf1f53627ae7934020991446816919b19e7c4de1bf3ba8686b3a1cf4c31616ccbcf3ebeaf5f0585a552b53395e295e9192df41ef50a12a0c98723fa15f2cd9ede372c1de358a46d08d7ab9047ce8e9030c790ef9e9df9afb18990d795e755c465f", "3059301306072a8648ce3d020106082a8648ce3d030107034200041a969b3cdd08cd234d8d6a7f1952f8d38abfd22c39abba1ee026379078ea94f46fb7bfa033a697732969f42f4c9f2495c43cec0057933e1555fff5c8239fa229")
	crossNodeKeys, _ := shared.NewNodeCrossKeyPair("30770201010420bf46e6915518dbca07d79b908499ab2bd2490470bf41b80e73eea0cd6de9f90ea00a06082a8648ce3d030107a1440342000476ab10e633bc8aa3d9225272237428a02b6011a3c5e9a81aff9cfca58ec491ba1d6e0b659fff27db4d11bcc72cb5d862ead6ea05e3cf99b1c70147963e25d9ab", "3059301306072a8648ce3d020106082a8648ce3d0301070342000476ab10e633bc8aa3d9225272237428a02b6011a3c5e9a81aff9cfca58ec491ba1d6e0b659fff27db4d11bcc72cb5d862ead6ea05e3cf99b1c70147963e25d9ab")

	//TODO: once the feature is implemented remove it
	return SharedDatabase{
		nodeCrossKeys:    []shared.NodeCrossKeyPair{crossNodeKeys},
		emitterCrossKeys: []shared.EmitterCrossKeyPair{crossEmKeys},
		authNodePublicKeys: []string{
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

//FirstCrossKeysNode retrieve the first shared crosskeys for the nodes
func (db SharedDatabase) FirstCrossKeysNode() (shared.NodeCrossKeyPair, error) {
	return db.nodeCrossKeys[0], nil
}

//LastNodeCrossKeypair retrieve the last shared crosskeys for the nodes
func (db SharedDatabase) LastNodeCrossKeypair() (shared.NodeCrossKeyPair, error) {
	return db.nodeCrossKeys[len(db.nodeCrossKeys)-1], nil
}

//AuthorizedNodesPublicKeys retrieves the list of public keys of the authorized nodes
func (db SharedDatabase) AuthorizedNodesPublicKeys() ([]string, error) {
	return db.authNodePublicKeys, nil
}
