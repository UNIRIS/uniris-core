package system

import (
	"errors"
	"log"
	"net"
	"time"

	"github.com/beevik/ntp"
	"github.com/uniris/uniris-core/autodiscovery/pkg/monitoring"
)

type privatePeerNetworker struct{}

//IP lookups the peer's IP
func (n privatePeerNetworker) IP() (net.IP, error) {
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

//CheckInternetConfig check internet configuration on the node
func (n privatePeerNetworker) CheckInternetState() error {
	return nil
}

//CheckNtp check time synchonization on the node
func (n privatePeerNetworker) CheckNtpState() error {
	for _, ntps := range cntp {
		r, err := ntp.QueryWithOptions(ntps, ntp.QueryOptions{Version: 4})
		if err == nil {
			if (int64(r.ClockOffset/time.Second) < downmaxOffset) || (int64(r.ClockOffset/time.Second) > upmaxOffset) {
				for i := 0; i < ntpRetry; i++ {
					r, err := ntp.QueryWithOptions(ntps, ntp.QueryOptions{Version: 4})
					if err == nil {
						if (int64(r.ClockOffset/time.Second) > downmaxOffset) || (int64(r.ClockOffset/time.Second) < upmaxOffset) {
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
func NewPrivateNetworker() monitoring.PeerNetworker {
	return privatePeerNetworker{}
}
