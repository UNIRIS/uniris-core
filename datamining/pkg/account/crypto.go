package account

//KeychainSignatureVerifier define methods to handle keychain signatures verification
type KeychainSignatureVerifier interface {

	//VerifyKeychainSignatures checks the signatures of the keychain
	VerifyKeychainSignatures(Keychain) error
}

//IDSignatureVerifier define methods to handle biometric signatures verification
type IDSignatureVerifier interface {

	//VerifyIDSignatures checks the signatures of the ID
	VerifyIDSignatures(ID) error
}

//KeychainHasher defines methods to handle keychain hashing
type KeychainHasher interface {

	//HashEndorsedKeychain produces hash of the endorsed keychain
	HashEndorsedKeychain(EndorsedKeychain) (string, error)

	//HashKeychain produces hash of the keychain
	HashKeychain(Keychain) (string, error)
}

//IDHasher defines methods to handle ID hashing
type IDHasher interface {

	//HashEndorsedID produces hash of the endorsed ID
	HashEndorsedID(EndorsedID) (string, error)

	//HashID produces hash of the ID
	HashID(ID) (string, error)
}
