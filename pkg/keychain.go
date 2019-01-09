package uniris

//Keychain represents a keychain transaction
type Keychain struct {
	addr      string
	encAddr   string
	encWallet string
	prev      *Keychain
	Transaction
}

//NewKeychain creates a transaction related to a keychain
func NewKeychain(addr string, encAddr string, encWallet string, tx Transaction) Keychain {
	return Keychain{
		addr:        addr,
		encAddr:     encAddr,
		encWallet:   encWallet,
		Transaction: tx,
	}
}

//EncryptedAddrByRobot returns the encrypted keychain address by the shared robot key
func (k Keychain) EncryptedAddrByRobot() string {
	return k.encAddr
}

//EncryptedWallet returns encrypted wallet by the person AES key
func (k Keychain) EncryptedWallet() string {
	return k.encWallet
}

//Address returns the keychain address
func (k Keychain) Address() string {
	return k.addr
}

//Chain links previous keychain to this one
func (k *Keychain) Chain(prevK Keychain) {
	k.prev = &prevK
}
