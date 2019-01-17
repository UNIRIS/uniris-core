package adding

import api "github.com/uniris/uniris-core/api/pkg"

//AccountCreationResult represents the result of the account creation
type AccountCreationResult interface {

	//ResultTransactions returns the result transactions of the account creation
	ResultTransactions() AccountCreationTransactionResult

	//Signature returns the signature of the result
	Signature() string
}

type accCreateRes struct {
	txs AccountCreationTransactionResult
	sig string
}

//NewAccountCreationResult creates a new account creation result
func NewAccountCreationResult(txs AccountCreationTransactionResult, sig string) AccountCreationResult {
	return accCreateRes{txs, sig}
}

func (r accCreateRes) ResultTransactions() AccountCreationTransactionResult {
	return r.txs
}

func (r accCreateRes) Signature() string {
	return r.sig
}

//AccountCreationTransactionResult represents the transactions for the account creation
type AccountCreationTransactionResult interface {

	//ID returns the ID transaction result
	ID() api.TransactionResult

	//Keychain returns the Keychain transaction result
	Keychain() api.TransactionResult
}

type accTxRes struct {
	id       api.TransactionResult
	keychain api.TransactionResult
}

//NewAccountCreationTransactionResult create a new creation transaction result
func NewAccountCreationTransactionResult(id api.TransactionResult, keychain api.TransactionResult) AccountCreationTransactionResult {
	return accTxRes{id, keychain}
}

func (r accTxRes) ID() api.TransactionResult {
	return r.id
}

func (r accTxRes) Keychain() api.TransactionResult {
	return r.keychain
}

//AccountCreationRequest represents the required data to create an account
type AccountCreationRequest interface {

	//EncryptedID returns the encrypted ID data to create
	EncryptedID() string

	//EncryptedKeychain returns the encrypted Keychain data to create
	EncryptedKeychain() string

	//Signature returns the signature of the
	Signature() string
}

type accCreateReq struct {
	encID       string
	encKeychain string
	sig         string
}

//NewAccountCreationRequest creates a new account creation request
func NewAccountCreationRequest(encID, encKeychain, sig string) AccountCreationRequest {
	return accCreateReq{encID, encKeychain, sig}
}

func (r accCreateReq) EncryptedID() string {
	return r.encID
}

func (r accCreateReq) EncryptedKeychain() string {
	return r.encKeychain
}

func (r accCreateReq) Signature() string {
	return r.sig
}
