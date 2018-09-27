package mock

import (
	"errors"
	"net"
)

type RobotWatcher struct{}

func (r RobotWatcher) CheckAutodiscoveryProcess(port int) error {
	return nil
}

func (r RobotWatcher) CheckDataProcess() error {
	return nil
}

func (r RobotWatcher) CheckMiningProcess() error {
	return nil
}

func (r RobotWatcher) CheckAIProcess() error {
	return nil
}

func (r RobotWatcher) CheckScyllaDbProcess() error {
	return nil
}

func (r RobotWatcher) CheckRedisProcess() error {
	return nil
}

func (r RobotWatcher) CheckRabbitmqProcess() error {
	return nil
}

type mockSystemNetworker2 struct{}

func (n mockSystemNetworker2) IP() (net.IP, error) {
	return net.ParseIP("127.0.0.1"), nil
}

func (n mockSystemNetworker2) CheckInternetState() error {
	return nil
}

func (n mockSystemNetworker2) CheckNtpState() error {
	return errors.New("System Clock have a big Offset check the ntp configuration of the system")
}

type mockSystemNetworker3 struct{}

func (n mockSystemNetworker3) IP() (net.IP, error) {
	return net.ParseIP("127.0.0.1"), nil
}

func (n mockSystemNetworker3) CheckInternetState() error {
	return errors.New("required processes are not running")
}

func (n mockSystemNetworker3) CheckNtpState() error {
	return nil
}
