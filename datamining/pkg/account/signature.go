package account

//Signatures represents needed signatures by the account data
type Signatures interface {

	//Biod returns the biometric device signature
	Biod() string

	//Person returns the person signature
	Person() string
}

type sig struct {
	biod   string
	person string
}

//NewSignatures creates a new signatures
func NewSignatures(biod, person string) Signatures {
	return sig{biod, person}
}

func (s sig) Biod() string {
	return s.biod
}

func (s sig) Person() string {
	return s.person
}
