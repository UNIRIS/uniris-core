package listing

//AccountResult defines the account's data returned from the robot
type AccountResult interface {

	//EncryptedAESKey returns the account's encrypted AES key
	EncryptedAESKey() string

	//EncryptedWallet returns account's wallet
	EncryptedWallet() string

	//Encrypted returns the account's address
	EncryptedAddress() string

	//Signature returns the signature of the response
	Signature() string
}

type accRes struct {
	encAesKey string
	encWallet string
	encAddr   string
	sig       string
}

//NewAccountResult creates a new account result
func NewAccountResult(aesKey, wallet, address, sig string) AccountResult {
	return accRes{
		encAesKey: aesKey,
		encWallet: wallet,
		encAddr:   address,
		sig:       sig,
	}
}

func (r accRes) EncryptedAESKey() string {
	return r.encAesKey
}

func (r accRes) EncryptedWallet() string {
	return r.encWallet
}

func (r accRes) EncryptedAddress() string {
	return r.encAddr
}

func (r accRes) Signature() string {
	return r.sig
}
