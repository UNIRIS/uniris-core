package mining

//PowSigVerifier define methods to handle signature verification of the transaction data
type PowSigVerifier interface {

	//VerifyTransactionDataSignature checks the transaction of data signature to perform the Pow
	VerifyTransactionDataSignature(txType TransactionType, pubKey string, data interface{}, sig string) error
}

//ValidationVerifier define methods to handle validation signature verification
type ValidationVerifier interface {

	//VerifyValidationSignature checks the signature of the validation
	VerifyValidationSignature(Validation) error
}

//ValidationSigner define methods to handle signing of the validation
type ValidationSigner interface {

	//SignValidation create signature for a validation data
	SignValidation(v Validation, pvKey string) (Validation, error)
}
