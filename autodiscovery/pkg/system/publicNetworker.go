package system

import (
	"errors"
	"io/ioutil"
	"net"
	"net/http"
	"time"

	"github.com/beevik/ntp"
	"github.com/uniris/uniris-core/autodiscovery/pkg/monitoring"
)

//ErrFailToGetIP is returned when the service to get IP does not respond
var ErrFailToGetIP = errors.New("Cannot get the peer IP. IP providers may failed")

type publicPeerNetworker struct{}

//IP lookups the peer's IP
func (n publicPeerNetworker) IP() (net.IP, error) {
	var ip net.IP
	ip, err := n.ipify()
	if err != nil {
		ip, err := n.myExternalIP()
		if err != nil {
			return nil, ErrFailToGetIP
		}
		return ip, nil
	}
	return ip, nil
}

func (n publicPeerNetworker) ipify() (net.IP, error) {
	resp, err := http.Get("https://api.ipify.org")
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	bytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return net.ParseIP(string(bytes)), nil
}

func (n publicPeerNetworker) myExternalIP() (net.IP, error) {
	resp, err := http.Get("http://www.myexternalip.com/raw")
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	bytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return net.ParseIP(string(bytes)), nil
}

//CheckInternetConfig check internet configuration on the node
func (n publicPeerNetworker) CheckInternetState() error {
	_, err := net.LookupIP(cdns)
	if err != nil {
		return err
	}
	return nil
}

//CheckNtp check time synchonization on the node
func (n publicPeerNetworker) CheckNtpState() error {
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

//NewPublicNetworker creates a new instance of the public implementation of the PeerNetworker interface
func NewPublicNetworker() monitoring.PeerNetworker {
	return publicPeerNetworker{}
}
