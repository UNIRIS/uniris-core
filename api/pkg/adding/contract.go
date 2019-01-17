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

type ContractMessageCreationRequest interface {
	ContractAddress() string
	EncryptedMessage() string
	Signature() string
}

type contractMsg struct {
	contractAddress string
	msg             string
	sig             string
}

func NewContractMessageCreationRequest(contractAddress, msg, sig string) ContractMessageCreationRequest {
	return contractMsg{contractAddress, msg, sig}
}

func (m contractMsg) ContractAddress() string {
	return m.contractAddress
}

func (m contractMsg) EncryptedMessage() string {
	return m.msg
}

func (m contractMsg) Signature() string {
	return m.sig
}
