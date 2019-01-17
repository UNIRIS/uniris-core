package adding

type ContractCreationRequest interface {
	Code() string
	Event() string
	PublicKey() string
	Signature() string
	EmitterSignature() string
	RequestSignature() string
}

func NewContractCreationRequest(code, event, pubK, sig, emSig, reqSig string) ContractCreationRequest {
	return contractCreationReq{code, event, pubK, sig, emSig, reqSig}
}

type contractCreationReq struct {
	code         string
	event        string
	publicKey    string
	signature    string
	emSignature  string
	reqSignature string
}

func NewContractCreationResponse(txHash, addr, masterPeer, sig string) ContractCreationResponse {
	return contractCreateRes{txHash, addr, masterPeer, sig}
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

type ContractCreationResponse interface {
	Address() string
	TransactionResult
}

type contractCreateRes struct {
	txHash       string
	address      string
	masterPeerIP string
	signature    string
}

func (c contractCreateRes) Address() string {
	return c.address
}

func (c contractCreateRes) TransactionHash() string {
	return c.txHash
}

func (c contractCreateRes) MasterPeerIP() string {
	return c.masterPeerIP
}

func (c contractCreateRes) Signature() string {
	return c.signature
}
