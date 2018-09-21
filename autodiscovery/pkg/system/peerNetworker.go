package system

import (
	"fmt"
	"net"

	"github.com/uniris/uniris-core/autodiscovery/pkg/bootstraping"
)

//SystemNetwork implements the PeerNetworker interface which provides the methods to get network peer's details
type systemNetworker struct {
	iface string
}

//IP lookups the peer's IP
func (n systemNetworker) IP() (net.IP, error) {
	ifaces, err := net.Interfaces()
	if err != nil {
		return nil, err
	}
	var iface net.Interface
	for _, i := range ifaces {
		if i.Name == n.iface {
			iface = i
			break
		}
	}
	addrs, err := iface.Addrs()
	if err != nil {
		return nil, err
	}
	var ip net.IP
	for _, addr := range addrs {
		switch v := addr.(type) {
		case *net.IPNet:
			ip = v.IP
		case *net.IPAddr:
			ip = v.IP
		}
		return ip, nil
	}
	return nil, fmt.Errorf("Cannot find a IP address from the interface %s", n.iface)
}

//NewPeerNetworker creates a new instance of the system implementation of the PeerNetworker interface
func NewPeerNetworker(iface string) bootstraping.PeerNetworker {
	return systemNetworker{iface}
}
