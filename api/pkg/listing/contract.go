package listing

type ContractState interface {
	Data() string
	Signature() string
}

type contractState struct {
	data string
	sig  string
}

func NewContractState(data, sig string) ContractState {
	return contractState{data, sig}
}

func (c contractState) Data() string {
	return c.data
}

func (c contractState) Signature() string {
	return c.sig
}
