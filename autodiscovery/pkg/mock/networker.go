package mock

import (
	"errors"
	"net"
)

type Networker struct{}

func (n Networker) IP() (net.IP, error) {
	return net.ParseIP("127.0.0.1"), nil
}

func (n Networker) CheckInternetState() error {
	return nil
}

func (n Networker) CheckNtpState() error {
	return nil
}

type mockRobotWatcher struct{}

func (r mockRobotWatcher) CheckAutodiscoveryProcess(port int) error {
	return nil
}

func (r mockRobotWatcher) CheckDataProcess() error {
	return nil
}

func (r mockRobotWatcher) CheckMiningProcess() error {
	return nil
}

func (r mockRobotWatcher) CheckAIProcess() error {
	return nil
}

func (r mockRobotWatcher) CheckScyllaDbProcess() error {
	return nil
}

func (r mockRobotWatcher) CheckRedisProcess() error {
	return nil
}

func (r mockRobotWatcher) CheckRabbitmqProcess() error {
	return nil
}

type NetworkerNTPFails struct{}

func (n NetworkerNTPFails) IP() (net.IP, error) {
	return net.ParseIP("127.0.0.1"), nil
}

func (n NetworkerNTPFails) CheckInternetState() error {
	return nil
}

func (n NetworkerNTPFails) CheckNtpState() error {
	return errors.New("System Clock have a big Offset check the ntp configuration of the system")
}

type NetworkerInternetFails struct{}

func (n NetworkerInternetFails) IP() (net.IP, error) {
	return net.ParseIP("127.0.0.1"), nil
}

func (n NetworkerInternetFails) CheckInternetState() error {
	return errors.New("required processes are not running")
}

func (n NetworkerInternetFails) CheckNtpState() error {
	return nil
}
