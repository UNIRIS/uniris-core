package adding

//Repository define methods to handle emitter keys storage
type Repository interface {

	//StoreEmitterSharedKey stores a shared public key into the database
	StoreEmitterSharedKey(pubKey string) error
}
