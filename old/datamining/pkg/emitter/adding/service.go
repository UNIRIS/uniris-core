package adding

import "github.com/uniris/uniris-core/datamining/pkg/emitter"

//Repository define methods to handle emitter keys storage
type Repository interface {

	//StoreSharedEmitterKeyPair stores a shared emitter keypair
	StoreSharedEmitterKeyPair(kp emitter.SharedKeyPair) error
}
