package adapters

//IdentityChecker checks if a peer is authorized to performs autodiscovery
type IdentityChecker struct {
}

//IsPublicKeyAuthorized checks if a public key is authorized
func (c IdentityChecker) IsPublicKeyAuthorized(publicKey []byte) bool {
	return true //TODO: implements true logic
}
