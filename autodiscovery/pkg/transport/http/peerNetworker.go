package http

import (
	"errors"
	"io/ioutil"
	"net"
	"net/http"

	"github.com/uniris/uniris-core/autodiscovery/pkg/bootstraping"
)

//ErrFailToGetIP is returned when the service to get IP does not respond
var ErrFailToGetIP = errors.New("Cannot get the peer IP. IP providers may failed")

type httpNetworker struct{}

//IP lookups the peer's IP
func (n httpNetworker) IP() (net.IP, error) {
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

func (n httpNetworker) ipify() (net.IP, error) {
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

func (n httpNetworker) myExternalIP() (net.IP, error) {
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

//NewPeerNetworker creates a new instance of the http implementation of the PeerNetworker interface
func NewPeerNetworker() bootstraping.PeerNetworker {
	return httpNetworker{}
}
