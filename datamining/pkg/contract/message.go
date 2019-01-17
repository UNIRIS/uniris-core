package contract

import "github.com/uniris/uniris-core/datamining/pkg/mining"

type Message interface {
	ContractAddress() string
	Method() string
	Parameters() []string
	PublicKey() string
	Signature() string
	EmitterSignature() string
}

type EndorsedMessage interface {
	Message
	Endorsement() mining.Endorsement
}

type msg struct {
	address string
	method  string
	params  []string
	pubk    string
	sig     string
	emSig   string
}

func NewMessage(addr, method string, params []string, pubK, sig, emSig string) Message {
	return msg{addr, method, params, pubK, sig, emSig}
}

func (m msg) ContractAddress() string {
	return m.address
}

func (m msg) Method() string {
	return m.method
}

func (m msg) Parameters() []string {
	return m.params
}

func (m msg) PublicKey() string {
	return m.pubk
}

func (m msg) Signature() string {
	return m.sig
}

func (m msg) EmitterSignature() string {
	return m.emSig
}

type endorsedMsg struct {
	msg Message
	end mining.Endorsement
}

func NewEndorsedContractMessage(msg Message, end mining.Endorsement) EndorsedMessage {
	return endorsedMsg{msg, end}
}

func (m endorsedMsg) ContractAddress() string {
	return m.msg.ContractAddress()
}

func (m endorsedMsg) Method() string {
	return m.msg.Method()
}

func (m endorsedMsg) Parameters() []string {
	return m.msg.Parameters()
}

func (m endorsedMsg) PublicKey() string {
	return m.msg.PublicKey()
}

func (m endorsedMsg) Signature() string {
	return m.msg.Signature()
}

func (m endorsedMsg) EmitterSignature() string {
	return m.msg.EmitterSignature()
}

func (m endorsedMsg) Endorsement() mining.Endorsement {
	return m.end
}
