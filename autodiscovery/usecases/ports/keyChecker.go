package ports

//PeerIdentityChecker checks if a peer is authorized to performs autodiscovery
type PeerIdentityChecker interface {
	IsPublicKeyAuthorized(publicKey []byte) bool
}
