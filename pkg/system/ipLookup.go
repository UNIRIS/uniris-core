package system

import (
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
)

func ipifyIP() (net.IP, error) {
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

func myExternalIP() (net.IP, error) {
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

func privateIP(localIface string) (net.IP, error) {
	ifaces, err := net.Interfaces()
	if err != nil {
		return nil, err
	}
	var iface net.Interface
	for _, i := range ifaces {
		if i.Name == localIface {
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
	return nil, fmt.Errorf("Cannot find a IP address from the interface %s", localIface)
}
