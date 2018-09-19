package system

import (
	"log"
	"net"

	"github.com/uniris/uniris-core/autodiscovery/pkg/bootstraping"
)

//SystemNetwork implements the PeerNetworker interface which provides the methods to get network peer's details
type systemNetworker struct{}

//IP lookups the peer's IP
func (n systemNetworker) IP() (net.IP, error) {
	ifaces, err := net.Interfaces()
	if err != nil {
		return nil, err
	}
	for _, i := range ifaces {
		addrs, err := i.Addrs()
		if err != nil {
			return nil, err
		}
		for _, addr := range addrs {
			log.Printf(addr.String())
			var ip net.IP
			switch v := addr.(type) {
			case *net.IPNet:
				ip = v.IP
			case *net.IPAddr:
				ip = v.IP
			}
			// process IP address
			return ip, nil
		}
	}
	return nil, nil
}

//NewPeerNetworker creates a new instance of the system implementation of the PeerNetworker interface
func NewPeerNetworker() bootstraping.PeerNetworker {
	return systemNetworker{}
}
