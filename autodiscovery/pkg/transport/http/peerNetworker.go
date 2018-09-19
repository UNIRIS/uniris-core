package http

import (
	"errors"
	"io/ioutil"
	"net"
	"net/http"
)

//ErrFailToGetIP is returned when the service to get IP does not respond
var ErrFailToGetIP = errors.New("Cannot get the peer IP. IP providers may failed")

//PeerNetworker implements the PeerNetworker interface which provides the methods to get network peer's details
type PeerNetworker struct{}

//IP lookups the peer's IP
func (pn PeerNetworker) IP() (net.IP, error) {
	var ip net.IP
	ip, err := pn.ipify()
	if err != nil {
		ip, err := pn.myExternalIP()
		if err != nil {
			return nil, ErrFailToGetIP
		}
		return ip, nil
	}
	return ip, nil
}

func (pn PeerNetworker) ipify() (net.IP, error) {
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

func (pn PeerNetworker) myExternalIP() (net.IP, error) {
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
