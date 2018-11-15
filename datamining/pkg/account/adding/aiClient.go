package adding

//AIClient handles AI requests
type AIClient interface {

	//CheckStorageAuthorization asks the AI service to determines
	//if the storage of this transaction is authorized to be done on this peer
	CheckStorageAuthorization(txHash string) error

	//GetMininumValidations returns the minimum required validations for this transaction hash
	GetMininumValidations(txHash string) (int, error)
}
