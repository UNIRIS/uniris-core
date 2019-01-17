package contract

import (
	"github.com/uniris/uniris-core/datamining/pkg/mining"
)

type Contract interface {
	Code() string
	Event() string
	PublicKey() string
	Signature() string
	EmitterSignature() string
}

func New(code, event, pubK, sig, emSig string) Contract {
	return contract{code, event, pubK, sig, emSig}
}

type contract struct {
	code        string
	event       string
	publicKey   string
	signature   string
	emSignature string
}

func (c contract) Code() string {
	return c.code
}

func (c contract) Event() string {
	return c.event
}

func (c contract) PublicKey() string {
	return c.publicKey
}

func (c contract) Signature() string {
	return c.signature
}

func (c contract) EmitterSignature() string {
	return c.emSignature
}

type EndorsedContract interface {
	Address() string
	Contract
	Endorsement() mining.Endorsement
}

type endorsedContract struct {
	address     string
	c           Contract
	endorsement mining.Endorsement
}

func NewEndorsedContract(address string, c Contract, end mining.Endorsement) EndorsedContract {
	return endorsedContract{address, c, end}
}

func (ec endorsedContract) Address() string {
	return ec.address
}

func (ec endorsedContract) Code() string {
	return ec.c.Code()
}

func (ec endorsedContract) Event() string {
	return ec.c.Event()
}

func (ec endorsedContract) PublicKey() string {
	return ec.c.PublicKey()
}

func (ec endorsedContract) Signature() string {
	return ec.c.Signature()
}

func (ec endorsedContract) EmitterSignature() string {
	return ec.c.EmitterSignature()
}

func (ec endorsedContract) Endorsement() mining.Endorsement {
	return ec.endorsement
}
