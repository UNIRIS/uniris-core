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
	EncryptedMessage() string
	Signature() string
}

type contractMsg struct {
	msg string
	sig string
}

func NewContractMessageCreationRequest(msg, sig string) ContractMessageCreationRequest {
	return contractMsg{msg, sig}
}

func (m contractMsg) EncryptedMessage() string {
	return m.msg
}

func (m contractMsg) Signature() string {
	return m.sig
}
