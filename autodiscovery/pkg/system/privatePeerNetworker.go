package system

import (
	"errors"
	"fmt"
	"net"
	"time"

	"github.com/beevik/ntp"
	"github.com/uniris/uniris-core/autodiscovery/pkg/monitoring"
)

type privatePeerNetworker struct {
	iface string
}

//IP lookups the peer's IP
func (n privatePeerNetworker) IP() (net.IP, error) {
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

//CheckInternetConfig check internet configuration on the node
func (n privatePeerNetworker) CheckInternetState() error {
	_, err := net.LookupIP(cdns)
	if err != nil {
		return err
	}

	return nil
}

//CheckNtp check time synchonization on the node
func (n privatePeerNetworker) CheckNtpState() error {
	for _, ntps := range cntp {
		r, err := ntp.QueryWithOptions(ntps, ntp.QueryOptions{Version: 4})
		if err == nil {
			if (int64(r.ClockOffset/time.Second) < downmaxoffset) || (int64(r.ClockOffset/time.Second) > upmaxoffset) {
				for i := 0; i < ntpretry; i++ {
					r, err := ntp.QueryWithOptions(ntps, ntp.QueryOptions{Version: 4})
					if err == nil {
						if (int64(r.ClockOffset/time.Second) > downmaxoffset) || (int64(r.ClockOffset/time.Second) < upmaxoffset) {
							return nil
						}
					}
				}
				return errors.New("System Clock have a big Offset check the ntp configuration of the system")
			}
			return nil
		}
	}
	return errors.New("Could not get reply from ntp servers")
}

//NewPrivateNetworker creates a new instance of the local implementation of the PeerNetworker interface
func NewPrivateNetworker(iface string) monitoring.PeerNetworker {
	return privatePeerNetworker{iface}
}
