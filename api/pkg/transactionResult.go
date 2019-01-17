package api

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
