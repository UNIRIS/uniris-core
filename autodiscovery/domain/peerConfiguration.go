package domain

//PeerConfiguration wraps the starting configuration
type PeerConfiguration struct {
	Version   string
	PublicKey []byte
	Port      int
	P2PFactor int
}

//NewPeerConfiguration creates a configuration for the peer starting up
func NewPeerConfiguration(ver string, pbKey []byte, port int, p2pFactor int) PeerConfiguration {
	return PeerConfiguration{
		Version:   ver,
		PublicKey: pbKey,
		Port:      port,
		P2PFactor: p2pFactor,
	}
}
