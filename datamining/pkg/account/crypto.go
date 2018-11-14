package account

//KeychainSignatureVerifier define methods to handle keychain signatures verification
type KeychainSignatureVerifier interface {
	//VerifyKeychainDataSignatures checks the signatures of the keychain data
	VerifyKeychainDataSignatures(KeychainData) error
}

//BiometricSignatureVerifier define methods to handle biometric signatures verification
type BiometricSignatureVerifier interface {
	//VerifyBiometricDataSignature checks the signatures
	VerifyBiometricDataSignatures(BiometricData) error
}

//KeychainHasher defines methods to handle keychain hashing
type KeychainHasher interface {

	//HashKeychain produces hash of the keychain (including data, address and endorsement)
	HashKeychain(Keychain) (string, error)

	//HashKeychainData produces hash of the keychain data
	HashKeychainData(KeychainData) (string, error)
}

//BiometricHasher defines methods to handle keychain hashing
type BiometricHasher interface {

	//HashBiometric produces hash of the biometric (including data and endorsement)
	HashBiometric(Biometric) (string, error)

	//HashBiometricData produces hash of the biometric data
	HashBiometricData(BiometricData) (string, error)
}
