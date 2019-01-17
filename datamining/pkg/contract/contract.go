package contract

import (
	"github.com/uniris/uniris-core/datamining/pkg/mining"
)

type Contract interface {
	Address() string
	Code() string
	Event() string
	PublicKey() string
	Signature() string
	EmitterSignature() string
}

func New(addr, code, event, pubK, sig, emSig string) Contract {
	return contract{addr, code, event, pubK, sig, emSig}
}

type contract struct {
	address     string
	code        string
	event       string
	publicKey   string
	signature   string
	emSignature string
}

func (c contract) Address() string {
	return c.address
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
	Contract
	Endorsement() mining.Endorsement
}

type endorsedContract struct {
	c           Contract
	endorsement mining.Endorsement
}

func NewEndorsedContract(c Contract, end mining.Endorsement) EndorsedContract {
	return endorsedContract{c, end}
}

func (ec endorsedContract) Address() string {
	return ec.c.Address()
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
