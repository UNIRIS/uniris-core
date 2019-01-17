package adding

type ContractCreationRequest interface {
	EncryptedContract() string
	Signature() string
}

func NewContractCreationRequest(encContract string, sig string) ContractCreationRequest {
	return contractCreationReq{encContract, sig}
}

type contractCreationReq struct {
	encContract string
	sig         string
}

func (c contractCreationReq) EncryptedContract() string {
	return c.encContract
}

func (c contractCreationReq) Signature() string {
	return c.sig
}
