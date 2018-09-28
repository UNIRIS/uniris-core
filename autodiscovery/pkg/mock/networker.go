package mock

import (
	"errors"
	"net"
)

//Networker mock
type Networker struct{}

//IP mock
func (n Networker) IP() (net.IP, error) {
	return net.ParseIP("127.0.0.1"), nil
}

//CheckInternetState mock
func (n Networker) CheckInternetState() error {
	return nil
}

//CheckNtpState mock
func (n Networker) CheckNtpState() error {
	return nil
}

//NetworkerNTPFails a Networker mock
type NetworkerNTPFails struct{}

//IP mock
func (n NetworkerNTPFails) IP() (net.IP, error) {
	return net.ParseIP("127.0.0.1"), nil
}

//CheckInternetState mock
func (n NetworkerNTPFails) CheckInternetState() error {
	return nil
}

//CheckNtpState mock
func (n NetworkerNTPFails) CheckNtpState() error {
	return errors.New("System Clock have a big Offset check the ntp configuration of the system")
}

//NetworkerInternetFails mock
type NetworkerInternetFails struct{}

//IP mock
func (n NetworkerInternetFails) IP() (net.IP, error) {
	return net.ParseIP("127.0.0.1"), nil
}

//CheckInternetState mock
func (n NetworkerInternetFails) CheckInternetState() error {
	return errors.New("required processes are not running")
}

//CheckNtpState mock
func (n NetworkerInternetFails) CheckNtpState() error {
	return nil
}
