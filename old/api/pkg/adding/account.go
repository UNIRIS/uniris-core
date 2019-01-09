package adding

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
	ID() TransactionResult

	//Keychain returns the Keychain transaction result
	Keychain() TransactionResult
}

type accTxRes struct {
	id       TransactionResult
	keychain TransactionResult
}

//NewAccountCreationTransactionResult create a new creation transaction result
func NewAccountCreationTransactionResult(id TransactionResult, keychain TransactionResult) AccountCreationTransactionResult {
	return accTxRes{id, keychain}
}

func (r accTxRes) ID() TransactionResult {
	return r.id
}

func (r accTxRes) Keychain() TransactionResult {
	return r.keychain
}

//TransactionResult represents the result for a transaction
type TransactionResult interface {

	//Transaction returns the transaction hash of the account data creation
	TransactionHash() string

	//MasterPeerIP returns the IP of the peer leading the transaction
	MasterPeerIP() string

	//Signature returns the signature of the transaction processing
	Signature() string
}

type txRes struct {
	txHash   string
	masterIP string
	sig      string
}

//NewTransactionResult creates a new transaction result
func NewTransactionResult(txHash string, masterIP string, sig string) TransactionResult {
	return txRes{txHash, masterIP, sig}
}

func (r txRes) TransactionHash() string {
	return r.txHash
}

func (r txRes) MasterPeerIP() string {
	return r.masterIP
}

func (r txRes) Signature() string {
	return r.sig
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
