package adding

type ContractCreationRequest interface {
	Address() string
	Code() string
	Event() string
	PublicKey() string
	Signature() string
	EmitterSignature() string
	RequestSignature() string
}

func NewContractCreationRequest(addr, code, event, pubK, sig, emSig, reqSig string) ContractCreationRequest {
	return contractCreationReq{addr, code, event, pubK, sig, emSig, reqSig}
}

type contractCreationReq struct {
	address      string
	code         string
	event        string
	publicKey    string
	signature    string
	emSignature  string
	reqSignature string
}

func (c contractCreationReq) Address() string {
	return c.address
}

func (c contractCreationReq) Code() string {
	return c.code
}

func (c contractCreationReq) Event() string {
	return c.event
}

func (c contractCreationReq) PublicKey() string {
	return c.publicKey
}

func (c contractCreationReq) Signature() string {
	return c.signature
}

func (c contractCreationReq) EmitterSignature() string {
	return c.emSignature
}

func (c contractCreationReq) RequestSignature() string {
	return c.reqSignature
}
